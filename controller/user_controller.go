package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"project/models"
	"regexp"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	DB *sql.DB
}

func NewUserController(db *sql.DB) *UserController {
	return &UserController{DB: db}
}

func (uc *UserController) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized - no session found", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid cookie", http.StatusBadRequest)
			return
		}

		sessionToken := cookie.Value
		secretKey := []byte(os.Getenv("Secret_Key"))
		token, err := jwt.Parse(sessionToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (uc *UserController) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "logout successfully"}`))
}

func (uc *UserController) AddUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/registerUser.html"))
		tmpl.Execute(w, nil)
		return
	}

	//  Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	// Create user object from form values
	user.Name = strings.TrimSpace(r.FormValue("name"))
	user.Email = strings.TrimSpace(r.FormValue("email"))
	user.Username = strings.TrimSpace(r.FormValue("username"))
	user.Password = strings.TrimSpace(r.FormValue("password"))
	user.Phone = strings.TrimSpace(r.FormValue("phone"))

	// validation
	if user.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if user.Email == "" || !emailRegex.MatchString(user.Email) {
		http.Error(w, "Invalid or missing email", http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	if len(user.Password) < 6 {
		http.Error(w, "Password must be at least 6 characters long", http.StatusBadRequest)
		return
	}

	phoneRegex := regexp.MustCompile(`^[0-9]{10}$`)
	if !phoneRegex.MatchString(user.Phone) {
		http.Error(w, "Invalid phone number", http.StatusBadRequest)
		return
	}

	//  Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Insert into DB
	res, err := uc.DB.Exec("INSERT INTO users (name, email, username, password, phone) VALUES (?, ?, ?, ?, ?)",
		user.Name, user.Email, user.Username, user.Password, user.Phone)
	if err != nil {
		http.Error(w, "Insert failed", http.StatusInternalServerError)
		return
	}
	// Get last inserted ID
	id, err := res.LastInsertId()
	if err == nil {
		user.Id = int(id)
	}
	fmt.Println("Inserted user ID:", user.Id)
	// Redirect
	http.Redirect(w, r, "/get_all_users", http.StatusSeeOther)

}

func (uc *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := uc.DB.Query("SELECT * FROM users")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Username, &user.Password, &user.Phone); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		users = append(users, user)
	}
	tmpl := template.Must(template.ParseFiles("templates/navbar.html", "templates/getUserList.html"))
	tmpl.Execute(w, users)
}

func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var s models.User
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	_, err = uc.DB.Exec("UPDATE users SET name=?, email=?, username=?, password=?, phone=? WHERE id=?",
		s.Name, s.Email, s.Username, s.Password, s.Phone, id)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "record updated"}`))
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	_, err := uc.DB.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		http.Error(w, "Could not delete user", http.StatusInternalServerError)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// w.Write([]byte(`{"message": "record deleted"}`))
	http.Redirect(w, r, "/get_all_users", http.StatusSeeOther)
}

func (uc *UserController) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	res, err := uc.DB.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)",
		post.User_Id, post.Title, post.Content)
	if err != nil {
		fmt.Println("error :", err)
		http.Error(w, "Insert failed", http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err == nil {
		post.Id = int(id)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (uc *UserController) GetAllPost(w http.ResponseWriter, r *http.Request) {
	posts := []models.Post{}
	post := models.Post{}
	rows, err := uc.DB.Query("select * from posts")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&post.Id, &post.User_Id, &post.Title, &post.Content)
		posts = append(posts, post)
	}
	// json.NewEncoder(w).Encode(posts)
	tmpl := template.Must(template.ParseFiles("templates/navbar.html", "templates/getPostList.html"))
	tmpl.Execute(w, posts)

}

func (uc *UserController) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	_, err = uc.DB.Exec("UPDATE posts SET user_id = ?, title = ?, content = ? WHERE id = ?",
		post.User_Id, post.Title, post.Content, id)
	if err != nil {
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Post Updated"}`))
}
func (uc *UserController) DeletePost(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := uc.DB.Exec("DELETE FROM posts WHERE id=?", id)
	if err != nil {
		http.Error(w, "Could not delete post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Post Deleted"}`))
}
func (uc *UserController) GetPostById(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	err = uc.DB.QueryRow("SELECT id, user_id, title, content FROM posts WHERE id = ?", id).
		Scan(&post.Id, &post.User_Id, &post.Title, &post.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)

}

func (uc *UserController) GetHtmlData(w http.ResponseWriter, r *http.Request) {
	// tmpl := template.Must(template.ParseFiles("templates/navbar.html"))
	tmpl := template.Must(template.ParseFiles("templates/navbar.html", "templates/index.html"))

	tmpl.Execute(w, nil)

}

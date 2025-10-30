package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"project/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

func (uc *UserController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, nil)
		return
	}
	details := models.ContactDetails{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	_ = details

	var user models.User

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	user.Email = r.FormValue("email")
	user.Password = r.FormValue("password")

	var storedUser models.User
	err := uc.DB.QueryRow("SELECT email, password FROM users WHERE email = ?", user.Email).
		Scan(&storedUser.Email, &storedUser.Password)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}
	// Create session
	session, _ := store.Get(r, "session")
	session.Values["authenticated"] = true
	session.Values["userEmail"] = storedUser.Email
	session.Values["userId"] = storedUser.Id

	// Save session
	session.Save(r, w)

	//  token
	token, err := createToken(user.Email)
	if err != nil {
		fmt.Println("error found", err)
		http.Error(w, "failed to create token", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "user_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})

	// Redirect
	http.Redirect(w, r, "/", http.StatusSeeOther)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": true,
		"message": "Login Successfully",
		"token":   token,
	}
	json.NewEncoder(w).Encode(response)
}

func createToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"useremail": email,
			"exp":       time.Now().Add(time.Hour * 24).Unix(),
		})
	secretKey := []byte(os.Getenv("Secret_Key"))

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

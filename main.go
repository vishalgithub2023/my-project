package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	routerhttp "project/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

// type User struct {
// 	Id       int    `json:"id"`
// 	Name     string `json:"name"`
// 	Email    string `json:"email"`
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// 	Phone    string `json:"phone"`
// }

func getMySqlDB() {
	dsn := "root:12345@tcp(127.0.0.1:3307)/admin?parseTime=true"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	fmt.Println("Successfully connected to MySQL!")

}

// func getUsers(w http.ResponseWriter, r *http.Request) {
// 	ss := []User{}
// 	s := User{}
// 	rows, err := db.Query("select * from users")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		rows.Scan(&s.Id, &s.Name, &s.Email, &s.Username, &s.Password, &s.Phone)
// 		ss = append(ss, s)
// 	}
// 	json.NewEncoder(w).Encode(ss)

// }

// func addUser(w http.ResponseWriter, r *http.Request) {
// 	var s User
// 	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
// 		http.Error(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	res, err := db.Exec("INSERT INTO users (name, email, username, password, phone) VALUES (?, ?, ?, ?, ?)",
// 		s.Name, s.Email, s.Username, s.Password, s.Phone)
// 	if err != nil {
// 		http.Error(w, "Insert failed", http.StatusInternalServerError)
// 		return
// 	}

// 	id, err := res.LastInsertId()
// 	if err == nil {
// 		s.Id = int(id)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(s)
// }

// func updateUser(w http.ResponseWriter, r *http.Request) {
// 	var s User
// 	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
// 		http.Error(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	idStr := mux.Vars(r)["id"]
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid user ID", http.StatusBadRequest)
// 		return
// 	}

// 	_, err = db.Exec("UPDATE users SET name=?, email=?, username=?, password=?, phone=? WHERE id=?",
// 		s.Name, s.Email, s.Username, s.Password, s.Phone, id)
// 	if err != nil {
// 		http.Error(w, "Failed to update user", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(s)
// }

// func deleteUser(w http.ResponseWriter, r *http.Request) {
// 	id, _ := strconv.Atoi(mux.Vars(r)["id"])

// 	_, err := db.Exec("DELETE FROM users WHERE id=?", id)
// 	if err != nil {
// 		http.Error(w, "Could not delete user", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write([]byte(`{"message": "record deleted"}`))
// }

func main() {
	getMySqlDB()
	defer db.Close()
	// _ = controller.NewUserController(db)

	r := mux.NewRouter()
	routerhttp.RegisterUserRoutes(r, db)
	fmt.Println("hello api")
	// r.HandleFunc("/users", getUsers).Methods("GET")
	// r.HandleFunc("/users", addUser).Methods("POST")
	// r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	// r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	http.ListenAndServe(":8080", r)

}

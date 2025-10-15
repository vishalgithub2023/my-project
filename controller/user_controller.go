package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"project/models"
	"strconv"

	"github.com/gorilla/mux"
)

type UserController struct {
	DB *sql.DB
}

func NewUserController(db *sql.DB) *UserController {
	return &UserController{DB: db}
}

func (uc *UserController) AddUser(w http.ResponseWriter, r *http.Request) {
	var s models.User
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	res, err := uc.DB.Exec("INSERT INTO users (name, email, username, password, phone) VALUES (?, ?, ?, ?, ?)",
		s.Name, s.Email, s.Username, s.Password, s.Phone)
	if err != nil {
		http.Error(w, "Insert failed", http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err == nil {
		s.Id = int(id)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func (uc *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("check new branch")
	fmt.Println("Get user called")
	ss := []models.User{}
	s := models.User{}
	rows, err := uc.DB.Query("select * from users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&s.Id, &s.Name, &s.Email, &s.Username, &s.Password, &s.Phone)
		ss = append(ss, s)
	}
	json.NewEncoder(w).Encode(ss)

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
	json.NewEncoder(w).Encode(s)
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := uc.DB.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		http.Error(w, "Could not delete user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "record deleted"}`))
}

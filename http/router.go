package http

import (
	"database/sql"
	"project/controller"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(r *mux.Router, db *sql.DB) {
	userController := controller.NewUserController(db)
	r.HandleFunc("/", userController.GetUsers).Methods("GET")
	r.HandleFunc("/users", userController.AddUser).Methods("POST")
	r.HandleFunc("/users/{id}", userController.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", userController.DeleteUser).Methods("DELETE")

}

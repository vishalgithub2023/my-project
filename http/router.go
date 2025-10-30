package http

import (
	"database/sql"
	"project/controller"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(r *mux.Router, db *sql.DB) {

	userController := controller.NewUserController(db)

	// Public Routes
	r.HandleFunc("/", userController.GetHtmlData).Methods("GET")
	r.HandleFunc("/login", userController.LoginHandler).Methods("GET", "POST")
	r.HandleFunc("/logout", userController.LogoutHandler).Methods("POST")

	r.HandleFunc("/add_user", userController.AddUser).Methods("POST", "GET")
	r.HandleFunc("/get_all_users", userController.GetUsers).Methods("GET")

	// User Routes
	r.HandleFunc("/update/user/{id}", userController.AuthMiddleware(userController.UpdateUser)).Methods("PUT")
	r.HandleFunc("/delete/user/{id}", userController.AuthMiddleware(userController.DeleteUser)).Methods("POST")
	// Post Route
	r.HandleFunc("/create_post", userController.AuthMiddleware(userController.CreatePost)).Methods("POST")
	r.HandleFunc("/get_all_post", userController.AuthMiddleware(userController.GetAllPost)).Methods("GET")
	r.HandleFunc("/update_post/{id}", userController.AuthMiddleware(userController.UpdatePost)).Methods("PUT")
	r.HandleFunc("/delete_post/{id}", userController.AuthMiddleware(userController.DeletePost)).Methods("DELETE")
	r.HandleFunc("/get_post_by_id/{id}", userController.AuthMiddleware(userController.GetPostById)).Methods("GET")
}

// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	routerhttp "project/api"

// 	_ "github.com/go-sql-driver/mysql"
// 	"github.com/gorilla/mux"
// )

// var db *sql.DB

// func getMySqlDB() {
// 	dsn := "root:12345@tcp(127.0.0.1:3307)/admin?parseTime=true"
// 	var err error
// 	db, err = sql.Open("mysql", dsn)
// 	if err != nil {
// 		log.Fatalf("Error opening database: %v", err)
// 	}
// 	if err = db.Ping(); err != nil {
// 		log.Fatalf("Error connecting to database: %v", err)
// 	}
// 	fmt.Println("Successfully connected to MySQL!")

// }

// func main() {
// 	getMySqlDB()
// 	defer db.Close()
// 	r := mux.NewRouter()
// 	routerhttp.RegisterUserRoutes(r, db)
// 	// fmt.Println("hello api")
// 	http.ListenAndServe(":8080", r)

// }

package main

import (
	"log"
	"net/http"
	"project/api"
	"project/internal/database"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	db := database.ConnectMySQL()
	defer db.Close()
	router := mux.NewRouter()
	// router := api.RegisterRoutes(db,_)
	api.RegisterRoutes(router, db)

	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", router)
}

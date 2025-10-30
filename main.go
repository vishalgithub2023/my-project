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

func main() {
	getMySqlDB()
	defer db.Close()
	r := mux.NewRouter()
	routerhttp.RegisterUserRoutes(r, db)
	// fmt.Println("hello api")
	http.ListenAndServe(":8080", r)

}

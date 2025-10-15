package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"project/controller"
// 	route "project/http"

// 	_ "github.com/go-sql-driver/mysql"
// )

// func ConnectToDB() {
// 	dsn := "root:12345@tcp(127.0.0.1:3307)/admin?parseTime=true"

// 	db, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		log.Fatalf("Error opening database: %v", err)
// 	}
// 	defer db.Close()

// 	if err = db.Ping(); err != nil {
// 		log.Fatalf("Error connecting to database: %v", err)
// 	}
// 	fmt.Println("Successfully connected to MySQL!")
// 	defer db.Close()

// 	controller.SetDB(db)
// 	route.RegisterUserRoutes()

// 	log.Println("Server starting on :8082")
// 	log.Fatal(http.ListenAndServe(":8082", nil))

// }

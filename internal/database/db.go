package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectMySQL() *sql.DB {
	dsn := "root:12345@tcp(127.0.0.1:3307)/admin?parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf(" Error opening DB: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf(" Error connecting DB: %v", err)
	}

	log.Println("Connected to MySQL database")
	return db
}

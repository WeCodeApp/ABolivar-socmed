package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func ConnectDB() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
        os.Getenv("DB_USER"),     
        os.Getenv("DB_PASSWORD"), 
        os.Getenv("DB_HOST"),     
        os.Getenv("DB_PORT"),     
        os.Getenv("DB_NAME"),     
    )

    DB, err = sql.Open("mysql", connStr)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    if DB == nil {
        log.Fatalf("Database connection is nil")
    }

    err = DB.Ping()
    if err != nil {
        log.Fatalf("Database connection is not alive: %v", err)
    }

    log.Println("MySQL database connection established")
}
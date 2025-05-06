package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JWTSecretKey string
var DBConnectionString string

func LoadConfig() {
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: No .env file found, using system environment variables")
    }

    JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
    if JWTSecretKey == "" {
        log.Fatal("JWT_SECRET_KEY is not set in environment variables")
	}

    DBConnectionString = os.Getenv("DB_CONNECTION_STRING")
    if DBConnectionString == "" {
        log.Fatal("DB_CONNECTION_STRING is not set in environment variables")
    }
}
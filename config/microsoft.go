package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
    MicrosoftClientID     string
    MicrosoftClientSecret string
    MicrosoftTenantID     string
    MicrosoftRedirectURL  string
)

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    MicrosoftClientID = os.Getenv("MICROSOFT_CLIENT_ID")
    MicrosoftClientSecret = os.Getenv("MICROSOFT_CLIENT_SECRET") 
    MicrosoftTenantID = os.Getenv("MICROSOFT_TENANT_ID")
    MicrosoftRedirectURL = os.Getenv("MICROSOFT_CALLBACK_URL")
}

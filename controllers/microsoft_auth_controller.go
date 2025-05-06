package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"social-media-app/config"
	database "social-media-app/db"
	"social-media-app/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var jwtKey []byte

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
}

var msConfig = &oauth2.Config{
	ClientID:     config.MicrosoftClientID,
	ClientSecret: config.MicrosoftClientSecret,
	RedirectURL:  config.MicrosoftRedirectURL,
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://login.microsoftonline.com/consumers/oauth2/v2.0/authorize",
		TokenURL: "https://login.microsoftonline.com/consumers/oauth2/v2.0/token",
	},
	Scopes: []string{
		"openid",
		"email",
		"profile",
		"offline_access",
		"https://graph.microsoft.com/user.read",
	},
}

type MicrosoftUser struct {
	DisplayName       string `json:"displayName"`
	Email             string `json:"mail"`
	UserPrincipalName string `json:"userPrincipalName"` 
	ID                string `json:"id"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func InitiateMicrosoftAuth(w http.ResponseWriter, r *http.Request) {
	url := msConfig.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleMicrosoftCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state != "state" { 
		log.Println("Invalid state parameter")
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	if errMsg := r.URL.Query().Get("error"); errMsg != "" {
		errDesc := r.URL.Query().Get("error_description")
		log.Printf("Authentication error: %s - %s\n", errMsg, errDesc)
		http.Error(w, "Authentication error: "+errDesc, http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		log.Println("Authorization code is missing")
		http.Error(w, "Authorization code is missing", http.StatusBadRequest)
		return
	}

	token, err := msConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Printf("Microsoft token exchange error: %v\n", err)
		http.Error(w, "Authentication failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := msConfig.Client(r.Context(), token)
	resp, err := client.Get("https://graph.microsoft.com/v1.0/me")
	if err != nil {
		log.Printf("Failed to get user info: %v\n", err)
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Microsoft Graph API error: Status %d, Body: %s\n", resp.StatusCode, string(body))
		http.Error(w, fmt.Sprintf("Microsoft API error: %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	var msUser MicrosoftUser
	if err := json.NewDecoder(resp.Body).Decode(&msUser); err != nil {
		log.Println("Failed to decode user info")
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	userEmail := msUser.Email
	if userEmail == "" {
		userEmail = msUser.UserPrincipalName
	}

	user := models.User{
		Username:     msUser.DisplayName,
		Email:        userEmail,
		PasswordHash: "", 
		CreatedAt:    time.Now(),
	}

	var existingUser models.User
	err = database.DB.QueryRow("SELECT id FROM users WHERE email = ?", user.Email).Scan(&existingUser.ID)
	if err != nil {
		_, err = database.DB.Exec(
			"INSERT INTO users (username, email, password_hash, created_at) VALUES (?, ?, ?, ?)",
			user.Username, user.Email, user.PasswordHash, user.CreatedAt,
		)
		if err != nil {
			log.Printf("Failed to create user: %v\n", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(jwtKey)
	if err != nil {
		log.Printf("Failed to generate token: %v\n", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
		Secure:   r.TLS != nil, 
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	http.Redirect(w, r, "http://localhost:5173", http.StatusTemporaryRedirect)
}

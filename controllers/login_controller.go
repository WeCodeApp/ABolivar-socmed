package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"social-media-app/config"
	database "social-media-app/db"
	"social-media-app/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(config.JWTSecretKey) 

type Credentials struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type Claims struct {
    Email string `json:"email"`
    jwt.RegisteredClaims
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
    var creds Credentials
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    var user models.User
    var createdAt string
    err = database.DB.QueryRow("SELECT id, username, email, password_hash, created_at FROM users WHERE email = ?", creds.Email).
        Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &createdAt)
    if err == sql.ErrNoRows {
        log.Printf("Login failed: User with email %s not found", creds.Email)
        http.Error(w, "User not found", http.StatusUnauthorized)
        return
    } else if err != nil {
        log.Printf("Login failed: Error querying database for email %s: %v", creds.Email, err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
    if err != nil {
        log.Printf("Login failed: Error parsing created_at for email %s: %v", creds.Email, err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password))
    if err != nil {
        log.Printf("Login failed: Invalid credentials for email %s", creds.Email)
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        Email: creds.Email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name:    "token",
        Value:   tokenString,
        Expires: expirationTime,
    })

    log.Printf("Login successful: User with email %s logged in", creds.Email)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Login successful"))
}
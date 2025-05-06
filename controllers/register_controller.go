package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"social-media-app/models"
	"time"

	database "social-media-app/db"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	user.PasswordHash = string(hashedPassword)
	user.CreatedAt = time.Now()

	_, err = database.DB.Exec(
		"INSERT INTO users (username, email, password_hash, created_at) VALUES (?, ?, ?, ?)",
		user.Username, user.Email, user.PasswordHash, user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User already exists", http.StatusConflict)
		} else {
			log.Printf("Error inserting user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

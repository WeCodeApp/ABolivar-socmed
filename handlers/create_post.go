package handlers

import (
	"encoding/json"
	"log" 
	"net/http"
	"social-media-app/helper"
	"social-media-app/models"
	"time"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
    var post models.Post
    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    post.CreatedAt = time.Now()
    post.UpdatedAt = time.Now()

    if err := helper.CreatePost(&post); err != nil {
        log.Printf("Error creating post: %v", err) 
        http.Error(w, "Failed to create post", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(post)
}
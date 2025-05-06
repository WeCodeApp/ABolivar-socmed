package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"social-media-app/models"
	"strconv"
	"time"

	database "social-media-app/db"

	"github.com/gorilla/mux"
)

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Error parsing post ID: %v", err) 
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var existingPost models.Post
	err = database.DB.QueryRow("SELECT id, user_id, title, content, image_url, likes, created_at, updated_at FROM post WHERE id = ?", postID).Scan(
		&existingPost.ID, &existingPost.UserID, &existingPost.Title, &existingPost.Content, &existingPost.ImageURL, &existingPost.Likes, &existingPost.CreatedAt, &existingPost.UpdatedAt,
	)
	if err != nil {
		log.Printf("Error fetching post with ID %d: %v", postID, err) 
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	var updatedPost models.Post
	if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
		log.Printf("Error decoding request body: %v", err) 
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec(
		"UPDATE post SET title = ?, content = ?, image_url = ?, likes = ?, updated_at = ? WHERE id = ?",
		updatedPost.Title, updatedPost.Content, updatedPost.ImageURL, updatedPost.Likes, time.Now(), postID,
	)
	if err != nil {
		log.Printf("Error updating post with ID %d: %v", postID, err) 
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	existingPost.Title = updatedPost.Title
	existingPost.Content = updatedPost.Content
	existingPost.ImageURL = updatedPost.ImageURL
	existingPost.Likes = updatedPost.Likes
	existingPost.UpdatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingPost)
}

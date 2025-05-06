package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"social-media-app/models"
	"strconv"
	"time"

	database "social-media-app/db"

	"github.com/gorilla/mux"
)

func ViewPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Error converting post ID: %v", err)
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var post models.Post
	var imageURL sql.NullString
	err = database.DB.QueryRow(
		"SELECT id, user_id, title, content, image_url, likes, created_at FROM post WHERE id = ?",
		postID,
	).Scan(
		&post.ID, &post.UserID, &post.Title, &post.Content, &imageURL,
		&post.Likes, &post.CreatedAt,
	)
	if err == sql.ErrNoRows {
		log.Printf("Post not found: ID %d", postID)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error querying post: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	post.ImageURL = imageURL.String

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func ViewAllPosts(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(
		"SELECT id, user_id, title, content, image_url, likes, created_at FROM post",
	)
	if err != nil {
		log.Printf("Error querying all posts: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var imageURL sql.NullString
		var createdAt []uint8 
		err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Content, &imageURL,
			&post.Likes, &createdAt,
		)
		if err != nil {
			log.Printf("Error scanning post row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		post.ImageURL = imageURL.String
		post.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			log.Printf("Error parsing created_at: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

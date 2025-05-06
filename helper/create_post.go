package helper

import (
	"log" 
	"social-media-app/models"

	database "social-media-app/db"
)

func CreatePost(post *models.Post) error {
	query := `
        INSERT INTO post (user_id, title, content, image_url, likes, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `
	_, err := database.DB.Exec(query, post.UserID, post.Title, post.Content, post.ImageURL, post.Likes, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		log.Printf("Error inserting post into database: %v", err) 
	}
	return err
}
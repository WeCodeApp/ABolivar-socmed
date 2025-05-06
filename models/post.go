package models

import (
	"time"
)

type Post struct {
    ID        uint        `json:"id" gorm:"primaryKey"`
    UserID    uint        `json:"user_id"`
    Title     string      `json:"title"`
    Content   string      `json:"content"`
    ImageURL  string      `json:"image_url"`
    Likes     int         `json:"likes"`
    CreatedAt time.Time   `json:"created_at"`
    UpdatedAt time.Time   `json:"updated_at"`
}
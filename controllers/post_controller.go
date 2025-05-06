package controllers

import (
	"net/http"

	"social-media-app/handlers"
)

func ViewPost(w http.ResponseWriter, r *http.Request) {
	handlers.ViewPost(w, r)
}

func ViewAllPosts(w http.ResponseWriter, r *http.Request) {
	handlers.ViewAllPosts(w, r)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	handlers.UpdatePost(w, r)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	handlers.DeletePost(w, r)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	handlers.CreatePost(w, r)
}

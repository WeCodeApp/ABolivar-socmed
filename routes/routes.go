package routes

import (
	"social-media-app/controllers"

	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(router *mux.Router) {
    authRouter := router.PathPrefix("/auth").Subrouter()
    authRouter.HandleFunc("/login", controllers.LoginUser).Methods("GET")
    authRouter.HandleFunc("/register", controllers.RegisterUser).Methods("GET")
    authRouter.HandleFunc("/microsoft", controllers.InitiateMicrosoftAuth).Methods("GET")
    authRouter.HandleFunc("/microsoft/callback", controllers.HandleMicrosoftCallback).Methods("GET")
}

func PostRoutes(router *mux.Router) {
    postRouter := router.PathPrefix("/posts").Subrouter()
    //view all posts
    postRouter.HandleFunc("/", controllers.ViewAllPosts).Methods("GET") 
    //view a single post
    postRouter.HandleFunc("/{id}", controllers.ViewPost).Methods("GET")
    //add a new post
    postRouter.HandleFunc("/", controllers.CreatePost).Methods("POST")
    //update a post
    postRouter.HandleFunc("/{id}", controllers.UpdatePost).Methods("PUT")
    //delete a post
    postRouter.HandleFunc("/{id}", controllers.DeletePost).Methods("DELETE")
}
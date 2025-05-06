package main

import (
	"log"
	"net/http"
	"social-media-app/config"
	database "social-media-app/db"
	"social-media-app/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	config.LoadConfig()

	database.ConnectDB()

	router := mux.NewRouter()

	routes.RegisterAuthRoutes(router)
	routes.PostRoutes(router)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
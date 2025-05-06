package controllers

import (
	"log"
	"net/http"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
    log.Println("RegisterUser called")
    InitiateMicrosoftAuth(w, r)
}

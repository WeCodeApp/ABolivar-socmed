package controllers

import (
	"log"
	"net/http"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
    log.Println("LoginUser called")
    InitiateMicrosoftAuth(w, r)
}
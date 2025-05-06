package middlewares

import (
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Unauthorized: Missing Authorization header", http.StatusUnauthorized)
            return
        }

    
        next.ServeHTTP(w, r)
    })
}
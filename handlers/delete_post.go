package handlers

import (
	"encoding/json"
	"net/http"
	database "social-media-app/db"
	"strconv"

	"github.com/gorilla/mux"
)

func DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM post WHERE id = ?"
	result, err := database.DB.Exec(query, postID)
	if err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":       "Post deleted successfully",
		"rows_affected": rowsAffected,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

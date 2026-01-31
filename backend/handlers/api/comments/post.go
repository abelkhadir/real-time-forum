package comments

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "real/backend/database"
)

type CommentReq struct {
	PostID  int    `json:"post_id"`
	Content string `json:"content"`
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	var req CommentReq

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON payload"})
		return
	}

	if req.Content == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Comment content cannot be empty"})
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthenticated"})
		return
	}

	username, err := db.GetUserBySession(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthenticated"})
		return
	}

	// Get user ID
	userID, err := db.GetUserIDByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get user ID"})
		return
	}

	_, err = db.InsertComment(req.PostID, userID, username, req.Content)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "Database insertion failed",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
	})
}

package posts

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "real/backend/database"
)

type PostReq struct {
	PostTitle      string   `json:"title"`
	PostContent    string   `json:"content"`
	PostCategories []string `json:"categories"`
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	// choose valid category, title + content non empty
	var req PostReq

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON payload"})
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

	err = db.InsertPost(username, req.PostTitle, req.PostContent, req.PostCategories)
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

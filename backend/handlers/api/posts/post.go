package posts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	db "real/backend/database"
	ws "real/backend/handlers/api/websocket"
)

type PostReq struct {
	PostTitle      string   `json:"title"`
	PostContent    string   `json:"content"`
	PostCategories []string `json:"categories"`
}

var allowedCategories = map[string]struct{}{
	"news":   {},
	"tech":   {},
	"sports": {},
	"gaming": {},
}

const (
	maxPostTitleLength   = 120
	maxPostContentLength = 1000
)

// CreatePost stores a new post for the authenticated user.
func CreatePost(w http.ResponseWriter, r *http.Request) {
	// choose valid category, title + content non empty
	var req PostReq

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON payload"})
		return
	}
	req.PostTitle = strings.TrimSpace(req.PostTitle)
	req.PostContent = strings.TrimSpace(req.PostContent)
	if req.PostTitle == "" || req.PostContent == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Title and content are required"})
		return
	}
	if len(req.PostTitle) > maxPostTitleLength {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Max title lenght is 120"})
		return
	}
	if len(req.PostContent) > maxPostContentLength {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Max content lenght is 1000"})
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

	if len(req.PostCategories) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "At least one category is required"})
		return
	}
	for _, c := range req.PostCategories {
		if _, ok := allowedCategories[c]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid category"})
			return
		}
	}

	postID, err := db.InsertPost(username, req.PostTitle, req.PostContent, req.PostCategories)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"error":   "Database insertion failed",
		})
		return
	}

	post, err := db.GetPost(int(postID))
	if err == nil {
		ws.BroadcastPost(post)
	}

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
	})
}

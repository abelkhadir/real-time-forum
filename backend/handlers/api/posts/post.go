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
	var req PostReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON payload"})
		return
	}

	UserID := "4041-3030"
	Username := UserID + "name"

	err := db.InsertPost(Username, req.PostTitle, req.PostContent, req.PostCategories)

	w.Header().Set("Content-Type", "application/json")
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

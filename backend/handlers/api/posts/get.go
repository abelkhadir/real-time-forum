package posts

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	db "real/backend/database"
	"strconv"
)

// Return number of posts determined by query, as json
func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	posts, err := db.GetPosts(page, limit)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"error": "Database read failed",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"posts":   posts,
	})
}

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	post, err := db.GetPost(id)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"error": "Database read failed",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"post":    post,
	})
}

func BroadcastOnlineUsers() {
	usernames := make([]string, 0, len(clients))
	for username := range clients {
		usernames = append(usernames, username)
	}

	msg := map[string]interface{}{
		"type":  "updatecontacts",
		"users": usernames,
	}

	for _, client := range clients {
		if err := client.Conn.WriteJSON(msg); err != nil {
			log.Println("Failed to send updatecontacts to", client.Username, err)
		}
	}
}

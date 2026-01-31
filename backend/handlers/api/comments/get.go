package comments

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	db "real/backend/database"
)

func GetComments(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil || postID < 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid post_id"})
		return
	}

	comments, err := db.GetCommentsByPost(postID)

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
		"success":  true,
		"comments": comments,
	})
}

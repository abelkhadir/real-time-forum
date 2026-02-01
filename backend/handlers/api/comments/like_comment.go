package comments

import (
	"encoding/json"
	"net/http"

	db "real/backend/database"
)

func LikeComment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CommentID int  `json:"comment_id"`
		IsLike    bool `json:"is_like"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	cookie, _ := r.Cookie("session_token")
	username, _ := db.GetUserBySession(cookie.Value)
	userID, _ := db.GetUserIDByUsername(username)

	err := db.UpdateCommentLike(userID, req.CommentID, req.IsLike)

	json.NewEncoder(w).Encode(map[string]any{
		"success": err == nil,
	})
}

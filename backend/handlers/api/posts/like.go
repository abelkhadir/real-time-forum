package posts

import (
	"encoding/json"
	"net/http"

	db "real/backend/database"
)

type LikeReq struct {
	PostID int  `json:"post_id"`
	IsLike bool `json:"is_like"`
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	var req LikeReq
	json.NewDecoder(r.Body).Decode(&req)

	cookie, _ := r.Cookie("session_token")
	username, _ := db.GetUserBySession(cookie.Value)
	userID, _ := db.GetUserIDByUsername(username)

	err := db.UpdatePostLike(userID, req.PostID, req.IsLike)

	json.NewEncoder(w).Encode(map[string]any{
		"success": err == nil,
	})
}

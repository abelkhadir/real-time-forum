package notifications

import (
	"encoding/json"
	"net/http"

	db "real/backend/database"
)

func MarkRead(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username, err := db.GetUserBySession(cookie.Value)
	if err != nil || username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := db.MarkNotificationsRead(username); err != nil {
		http.Error(w, "Could not update notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

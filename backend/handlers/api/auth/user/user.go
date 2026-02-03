package user

import (
	"encoding/json"
	"net/http"

	db "real/backend/database"
)

func GetContactsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	username := ""
	cookie, err := r.Cookie("session_token")
	if err == nil {
		username, _ = db.GetUserBySession(cookie.Value)
	}

	contacts, err := db.GetContacts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Could not retrieve user information"})
		return
	}
	email, _ := db.GetEmailBySession(username)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"success":  true,
		"username": username,
		"email":    email,
		"contacts": contacts,
	})
}

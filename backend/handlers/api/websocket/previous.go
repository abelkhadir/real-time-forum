package ws

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "real/backend/database"
)

func PreviousMessagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

	messages, err := db.ReadMessages(username, r.URL.Query().Get("id"))
	fmt.Println(messages)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Could not retrieve messages"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"success":  true,
		"messages": messages,
	})
}

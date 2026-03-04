package home

import (
	"net/http"
	"time"

	db "real/backend/database"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.ServeFile(w, r, "./frontend/login.html")
			return
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if _, err := db.GetUserBySession(cookie.Value); err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   "",
			Expires: time.Unix(0, 0),
			Path:    "/",
			MaxAge:  -1,
		})
		http.ServeFile(w, r, "./frontend/login.html")
		return
	}

	http.ServeFile(w, r, "./frontend/index.html")
}

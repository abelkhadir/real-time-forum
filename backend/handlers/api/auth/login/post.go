package login

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	db "real/backend/database"
	ws "real/backend/handlers/api/websocket"

	"github.com/gofrs/uuid"
)

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON payload"})
		return
	}

	req.Identifier = strings.TrimSpace(req.Identifier)
	if req.Identifier == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Identifier and password are required"})
		return
	}

	username := req.Identifier
	if strings.Contains(req.Identifier, "@") {
		if err := db.CheckCreds_email(req.Identifier, req.Password); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
			return
		}

		var err error
		username, err = db.GetUserByEmail(req.Identifier)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
			return
		}
	} else {
		if err := db.CheckCreds_user(req.Identifier, req.Password); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
			return
		}
	}

	// 3. Generer Session Token (UUID)
	sessionToken, err := uuid.NewV4()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	tokenString := sessionToken.String()
	expiresAt := time.Now().Add(24 * time.Hour)
	if err := db.InsertSession(username, tokenString, expiresAt); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Could not create session"})
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   tokenString,
		Expires: expiresAt,
		Path:    "/",
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"success":  true,
		"username": username,
	})
}

// ==========================
// LOGOUT HANDLER
// ==========================
func Logout(w http.ResponseWriter, r *http.Request) {
	// 1. Jbed cookie

	cookie, err := r.Cookie("session_token")
	if err != nil {
		// Ila makanch cookie, aslan howa logout
		w.WriteHeader(http.StatusOK)
		return
	}

	username, _ := db.GetUserBySession(cookie.Value)

	db.DeleteSess(cookie.Value)

	// 2. Mse7 session mn Database
	if username != "" {
		_ = db.RemoveOnline(username)
		ws.BroadcastContacts(username)
	}

	// 3. 9tel l-Cookie f Browser (Set expired date)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Unix(0, 0), // Date qdima
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1, // Delete immediately
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		_, err = db.GetUserBySession(c.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func NoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}

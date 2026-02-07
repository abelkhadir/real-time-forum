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
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON payload"})
		return
	}

	if strings.Contains(req.Identifier, "@") {
		if !db.CheckCreds_email(req.Identifier, req.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
			return
		}
	}

	if strings.Contains(req.Identifier, "@") == false {
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
	db.InsertSession(req.Identifier, tokenString, expiresAt)

	user := ""
	if strings.Contains(req.Identifier, "@") {
		user, err = db.GetUserByEmail(req.Identifier)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Could not retrieve user information"})
			return
		}
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Could not retrieve user information"})
		return
	}

	// 5. Sifet HTTP Only Cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   tokenString,
		Expires: expiresAt,
		Path:    "/",
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"success":  true,
		"username": user,
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

// CheckAuth Middleware
// Ay route ghat-dwr 3liha had fonction, maghadich tkhdem illa ila kan user m-connecté
func CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
	/*

		return func(w http.ResponseWriter, r *http.Request) {
			// 1. Jib Cookie
			c, err := r.Cookie("session_token")
			if err != nil {
				if err == http.ErrNoCookie {
					// Ma3ndoch cookie -> Error 401
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			sessionToken := c.Value

			// 2. Verifier wach had token kayn f DB w ma-mitsh (Expired)
			var userID int
			var expiry time.Time

			row := database.Db.QueryRow("SELECT user_id, expiry FROM sessions WHERE token = ?", sessionToken)
			err = row.Scan(&userID, &expiry)
			if err != nil {
				http.Error(w, "Unauthorized (Invalid Token)", http.StatusUnauthorized)
				return
			}

			// 3. Check Expiry
			if expiry.Before(time.Now()) {
				database.Db.Exec("DELETE FROM sessions WHERE token = ?", sessionToken)
				http.Error(w, "Session expired", http.StatusUnauthorized)
				return
			}

			// Hna tqder tzid userID f Context ila bghiti tsta3mlo mn b3d
			// ...

			// User Clean -> Doz l-fonction l-asliya
			next(w, r)
		}
	*/
}

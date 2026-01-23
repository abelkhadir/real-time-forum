package login

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"real/database"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// 1. Jbed l-User mn Database (b Email AW Username)
	var user struct {
		ID           int
		Username     string
		PasswordHash string
	}

	query := `SELECT id, username, password_hash FROM users WHERE username = ? OR email = ?`
	err := database.Db.QueryRow(query, req.Identifier, req.Identifier).Scan(&user.ID, &user.Username, &user.PasswordHash)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// 2. Tcheki Password (Compare Hash)
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 3. Generer Session Token (UUID)
	sessionToken, err := uuid.NewV4()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	tokenString := sessionToken.String()
	expiresAt := time.Now().Add(24 * time.Hour) // Session kat-mout mor 24 sa3a

	// 4. Sjel Session f Database
	// Kan-ms7o ay session qdima d-had l-user bach nbqaw clean (Optionnel)
	database.Db.Exec("DELETE FROM sessions WHERE username = ?", user.Username)

	_, err = database.Db.Exec("INSERT INTO sessions (id, username, expires_at) VALUES (?, ?, ?)",
		tokenString, user.Username, expiresAt)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// 5. Sifet HTTP Only Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    tokenString,
		Expires:  expiresAt,
		HttpOnly: true, // Mohim bzaf: JS mayqdrch yqrah (Security)
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	// 6. Jaweb b Success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Logged in successfully",
		"username": user.Username,
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

	// 2. Mse7 session mn Database
	database.Db.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)

	// 3. 9tel l-Cookie f Browser (Set expired date)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Unix(0, 0), // Date qdima
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1, // Delete immediately
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out"})
}

// CheckAuth Middleware
// Ay route ghat-dwr 3liha had fonction, maghadich tkhdem illa ila kan user m-connecté
func CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}

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
}

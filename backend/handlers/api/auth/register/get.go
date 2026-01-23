package register

import (
	"encoding/json"
	"net/http"
	"real/database"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// --- Structs ---
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Decode JSON body
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// 3. Simple Validation
	if len(req.Username) < 4 || len(req.Password) < 8 || !strings.Contains(req.Email, "@") {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	// 4. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// 5. Insert f Database

	stmt := `INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)`
	_, err = database.Db.Exec(stmt, req.Username, req.Email, string(hashedPassword))

	if err != nil {
		// Ila l-error fih "UNIQUE constraint", ya3ni username/email deja kayn
		if strings.Contains(err.Error(), "UNIQUE") {
			http.Error(w, "Username or Email already taken", http.StatusConflict)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// 6. Jaweb b Success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

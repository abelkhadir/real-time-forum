package register

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	db "real/backend/database"
)

// --- Structs ---
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	// TODO: test if email and passwords cant be repeated
	// TODO: sanitize input

	var req RegisterRequest
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON payload"})
		return
	}

	if err := validateInput(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// 5. Insert f Database
	err := db.InsertUser(req.Username, req.Email, req.Password)
	if err != nil {
		log.Printf("DB insertion error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "database error, "})
		return
	}

	// 6. Jaweb b Success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully", "success": "true"})
}

func validateInput(req RegisterRequest) error {
	email := req.Email
	username := req.Username
	password := req.Password

	if len(email) < 7 || len(email) > 40 {
		return fmt.Errorf("email length must be between 7 and 40 characters")
	}

	t, _ := db.DoesEmailExist(email)
	if t {
		return fmt.Errorf("this email is already used")
	}

	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if matched, err := regexp.MatchString(regex, email); err != nil {
		return fmt.Errorf("interserver server error")
	} else if !matched {
		return fmt.Errorf("invalid email, try again")
	}

	if len(username) < 4 || len(username) > 20 {
		return fmt.Errorf("username must be between 4 and 20 characters")
	}
	regexusername := `^[a-zA-Z0-9_-]+$`
	if matched, err := regexp.MatchString(regexusername, username); err != nil {
		return fmt.Errorf("interserver server error")
	} else if !matched {
		return fmt.Errorf("username can only use letters, numbers, - and _")
	}

	t2, _ := db.DoesUserExist(username)
	if t2 {
		return fmt.Errorf("this username is already used")
	}

	if len(password) < 8 || len(password) > 64 {
		return fmt.Errorf("password must be between 8 and 64 characters")
	}

	var hasLetter, hasNumber bool
	for _, c := range password {
		switch {
		case (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z'):
			hasLetter = true
		case c >= '0' && c <= '9':
			hasNumber = true
		}
	}
	if !hasLetter || !hasNumber {
		return fmt.Errorf("password must contain at least one letter and one number")
	}

	return nil
}


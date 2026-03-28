package register

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	db "real/backend/database"
)

// --- Structs ---
type RegisterRequest struct {
	Username  string `json:"username"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// Register validates input and creates a new account.
func Register(w http.ResponseWriter, r *http.Request) {

	var req RegisterRequest
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON payload"})
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Gender = strings.ToLower(strings.TrimSpace(req.Gender))
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Email = strings.TrimSpace(req.Email)

	if err := validateInput(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// 5. Insert f Database
	err := db.InsertUser(req.Username, req.Email, req.Age, req.Gender, req.FirstName, req.LastName, req.Password)
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

// validateInput checks the registration payload against basic rules.
func validateInput(req RegisterRequest) error {
	email := req.Email
	username := req.Username
	age := req.Age
	gender := req.Gender
	firstName := req.FirstName
	lastName := req.LastName
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

	if age < 1 || age > 130 {
		return fmt.Errorf("age must be between 1 and 130")
	}

	if gender != "male" && gender != "female" {
		return fmt.Errorf("gender must be male or female")
	}

	if len(firstName) < 1 || len(firstName) > 20 {
		return fmt.Errorf("first name must be between 1 and 50 characters")
	}
	if matched, err := regexp.MatchString(`^[A-Za-z]+$`, firstName); err != nil {
		return fmt.Errorf("interserver server error")
	} else if !matched {
		return fmt.Errorf("first name must contain letters only")
	}

	if len(lastName) < 1 || len(lastName) > 20 {
		return fmt.Errorf("last name must be between 1 and 50 characters")
	}
	if matched, err := regexp.MatchString(`^[A-Za-z]+$`, lastName); err != nil {
		return fmt.Errorf("interserver server error")
	} else if !matched {
		return fmt.Errorf("last name must contain letters only")
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

package db

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// InsertUser creates a new user with profile details and a hashed password.
func InsertUser(Username, Email string, Age int, Gender, FirstName, LastName, Password string) error {
	PasswordHash, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (username, email, age, gender, first_name, last_name, password_hash, is_online) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(stmt, Username, Email, Age, Gender, FirstName, LastName, string(PasswordHash), false)
	return err
}

// DoesEmailExist checks whether an email is already registered.
func DoesEmailExist(email string) (bool, error) {
	var existingEmail string
	err := db.QueryRow(`SELECT email FROM users WHERE email = ?`, email).Scan(&existingEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// DoesUserExist checks whether a username is already registered.
func DoesUserExist(username string) (bool, error) {
	var existingUsername string
	err := db.QueryRow(`SELECT username FROM users WHERE username = ?`, username).Scan(&existingUsername)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CheckCreds_user validates a username and password login attempt.
func CheckCreds_user(user, password string) error {
	var password_hash string

	query := `SELECT password_hash FROM users WHERE username = ?`
	err := db.QueryRow(query, user).Scan(&password_hash)

	if err == sql.ErrNoRows {
		return fmt.Errorf("Invalid credentials")
	} else if err != nil {
		fmt.Printf("Database error: %v", err)
		return fmt.Errorf("Database error: %v", err)
	}

	// 2. Tcheki Password (Compare Hash)
	err = bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password))
	if err != nil {
		return fmt.Errorf("Invalid credentials")
	}

	return nil
}

// CheckCreds_email validates an email and password login attempt.
func CheckCreds_email(email, password string) error {
	var passwordHash string

	query := `SELECT password_hash FROM users WHERE email = ?`
	err := db.QueryRow(query, email).Scan(&passwordHash)
	if err == sql.ErrNoRows {
		return fmt.Errorf("Invalid credentials")
	} else if err != nil {
		fmt.Printf("Database error: %v", err)
		return fmt.Errorf("Database error: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return fmt.Errorf("Invalid credentials")
	}

	return nil
}

// InsertSession replaces any existing session and stores a new one.
func InsertSession(username, tokenString string, expiresAt time.Time) error {
	_, err := db.Exec("DELETE FROM sessions WHERE username = ?", username)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO sessions (id, username, expires_at) VALUES (?, ?, ?)",
		tokenString, username, expiresAt)
	if err != nil {
		return err
	}

	return nil
}

// GetUserByEmail returns the username associated with an email.
func GetUserByEmail(email string) (string, error) {
	var user string

	query := `SELECT username FROM users WHERE email = ?`
	err := db.QueryRow(query, email).Scan(&user)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("no email found")
	} else if err != nil {
		fmt.Printf("Database error: %v", err)
		return "", fmt.Errorf("Database error: %v", err)
	}

	return user, nil
}

// GetEmailBySession returns the email associated with a username.
func GetEmailBySession(username string) (string, error) {
	var email string

	query := `SELECT email FROM users WHERE username = ?`
	err := db.QueryRow(query, username).Scan(&email)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("no email found")
	} else if err != nil {
		fmt.Printf("Database error: %v", err)
		return "", fmt.Errorf("Database error: %v", err)
	}

	return email, nil
}

// GetUserBySession returns the username for a valid session token.
func GetUserBySession(token string) (string, error) {
	var username string

	query := `SELECT username FROM sessions WHERE id = ? AND expires_at > ?`
	err := db.QueryRow(query, token, time.Now()).Scan(&username)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("no session found")
	} else if err != nil {
		fmt.Printf("Database error: %v", err)
		return "", fmt.Errorf("Database error: %v", err)
	}

	return username, nil
}

// DeleteSess removes a session by token.
func DeleteSess(cookie string) {
	db.Exec("DELETE FROM sessions WHERE id = ?", cookie)
}

type User struct {
	Username string
	Online   string
}

type Contact struct {
	Username string
	Online   bool
}

// GetContacts returns contacts ordered by recent conversation activity.
func GetContacts(currentUsername string) ([]Contact, error) {
	query := `
		SELECT u.username, u.is_online
		FROM users u
		LEFT JOIN messages m
			ON (
				(m.from_username = ? AND m.to_username = u.username)
				OR
				(m.from_username = u.username AND m.to_username = ?)
			)
		GROUP BY u.username, u.is_online
		ORDER BY
			CASE WHEN MAX(m.created_at) IS NULL THEN 1 ELSE 0 END,
			MAX(m.created_at) DESC,
			LOWER(u.username) ASC
	`

	rows, err := db.Query(query, currentUsername, currentUsername)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []Contact
	for rows.Next() {
		var c Contact
		if err := rows.Scan(&c.Username, &c.Online); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}

	return contacts, rows.Err()
}

// GetUserIDByUsername returns a user's numeric ID.
func GetUserIDByUsername(username string) (int, error) {
	var userID int

	query := `SELECT id FROM users WHERE username = ?`
	err := db.QueryRow(query, username).Scan(&userID)

	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("user not found")
	} else if err != nil {
		return 0, err
	}

	return userID, nil
}

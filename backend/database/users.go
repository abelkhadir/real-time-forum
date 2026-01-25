package db

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func InsertUser(Username, Email, Password string) error {
	fmt.Println(Username, Email, Password)
	PasswordHash, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)`
	_, err = db.Exec(stmt, Username, Email, string(PasswordHash))
	return err
}

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

func CheckCreds_email(email, password string) bool {
	return true
}

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

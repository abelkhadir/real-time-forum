package db

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return err
	}
	return db.Ping()
}

func Migrate() error {
	schema := `

	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL CHECK(length(username) BETWEEN 4 AND 24),
		email TEXT UNIQUE NOT NULL CHECK(length(email) <= 100),
		age INTEGER NOT NULL CHECK(age BETWEEN 1 AND 130),
		gender TEXT NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		is_online BOOLEAN NOT NULL DEFAULT 0,
		password_hash TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL,
		expires_at DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		CHECK(expires_at > created_at)
	);
	
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		username TEXT NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		comments_num INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS post_categories (
		post_id INTEGER,
		category_id INTEGER,
		PRIMARY KEY(post_id, category_id),
		FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY(category_id) REFERENCES categories(id)
	);

	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY,
		from_username TEXT NOT NULL,
		to_username TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS notifications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		from_username TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		read_at TIMESTAMP
	);
		CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER,
		user_id INTEGER,
		username TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);

	`
	if _, err := db.Exec(schema); err != nil {
		return err
	}

	// Keep existing local DBs compatible with the newer required user fields.
	if err := addColumnIfMissing(`ALTER TABLE users ADD COLUMN age INTEGER NOT NULL DEFAULT 18`); err != nil {
		return err
	}
	if err := addColumnIfMissing(`ALTER TABLE users ADD COLUMN gender TEXT NOT NULL DEFAULT 'other'`); err != nil {
		return err
	}
	if err := addColumnIfMissing(`ALTER TABLE users ADD COLUMN first_name TEXT NOT NULL DEFAULT ''`); err != nil {
		return err
	}
	if err := addColumnIfMissing(`ALTER TABLE users ADD COLUMN last_name TEXT NOT NULL DEFAULT ''`); err != nil {
		return err
	}

	return nil
}

func addColumnIfMissing(query string) error {
	_, err := db.Exec(query)
	if err == nil {
		return nil
	}

	if strings.Contains(strings.ToLower(err.Error()), "duplicate column name") {
		return nil
	}

	return err
}

func AddOnline(username string) error {
	_, err := db.Exec("UPDATE users SET is_online = 1 WHERE username = ?", username)
	return err
}

func RemoveOnline(username string) error {
	_, err := db.Exec("UPDATE users SET is_online = 0 WHERE username = ?", username)
	return err
}

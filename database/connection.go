package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

// Init initializes the database connection and creates all tables and defaults.
func Init() {
	openDatabase()
	if err := setupDatabase(); err != nil {
		log.Fatalf("Database setup failed: %v", err)
	}
}

// openDatabase opens the SQLite database connection.
func openDatabase() {
	var err error
	Db, err = sql.Open("sqlite3", "./backend/database/sqlite.db")
	if err != nil {
		log.Fatal("DB open error:", err)
	}

	if err = Db.Ping(); err != nil {
		log.Fatal("DB connection failed:", err)
	}
	log.Println("Database connection established successfully.")
}

// setupDatabase creates all tables and inserts defaults in one transaction.
func setupDatabase() error {
	tx, err := Db.Begin()
	if err != nil {
		return err
	}
	
	
	defer tx.Rollback()

	// Handle Panics
	defer func() {
		if p := recover(); p != nil {
			log.Fatalf("panic during setup: %v", p)
		}
	}()

	// Enable foreign keys
	if _, err := tx.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return err
	}

	// Helper function
	createTable := func(query, name string) error {
		if _, err := tx.Exec(query); err != nil {
			log.Printf("Failed to create %s table: %v", name, err)
			return err
		}
		log.Printf("Table '%s' ensured.", name)
		return nil
	}

	// === Creation des Tables ===
	
	if err := createTable(usersTable, "users"); err != nil {
		return err
	}
	if err := createTable(postsTable, "posts"); err != nil {
		return err
	}
	if err := createTable(categoriesTable, "categories"); err != nil {
		return err
	}
	if err := createTable(commentsTable, "comments"); err != nil {
		return err
	}
	if err := createTable(postCategoriesTable, "post_categories"); err != nil {
		return err
	}
	if err := createTable(likesTable, "likes"); err != nil {
		return err
	}
	if err := createTable(sessionsTable, "sessions"); err != nil {
		return err
	}
	if err := createTable(privateMessagesTable, "private_messages"); err != nil {
		return err
	}

	// Insert default categories
	for _, c := range []string{"All", "Programming", "Cybersecurity", "Gadgets & Hardware", "Web Development"} {
		if _, err := tx.Exec(`INSERT OR IGNORE INTO categories (name) VALUES (?)`, c); err != nil {
			log.Printf("Failed to insert category '%s': %v", c, err)
		}
	}

	// Create Indexes
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_likes_post_id ON likes(post_id);`,
		`CREATE INDEX IF NOT EXISTS idx_likes_user_value ON likes(username, value);`,
		`CREATE INDEX IF NOT EXISTS idx_likes_created_at ON likes(created_at);`,
		`CREATE INDEX IF NOT EXISTS idx_messages_sender ON private_messages(sender_username);`,
		`CREATE INDEX IF NOT EXISTS idx_messages_receiver ON private_messages(receiver_username);`,
		`CREATE INDEX IF NOT EXISTS idx_messages_conversation ON private_messages(sender_username, receiver_username, created_at);`,
	}

	for _, q := range indexes {
		if _, err := tx.Exec(q); err != nil {
			log.Printf("Failed to create index: %v", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Println("Database setup completed successfully.")
	return nil
}

// === TABLE DEFINITIONS ===
var (
	usersTable = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL CHECK(length(username) BETWEEN 4 AND 24),
		email TEXT UNIQUE NOT NULL CHECK(length(email) <= 100),
		password_hash TEXT NOT NULL
	);`

	postsTable = `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		title TEXT NOT NULL CHECK(trim(title) != '' AND length(title) <= 30),
		content TEXT NOT NULL CHECK(trim(content) != '' AND length(content) <= 1000),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		likes_count INTEGER DEFAULT 0,
		dislikes_count INTEGER DEFAULT 0,
		comments_count INTEGER DEFAULT 0,
		FOREIGN KEY(username) REFERENCES users(username) ON DELETE CASCADE
	);`

	categoriesTable = `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	);`

	commentsTable = `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		content TEXT NOT NULL CHECK(trim(content) != '' AND length(content) <= 300),
		likes_count INTEGER DEFAULT 0,
		dislikes_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (username) REFERENCES users(username) ON DELETE CASCADE
	);`

	postCategoriesTable = `
	CREATE TABLE IF NOT EXISTS post_categories (
		post_id INTEGER NOT NULL REFERENCES posts(id),
		category_id INTEGER NOT NULL REFERENCES categories(id),
		UNIQUE(post_id, category_id),
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
	);`

	likesTable = `
	CREATE TABLE IF NOT EXISTS likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		post_id INTEGER,
		comment_id INTEGER,
		value INTEGER NOT NULL CHECK(value IN (1, -1)),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY(comment_id) REFERENCES comments(id) ON DELETE CASCADE,
		UNIQUE(username, post_id, comment_id)
	);`

	sessionsTable = `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL,
		expires_at DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		CHECK(expires_at > created_at),
		FOREIGN KEY(username) REFERENCES users(username) ON DELETE CASCADE
	);`

	privateMessagesTable = `
	CREATE TABLE IF NOT EXISTS private_messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sender_username TEXT NOT NULL,
		receiver_username TEXT NOT NULL,
		content TEXT NOT NULL CHECK(trim(content) != '' AND length(content) <= 2000),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(sender_username) REFERENCES users(username) ON DELETE CASCADE,
		FOREIGN KEY(receiver_username) REFERENCES users(username) ON DELETE CASCADE
	);`
)
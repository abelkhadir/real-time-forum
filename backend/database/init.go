package db

import (
	"database/sql"

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
  CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    username TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    likes_num INTEGER NOT NULL DEFAULT 0,
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
  `
	_, err := db.Exec(schema)
	return err
}

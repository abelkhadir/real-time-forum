package main

import (
	"fmt"
	"net/http"

	db "real/backend/database"
	"real/backend/handlers/api/auth/login"
	"real/backend/handlers/api/auth/register"
	"real/backend/handlers/api/home"
	"real/backend/handlers/api/posts"
	"real/backend/handlers/api/ws"
)

func main() {
	err := db.InitDB()
	if err != nil {
		return
	}

	err = db.Migrate()
	if err != nil {
		return
	}

	mux := http.NewServeMux()

	// backend routes
	mux.HandleFunc("/", home.HomeHandler)
	mux.HandleFunc("/api/register", register.Register)
	mux.HandleFunc("/api/login", login.Login)
	mux.HandleFunc("/api/logout", login.Logout)

	mux.HandleFunc("GET /api/posts", login.CheckAuth(posts.GetPostsHandler))
	mux.HandleFunc("POST /api/posts/create", login.CheckAuth(posts.CreatePost))
	mux.HandleFunc("GET /api/posts/read", login.CheckAuth(posts.GetPostHandler))

	mux.HandleFunc("/ws", (ws.WebSocketsHandler))
	// frontend (HTML, CSS, JS)
	fs := http.FileServer(http.Dir("./frontend/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

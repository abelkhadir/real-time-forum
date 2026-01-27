package main

import (
	"fmt"
	"log"
	"net/http"

	db "real/backend/database"
	"real/backend/handlers/api/auth/login"
	"real/backend/handlers/api/auth/register"
	"real/backend/handlers/api/auth/user"
	"real/backend/handlers/api/home"
	ws "real/backend/handlers/api/messages"
	"real/backend/handlers/api/posts"
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
	mux.HandleFunc("POST /api/register", register.Register)
	mux.HandleFunc("POST /api/login", login.Login)
	mux.HandleFunc("/api/logout", login.Logout)

	mux.HandleFunc("GET /api/posts", posts.GetPostsHandler)
	mux.HandleFunc("POST /api/posts/create", login.CheckAuth(posts.CreatePost))
	mux.HandleFunc("GET /api/posts/read", posts.GetPostHandler)

	mux.HandleFunc("GET /api/contacts", user.GetContactsHandler)

	mux.HandleFunc("/ws", ws.WebSocketsHandler)
	mux.HandleFunc("/api/conversations/messages", ws.PreviousMessagesHandler)

	// frontend (HTML, CSS, JS)
	fs := http.FileServer(http.Dir("./frontend/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server started at http://localhost:8080")
	log.Panic(http.ListenAndServe(":8080", mux))
}

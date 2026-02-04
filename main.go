package main

import (
	"fmt"
	"log"
	"net/http"

	db "real/backend/database"
	"real/backend/handlers/api/auth/login"
	"real/backend/handlers/api/auth/register"
	"real/backend/handlers/api/auth/user"
	"real/backend/handlers/api/comments"
	"real/backend/handlers/api/home"
	"real/backend/handlers/api/posts"
	ws "real/backend/handlers/api/websocket"
)

func main() {
	err := db.InitDB()
	if err != nil {
		fmt.Println("Failed to initialize database:", err)
		return
	}

	err = db.Migrate()
	if err != nil {
		fmt.Println("Failed to migrate database:", err)
		return
	}

	mux := http.NewServeMux()

	// backend routes
	mux.HandleFunc("/", home.HomeHandler)
	mux.HandleFunc("POST /api/register", register.Register)
	mux.HandleFunc("POST /api/login", login.Login)
	mux.HandleFunc("POST /api/logout", login.Logout)

	mux.HandleFunc("GET /api/posts", posts.GetPostsHandler)
	mux.HandleFunc("POST /api/posts/create", login.CheckAuth(posts.CreatePost))
	mux.HandleFunc("GET /api/posts/read", posts.GetPostHandler)
	mux.HandleFunc("POST /api/posts/like", login.CheckAuth(posts.LikePost))

	mux.HandleFunc("GET /api/comments", comments.GetComments)
	mux.HandleFunc("POST /api/comments/create", login.CheckAuth(comments.CreateComment))
	mux.HandleFunc("POST /api/comments/like", login.CheckAuth(comments.LikeComment))

	mux.HandleFunc("GET /api/me", user.GetMeHandler)

	mux.HandleFunc("GET /ws", ws.WebSocketsHandler)
	mux.HandleFunc("GET /api/conversations/messages", ws.PreviousMessagesHandler)

	// frontend (HTML, CSS, JS)
	fs := http.FileServer(http.Dir("./frontend/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server started at http://localhost:8080")
	log.Panic(http.ListenAndServe(":8080", mux))
}

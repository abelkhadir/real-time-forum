package main

import (
	"fmt"
	"net/http"

	"real/backend/handlers"
)

func main() {
	mux := http.NewServeMux()

	// backend routes
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)

	// frontend (HTML, CSS, JS)
	fs := http.FileServer(http.Dir("./frontend/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

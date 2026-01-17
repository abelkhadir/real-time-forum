package main

import (
	"net/http"
	"fmt"

	"real/backend/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)

	fmt.Println("Server has started on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

package home

import "net/http"

// HomeHandler serves the single-page application entry point.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/index.html")
}

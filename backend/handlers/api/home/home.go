package home

import (
	"net/http"
	error_handler "real/backend/handlers/api/error"
)

// HomeHandler serves the single-page application entry point.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		error_handler.ErrorPage(w, "This page does not exist, or it may have been moved.", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, "./frontend/index.html")
}

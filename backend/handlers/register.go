package handlers

import "net/http"

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("register"))
}

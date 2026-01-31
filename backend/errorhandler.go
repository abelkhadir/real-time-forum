package backend

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func writeError(w http.ResponseWriter, status int, msg string, err error) {
	// log server-side for debugging
	if err != nil {
		log.Printf("%s: %v\n", msg, err)
	} else {
		log.Printf("%s\n", msg)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	er := errorResponse{Success: false, Message: msg}
	if err != nil {
		er.Error = err.Error()
	}
	json.NewEncoder(w).Encode(er)
}

func HandleBadRequest(w http.ResponseWriter, msg string, err error) {
	writeError(w, http.StatusBadRequest, msg, err)
}

func HandleInternalServerError(w http.ResponseWriter, msg string, err error) {
	writeError(w, http.StatusInternalServerError, msg, err)
}

func HandleNotFound(w http.ResponseWriter, msg string) {
	writeError(w, http.StatusNotFound, msg, nil)
}

func HandleMethodNotAllowed(w http.ResponseWriter) {
	writeError(w, http.StatusMethodNotAllowed, "method not allowed", nil)
}

func HandleUnauthorized(w http.ResponseWriter, msg string) {
	writeError(w, http.StatusUnauthorized, msg, nil)
}
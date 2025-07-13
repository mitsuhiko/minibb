package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type APIError struct {
	Detail string `json:"detail"`
}

func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func RespondWithError(w http.ResponseWriter, status int, err APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

func InternalServerError(w http.ResponseWriter, originalErr error) {
	log.Printf("Internal server error: %v", originalErr)
	RespondWithError(w, http.StatusInternalServerError, APIError{Detail: "encountered an unexpected internal failure on the backend server"})
}

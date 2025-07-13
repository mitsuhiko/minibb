package handlers

import (
	"minibb/internal/utils"
	"net/http"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "ok",
		Message: "MiniBB server is running",
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// TODO: Add other API handlers for boards, topics, posts

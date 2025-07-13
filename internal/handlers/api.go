package handlers

import (
	"net/http"

	"minibb/internal/utils"
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

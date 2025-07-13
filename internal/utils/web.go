package utils

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
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

type PaginationParams struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Offset  int `json:"offset"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalPages int `json:"total_pages"`
	Total      int `json:"total"`
}

func ParsePaginationParams(r *http.Request) PaginationParams {
	page := 1
	perPage := 50

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if perPageStr := r.URL.Query().Get("per_page"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 100 {
			perPage = pp
		}
	}

	offset := (page - 1) * perPage

	return PaginationParams{
		Page:    page,
		PerPage: perPage,
		Offset:  offset,
	}
}

func CalculatePaginationMeta(page, perPage, total int) PaginationMeta {
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	return PaginationMeta{
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
		Total:      total,
	}
}

func ParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

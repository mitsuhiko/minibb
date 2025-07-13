package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"minibb/internal/db"
	"minibb/internal/models"
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

type CreateTopicRequest struct {
	Title   string `json:"title"`
	Author  string `json:"author"`
	BoardID int    `json:"board_id"`
}

func CreateTopic(w http.ResponseWriter, r *http.Request) {
	db := db.FromContext(r.Context())

	var req CreateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "Invalid JSON"})
		return
	}

	if req.Title == "" || req.Author == "" || req.BoardID == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "Title, author, and board_id are required"})
		return
	}

	topic, err := models.CreateTopic(db, req.BoardID, req.Title, req.Author)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, topic)
}

type CreatePostRequest struct {
	TopicID int    `json:"topic_id"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	db := db.FromContext(r.Context())

	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "Invalid JSON"})
		return
	}

	if req.TopicID == 0 || req.Author == "" || req.Content == "" {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "Topic ID, author, and content are required"})
		return
	}

	_, err := models.GetTopicByID(db, req.TopicID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, utils.APIError{Detail: "Topic not found"})
		} else {
			utils.InternalServerError(w, err)
		}
		return
	}

	post, err := models.CreatePost(db, req.TopicID, req.Author, req.Content)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, post)
}

type ListTopicsRequest struct {
	BoardID int  `json:"board_id"`
	Limit   int  `json:"limit,omitempty"`
	Cursor  *int `json:"cursor,omitempty"`
}

func ListTopics(w http.ResponseWriter, r *http.Request) {
	db := db.FromContext(r.Context())

	var req ListTopicsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "Invalid JSON"})
		return
	}

	if req.BoardID == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "board_id is required"})
		return
	}

	if req.Limit == 0 {
		req.Limit = 20
	}

	opts := models.TopicsListOptions{
		BoardID: req.BoardID,
		Limit:   req.Limit,
		Cursor:  req.Cursor,
	}

	topics, err := models.ListTopics(db, opts)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	response := utils.CreatePaginatedResponse(topics, req.Limit)
	utils.RespondWithJSON(w, http.StatusOK, response)
}

type ListPostsRequest struct {
	TopicID int  `json:"topic_id"`
	Limit   int  `json:"limit,omitempty"`
	Cursor  *int `json:"cursor,omitempty"`
}

func ListPosts(w http.ResponseWriter, r *http.Request) {
	db := db.FromContext(r.Context())

	var req ListPostsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "Invalid JSON"})
		return
	}

	if req.TopicID == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "topic_id is required"})
		return
	}

	_, err := models.GetTopicByID(db, req.TopicID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, utils.APIError{Detail: "Topic not found"})
		} else {
			utils.InternalServerError(w, err)
		}
		return
	}

	if req.Limit == 0 {
		req.Limit = 20
	}

	opts := models.PostsListOptions{
		TopicID: req.TopicID,
		Limit:   req.Limit,
		Cursor:  req.Cursor,
	}

	posts, err := models.ListPosts(db, opts)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	response := utils.CreatePaginatedResponse(posts, req.Limit)
	utils.RespondWithJSON(w, http.StatusOK, response)
}

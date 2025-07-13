package handlers

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"

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

type BoardsResponse struct {
	Boards []BoardWithRecent `json:"boards"`
}

type BoardWithRecent struct {
	models.Board
	RecentTopic *models.Topic `json:"recent_topic"`
	RecentPost  *models.Post  `json:"recent_post"`
}

func ListBoards(w http.ResponseWriter, r *http.Request) {
	database := db.FromContext(r.Context())

	boards, err := getBoardsWithRecent(database)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	response := BoardsResponse{Boards: boards}
	utils.RespondWithJSON(w, http.StatusOK, response)
}

func getBoardsWithRecent(database *sql.DB) ([]BoardWithRecent, error) {
	// Get all boards using the model
	boards, err := models.GetAllBoards(database)
	if err != nil {
		return nil, err
	}

	// For each board, get the most recent topic and post
	var boardsWithRecent []BoardWithRecent
	for _, board := range boards {
		boardWithRecent := BoardWithRecent{Board: board}

		// Get most recent topic for this board
		recentTopic, err := models.GetMostRecentTopicByBoardID(database, board.ID)
		if err != nil {
			return nil, err
		}
		boardWithRecent.RecentTopic = recentTopic

		// Get most recent post for this board
		recentPost, err := models.GetMostRecentPostByBoardID(database, board.ID)
		if err != nil {
			return nil, err
		}
		boardWithRecent.RecentPost = recentPost

		boardsWithRecent = append(boardsWithRecent, boardWithRecent)
	}

	return boardsWithRecent, nil
}

type TopicsResponse struct {
	Topics     []models.Topic       `json:"topics"`
	Pagination utils.PaginationMeta `json:"pagination"`
}

func ListTopics(w http.ResponseWriter, r *http.Request) {
	database := db.FromContext(r.Context())
	boardSlug := chi.URLParam(r, "board")

	board, err := models.GetBoardBySlug(database, boardSlug)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}
	if board == nil {
		utils.RespondWithError(w, http.StatusNotFound, utils.APIError{Detail: "board not found"})
		return
	}

	params := utils.ParsePaginationParams(r)

	topics, err := models.GetTopicsByBoardIDWithPagination(database, board.ID, params.PerPage, params.Offset)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	total, err := models.CountTopicsByBoardID(database, board.ID)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	meta := utils.CalculatePaginationMeta(params.Page, params.PerPage, total)
	response := TopicsResponse{
		Topics:     topics,
		Pagination: meta,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

type PostsResponse struct {
	Posts      []models.Post        `json:"posts"`
	Topic      *models.Topic        `json:"topic"`
	Pagination utils.PaginationMeta `json:"pagination"`
}

func ListPosts(w http.ResponseWriter, r *http.Request) {
	database := db.FromContext(r.Context())
	topicIDStr := chi.URLParam(r, "topicId")

	topicID, err := utils.ParseInt(topicIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.APIError{Detail: "invalid topic ID"})
		return
	}

	topic, err := models.GetTopicByID(database, topicID)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}
	if topic == nil {
		utils.RespondWithError(w, http.StatusNotFound, utils.APIError{Detail: "topic not found"})
		return
	}

	params := utils.ParsePaginationParams(r)

	posts, err := models.GetPostsByTopicIDWithPagination(database, topicID, params.PerPage, params.Offset)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	total, err := models.CountPostsByTopicID(database, topicID)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}

	meta := utils.CalculatePaginationMeta(params.Page, params.PerPage, total)
	response := PostsResponse{
		Posts:      posts,
		Topic:      topic,
		Pagination: meta,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

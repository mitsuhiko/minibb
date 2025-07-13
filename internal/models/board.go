package models

import (
	"database/sql"

	"minibb/internal/db"
)

type Board struct {
	ID          int    `json:"id"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

func GetAllBoards(q db.Querier) ([]Board, error) {
	query := `SELECT id, slug, description FROM boards ORDER BY id`
	rows, err := q.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []Board
	for rows.Next() {
		var board Board
		if err := rows.Scan(&board.ID, &board.Slug, &board.Description); err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}

	return boards, rows.Err()
}

func GetBoardByID(q db.Querier, id int) (*Board, error) {
	query := `SELECT id, slug, description FROM boards WHERE id = ?`
	var board Board
	err := q.QueryRow(query, id).Scan(&board.ID, &board.Slug, &board.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &board, nil
}

func GetBoardBySlug(q db.Querier, slug string) (*Board, error) {
	query := `SELECT id, slug, description FROM boards WHERE slug = ?`
	var board Board
	err := q.QueryRow(query, slug).Scan(&board.ID, &board.Slug, &board.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &board, nil
}

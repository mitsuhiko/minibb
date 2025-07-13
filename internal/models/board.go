package models

import (
	"database/sql"
)

type Board struct {
	ID          int    `json:"id"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

func GetAllBoards(db *sql.DB) ([]Board, error) {
	query := `SELECT id, slug, description FROM boards ORDER BY id`
	rows, err := db.Query(query)
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

func GetBoardByID(db *sql.DB, id int) (*Board, error) {
	query := `SELECT id, slug, description FROM boards WHERE id = ?`
	var board Board
	err := db.QueryRow(query, id).Scan(&board.ID, &board.Slug, &board.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &board, nil
}

func GetBoardBySlug(db *sql.DB, slug string) (*Board, error) {
	query := `SELECT id, slug, description FROM boards WHERE slug = ?`
	var board Board
	err := db.QueryRow(query, slug).Scan(&board.ID, &board.Slug, &board.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &board, nil
}

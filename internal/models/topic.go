package models

import (
	"database/sql"
	"time"
)

type Topic struct {
	ID         int       `json:"id"`
	BoardID    int       `json:"board_id"`
	Title      string    `json:"title"`
	Author     string    `json:"author"`
	PubDate    time.Time `json:"pub_date"`
	Status     string    `json:"status"`
	LastPostID *int      `json:"last_post_id"`
	PostCount  int       `json:"post_count"`
}

func GetTopicByID(db *sql.DB, id int) (*Topic, error) {
	query := `SELECT id, board_id, title, author, pub_date, status, last_post_id, post_count 
		FROM topics WHERE id = ?`
	var topic Topic
	err := db.QueryRow(query, id).Scan(
		&topic.ID, &topic.BoardID, &topic.Title, &topic.Author,
		&topic.PubDate, &topic.Status, &topic.LastPostID, &topic.PostCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &topic, nil
}

func GetTopicsByBoardID(db *sql.DB, boardID int) ([]Topic, error) {
	query := `SELECT id, board_id, title, author, pub_date, status, last_post_id, post_count 
		FROM topics WHERE board_id = ? ORDER BY pub_date DESC`
	rows, err := db.Query(query, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []Topic
	for rows.Next() {
		var topic Topic
		if err := rows.Scan(
			&topic.ID, &topic.BoardID, &topic.Title, &topic.Author,
			&topic.PubDate, &topic.Status, &topic.LastPostID, &topic.PostCount,
		); err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}

	return topics, rows.Err()
}

func GetMostRecentTopicByBoardID(db *sql.DB, boardID int) (*Topic, error) {
	query := `SELECT id, board_id, title, author, pub_date, status, last_post_id, post_count 
		FROM topics WHERE board_id = ? ORDER BY pub_date DESC LIMIT 1`
	var topic Topic
	err := db.QueryRow(query, boardID).Scan(
		&topic.ID, &topic.BoardID, &topic.Title, &topic.Author,
		&topic.PubDate, &topic.Status, &topic.LastPostID, &topic.PostCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &topic, nil
}

func GetTopicsByBoardIDWithPagination(db *sql.DB, boardID int, limit, offset int) ([]Topic, error) {
	query := `SELECT id, board_id, title, author, pub_date, status, last_post_id, post_count 
		FROM topics WHERE board_id = ? ORDER BY pub_date DESC LIMIT ? OFFSET ?`
	rows, err := db.Query(query, boardID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []Topic
	for rows.Next() {
		var topic Topic
		if err := rows.Scan(
			&topic.ID, &topic.BoardID, &topic.Title, &topic.Author,
			&topic.PubDate, &topic.Status, &topic.LastPostID, &topic.PostCount,
		); err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}

	return topics, rows.Err()
}

func CountTopicsByBoardID(db *sql.DB, boardID int) (int, error) {
	query := `SELECT COUNT(*) FROM topics WHERE board_id = ?`
	var count int
	err := db.QueryRow(query, boardID).Scan(&count)
	return count, err
}

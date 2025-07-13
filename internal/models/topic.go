package models

import (
	"database/sql"
	"time"

	"minibb/internal/db"
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

func GetTopicByID(q db.Querier, id int) (*Topic, error) {
	query := `SELECT id, board_id, title, author, pub_date, status, last_post_id, post_count 
		FROM topics WHERE id = ?`
	var topic Topic
	err := q.QueryRow(query, id).Scan(
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

func GetTopicsByBoardID(q db.Querier, boardID int) ([]Topic, error) {
	query := `SELECT id, board_id, title, author, pub_date, status, last_post_id, post_count 
		FROM topics WHERE board_id = ? ORDER BY pub_date DESC`
	rows, err := q.Query(query, boardID)
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

func GetMostRecentTopicByBoardID(q db.Querier, boardID int) (*Topic, error) {
	query := `SELECT id, board_id, title, author, pub_date, status, last_post_id, post_count 
		FROM topics WHERE board_id = ? ORDER BY pub_date DESC LIMIT 1`
	var topic Topic
	err := q.QueryRow(query, boardID).Scan(
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

func GetTopicsByBoardIDWithPagination(q db.Querier, boardID int, limit, offset int) ([]Topic, error) {
	query := `SELECT id, board_id, title, author, pub_date, status, last_post_id, post_count 
		FROM topics WHERE board_id = ? ORDER BY pub_date DESC LIMIT ? OFFSET ?`
	rows, err := q.Query(query, boardID, limit, offset)
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

func CountTopicsByBoardID(q db.Querier, boardID int) (int, error) {
	query := `SELECT COUNT(*) FROM topics WHERE board_id = ?`
	var count int
	err := q.QueryRow(query, boardID).Scan(&count)
	return count, err
}

// CreateTopic creates a new topic with an initial post atomically
func CreateTopic(tm *db.TxManager, boardID int, title, author, content string) (*Topic, error) {
	var topicID int64

	err := tm.WithTx(func(q db.Querier) error {
		// Create the topic
		topicQuery := `INSERT INTO topics (board_id, title, author) VALUES (?, ?, ?)`
		topicResult, err := q.Exec(topicQuery, boardID, title, author)
		if err != nil {
			return err
		}

		topicID, err = topicResult.LastInsertId()
		if err != nil {
			return err
		}

		// Create the initial post
		postQuery := `INSERT INTO posts (topic_id, author, content) VALUES (?, ?, ?)`
		postResult, err := q.Exec(postQuery, topicID, author, content)
		if err != nil {
			return err
		}

		postID, err := postResult.LastInsertId()
		if err != nil {
			return err
		}

		// Update topic with the post ID
		updateTopicQuery := `UPDATE topics SET last_post_id = ? WHERE id = ?`
		_, err = q.Exec(updateTopicQuery, postID, topicID)
		return err
	})

	if err != nil {
		return nil, err
	}

	// Return the created topic using the transaction manager's querier
	return GetTopicByID(tm.Querier(), int(topicID))
}

// CreateTopicWithQuerier creates a topic assuming we're already in the right transactional context
func CreateTopicWithQuerier(q db.Querier, boardID int, title, author, content string) (*Topic, error) {
	// Create the topic
	topicQuery := `INSERT INTO topics (board_id, title, author) VALUES (?, ?, ?)`
	topicResult, err := q.Exec(topicQuery, boardID, title, author)
	if err != nil {
		return nil, err
	}

	topicID, err := topicResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Create the initial post
	postQuery := `INSERT INTO posts (topic_id, author, content) VALUES (?, ?, ?)`
	postResult, err := q.Exec(postQuery, topicID, author, content)
	if err != nil {
		return nil, err
	}

	postID, err := postResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Update topic with the post ID
	updateTopicQuery := `UPDATE topics SET last_post_id = ? WHERE id = ?`
	_, err = q.Exec(updateTopicQuery, postID, topicID)
	if err != nil {
		return nil, err
	}

	return GetTopicByID(q, int(topicID))
}

package models

import (
	"database/sql"
	"time"
)

type Topic struct {
	ID         int       `json:"id"`
	BoardID    int       `json:"board_id"`
	PubDate    time.Time `json:"pub_date"`
	Title      string    `json:"title"`
	Status     string    `json:"status"`
	Author     string    `json:"author"`
	LastPostID *int      `json:"last_post_id"`
	PostCount  int       `json:"post_count"`
}

func (t *Topic) GetID() int {
	return t.ID
}

func CreateTopic(db *sql.DB, boardID int, title, author string) (*Topic, error) {
	query := `
		INSERT INTO topics (board_id, title, author)
		VALUES (?, ?, ?)
		RETURNING id, board_id, pub_date, title, status, author, last_post_id, post_count
	`
	
	row := db.QueryRow(query, boardID, title, author)
	
	topic := &Topic{}
	err := row.Scan(
		&topic.ID,
		&topic.BoardID,
		&topic.PubDate,
		&topic.Title,
		&topic.Status,
		&topic.Author,
		&topic.LastPostID,
		&topic.PostCount,
	)
	if err != nil {
		return nil, err
	}
	
	return topic, nil
}

func GetTopicByID(db *sql.DB, id int) (*Topic, error) {
	query := `
		SELECT id, board_id, pub_date, title, status, author, last_post_id, post_count
		FROM topics
		WHERE id = ?
	`
	
	row := db.QueryRow(query, id)
	
	topic := &Topic{}
	err := row.Scan(
		&topic.ID,
		&topic.BoardID,
		&topic.PubDate,
		&topic.Title,
		&topic.Status,
		&topic.Author,
		&topic.LastPostID,
		&topic.PostCount,
	)
	if err != nil {
		return nil, err
	}
	
	return topic, nil
}

type TopicsListOptions struct {
	BoardID int
	Limit   int
	Cursor  *int
}

func ListTopics(db *sql.DB, opts TopicsListOptions) ([]*Topic, error) {
	query := `
		SELECT id, board_id, pub_date, title, status, author, last_post_id, post_count
		FROM topics
		WHERE board_id = ?
	`
	args := []interface{}{opts.BoardID}
	
	if opts.Cursor != nil {
		query += " AND id < ?"
		args = append(args, *opts.Cursor)
	}
	
	query += " ORDER BY id DESC LIMIT ?"
	args = append(args, opts.Limit)
	
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var topics []*Topic
	for rows.Next() {
		topic := &Topic{}
		err := rows.Scan(
			&topic.ID,
			&topic.BoardID,
			&topic.PubDate,
			&topic.Title,
			&topic.Status,
			&topic.Author,
			&topic.LastPostID,
			&topic.PostCount,
		)
		if err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}
	
	return topics, rows.Err()
}


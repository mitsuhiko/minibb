package models

import (
	"database/sql"
	"time"

	"minibb/internal/db"
)

type Post struct {
	ID      int       `json:"id"`
	TopicID int       `json:"topic_id"`
	Author  string    `json:"author"`
	Content string    `json:"content"`
	PubDate time.Time `json:"pub_date"`
}

func GetPostByID(q db.Querier, id int) (*Post, error) {
	query := `SELECT id, topic_id, author, content, pub_date FROM posts WHERE id = ?`
	var post Post
	err := q.QueryRow(query, id).Scan(
		&post.ID, &post.TopicID, &post.Author, &post.Content, &post.PubDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func GetPostsByTopicID(q db.Querier, topicID int) ([]Post, error) {
	query := `SELECT id, topic_id, author, content, pub_date 
		FROM posts WHERE topic_id = ? ORDER BY pub_date ASC`
	rows, err := q.Query(query, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(
			&post.ID, &post.TopicID, &post.Author, &post.Content, &post.PubDate,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, rows.Err()
}

func GetMostRecentPostByBoardID(q db.Querier, boardID int) (*Post, error) {
	query := `SELECT p.id, p.topic_id, p.author, p.content, p.pub_date
		FROM posts p
		JOIN topics t ON p.topic_id = t.id
		WHERE t.board_id = ?
		ORDER BY p.pub_date DESC
		LIMIT 1`
	var post Post
	err := q.QueryRow(query, boardID).Scan(
		&post.ID, &post.TopicID, &post.Author, &post.Content, &post.PubDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func GetMostRecentPostByTopicID(q db.Querier, topicID int) (*Post, error) {
	query := `SELECT id, topic_id, author, content, pub_date
		FROM posts WHERE topic_id = ?
		ORDER BY pub_date DESC
		LIMIT 1`
	var post Post
	err := q.QueryRow(query, topicID).Scan(
		&post.ID, &post.TopicID, &post.Author, &post.Content, &post.PubDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func GetPostsByTopicIDWithPagination(q db.Querier, topicID int, limit, offset int) ([]Post, error) {
	query := `SELECT id, topic_id, author, content, pub_date 
		FROM posts WHERE topic_id = ? ORDER BY pub_date ASC LIMIT ? OFFSET ?`
	rows, err := q.Query(query, topicID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(
			&post.ID, &post.TopicID, &post.Author, &post.Content, &post.PubDate,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, rows.Err()
}

func CountPostsByTopicID(q db.Querier, topicID int) (int, error) {
	query := `SELECT COUNT(*) FROM posts WHERE topic_id = ?`
	var count int
	err := q.QueryRow(query, topicID).Scan(&count)
	return count, err
}

func CreatePost(q db.Querier, topicID int, author, content string) (*Post, error) {
	query := `INSERT INTO posts (topic_id, author, content) VALUES (?, ?, ?)`
	result, err := q.Exec(query, topicID, author, content)
	if err != nil {
		return nil, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetPostByID(q, int(postID))
}

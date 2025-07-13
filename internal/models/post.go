package models

import (
	"database/sql"
	"time"
)

type Post struct {
	ID      int       `json:"id"`
	TopicID int       `json:"topic_id"`
	PubDate time.Time `json:"pub_date"`
	Author  string    `json:"author"`
	Content string    `json:"content"`
}

func (p *Post) GetID() int {
	return p.ID
}

func CreatePost(db *sql.DB, topicID int, author, content string) (*Post, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO posts (topic_id, author, content)
		VALUES (?, ?, ?)
		RETURNING id, topic_id, pub_date, author, content
	`
	row := tx.QueryRow(query, topicID, author, content)

	post := &Post{}
	err = row.Scan(
		&post.ID,
		&post.TopicID,
		&post.PubDate,
		&post.Author,
		&post.Content,
	)
	if err != nil {
		return nil, err
	}

	query = `
		UPDATE topics 
		SET last_post_id = ?, post_count = post_count + 1
		WHERE id = ?
	`
	_, err = tx.Exec(query, post.ID, topicID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return post, nil
}

func GetPostByID(db *sql.DB, id int) (*Post, error) {
	query := `
		SELECT id, topic_id, pub_date, author, content
		FROM posts
		WHERE id = ?
	`

	row := db.QueryRow(query, id)

	post := &Post{}
	err := row.Scan(
		&post.ID,
		&post.TopicID,
		&post.PubDate,
		&post.Author,
		&post.Content,
	)
	if err != nil {
		return nil, err
	}

	return post, nil
}

type PostsListOptions struct {
	TopicID int
	Limit   int
	Cursor  *int
}

func ListPosts(db *sql.DB, opts PostsListOptions) ([]*Post, error) {
	query := `
		SELECT id, topic_id, pub_date, author, content
		FROM posts
		WHERE topic_id = ?
	`
	args := []interface{}{opts.TopicID}

	if opts.Cursor != nil {
		query += " AND id > ?"
		args = append(args, *opts.Cursor)
	}

	query += " ORDER BY id ASC LIMIT ?"
	args = append(args, opts.Limit)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		err := rows.Scan(
			&post.ID,
			&post.TopicID,
			&post.PubDate,
			&post.Author,
			&post.Content,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, rows.Err()
}

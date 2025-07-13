package database

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

func Init() (*sql.DB, error) {
	dbPath := getDBPath()

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// TODO: Run migrations
	if err := createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func getDBPath() string {
	if dbPath := os.Getenv("DATABASE_PATH"); dbPath != "" {
		return dbPath
	}
	return "minibb.db"
}

func createTables(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS boards (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		slug TEXT UNIQUE NOT NULL,
		description TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS topics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		board_id INTEGER NOT NULL,
		pub_date DATETIME DEFAULT CURRENT_TIMESTAMP,
		title TEXT NOT NULL,
		status TEXT DEFAULT 'open',
		author TEXT NOT NULL,
		last_post_id INTEGER,
		post_count INTEGER DEFAULT 1,
		FOREIGN KEY (board_id) REFERENCES boards(id)
	);

	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		topic_id INTEGER NOT NULL,
		pub_date DATETIME DEFAULT CURRENT_TIMESTAMP,
		author TEXT NOT NULL,
		content TEXT NOT NULL,
		FOREIGN KEY (topic_id) REFERENCES topics(id)
	);

	CREATE INDEX IF NOT EXISTS idx_topics_board_id ON topics(board_id);
	CREATE INDEX IF NOT EXISTS idx_posts_topic_id ON posts(topic_id);
	`

	_, err := db.Exec(schema)
	return err
}

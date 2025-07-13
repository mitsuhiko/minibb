package db

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func Init() (*sql.DB, error) {
	dbPath := getDBPath()

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := runMigrations(db); err != nil {
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

func createMigrationsTable(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS migrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		migration_number INTEGER NOT NULL UNIQUE,
		filename TEXT NOT NULL,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(schema)
	return err
}

func runMigrations(db *sql.DB) error {
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	files, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	sort.Strings(migrationFiles)

	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	for _, filename := range migrationFiles {
		migrationNumber, err := extractMigrationNumber(filename)
		if err != nil {
			return fmt.Errorf("invalid migration filename %s: %w", filename, err)
		}

		if appliedMigrations[migrationNumber] {
			continue
		}

		if err := applyMigration(db, migrationNumber, filename); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", filename, err)
		}
	}

	return nil
}

func getAppliedMigrations(db *sql.DB) (map[int]bool, error) {
	applied := make(map[int]bool)

	rows, err := db.Query("SELECT migration_number FROM migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var migrationNumber int
		if err := rows.Scan(&migrationNumber); err != nil {
			return nil, err
		}
		applied[migrationNumber] = true
	}

	return applied, rows.Err()
}

func extractMigrationNumber(filename string) (int, error) {
	parts := strings.Split(filename, "_")
	if len(parts) < 2 {
		return 0, fmt.Errorf("filename should be in format NNNN_description.sql")
	}

	migrationNumber, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("migration number should be numeric: %w", err)
	}

	return migrationNumber, nil
}

func applyMigration(db *sql.DB, migrationNumber int, filename string) error {
	migrationPath := "migrations/" + filename
	content, err := fs.ReadFile(migrationsFS, migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	if _, err := tx.Exec(
		"INSERT INTO migrations (migration_number, filename) VALUES (?, ?)",
		migrationNumber, filename,
	); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	fmt.Printf("Applied migration: %s\n", filename)
	return nil
}

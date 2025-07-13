package models

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"

	"minibb/internal/db"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	database, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Run the real migrations to ensure test schema matches production
	if err := db.RunMigrationsForTesting(database); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	cleanup := func() {
		database.Close()
	}

	return database, cleanup
}

func TestGetAllBoards(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	boards, err := GetAllBoards(database)
	if err != nil {
		t.Fatalf("GetAllBoards failed: %v", err)
	}

	if len(boards) != 2 {
		t.Errorf("Expected 2 boards, got %d", len(boards))
	}

	expectedSlugs := map[string]bool{"general": true, "watercooler": true}
	if !expectedSlugs[boards[0].Slug] || !expectedSlugs[boards[1].Slug] {
		t.Errorf("Expected boards to have slugs 'general' and 'watercooler', got '%s' and '%s'", boards[0].Slug, boards[1].Slug)
	}
}

func TestGetAllBoardsWithTransaction(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	// Start a transaction for testing
	tx, err := database.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback() // Always rollback in tests

	// Test reading within transaction
	boards, err := GetAllBoards(tx)
	if err != nil {
		t.Fatalf("GetAllBoards with transaction failed: %v", err)
	}

	if len(boards) != 2 {
		t.Errorf("Expected 2 boards, got %d", len(boards))
	}

	// Insert a new board within the transaction
	_, err = tx.Exec("INSERT INTO boards (slug, description) VALUES (?, ?)", "test", "Test Board")
	if err != nil {
		t.Fatalf("Failed to insert test board: %v", err)
	}

	// Should see 3 boards within the transaction
	boards, err = GetAllBoards(tx)
	if err != nil {
		t.Fatalf("GetAllBoards after insert failed: %v", err)
	}

	if len(boards) != 3 {
		t.Errorf("Expected 3 boards within transaction, got %d", len(boards))
	}

	// After rollback, should still see only 2 boards
	tx.Rollback()

	boards, err = GetAllBoards(database)
	if err != nil {
		t.Fatalf("GetAllBoards after rollback failed: %v", err)
	}

	if len(boards) != 2 {
		t.Errorf("Expected 2 boards after rollback, got %d", len(boards))
	}
}

func TestTransactionManager(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	tm := db.NewTxManager(database)

	// Test nested transactions with rollback
	err := tm.WithTx(func(q db.Querier) error {
		// Insert a board in outer transaction
		_, err := q.Exec("INSERT INTO boards (slug, description) VALUES (?, ?)", "outer", "Outer Board")
		if err != nil {
			return err
		}

		// Verify we can see it
		boards, err := GetAllBoards(q)
		if err != nil {
			return err
		}
		if len(boards) != 3 {
			t.Errorf("Expected 3 boards in outer transaction, got %d", len(boards))
		}

		// Nested transaction that will be rolled back
		return tm.WithTx(func(q2 db.Querier) error {
			// Insert another board in nested transaction
			_, err := q2.Exec("INSERT INTO boards (slug, description) VALUES (?, ?)", "nested", "Nested Board")
			if err != nil {
				return err
			}

			// Verify we can see both
			boards, err := GetAllBoards(q2)
			if err != nil {
				return err
			}
			if len(boards) != 4 {
				t.Errorf("Expected 4 boards in nested transaction, got %d", len(boards))
			}

			// Return error to trigger rollback of nested transaction
			return sql.ErrNoRows
		})
		// The nested transaction should be rolled back, but outer should continue
	})

	// The entire outer transaction should be rolled back due to nested error
	if err == nil {
		t.Error("Expected error from nested transaction rollback")
	}

	// Should still have only original 2 boards
	boards, err := GetAllBoards(database)
	if err != nil {
		t.Fatalf("GetAllBoards after nested rollback failed: %v", err)
	}

	if len(boards) != 2 {
		t.Errorf("Expected 2 boards after full rollback, got %d", len(boards))
	}
}

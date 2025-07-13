package db

import (
	"database/sql"
	"fmt"
)

// TxManager handles nested transactions using savepoints
type TxManager struct {
	db           *sql.DB
	tx           *sql.Tx
	nestingLevel int
	savepoints   []string
}

// NewTxManager creates a new transaction manager
func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{
		db:           db,
		nestingLevel: 0,
		savepoints:   make([]string, 0),
	}
}

// Begin starts a transaction or creates a savepoint if already in a transaction
func (tm *TxManager) Begin() error {
	if tm.nestingLevel == 0 {
		// Start the root transaction
		tx, err := tm.db.Begin()
		if err != nil {
			return err
		}
		tm.tx = tx
		tm.nestingLevel = 1
		return nil
	}

	// Create a savepoint for nested transaction
	tm.nestingLevel++
	savepointName := fmt.Sprintf("sp_%d", tm.nestingLevel)
	tm.savepoints = append(tm.savepoints, savepointName)

	_, err := tm.tx.Exec(fmt.Sprintf("SAVEPOINT %s", savepointName))
	return err
}

// Commit commits the transaction or releases a savepoint
func (tm *TxManager) Commit() error {
	if tm.nestingLevel == 0 {
		return fmt.Errorf("no active transaction to commit")
	}

	if tm.nestingLevel == 1 {
		// Commit the root transaction
		err := tm.tx.Commit()
		tm.tx = nil
		tm.nestingLevel = 0
		tm.savepoints = tm.savepoints[:0]
		return err
	}

	// Release the savepoint
	savepointName := tm.savepoints[len(tm.savepoints)-1]
	tm.savepoints = tm.savepoints[:len(tm.savepoints)-1]
	tm.nestingLevel--

	_, err := tm.tx.Exec(fmt.Sprintf("RELEASE SAVEPOINT %s", savepointName))
	return err
}

// Rollback rolls back the transaction or to a savepoint
func (tm *TxManager) Rollback() error {
	if tm.nestingLevel == 0 {
		return fmt.Errorf("no active transaction to rollback")
	}

	if tm.nestingLevel == 1 {
		// Rollback the root transaction
		err := tm.tx.Rollback()
		tm.tx = nil
		tm.nestingLevel = 0
		tm.savepoints = tm.savepoints[:0]
		return err
	}

	// Rollback to savepoint
	savepointName := tm.savepoints[len(tm.savepoints)-1]
	tm.savepoints = tm.savepoints[:len(tm.savepoints)-1]
	tm.nestingLevel--

	_, err := tm.tx.Exec(fmt.Sprintf("ROLLBACK TO SAVEPOINT %s", savepointName))
	return err
}

// Querier returns the current querier (either db or tx)
func (tm *TxManager) Querier() Querier {
	if tm.tx != nil {
		return tm.tx
	}
	return tm.db
}

// InTransaction returns true if currently in a transaction
func (tm *TxManager) InTransaction() bool {
	return tm.nestingLevel > 0
}

// WithTx executes a function within a transaction, handling nesting automatically
func (tm *TxManager) WithTx(fn func(Querier) error) error {
	if err := tm.Begin(); err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tm.Rollback()
			panic(r)
		}
	}()

	if err := fn(tm.Querier()); err != nil {
		if rollbackErr := tm.Rollback(); rollbackErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %v", err, rollbackErr)
		}
		return err
	}

	return tm.Commit()
}

// WithNewTxManager executes a function with a new transaction manager
func WithNewTxManager(db *sql.DB, fn func(*TxManager) error) error {
	tm := NewTxManager(db)
	return fn(tm)
}

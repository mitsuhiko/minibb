package db

import (
	"context"
	"database/sql"
)

type contextKey string

const (
	dbContextKey        contextKey = "database"
	txManagerContextKey contextKey = "txmanager"
)

func WithDB(ctx context.Context, db *sql.DB) context.Context {
	return context.WithValue(ctx, dbContextKey, db)
}

func WithQuerier(ctx context.Context, q Querier) context.Context {
	return context.WithValue(ctx, dbContextKey, q)
}

func WithTxManager(ctx context.Context, tm *TxManager) context.Context {
	return context.WithValue(ctx, txManagerContextKey, tm)
}

func FromContext(ctx context.Context) *sql.DB {
	db, ok := ctx.Value(dbContextKey).(*sql.DB)
	if !ok {
		panic("database not found in context")
	}
	return db
}

func QuerierFromContext(ctx context.Context) Querier {
	q, ok := ctx.Value(dbContextKey).(Querier)
	if !ok {
		panic("querier not found in context")
	}
	return q
}

func TxManagerFromContext(ctx context.Context) *TxManager {
	tm, ok := ctx.Value(txManagerContextKey).(*TxManager)
	if !ok {
		panic("transaction manager not found in context")
	}
	return tm
}

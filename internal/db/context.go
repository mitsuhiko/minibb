package db

import (
	"context"
	"database/sql"
)

type contextKey string

const dbContextKey contextKey = "database"

func WithDB(ctx context.Context, db *sql.DB) context.Context {
	return context.WithValue(ctx, dbContextKey, db)
}

func FromContext(ctx context.Context) *sql.DB {
	db, ok := ctx.Value(dbContextKey).(*sql.DB)
	if !ok {
		panic("database not found in context")
	}
	return db
}

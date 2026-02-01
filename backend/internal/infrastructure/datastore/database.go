package datastore

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Database defines the interface for database operations
type Database interface {
	// Connect establishes database connection
	Connect(ctx context.Context) error

	// Query executes a query that returns rows
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)

	// QueryRow executes a query that returns at most one row
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row

	// Execute executes INSERT/UPDATE/DELETE operations
	Execute(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)

	// Begin starts a transaction
	Begin(ctx context.Context) (pgx.Tx, error)

	// Close closes the database connection
	Close()
}

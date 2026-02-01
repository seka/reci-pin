package datastore

import "context"

// Database defines the interface for database operations
type Database interface {
	// Connect establishes database connection
	Connect(ctx context.Context) error

	// Query executes a query and returns results as []T
	// Uses reflection to map columns to struct fields via `db` tags
	Query[T any](ctx context.Context, query string, args ...any) ([]T, error)

	// Execute executes INSERT/UPDATE/DELETE and returns rows affected
	Execute(ctx context.Context, query string, args ...any) (int64, error)

	// Begin starts a transaction
	Begin(ctx context.Context) (Tx, error)

	// Close closes the database connection
	Close()
}

// Tx represents a database transaction
type Tx interface {
	// Query executes a query within the transaction
	Query[T any](ctx context.Context, query string, args ...any) ([]T, error)

	// Execute executes a command within the transaction
	Execute(ctx context.Context, query string, args ...any) (int64, error)

	// Commit commits the transaction
	Commit(ctx context.Context) error

	// Rollback rolls back the transaction
	Rollback(ctx context.Context) error
}

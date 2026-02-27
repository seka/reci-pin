package database

import (
	"context"

	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

// Database defines the interface for database operations
//
//go:generate mockgen -source=$GOFILE -destination=mock/database_mock.go -package=mock
type Database interface {
	Connect(ctx context.Context) error

	// Query executes a query and returns abstract Rows interface
	Query(ctx context.Context, query string, args ...any) (Rows, error)

	// Execute executes INSERT/UPDATE/DELETE operations and returns affected rows count
	Execute(ctx context.Context, query string, args ...any) (int64, error)

	// TransactionManager returns the transaction manager for this database
	TransactionManager() repository.TransactionManager

	// Close closes the database connection
	Close()
}

// Rows defines the interface for iterating over query results
type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close()
	Err() error
}

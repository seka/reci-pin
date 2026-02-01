package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/seka/reci-pin/backend/internal/infrastructure/datastore"
)

type DB struct {
	dsn  string
	pool *pgxpool.Pool
}

// New creates a new PostgreSQL database instance with configuration
func New(dsn string) datastore.Database {
	return &DB{
		dsn: dsn,
	}
}

// Connect establishes the database connection
func (db *DB) Connect(ctx context.Context) error {
	pool, err := pgxpool.New(ctx, db.dsn)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("unable to ping database: %w", err)
	}

	db.pool = pool
	return nil
}

// Query executes a query that returns rows
func (db *DB) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return db.pool.Query(ctx, query, args...)
}

// QueryRow executes a query that returns at most one row
func (db *DB) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return db.pool.QueryRow(ctx, query, args...)
}

// Execute executes INSERT/UPDATE/DELETE operations
func (db *DB) Execute(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return db.pool.Exec(ctx, query, args...)
}

// Begin starts a transaction
func (db *DB) Begin(ctx context.Context) (pgx.Tx, error) {
	return db.pool.Begin(ctx)
}

// Close closes the database connection
func (db *DB) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

// Compile-time check that DB implements Database interface
var _ datastore.Database = (*DB)(nil)

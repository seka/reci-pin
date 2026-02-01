package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Database defines the interface for database operations
type Database interface {
	Connect(ctx context.Context) error

	// Query executes a query and returns abstract Rows interface
	Query(ctx context.Context, query string, args ...any) (Rows, error)

	// Execute executes INSERT/UPDATE/DELETE operations and returns affected rows count
	Execute(ctx context.Context, query string, args ...any) (int64, error)

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

type Postgres struct {
	dsn  string
	pool *pgxpool.Pool
}

func New(dsn string) Database {
	return &Postgres{dsn: dsn}
}

func (p *Postgres) Connect(ctx context.Context) error {
	pool, err := pgxpool.New(ctx, p.dsn)
	if err != nil {
		return err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return err
	}
	p.pool = pool
	return nil
}

func (p *Postgres) Query(ctx context.Context, query string, args ...any) (Rows, error) {
	// pgx.Rows implements datastore.Rows implicitly
	return p.pool.Query(ctx, query, args...)
}

func (p *Postgres) Execute(ctx context.Context, query string, args ...any) (int64, error) {
	tag, err := p.pool.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func (p *Postgres) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}

// Check if Postgres implements Database interface
var _ Database = (*Postgres)(nil)

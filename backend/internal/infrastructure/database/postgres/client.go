package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database"
)

type Client struct {
	dsn  string
	pool *pgxpool.Pool
}

func New(dsn string) database.Database {
	return &Client{dsn: dsn}
}

func (p *Client) Connect(ctx context.Context) error {
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

func (p *Client) Query(ctx context.Context, query string, args ...any) (database.Rows, error) {
	// pgx.Rows implements datastore.Rows implicitly
	return p.pool.Query(ctx, query, args...)
}

func (p *Client) Execute(ctx context.Context, query string, args ...any) (int64, error) {
	tag, err := p.pool.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func (p *Client) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}

// Check if Postgres implements Database interface
var _ database.Database = (*Client)(nil)

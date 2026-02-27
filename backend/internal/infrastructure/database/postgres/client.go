package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database"
)

type Client struct {
	cfg  config.Database
	pool *pgxpool.Pool
}

func NewClient(cfg config.Database) database.Database {
	return &Client{cfg: cfg}
}

func (p *Client) Connect(ctx context.Context) error {
	pool, err := pgxpool.New(ctx, p.cfg.DSN())
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
	if tx, ok := getTx(ctx); ok {
		return tx.Query(ctx, query, args...)
	}
	return p.pool.Query(ctx, query, args...)
}

func (p *Client) Execute(ctx context.Context, query string, args ...any) (int64, error) {
	var tag pgconn.CommandTag
	var err error

	if tx, ok := getTx(ctx); ok {
		tag, err = tx.Exec(ctx, query, args...)
	} else {
		tag, err = p.pool.Exec(ctx, query, args...)
	}

	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func (p *Client) TransactionManager() repository.TransactionManager {
	return NewTransactionManager(p.pool)
}

func (p *Client) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}

// Check if Postgres implements Database interface
var _ database.Database = (*Client)(nil)

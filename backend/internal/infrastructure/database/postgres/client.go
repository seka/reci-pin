package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database"
)

type pgxExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type Client struct {
	cfg       config.Database
	pool      *pgxpool.Pool
	txManager repository.TransactionManager
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
	p.txManager = NewTransactionManager(pool)
	return nil
}

func (p *Client) executor(ctx context.Context) pgxExecutor {
	if tm, ok := p.txManager.(*transactionManager); ok {
		if tx, ok := tm.getTx(ctx); ok {
			return tx
		}
	}
	return p.pool
}

func (p *Client) Query(ctx context.Context, query string, args ...any) (database.Rows, error) {
	return p.executor(ctx).Query(ctx, query, args...)
}

func (p *Client) Execute(ctx context.Context, query string, args ...any) (int64, error) {
	tag, err := p.executor(ctx).Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func (p *Client) TransactionManager() repository.TransactionManager {
	return p.txManager
}

func (p *Client) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}

// Check if Postgres implements Database interface
var _ database.Database = (*Client)(nil)

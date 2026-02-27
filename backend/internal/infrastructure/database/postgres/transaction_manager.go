package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

type txKey struct{}

type transactionManager struct {
	pool *pgxpool.Pool
}

// NewTransactionManager creates a new repository.TransactionManager implementation for Postgres.
func NewTransactionManager(pool *pgxpool.Pool) repository.TransactionManager {
	return &transactionManager{pool: pool}
}

// WithTransaction executes the given function within a database transaction.
func (tm *transactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(txCtx); err != nil {
		// Use background context for rollback to ensure it completes even if the original context is canceled.
		if rbErr := tx.Rollback(context.Background()); rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

// getTx retrieves the transaction from the context, if any.
func (tm *transactionManager) getTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

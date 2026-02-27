package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

// mockTx is a handwritten mock for pgx.Tx
type mockTx struct {
	pgx.Tx
	committed  bool
	rolledBack bool
}

func (m *mockTx) Commit(ctx context.Context) error {
	m.committed = true
	return nil
}

func (m *mockTx) Rollback(ctx context.Context) error {
	m.rolledBack = true
	return nil
}

// Implement other minimal required methods for pgx.Tx if needed, 
// but for WithTransaction only Commit and Rollback are used on the tx object directly.
func (m *mockTx) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

// mockBeginner is a handwritten mock for txBeginner
type mockBeginner struct {
	tx        *mockTx
	err       error
	BeginFunc func(ctx context.Context) (pgx.Tx, error)
}

func (m *mockBeginner) Begin(ctx context.Context) (pgx.Tx, error) {
	if m.BeginFunc != nil {
		return m.BeginFunc(ctx)
	}
	return m.tx, m.err
}

func TestWithTransaction(t *testing.T) {
	t.Run("Success_Commits", func(t *testing.T) {
		tx := &mockTx{}
		beginner := &mockBeginner{tx: tx}
		tm := &transactionManager{pool: beginner}

		err := tm.WithTransaction(context.Background(), func(ctx context.Context) error {
			return nil
		})

		assert.NoError(t, err)
		assert.True(t, tx.committed)
		assert.False(t, tx.rolledBack)
	})

	t.Run("Failure_Rollbacks", func(t *testing.T) {
		tx := &mockTx{}
		beginner := &mockBeginner{tx: tx}
		tm := &transactionManager{pool: beginner}

		expectedErr := errors.New("business logic error")
		err := tm.WithTransaction(context.Background(), func(ctx context.Context) error {
			return expectedErr
		})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.False(t, tx.committed)
		assert.True(t, tx.rolledBack)
	})

	t.Run("BeginFailure_ReturnsError", func(t *testing.T) {
		expectedErr := errors.New("begin error")
		beginner := &mockBeginner{err: expectedErr}
		tm := &transactionManager{pool: beginner}

		err := tm.WithTransaction(context.Background(), func(ctx context.Context) error {
			return nil
		})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Nested_ReusesTransaction", func(t *testing.T) {
		tx := &mockTx{}
		beginner := &mockBeginner{tx: tx}
		tm := &transactionManager{pool: beginner}

		// Counter to verify Begin is only called once
		beginCount := 0
		beginner.BeginFunc = func(ctx context.Context) (pgx.Tx, error) {
			beginCount++
			return tx, nil
		}

		err := tm.WithTransaction(context.Background(), func(ctx context.Context) error {
			// Second call should reuse the transaction
			return tm.WithTransaction(ctx, func(nestedCtx context.Context) error {
				return nil
			})
		})

		assert.NoError(t, err)
		assert.Equal(t, 1, beginCount, "Begin should be called exactly once")
		assert.True(t, tx.committed)
	})
}

package repository

//go:generate mockgen -source=$GOFILE -destination=mock/transaction_mock.go -package=mock

import "context"

// TransactionManager handles database transactions using a closure pattern.
type TransactionManager interface {
	// WithTransaction executes the given function within a transaction.
	// If the function returns an error, the transaction is rolled back.
	// Otherwise, the transaction is committed.
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

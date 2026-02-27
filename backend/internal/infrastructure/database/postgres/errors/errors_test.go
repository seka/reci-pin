package errors

import (
	"errors"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	domainErrors "github.com/seka/reci-pin/backend/internal/domain/errors"
	"github.com/stretchr/testify/assert"
)

func TestAs(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected error
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: nil,
		},
		{
			name:     "pgx.ErrNoRows",
			err:      pgx.ErrNoRows,
			expected: domainErrors.ErrNotFound,
		},
		{
			name: "unique violation",
			err: &pgconn.PgError{
				Code: pgerrcode.UniqueViolation,
			},
			expected: domainErrors.ErrAlreadyExists,
		},
		{
			name: "foreign key violation",
			err: &pgconn.PgError{
				Code: pgerrcode.ForeignKeyViolation,
			},
			expected: domainErrors.ErrNotFound,
		},
		{
			name: "not null violation",
			err: &pgconn.PgError{
				Code: pgerrcode.NotNullViolation,
			},
			expected: domainErrors.ErrIntegrityConstraint,
		},
		{
			name: "data exception",
			err: &pgconn.PgError{
				Code: pgerrcode.StringDataRightTruncationDataException,
			},
			expected: domainErrors.ErrBadRequest,
		},
		{
			name: "invalid transaction initiation",
			err: &pgconn.PgError{
				Code: pgerrcode.InvalidTransactionInitiation,
			},
			expected: domainErrors.ErrTransaction,
		},
		{
			name: "unclassified pg error",
			err: &pgconn.PgError{
				Code: "99999", // Unknown code
			},
			expected: domainErrors.ErrInternal,
		},
		{
			name:     "other error",
			err:      errors.New("db error"),
			expected: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := As(tt.err)
			if tt.expected == nil {
				assert.NoError(t, actual)
			} else {
				// Use errors.Is for domain errors (sentinels/wrapped)
				if errors.Is(actual, tt.expected) {
					return
				}
				// Fallback to string comparison for generic errors
				assert.Equal(t, tt.expected.Error(), actual.Error(), "expected error message to match")
			}
		})
	}
}

package postgres

import (
	"errors"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/stretchr/testify/assert"
)

func TestMapError(t *testing.T) {
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
			expected: repository.ErrNotFound,
		},
		{
			name: "unique violation",
			err: &pgconn.PgError{
				Code: pgerrcode.UniqueViolation,
			},
			expected: repository.ErrAlreadyExists,
		},
		{
			name: "foreign key violation",
			err: &pgconn.PgError{
				Code: pgerrcode.ForeignKeyViolation,
			},
			expected: repository.ErrNotFound,
		},
		{
			name:     "other error",
			err:      errors.New("db error"),
			expected: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := MapError(tt.err)
			if tt.expected == nil {
				assert.NoError(t, actual)
			} else {
				assert.Equal(t, tt.expected.Error(), actual.Error())
				assert.True(t, errors.Is(actual, tt.expected) || actual.Error() == tt.expected.Error())
			}
		})
	}
}

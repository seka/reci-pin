package postgres

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

// MapError translates PostgreSQL specific errors into domain-level repository errors.
func MapError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return repository.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return repository.ErrAlreadyExists
		case pgerrcode.ForeignKeyViolation:
			return repository.ErrNotFound // Or defined ErrForeignKeyViolation
		}
	}

	return err
}

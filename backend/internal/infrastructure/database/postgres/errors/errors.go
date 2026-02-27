package errors

import (
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	domainErrors "github.com/seka/reci-pin/backend/internal/domain/errors"
)

// As は、PostgreSQL 特有のエラーをドメイン層のリポジトリエラーに変換します。
func As(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return domainErrors.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		code := pgErr.Code

		if pgerrcode.IsCaseNotFound(code) || pgerrcode.IsNoData(code) {
			return domainErrors.ErrNotFound
		}

		if pgerrcode.IsDataException(code) {
			return domainErrors.ErrBadRequest
		}

		if pgerrcode.IsIntegrityConstraintViolation(code) {
			switch code {
			case pgerrcode.UniqueViolation:
				return domainErrors.ErrAlreadyExists
			case pgerrcode.ForeignKeyViolation:
				return domainErrors.ErrNotFound
			default:
				return domainErrors.ErrIntegrityConstraint
			}
		}

		if pgerrcode.IsInvalidTransactionInitiation(code) {
			return domainErrors.ErrTransaction
		}

		// 分類されていない PostgreSQL エラーのフォールバック
		return fmt.Errorf("%w: %v", domainErrors.ErrInternal, err)
	}

	return err
}

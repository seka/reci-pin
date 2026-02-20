package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/seka/reci-pin/backend/internal/infrastructure/database"

	"github.com/seka/reci-pin/backend/internal/domain/model"
)

type PasswordResetTokenRepository struct {
	db database.Database
}

func NewPasswordResetTokenRepository(db database.Database) *PasswordResetTokenRepository {
	return &PasswordResetTokenRepository{db: db}
}

func (r *PasswordResetTokenRepository) Save(ctx context.Context, token string, userID int64, expiresAt time.Time) error {
	query := `
		INSERT INTO password_reset_tokens (token, user_id, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Execute(ctx, query, token, userID, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to save password reset token: %w", err)
	}
	return nil
}

func (r *PasswordResetTokenRepository) GetByToken(ctx context.Context, token string) (*model.PasswordResetToken, error) {
	query := `
		SELECT token, user_id, expires_at, created_at
		FROM password_reset_tokens
		WHERE token = $1
	`
	var t model.PasswordResetToken
	rows, err := r.db.Query(ctx, query, token)
	if err != nil {
		return nil, fmt.Errorf("failed to find password reset token: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil // Not found
	}

	err = rows.Scan(&t.Token, &t.UserID, &t.ExpiresAt, &t.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to scan password reset token: %w", err)
	}
	return &t, nil
}

func (r *PasswordResetTokenRepository) Delete(ctx context.Context, token string) error {
	query := `DELETE FROM password_reset_tokens WHERE token = $1`
	_, err := r.db.Execute(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to delete password reset token: %w", err)
	}
	return nil
}

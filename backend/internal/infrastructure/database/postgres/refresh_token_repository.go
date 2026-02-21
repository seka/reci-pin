package postgres

import (
	"context"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database"
	"github.com/seka/reci-pin/backend/internal/infrastructure/entity"
)

type RefreshTokenRepository struct {
	db database.Database
}

func NewRefreshTokenRepository(db database.Database) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Save(ctx context.Context, token *model.RefreshToken) error {
	e := entity.NewRefreshToken(token)
	query := `
		INSERT INTO user_refresh_tokens (user_id, token_hash, expires_at, created_at, revoked_at, user_agent, ip_address)
		VALUES ($1, $2, $3, NOW(), $4, $5, $6)
		RETURNING id, created_at
	`

	rows, err := r.db.Query(ctx, query, e.UserID, e.TokenHash, e.ExpiresAt, e.RevokedAt, e.UserAgent, e.IPAddress)
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return fmt.Errorf("failed to save refresh token: no rows returned")
	}

	err = rows.Scan(&e.ID, &e.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to scan refresh token: %w", err)
	}

	token.ID = e.ID
	token.CreatedAt = e.CreatedAt
	return nil
}

func (r *RefreshTokenRepository) GetByHash(ctx context.Context, hash string) (*model.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at, revoked_at, user_agent, ip_address
		FROM user_refresh_tokens
		WHERE token_hash = $1
	`
	var e entity.RefreshToken

	rows, err := r.db.Query(ctx, query, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token by hash: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil // Not found
	}

	err = rows.Scan(
		&e.ID,
		&e.UserID,
		&e.TokenHash,
		&e.ExpiresAt,
		&e.CreatedAt,
		&e.RevokedAt,
		&e.UserAgent,
		&e.IPAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan refresh token: %w", err)
	}

	return e.ToModel(), nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, id int64) error {
	query := `UPDATE user_refresh_tokens SET revoked_at = NOW() WHERE id = $1`
	_, err := r.db.Execute(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID int64) error {
	query := `UPDATE user_refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`
	_, err := r.db.Execute(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke all refresh tokens for user: %w", err)
	}
	return nil
}

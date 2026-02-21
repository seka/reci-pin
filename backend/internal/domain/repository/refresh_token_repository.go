package repository

//go:generate mockgen -source=$GOFILE -destination=mock/refresh_token_repository_mock.go -package=mock

import (
	"context"

	"github.com/seka/reci-pin/backend/internal/domain/model"
)

type RefreshTokenRepository interface {
	Save(ctx context.Context, token *model.RefreshToken) error
	GetByHash(ctx context.Context, hash string) (*model.RefreshToken, error)
	Revoke(ctx context.Context, id int64) error
	RevokeAllByUserID(ctx context.Context, userID int64) error
}

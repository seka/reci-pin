package repository

import (
	"context"
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/model"
)

type PasswordResetTokenRepository interface {
	Save(ctx context.Context, token string, userID int64, expiresAt time.Time) error
	Find(ctx context.Context, token string) (*model.PasswordResetToken, error)
	Delete(ctx context.Context, token string) error
}

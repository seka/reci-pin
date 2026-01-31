package repository

import (
	"context"

	"github.com/seka/reci-pin/backend/internal/domain/model"
)

type UserEmailCredentialRepository interface {
	Create(ctx context.Context, credential *model.UserEmailCredential) error
	GetByEmail(ctx context.Context, email string) (*model.UserEmailCredential, error)
	GetByUserID(ctx context.Context, userID int64) (*model.UserEmailCredential, error)
	GetByToken(ctx context.Context, token string) (*model.UserEmailCredential, error)
	Update(ctx context.Context, credential *model.UserEmailCredential) error
}

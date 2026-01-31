package repository

import (
	"context"

	"github.com/seka/reci-pin/backend/internal/domain/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User, passwordHash string) error
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, string, error) // returns user and passwordHash
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
}

package repository

import (
	"context"

	"github.com/seka/reci-pin/backend/internal/domain/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetAll(ctx context.Context) ([]*model.User, error)
	// Update, Delete will be added if needed for profile management
}

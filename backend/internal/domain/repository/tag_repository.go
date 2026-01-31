package repository

import (
	"context"

	"github.com/seka/reci-pin/backend/internal/domain/model"
)

type TagRepository interface {
	Create(ctx context.Context, tag *model.Tag) error
	GetByID(ctx context.Context, id int64) (*model.Tag, error)
	GetByName(ctx context.Context, name string) (*model.Tag, error)
	GetAll(ctx context.Context) ([]model.Tag, error)
	Delete(ctx context.Context, id int64) error
}

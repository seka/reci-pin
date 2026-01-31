package repository

import (
	"context"

	"github.com/seka/reci-pin/backend/internal/domain/entity"
)

type TagRepository interface {
	Create(ctx context.Context, tag *entity.Tag) error
	GetByID(ctx context.Context, id int64) (*entity.Tag, error)
	GetByName(ctx context.Context, name string) (*entity.Tag, error)
	GetAll(ctx context.Context) ([]entity.Tag, error)
	Delete(ctx context.Context, id int64) error
}

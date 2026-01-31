package repository

import (
	"context"

	"github.com/seka/reci-pin/backend/internal/domain/entity"
)

type RecipeRepository interface {
	Create(ctx context.Context, recipe *entity.Recipe) error
	GetByID(ctx context.Context, id int64) (*entity.Recipe, error)
	GetByUserID(ctx context.Context, userID int64) ([]entity.Recipe, error)
	Search(ctx context.Context, userID int64, query string, tagIDs []int64) ([]entity.Recipe, error)
	Update(ctx context.Context, recipe *entity.Recipe) error
	Delete(ctx context.Context, id int64) error
	AddTags(ctx context.Context, recipeID int64, tagIDs []int64) error
	RemoveTags(ctx context.Context, recipeID int64, tagIDs []int64) error
	GetTags(ctx context.Context, recipeID int64) ([]entity.Tag, error)
}

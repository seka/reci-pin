package repository

import (
	"context"

	"github.com/seka/reci-pin/backend/internal/domain/model"
)

//go:generate mockgen -source=$GOFILE -destination=./mock/recipe_repository_mock.go -package=mock
type RecipeRepository interface {
	Create(ctx context.Context, recipe *model.Recipe) error
	GetByID(ctx context.Context, id int64) (*model.Recipe, error)
	GetByIDs(ctx context.Context, ids []int64) ([]model.Recipe, error)
	GetByUserID(ctx context.Context, userID int64) ([]model.Recipe, error)
	GetAll(ctx context.Context) ([]model.Recipe, error)
	Update(ctx context.Context, recipe *model.Recipe) error
	Delete(ctx context.Context, id int64) error
	Search(ctx context.Context, userID int64, query string, tagIDs []int64) ([]model.Recipe, error)

	// Tag operations
	GetTags(ctx context.Context, recipeID int64) ([]model.Tag, error)
	GetTagsBatch(ctx context.Context, recipeIDs []int64) (map[int64][]model.Tag, error)
	AddTags(ctx context.Context, recipeID int64, tagIDs []int64) error
	RemoveTags(ctx context.Context, recipeID int64, tagIDs []int64) error
}

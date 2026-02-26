package repository

import (
	"context"

	"github.com/seka/reci-pin/backend/internal/domain/model"
)

//go:generate mockgen -source=$GOFILE -destination=./mock/recipe_image_repository_mock.go -package=mock
type RecipeImageRepository interface {
	Create(ctx context.Context, image *model.RecipeImage) error
	GetByRecipeID(ctx context.Context, recipeID int64) ([]model.RecipeImage, error)
	GetByRecipeIDs(ctx context.Context, recipeIDs []int64) (map[int64][]model.RecipeImage, error)
	Delete(ctx context.Context, id int64) error
}

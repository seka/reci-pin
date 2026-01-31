package repository

import (
	"context"

	"github.com/seka/reci-pin/backend/internal/domain/entity"
)

type RecipeImageRepository interface {
	Create(ctx context.Context, image *entity.RecipeImage) error
	GetByRecipeID(ctx context.Context, recipeID int64) ([]entity.RecipeImage, error)
	Delete(ctx context.Context, id int64) error
}

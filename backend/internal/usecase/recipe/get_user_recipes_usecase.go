package recipe

import (
	"context"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/entity"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

type GetUserRecipesUseCase struct {
	recipeRepo      repository.RecipeRepository
	recipeImageRepo repository.RecipeImageRepository
}

func NewGetUserRecipesUseCase(
	recipeRepo repository.RecipeRepository,
	recipeImageRepo repository.RecipeImageRepository,
) *GetUserRecipesUseCase {
	return &GetUserRecipesUseCase{
		recipeRepo:      recipeRepo,
		recipeImageRepo: recipeImageRepo,
	}
}

func (uc *GetUserRecipesUseCase) Execute(ctx context.Context, userID int64) ([]entity.Recipe, error) {
	recipes, err := uc.recipeRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user recipes: %w", err)
	}

	// Load tags and images for each recipe
	for i := range recipes {
		tags, err := uc.recipeRepo.GetTags(ctx, recipes[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get recipe tags: %w", err)
		}
		recipes[i].Tags = tags

		images, err := uc.recipeImageRepo.GetByRecipeID(ctx, recipes[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get recipe images: %w", err)
		}
		recipes[i].Images = images
	}

	return recipes, nil
}

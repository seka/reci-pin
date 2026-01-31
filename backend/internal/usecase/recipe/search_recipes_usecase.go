package recipe

import (
	"context"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

type SearchRecipesUseCase struct {
	recipeRepo      repository.RecipeRepository
	recipeImageRepo repository.RecipeImageRepository
}

func NewSearchRecipesUseCase(
	recipeRepo repository.RecipeRepository,
	recipeImageRepo repository.RecipeImageRepository,
) *SearchRecipesUseCase {
	return &SearchRecipesUseCase{
		recipeRepo:      recipeRepo,
		recipeImageRepo: recipeImageRepo,
	}
}

type SearchRecipesInput struct {
	UserID int64
	Query  string
	TagIDs []int64
}

func (uc *SearchRecipesUseCase) Execute(ctx context.Context, input SearchRecipesInput) ([]model.Recipe, error) {
	recipes, err := uc.recipeRepo.Search(ctx, input.UserID, input.Query, input.TagIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to search recipes: %w", err)
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

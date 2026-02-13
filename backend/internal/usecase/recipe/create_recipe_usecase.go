package recipe

import (
	"context"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/domain/validation"
)

type CreateRecipeUseCase interface {
	Execute(ctx context.Context, input CreateRecipeInput) (*model.Recipe, error)
}

type createRecipeInteractor struct {
	recipeRepo repository.RecipeRepository
}

func NewCreateRecipeUseCase(recipeRepo repository.RecipeRepository) CreateRecipeUseCase {
	return &createRecipeInteractor{recipeRepo: recipeRepo}
}

type CreateRecipeInput struct {
	UserID int64
	Name   string
	URL    string
	Memo   string
	TagIDs []int64
}

func (uc *createRecipeInteractor) Execute(ctx context.Context, input CreateRecipeInput) (*model.Recipe, error) {
	if err := validation.ValidateRecipe(input.Name, input.URL); err != nil {
		return nil, err
	}

	recipe := &model.Recipe{
		UserID: input.UserID,
		Name:   input.Name,
		URL:    input.URL,
		Memo:   input.Memo,
	}

	if err := uc.recipeRepo.Create(ctx, recipe); err != nil {
		return nil, fmt.Errorf("failed to create recipe: %w", err)
	}

	// Add tags if provided
	if len(input.TagIDs) > 0 {
		if err := uc.recipeRepo.AddTags(ctx, recipe.ID, input.TagIDs); err != nil {
			return nil, fmt.Errorf("failed to add tags to recipe: %w", err)
		}
	}

	return recipe, nil
}

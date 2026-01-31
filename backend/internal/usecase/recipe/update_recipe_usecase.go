package recipe

import (
	"context"
	"errors"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

type UpdateRecipeUseCase struct {
	recipeRepo repository.RecipeRepository
}

func NewUpdateRecipeUseCase(recipeRepo repository.RecipeRepository) *UpdateRecipeUseCase {
	return &UpdateRecipeUseCase{recipeRepo: recipeRepo}
}

type UpdateRecipeInput struct {
	ID     int64
	UserID int64
	Name   string
	URL    string
	Memo   string
}

func (uc *UpdateRecipeUseCase) Execute(ctx context.Context, input UpdateRecipeInput) (*model.Recipe, error) {
	// Get existing recipe to verify ownership
	recipe, err := uc.recipeRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe: %w", err)
	}

	if recipe.UserID != input.UserID {
		return nil, errors.New("unauthorized access to recipe")
	}

	// Update recipe
	recipe.Name = input.Name
	recipe.URL = input.URL
	recipe.Memo = input.Memo

	if err := uc.recipeRepo.Update(ctx, recipe); err != nil {
		return nil, fmt.Errorf("failed to update recipe: %w", err)
	}

	return recipe, nil
}

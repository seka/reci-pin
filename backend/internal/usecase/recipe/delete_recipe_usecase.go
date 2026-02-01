package recipe

import (
	"context"
	"errors"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

type DeleteRecipeUseCase interface {
	Execute(ctx context.Context, id, userID int64) error
}

type deleteRecipeInteractor struct {
	recipeRepo repository.RecipeRepository
}

func NewDeleteRecipeUseCase(recipeRepo repository.RecipeRepository) DeleteRecipeUseCase {
	return &deleteRecipeInteractor{recipeRepo: recipeRepo}
}

func (uc *deleteRecipeInteractor) Execute(ctx context.Context, id, userID int64) error {
	// Get existing recipe to verify ownership
	recipe, err := uc.recipeRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get recipe: %w", err)
	}

	if recipe.UserID != userID {
		return errors.New("unauthorized access to recipe")
	}

	if err := uc.recipeRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete recipe: %w", err)
	}

	return nil
}

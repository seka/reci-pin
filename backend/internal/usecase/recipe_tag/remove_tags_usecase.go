package recipe_tag

import (
	"context"
	"errors"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

type RemoveTagsUseCase interface {
	Execute(ctx context.Context, recipeID, userID int64, tagIDs []int64) error
}

type removeTagsInteractor struct {
	recipeRepo repository.RecipeRepository
}

func NewRemoveTagsUseCase(recipeRepo repository.RecipeRepository) RemoveTagsUseCase {
	return &removeTagsInteractor{recipeRepo: recipeRepo}
}

func (uc *removeTagsInteractor) Execute(ctx context.Context, recipeID, userID int64, tagIDs []int64) error {
	// Verify ownership
	recipe, err := uc.recipeRepo.GetByID(ctx, recipeID)
	if err != nil {
		return fmt.Errorf("failed to get recipe: %w", err)
	}

	if recipe.UserID != userID {
		return errors.New("unauthorized access to recipe")
	}

	return uc.recipeRepo.RemoveTags(ctx, recipeID, tagIDs)
}

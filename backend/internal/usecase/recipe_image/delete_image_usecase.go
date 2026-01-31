package recipe_image

import (
	"context"
	"errors"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

type DeleteImageUseCase struct {
	recipeRepo      repository.RecipeRepository
	recipeImageRepo repository.RecipeImageRepository
}

func NewDeleteImageUseCase(
	recipeRepo repository.RecipeRepository,
	recipeImageRepo repository.RecipeImageRepository,
) *DeleteImageUseCase {
	return &DeleteImageUseCase{
		recipeRepo:      recipeRepo,
		recipeImageRepo: recipeImageRepo,
	}
}

func (uc *DeleteImageUseCase) Execute(ctx context.Context, imageID, userID int64) error {
	// Get image to find recipe
	images, err := uc.recipeImageRepo.GetByRecipeID(ctx, imageID)
	if err != nil || len(images) == 0 {
		return fmt.Errorf("image not found")
	}

	// Verify ownership through recipe
	recipe, err := uc.recipeRepo.GetByID(ctx, images[0].RecipeID)
	if err != nil {
		return fmt.Errorf("failed to get recipe: %w", err)
	}

	if recipe.UserID != userID {
		return errors.New("unauthorized access to recipe")
	}

	return uc.recipeImageRepo.Delete(ctx, imageID)
}

package recipe_image

import (
	"context"
	"errors"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/entity"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

type AddImageUseCase struct {
	recipeRepo      repository.RecipeRepository
	recipeImageRepo repository.RecipeImageRepository
}

func NewAddImageUseCase(
	recipeRepo repository.RecipeRepository,
	recipeImageRepo repository.RecipeImageRepository,
) *AddImageUseCase {
	return &AddImageUseCase{
		recipeRepo:      recipeRepo,
		recipeImageRepo: recipeImageRepo,
	}
}

func (uc *AddImageUseCase) Execute(ctx context.Context, recipeID, userID int64, imagePath string) (*entity.RecipeImage, error) {
	// Verify ownership
	recipe, err := uc.recipeRepo.GetByID(ctx, recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe: %w", err)
	}

	if recipe.UserID != userID {
		return nil, errors.New("unauthorized access to recipe")
	}

	image := &entity.RecipeImage{
		RecipeID:  recipeID,
		ImagePath: imagePath,
	}

	if err := uc.recipeImageRepo.Create(ctx, image); err != nil {
		return nil, fmt.Errorf("failed to add image to recipe: %w", err)
	}

	return image, nil
}

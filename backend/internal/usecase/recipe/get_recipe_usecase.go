package recipe

import (
	"context"
	"errors"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

type GetRecipeUseCase interface {
	Execute(ctx context.Context, id, userID int64) (*model.Recipe, error)
}

type getRecipeInteractor struct {
	recipeRepo      repository.RecipeRepository
	recipeImageRepo repository.RecipeImageRepository
}

func NewGetRecipeUseCase(
	recipeRepo repository.RecipeRepository,
	recipeImageRepo repository.RecipeImageRepository,
) GetRecipeUseCase {
	return &getRecipeInteractor{
		recipeRepo:      recipeRepo,
		recipeImageRepo: recipeImageRepo,
	}
}

func (uc *getRecipeInteractor) Execute(ctx context.Context, id, userID int64) (*model.Recipe, error) {
	recipe, err := uc.recipeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe: %w", err)
	}

	// Verify ownership
	if recipe.UserID != userID {
		return nil, errors.New("unauthorized access to recipe")
	}

	// Load tags
	tags, err := uc.recipeRepo.GetTags(ctx, recipe.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe tags: %w", err)
	}
	recipe.Tags = tags

	// Load images
	images, err := uc.recipeImageRepo.GetByRecipeID(ctx, recipe.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe images: %w", err)
	}
	recipe.Images = images

	return recipe, nil
}

package recipe

import (
	"context"
	"errors"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/domain/validation"
)

type UpdateRecipeUseCase interface {
	Execute(ctx context.Context, input UpdateRecipeInput) (*model.Recipe, error)
}

type updateRecipeInteractor struct {
	recipeRepo repository.RecipeRepository
	searchRepo repository.RecipeSearchRepository
}

func NewUpdateRecipeUseCase(
	recipeRepo repository.RecipeRepository,
	searchRepo repository.RecipeSearchRepository,
) UpdateRecipeUseCase {
	return &updateRecipeInteractor{
		recipeRepo: recipeRepo,
		searchRepo: searchRepo,
	}
}

type UpdateRecipeInput struct {
	ID     int64
	UserID int64
	Name   string
	URL    string
	Memo   string
}

func (uc *updateRecipeInteractor) Execute(ctx context.Context, input UpdateRecipeInput) (*model.Recipe, error) {
	// Validation
	if err := validation.ValidateRecipe(input.Name, input.URL); err != nil {
		return nil, err
	}

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

	// Fetch tags for indexing
	tags, err := uc.recipeRepo.GetTags(ctx, recipe.ID)
	if err != nil {
		// Log warning but proceed with indexing without tags if failed
		fmt.Printf("failed to get tags for indexing: %v\n", err)
	} else {
		recipe.Tags = tags
	}

	// Update index
	if err := uc.searchRepo.Index(ctx, recipe); err != nil {
		fmt.Printf("failed to update recipe index: %v\n", err)
	}

	return recipe, nil
}

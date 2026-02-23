package recipe

import (
	"context"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/domain/searcher"
	"github.com/seka/reci-pin/backend/internal/domain/validation"
)

type CreateRecipeUseCase interface {
	Execute(ctx context.Context, input CreateRecipeInput) (*model.Recipe, error)
}

type createRecipeInteractor struct {
	recipeRepo repository.RecipeRepository
	searchRepo searcher.RecipeSearcher
}

func NewCreateRecipeUseCase(
	recipeRepo repository.RecipeRepository,
	searchRepo searcher.RecipeSearcher,
) CreateRecipeUseCase {
	return &createRecipeInteractor{
		recipeRepo: recipeRepo,
		searchRepo: searchRepo,
	}
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
		// Set tags to recipe for indexing (only ID is needed for current implementation)
		for _, tagID := range input.TagIDs {
			recipe.Tags = append(recipe.Tags, model.Tag{ID: tagID})
		}
	}

	// Index recipe for search
	if err := uc.searchRepo.Index(ctx, recipe); err != nil {
		// Log error but do not fail the request (dual write best effort)
		fmt.Printf("failed to index recipe: %v\n", err)
	}

	return recipe, nil
}

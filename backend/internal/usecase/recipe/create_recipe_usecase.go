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
	txManager  repository.TransactionManager
}

func NewCreateRecipeUseCase(
	recipeRepo repository.RecipeRepository,
	searchRepo searcher.RecipeSearcher,
	txManager repository.TransactionManager,
) CreateRecipeUseCase {
	return &createRecipeInteractor{
		recipeRepo: recipeRepo,
		searchRepo: searchRepo,
		txManager:  txManager,
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

	err := uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.recipeRepo.Create(txCtx, recipe); err != nil {
			return fmt.Errorf("failed to create recipe: %w", err)
		}

		// Add tags if provided
		if len(input.TagIDs) > 0 {
			if err := uc.recipeRepo.AddTags(txCtx, recipe.ID, input.TagIDs); err != nil {
				return fmt.Errorf("failed to add tags to recipe: %w", err)
			}
			// Set tags to recipe for indexing (only ID is needed for current implementation)
			for _, tagID := range input.TagIDs {
				recipe.Tags = append(recipe.Tags, model.Tag{ID: tagID})
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Index recipe for search
	if err := uc.searchRepo.Index(ctx, recipe); err != nil {
		// Log error but do not fail the request (dual write best effort)
		fmt.Printf("failed to index recipe: %v\n", err)
	}

	return recipe, nil
}

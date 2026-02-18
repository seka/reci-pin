package recipe

import (
	"context"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

type SearchRecipesUseCase interface {
	Execute(ctx context.Context, input SearchRecipesInput) ([]model.Recipe, error)
}

type searchRecipesInteractor struct {
	recipeRepo      repository.RecipeRepository
	recipeImageRepo repository.RecipeImageRepository
	searchRepo      repository.RecipeSearchRepository
}

func NewSearchRecipesUseCase(
	recipeRepo repository.RecipeRepository,
	recipeImageRepo repository.RecipeImageRepository,
	searchRepo repository.RecipeSearchRepository,
) SearchRecipesUseCase {
	return &searchRecipesInteractor{
		recipeRepo:      recipeRepo,
		recipeImageRepo: recipeImageRepo,
		searchRepo:      searchRepo,
	}
}

type SearchRecipesInput struct {
	UserID int64
	Query  string
	TagIDs []int64
}

func (uc *searchRecipesInteractor) Execute(ctx context.Context, input SearchRecipesInput) ([]model.Recipe, error) {
	// Search in ElasticSearch
	criteria := repository.SearchCriteria{
		UserID:   input.UserID,
		Keyword:  input.Query,
		TagIDs:   input.TagIDs,
		Page:     1,  // Default to page 1
		PageSize: 50, // Default distinct fetch size
	}

	ids, _, err := uc.searchRepo.Search(ctx, criteria)
	if err != nil {
		// Fallback to DB search if ES fails? Or return error.
		// For now, return error to make sure ES is working.
		// Alternatively, we could log and fallback to recipeRepo.Search(ctx, input.UserID, input.Query, input.TagIDs)
		// But recipeRepo.Search might be deprecated or removed.
		// Let's fallback for robustness during migration.
		fmt.Printf("ES search failed, falling back to DB: %v\n", err)
		return uc.recipeRepo.Search(ctx, input.UserID, input.Query, input.TagIDs)
	}

	if len(ids) == 0 {
		return []model.Recipe{}, nil
	}

	// Fetch details from DB
	// TODO: Add GetByIDs to repository for better performance
	recipes := make([]model.Recipe, 0, len(ids))
	for _, id := range ids {
		recipe, err := uc.recipeRepo.GetByID(ctx, id)
		if err != nil {
			// Start logging warning. Maybe the recipe was deleted in DB but not in ES?
			// Skip this recipe.
			fmt.Printf("failed to get recipe details for id %d: %v\n", id, err)
			continue
		}
		recipes = append(recipes, *recipe)
	}

	// Load tags and images for each recipe
	for i := range recipes {
		tags, err := uc.recipeRepo.GetTags(ctx, recipes[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get recipe tags: %w", err)
		}
		recipes[i].Tags = tags

		images, err := uc.recipeImageRepo.GetByRecipeID(ctx, recipes[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get recipe images: %w", err)
		}
		recipes[i].Images = images
	}

	return recipes, nil
}

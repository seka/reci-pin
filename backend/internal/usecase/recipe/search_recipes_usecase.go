package recipe

import (
	"context"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/domain/storage"
)

type SearchRecipesUseCase interface {
	Execute(ctx context.Context, input SearchRecipesInput) ([]model.Recipe, error)
}

type searchRecipesInteractor struct {
	recipeRepo      repository.RecipeRepository
	recipeImageRepo repository.RecipeImageRepository
	searchRepo      repository.RecipeSearchRepository
	storageService  storage.Client
}

func NewSearchRecipesUseCase(
	recipeRepo repository.RecipeRepository,
	recipeImageRepo repository.RecipeImageRepository,
	searchRepo repository.RecipeSearchRepository,
	storageService storage.Client,
) SearchRecipesUseCase {
	return &searchRecipesInteractor{
		recipeRepo:      recipeRepo,
		recipeImageRepo: recipeImageRepo,
		searchRepo:      searchRepo,
		storageService:  storageService,
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
		fmt.Printf("ES search failed, falling back to DB: %v\n", err)
		// Search in DB
		recipes, err := uc.recipeRepo.Search(ctx, input.UserID, input.Query, input.TagIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to search recipes in DB: %w", err)
		}
		// Extract IDs from DB results to use common loading logic
		ids = make([]int64, len(recipes))
		for i, r := range recipes {
			ids[i] = r.ID
		}
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
			return nil, fmt.Errorf("failed to get images for recipe %d: %w", recipes[i].ID, err)
		}

		// Orchestrate: Convert to PublicRecipeImage
		baseURL := uc.storageService.GetPublicURL()
		publicImages := make([]model.PublicRecipeImage, len(images))
		for j, img := range images {
			u := baseURL.JoinPath(img.ImagePath)
			publicImages[j] = model.PublicRecipeImage{
				RecipeImage: img,
				ImageURL:    *u,
			}
		}
		recipes[i].Images = publicImages
	}

	return recipes, nil
}

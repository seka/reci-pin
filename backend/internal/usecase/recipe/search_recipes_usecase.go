package recipe

import (
	"context"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/domain/searcher"
	"github.com/seka/reci-pin/backend/internal/domain/storage"
)

type SearchRecipesUseCase interface {
	Execute(ctx context.Context, input SearchRecipesInput) ([]model.Recipe, error)
}

type searchRecipesInteractor struct {
	recipeRepo      repository.RecipeRepository
	recipeImageRepo repository.RecipeImageRepository
	searchRepo      searcher.RecipeSearcher
	storageService  storage.Client
}

func NewSearchRecipesUseCase(
	recipeRepo repository.RecipeRepository,
	recipeImageRepo repository.RecipeImageRepository,
	recipeSearcher searcher.RecipeSearcher,
	storageService storage.Client,
) SearchRecipesUseCase {
	return &searchRecipesInteractor{
		recipeRepo:      recipeRepo,
		recipeImageRepo: recipeImageRepo,
		searchRepo:      recipeSearcher,
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
	criteria := model.RecipeSearchCriteria{
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

	// 2. Batch fetch recipes by IDs
	recipes, err := uc.recipeRepo.BulkGetByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipes: %w", err)
	}

	// 3. Batch fetch tags and images for the retrieved recipes
	tagsMap, err := uc.recipeRepo.BulkGetTags(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	imagesMap, err := uc.recipeImageRepo.BulkGetByRecipeIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe images batch: %w", err)
	}

	baseURL := uc.storageService.GetPublicURL()

	// Orchestrate: Assemble tags and public images for each recipe
	for i := range recipes {
		id := recipes[i].ID
		
		// Set tags
		if tags, ok := tagsMap[id]; ok {
			recipes[i].Tags = tags
		} else {
			recipes[i].Tags = []model.Tag{}
		}

		// Set public images
		if images, ok := imagesMap[id]; ok {
			publicImages := make([]model.PublicRecipeImage, len(images))
			for j, img := range images {
				u := baseURL.JoinPath(img.ImagePath)
				publicImages[j] = model.PublicRecipeImage{
					RecipeImage: img,
					ImageURL:    *u,
				}
			}
			recipes[i].Images = publicImages
		} else {
			recipes[i].Images = []model.PublicRecipeImage{}
		}
	}

	return recipes, nil
}

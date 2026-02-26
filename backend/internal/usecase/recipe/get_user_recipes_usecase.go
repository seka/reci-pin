package recipe

import (
	"context"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/domain/storage"
)

type GetUserRecipesUseCase interface {
	Execute(ctx context.Context, userID int64) ([]model.Recipe, error)
}

type getUserRecipesInteractor struct {
	recipeRepo      repository.RecipeRepository
	recipeImageRepo repository.RecipeImageRepository
	storageService  storage.Client
}

func NewGetUserRecipesUseCase(
	recipeRepo repository.RecipeRepository,
	recipeImageRepo repository.RecipeImageRepository,
	storageService storage.Client,
) GetUserRecipesUseCase {
	return &getUserRecipesInteractor{
		recipeRepo:      recipeRepo,
		recipeImageRepo: recipeImageRepo,
		storageService:  storageService,
	}
}

func (uc *getUserRecipesInteractor) Execute(ctx context.Context, userID int64) ([]model.Recipe, error) {
	recipes, err := uc.recipeRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user recipes: %w", err)
	}

	// Load tags and images in batch
	recipeIDs := make([]int64, len(recipes))
	for i, r := range recipes {
		recipeIDs[i] = r.ID
	}

	tagsMap, err := uc.recipeRepo.GetTagsBatch(ctx, recipeIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe tags batch: %w", err)
	}

	imagesMap, err := uc.recipeImageRepo.GetByRecipeIDs(ctx, recipeIDs)
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

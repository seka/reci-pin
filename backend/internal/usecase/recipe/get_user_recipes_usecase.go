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

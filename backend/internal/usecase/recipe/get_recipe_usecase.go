package recipe

import (
	"context"
	"errors"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/domain/storage"
)

type GetRecipeUseCase interface {
	Execute(ctx context.Context, id, userID int64) (*model.Recipe, error)
}

type getRecipeInteractor struct {
	recipeRepo      repository.RecipeRepository
	recipeImageRepo repository.RecipeImageRepository
	storageService  storage.Client
}

func NewGetRecipeUseCase(
	recipeRepo repository.RecipeRepository,
	recipeImageRepo repository.RecipeImageRepository,
	storageService storage.Client,
) GetRecipeUseCase {
	return &getRecipeInteractor{
		recipeRepo:      recipeRepo,
		recipeImageRepo: recipeImageRepo,
		storageService:  storageService,
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

	// Orchestrate: Convert to PublicRecipeImage
	baseURL := uc.storageService.GetPublicURL()
	publicImages := make([]model.PublicRecipeImage, len(images))
	for i, img := range images {
		u := baseURL.JoinPath(img.ImagePath)
		publicImages[i] = model.PublicRecipeImage{
			RecipeImage: img,
			ImageURL:    *u,
		}
	}
	recipe.Images = publicImages

	return recipe, nil
}

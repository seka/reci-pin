package recipe_image

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/domain/storage"
	"github.com/seka/reci-pin/backend/internal/domain/validation"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/create_recipe_image_usecase_mock.go -package=mock
type CreateRecipeImageUseCase interface {
	Execute(ctx context.Context, input CreateRecipeImageInput) (*model.RecipeImage, string, error)
}

type CreateRecipeImageInput struct {
	RecipeID    int64
	UserID      int64
	Filename    string
	ContentType string
	Size        int64
}

type createRecipeImageInteractor struct {
	recipeRepo      repository.RecipeRepository
	recipeImageRepo repository.RecipeImageRepository
	storageService  storage.Storage
}

func NewCreateRecipeImageUseCase(
	recipeRepo repository.RecipeRepository,
	recipeImageRepo repository.RecipeImageRepository,
	storageService storage.Storage,
) CreateRecipeImageUseCase {
	return &createRecipeImageInteractor{
		recipeRepo:      recipeRepo,
		recipeImageRepo: recipeImageRepo,
		storageService:  storageService,
	}
}

func (uc *createRecipeImageInteractor) Execute(ctx context.Context, input CreateRecipeImageInput) (*model.RecipeImage, string, error) {
	recipe, err := uc.recipeRepo.GetByID(ctx, input.RecipeID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get recipe: %w", err)
	}

	if recipe.UserID != input.UserID {
		return nil, "", errors.New("unauthorized access to recipe")
	}

	if input.Size > validation.ImageMaxFileSize {
		return nil, "", fmt.Errorf("file size exceeds limit of %d bytes", validation.ImageMaxFileSize)
	}

	if !validation.ImageAllowedContentTypes[input.ContentType] {
		return nil, "", fmt.Errorf("unsupported content type: %s", input.ContentType)
	}

	ext := strings.ToLower(filepath.Ext(input.Filename))
	if !validation.ImageAllowedExtensions[ext] {
		return nil, "", fmt.Errorf("unsupported file extension: %s", ext)
	}

	timestamp := time.Now().UnixNano()
	key := fmt.Sprintf("recipes/%d/%d_%s", input.RecipeID, timestamp, input.Filename)

	url, err := uc.storageService.GeneratePresignedURL(ctx, key, input.ContentType, input.Size, 15*time.Minute)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	image := &model.RecipeImage{
		RecipeID:  input.RecipeID,
		ImagePath: key,
	}

	if err := uc.recipeImageRepo.Create(ctx, image); err != nil {
		return nil, "", fmt.Errorf("failed to add image to recipe: %w", err)
	}

	return image, url, nil
}

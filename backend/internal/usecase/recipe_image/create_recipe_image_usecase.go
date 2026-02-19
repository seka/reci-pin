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
	// 1. Verify ownership
	recipe, err := uc.recipeRepo.GetByID(ctx, input.RecipeID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get recipe: %w", err)
	}

	if recipe.UserID != input.UserID {
		return nil, "", errors.New("unauthorized access to recipe")
	}

	// 最近のスマホ(ProRAW/高画素モード等)では10MBを超える画像が生成されるため、余裕を持たせて50MBに設定
	const maxFileSize = 50 * 1024 * 1024 // 50MB
	if input.Size > maxFileSize {
		return nil, "", fmt.Errorf("file size exceeds limit of %d bytes", maxFileSize)
	}

	// 3. Validate ContentType/Extension
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	if !allowedTypes[input.ContentType] {
		return nil, "", fmt.Errorf("unsupported content type: %s", input.ContentType)
	}

	ext := strings.ToLower(filepath.Ext(input.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		return nil, "", fmt.Errorf("unsupported file extension: %s", ext)
	}

	// 4. Generate S3 Key
	timestamp := time.Now().UnixNano()
	key := fmt.Sprintf("recipes/%d/%d_%s", input.RecipeID, timestamp, input.Filename)

	// 5. Generate Presigned URL
	url, err := uc.storageService.GeneratePresignedURL(ctx, key, input.ContentType, input.Size, 15*time.Minute)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	// 6. Create RecipeImage record
	image := &model.RecipeImage{
		RecipeID:  input.RecipeID,
		ImagePath: key,
	}

	if err := uc.recipeImageRepo.Create(ctx, image); err != nil {
		return nil, "", fmt.Errorf("failed to add image to recipe: %w", err)
	}

	return image, url, nil
}

package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/entity"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

type RecipeUseCase struct {
	recipeRepo      repository.RecipeRepository
	tagRepo         repository.TagRepository
	recipeImageRepo repository.RecipeImageRepository
}

func NewRecipeUseCase(
	recipeRepo repository.RecipeRepository,
	tagRepo repository.TagRepository,
	recipeImageRepo repository.RecipeImageRepository,
) *RecipeUseCase {
	return &RecipeUseCase{
		recipeRepo:      recipeRepo,
		tagRepo:         tagRepo,
		recipeImageRepo: recipeImageRepo,
	}
}

type CreateRecipeInput struct {
	UserID int64
	Name   string
	URL    string
	Memo   string
	TagIDs []int64
}

type UpdateRecipeInput struct {
	ID     int64
	UserID int64
	Name   string
	URL    string
	Memo   string
}

type SearchRecipeInput struct {
	UserID int64
	Query  string
	TagIDs []int64
}

func (uc *RecipeUseCase) CreateRecipe(ctx context.Context, input CreateRecipeInput) (*entity.Recipe, error) {
	recipe := &entity.Recipe{
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
	}

	return recipe, nil
}

func (uc *RecipeUseCase) GetRecipe(ctx context.Context, id, userID int64) (*entity.Recipe, error) {
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
	recipe.Images = images

	return recipe, nil
}

func (uc *RecipeUseCase) GetUserRecipes(ctx context.Context, userID int64) ([]entity.Recipe, error) {
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
		recipes[i].Images = images
	}

	return recipes, nil
}

func (uc *RecipeUseCase) SearchRecipes(ctx context.Context, input SearchRecipeInput) ([]entity.Recipe, error) {
	recipes, err := uc.recipeRepo.Search(ctx, input.UserID, input.Query, input.TagIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to search recipes: %w", err)
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

func (uc *RecipeUseCase) UpdateRecipe(ctx context.Context, input UpdateRecipeInput) (*entity.Recipe, error) {
	// Get existing recipe to verify ownership
	recipe, err := uc.recipeRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe: %w", err)
	}

	if recipe.UserID != input.UserID {
		return nil, errors.New("unauthorized access to recipe")
	}

	// Update recipe
	recipe.Name = input.Name
	recipe.URL = input.URL
	recipe.Memo = input.Memo

	if err := uc.recipeRepo.Update(ctx, recipe); err != nil {
		return nil, fmt.Errorf("failed to update recipe: %w", err)
	}

	return recipe, nil
}

func (uc *RecipeUseCase) DeleteRecipe(ctx context.Context, id, userID int64) error {
	// Get existing recipe to verify ownership
	recipe, err := uc.recipeRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get recipe: %w", err)
	}

	if recipe.UserID != userID {
		return errors.New("unauthorized access to recipe")
	}

	if err := uc.recipeRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete recipe: %w", err)
	}

	return nil
}

func (uc *RecipeUseCase) AddTagsToRecipe(ctx context.Context, recipeID, userID int64, tagIDs []int64) error {
	// Verify ownership
	recipe, err := uc.recipeRepo.GetByID(ctx, recipeID)
	if err != nil {
		return fmt.Errorf("failed to get recipe: %w", err)
	}

	if recipe.UserID != userID {
		return errors.New("unauthorized access to recipe")
	}

	return uc.recipeRepo.AddTags(ctx, recipeID, tagIDs)
}

func (uc *RecipeUseCase) RemoveTagsFromRecipe(ctx context.Context, recipeID, userID int64, tagIDs []int64) error {
	// Verify ownership
	recipe, err := uc.recipeRepo.GetByID(ctx, recipeID)
	if err != nil {
		return fmt.Errorf("failed to get recipe: %w", err)
	}

	if recipe.UserID != userID {
		return errors.New("unauthorized access to recipe")
	}

	return uc.recipeRepo.RemoveTags(ctx, recipeID, tagIDs)
}

func (uc *RecipeUseCase) AddImageToRecipe(ctx context.Context, recipeID, userID int64, imagePath string) (*entity.RecipeImage, error) {
	// Verify ownership
	recipe, err := uc.recipeRepo.GetByID(ctx, recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe: %w", err)
	}

	if recipe.UserID != userID {
		return nil, errors.New("unauthorized access to recipe")
	}

	image := &entity.RecipeImage{
		RecipeID:  recipeID,
		ImagePath: imagePath,
	}

	if err := uc.recipeImageRepo.Create(ctx, image); err != nil {
		return nil, fmt.Errorf("failed to add image to recipe: %w", err)
	}

	return image, nil
}

func (uc *RecipeUseCase) DeleteRecipeImage(ctx context.Context, imageID, userID int64) error {
	// Get image to find recipe
	images, err := uc.recipeImageRepo.GetByRecipeID(ctx, imageID)
	if err != nil || len(images) == 0 {
		return fmt.Errorf("image not found")
	}

	// Verify ownership through recipe
	recipe, err := uc.recipeRepo.GetByID(ctx, images[0].RecipeID)
	if err != nil {
		return fmt.Errorf("failed to get recipe: %w", err)
	}

	if recipe.UserID != userID {
		return errors.New("unauthorized access to recipe")
	}

	return uc.recipeImageRepo.Delete(ctx, imageID)
}

// Tag management

func (uc *RecipeUseCase) CreateTag(ctx context.Context, name string) (*entity.Tag, error) {
	// Check if tag already exists
	existingTag, err := uc.tagRepo.GetByName(ctx, name)
	if err == nil && existingTag != nil {
		return existingTag, nil // Return existing tag
	}

	tag := &entity.Tag{Name: name}
	if err := uc.tagRepo.Create(ctx, tag); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return tag, nil
}

func (uc *RecipeUseCase) GetAllTags(ctx context.Context) ([]entity.Tag, error) {
	return uc.tagRepo.GetAll(ctx)
}

func (uc *RecipeUseCase) DeleteTag(ctx context.Context, id int64) error {
	return uc.tagRepo.Delete(ctx, id)
}

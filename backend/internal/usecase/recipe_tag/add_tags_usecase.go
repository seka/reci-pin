package recipe_tag

import (
	"context"
	"errors"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

type AddTagsUseCase interface {
	Execute(ctx context.Context, recipeID, userID int64, tagIDs []int64) error
}

type addTagsInteractor struct {
	recipeRepo repository.RecipeRepository
	txManager  repository.TransactionManager
}

func NewAddTagsUseCase(recipeRepo repository.RecipeRepository, txManager repository.TransactionManager) AddTagsUseCase {
	return &addTagsInteractor{
		recipeRepo: recipeRepo,
		txManager:  txManager,
	}
}

func (uc *addTagsInteractor) Execute(ctx context.Context, recipeID, userID int64, tagIDs []int64) error {
	return uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// Verify ownership
		recipe, err := uc.recipeRepo.GetByID(txCtx, recipeID)
		if err != nil {
			return fmt.Errorf("failed to get recipe: %w", err)
		}

		if recipe.UserID != userID {
			return errors.New("unauthorized access to recipe")
		}

		return uc.recipeRepo.AddTags(txCtx, recipeID, tagIDs)
	})
}

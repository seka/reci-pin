package auth

import (
	"context"

	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

// WithdrawUseCase handles user account deletion.
// Note: Database CASCADE constraints (ON DELETE CASCADE) automatically
// remove all related data (recipes, tags, credentials, etc.) when a user is deleted.
type WithdrawUseCase interface {
	Execute(ctx context.Context, userID int64) error
}

type withdrawUseCase struct {
	userRepo repository.UserRepository
}

func NewWithdrawUseCase(userRepo repository.UserRepository) WithdrawUseCase {
	return &withdrawUseCase{userRepo: userRepo}
}

// Execute deletes the user account and all related data.
// Related data (recipes, recipe_images, recipe_tags, user_email_credentials)
// are automatically deleted by database CASCADE constraints.
func (u *withdrawUseCase) Execute(ctx context.Context, userID int64) error {
	return u.userRepo.Delete(ctx, userID)
}

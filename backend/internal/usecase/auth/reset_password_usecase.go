package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database/postgres"
)

type ResetPasswordInput struct {
	Token       string
	NewPassword string
}

type ResetPasswordUseCase interface {
	Execute(ctx context.Context, input ResetPasswordInput) error
}

type resetPasswordInteractor struct {
	tokenRepo repository.PasswordResetTokenRepository
	userRepo  repository.UserEmailCredentialRepository
	txManager repository.TransactionManager
}

func NewResetPasswordUseCase(
	tokenRepo repository.PasswordResetTokenRepository,
	userRepo repository.UserEmailCredentialRepository,
	txManager repository.TransactionManager,
) ResetPasswordUseCase {
	return &resetPasswordInteractor{
		tokenRepo: tokenRepo,
		userRepo:  userRepo,
		txManager: txManager,
	}
}

func (uc *resetPasswordInteractor) Execute(ctx context.Context, input ResetPasswordInput) error {
	// 1. Validate new password
	if err := ValidatePassword(input.NewPassword); err != nil {
		return err
	}

	return uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// 2. Find token
		t, err := uc.tokenRepo.GetByToken(txCtx, input.Token)
		if err != nil {
			return fmt.Errorf("failed to find token: %w", err)
		}
		if t == nil {
			return errors.New("invalid token")
		}

		// 3. Check expiration
		if t.ExpiresAt.Before(time.Now()) {
			return errors.New("token expired")
		}

		// 4. Hash new password
		hashedPassword, err := postgres.HashPassword(input.NewPassword)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		// 5. Update user password
		if err := uc.userRepo.UpdatePassword(txCtx, t.UserID, hashedPassword); err != nil {
			return fmt.Errorf("failed to update password: %w", err)
		}

		// 6. Delete token
		if err := uc.tokenRepo.Delete(txCtx, input.Token); err != nil {
			return fmt.Errorf("failed to delete token: %w", err)
		}

		return nil
	})
}

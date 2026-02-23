package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/notification"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database/postgres"
)

type ChangePasswordUseCase interface {
	Execute(ctx context.Context, userID int64, input ChangePasswordInput) error
}

type ChangePasswordInput struct {
	CurrentPassword string
	NewPassword     string
}

type changePasswordInteractor struct {
	credentialRepo repository.UserEmailCredentialRepository
	emailSender    notification.EmailClient
}

func NewChangePasswordUseCase(
	credentialRepo repository.UserEmailCredentialRepository,
	emailSender notification.EmailClient,
) ChangePasswordUseCase {
	return &changePasswordInteractor{
		credentialRepo: credentialRepo,
		emailSender:    emailSender,
	}
}

func (uc *changePasswordInteractor) Execute(ctx context.Context, userID int64, input ChangePasswordInput) error {
	// 1. Validate new password
	if err := ValidatePassword(input.NewPassword); err != nil {
		return err
	}

	// 2. Get current credential
	cred, err := uc.credentialRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user credential: %w", err)
	}
	if cred == nil {
		return errors.New("user credential not found")
	}

	// 3. Verify current password
	if !postgres.CheckPasswordHash(input.CurrentPassword, cred.PasswordHash) {
		return errors.New("invalid current password")
	}

	// 4. Hash new password
	hashedPassword, err := postgres.HashPassword(input.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 5. Update password
	if err := uc.credentialRepo.UpdatePassword(ctx, userID, hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// 6. Send notification email
	// Fire and forget, or handle error? For now, log error if fails but don't fail the request
	if err := uc.emailSender.SendPasswordChangeNotification(cred.Email); err != nil {
		// Log error (should use logger here)
		fmt.Printf("failed to send password change notification: %v\n", err)
	}

	return nil
}

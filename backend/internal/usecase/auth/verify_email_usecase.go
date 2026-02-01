package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

type VerifyEmailUseCase interface {
	Execute(ctx context.Context, token string) error
}

type verifyEmailInteractor struct {
	credentialRepo repository.UserEmailCredentialRepository
}

func NewVerifyEmailUseCase(credentialRepo repository.UserEmailCredentialRepository) VerifyEmailUseCase {
	return &verifyEmailInteractor{credentialRepo: credentialRepo}
}

func (uc *verifyEmailInteractor) Execute(ctx context.Context, token string) error {
	credential, err := uc.credentialRepo.GetByToken(ctx, token)
	if err != nil {
		return errors.New("invalid token")
	}

	if credential.IsVerified() {
		return errors.New("email already verified")
	}

	if credential.VerificationTokenExpiresAt != nil && credential.VerificationTokenExpiresAt.Before(time.Now()) {
		return errors.New("token expired")
	}

	// Update verification status
	now := time.Now()
	credential.EmailVerifiedAt = &now
	credential.VerificationToken = ""
	credential.VerificationTokenExpiresAt = nil

	if err := uc.credentialRepo.Update(ctx, credential); err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	return nil
}

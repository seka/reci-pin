package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/infrastructure/email"
)

type RequestPasswordResetUseCase interface {
	Execute(ctx context.Context, email string) error
}

type requestPasswordResetInteractor struct {
	userRepo    repository.UserEmailCredentialRepository
	tokenRepo   repository.PasswordResetTokenRepository
	emailSender email.EmailSender
}

func NewRequestPasswordResetUseCase(
	userRepo repository.UserEmailCredentialRepository,
	tokenRepo repository.PasswordResetTokenRepository,
	emailSender email.EmailSender,
) RequestPasswordResetUseCase {
	return &requestPasswordResetInteractor{
		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		emailSender: emailSender,
	}
}

func (uc *requestPasswordResetInteractor) Execute(ctx context.Context, emailAddr string) error {
	// 1. Check if user exists
	user, err := uc.userRepo.GetByEmail(ctx, emailAddr)
	if err != nil {
		// Return nil to avoid email enumeration attacks
		return nil
	}
	if user == nil {
		return nil
	}

	// 2. Generate token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// 3. Save token
	expiresAt := time.Now().Add(30 * time.Minute)
	if err := uc.tokenRepo.Save(ctx, token, user.UserID, expiresAt); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	// 4. Send email
	if err := uc.emailSender.SendPasswordReset(emailAddr, token); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

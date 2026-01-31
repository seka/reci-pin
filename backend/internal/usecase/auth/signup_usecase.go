package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/infrastructure/datastore/postgres"
)

type SignupUseCase struct {
	userRepo       repository.UserRepository
	credentialRepo repository.UserEmailCredentialRepository
}

func NewSignupUseCase(
	userRepo repository.UserRepository,
	credentialRepo repository.UserEmailCredentialRepository,
) *SignupUseCase {
	return &SignupUseCase{
		userRepo:       userRepo,
		credentialRepo: credentialRepo,
	}
}

type SignupInput struct {
	Email    string
	Password string
	Name     string
}

func (uc *SignupUseCase) Execute(ctx context.Context, input SignupInput) (int64, error) {
	// Check if email already exists
	existingCred, err := uc.credentialRepo.GetByEmail(ctx, input.Email)
	if err == nil && existingCred != nil {
		if existingCred.IsVerified() {
			return 0, errors.New("user with this email already exists")
		}
		// TODO: 未認証ユーザーへの確認メール再送フローを実装する
		return 0, errors.New("registration pending for this email")
	}

	// Hash password
	hashedPassword, err := postgres.HashPassword(input.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user (profile)
	user := &model.User{
		Name: input.Name,
	}
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return 0, fmt.Errorf("failed to create user profile: %w", err)
	}

	// Generate verification token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return 0, fmt.Errorf("failed to generate random token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)
	expiresAt := time.Now().Add(24 * time.Hour)

	// Create credential (unverified)
	credential := &model.UserEmailCredential{
		UserID:                     user.ID,
		Email:                      input.Email,
		PasswordHash:               hashedPassword,
		EmailVerifiedAt:            nil, // Unverified
		VerificationToken:          token,
		VerificationTokenExpiresAt: &expiresAt,
	}

	if err := uc.credentialRepo.Create(ctx, credential); err != nil {
		// TODO: ユーザー作成失敗時のロールバック（トランザクション）を実装する
		return 0, fmt.Errorf("failed to create user credential: %w", err)
	}

	// 検証用メール送信（現在はログ出力のみ）
	// TODO: EmailSenderサービスの実装
	log.Printf("==============================================")
	log.Printf("Email Sent to: %s", input.Email)
	log.Printf("Verification Link: /verify?token=%s", token)
	log.Printf("==============================================")

	return user.ID, nil
}

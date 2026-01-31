package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/infrastructure/datastore/postgres"
)

type SignupUseCase struct {
	userRepo repository.UserRepository
}

func NewSignupUseCase(userRepo repository.UserRepository) *SignupUseCase {
	return &SignupUseCase{userRepo: userRepo}
}

type SignupInput struct {
	Email    string
	Password string
	Name     string
}

func (uc *SignupUseCase) Execute(ctx context.Context, input SignupInput) (int64, error) {
	// Check if user already exists
	existingUser, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return 0, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := postgres.HashPassword(input.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &struct {
		ID           int64
		Email        string
		PasswordHash string
		Name         string
	}{
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Name:         input.Name,
	}

	if err := uc.userRepo.Create(ctx, (*struct {
		ID           int64
		Email        string
		PasswordHash string
		Name         string
		CreatedAt    string
		UpdatedAt    string
	})(user)); err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return user.ID, nil
}

package auth

import (
	"context"
	"errors"

	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/infrastructure/datastore/postgres"
)

type LoginUseCase struct {
	userRepo repository.UserRepository
}

func NewLoginUseCase(userRepo repository.UserRepository) *LoginUseCase {
	return &LoginUseCase{userRepo: userRepo}
}

type LoginInput struct {
	Email    string
	Password string
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (int64, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return 0, errors.New("invalid email or password")
	}

	// Check password
	if !postgres.CheckPasswordHash(input.Password, user.PasswordHash) {
		return 0, errors.New("invalid email or password")
	}

	return user.ID, nil
}

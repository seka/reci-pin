package auth

import (
	"context"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

type GetUserUseCase interface {
	Execute(ctx context.Context, userID int64) (*model.User, error)
}

type getUserInteractor struct {
	userRepo repository.UserRepository
}

func NewGetUserUseCase(userRepo repository.UserRepository) GetUserUseCase {
	return &getUserInteractor{userRepo: userRepo}
}

func (uc *getUserInteractor) Execute(ctx context.Context, userID int64) (*model.User, error) {
	return uc.userRepo.GetByID(ctx, userID)
}

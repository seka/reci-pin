package auth

import (
	"context"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

type GetUserUseCase struct {
	userRepo repository.UserRepository
}

func NewGetUserUseCase(userRepo repository.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{userRepo: userRepo}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, userID int64) (*model.User, error) {
	return uc.userRepo.GetByID(ctx, userID)
}

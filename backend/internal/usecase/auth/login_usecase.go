package auth

import (
	"context"
	"errors"

	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/infrastructure/datastore/postgres"
)

type LoginUseCase interface {
	Execute(ctx context.Context, input LoginInput) (int64, error)
}

type loginInteractor struct {
	credentialRepo repository.UserEmailCredentialRepository
}

func NewLoginUseCase(credentialRepo repository.UserEmailCredentialRepository) LoginUseCase {
	return &loginInteractor{credentialRepo: credentialRepo}
}

type LoginInput struct {
	Email    string
	Password string
}

func (uc *loginInteractor) Execute(ctx context.Context, input LoginInput) (int64, error) {
	// メールアドレスから認証情報を取得
	credential, err := uc.credentialRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return 0, errors.New("invalid email or password")
	}

	// パスワードチェック
	if !postgres.CheckPasswordHash(input.Password, credential.PasswordHash) {
		return 0, errors.New("invalid email or password")
	}

	// 認証済みかチェック
	if !credential.IsVerified() {
		return 0, errors.New("email not verified")
	}

	return credential.UserID, nil
}

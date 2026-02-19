package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository/mock"

	"github.com/seka/reci-pin/backend/internal/usecase/auth"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSignupUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		input   auth.SignupInput
		setup   func(*mock.MockUserRepository, *mock.MockUserEmailCredentialRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "正常系_ユーザー作成成功",
			input: auth.SignupInput{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			setup: func(mr *mock.MockUserRepository, mc *mock.MockUserEmailCredentialRepository) {
				// Email credential check (not exists)
				mc.EXPECT().GetByEmail(gomock.Any(), "test@example.com").Return(nil, errors.New("not found"))

				// Create User (Profile)
				mr.EXPECT().Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, user *model.User) error {
						user.ID = 1
						return nil
					})

				// Create Credential
				mc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "異常系_メール重複エラー（Verified）",
			input: auth.SignupInput{
				Email:    "duplicate@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			setup: func(mr *mock.MockUserRepository, mc *mock.MockUserEmailCredentialRepository) {
				now := time.Now()
				mc.EXPECT().GetByEmail(gomock.Any(), "duplicate@example.com").
					Return(&model.UserEmailCredential{
						Email:           "duplicate@example.com",
						EmailVerifiedAt: &now,
					}, nil)
			},
			wantErr: true,
			errMsg:  "user with this email already exists",
		},
		{
			name: "異常系_メール重複エラー（Unverified - 今回の仕様ではエラー）",
			input: auth.SignupInput{
				Email:    "pending@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			setup: func(mr *mock.MockUserRepository, mc *mock.MockUserEmailCredentialRepository) {
				mc.EXPECT().GetByEmail(gomock.Any(), "pending@example.com").
					Return(&model.UserEmailCredential{
						Email:           "pending@example.com",
						EmailVerifiedAt: nil,
					}, nil)
			},
			wantErr: true,
			errMsg:  "registration pending",
		},
		{
			name: "異常系_ユーザー作成失敗_DBエラー",
			input: auth.SignupInput{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			setup: func(mr *mock.MockUserRepository, mc *mock.MockUserEmailCredentialRepository) {
				mc.EXPECT().GetByEmail(gomock.Any(), "test@example.com").Return(nil, errors.New("not found"))
				mr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
			errMsg:  "failed to create user profile",
		},
		{
			name: "異常系_ユーザー作成失敗_パスワード要件満たさず",
			input: auth.SignupInput{
				Email:    "test@example.com",
				Password: "weak",
				Name:     "Test User",
			},
			setup: func(mr *mock.MockUserRepository, mc *mock.MockUserEmailCredentialRepository) {
				// No repository calls expected as validation fails first
			},
			wantErr: true,
			errMsg:  "validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := mock.NewMockUserRepository(ctrl)
			mockCredRepo := mock.NewMockUserEmailCredentialRepository(ctrl)
			tt.setup(mockUserRepo, mockCredRepo)

			uc := auth.NewSignupUseCase(mockUserRepo, mockCredRepo)
			userID, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Equal(t, int64(0), userID)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, int64(0), userID)
			}
		})
	}
}

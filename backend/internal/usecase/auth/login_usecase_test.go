package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/notification/mock"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database/postgres"
	"github.com/seka/reci-pin/backend/internal/usecase/auth"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLoginUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validPasswordHash, _ := postgres.HashPassword("password123")
	now := time.Now()

	tests := []struct {
		name    string
		input   auth.LoginInput
		setup   func(*mock.MockUserEmailCredentialRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "正常系_ログイン成功",
			input: auth.LoginInput{
				Email:    "test@example.com",
				Password: "password123",
			},
			setup: func(m *mock.MockUserEmailCredentialRepository) {
				m.EXPECT().GetByEmail(gomock.Any(), "test@example.com").
					Return(&model.UserEmailCredential{
						UserID:          1,
						Email:           "test@example.com",
						PasswordHash:    validPasswordHash,
						EmailVerifiedAt: &now,
					}, nil)
			},
			wantErr: false,
		},
		{
			name: "異常系_未認証ユーザー",
			input: auth.LoginInput{
				Email:    "unverified@example.com",
				Password: "password123",
			},
			setup: func(m *mock.MockUserEmailCredentialRepository) {
				m.EXPECT().GetByEmail(gomock.Any(), "unverified@example.com").
					Return(&model.UserEmailCredential{
						UserID:          2,
						Email:           "unverified@example.com",
						PasswordHash:    validPasswordHash,
						EmailVerifiedAt: nil,
					}, nil)
			},
			wantErr: true,
			errMsg:  "email not verified",
		},
		{
			name: "異常系_パスワード不一致",
			input: auth.LoginInput{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setup: func(m *mock.MockUserEmailCredentialRepository) {
				m.EXPECT().GetByEmail(gomock.Any(), "test@example.com").
					Return(&model.UserEmailCredential{
						UserID:       1,
						PasswordHash: validPasswordHash,
					}, nil)
			},
			wantErr: true,
			errMsg:  "invalid email or password",
		},
		{
			name: "異常系_ユーザー不在",
			input: auth.LoginInput{
				Email:    "notfound@example.com",
				Password: "password123",
			},
			setup: func(m *mock.MockUserEmailCredentialRepository) {
				m.EXPECT().GetByEmail(gomock.Any(), "notfound@example.com").
					Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "invalid email or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.NewMockUserEmailCredentialRepository(ctrl)
			tt.setup(mockRepo)

			uc := auth.NewLoginUseCase(mockRepo)
			userID, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Equal(t, int64(0), userID)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, int64(0), userID)
			}
		})
	}
}

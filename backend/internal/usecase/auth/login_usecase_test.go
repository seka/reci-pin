package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository/mock"
	"github.com/seka/reci-pin/backend/internal/infrastructure/datastore/postgres"
	"github.com/seka/reci-pin/backend/internal/usecase/auth"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLoginUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// テスト用パスワードハッシュを生成
	validPasswordHash, err := postgres.HashPassword("password123")
	assert.NoError(t, err)

	tests := []struct {
		name    string
		input   auth.LoginInput
		setup   func(*mock.MockUserRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "正常系_ログイン成功",
			input: auth.LoginInput{
				Email:    "test@example.com",
				Password: "password123",
			},
			setup: func(m *mock.MockUserRepository) {
				m.EXPECT().
					GetByEmail(gomock.Any(), "test@example.com").
					Return(&model.User{
						ID:    1,
						Email: "test@example.com",
						Name:  "Test User",
					}, validPasswordHash, nil)
			},
			wantErr: false,
		},
		{
			name: "異常系_ユーザー不在",
			input: auth.LoginInput{
				Email:    "notfound@example.com",
				Password: "password123",
			},
			setup: func(m *mock.MockUserRepository) {
				m.EXPECT().
					GetByEmail(gomock.Any(), "notfound@example.com").
					Return(nil, "", errors.New("user not found"))
			},
			wantErr: true,
			errMsg:  "invalid email or password",
		},
		{
			name: "異常系_パスワード不一致",
			input: auth.LoginInput{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setup: func(m *mock.MockUserRepository) {
				m.EXPECT().
					GetByEmail(gomock.Any(), "test@example.com").
					Return(&model.User{
						ID:    1,
						Email: "test@example.com",
						Name:  "Test User",
					}, validPasswordHash, nil)
			},
			wantErr: true,
			errMsg:  "invalid email or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.NewMockUserRepository(ctrl)
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

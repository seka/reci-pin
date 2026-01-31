package auth_test

import (
	"context"
	"errors"
	"testing"

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
		setup   func(*mock.MockUserRepository)
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
			setup: func(m *mock.MockUserRepository) {
				// メール重複チェック（存在しない）
				m.EXPECT().
					GetByEmail(gomock.Any(), "test@example.com").
					Return(nil, "", errors.New("not found"))

				// ユーザー作成
				m.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, user *model.User, passwordHash string) error {
						user.ID = 1
						return nil
					})
			},
			wantErr: false,
		},
		{
			name: "異常系_メール重複エラー",
			input: auth.SignupInput{
				Email:    "duplicate@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			setup: func(m *mock.MockUserRepository) {
				// メール重複チェック（既存ユーザーあり）
				m.EXPECT().
					GetByEmail(gomock.Any(), "duplicate@example.com").
					Return(&model.User{ID: 1, Email: "duplicate@example.com"}, "hash", nil)
			},
			wantErr: true,
			errMsg:  "user with this email already exists",
		},
		{
			name: "異常系_ユーザー作成失敗",
			input: auth.SignupInput{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			setup: func(m *mock.MockUserRepository) {
				m.EXPECT().
					GetByEmail(gomock.Any(), "test@example.com").
					Return(nil, "", errors.New("not found"))

				m.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "failed to create user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.NewMockUserRepository(ctrl)
			tt.setup(mockRepo)

			uc := auth.NewSignupUseCase(mockRepo)
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

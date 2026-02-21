package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	repo "github.com/seka/reci-pin/backend/internal/domain/repository/mock"
	"github.com/seka/reci-pin/backend/internal/usecase/auth"
	usecasemock "github.com/seka/reci-pin/backend/internal/usecase/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGenerateTokenUseCase_Execute(t *testing.T) {
	tests := []struct {
		name            string
		jwtSecret       string
		expirationHours int
		userID          int64
		wantErr         bool
	}{
		{
			name:            "正常系_トークン生成成功",
			jwtSecret:       "test-secret-key",
			expirationHours: 24,
			userID:          1,
			wantErr:         false,
		},
		{
			name:            "正常系_異なるユーザーID",
			jwtSecret:       "test-secret-key",
			expirationHours: 24,
			userID:          999,
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repo.NewMockRefreshTokenRepository(ctrl)
			if !tt.wantErr {
				mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
			}

			uc := auth.NewGenerateTokenUseCase(
				tt.jwtSecret,
				time.Duration(tt.expirationHours)*time.Hour,
				mockRepo,
				7*24*time.Hour,
			)

			result, err := uc.Execute(context.Background(), tt.userID, "test-agent", "127.0.0.1")

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken.Token)
				assert.NotEmpty(t, result.RefreshToken.Token)
				assert.WithinDuration(t, time.Now().Add(time.Duration(tt.expirationHours)*time.Hour), result.AccessToken.ExpiresAt, time.Minute)
			}
		})
	}
}

func TestValidateTokenUseCase_Execute(t *testing.T) {
	jwtSecret := "test-secret-key"
	expirationHours := 24

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo.NewMockRefreshTokenRepository(ctrl)
	mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// 有効なトークンを生成
	genUC := auth.NewGenerateTokenUseCase(jwtSecret, time.Duration(expirationHours)*time.Hour, mockRepo, 7*24*time.Hour)
	result, err := genUC.Execute(context.Background(), 1, "test-agent", "127.0.0.1")
	assert.NoError(t, err)
	validToken := result.AccessToken.Token

	tests := []struct {
		name       string
		jwtSecret  string
		token      string
		wantErr    bool
		wantUserID int64
	}{
		{
			name:       "正常系_有効なトークン",
			jwtSecret:  jwtSecret,
			token:      validToken,
			wantErr:    false,
			wantUserID: 1,
		},
		{
			name:      "異常系_不正なトークン",
			jwtSecret: jwtSecret,
			token:     "invalid.token.here",
			wantErr:   true,
		},
		{
			name:      "異常系_異なるシークレット",
			jwtSecret: "different-secret",
			token:     validToken,
			wantErr:   true,
		},
		{
			name:      "異常系_空のトークン",
			jwtSecret: jwtSecret,
			token:     "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := auth.NewValidateTokenUseCase(tt.jwtSecret)
			userID, err := uc.Execute(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUserID, userID)
			}
		})
	}
}

func TestTokenExpiration(t *testing.T) {
	jwtSecret := "test-secret-key"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repo.NewMockRefreshTokenRepository(ctrl)
	mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// 有効期限が極めて短いトークンを生成（テスト用）
	genUC := auth.NewGenerateTokenUseCase(jwtSecret, -1*time.Hour, mockRepo, 7*24*time.Hour) // -1時間 = 既に期限切れ
	result, err := genUC.Execute(context.Background(), 1, "test-agent", "127.0.0.1")
	assert.NoError(t, err)
	expiredToken := result.AccessToken.Token

	// 少し待機
	time.Sleep(100 * time.Millisecond)

	validUC := auth.NewValidateTokenUseCase(jwtSecret)
	_, err = validUC.Execute(expiredToken)

	// 期限切れトークンはエラーになるはず
	assert.Error(t, err)
}

func TestRefreshTokenUseCase_Execute(t *testing.T) {
	userID := int64(1)

	tests := []struct {
		name      string
		token     string
		setupMock func(m *repo.MockRefreshTokenRepository, mg *usecasemock.MockGenerateTokenUseCase)
		wantErr   bool
	}{
		{
			name:  "正常系_有効なトークンでリフレッシュ成功",
			token: "valid-token",
			setupMock: func(m *repo.MockRefreshTokenRepository, mg *usecasemock.MockGenerateTokenUseCase) {
				m.EXPECT().GetByHash(gomock.Any(), gomock.Any()).Return(&model.RefreshToken{
					ID:        1,
					UserID:    userID,
					ExpiresAt: time.Now().Add(24 * time.Hour),
				}, nil)
				m.EXPECT().Revoke(gomock.Any(), int64(1)).Return(nil)
				mg.EXPECT().Execute(gomock.Any(), userID, gomock.Any(), gomock.Any()).Return(&model.TokenResult{
					AccessToken: model.AuthToken{
						Token: "new-access",
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name:  "異常系_失効済みトークン使用で盗難検知",
			token: "revoked-token",
			setupMock: func(m *repo.MockRefreshTokenRepository, mg *usecasemock.MockGenerateTokenUseCase) {
				now := time.Now()
				m.EXPECT().GetByHash(gomock.Any(), gomock.Any()).Return(&model.RefreshToken{
					ID:        2,
					UserID:    userID,
					ExpiresAt: now.Add(24 * time.Hour),
					RevokedAt: &now, // 失効済み
				}, nil)
				// 全セッション失効が呼ばれるはず
				m.EXPECT().RevokeAllByUserID(gomock.Any(), userID).Return(nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repo.NewMockRefreshTokenRepository(ctrl)
			mockGenUC := usecasemock.NewMockGenerateTokenUseCase(ctrl)
			tt.setupMock(mockRepo, mockGenUC)

			uc := auth.NewRefreshTokenUseCase(mockGenUC, mockRepo)
			result, err := uc.Execute(context.Background(), tt.token, "agent", "1.1.1.1")

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

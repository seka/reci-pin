package auth_test

import (
	"testing"
	"time"

	"github.com/seka/reci-pin/backend/internal/usecase/auth"
	"github.com/stretchr/testify/assert"
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
			uc := auth.NewGenerateTokenUseCase(tt.jwtSecret, time.Duration(tt.expirationHours)*time.Hour)
			token, expiresAt, err := uc.Execute(tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.WithinDuration(t, time.Now().Add(time.Duration(tt.expirationHours)*time.Hour), expiresAt, time.Minute)
			}
		})
	}
}

func TestValidateTokenUseCase_Execute(t *testing.T) {
	jwtSecret := "test-secret-key"
	expirationHours := 24

	// 有効なトークンを生成
	genUC := auth.NewGenerateTokenUseCase(jwtSecret, time.Duration(expirationHours)*time.Hour)
	validToken, _, err := genUC.Execute(1)
	assert.NoError(t, err)

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

	// 有効期限が極めて短いトークンを生成（テスト用）
	genUC := auth.NewGenerateTokenUseCase(jwtSecret, -1*time.Hour) // -1時間 = 既に期限切れ
	expiredToken, _, err := genUC.Execute(1)
	assert.NoError(t, err)

	// 少し待機
	time.Sleep(100 * time.Millisecond)

	validUC := auth.NewValidateTokenUseCase(jwtSecret)
	_, err = validUC.Execute(expiredToken)

	// 期限切れトークンはエラーになるはず
	assert.Error(t, err)
}

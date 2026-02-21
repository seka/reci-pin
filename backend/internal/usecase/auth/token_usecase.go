package auth

//go:generate mockgen -source=$GOFILE -destination=../mock/token_usecase_mock.go -package=mock

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

type TokenResult struct {
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

type ValidateTokenUseCase interface {
	Execute(tokenString string) (int64, error)
}

type validateTokenInteractor struct {
	jwtSecret string
}

func NewValidateTokenUseCase(jwtSecret string) ValidateTokenUseCase {
	return &validateTokenInteractor{jwtSecret: jwtSecret}
}

func (uc *validateTokenInteractor) Execute(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(uc.jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int64(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, errors.New("invalid token")
}

type GenerateTokenUseCase interface {
	Execute(userID int64, userAgent, ipAddress string) (*TokenResult, error)
}

type generateTokenInteractor struct {
	jwtSecret              string
	jwtExpiration          time.Duration
	refreshTokenRepo       repository.RefreshTokenRepository
	refreshTokenExpiration time.Duration
}

func NewGenerateTokenUseCase(
	jwtSecret string,
	jwtExpiration time.Duration,
	refreshTokenRepo repository.RefreshTokenRepository,
	refreshTokenExpiration time.Duration,
) GenerateTokenUseCase {
	return &generateTokenInteractor{
		jwtSecret:              jwtSecret,
		jwtExpiration:          jwtExpiration,
		refreshTokenRepo:       refreshTokenRepo,
		refreshTokenExpiration: refreshTokenExpiration,
	}
}

func (uc *generateTokenInteractor) Execute(userID int64, userAgent, ipAddress string) (*TokenResult, error) {
	// Generate Access Token
	accessTokenExpiresAt := time.Now().Add(uc.jwtExpiration)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     accessTokenExpiresAt.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedAccessToken, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate Refresh Token
	refreshTokenString, err := generateRandomToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	refreshTokenExpiresAt := time.Now().Add(uc.refreshTokenExpiration)
	rtModel := &model.RefreshToken{
		UserID:    userID,
		TokenHash: HashToken(refreshTokenString),
		ExpiresAt: refreshTokenExpiresAt,
		UserAgent: userAgent,
		IPAddress: ipAddress,
	}

	if err := uc.refreshTokenRepo.Save(context.TODO(), rtModel); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &TokenResult{
		AccessToken:           signedAccessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          refreshTokenString,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}, nil
}

type RefreshTokenUseCase interface {
	Execute(refreshToken string, userAgent, ipAddress string) (*TokenResult, error)
}

type refreshTokenInteractor struct {
	genTokenUC GenerateTokenUseCase
	repo       repository.RefreshTokenRepository
}

func NewRefreshTokenUseCase(genTokenUC GenerateTokenUseCase, repo repository.RefreshTokenRepository) RefreshTokenUseCase {
	return &refreshTokenInteractor{
		genTokenUC: genTokenUC,
		repo:       repo,
	}
}

func (uc *refreshTokenInteractor) Execute(refreshToken string, userAgent, ipAddress string) (*TokenResult, error) {
	hash := HashToken(refreshToken)
	rt, err := uc.repo.GetByHash(context.TODO(), hash)
	if err != nil {
		return nil, fmt.Errorf("failed to find refresh token: %w", err)
	}

	if rt == nil {
		return nil, errors.New("invalid refresh token")
	}

	if rt.IsRevoked() {
		// 盗難検知: 使用済みのトークンが再利用された場合、該当ユーザーの全セッションを無効化する事も検討
		// 今回は一旦このトークンチェーンのみ失効とする（既に失効済み）
		return nil, errors.New("refresh token has been revoked")
	}

	if rt.IsExpired() {
		return nil, errors.New("refresh token has expired")
	}

	// トークンの回転 (Rotation): 現在のトークンを失効させる
	if err := uc.repo.Revoke(context.TODO(), rt.ID); err != nil {
		return nil, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	// 新しいトークンペアを発行
	return uc.genTokenUC.Execute(rt.UserID, userAgent, ipAddress)
}

type LogoutUseCase interface {
	Execute(refreshToken string) error
}

type logoutInteractor struct {
	repo repository.RefreshTokenRepository
}

func NewLogoutUseCase(repo repository.RefreshTokenRepository) LogoutUseCase {
	return &logoutInteractor{repo: repo}
}

func (uc *logoutInteractor) Execute(refreshToken string) error {
	hash := HashToken(refreshToken)
	rt, err := uc.repo.GetByHash(context.TODO(), hash)
	if err != nil {
		return fmt.Errorf("failed to find refresh token: %w", err)
	}
	if rt != nil {
		return uc.repo.Revoke(context.TODO(), rt.ID)
	}
	return nil
}

// Helpers
func generateRandomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func HashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

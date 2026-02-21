package auth

//go:generate mockgen -source=$GOFILE -destination=../mock/token_usecase_mock.go -package=mock

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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
	Execute(userID int64) (string, time.Time, error)
}

type generateTokenInteractor struct {
	jwtSecret     string
	jwtExpiration time.Duration
}

func NewGenerateTokenUseCase(jwtSecret string, expiration time.Duration) GenerateTokenUseCase {
	return &generateTokenInteractor{
		jwtSecret:     jwtSecret,
		jwtExpiration: expiration,
	}
}

func (uc *generateTokenInteractor) Execute(userID int64) (string, time.Time, error) {
	expiresAt := time.Now().Add(uc.jwtExpiration)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}
	return signedToken, expiresAt, nil
}

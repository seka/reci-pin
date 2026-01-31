package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ValidateTokenUseCase struct {
	jwtSecret string
}

func NewValidateTokenUseCase(jwtSecret string) *ValidateTokenUseCase {
	return &ValidateTokenUseCase{jwtSecret: jwtSecret}
}

func (uc *ValidateTokenUseCase) Execute(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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

type GenerateTokenUseCase struct {
	jwtSecret     string
	jwtExpiration time.Duration
}

func NewGenerateTokenUseCase(jwtSecret string, expirationHours int) *GenerateTokenUseCase {
	return &GenerateTokenUseCase{
		jwtSecret:     jwtSecret,
		jwtExpiration: time.Duration(expirationHours) * time.Hour,
	}
}

func (uc *GenerateTokenUseCase) Execute(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(uc.jwtExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}

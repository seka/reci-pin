package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/seka/reci-pin/backend/internal/domain/entity"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/infrastructure/datastore/postgres"
)

type AuthUseCase struct {
	userRepo      repository.UserRepository
	jwtSecret     string
	jwtExpiration time.Duration
}

func NewAuthUseCase(userRepo repository.UserRepository, jwtSecret string, expirationHours int) *AuthUseCase {
	return &AuthUseCase{
		userRepo:      userRepo,
		jwtSecret:     jwtSecret,
		jwtExpiration: time.Duration(expirationHours) * time.Hour,
	}
}

type SignupInput struct {
	Email    string
	Password string
	Name     string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *entity.User `json:"user"`
}

func (uc *AuthUseCase) Signup(ctx context.Context, input SignupInput) (*AuthResponse, error) {
	// Check if user already exists
	existingUser, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := postgres.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &entity.User{
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Name:         input.Name,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := uc.generateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	if !postgres.CheckPasswordHash(input.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := uc.generateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (uc *AuthUseCase) generateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(uc.jwtExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}

func (uc *AuthUseCase) ValidateToken(tokenString string) (int64, error) {
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

func (uc *AuthUseCase) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}

package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/server/handler"
	"github.com/seka/reci-pin/backend/internal/usecase/auth"
	"github.com/stretchr/testify/assert"
)

// Mock implementations

type mockSignupUseCase struct {
	ExecuteFunc func(ctx context.Context, input auth.SignupInput) (int64, error)
}

func (m *mockSignupUseCase) Execute(ctx context.Context, input auth.SignupInput) (int64, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return 1, nil
}

type mockLoginUseCase struct {
	ExecuteFunc func(ctx context.Context, input auth.LoginInput) (int64, error)
}

func (m *mockLoginUseCase) Execute(ctx context.Context, input auth.LoginInput) (int64, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return 1, nil
}

type mockGenerateTokenUseCase struct {
	ExecuteFunc func(userID int64) (string, error)
}

func (m *mockGenerateTokenUseCase) Execute(userID int64) (string, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(userID)
	}
	return "mock-token", nil
}

type mockGetUserUseCase struct {
	ExecuteFunc func(ctx context.Context, userID int64) (*model.User, error)
}

func (m *mockGetUserUseCase) Execute(ctx context.Context, userID int64) (*model.User, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, userID)
	}
	return &model.User{ID: userID, Name: "Test User"}, nil
}

type mockVerifyEmailUseCase struct {
	ExecuteFunc func(ctx context.Context, token string) error
}

func (m *mockVerifyEmailUseCase) Execute(ctx context.Context, token string) error {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, token)
	}
	return nil
}

// Tests

func TestAuthHandler_Signup(t *testing.T) {
	mockSignup := &mockSignupUseCase{
		ExecuteFunc: func(ctx context.Context, input auth.SignupInput) (int64, error) {
			assert.Equal(t, "test@example.com", input.Email)
			assert.Equal(t, "password123", input.Password)
			assert.Equal(t, "Test User", input.Name)
			return 123, nil
		},
	}

	h := handler.NewAuthHandler(
		mockSignup,
		&mockLoginUseCase{},
		&mockGenerateTokenUseCase{},
		&mockGetUserUseCase{},
		&mockVerifyEmailUseCase{},
	)

	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
		"name":     "Test User",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Signup(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	assert.Contains(t, response["message"], "Verification email sent")
}

func TestAuthHandler_Signup_InvalidJSON(t *testing.T) {
	h := handler.NewAuthHandler(
		&mockSignupUseCase{},
		&mockLoginUseCase{},
		&mockGenerateTokenUseCase{},
		&mockGetUserUseCase{},
		&mockVerifyEmailUseCase{},
	)

	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Signup(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Login(t *testing.T) {
	mockLogin := &mockLoginUseCase{
		ExecuteFunc: func(ctx context.Context, input auth.LoginInput) (int64, error) {
			assert.Equal(t, "test@example.com", input.Email)
			assert.Equal(t, "password123", input.Password)
			return 42, nil
		},
	}

	mockGetUser := &mockGetUserUseCase{
		ExecuteFunc: func(ctx context.Context, userID int64) (*model.User, error) {
			assert.Equal(t, int64(42), userID)
			return &model.User{ID: 42, Name: "Test User"}, nil
		},
	}

	mockGenerateToken := &mockGenerateTokenUseCase{
		ExecuteFunc: func(userID int64) (string, error) {
			assert.Equal(t, int64(42), userID)
			return "test-jwt-token", nil
		},
	}

	h := handler.NewAuthHandler(
		&mockSignupUseCase{},
		mockLogin,
		mockGenerateToken,
		mockGetUser,
		&mockVerifyEmailUseCase{},
	)

	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, "test-jwt-token", response["token"])
	assert.NotNil(t, response["user"])
}

func TestAuthHandler_Verify(t *testing.T) {
	mockVerify := &mockVerifyEmailUseCase{
		ExecuteFunc: func(ctx context.Context, token string) error {
			assert.Equal(t, "valid-token", token)
			return nil
		},
	}

	h := handler.NewAuthHandler(
		&mockSignupUseCase{},
		&mockLoginUseCase{},
		&mockGenerateTokenUseCase{},
		&mockGetUserUseCase{},
		mockVerify,
	)

	reqBody := map[string]string{
		"token": "valid-token",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/auth/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Verify(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	assert.Contains(t, response["message"], "verified successfully")
}

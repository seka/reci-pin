package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/reci-pin/backend/internal/server/middleware"
	"github.com/stretchr/testify/assert"
)

// Mock implementation

type mockValidateTokenUseCase struct {
	ExecuteFunc func(tokenString string) (int64, error)
}

func (m *mockValidateTokenUseCase) Execute(tokenString string) (int64, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(tokenString)
	}
	return 1, nil
}

// Tests

func TestAuthMiddleware_Authenticate_MissingToken(t *testing.T) {
	validateUC := &mockValidateTokenUseCase{}
	authMW := middleware.NewAuthMiddleware(validateUC)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("protected"))
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	handler := authMW.Authenticate(nextHandler)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "missing authentication")
}

func TestAuthMiddleware_Authenticate_ValidToken(t *testing.T) {
	validateUC := &mockValidateTokenUseCase{
		ExecuteFunc: func(tokenString string) (int64, error) {
			assert.Equal(t, "valid-token", tokenString)
			return 42, nil
		},
	}
	authMW := middleware.NewAuthMiddleware(validateUC)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserIDFromContext(r.Context())
		assert.True(t, ok)
		assert.Equal(t, int64(42), userID)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("protected"))
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	handler := authMW.Authenticate(nextHandler)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "protected", w.Body.String())
}

func TestAuthMiddleware_Authenticate_ValidCookie(t *testing.T) {
	validateUC := &mockValidateTokenUseCase{
		ExecuteFunc: func(tokenString string) (int64, error) {
			assert.Equal(t, "valid-cookie-token", tokenString)
			return 42, nil
		},
	}
	authMW := middleware.NewAuthMiddleware(validateUC)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserIDFromContext(r.Context())
		assert.True(t, ok)
		assert.Equal(t, int64(42), userID)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("protected"))
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: "valid-cookie-token"})
	w := httptest.NewRecorder()

	handler := authMW.Authenticate(nextHandler)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "protected", w.Body.String())
}

func TestAuthMiddleware_Authenticate_InvalidToken(t *testing.T) {
	validateUC := &mockValidateTokenUseCase{
		ExecuteFunc: func(tokenString string) (int64, error) {
			return 0, assert.AnError
		},
	}
	authMW := middleware.NewAuthMiddleware(validateUC)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Should not reach next handler")
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	handler := authMW.Authenticate(nextHandler)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid or expired token")
}

func TestAuthMiddleware_Authenticate_InvalidTokenFormat(t *testing.T) {
	validateUC := &mockValidateTokenUseCase{
		ExecuteFunc: func(tokenString string) (int64, error) {
			// Empty token should fail validation
			if tokenString == "" {
				return 0, assert.AnError
			}
			return 0, assert.AnError
		},
	}
	authMW := middleware.NewAuthMiddleware(validateUC)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	testCases := []struct {
		name   string
		header string
	}{
		{"No Bearer prefix", "just-a-token"},
		{"Empty token after Bearer", "Bearer "},
		{"Wrong prefix", "Token abc123"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", tc.header)
			w := httptest.NewRecorder()

			handler := authMW.Authenticate(nextHandler)
			handler.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}

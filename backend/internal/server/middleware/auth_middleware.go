package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/seka/reci-pin/backend/internal/usecase/auth"
)

type contextKey string

const UserIDKey contextKey = "user_id"

type AuthMiddleware struct {
	validateTokenUseCase *auth.ValidateTokenUseCase
}

func NewAuthMiddleware(validateTokenUseCase *auth.ValidateTokenUseCase) *AuthMiddleware {
	return &AuthMiddleware{validateTokenUseCase: validateTokenUseCase}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		userID, err := m.validateTokenUseCase.Execute(token)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}

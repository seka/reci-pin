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
	validateTokenUseCase auth.ValidateTokenUseCase
}

func NewAuthMiddleware(validateTokenUseCase auth.ValidateTokenUseCase) *AuthMiddleware {
	return &AuthMiddleware{validateTokenUseCase: validateTokenUseCase}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		// 1. Try to get token from Cookie (Preferred for SSR/Web)
		if cookie, err := r.Cookie("auth_token"); err == nil {
			token = cookie.Value
		}

		// 2. Try to get token from Authorization header (Fallback for API/Mobile)
		if token == "" {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					token = parts[1]
				}
			}
		}

		if token == "" {
			http.Error(w, "missing authentication", http.StatusUnauthorized)
			return
		}

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

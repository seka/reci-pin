package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/server/handler"
	"github.com/seka/reci-pin/backend/internal/server/middleware"
	usecasemock "github.com/seka/reci-pin/backend/internal/usecase/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthHandler_Signup(t *testing.T) {
	type args struct {
		body map[string]any
	}
	type mocks struct {
		setup func(m *usecasemock.MockSignupUseCase)
	}
	tests := []struct {
		name       string
		args       args
		mocks      mocks
		wantStatus int
	}{
		{
			name: "Success",
			args: args{
				body: map[string]any{
					"email":    "test@example.com",
					"password": "password",
					"name":     "Test",
				},
			},
			mocks: mocks{
				setup: func(m *usecasemock.MockSignupUseCase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(int64(1), nil)
				},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Invalid JSON",
			args: args{body: nil},
			mocks: mocks{
				setup: func(m *usecasemock.MockSignupUseCase) {},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "UseCase Error",
			args: args{
				body: map[string]any{
					"email":    "exists@example.com",
					"password": "password",
				},
			},
			mocks: mocks{
				setup: func(m *usecasemock.MockSignupUseCase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any()).
						Return(int64(0), errors.New("exists"))
				},
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSignup := usecasemock.NewMockSignupUseCase(ctrl)
			tt.mocks.setup(mockSignup)

			h := handler.NewAuthHandler(
				mockSignup,
				usecasemock.NewMockLoginUseCase(ctrl),
				usecasemock.NewMockGenerateTokenUseCase(ctrl),
				usecasemock.NewMockGetUserUseCase(ctrl),
				usecasemock.NewMockVerifyEmailUseCase(ctrl),
				usecasemock.NewMockWithdrawUseCase(ctrl),
				usecasemock.NewMockChangePasswordUseCase(ctrl),
				usecasemock.NewMockRequestPasswordResetUseCase(ctrl),
				usecasemock.NewMockResetPasswordUseCase(ctrl),
			)

			var req *http.Request
			if tt.name == "Invalid JSON" {
				req = httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader([]byte("invalid")))
			} else {
				body, _ := json.Marshal(tt.args.body)
				req = httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader(body))
			}
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			h.Signup(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	type mocks struct {
		login    func(m *usecasemock.MockLoginUseCase)
		getUser  func(m *usecasemock.MockGetUserUseCase)
		genToken func(m *usecasemock.MockGenerateTokenUseCase)
	}
	tests := []struct {
		name       string
		body       map[string]any
		mocks      mocks
		wantStatus int
	}{
		{
			name: "Success",
			body: map[string]any{"email": "a@b.com", "password": "p"},
			mocks: mocks{
				login: func(m *usecasemock.MockLoginUseCase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(int64(1), nil)
				},
				getUser: func(m *usecasemock.MockGetUserUseCase) {
					m.EXPECT().Execute(gomock.Any(), int64(1)).Return(&model.User{ID: 1}, nil)
				},
				genToken: func(m *usecasemock.MockGenerateTokenUseCase) {
					m.EXPECT().Execute(int64(1)).Return("token", time.Now().Add(24*time.Hour), nil)
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Login Failed",
			body: map[string]any{"email": "a@b.com", "password": "p"},
			mocks: mocks{
				login: func(m *usecasemock.MockLoginUseCase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("auth failed"))
				},
				getUser:  func(m *usecasemock.MockGetUserUseCase) {},
				genToken: func(m *usecasemock.MockGenerateTokenUseCase) {},
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "GetUser Error",
			body: map[string]any{"email": "a@b.com", "password": "p"},
			mocks: mocks{
				login: func(m *usecasemock.MockLoginUseCase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(int64(1), nil)
				},
				getUser: func(m *usecasemock.MockGetUserUseCase) {
					m.EXPECT().Execute(gomock.Any(), int64(1)).Return(nil, errors.New("db error"))
				},
				genToken: func(m *usecasemock.MockGenerateTokenUseCase) {},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "GenToken Error",
			body: map[string]any{"email": "a@b.com", "password": "p"},
			mocks: mocks{
				login: func(m *usecasemock.MockLoginUseCase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(int64(1), nil)
				},
				getUser: func(m *usecasemock.MockGetUserUseCase) {
					m.EXPECT().Execute(gomock.Any(), int64(1)).Return(&model.User{ID: 1}, nil)
				},
				genToken: func(m *usecasemock.MockGenerateTokenUseCase) {
					m.EXPECT().Execute(int64(1)).Return("", time.Time{}, errors.New("token error"))
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLogin := usecasemock.NewMockLoginUseCase(ctrl)
			mockGetUser := usecasemock.NewMockGetUserUseCase(ctrl)
			mockGenToken := usecasemock.NewMockGenerateTokenUseCase(ctrl)

			tt.mocks.login(mockLogin)
			tt.mocks.getUser(mockGetUser)
			tt.mocks.genToken(mockGenToken)

			h := handler.NewAuthHandler(
				usecasemock.NewMockSignupUseCase(ctrl),
				mockLogin,
				mockGenToken,
				mockGetUser,
				usecasemock.NewMockVerifyEmailUseCase(ctrl),
				usecasemock.NewMockWithdrawUseCase(ctrl),
				usecasemock.NewMockChangePasswordUseCase(ctrl),
				usecasemock.NewMockRequestPasswordResetUseCase(ctrl),
				usecasemock.NewMockResetPasswordUseCase(ctrl),
			)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.Login(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestAuthHandler_Verify(t *testing.T) {
	tests := []struct {
		name       string
		body       map[string]string
		setupMock  func(m *usecasemock.MockVerifyEmailUseCase)
		wantStatus int
	}{
		{
			name: "Success",
			body: map[string]string{"token": "valid"},
			setupMock: func(m *usecasemock.MockVerifyEmailUseCase) {
				m.EXPECT().Execute(gomock.Any(), "valid").Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Invalid Token",
			body: map[string]string{"token": "invalid"},
			setupMock: func(m *usecasemock.MockVerifyEmailUseCase) {
				m.EXPECT().Execute(gomock.Any(), "invalid").Return(errors.New("invalid token"))
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockVerify := usecasemock.NewMockVerifyEmailUseCase(ctrl)
			tt.setupMock(mockVerify)

			h := handler.NewAuthHandler(
				usecasemock.NewMockSignupUseCase(ctrl),
				usecasemock.NewMockLoginUseCase(ctrl),
				usecasemock.NewMockGenerateTokenUseCase(ctrl),
				usecasemock.NewMockGetUserUseCase(ctrl),
				mockVerify,
				usecasemock.NewMockWithdrawUseCase(ctrl),
				usecasemock.NewMockChangePasswordUseCase(ctrl),
				usecasemock.NewMockRequestPasswordResetUseCase(ctrl),
				usecasemock.NewMockResetPasswordUseCase(ctrl),
			)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/auth/verify", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.Verify(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestAuthHandler_ChangePassword(t *testing.T) {
	tests := []struct {
		name       string
		body       map[string]any
		setupMock  func(m *usecasemock.MockChangePasswordUseCase)
		wantStatus int
	}{
		{
			name: "Success",
			body: map[string]any{
				"current_password": "old",
				"new_password":     "new",
			},
			setupMock: func(m *usecasemock.MockChangePasswordUseCase) {
				m.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Invalid JSON",
			body: nil,
			setupMock: func(m *usecasemock.MockChangePasswordUseCase) {
				// No call expected
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "UseCase Error",
			body: map[string]any{
				"current_password": "old",
				"new_password":     "new",
			},
			setupMock: func(m *usecasemock.MockChangePasswordUseCase) {
				m.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockChangePassword := usecasemock.NewMockChangePasswordUseCase(ctrl)
			tt.setupMock(mockChangePassword)

			h := handler.NewAuthHandler(
				usecasemock.NewMockSignupUseCase(ctrl),
				usecasemock.NewMockLoginUseCase(ctrl),
				usecasemock.NewMockGenerateTokenUseCase(ctrl),
				usecasemock.NewMockGetUserUseCase(ctrl),
				usecasemock.NewMockVerifyEmailUseCase(ctrl),
				usecasemock.NewMockWithdrawUseCase(ctrl),
				mockChangePassword,
				usecasemock.NewMockRequestPasswordResetUseCase(ctrl),
				usecasemock.NewMockResetPasswordUseCase(ctrl),
			)

			var req *http.Request
			if tt.name == "Invalid JSON" {
				req = httptest.NewRequest(http.MethodPut, "/auth/password", bytes.NewReader([]byte("invalid")))
			} else {
				body, _ := json.Marshal(tt.body)
				req = httptest.NewRequest(http.MethodPut, "/auth/password", bytes.NewReader(body))
			}
			req.Header.Set("Content-Type", "application/json")
			// Add dummy user ID to context
			ctx := context.WithValue(req.Context(), middleware.UserIDKey, int64(1))
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.ChangePassword(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	tests := []struct {
		name       string
		wantStatus int
	}{
		{
			name:       "Success",
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := handler.NewAuthHandler(
				usecasemock.NewMockSignupUseCase(ctrl),
				usecasemock.NewMockLoginUseCase(ctrl),
				usecasemock.NewMockGenerateTokenUseCase(ctrl),
				usecasemock.NewMockGetUserUseCase(ctrl),
				usecasemock.NewMockVerifyEmailUseCase(ctrl),
				usecasemock.NewMockWithdrawUseCase(ctrl),
				usecasemock.NewMockChangePasswordUseCase(ctrl),
				usecasemock.NewMockRequestPasswordResetUseCase(ctrl),
				usecasemock.NewMockResetPasswordUseCase(ctrl),
			)

			req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
			w := httptest.NewRecorder()

			h.Logout(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestAuthHandler_RequestPasswordReset(t *testing.T) {
	tests := []struct {
		name       string
		body       map[string]string
		setupMock  func(m *usecasemock.MockRequestPasswordResetUseCase)
		wantStatus int
	}{
		{
			name: "Success",
			body: map[string]string{"email": "test@example.com"},
			setupMock: func(m *usecasemock.MockRequestPasswordResetUseCase) {
				m.EXPECT().Execute(gomock.Any(), "test@example.com").Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "UseCase Error",
			body: map[string]string{"email": "error@example.com"},
			setupMock: func(m *usecasemock.MockRequestPasswordResetUseCase) {
				m.EXPECT().Execute(gomock.Any(), "error@example.com").Return(errors.New("error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRequest := usecasemock.NewMockRequestPasswordResetUseCase(ctrl)
			tt.setupMock(mockRequest)

			h := handler.NewAuthHandler(
				usecasemock.NewMockSignupUseCase(ctrl),
				usecasemock.NewMockLoginUseCase(ctrl),
				usecasemock.NewMockGenerateTokenUseCase(ctrl),
				usecasemock.NewMockGetUserUseCase(ctrl),
				usecasemock.NewMockVerifyEmailUseCase(ctrl),
				usecasemock.NewMockWithdrawUseCase(ctrl),
				usecasemock.NewMockChangePasswordUseCase(ctrl),
				mockRequest,
				usecasemock.NewMockResetPasswordUseCase(ctrl),
			)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/auth/password-reset/request", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.RequestPasswordReset(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestAuthHandler_ResetPassword(t *testing.T) {
	tests := []struct {
		name       string
		body       map[string]string
		setupMock  func(m *usecasemock.MockResetPasswordUseCase)
		wantStatus int
	}{
		{
			name: "Success",
			body: map[string]string{
				"token":        "valid_token",
				"new_password": "new_password",
			},
			setupMock: func(m *usecasemock.MockResetPasswordUseCase) {
				m.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "UseCase Error",
			body: map[string]string{
				"token":        "invalid_token",
				"new_password": "new_password",
			},
			setupMock: func(m *usecasemock.MockResetPasswordUseCase) {
				m.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockReset := usecasemock.NewMockResetPasswordUseCase(ctrl)
			tt.setupMock(mockReset)

			h := handler.NewAuthHandler(
				usecasemock.NewMockSignupUseCase(ctrl),
				usecasemock.NewMockLoginUseCase(ctrl),
				usecasemock.NewMockGenerateTokenUseCase(ctrl),
				usecasemock.NewMockGetUserUseCase(ctrl),
				usecasemock.NewMockVerifyEmailUseCase(ctrl),
				usecasemock.NewMockWithdrawUseCase(ctrl),
				usecasemock.NewMockChangePasswordUseCase(ctrl),
				usecasemock.NewMockRequestPasswordResetUseCase(ctrl),
				mockReset,
			)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/auth/password-reset", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.ResetPassword(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/reci-pin/backend/config"
	registrymock "github.com/seka/reci-pin/backend/internal/registry/mock"
	usecasemock "github.com/seka/reci-pin/backend/internal/usecase/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupMockRegistry(ctrl *gomock.Controller, m *registrymock.MockUseCase) {
	// Auth
	m.EXPECT().NewValidateTokenUseCase().Return(usecasemock.NewMockValidateTokenUseCase(ctrl))
	m.EXPECT().NewSignupUseCase().Return(usecasemock.NewMockSignupUseCase(ctrl))
	m.EXPECT().NewLoginUseCase().Return(usecasemock.NewMockLoginUseCase(ctrl))
	m.EXPECT().NewGenerateTokenUseCase().Return(usecasemock.NewMockGenerateTokenUseCase(ctrl))
	m.EXPECT().NewGetUserUseCase().Return(usecasemock.NewMockGetUserUseCase(ctrl))
	m.EXPECT().NewVerifyEmailUseCase().Return(usecasemock.NewMockVerifyEmailUseCase(ctrl))
	m.EXPECT().NewWithdrawUseCase().Return(usecasemock.NewMockWithdrawUseCase(ctrl))
	m.EXPECT().NewChangePasswordUseCase().Return(usecasemock.NewMockChangePasswordUseCase(ctrl))
	m.EXPECT().NewRequestPasswordResetUseCase().Return(usecasemock.NewMockRequestPasswordResetUseCase(ctrl))
	m.EXPECT().NewResetPasswordUseCase().Return(usecasemock.NewMockResetPasswordUseCase(ctrl))
	m.EXPECT().NewRefreshTokenUseCase().Return(usecasemock.NewMockRefreshTokenUseCase(ctrl))
	m.EXPECT().NewLogoutUseCase().Return(usecasemock.NewMockLogoutUseCase(ctrl))

	// Recipe
	m.EXPECT().NewCreateRecipeUseCase().Return(usecasemock.NewMockCreateRecipeUseCase(ctrl))
	m.EXPECT().NewGetRecipeUseCase().Return(usecasemock.NewMockGetRecipeUseCase(ctrl))
	m.EXPECT().NewGetUserRecipesUseCase().Return(usecasemock.NewMockGetUserRecipesUseCase(ctrl))
	m.EXPECT().NewUpdateRecipeUseCase().Return(usecasemock.NewMockUpdateRecipeUseCase(ctrl))
	m.EXPECT().NewDeleteRecipeUseCase().Return(usecasemock.NewMockDeleteRecipeUseCase(ctrl))
	m.EXPECT().NewSearchRecipesUseCase().Return(usecasemock.NewMockSearchRecipesUseCase(ctrl))

	// Recipe Tag
	m.EXPECT().NewAddTagsUseCase().Return(usecasemock.NewMockAddTagsUseCase(ctrl))
	m.EXPECT().NewRemoveTagsUseCase().Return(usecasemock.NewMockRemoveTagsUseCase(ctrl))

	// Recipe Image
	m.EXPECT().NewCreateRecipeImageUseCase().Return(usecasemock.NewMockCreateRecipeImageUseCase(ctrl))

	// Tag
	m.EXPECT().NewCreateTagUseCase().Return(usecasemock.NewMockCreateTagUseCase(ctrl))
	m.EXPECT().NewGetAllTagsUseCase().Return(usecasemock.NewMockGetAllTagsUseCase(ctrl))
	m.EXPECT().NewDeleteTagUseCase().Return(usecasemock.NewMockDeleteTagUseCase(ctrl))
}

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{Port: "8080"},
	}
	mockRegistry := registrymock.NewMockUseCase(ctrl)
	setupMockRegistry(ctrl, mockRegistry)

	srv := New(cfg, mockRegistry)

	assert.NotNil(t, srv)
	assert.NotNil(t, srv.router)
	assert.Equal(t, cfg, srv.cfg)
}

func TestServer_HealthEndpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	mockRegistry := registrymock.NewMockUseCase(ctrl)
	setupMockRegistry(ctrl, mockRegistry)

	srv := New(cfg, mockRegistry)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

func TestServer_CORSHeaders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{}
	mockRegistry := registrymock.NewMockUseCase(ctrl)
	setupMockRegistry(ctrl, mockRegistry)

	srv := New(cfg, mockRegistry)

	req := httptest.NewRequest("OPTIONS", "/health", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestServer_RoutingExists(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"Health", "GET", "/health"},
		{"Signup", "POST", "/api/auth/signup"},
		{"Login", "POST", "/api/auth/login"},
		// We can add more route checks here
		{"GetRecipe", "GET", "/api/recipes/1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cfg := &config.Config{}
			mockRegistry := registrymock.NewMockUseCase(ctrl)
			setupMockRegistry(ctrl, mockRegistry)

			srv := New(cfg, mockRegistry)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			srv.ServeHTTP(w, req)

			// Route exists if it's NOT 404
			assert.NotEqual(t, http.StatusNotFound, w.Code, "Route %s %s should exist", tt.method, tt.path)
		})
	}
}

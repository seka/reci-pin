package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/reci-pin/backend/config"
)

// mockRepository implements registry.Repository for testing
type mockRepository struct{}

func (m *mockRepository) NewUserRepository() interface{}      { return nil }
func (m *mockRepository) NewRecipeRepository() interface{}    { return nil }
func (m *mockRepository) NewTagRepository() interface{}       { return nil }
func (m *mockRepository) NewRecipeTagRepository() interface{} { return nil }
func (m *mockRepository) NewImageRepository() interface{}     { return nil }
func (m *mockRepository) Close() error                        { return nil }

// mockUseCase implements registry.UseCase for testing
type mockUseCase struct{}

func (m *mockUseCase) NewSignupUseCase() interface{}         { return nil }
func (m *mockUseCase) NewLoginUseCase() interface{}          { return nil }
func (m *mockUseCase) NewGenerateTokenUseCase() interface{}  { return nil }
func (m *mockUseCase) NewValidateTokenUseCase() interface{}  { return nil }
func (m *mockUseCase) NewGetUserUseCase() interface{}        { return nil }
func (m *mockUseCase) NewVerifyEmailUseCase() interface{}    { return nil }
func (m *mockUseCase) NewCreateRecipeUseCase() interface{}   { return nil }
func (m *mockUseCase) NewGetRecipeUseCase() interface{}      { return nil }
func (m *mockUseCase) NewGetUserRecipesUseCase() interface{} { return nil }
func (m *mockUseCase) NewUpdateRecipeUseCase() interface{}   { return nil }
func (m *mockUseCase) NewDeleteRecipeUseCase() interface{}   { return nil }
func (m *mockUseCase) NewSearchRecipesUseCase() interface{}  { return nil }
func (m *mockUseCase) NewAddTagsUseCase() interface{}        { return nil }
func (m *mockUseCase) NewRemoveTagsUseCase() interface{}     { return nil }
func (m *mockUseCase) NewAddImageUseCase() interface{}       { return nil }
func (m *mockUseCase) NewCreateTagUseCase() interface{}      { return nil }
func (m *mockUseCase) NewGetAllTagsUseCase() interface{}     { return nil }
func (m *mockUseCase) NewDeleteTagUseCase() interface{}      { return nil }

func TestNew(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Port: 8080},
	}
	mockRepo := &mockRepository{}
	mockUC := &mockUseCase{}

	srv := New(cfg, mockRepo, mockUC)

	if srv == nil {
		t.Fatal("Server should not be nil")
	}
	if srv.router == nil {
		t.Fatal("Router should be initialized")
	}
	if srv.cfg != cfg {
		t.Error("Config should be set")
	}
}

func TestServer_HealthEndpoint(t *testing.T) {
	cfg := &config.Config{}
	mockRepo := &mockRepository{}
	mockUC := &mockUseCase{}

	srv := New(cfg, mockRepo, mockUC)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
	if w.Body.String() != "OK" {
		t.Errorf("Expected 'OK', got %s", w.Body.String())
	}
}

func TestServer_CORSHeaders(t *testing.T) {
	cfg := &config.Config{}
	mockRepo := &mockRepository{}
	mockUC := &mockUseCase{}

	srv := New(cfg, mockRepo, mockUC)

	req := httptest.NewRequest("OPTIONS", "/health", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS origin header not set correctly")
	}
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected 204 for OPTIONS, got %d", w.Code)
	}
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
	}

	cfg := &config.Config{}
	mockRepo := &mockRepository{}
	mockUC := &mockUseCase{}

	srv := New(cfg, mockRepo, mockUC)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			srv.ServeHTTP(w, req)

			// ルートが存在することを確認（404でないこと）
			// ただし Handler が nil なので 500 などになる可能性あり
			if w.Code == http.StatusNotFound {
				t.Errorf("Route %s %s not found", tt.method, tt.path)
			}
		})
	}
}

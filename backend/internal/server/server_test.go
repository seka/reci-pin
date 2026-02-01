package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/usecase/auth"
	"github.com/seka/reci-pin/backend/internal/usecase/recipe"
	"github.com/seka/reci-pin/backend/internal/usecase/recipe_image"
	"github.com/seka/reci-pin/backend/internal/usecase/recipe_tag"
	"github.com/seka/reci-pin/backend/internal/usecase/tag"
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

func (m *mockUseCase) NewSignupUseCase() auth.SignupUseCase               { return nil }
func (m *mockUseCase) NewLoginUseCase() auth.LoginUseCase                 { return nil }
func (m *mockUseCase) NewGenerateTokenUseCase() auth.GenerateTokenUseCase { return nil }
func (m *mockUseCase) NewValidateTokenUseCase() auth.ValidateTokenUseCase { return nil }
func (m *mockUseCase) NewGetUserUseCase() auth.GetUserUseCase             { return nil }
func (m *mockUseCase) NewVerifyEmailUseCase() auth.VerifyEmailUseCase     { return nil }
func (m *mockUseCase) NewCreateRecipeUseCase() recipe.CreateRecipeUseCase { return nil }
func (m *mockUseCase) NewGetRecipeUseCase() recipe.GetRecipeUseCase       { return nil }
func (m *mockUseCase) NewGetUserRecipesUseCase() recipe.GetUserRecipesUseCase {
	return nil
}
func (m *mockUseCase) NewUpdateRecipeUseCase() recipe.UpdateRecipeUseCase { return nil }
func (m *mockUseCase) NewDeleteRecipeUseCase() recipe.DeleteRecipeUseCase { return nil }
func (m *mockUseCase) NewSearchRecipesUseCase() recipe.SearchRecipesUseCase {
	return nil
}
func (m *mockUseCase) NewAddTagsUseCase() recipe_tag.AddTagsUseCase       { return nil }
func (m *mockUseCase) NewRemoveTagsUseCase() recipe_tag.RemoveTagsUseCase { return nil }
func (m *mockUseCase) NewAddImageUseCase() recipe_image.AddImageUseCase   { return nil }
func (m *mockUseCase) NewCreateTagUseCase() tag.CreateTagUseCase          { return nil }
func (m *mockUseCase) NewGetAllTagsUseCase() tag.GetAllTagsUseCase        { return nil }
func (m *mockUseCase) NewDeleteTagUseCase() tag.DeleteTagUseCase          { return nil }

func TestNew(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{Port: 8080},
	}
	// mockRepo := &mockRepository{}
	mockUC := &mockUseCase{}

	srv := New(cfg, mockUC)

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
	// mockRepo := &mockRepository{}
	mockUC := &mockUseCase{}

	srv := New(cfg, mockUC)

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
	// mockRepo := &mockRepository{}
	mockUC := &mockUseCase{}

	srv := New(cfg, mockUC)

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
	// mockRepo := &mockRepository{}
	mockUC := &mockUseCase{}

	srv := New(cfg, mockUC)

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

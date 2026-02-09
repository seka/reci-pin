package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/registry"
	"github.com/seka/reci-pin/backend/internal/server/handler"
	"github.com/seka/reci-pin/backend/internal/server/middleware"
)

type Server struct {
	router          *chi.Mux
	httpServer      *http.Server
	useCaseRegistry registry.UseCase
	cfg             *config.Config
}

// New creates a new server instance with the given dependencies
func New(cfg *config.Config, useCase registry.UseCase) *Server {
	s := &Server{
		router:          chi.NewRouter(),
		useCaseRegistry: useCase,
		cfg:             cfg,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

func (s *Server) setupMiddleware() {
	s.router.Use(chimiddleware.Logger)
	s.router.Use(chimiddleware.Recoverer)
	s.router.Use(middleware.CORS)
}

func (s *Server) setupRoutes() {
	authMiddleware := middleware.NewAuthMiddleware(s.useCaseRegistry.NewValidateTokenUseCase())

	// Public routes
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Auth public routes
	authHandler := handler.NewAuthHandler(
		s.useCaseRegistry.NewSignupUseCase(),
		s.useCaseRegistry.NewLoginUseCase(),
		s.useCaseRegistry.NewGenerateTokenUseCase(),
		s.useCaseRegistry.NewGetUserUseCase(),
		s.useCaseRegistry.NewVerifyEmailUseCase(),
		s.useCaseRegistry.NewWithdrawUseCase(),
		s.useCaseRegistry.NewChangePasswordUseCase(),
	)

	s.router.Post("/api/auth/signup", authHandler.Signup)
	s.router.Post("/api/auth/login", authHandler.Login)
	s.router.Post("/api/auth/verify", authHandler.Verify)

	s.router.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)

		// Auth (authenticated routes)
		r.Delete("/api/auth/withdraw", authHandler.Withdraw)
		r.Put("/api/auth/password", authHandler.ChangePassword)

		// Recipes
		recipeHandler := handler.NewRecipeHandler(
			s.useCaseRegistry.NewCreateRecipeUseCase(),
			s.useCaseRegistry.NewGetRecipeUseCase(),
			s.useCaseRegistry.NewGetUserRecipesUseCase(),
			s.useCaseRegistry.NewUpdateRecipeUseCase(),
			s.useCaseRegistry.NewDeleteRecipeUseCase(),
			s.useCaseRegistry.NewSearchRecipesUseCase(),
			s.useCaseRegistry.NewAddTagsUseCase(),
			s.useCaseRegistry.NewRemoveTagsUseCase(),
			s.useCaseRegistry.NewAddImageUseCase(),
			s.useCaseRegistry.NewCreateTagUseCase(),
			s.useCaseRegistry.NewGetAllTagsUseCase(),
			s.useCaseRegistry.NewDeleteTagUseCase(),
		)
		r.Post("/api/recipes", recipeHandler.CreateRecipe)
		r.Get("/api/recipes", recipeHandler.GetUserRecipes)
		r.Get("/api/recipes/{id}", recipeHandler.GetRecipe)
		r.Get("/api/users/{user_id}/recipes", recipeHandler.GetUserRecipes)
		r.Put("/api/recipes/{id}", recipeHandler.UpdateRecipe)
		r.Delete("/api/recipes/{id}", recipeHandler.DeleteRecipe)
		r.Post("/api/recipes/search", recipeHandler.SearchRecipes)

		// Recipe tags
		r.Post("/api/recipes/{id}/tags", recipeHandler.AddTags)
		r.Delete("/api/recipes/{id}/tags", recipeHandler.RemoveTags)

		// Recipe images
		r.Post("/api/recipes/{id}/images", recipeHandler.AddImage)

		// Tags
		r.Post("/api/tags", recipeHandler.CreateTag)
		r.Get("/api/tags", recipeHandler.GetAllTags)
		r.Delete("/api/tags/{id}", recipeHandler.DeleteTag)
	})
}

// Run starts the HTTP server (blocking)
func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	log.Printf("Server starting on %s\n", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// ServeHTTP implements http.Handler for testing
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Shutdown performs graceful shutdown of the server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exited")
	return nil
}

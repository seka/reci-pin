package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	repoRegistry    registry.Repository
	useCaseRegistry registry.UseCase
	cfg             *config.Config
}

// New creates a new server with all dependencies initialized
func New(cfg *config.Config) (*Server, error) {
	// Initialize Repository Registry
	repoRegistry, err := registry.NewRepository(context.Background(), cfg.Database.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository registry: %w", err)
	}

	log.Println("Successfully connected to database")

	// Initialize UseCase Registry
	useCaseRegistry := registry.NewUseCase(repoRegistry, cfg)

	s := &Server{
		router:          chi.NewRouter(),
		repoRegistry:    repoRegistry,
		useCaseRegistry: useCaseRegistry,
		cfg:             cfg,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s, nil
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
		w.Write([]byte("OK"))
	})

	// Auth public routes
	authHandler := handler.NewAuthHandler(
		s.useCaseRegistry.NewSignupUseCase(),
		s.useCaseRegistry.NewLoginUseCase(),
		s.useCaseRegistry.NewGenerateTokenUseCase(),
		s.useCaseRegistry.NewGetUserUseCase(),
		s.useCaseRegistry.NewVerifyEmailUseCase(),
	)

	s.router.Post("/api/auth/signup", authHandler.Signup)
	s.router.Post("/api/auth/login", authHandler.Login)
	s.router.Post("/api/auth/verify", authHandler.Verify)

	s.router.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)

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

// Run starts the HTTP server and handles graceful shutdown
func (s *Server) Run() error {
	defer s.repoRegistry.Close()

	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	go func() {
		log.Printf("Server starting on %s\n", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	return s.Shutdown()
}

// Shutdown performs graceful shutdown of the server
func (s *Server) Shutdown() error {
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exited")
	return nil
}

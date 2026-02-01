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
	router *chi.Mux
}

// New creates a new server with routing configured
func New(
	authHandler *handler.AuthHandler,
	recipeHandler *handler.RecipeHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Server {
	s := &Server{
		router: chi.NewRouter(),
	}
	s.setupMiddleware()
	s.setupRoutes(authHandler, recipeHandler, authMiddleware)
	return s
}

func (s *Server) setupMiddleware() {
	s.router.Use(chimiddleware.RequestID)
	s.router.Use(chimiddleware.RealIP)
	s.router.Use(chimiddleware.Logger)
	s.router.Use(chimiddleware.Recoverer)
	s.router.Use(chimiddleware.Timeout(60 * time.Second))

	// CORS middleware
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})
}

func (s *Server) setupRoutes(
	authHandler *handler.AuthHandler,
	recipeHandler *handler.RecipeHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	// Public routes
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Auth public routes
	s.router.Post("/api/auth/signup", authHandler.Signup)
	s.router.Post("/api/auth/login", authHandler.Login)
	s.router.Post("/api/auth/verify", authHandler.Verify)

	// Protected routes
	s.router.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)

		// Recipes
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

// ServeHTTP allows Server to implement http.Handler for testing
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Run initializes all dependencies and starts the server
func Run(ctx context.Context, cfg *config.Config) error {
	// Initialize Repository Registry
	repoRegistry, err := registry.NewRepository(ctx, cfg.Database.DSN())
	if err != nil {
		return fmt.Errorf("failed to initialize repository registry: %w", err)
	}
	defer repoRegistry.Close()

	log.Println("Successfully connected to database")

	// Initialize UseCase Registry
	useCaseRegistry := registry.NewUseCase(repoRegistry, cfg)

	// Create Handlers
	authHandler := handler.NewAuthHandler(
		useCaseRegistry.NewSignupUseCase(),
		useCaseRegistry.NewLoginUseCase(),
		useCaseRegistry.NewGenerateTokenUseCase(),
		useCaseRegistry.NewGetUserUseCase(),
		useCaseRegistry.NewVerifyEmailUseCase(),
	)

	recipeHandler := handler.NewRecipeHandler(
		useCaseRegistry.NewCreateRecipeUseCase(),
		useCaseRegistry.NewGetRecipeUseCase(),
		useCaseRegistry.NewGetUserRecipesUseCase(),
		useCaseRegistry.NewUpdateRecipeUseCase(),
		useCaseRegistry.NewDeleteRecipeUseCase(),
		useCaseRegistry.NewSearchRecipesUseCase(),
		useCaseRegistry.NewAddTagsUseCase(),
		useCaseRegistry.NewRemoveTagsUseCase(),
		useCaseRegistry.NewAddImageUseCase(),
		useCaseRegistry.NewCreateTagUseCase(),
		useCaseRegistry.NewGetAllTagsUseCase(),
		useCaseRegistry.NewDeleteTagUseCase(),
	)

	authMiddleware := middleware.NewAuthMiddleware(useCaseRegistry.NewValidateTokenUseCase())

	// Create server
	server := New(authHandler, recipeHandler, authMiddleware)

	// Start HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: server.router,
	}

	startServer(httpServer)
	return gracefulShutdown(httpServer)
}

// startServer starts the HTTP server in a goroutine
func startServer(httpServer *http.Server) {
	go func() {
		log.Printf("Server starting on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

// gracefulShutdown waits for interrupt signal and performs graceful shutdown
func gracefulShutdown(httpServer *http.Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exiting")
	return nil
}

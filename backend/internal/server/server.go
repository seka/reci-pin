package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/infrastructure/datastore/postgres"
	"github.com/seka/reci-pin/backend/internal/server/handler"
	"github.com/seka/reci-pin/backend/internal/server/middleware"
	authUC "github.com/seka/reci-pin/backend/internal/usecase/auth"
	recipeUC "github.com/seka/reci-pin/backend/internal/usecase/recipe"
	recipeImageUC "github.com/seka/reci-pin/backend/internal/usecase/recipe_image"
	recipeTagUC "github.com/seka/reci-pin/backend/internal/usecase/recipe_tag"
	tagUC "github.com/seka/reci-pin/backend/internal/usecase/tag"
)

type Server struct {
	router         *chi.Mux
	authHandler    *handler.AuthHandler
	recipeHandler  *handler.RecipeHandler
	authMiddleware *middleware.AuthMiddleware
}

func Run(ctx context.Context, cfg *config.Config) error {
	// Initialize database
	db, err := postgres.New(ctx, cfg.Database.DSN())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	log.Println("Successfully connected to database")

	// Repositories
	userRepo := postgres.NewUserRepository(db)
	recipeRepo := postgres.NewRecipeRepository(db)
	tagRepo := postgres.NewTagRepository(db)
	recipeImageRepo := postgres.NewRecipeImageRepository(db)
	credentialRepo := postgres.NewUserEmailCredentialRepository(db)

	// UseCases - Auth
	signupUseCase := authUC.NewSignupUseCase(userRepo, credentialRepo)
	loginUseCase := authUC.NewLoginUseCase(credentialRepo)
	generateTokenUseCase := authUC.NewGenerateTokenUseCase(cfg.JWT.Secret, time.Duration(cfg.JWT.ExpirationHours)*time.Hour)
	getUserUseCase := authUC.NewGetUserUseCase(userRepo)
	validateTokenUseCase := authUC.NewValidateTokenUseCase(cfg.JWT.Secret)
	verifyEmailUseCase := authUC.NewVerifyEmailUseCase(credentialRepo)

	// UseCases - Recipe
	createRecipeUseCase := recipeUC.NewCreateRecipeUseCase(recipeRepo)
	getRecipeUseCase := recipeUC.NewGetRecipeUseCase(recipeRepo, recipeImageRepo)
	getUserRecipesUseCase := recipeUC.NewGetUserRecipesUseCase(recipeRepo, recipeImageRepo)
	updateRecipeUseCase := recipeUC.NewUpdateRecipeUseCase(recipeRepo)
	deleteRecipeUseCase := recipeUC.NewDeleteRecipeUseCase(recipeRepo)
	searchRecipesUseCase := recipeUC.NewSearchRecipesUseCase(recipeRepo, recipeImageRepo)
	addImageUseCase := recipeImageUC.NewAddImageUseCase(recipeRepo, recipeImageRepo)

	// Usecases - Tag
	createTagUseCase := tagUC.NewCreateTagUseCase(tagRepo)
	getAllTagsUseCase := tagUC.NewGetAllTagsUseCase(tagRepo)
	deleteTagUseCase := tagUC.NewDeleteTagUseCase(tagRepo)
	addTagsUseCase := recipeTagUC.NewAddTagsUseCase(recipeRepo)
	removeTagsUseCase := recipeTagUC.NewRemoveTagsUseCase(recipeRepo)

	s := &Server{
		router: chi.NewRouter(),
		authHandler: handler.NewAuthHandler(
			signupUseCase,
			loginUseCase,
			generateTokenUseCase,
			getUserUseCase,
			verifyEmailUseCase,
		),
		recipeHandler: handler.NewRecipeHandler(
			createRecipeUseCase,
			getRecipeUseCase,
			getUserRecipesUseCase,
			updateRecipeUseCase,
			deleteRecipeUseCase,
			searchRecipesUseCase,
			addTagsUseCase,
			removeTagsUseCase,
			addImageUseCase,
			createTagUseCase,
			getAllTagsUseCase,
			deleteTagUseCase,
		),
		authMiddleware: middleware.NewAuthMiddleware(validateTokenUseCase),
	}

	s.setupMiddleware()
	s.setupRoutes()

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: s,
	}

	go func() {
		log.Printf("Server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server listening error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Println("Server shutdown successfully")
	return nil
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

func (s *Server) setupRoutes() {
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Public routes
	s.router.Post("/auth/signup", s.authHandler.Signup)
	s.router.Post("/auth/login", s.authHandler.Login)
	s.router.Post("/auth/verify", s.authHandler.Verify)

	// Protected routes
	s.router.Group(func(r chi.Router) {
		r.Use(s.authMiddleware.Authenticate)

		// Auth
		r.Post("/auth/logout", s.authHandler.Logout)

		// Recipes
		r.Post("/recipes", s.recipeHandler.CreateRecipe)
		r.Get("/recipes", s.recipeHandler.GetUserRecipes)
		r.Post("/recipes/search", s.recipeHandler.SearchRecipes)
		r.Get("/recipes/{id}", s.recipeHandler.GetRecipe)
		r.Put("/recipes/{id}", s.recipeHandler.UpdateRecipe)
		r.Delete("/recipes/{id}", s.recipeHandler.DeleteRecipe)

		// Recipe tags
		r.Post("/recipes/{id}/tags", s.recipeHandler.AddTags)
		r.Delete("/recipes/{id}/tags", s.recipeHandler.RemoveTags)

		// Recipe images
		r.Post("/recipes/{id}/images", s.recipeHandler.AddImage)

		// Tags
		r.Post("/tags", s.recipeHandler.CreateTag)
		r.Get("/tags", s.recipeHandler.GetAllTags)
		r.Delete("/tags/{id}", s.recipeHandler.DeleteTag)
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

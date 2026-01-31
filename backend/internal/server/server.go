package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
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

func New(
	signupUseCase *authUC.SignupUseCase,
	loginUseCase *authUC.LoginUseCase,
	generateTokenUseCase *authUC.GenerateTokenUseCase,
	validateTokenUseCase *authUC.ValidateTokenUseCase,
	getUserUseCase *authUC.GetUserUseCase,
	verifyEmailUseCase *authUC.VerifyEmailUseCase, // Added
	createRecipeUseCase *recipeUC.CreateRecipeUseCase,
	getRecipeUseCase *recipeUC.GetRecipeUseCase,
	getUserRecipesUseCase *recipeUC.GetUserRecipesUseCase,
	updateRecipeUseCase *recipeUC.UpdateRecipeUseCase,
	deleteRecipeUseCase *recipeUC.DeleteRecipeUseCase,
	searchRecipesUseCase *recipeUC.SearchRecipesUseCase,
	addTagsUseCase *recipeTagUC.AddTagsUseCase,
	removeTagsUseCase *recipeTagUC.RemoveTagsUseCase,
	addImageUseCase *recipeImageUC.AddImageUseCase,
	createTagUseCase *tagUC.CreateTagUseCase,
	getAllTagsUseCase *tagUC.GetAllTagsUseCase,
	deleteTagUseCase *tagUC.DeleteTagUseCase,
) *Server {
	s := &Server{
		router: chi.NewRouter(),
		authHandler: handler.NewAuthHandler(
			signupUseCase,
			loginUseCase,
			generateTokenUseCase,
			getUserUseCase,
			verifyEmailUseCase, // Added
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

func (s *Server) setupRoutes() {
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Public routes
	s.router.Post("/auth/signup", s.authHandler.Signup)
	s.router.Post("/auth/login", s.authHandler.Login)
	s.router.Post("/auth/verify", s.authHandler.Verify) // Added

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

func (s *Server) Start(port int) error {
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server starting on %s\n", addr)
	return http.ListenAndServe(addr, s.router)
}

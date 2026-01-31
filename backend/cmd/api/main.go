package main

import (
	"context"
	"fmt"
	"log"

	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/infrastructure/datastore/postgres"
	"github.com/seka/reci-pin/backend/internal/server"
	authUC "github.com/seka/reci-pin/backend/internal/usecase/auth"
	recipeUC "github.com/seka/reci-pin/backend/internal/usecase/recipe"
	recipeImageUC "github.com/seka/reci-pin/backend/internal/usecase/recipe_image"
	recipeTagUC "github.com/seka/reci-pin/backend/internal/usecase/recipe_tag"
	tagUC "github.com/seka/reci-pin/backend/internal/usecase/tag"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	ctx := context.Background()
	db, err := postgres.New(ctx, cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to database")

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	recipeRepo := postgres.NewRecipeRepository(db)
	tagRepo := postgres.NewTagRepository(db)
	recipeImageRepo := postgres.NewRecipeImageRepository(db)

	// Initialize auth use cases
	signupUseCase := authUC.NewSignupUseCase(userRepo)
	loginUseCase := authUC.NewLoginUseCase(userRepo)
	generateTokenUseCase := authUC.NewGenerateTokenUseCase(cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	validateTokenUseCase := authUC.NewValidateTokenUseCase(cfg.JWT.Secret)
	getUserUseCase := authUC.NewGetUserUseCase(userRepo)

	// Initialize recipe use cases
	createRecipeUseCase := recipeUC.NewCreateRecipeUseCase(recipeRepo)
	getRecipeUseCase := recipeUC.NewGetRecipeUseCase(recipeRepo, recipeImageRepo)
	getUserRecipesUseCase := recipeUC.NewGetUserRecipesUseCase(recipeRepo, recipeImageRepo)
	updateRecipeUseCase := recipeUC.NewUpdateRecipeUseCase(recipeRepo)
	deleteRecipeUseCase := recipeUC.NewDeleteRecipeUseCase(recipeRepo)
	searchRecipesUseCase := recipeUC.NewSearchRecipesUseCase(recipeRepo, recipeImageRepo)

	// Initialize tag use cases
	createTagUseCase := tagUC.NewCreateTagUseCase(tagRepo)
	getAllTagsUseCase := tagUC.NewGetAllTagsUseCase(tagRepo)
	deleteTagUseCase := tagUC.NewDeleteTagUseCase(tagRepo)

	// Initialize recipe_tag use cases
	addTagsUseCase := recipeTagUC.NewAddTagsUseCase(recipeRepo)
	removeTagsUseCase := recipeTagUC.NewRemoveTagsUseCase(recipeRepo)

	// Initialize recipe_image use cases
	addImageUseCase := recipeImageUC.NewAddImageUseCase(recipeRepo, recipeImageRepo)

	// Initialize server
	srv := server.New(
		signupUseCase,
		loginUseCase,
		generateTokenUseCase,
		validateTokenUseCase,
		getUserUseCase,
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
	)

	// Start server
	if err := srv.Start(cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

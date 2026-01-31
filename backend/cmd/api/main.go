package main

import (
	"context"
	"fmt"
	"log"

	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/infrastructure/datastore/postgres"
	"github.com/seka/reci-pin/backend/internal/server"
	"github.com/seka/reci-pin/backend/internal/usecase"
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

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	recipeUseCase := usecase.NewRecipeUseCase(recipeRepo, tagRepo, recipeImageRepo)

	// Initialize server
	srv := server.New(authUseCase, recipeUseCase)

	// Start server
	if err := srv.Start(cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

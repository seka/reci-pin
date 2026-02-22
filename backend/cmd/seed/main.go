package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"path"

	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database/postgres"
	"github.com/seka/reci-pin/backend/internal/registry"
)

var (
	cfg     config.Config
	doClean bool
)

func init() {
	flag.StringVar(&cfg.Database.Host, "db-host", "localhost", "Database host")
	flag.IntVar(&cfg.Database.Port, "db-port", 5432, "Database port")
	flag.StringVar(&cfg.Database.User, "db-user", "postgres", "Database user")
	flag.StringVar(&cfg.Database.Password, "db-password", "postgres", "Database password")
	flag.StringVar(&cfg.Database.DBName, "db-name", "recipin_dev", "Database name")
	flag.StringVar(&cfg.Database.SSLMode, "db-sslmode", "disable", "Database SSL mode")
	flag.BoolVar(&doClean, "clean", false, "Clean existing data before seeing")
}

func main() {
	flag.Parse()

	ctx := context.Background()

	// Connect to Database
	db := postgres.New(cfg.Database.DSN())
	if err := db.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to database")

	// Clean data if requested
	if doClean {
		log.Println("Cleaning existing data...")
		if err := cleanData(ctx, db); err != nil {
			log.Fatalf("Failed to clean data: %v", err)
		}
	}

	// Initialize Registry
	repoReg := registry.NewRepository(db, nil)

	// Seed Data
	if err := seedData(ctx, repoReg); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}

	log.Println("Seeding completed successfully")
}

func cleanData(ctx context.Context, db database.Database) error {
	// Execute raw SQL to truncate tables
	queries := []string{
		"TRUNCATE TABLE recipe_tags RESTART IDENTITY CASCADE",
		"TRUNCATE TABLE recipe_images RESTART IDENTITY CASCADE",
		"TRUNCATE TABLE recipes RESTART IDENTITY CASCADE",
		"TRUNCATE TABLE tags RESTART IDENTITY CASCADE",
		"TRUNCATE TABLE user_email_credentials RESTART IDENTITY CASCADE",
		"TRUNCATE TABLE users RESTART IDENTITY CASCADE",
	}

	for _, query := range queries {
		if _, err := db.Execute(ctx, query); err != nil {
			// Some tables might not exist or be empty, but TRUNCATE usually works if table exists.
			// Ignoring errors for simplicity? No, better warn.
			log.Printf("Warning: failed to execute %s: %v", query, err)
		}
	}

	return nil
}

func seedData(ctx context.Context, repoReg registry.Repository) error {
	// 1. Create Users
	if err := createUsers(ctx, repoReg); err != nil {
		return err
	}
	// 2. Create Tags
	if err := createTags(ctx, repoReg); err != nil {
		return err
	}
	// 3. Create Recipes
	if err := createRecipes(ctx, repoReg); err != nil {
		return err
	}
	// 4. Assign Tags to Recipes
	if err := assignTagsToRecipes(ctx, repoReg); err != nil {
		return err
	}
	// 5. Create Recipe Images
	if err := createRecipeImages(ctx, repoReg); err != nil {
		return err
	}
	return nil
}

func createUsers(ctx context.Context, repoReg registry.Repository) error {
	userRepo := repoReg.NewUserRepository()
	credRepo := repoReg.NewUserEmailCredentialRepository()

	users := []struct {
		Name     string
		Email    string
		Password string
	}{
		{"John", "john@example.com", "password"},
		{"Jane", "jane@example.com", "password"},
	}

	for _, u := range users {
		// Create User
		user := &model.User{Name: u.Name}
		if err := userRepo.Create(ctx, user); err != nil {
			return fmt.Errorf("creating user %s: %w", u.Name, err)
		}

		// Create Credential
		if err := createCredential(ctx, credRepo, user.ID, u.Email, u.Password); err != nil {
			return err
		}

		log.Printf("Created user: %s (%s)", u.Name, u.Email)
	}
	return nil
}

func createCredential(ctx context.Context, credRepo repository.UserEmailCredentialRepository, userID int64, email, password string) error {
	hashed, err := postgres.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	now := time.Now()
	cred := &model.UserEmailCredential{
		UserID:          userID,
		Email:           email,
		PasswordHash:    hashed,
		EmailVerifiedAt: &now,
	}
	if err := credRepo.Create(ctx, cred); err != nil {
		return fmt.Errorf("creating credential for %s: %w", email, err)
	}
	return nil
}

func createRecipes(ctx context.Context, repoReg registry.Repository) error {
	userRepo := repoReg.NewUserRepository()
	recipeRepo := repoReg.NewRecipeRepository()

	users, err := userRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("getting users: %w", err)
	}

	for i, user := range users {
		recipe := &model.Recipe{
			UserID: user.ID,
			Name:   fmt.Sprintf("Delicious Fish %d", i),
			Memo:   "Freshly caught fish recipe.",
			URL:    fmt.Sprintf("http://example.com/recipes/%d", i),
		}
		if err := recipeRepo.Create(ctx, recipe); err != nil {
			return fmt.Errorf("creating recipe %d: %w", i, err)
		}
		log.Printf("Created recipe: %s", recipe.Name)
	}
	return nil
}

func createTags(ctx context.Context, repoReg registry.Repository) error {
	tagRepo := repoReg.NewTagRepository()

	tagNames := []string{"魚料理", "簡単", "和食", "洋食", "煮物", "焼き物"}

	for _, name := range tagNames {
		tag := &model.Tag{Name: name}
		if err := tagRepo.Create(ctx, tag); err != nil {
			return fmt.Errorf("creating tag %s: %w", name, err)
		}
		log.Printf("Created tag: %s", name)
	}
	return nil
}
func assignTagsToRecipes(ctx context.Context, repoReg registry.Repository) error {
	recipeRepo := repoReg.NewRecipeRepository()
	tagRepo := repoReg.NewTagRepository()

	// Get all recipes and tags
	recipes, err := recipeRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("getting all recipes: %w", err)
	}

	tags, err := tagRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("getting all tags: %w", err)
	}

	if len(tags) == 0 {
		log.Println("No tags available to assign")
		return nil
	}

	// Assign 2-3 random tags to each recipe
	for _, recipe := range recipes {
		numTags := 2 + (int(recipe.ID) % 2) // 2 or 3 tags
		var tagIDs []int64

		// Select random tags
		for i := 0; i < numTags && i < len(tags); i++ {
			tagIdx := (int(recipe.ID) + i) % len(tags)
			tagIDs = append(tagIDs, tags[tagIdx].ID)
		}

		if err := recipeRepo.AddTags(ctx, recipe.ID, tagIDs); err != nil {
			return fmt.Errorf("assigning tags to recipe %d: %w", recipe.ID, err)
		}
		log.Printf("Assigned %d tags to recipe: %s", len(tagIDs), recipe.Name)
	}
	return nil
}

func createRecipeImages(ctx context.Context, repoReg registry.Repository) error {
	recipeRepo := repoReg.NewRecipeRepository()
	imageRepo := repoReg.NewRecipeImageRepository()

	// Get all recipes
	recipes, err := recipeRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("getting all recipes: %w", err)
	}

	// Add 1-2 images per recipe
	for i, recipe := range recipes {
		numImages := (i % 2) + 1 // 1 or 2 images
		for j := range numImages {
			image := &model.RecipeImage{
				RecipeID:  recipe.ID,
				ImagePath: path.Join("recipes", fmt.Sprintf("%d", recipe.ID), fmt.Sprintf("seed_%d.jpg", j+1)),
			}
			if err := imageRepo.Create(ctx, image); err != nil {
				return fmt.Errorf("creating image for recipe %d: %w", recipe.ID, err)
			}
			log.Printf("  Created image: %s", image.ImagePath)
		}
	}
	return nil
}

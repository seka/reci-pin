package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/domain/storage"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database/postgres"
	"github.com/seka/reci-pin/backend/internal/infrastructure/storage/s3"
	"github.com/seka/reci-pin/backend/internal/registry"
)

var (
	cfg     config.Config
	doClean bool
)

func init() {
	// DB
	flag.StringVar(&cfg.Database.Host, "db-host", "localhost", "Database host")
	flag.StringVar(&cfg.Database.Port, "db-port", "5432", "Database port")
	flag.StringVar(&cfg.Database.User, "db-user", "postgres", "Database user")
	flag.StringVar(&cfg.Database.Password, "db-password", "postgres", "Database password")
	flag.StringVar(&cfg.Database.DBName, "db-name", "recipin_dev", "Database name")
	flag.StringVar(&cfg.Database.SSLMode, "db-sslmode", "disable", "Database SSL mode")

	// Storage
	flag.StringVar(&cfg.Storage.Bucket, "storage-bucket", "recipin-bucket", "Storage bucket name")
	flag.StringVar(&cfg.Storage.Endpoint, "storage-endpoint", "", "Storage endpoint (e.g. LocalStack)")
	flag.StringVar(&cfg.Storage.Region, "storage-region", "us-east-1", "Storage region")
	flag.StringVar(&cfg.Storage.AccessKey, "storage-access-key", "", "Storage access key")
	flag.StringVar(&cfg.Storage.SecretKey, "storage-secret-key", "", "Storage secret key")
	flag.StringVar(&cfg.Storage.PublicBaseURL, "storage-public-url", "", "Storage public base URL")

	flag.BoolVar(&doClean, "clean", false, "Clean existing data before seeing")
}

func main() {
	flag.Parse()

	ctx := context.Background()

	// Connect to Database
	db := postgres.NewClient(cfg.Database)
	if err := db.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to database")

	// Connect to Storage
	storageClient, err := s3.NewClient(ctx, cfg.Storage)
	if err != nil {
		log.Fatalf("Failed to initialize storage client: %v", err)
	}

	// Clean data if requested
	if doClean {
		log.Println("Cleaning existing data...")
		if err := cleanData(ctx, db); err != nil {
			log.Fatalf("Failed to clean data: %v", err)
		}
	}

	// Initialize Registry
	repoReg := registry.NewRepository(db)

	// Seed Data
	if err := seedData(ctx, repoReg, storageClient); err != nil {
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

func seedData(ctx context.Context, repoReg registry.Repository, storageClient storage.Client) error {
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
	if err := createRecipeImages(ctx, repoReg, storageClient); err != nil {
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
		recipePath, err := url.JoinPath("recipes", strconv.Itoa(i))
		if err != nil {
			return fmt.Errorf("creating recipe path: %w", err)
		}
		recipeURL := (&url.URL{
			Scheme: "http",
			Host:   "example.com",
			Path:   recipePath,
		}).String()

		recipe := &model.Recipe{
			UserID: user.ID,
			Name:   fmt.Sprintf("Delicious Fish %d", i),
			Memo:   "Freshly caught fish recipe.",
			URL:    recipeURL,
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

func createRecipeImages(ctx context.Context, repoReg registry.Repository, storageClient storage.Client) error {
	recipeRepo := repoReg.NewRecipeRepository()
	imageRepo := repoReg.NewRecipeImageRepository()

	// Get all recipes
	recipes, err := recipeRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("getting all recipes: %w", err)
	}

	// Read placeholder image
	placeholderPath := filepath.Join("internal", "infrastructure", "storage", "mock", "seed_placeholder.png")
	// Check relative paths based on where it's run. If run from backend/ it works.
	// In Docker, it's run from /app which is backend/.
	imageData, err := os.ReadFile(placeholderPath)
	if err != nil {
		return fmt.Errorf("reading placeholder image: %w", err)
	}

	// Add 1-2 images per recipe
	for i, recipe := range recipes {
		numImages := (i % 2) + 1 // 1 or 2 images
		for j := 0; j < numImages; j++ {
			fileName := fmt.Sprintf("seed_%d.png", j+1)
			imageKey, err := url.JoinPath("recipes", strconv.FormatInt(recipe.ID, 10), fileName)
			if err != nil {
				return fmt.Errorf("joining image path: %w", err)
			}

			// Upload to Storage
			if err := storageClient.Upload(ctx, imageKey, bytes.NewReader(imageData), model.RecipeImageTypePNG); err != nil {
				return fmt.Errorf("uploading image %s: %w", imageKey, err)
			}

			image := &model.RecipeImage{
				RecipeID:  recipe.ID,
				ImagePath: imageKey,
			}
			if err := imageRepo.Create(ctx, image); err != nil {
				return fmt.Errorf("creating image for recipe %d: %w", recipe.ID, err)
			}
			log.Printf("  Created and uploaded image: %s", image.ImagePath)
		}
	}
	return nil
}

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"sync"
	"time"

	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/domain/model"
	es "github.com/seka/reci-pin/backend/internal/infrastructure/database/elasticsearch"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database/postgres"
	"github.com/seka/reci-pin/backend/internal/registry"
)

var (
	cfg config.Config
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func init() {
	flag.StringVar(&cfg.Database.Host, "db-host", getEnv("DB_HOST", "localhost"), "Database host")
	flag.IntVar(&cfg.Database.Port, "db-port", 5432, "Database port")
	flag.StringVar(&cfg.Database.User, "db-user", getEnv("DB_USER", "postgres"), "Database user")
	flag.StringVar(&cfg.Database.Password, "db-password", getEnv("DB_PASSWORD", "postgres"), "Database password")
	flag.StringVar(&cfg.Database.DBName, "db-name", getEnv("DB_NAME", "recipin_dev"), "Database name")
	flag.StringVar(&cfg.Database.SSLMode, "db-sslmode", getEnv("DB_SSLMODE", "disable"), "Database SSL mode")
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

	// Connect to Elasticsearch
	esClient, err := es.NewClient()
	if err != nil {
		log.Fatalf("Failed to connect to elasticsearch: %v", err)
	}
	log.Println("Connected to elasticsearch")

	// Initialize Registry
	repoReg := registry.NewRepository(db, esClient)
	recipeRepo := repoReg.NewRecipeRepository()
	searchRepo := repoReg.NewRecipeSearchRepository()

	// 1. Get All Recipes
	log.Println("Fetching all recipes...")
	recipes, err := recipeRepo.GetAll(ctx)
	if err != nil {
		log.Fatalf("Failed to get recipes: %v", err)
	}
	log.Printf("Found %d recipes", len(recipes))

	// 2. Index each recipe
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10) // Limit concurrency
	// errCount := 0
	// successCount := 0

	startTime := time.Now()

	for _, recipe := range recipes {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(r model.Recipe) {
			defer wg.Done()
			defer func() { <-semaphore }()

			// Fetch tags
			tags, err := recipeRepo.GetTags(ctx, r.ID)
			if err != nil {
				log.Printf("Error fetching tags for recipe %d: %v", r.ID, err)
				// Continue indexing without tags? Or count as error?
				// Let's count as error for now to be safe
				// mu.Lock(); errCount++; mu.Unlock()
				// return
			}
			r.Tags = tags

			if err := searchRepo.Index(ctx, &r); err != nil {
				log.Printf("Error indexing recipe %d: %v", r.ID, err)
				// mu.Lock(); errCount++; mu.Unlock()
			}
		}(recipe)
	}

	wg.Wait()
	duration := time.Since(startTime)

	log.Printf("Sync completed in %v", duration)
	// log.Printf("Success: %d, Failed: %d", successCount, errCount) // Need mutex for counters
}

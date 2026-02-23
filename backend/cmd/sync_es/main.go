package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/domain/model"
	postgres "github.com/seka/reci-pin/backend/internal/infrastructure/database/postgres"
	es "github.com/seka/reci-pin/backend/internal/infrastructure/searchengine/elasticsearch"
	"github.com/seka/reci-pin/backend/internal/registry"
)

var (
	cfg    config.Config
	esHost string
	esPort string
)

func init() {
	flag.StringVar(&cfg.Database.Host, "db-host", "localhost", "Database host")
	flag.StringVar(&cfg.Database.Port, "db-port", "5432", "Database port")
	flag.StringVar(&cfg.Database.User, "db-user", "postgres", "Database user")
	flag.StringVar(&cfg.Database.Password, "db-password", "postgres", "Database password")
	flag.StringVar(&cfg.Database.DBName, "db-name", "recipin_dev", "Database name")
	flag.StringVar(&cfg.Database.SSLMode, "db-sslmode", "disable", "Database SSL mode")

	// Search Engine configuration
	flag.StringVar(&esHost, "es-host", "localhost", "Elasticsearch host")
	flag.StringVar(&esPort, "es-port", "9200", "Elasticsearch port")
	cfg.SearchEngine.Addresses = []string{}
	flag.Func("es-addresses", "Elasticsearch addresses (comma separated)", func(v string) error {
		cfg.SearchEngine.Addresses = parseAddresses(v)
		return nil
	})

	// Storage configuration
	flag.StringVar(&cfg.Storage.Bucket, "storage-bucket", "recipin-bucket", "S3 bucket name")
	flag.StringVar(&cfg.Storage.Endpoint, "storage-endpoint", "", "S3 endpoint URL (for LocalStack)")
	flag.StringVar(&cfg.Storage.PublicBaseURL, "storage-public-url", "", "Base URL for public access")
}

func main() {
	flag.Parse()

	if len(cfg.SearchEngine.Addresses) == 0 {
		cfg.SearchEngine.Addresses = []string{"http://" + esHost + ":" + esPort}
	}

	ctx := context.Background()

	// Connect to Database
	db := postgres.NewClient(cfg.Database)
	if err := db.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to database")

	// Start SearchEngine
	esClient, err := es.NewClient(cfg.SearchEngine)
	if err != nil {
		log.Fatalf("Failed to connect to elasticsearch: %v", err)
	}
	log.Println("Connected to elasticsearch")

	// Initialize Registry
	repoReg := registry.NewRepository(db)
	searcherReg := registry.NewSearcher(esClient)
	recipeRepo := repoReg.NewRecipeRepository()
	searchRepo := searcherReg.NewRecipeSearchRepository()

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

func parseAddresses(raw string) []string {
	parts := strings.Split(raw, ",")
	addresses := make([]string, 0, len(parts))
	for _, part := range parts {
		addr := strings.TrimSpace(part)
		if addr == "" {
			continue
		}
		addresses = append(addresses, addr)
	}
	return addresses
}

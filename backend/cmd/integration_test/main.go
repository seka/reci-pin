package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"

	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database"
	postgres "github.com/seka/reci-pin/backend/internal/infrastructure/database/postgres"
	"github.com/seka/reci-pin/backend/internal/infrastructure/notification/mailhog"
	es "github.com/seka/reci-pin/backend/internal/infrastructure/searchengine/elasticsearch"
	"github.com/seka/reci-pin/backend/internal/infrastructure/storage/s3"
	"github.com/seka/reci-pin/backend/internal/registry"
	"github.com/seka/reci-pin/backend/internal/server"
)

var (
	cfg config.Config
)

func init() {
	flag.StringVar(&cfg.Database.Port, "db-port", "5432", "Database port")
	flag.StringVar(&cfg.Database.User, "db-user", "postgres", "Database user")
	flag.StringVar(&cfg.Database.Password, "db-password", "postgres", "Database password")
	flag.StringVar(&cfg.Database.SSLMode, "db-sslmode", "disable", "Database SSL mode")
	// Note: db-name will be generated dynamically
	flag.StringVar(&cfg.ApiServer.Port, "port", "0", "Server port (0 for random)")
	flag.StringVar(&cfg.ApiServer.JWT.Secret, "jwt-secret", "test-secret", "JWT secret")
	cfg.ApiServer.JWT.ExpirationHours = 24

	// Initialize other configs with defaults
	cfg.Storage.Bucket = "test-bucket"
	cfg.Storage.Endpoint = "http://localhost:4566"
	cfg.Storage.PublicBaseURL = "http://localhost:4566/test-bucket"
	cfg.SearchEngine.Addresses = []string{"http://localhost:9200"}
	cfg.Email.Host = "localhost"
	cfg.Email.Port = "1025"
	cfg.Email.From = "no-reply@reci-pin.com"
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("Integration test failed: %v", err)
	}
	log.Println("Integration test passed successfully!")
}

func run() error {
	ctx := context.Background()

	// 1. Setup Test Database
	testDBName := fmt.Sprintf("recipin_test_%d", time.Now().UnixNano())
	log.Printf("Setting up test database: %s", testDBName)

	if err := createDatabase(ctx, testDBName); err != nil {
		return fmt.Errorf("creating test database: %w", err)
	}
	defer func() {
		if err := dropDatabase(ctx, testDBName); err != nil {
			log.Printf("Failed to drop test database: %v", err)
		} else {
			log.Println("Test database dropped")
		}
	}()

	// 2. Start Server
	// Override DB Name in config
	testCfg := cfg
	testCfg.Database.DBName = testDBName
	// Set specific port for test, e.g. 8081
	testCfg.ApiServer.Port = "8081"

	// Connect to the NEW test database
	db := postgres.NewClient(testCfg.Database)
	if err := db.Connect(ctx); err != nil {
		return fmt.Errorf("connecting to test database: %w", err)
	}
	defer db.Close()

	// Connect to Elasticsearch
	esClient, err := es.NewClient(testCfg.SearchEngine)
	if err != nil {
		return fmt.Errorf("creating elasticsearch client: %w", err)
	}

	// Initialize Storage Service
	storageService, err := s3.NewClient(ctx, testCfg.Storage)
	if err != nil {
		return fmt.Errorf("creating storage service: %w", err)
	}

	// Initialize Server
	mailClient := mailhog.NewClient(testCfg.Email)
	repoReg := registry.NewRepository(db)
	searcherReg := registry.NewSearcher(esClient)
	useCaseReg := registry.NewUseCase(repoReg, storageService, searcherReg, mailClient, &testCfg)
	srv := server.New(&testCfg.ApiServer, useCaseReg)

	// Run Server in Goroutine
	serverErrCh := make(chan error, 1)
	go func() {
		log.Printf("Starting test server on port %s", testCfg.ApiServer.Port)
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			serverErrCh <- err
		}
	}()

	// Wait for server to be ready
	// Health check is at root /health, not /api/health
	healthPath, err := url.JoinPath("health")
	if err != nil {
		return fmt.Errorf("creating health path: %w", err)
	}
	healthURL := &url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort("localhost", testCfg.ApiServer.Port),
		Path:   healthPath,
	}
	if err := waitForServer(healthURL); err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}

	apiPath, err := url.JoinPath("api")
	if err != nil {
		return fmt.Errorf("creating api path: %w", err)
	}
	baseURL := &url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort("localhost", testCfg.ApiServer.Port),
		Path:   apiPath,
	}

	// 3. Run Scenarios
	log.Println("Starting api scenarios...")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	// Use CookieJar if needed, but for now assuming JWT in header or cookie?
	// If cookie, need jar.
	jar, _ := cookiejar.New(nil)
	client.Jar = jar

	if err := runScenario(ctx, client, baseURL.String(), db); err != nil {
		return fmt.Errorf("scenario failed: %w", err)
	}

	// Graceful shutdown not strictly needed for test runner exit, but nice
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	return nil
}

func createDatabase(ctx context.Context, dbName string) error {
	adminCfg := cfg.Database
	adminCfg.DBName = "postgres"
	adminDB := postgres.NewClient(adminCfg)
	if err := adminDB.Connect(ctx); err != nil {
		return err
	}
	defer adminDB.Close()

	// EXECUTE CREATE DATABASE
	// Note: Parametrized queries cannot be used for identifiers (DB name).
	// Must validate dbName first to prevent injection (though it's test code).
	query := fmt.Sprintf("CREATE DATABASE %s", dbName)
	if _, err := adminDB.Execute(ctx, query); err != nil {
		return err
	}

	// We also need to run migrations.
	// Since we don't have a migration tool integrated in code (external sql files),
	// we assume 'schema' is needed.
	// If migrations are in 001_initial.sql, we need to read and execute it against the NEW DB.

	return runMigrations(ctx, dbName)
}

func dropDatabase(ctx context.Context, dbName string) error {
	adminCfg := cfg.Database
	adminCfg.DBName = "postgres"
	adminDB := postgres.NewClient(adminCfg)
	if err := adminDB.Connect(ctx); err != nil {
		return err
	}
	defer adminDB.Close()

	// Force drop - using DROP DATABASE IF EXISTS WITH (FORCE) for modern PG
	query := fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE)", dbName)

	if _, err := adminDB.Execute(ctx, query); err != nil {
		return err
	}
	return nil
}

func runMigrations(ctx context.Context, dbName string) error {
	// Connect to the new DB
	migrationCfg := cfg.Database
	migrationCfg.DBName = dbName
	db := postgres.NewClient(migrationCfg)
	if err := db.Connect(ctx); err != nil {
		return err
	}
	defer db.Close()

	// Read migration file
	// Assuming migrations are at backend/migrations/001_init.sql or similar?
	// Need to find where the migrations are.
	// I recall user consolidating them into '001_init.sql'.
	// Path relative to execution? 'migrations/001_init.sql'.
	// I'll try to read it.

	content, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		// Try absolute path or search?
		// Fallback to checking typical locations or fail.
		return fmt.Errorf("reading migration file: %w", err)
	}

	if _, err := db.Execute(ctx, string(content)); err != nil {
		return fmt.Errorf("executing migration: %w", err)
	}
	return nil
}

func waitForServer(url *url.URL) error {
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url.String())
		if err == nil && resp.StatusCode == http.StatusOK {
			_ = resp.Body.Close()
			return nil
		}
		if resp != nil {
			_ = resp.Body.Close()
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for server %s", url)
}

func runScenario(ctx context.Context, client *http.Client, baseURL string, db database.Database) error {
	// 1. Sign Up
	email := fmt.Sprintf("test_%d@example.com", time.Now().Unix())
	password := "password123"

	log.Println("1. Signing up...")
	payload := map[string]string{
		"name":     "Integration User",
		"email":    email,
		"password": password,
	}
	resp, err := postJSON(client, baseURL+"/auth/signup", payload)
	if err != nil {
		return fmt.Errorf("signup request: %w", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("signup failed status: %d", resp.StatusCode)
	}

	// 1.5 Verify Email
	log.Println("1.5 Verifying email...")
	// Fetch token from DB
	rows, err := db.Query(ctx, "SELECT verification_token FROM user_email_credentials WHERE email = $1", email)
	if err != nil {
		return fmt.Errorf("querying verification token: %w", err)
	}
	defer rows.Close()

	var token string
	if rows.Next() {
		if err := rows.Scan(&token); err != nil {
			return fmt.Errorf("scanning verification token: %w", err)
		}
	} else {
		return fmt.Errorf("verification token not found for email: %s", email)
	}
	rows.Close()

	// Call Verify API
	verifyPayload := map[string]string{
		"token": token,
	}
	resp, err = postJSON(client, baseURL+"/auth/verify", verifyPayload)
	if err != nil {
		return fmt.Errorf("verify request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("verify failed status: %d", resp.StatusCode)
	}

	// 2. Login
	log.Println("2. Logging in...")
	loginPayload := map[string]string{
		"email":    email,
		"password": password,
	}
	resp, err = postJSON(client, baseURL+"/auth/login", loginPayload)
	if err != nil {
		return fmt.Errorf("login request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed status: %d", resp.StatusCode)
	}

	// Extract token
	var loginResp struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return fmt.Errorf("decoding login response: %w", err)
	}
	token = loginResp.Token
	if token == "" {
		return fmt.Errorf("login succeeded but token is empty")
	}

	// 3. Create Recipe
	log.Println("3. Creating recipe...")
	recipePayload := map[string]any{
		"name":    "Integration Recipe",
		"url":     "http://example.com/integration",
		"memo":    "Created by integration test",
		"tag_ids": []int{}, // Empty for now
	}
	resp, err = postJSONWithAuth(client, baseURL+"/recipes", recipePayload, token)
	if err != nil {
		return fmt.Errorf("create recipe request: %w", err)
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK { // Adjust expected status
		return fmt.Errorf("create recipe failed status: %d", resp.StatusCode)
	}

	// 4. Search Recipe (Wait for indexing)
	log.Println("4. Searching recipe...")
	// Elasticsearch indexing is async, so we might need a retry loop
	found := false
	searchPayload := map[string]string{
		"query": "Integration",
	}

	for range 20 { // Increased retries
		time.Sleep(500 * time.Millisecond) // Wait a bit

		// Search is POST /api/recipes/search
		resp, err := postJSONWithAuth(client, baseURL+"/recipes/search", searchPayload, token)
		if err != nil {
			log.Printf("Search request failed: %v", err)
			continue
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode == http.StatusOK {
			// Note: Response structure might differ.
			// Check handler response: respondJSON(w, http.StatusOK, response.NewRecipes(recipes))
			// NewRecipes returns []RecipeResponse.
			// If it's a list directly, `var recipes []map...` works.
			// If it's wrapped in object, we need that.
			// response.NewRecipes returns definition in response package.
			// Usually it is `type RecipesResponse struct { Recipes []RecipeResponse }` or just `[]RecipeResponse`.
			// Let's assume []RecipeResponse based on common Go patterns if not wrapped.
			// But looking at code `respondJSON(w, http.StatusOK, response.NewRecipes(recipes))`
			// I should check `response.NewRecipes`.
			// Failing that, decode interface{} and inspect.

			var body any
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				log.Printf("Failed to decode search response: %v", err)
				continue
			}

			// Assuming it returns a JSON array of recipes or object with "recipes" key
			// internal/server/handler/response/recipe.go likely defines it.
			// If I look at `GetRecipe` -> `response.NewRecipe(result)`.
			// `GetUserRecipes` -> `response.NewRecipes(recipes)`.
			// Standard is usually object `{"recipes": [...]}` or just `[...]`.

			// Let's try to handle both or inspect.
			// For simplicity in test, verify simply by looking for the name in the whole JSON string dump
			// if strict typing is hard without seeing response package.
			// But clean way is better.
			// I'll assume array []map[string]interface{} first.
			// Wait, if I Decode to interface{}, I can assert type.

			if recipesList, ok := body.([]any); ok {
				for _, r := range recipesList {
					if rMap, ok := r.(map[string]any); ok {
						if name, ok := rMap["name"].(string); ok && name == "Integration Recipe" {
							found = true
							break
						}
					}
				}
			} else if recipesObj, ok := body.(map[string]any); ok {
				// Maybe wrapped in {"recipes": [...]}
				if list, ok := recipesObj["recipes"].([]any); ok {
					for _, r := range list {
						if rMap, ok := r.(map[string]any); ok {
							if name, ok := rMap["name"].(string); ok && name == "Integration Recipe" {
								found = true
								break
							}
						}
					}
				}
			}
		}
		if found {
			break
		}
	}

	if !found {
		return fmt.Errorf("search failed: recipe not found after retries")
	}

	log.Println("Scenario completed!")
	return nil
}

func postJSON(client *http.Client, url string, data any) (*http.Response, error) {
	return postJSONWithAuth(client, url, data, "")
}

func postJSONWithAuth(client *http.Client, url string, data any, token string) (*http.Response, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return client.Do(req)
}

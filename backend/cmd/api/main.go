package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/domain/storage"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database"
	es "github.com/seka/reci-pin/backend/internal/infrastructure/database/elasticsearch"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database/postgres"
	"github.com/seka/reci-pin/backend/internal/infrastructure/storage/s3"
	"github.com/seka/reci-pin/backend/internal/registry"
	"github.com/seka/reci-pin/backend/internal/server"
)

var (
	cfg config.Config
)

func init() {
	flag.StringVar(&cfg.ApiServer.Port, "port", "8080", "Server port")
	flag.StringVar(&cfg.Database.Host, "db-host", "localhost", "Database host")
	flag.StringVar(&cfg.Database.Port, "db-port", "5432", "Database port")
	flag.StringVar(&cfg.Database.User, "db-user", "postgres", "Database user")
	flag.StringVar(&cfg.Database.Password, "db-password", "postgres", "Database password")
	flag.StringVar(&cfg.Database.DBName, "db-name", "recipin_dev", "Database name")
	flag.StringVar(&cfg.Database.SSLMode, "db-sslmode", "disable", "Database SSL mode")
	flag.StringVar(&cfg.ApiServer.JWT.Secret, "jwt-secret", "change-me", "JWT secret key")
	flag.IntVar(&cfg.ApiServer.JWT.ExpirationHours, "jwt-expiration", 24, "JWT expiration hours")

	// Storage configuration
	flag.StringVar(&cfg.Storage.Bucket, "storage-bucket", "recipin-bucket", "S3 bucket name")
	flag.StringVar(&cfg.Storage.Endpoint, "storage-endpoint", "", "S3 endpoint URL (for LocalStack)")
	flag.StringVar(&cfg.Storage.PublicBaseURL, "storage-public-url", "", "Base URL for public access")
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	// Start Database
	db := postgres.New(cfg.Database.DSN())
	if err := connectDatabase(ctx, db); err != nil {
		log.Fatalf("database error: %v", err)
	}
	defer db.Close()

	// Start Elasticsearch
	esClient, err := es.NewClient()
	if err != nil {
		log.Fatalf("elasticsearch error: %v", err)
	}

	// Start Storage Service
	storageService, err := s3.NewClient(ctx, cfg.Storage)
	if err != nil {
		log.Fatalf("storage service error: %v", err)
	}

	// Start Server
	srv := createServer(&cfg, db, esClient, storageService)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := runServer(ctx, srv); err != nil {
			errCh <- fmt.Errorf("server error: %w", err)
			cancel()
		}
	}()

	// Wait for signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		// Error occurred in one of the components
		log.Printf("Shutting down due to error: %v", err)
		cancel() // Ensure everyone gets the signal
	case sig := <-quit:
		log.Printf("Received signal: %v", sig)
		cancel()
	case <-ctx.Done():
		// Context cancelled elsewhere
		log.Println("Context cancelled")
	}

	wg.Wait()
	log.Println("Shutdown complete")
}

func connectDatabase(ctx context.Context, db database.Database) error {
	if err := db.Connect(ctx); err != nil {
		return err
	}
	log.Println("Successfully connected to database")
	return nil
}

func createServer(
	cfg *config.Config,
	db database.Database,
	esClient *elasticsearch.TypedClient,
	storageService storage.Client,
) *server.Server {
	repoRegistry := registry.NewRepository(db, esClient)
	useCaseRegistry := registry.NewUseCase(repoRegistry, storageService, cfg)
	return server.New(cfg, useCaseRegistry)
}

func runServer(ctx context.Context, srv *server.Server) error {
	serverErrCh := make(chan error, 1)
	go func() {
		serverErrCh <- srv.Run()
	}()

	select {
	case err := <-serverErrCh:
		return err
	case <-ctx.Done():
		log.Println("Shutting down server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return err
		}
		log.Println("Server exited")
		return nil
	}
}

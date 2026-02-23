package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/domain/notification"
	"github.com/seka/reci-pin/backend/internal/domain/storage"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database/postgres"
	"github.com/seka/reci-pin/backend/internal/infrastructure/notification/mailhog"
	"github.com/seka/reci-pin/backend/internal/infrastructure/searchengine"
	es "github.com/seka/reci-pin/backend/internal/infrastructure/searchengine/elasticsearch"
	"github.com/seka/reci-pin/backend/internal/infrastructure/storage/s3"
	"github.com/seka/reci-pin/backend/internal/registry"
	"github.com/seka/reci-pin/backend/internal/server"
)

var (
	cfg    config.Config
	esHost string
	esPort string
)

func init() {
	flag.StringVar(&cfg.ApiServer.Port, "port", "8080", "Server port")
	flag.StringVar(&cfg.ApiServer.JWT.Secret, "jwt-secret", "change-me", "JWT secret key")
	flag.IntVar(&cfg.ApiServer.JWT.ExpirationHours, "jwt-expiration", 24, "JWT expiration hours")
	flag.IntVar(&cfg.ApiServer.JWT.RefreshTokenExpirationDays, "jwt-refresh-expiration", 30, "JWT refresh token expiration days")

	// Database configuration
	flag.StringVar(&cfg.Database.Host, "db-host", "localhost", "Database host")
	flag.StringVar(&cfg.Database.Port, "db-port", "5432", "Database port")
	flag.StringVar(&cfg.Database.User, "db-user", "postgres", "Database user")
	flag.StringVar(&cfg.Database.Password, "db-password", "postgres", "Database password")
	flag.StringVar(&cfg.Database.DBName, "db-name", "recipin_dev", "Database name")
	flag.StringVar(&cfg.Database.SSLMode, "db-sslmode", "disable", "Database SSL mode")

	// Storage configuration
	flag.StringVar(&cfg.Storage.Bucket, "storage-bucket", "recipin-bucket", "S3 bucket name")
	flag.StringVar(&cfg.Storage.Endpoint, "storage-endpoint", "", "S3 endpoint URL (for LocalStack)")
	flag.StringVar(&cfg.Storage.PublicBaseURL, "storage-public-url", "", "Base URL for public access")

	// Search Engine configuration
	flag.StringVar(&esHost, "es-host", "localhost", "Elasticsearch host")
	flag.StringVar(&esPort, "es-port", "9200", "Elasticsearch port")
	cfg.SearchEngine.Addresses = []string{}
	flag.Func("es-addresses", "Elasticsearch addresses (comma separated)", func(s string) error {
		cfg.SearchEngine.Addresses = parseAddresses(s)
		return nil
	})

	// Mail configuration
	flag.StringVar(&cfg.Email.Host, "mail-host", "localhost", "MailHog host")
	flag.StringVar(&cfg.Email.Port, "mail-port", "1025", "MailHog port")
	flag.StringVar(&cfg.Email.From, "mail-from", "no-reply@reci-pin.com", "Mail from address")
}

func main() {
	flag.Parse()

	if len(cfg.SearchEngine.Addresses) == 0 {
		cfg.SearchEngine.Addresses = []string{"http://" + esHost + ":" + esPort}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	// Start Database
	db := postgres.NewClient(cfg.Database)
	if err := connectDatabase(ctx, db); err != nil {
		log.Fatalf("database error: %v", err)
	}
	defer db.Close()

	// Start Elasticsearch
	esClient, err := es.NewClient(cfg.SearchEngine)
	if err != nil {
		log.Fatalf("search engine error: %v", err)
	}

	// Start Storage Service
	storage, err := s3.NewClient(ctx, cfg.Storage)
	if err != nil {
		log.Fatalf("storage service error: %v", err)
	}

	// Start Mail Service
	mailClient := mailhog.NewClient(cfg.Email)

	// Start Server
	srv := createServer(&cfg, db, esClient, storage, mailClient)
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
	esClient searchengine.SearchEngine,
	storageService storage.Client,
	emailClient notification.EmailClient,
) *server.Server {
	repoRegistry := registry.NewRepository(db)
	searcherRegistry := registry.NewSearcher(esClient)
	useCaseRegistry := registry.NewUseCase(repoRegistry, storageService, searcherRegistry, emailClient, cfg)
	return server.New(&cfg.ApiServer, useCaseRegistry)
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

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/infrastructure/postgres"
	"github.com/seka/reci-pin/backend/internal/registry"
	"github.com/seka/reci-pin/backend/internal/server"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Application error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := parseConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := connectDB(ctx, cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	useCaseRegistry := createRegistry(db, cfg)
	srv := createServer(cfg, useCaseRegistry)

	return startServer(srv)
}

func parseConfig() *config.Config {
	cfg := &config.Config{}

	flag.IntVar(&cfg.Server.Port, "port", 8080, "Server port")
	flag.StringVar(&cfg.Database.Host, "db-host", "localhost", "Database host")
	flag.IntVar(&cfg.Database.Port, "db-port", 5432, "Database port")
	flag.StringVar(&cfg.Database.User, "db-user", "postgres", "Database user")
	flag.StringVar(&cfg.Database.Password, "db-password", "postgres", "Database password")
	flag.StringVar(&cfg.Database.DBName, "db-name", "recipin_dev", "Database name")
	flag.StringVar(&cfg.Database.SSLMode, "db-sslmode", "disable", "Database SSL mode")
	flag.StringVar(&cfg.JWT.Secret, "jwt-secret", "change-me", "JWT secret key")
	flag.IntVar(&cfg.JWT.ExpirationHours, "jwt-expiration", 24, "JWT expiration hours")

	flag.Parse()

	return cfg
}

func connectDB(ctx context.Context, cfg *config.Config) (postgres.Database, error) {
	db := postgres.New(cfg.Database.DSN())
	if err := db.Connect(ctx); err != nil {
		return nil, err
	}
	log.Println("Successfully connected to database")
	return db, nil
}

func createRegistry(db postgres.Database, cfg *config.Config) registry.UseCase {
	repoRegistry := registry.NewRepository(db)
	return registry.NewUseCase(repoRegistry, cfg)
}

func createServer(cfg *config.Config, useCaseRegistry registry.UseCase) *server.Server {
	return server.New(cfg, useCaseRegistry)
}

func startServer(srv *server.Server) error {
	// Start server in goroutine
	serverErrCh := make(chan error, 1)
	go func() {
		serverErrCh <- srv.Run()
	}()

	// Wait for signal or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrCh:
		return err
	case sig := <-quit:
		log.Printf("Received signal: %v", sig)
	}

	log.Println("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		// Just log error for shutdown, don't necessarily fail Main Run if expected
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Ensure DB is closed by defer above (in run function)

	log.Println("Server exited")
	return nil
}

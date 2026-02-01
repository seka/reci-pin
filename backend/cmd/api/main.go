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

var (
	args Arguments
)

func init() {
	flag.IntVar(&args.ServerPort, "port", 8080, "Server port")
	flag.StringVar(&args.DBHost, "db-host", "localhost", "Database host")
	flag.IntVar(&args.DBPort, "db-port", 5432, "Database port")
	flag.StringVar(&args.DBUser, "db-user", "postgres", "Database user")
	flag.StringVar(&args.DBPassword, "db-password", "postgres", "Database password")
	flag.StringVar(&args.DBName, "db-name", "recipin_dev", "Database name")
	flag.StringVar(&args.DBSSLMode, "db-sslmode", "disable", "Database SSL mode")
	flag.StringVar(&args.JWTSecret, "jwt-secret", "change-me", "JWT secret key")
	flag.IntVar(&args.JWTExpiration, "jwt-expiration", 24, "JWT expiration hours")
}

func main() {
	flag.Parse()
	if err := run(args); err != nil {
		log.Printf("Application error: %v", err)
		os.Exit(1)
	}
}

type Arguments struct {
	ServerPort    int
	DBHost        string
	DBPort        int
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	JWTSecret     string
	JWTExpiration int
}

func run(args Arguments) error {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     args.DBHost,
			Port:     args.DBPort,
			User:     args.DBUser,
			Password: args.DBPassword,
			DBName:   args.DBName,
			SSLMode:  args.DBSSLMode,
		},
		Server: config.ServerConfig{
			Port: args.ServerPort,
		},
		JWT: config.JWTConfig{
			Secret:          args.JWTSecret,
			ExpirationHours: args.JWTExpiration,
		},
	}

	db := postgres.New(cfg.Database.DSN())
	repoRegistry := registry.NewRepository(db)
	useCaseRegistry := registry.NewUseCase(repoRegistry, cfg)
	srv := server.New(cfg, useCaseRegistry)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to database
	if err := db.Connect(ctx); err != nil {
		return err
	}
	defer db.Close()
	log.Println("Successfully connected to database")

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

	// Ensure DB is closed by defer above

	log.Println("Server exited")
	return nil
}

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/infrastructure/postgres"
	"github.com/seka/reci-pin/backend/internal/registry"
	"github.com/seka/reci-pin/backend/internal/server"
)

func main() {
	// Parse command line flags
	var (
		serverPort    = flag.Int("port", 8080, "Server port")
		dbHost        = flag.String("db-host", "localhost", "Database host")
		dbPort        = flag.Int("db-port", 5432, "Database port")
		dbUser        = flag.String("db-user", "postgres", "Database user")
		dbPassword    = flag.String("db-password", "postgres", "Database password")
		dbName        = flag.String("db-name", "recipin_dev", "Database name")
		dbSSLMode     = flag.String("db-sslmode", "disable", "Database SSL mode")
		jwtSecret     = flag.String("jwt-secret", "change-me", "JWT secret key")
		jwtExpiration = flag.Int("jwt-expiration", 24, "JWT expiration hours")
	)
	flag.Parse()

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     *dbHost,
			Port:     *dbPort,
			User:     *dbUser,
			Password: *dbPassword,
			DBName:   *dbName,
			SSLMode:  *dbSSLMode,
		},
		Server: config.ServerConfig{
			Port: *serverPort,
		},
		JWT: config.JWTConfig{
			Secret:          *jwtSecret,
			ExpirationHours: *jwtExpiration,
		},
	}

	// Create database instance and connect
	db := postgres.New(cfg.Database.DSN())
	if err := db.Connect(context.Background()); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Successfully connected to database")

	// Initialize Repository Registry
	repoRegistry := registry.NewRepository(db)
	useCaseRegistry := registry.NewUseCase(repoRegistry, cfg)

	// Create server
	srv := server.New(cfg, useCaseRegistry)

	// Start server in goroutine
	go func() {
		if err := srv.Run(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	if err := srv.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

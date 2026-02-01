package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/infrastructure/datastore/postgres"
	"github.com/seka/reci-pin/backend/internal/registry"
	"github.com/seka/reci-pin/backend/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
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

	// Initialize UseCase Registry
	useCaseRegistry := registry.NewUseCase(repoRegistry, cfg)

	// Create server
	srv := server.New(cfg, repoRegistry, useCaseRegistry)

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

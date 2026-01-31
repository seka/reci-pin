package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/infrastructure/datastore/postgres"
	"github.com/seka/reci-pin/backend/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	ctx := context.Background()
	db, err := postgres.New(ctx, cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to database")

	// Initialize server
	srv := server.NewServer(db)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, srv); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

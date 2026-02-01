package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := server.Run(ctx, cfg); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

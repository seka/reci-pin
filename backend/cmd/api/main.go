package main

import (
	"context"
	"log"

	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()

	if err := server.Run(ctx, cfg); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

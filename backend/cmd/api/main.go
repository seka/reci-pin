package main

import (
	"context"
	"fmt"
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

	srv, err := server.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	if err := srv.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

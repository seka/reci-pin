package main

import (
	"log"

	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := srv.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

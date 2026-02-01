package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/infrastructure/postgres"
	"github.com/seka/reci-pin/backend/internal/registry"
	"github.com/seka/reci-pin/backend/internal/server"
)

var (
	cfg config.Config
)

func init() {
	flag.IntVar(&cfg.Server.Port, "port", 8080, "Server port")
	flag.StringVar(&cfg.Database.Host, "db-host", "localhost", "Database host")
	flag.IntVar(&cfg.Database.Port, "db-port", 5432, "Database port")
	flag.StringVar(&cfg.Database.User, "db-user", "postgres", "Database user")
	flag.StringVar(&cfg.Database.Password, "db-password", "postgres", "Database password")
	flag.StringVar(&cfg.Database.DBName, "db-name", "recipin_dev", "Database name")
	flag.StringVar(&cfg.Database.SSLMode, "db-sslmode", "disable", "Database SSL mode")
	flag.StringVar(&cfg.JWT.Secret, "jwt-secret", "change-me", "JWT secret key")
	flag.IntVar(&cfg.JWT.ExpirationHours, "jwt-expiration", 24, "JWT expiration hours")
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	// Start Database
	db := postgres.New(cfg.Database.DSN())
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := runDatabase(ctx, db); err != nil {
			errCh <- fmt.Errorf("database error: %w", err)
			cancel() // Cancel context to stop other components
		}
	}()

	// Start Server
	srv := createServer(&cfg, db)
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

func runDatabase(ctx context.Context, db postgres.Database) error {
	if err := db.Connect(ctx); err != nil {
		return err
	}
	defer db.Close()
	log.Println("Successfully connected to database")

	<-ctx.Done()
	return nil
}

func createServer(cfg *config.Config, db postgres.Database) *server.Server {
	repoRegistry := registry.NewRepository(db)
	useCaseRegistry := registry.NewUseCase(repoRegistry, cfg)
	return server.New(cfg, useCaseRegistry)
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

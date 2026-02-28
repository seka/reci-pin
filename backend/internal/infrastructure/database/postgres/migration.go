package postgres

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const defaultMigrationsDir = "migrations"

// RunMigrations executes all up migrations in the migrations directory.
func RunMigrations(dsn string) error {
	// Find migrations directory
	// In development, it's usually in the project root/backend/migrations
	// In Docker, it's in /app/migrations
	migrationsPath := findMigrationsPath()
	if migrationsPath == "" {
		return fmt.Errorf("migrations directory not found")
	}

	sourceURL := fmt.Sprintf("file://%s", migrationsPath)
	
	m, err := migrate.New(sourceURL, dsn)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Migrations executed successfully")
	return nil
}

func findMigrationsPath() string {
	// 1. Check current directory
	// Use filepath.Join for cross-platform compatibility
	if _, err := os.Stat(defaultMigrationsDir); err == nil {
		path, _ := filepath.Abs(defaultMigrationsDir)
		return path
	}

	// 2. Check parent directory (if run from cmd/api)
	parentMigrationsDir := filepath.Join("..", "..", defaultMigrationsDir)
	if _, err := os.Stat(parentMigrationsDir); err == nil {
		path, _ := filepath.Abs(parentMigrationsDir)
		return path
	}

	// 3. Environment variable or fallback
	return ""
}

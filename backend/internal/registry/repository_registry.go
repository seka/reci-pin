package registry

import (
	"context"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/infrastructure/datastore/postgres"
)

// Repository defines the interface for creating repositories
type Repository interface {
	NewUserRepository() repository.UserRepository
	NewRecipeRepository() repository.RecipeRepository
	NewTagRepository() repository.TagRepository
	NewRecipeImageRepository() repository.RecipeImageRepository
	NewUserEmailCredentialRepository() repository.UserEmailCredentialRepository
	Close() error
}

// repositoryRegistry implements the Repository interface
type repositoryRegistry struct {
	db *postgres.DB
}

// NewRepository creates a new Repository registry
func NewRepository(ctx context.Context, dsn string) (Repository, error) {
	db, err := postgres.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &repositoryRegistry{db: db}, nil
}

func (r *repositoryRegistry) NewUserRepository() repository.UserRepository {
	return postgres.NewUserRepository(r.db)
}

func (r *repositoryRegistry) NewRecipeRepository() repository.RecipeRepository {
	return postgres.NewRecipeRepository(r.db)
}

func (r *repositoryRegistry) NewTagRepository() repository.TagRepository {
	return postgres.NewTagRepository(r.db)
}

func (r *repositoryRegistry) NewRecipeImageRepository() repository.RecipeImageRepository {
	return postgres.NewRecipeImageRepository(r.db)
}

func (r *repositoryRegistry) NewUserEmailCredentialRepository() repository.UserEmailCredentialRepository {
	return postgres.NewUserEmailCredentialRepository(r.db)
}

func (r *repositoryRegistry) Close() error {
	if r.db != nil {
		r.db.Close()
	}
	return nil
}

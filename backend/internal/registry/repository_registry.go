package registry

import (
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/infrastructure/postgres"
)

// Repository defines the interface for creating repositories
type Repository interface {
	NewUserRepository() repository.UserRepository
	NewRecipeRepository() repository.RecipeRepository
	NewTagRepository() repository.TagRepository
	NewRecipeImageRepository() repository.RecipeImageRepository
	NewUserEmailCredentialRepository() repository.UserEmailCredentialRepository
}

// repositoryRegistry implements the Repository interface
type repositoryRegistry struct {
	db postgres.Database
}

// NewRepository creates a new Repository registry
func NewRepository(db postgres.Database) Repository {
	return &repositoryRegistry{db: db}
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

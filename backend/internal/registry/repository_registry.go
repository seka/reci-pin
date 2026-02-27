package registry

import (
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database"
	postgres "github.com/seka/reci-pin/backend/internal/infrastructure/database/postgres"
)

// Repository defines the interface for creating repositories
type Repository interface {
	NewUserRepository() repository.UserRepository
	NewRecipeRepository() repository.RecipeRepository
	NewTagRepository() repository.TagRepository
	NewRecipeImageRepository() repository.RecipeImageRepository
	NewUserEmailCredentialRepository() repository.UserEmailCredentialRepository
	NewPasswordResetTokenRepository() repository.PasswordResetTokenRepository
	NewRefreshTokenRepository() repository.RefreshTokenRepository
	TransactionManager() repository.TransactionManager
}

// repositoryRegistry implements the Repository interface
type repositoryRegistry struct {
	db database.Database
}

// NewRepository creates a new Repository registry
func NewRepository(db database.Database) Repository {
	return &repositoryRegistry{
		db: db,
	}
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

func (r *repositoryRegistry) NewPasswordResetTokenRepository() repository.PasswordResetTokenRepository {
	return postgres.NewPasswordResetTokenRepository(r.db)
}

func (r *repositoryRegistry) NewRefreshTokenRepository() repository.RefreshTokenRepository {
	return postgres.NewRefreshTokenRepository(r.db)
}

func (r *repositoryRegistry) TransactionManager() repository.TransactionManager {
	return r.db.TransactionManager()
}

package registry

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database"
	es "github.com/seka/reci-pin/backend/internal/infrastructure/database/elasticsearch"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database/postgres"
)

// Repository defines the interface for creating repositories
type Repository interface {
	NewUserRepository() repository.UserRepository
	NewRecipeRepository() repository.RecipeRepository
	NewTagRepository() repository.TagRepository
	NewRecipeImageRepository() repository.RecipeImageRepository
	NewUserEmailCredentialRepository() repository.UserEmailCredentialRepository
	NewPasswordResetTokenRepository() repository.PasswordResetTokenRepository
	NewRecipeSearchRepository() repository.RecipeSearchRepository
	NewRefreshTokenRepository() repository.RefreshTokenRepository
}

// repositoryRegistry implements the Repository interface
type repositoryRegistry struct {
	db           database.Database
	searchEngine *elasticsearch.TypedClient
}

// NewRepository creates a new Repository registry
func NewRepository(db database.Database, searchEngine *elasticsearch.TypedClient) Repository {
	return &repositoryRegistry{
		db:           db,
		searchEngine: searchEngine,
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

func (r *repositoryRegistry) NewRecipeSearchRepository() repository.RecipeSearchRepository {
	return es.NewRecipeSearcher(r.searchEngine)
}

func (r *repositoryRegistry) NewRefreshTokenRepository() repository.RefreshTokenRepository {
	return postgres.NewRefreshTokenRepository(r.db)
}

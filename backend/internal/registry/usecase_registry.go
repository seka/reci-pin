package registry

//go:generate mockgen -source=$GOFILE -destination=mock/usecase_registry_mock.go -package=mock

import (
	"time"

	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/domain/notification"
	"github.com/seka/reci-pin/backend/internal/domain/storage"
	"github.com/seka/reci-pin/backend/internal/usecase/auth"
	"github.com/seka/reci-pin/backend/internal/usecase/recipe"
	"github.com/seka/reci-pin/backend/internal/usecase/recipe_image"
	"github.com/seka/reci-pin/backend/internal/usecase/recipe_tag"
	"github.com/seka/reci-pin/backend/internal/usecase/tag"
)

// UseCase defines the interface for creating use cases
type UseCase interface {
	// Auth
	NewSignupUseCase() auth.SignupUseCase
	NewLoginUseCase() auth.LoginUseCase
	NewGenerateTokenUseCase() auth.GenerateTokenUseCase
	NewValidateTokenUseCase() auth.ValidateTokenUseCase
	NewRefreshTokenUseCase() auth.RefreshTokenUseCase
	NewGetUserUseCase() auth.GetUserUseCase
	NewVerifyEmailUseCase() auth.VerifyEmailUseCase
	NewWithdrawUseCase() auth.WithdrawUseCase
	NewChangePasswordUseCase() auth.ChangePasswordUseCase
	NewRequestPasswordResetUseCase() auth.RequestPasswordResetUseCase
	NewResetPasswordUseCase() auth.ResetPasswordUseCase
	NewLogoutUseCase() auth.LogoutUseCase

	// Recipe
	NewCreateRecipeUseCase() recipe.CreateRecipeUseCase
	NewGetRecipeUseCase() recipe.GetRecipeUseCase
	NewGetUserRecipesUseCase() recipe.GetUserRecipesUseCase
	NewUpdateRecipeUseCase() recipe.UpdateRecipeUseCase
	NewDeleteRecipeUseCase() recipe.DeleteRecipeUseCase
	NewSearchRecipesUseCase() recipe.SearchRecipesUseCase

	// Recipe Tag
	NewAddTagsUseCase() recipe_tag.AddTagsUseCase
	NewRemoveTagsUseCase() recipe_tag.RemoveTagsUseCase

	// Recipe Image
	NewCreateRecipeImageUseCase() recipe_image.CreateRecipeImageUseCase

	// Tag
	NewCreateTagUseCase() tag.CreateTagUseCase
	NewGetAllTagsUseCase() tag.GetAllTagsUseCase
	NewDeleteTagUseCase() tag.DeleteTagUseCase
}

// useCaseRegistry implements the UseCase interface
type useCaseRegistry struct {
	repo        Repository
	storage     storage.Client
	searcher    Searcher
	emailClient notification.EmailClient
	cfg         *config.Config
}

// NewUseCase creates a new UseCase registry
func NewUseCase(repo Repository, storage storage.Client, searcher Searcher, email notification.EmailClient, cfg *config.Config) UseCase {
	return &useCaseRegistry{
		repo:        repo,
		storage:     storage,
		searcher:    searcher,
		emailClient: email,
		cfg:         cfg,
	}
}

// Auth UseCases
func (u *useCaseRegistry) NewSignupUseCase() auth.SignupUseCase {
	return auth.NewSignupUseCase(u.repo.NewUserRepository(), u.repo.NewUserEmailCredentialRepository(), u.repo.TransactionManager())
}

func (u *useCaseRegistry) NewLoginUseCase() auth.LoginUseCase {
	return auth.NewLoginUseCase(u.repo.NewUserEmailCredentialRepository())
}

func (u *useCaseRegistry) NewGenerateTokenUseCase() auth.GenerateTokenUseCase {
	return auth.NewGenerateTokenUseCase(
		u.cfg.ApiServer.JWT.Secret,
		time.Duration(u.cfg.ApiServer.JWT.ExpirationHours)*time.Hour,
		u.repo.NewRefreshTokenRepository(),
		time.Duration(u.cfg.ApiServer.JWT.RefreshTokenExpirationDays)*24*time.Hour,
	)
}

func (u *useCaseRegistry) NewRefreshTokenUseCase() auth.RefreshTokenUseCase {
	return auth.NewRefreshTokenUseCase(u.NewGenerateTokenUseCase(), u.repo.NewRefreshTokenRepository())
}

func (u *useCaseRegistry) NewLogoutUseCase() auth.LogoutUseCase {
	return auth.NewLogoutUseCase(u.repo.NewRefreshTokenRepository())
}

func (u *useCaseRegistry) NewValidateTokenUseCase() auth.ValidateTokenUseCase {
	return auth.NewValidateTokenUseCase(u.cfg.ApiServer.JWT.Secret)
}

func (u *useCaseRegistry) NewGetUserUseCase() auth.GetUserUseCase {
	return auth.NewGetUserUseCase(u.repo.NewUserRepository())
}

func (u *useCaseRegistry) NewVerifyEmailUseCase() auth.VerifyEmailUseCase {
	return auth.NewVerifyEmailUseCase(u.repo.NewUserEmailCredentialRepository(), u.repo.TransactionManager())
}

func (u *useCaseRegistry) NewWithdrawUseCase() auth.WithdrawUseCase {
	return auth.NewWithdrawUseCase(u.repo.NewUserRepository())
}

func (u *useCaseRegistry) NewChangePasswordUseCase() auth.ChangePasswordUseCase {
	return auth.NewChangePasswordUseCase(u.repo.NewUserEmailCredentialRepository(), u.emailClient, u.repo.TransactionManager())
}

func (u *useCaseRegistry) NewRequestPasswordResetUseCase() auth.RequestPasswordResetUseCase {
	return auth.NewRequestPasswordResetUseCase(u.repo.NewUserEmailCredentialRepository(), u.repo.NewPasswordResetTokenRepository(), u.emailClient)
}

func (u *useCaseRegistry) NewResetPasswordUseCase() auth.ResetPasswordUseCase {
	return auth.NewResetPasswordUseCase(u.repo.NewPasswordResetTokenRepository(), u.repo.NewUserEmailCredentialRepository(), u.repo.TransactionManager())
}

// Recipe UseCases
func (u *useCaseRegistry) NewCreateRecipeUseCase() recipe.CreateRecipeUseCase {
	return recipe.NewCreateRecipeUseCase(u.repo.NewRecipeRepository(), u.searcher.NewRecipeSearchRepository(), u.repo.TransactionManager())
}

func (u *useCaseRegistry) NewGetRecipeUseCase() recipe.GetRecipeUseCase {
	return recipe.NewGetRecipeUseCase(u.repo.NewRecipeRepository(), u.repo.NewRecipeImageRepository(), u.storage)
}

func (u *useCaseRegistry) NewGetUserRecipesUseCase() recipe.GetUserRecipesUseCase {
	return recipe.NewGetUserRecipesUseCase(u.repo.NewRecipeRepository(), u.repo.NewRecipeImageRepository(), u.storage)
}

func (u *useCaseRegistry) NewUpdateRecipeUseCase() recipe.UpdateRecipeUseCase {
	return recipe.NewUpdateRecipeUseCase(u.repo.NewRecipeRepository(), u.repo.NewRecipeImageRepository(), u.searcher.NewRecipeSearchRepository(), u.storage)
}

func (u *useCaseRegistry) NewDeleteRecipeUseCase() recipe.DeleteRecipeUseCase {
	return recipe.NewDeleteRecipeUseCase(u.repo.NewRecipeRepository(), u.searcher.NewRecipeSearchRepository())
}

func (u *useCaseRegistry) NewSearchRecipesUseCase() recipe.SearchRecipesUseCase {
	return recipe.NewSearchRecipesUseCase(
		u.repo.NewRecipeRepository(),
		u.repo.NewRecipeImageRepository(),
		u.searcher.NewRecipeSearchRepository(),
		u.storage,
	)
}

// Recipe Tag UseCases
func (u *useCaseRegistry) NewAddTagsUseCase() recipe_tag.AddTagsUseCase {
	return recipe_tag.NewAddTagsUseCase(u.repo.NewRecipeRepository(), u.repo.TransactionManager())
}

func (u *useCaseRegistry) NewRemoveTagsUseCase() recipe_tag.RemoveTagsUseCase {
	return recipe_tag.NewRemoveTagsUseCase(u.repo.NewRecipeRepository(), u.repo.TransactionManager())
}

// Recipe Image UseCase
func (u *useCaseRegistry) NewCreateRecipeImageUseCase() recipe_image.CreateRecipeImageUseCase {
	return recipe_image.NewCreateRecipeImageUseCase(
		u.repo.NewRecipeRepository(),
		u.repo.NewRecipeImageRepository(),
		u.storage,
	)
}

// Tag UseCases
func (u *useCaseRegistry) NewCreateTagUseCase() tag.CreateTagUseCase {
	return tag.NewCreateTagUseCase(u.repo.NewTagRepository())
}

func (u *useCaseRegistry) NewGetAllTagsUseCase() tag.GetAllTagsUseCase {
	return tag.NewGetAllTagsUseCase(u.repo.NewTagRepository())
}

func (u *useCaseRegistry) NewDeleteTagUseCase() tag.DeleteTagUseCase {
	return tag.NewDeleteTagUseCase(u.repo.NewTagRepository())
}

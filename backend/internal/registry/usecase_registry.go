package registry

import (
	"time"

	"github.com/seka/reci-pin/backend/config"
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
	NewGetUserUseCase() auth.GetUserUseCase
	NewVerifyEmailUseCase() auth.VerifyEmailUseCase
	NewWithdrawUseCase() auth.WithdrawUseCase

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
	NewAddImageUseCase() recipe_image.AddImageUseCase

	// Tag
	NewCreateTagUseCase() tag.CreateTagUseCase
	NewGetAllTagsUseCase() tag.GetAllTagsUseCase
	NewDeleteTagUseCase() tag.DeleteTagUseCase
}

// useCaseRegistry implements the UseCase interface
type useCaseRegistry struct {
	repo Repository
	cfg  *config.Config
}

// NewUseCase creates a new UseCase registry
func NewUseCase(repo Repository, cfg *config.Config) UseCase {
	return &useCaseRegistry{
		repo: repo,
		cfg:  cfg,
	}
}

// Auth UseCases
func (u *useCaseRegistry) NewSignupUseCase() auth.SignupUseCase {
	return auth.NewSignupUseCase(u.repo.NewUserRepository(), u.repo.NewUserEmailCredentialRepository())
}

func (u *useCaseRegistry) NewLoginUseCase() auth.LoginUseCase {
	return auth.NewLoginUseCase(u.repo.NewUserEmailCredentialRepository())
}

func (u *useCaseRegistry) NewGenerateTokenUseCase() auth.GenerateTokenUseCase {
	return auth.NewGenerateTokenUseCase(
		u.cfg.JWT.Secret,
		time.Duration(u.cfg.JWT.ExpirationHours)*time.Hour,
	)
}

func (u *useCaseRegistry) NewValidateTokenUseCase() auth.ValidateTokenUseCase {
	return auth.NewValidateTokenUseCase(u.cfg.JWT.Secret)
}

func (u *useCaseRegistry) NewGetUserUseCase() auth.GetUserUseCase {
	return auth.NewGetUserUseCase(u.repo.NewUserRepository())
}

func (u *useCaseRegistry) NewVerifyEmailUseCase() auth.VerifyEmailUseCase {
	return auth.NewVerifyEmailUseCase(u.repo.NewUserEmailCredentialRepository())
}

func (u *useCaseRegistry) NewWithdrawUseCase() auth.WithdrawUseCase {
	return auth.NewWithdrawUseCase(u.repo.NewUserRepository())
}

// Recipe UseCases
func (u *useCaseRegistry) NewCreateRecipeUseCase() recipe.CreateRecipeUseCase {
	return recipe.NewCreateRecipeUseCase(u.repo.NewRecipeRepository())
}

func (u *useCaseRegistry) NewGetRecipeUseCase() recipe.GetRecipeUseCase {
	return recipe.NewGetRecipeUseCase(u.repo.NewRecipeRepository(), u.repo.NewRecipeImageRepository())
}

func (u *useCaseRegistry) NewGetUserRecipesUseCase() recipe.GetUserRecipesUseCase {
	return recipe.NewGetUserRecipesUseCase(u.repo.NewRecipeRepository(), u.repo.NewRecipeImageRepository())
}

func (u *useCaseRegistry) NewUpdateRecipeUseCase() recipe.UpdateRecipeUseCase {
	return recipe.NewUpdateRecipeUseCase(u.repo.NewRecipeRepository())
}

func (u *useCaseRegistry) NewDeleteRecipeUseCase() recipe.DeleteRecipeUseCase {
	return recipe.NewDeleteRecipeUseCase(u.repo.NewRecipeRepository())
}

func (u *useCaseRegistry) NewSearchRecipesUseCase() recipe.SearchRecipesUseCase {
	return recipe.NewSearchRecipesUseCase(u.repo.NewRecipeRepository(), u.repo.NewRecipeImageRepository())
}

// Recipe Tag UseCases
func (u *useCaseRegistry) NewAddTagsUseCase() recipe_tag.AddTagsUseCase {
	return recipe_tag.NewAddTagsUseCase(u.repo.NewRecipeRepository())
}

func (u *useCaseRegistry) NewRemoveTagsUseCase() recipe_tag.RemoveTagsUseCase {
	return recipe_tag.NewRemoveTagsUseCase(u.repo.NewRecipeRepository())
}

// Recipe Image UseCase
func (u *useCaseRegistry) NewAddImageUseCase() recipe_image.AddImageUseCase {
	return recipe_image.NewAddImageUseCase(u.repo.NewRecipeRepository(), u.repo.NewRecipeImageRepository())
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

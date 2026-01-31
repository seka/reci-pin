package postgres

import (
	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/infrastructure/entity"
)

// User conversions

func userEntityToModel(e *entity.User) *model.User {
	if e == nil {
		return nil
	}
	return &model.User{
		ID:   e.ID,
		Name: e.Name,
	}
}

func userModelToEntity(m *model.User) *entity.User {
	if m == nil {
		return nil
	}
	return &entity.User{
		ID:   m.ID,
		Name: m.Name,
	}
}

// UserEmailCredential conversions

func userEmailCredentialEntityToModel(e *entity.UserEmailCredential) *model.UserEmailCredential {
	if e == nil {
		return nil
	}
	return &model.UserEmailCredential{
		UserID:                     e.UserID,
		Email:                      e.Email,
		PasswordHash:               e.PasswordHash,
		EmailVerifiedAt:            e.EmailVerifiedAt,
		VerificationToken:          e.VerificationToken,
		VerificationTokenExpiresAt: e.VerificationTokenExpiresAt,
	}
}

func userEmailCredentialModelToEntity(m *model.UserEmailCredential) *entity.UserEmailCredential {
	if m == nil {
		return nil
	}
	return &entity.UserEmailCredential{
		UserID:                     m.UserID,
		Email:                      m.Email,
		PasswordHash:               m.PasswordHash,
		EmailVerifiedAt:            m.EmailVerifiedAt,
		VerificationToken:          m.VerificationToken,
		VerificationTokenExpiresAt: m.VerificationTokenExpiresAt,
	}
}

// Recipe conversions

func recipeEntityToModel(e *entity.Recipe) *model.Recipe {
	if e == nil {
		return nil
	}
	return &model.Recipe{
		ID:     e.ID,
		UserID: e.UserID,
		Name:   e.Name,
		URL:    e.URL,
		Memo:   e.Memo,
	}
}

func recipeModelToEntity(m *model.Recipe) *entity.Recipe {
	if m == nil {
		return nil
	}
	return &entity.Recipe{
		ID:     m.ID,
		UserID: m.UserID,
		Name:   m.Name,
		URL:    m.URL,
		Memo:   m.Memo,
	}
}

func recipeEntitiesToModels(entities []entity.Recipe) []model.Recipe {
	models := make([]model.Recipe, len(entities))
	for i, e := range entities {
		models[i] = *recipeEntityToModel(&e)
	}
	return models
}

// Tag conversions

func tagEntityToModel(e *entity.Tag) *model.Tag {
	if e == nil {
		return nil
	}
	return &model.Tag{
		ID:   e.ID,
		Name: e.Name,
	}
}

func tagModelToEntity(m *model.Tag) *entity.Tag {
	if m == nil {
		return nil
	}
	return &entity.Tag{
		ID:   m.ID,
		Name: m.Name,
	}
}

func tagEntitiesToModels(entities []entity.Tag) []model.Tag {
	models := make([]model.Tag, len(entities))
	for i, e := range entities {
		models[i] = *tagEntityToModel(&e)
	}
	return models
}

// RecipeImage conversions

func recipeImageEntityToModel(e *entity.RecipeImage) *model.RecipeImage {
	if e == nil {
		return nil
	}
	return &model.RecipeImage{
		ID:        e.ID,
		RecipeID:  e.RecipeID,
		ImagePath: e.ImagePath,
	}
}

func recipeImageModelToEntity(m *model.RecipeImage) *entity.RecipeImage {
	if m == nil {
		return nil
	}
	return &entity.RecipeImage{
		ID:        m.ID,
		RecipeID:  m.RecipeID,
		ImagePath: m.ImagePath,
	}
}

func recipeImageEntitiesToModels(entities []entity.RecipeImage) []model.RecipeImage {
	models := make([]model.RecipeImage, len(entities))
	for i, e := range entities {
		models[i] = *recipeImageEntityToModel(&e)
	}
	return models
}

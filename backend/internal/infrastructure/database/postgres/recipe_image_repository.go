package postgres

import (
	"context"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database"
	"github.com/seka/reci-pin/backend/internal/infrastructure/entity"
)

type RecipeImageRepository struct {
	db database.Database
}

func NewRecipeImageRepository(db database.Database) repository.RecipeImageRepository {
	return &RecipeImageRepository{db: db}
}

func (r *RecipeImageRepository) Create(ctx context.Context, image *model.RecipeImage) error {
	e := entity.NewRecipeImage(image)
	query := `
		INSERT INTO recipe_images (recipe_id, image_path, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, created_at
	`
	rows, err := r.db.Query(ctx, query, e.RecipeID, e.ImagePath)
	if err != nil {
		return fmt.Errorf("failed to create recipe image: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return fmt.Errorf("failed to create recipe image: no rows returned")
	}

	err = rows.Scan(&e.ID, &e.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to scan recipe image: %w", err)
	}

	image.ID = e.ID
	return nil
}

func (r *RecipeImageRepository) GetByRecipeID(ctx context.Context, recipeID int64) ([]model.RecipeImage, error) {
	query := `
		SELECT id, recipe_id, image_path, created_at
		FROM recipe_images
		WHERE recipe_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe images: %w", err)
	}
	defer rows.Close()

	var entities []entity.RecipeImage
	for rows.Next() {
		var e entity.RecipeImage
		if err := rows.Scan(&e.ID, &e.RecipeID, &e.ImagePath, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan recipe image: %w", err)
		}
		entities = append(entities, e)
	}

	return r.toModels(entities), nil
}

func (r *RecipeImageRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM recipe_images WHERE id = $1`
	_, err := r.db.Execute(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete recipe image: %w", err)
	}
	return nil
}

func (r *RecipeImageRepository) toModels(entities []entity.RecipeImage) []model.RecipeImage {
	models := make([]model.RecipeImage, len(entities))
	for i, e := range entities {
		models[i] = *e.ToModel()
	}
	return models
}

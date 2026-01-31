package postgres

import (
	"context"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/entity"
	"github.com/seka/reci-pin/backend/internal/domain/model"
)

type RecipeImageRepository struct {
	db *DB
}

func NewRecipeImageRepository(db *DB) *RecipeImageRepository {
	return &RecipeImageRepository{db: db}
}

func (r *RecipeImageRepository) Create(ctx context.Context, image *model.RecipeImage) error {
	query := `
		INSERT INTO recipe_images (recipe_id, image_path, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id
	`
	err := r.db.Pool.QueryRow(ctx, query, image.RecipeID, image.ImagePath).Scan(&image.ID)
	if err != nil {
		return fmt.Errorf("failed to create recipe image: %w", err)
	}
	return nil
}

func (r *RecipeImageRepository) GetByRecipeID(ctx context.Context, recipeID int64) ([]model.RecipeImage, error) {
	query := `
		SELECT id, recipe_id, image_path
		FROM recipe_images
		WHERE recipe_id = $1
		ORDER BY created_at
	`
	rows, err := r.db.Pool.Query(ctx, query, recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe images: %w", err)
	}
	defer rows.Close()

	var imageEntities []entity.RecipeImage
	for rows.Next() {
		var img entity.RecipeImage
		if err := rows.Scan(&img.ID, &img.RecipeID, &img.ImagePath); err != nil {
			return nil, fmt.Errorf("failed to scan recipe image: %w", err)
		}
		imageEntities = append(imageEntities, img)
	}
	return recipeImageEntitiesToModels(imageEntities), nil
}

func (r *RecipeImageRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM recipe_images WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete recipe image: %w", err)
	}
	return nil
}

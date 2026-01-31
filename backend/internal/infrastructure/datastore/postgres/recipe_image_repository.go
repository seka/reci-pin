package postgres

import (
	"context"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/entity"
)

type RecipeImageRepository struct {
	db *DB
}

func NewRecipeImageRepository(db *DB) *RecipeImageRepository {
	return &RecipeImageRepository{db: db}
}

func (r *RecipeImageRepository) Create(ctx context.Context, image *entity.RecipeImage) error {
	query := `
		INSERT INTO recipe_images (recipe_id, image_path, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, created_at
	`
	err := r.db.Pool.QueryRow(ctx, query, image.RecipeID, image.ImagePath).
		Scan(&image.ID, &image.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create recipe image: %w", err)
	}
	return nil
}

func (r *RecipeImageRepository) GetByRecipeID(ctx context.Context, recipeID int64) ([]entity.RecipeImage, error) {
	query := `
		SELECT id, recipe_id, image_path, created_at
		FROM recipe_images
		WHERE recipe_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe images: %w", err)
	}
	defer rows.Close()

	var images []entity.RecipeImage
	for rows.Next() {
		var image entity.RecipeImage
		if err := rows.Scan(&image.ID, &image.RecipeID, &image.ImagePath, &image.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan recipe image: %w", err)
		}
		images = append(images, image)
	}
	return images, nil
}

func (r *RecipeImageRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM recipe_images WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete recipe image: %w", err)
	}
	return nil
}

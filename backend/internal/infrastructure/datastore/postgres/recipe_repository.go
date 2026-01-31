package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/seka/reci-pin/backend/internal/domain/entity"
)

type RecipeRepository struct {
	db *DB
}

func NewRecipeRepository(db *DB) *RecipeRepository {
	return &RecipeRepository{db: db}
}

func (r *RecipeRepository) Create(ctx context.Context, recipe *entity.Recipe) error {
	query := `
		INSERT INTO recipes (user_id, name, url, memo, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.Pool.QueryRow(ctx, query, recipe.UserID, recipe.Name, recipe.URL, recipe.Memo).
		Scan(&recipe.ID, &recipe.CreatedAt, &recipe.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create recipe: %w", err)
	}
	return nil
}

func (r *RecipeRepository) GetByID(ctx context.Context, id int64) (*entity.Recipe, error) {
	query := `
		SELECT id, user_id, name, url, memo, created_at, updated_at
		FROM recipes
		WHERE id = $1
	`
	recipe := &entity.Recipe{}
	err := r.db.Pool.QueryRow(ctx, query, id).
		Scan(&recipe.ID, &recipe.UserID, &recipe.Name, &recipe.URL, &recipe.Memo, &recipe.CreatedAt, &recipe.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe by id: %w", err)
	}
	return recipe, nil
}

func (r *RecipeRepository) GetByUserID(ctx context.Context, userID int64) ([]entity.Recipe, error) {
	query := `
		SELECT id, user_id, name, url, memo, created_at, updated_at
		FROM recipes
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipes by user id: %w", err)
	}
	defer rows.Close()

	var recipes []entity.Recipe
	for rows.Next() {
		var recipe entity.Recipe
		if err := rows.Scan(&recipe.ID, &recipe.UserID, &recipe.Name, &recipe.URL, &recipe.Memo, &recipe.CreatedAt, &recipe.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan recipe: %w", err)
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (r *RecipeRepository) Search(ctx context.Context, userID int64, query string, tagIDs []int64) ([]entity.Recipe, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	conditions = append(conditions, fmt.Sprintf("r.user_id = $%d", argIndex))
	args = append(args, userID)
	argIndex++

	if query != "" {
		conditions = append(conditions, fmt.Sprintf("(r.name ILIKE $%d OR r.memo ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+query+"%")
		argIndex++
	}

	sqlQuery := `
		SELECT DISTINCT r.id, r.user_id, r.name, r.url, r.memo, r.created_at, r.updated_at
		FROM recipes r
	`

	if len(tagIDs) > 0 {
		sqlQuery += ` INNER JOIN recipe_tags rt ON r.id = rt.recipe_id`
		placeholders := make([]string, len(tagIDs))
		for i, tagID := range tagIDs {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, tagID)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("rt.tag_id IN (%s)", strings.Join(placeholders, ",")))
	}

	sqlQuery += " WHERE " + strings.Join(conditions, " AND ")
	sqlQuery += " ORDER BY r.created_at DESC"

	rows, err := r.db.Pool.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search recipes: %w", err)
	}
	defer rows.Close()

	var recipes []entity.Recipe
	for rows.Next() {
		var recipe entity.Recipe
		if err := rows.Scan(&recipe.ID, &recipe.UserID, &recipe.Name, &recipe.URL, &recipe.Memo, &recipe.CreatedAt, &recipe.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan recipe: %w", err)
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (r *RecipeRepository) Update(ctx context.Context, recipe *entity.Recipe) error {
	query := `
		UPDATE recipes
		SET name = $1, url = $2, memo = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at
	`
	err := r.db.Pool.QueryRow(ctx, query, recipe.Name, recipe.URL, recipe.Memo, recipe.ID).
		Scan(&recipe.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update recipe: %w", err)
	}
	return nil
}

func (r *RecipeRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM recipes WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete recipe: %w", err)
	}
	return nil
}

func (r *RecipeRepository) AddTags(ctx context.Context, recipeID int64, tagIDs []int64) error {
	for _, tagID := range tagIDs {
		query := `
			INSERT INTO recipe_tags (recipe_id, tag_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`
		_, err := r.db.Pool.Exec(ctx, query, recipeID, tagID)
		if err != nil {
			return fmt.Errorf("failed to add tag to recipe: %w", err)
		}
	}
	return nil
}

func (r *RecipeRepository) RemoveTags(ctx context.Context, recipeID int64, tagIDs []int64) error {
	if len(tagIDs) == 0 {
		return nil
	}

	placeholders := make([]string, len(tagIDs))
	args := []interface{}{recipeID}
	for i, tagID := range tagIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args = append(args, tagID)
	}

	query := fmt.Sprintf(`
		DELETE FROM recipe_tags
		WHERE recipe_id = $1 AND tag_id IN (%s)
	`, strings.Join(placeholders, ","))

	_, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to remove tags from recipe: %w", err)
	}
	return nil
}

func (r *RecipeRepository) GetTags(ctx context.Context, recipeID int64) ([]entity.Tag, error) {
	query := `
		SELECT t.id, t.name
		FROM tags t
		INNER JOIN recipe_tags rt ON t.id = rt.tag_id
		WHERE rt.recipe_id = $1
		ORDER BY t.name
	`
	rows, err := r.db.Pool.Query(ctx, query, recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags for recipe: %w", err)
	}
	defer rows.Close()

	var tags []entity.Tag
	for rows.Next() {
		var tag entity.Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/infrastructure/entity"
)

type RecipeRepository struct {
	db Database
}

func NewRecipeRepository(db Database) *RecipeRepository {
	return &RecipeRepository{db: db}
}

func (r *RecipeRepository) Create(ctx context.Context, recipe *model.Recipe) error {
	e := entity.NewRecipe(recipe)
	query := `
		INSERT INTO recipes (user_id, name, url, memo, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	rows, err := r.db.Query(ctx, query,
		e.UserID, e.Name, e.URL, e.Memo,
	)
	if err != nil {
		return fmt.Errorf("failed to create recipe: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return fmt.Errorf("failed to create recipe: no rows returned")
	}

	err = rows.Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to scan recipe: %w", err)
	}

	recipe.ID = e.ID
	return nil
}

func (r *RecipeRepository) GetByID(ctx context.Context, id int64) (*model.Recipe, error) {
	query := `
		SELECT id, user_id, name, url, memo, created_at, updated_at
		FROM recipes
		WHERE id = $1
	`
	var e entity.Recipe
	// Manual Scan
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("no rows found")
	}

	err = rows.Scan(
		&e.ID, &e.UserID, &e.Name, &e.URL, &e.Memo, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe by id: %w", err)
	}
	return e.ToModel(), nil
}

func (r *RecipeRepository) GetByUserID(ctx context.Context, userID int64) ([]model.Recipe, error) {
	query := `
		SELECT id, user_id, name, url, memo, created_at, updated_at
		FROM recipes
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipes by user id: %w", err)
	}
	defer rows.Close()

	var recipeEntities []entity.Recipe
	for rows.Next() {
		var r entity.Recipe
		if err := rows.Scan(&r.ID, &r.UserID, &r.Name, &r.URL, &r.Memo, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan recipe: %w", err)
		}
		recipeEntities = append(recipeEntities, r)
	}
	return r.toRecipeModels(recipeEntities), nil
}

func (r *RecipeRepository) Search(ctx context.Context, userID int64, query string, tagIDs []int64) ([]model.Recipe, error) {
	sqlQuery := `
		SELECT DISTINCT r.id, r.user_id, r.name, r.url, r.memo, r.created_at, r.updated_at
		FROM recipes r
	`
	args := []interface{}{userID}
	whereConditions := []string{"r.user_id = $1"}
	paramIndex := 2

	if len(tagIDs) > 0 {
		sqlQuery += ` INNER JOIN recipe_tags rt ON r.id = rt.recipe_id`
		placeholders := make([]string, len(tagIDs))
		for i, tagID := range tagIDs {
			placeholders[i] = fmt.Sprintf("$%d", paramIndex)
			args = append(args, tagID)
			paramIndex++
		}
		whereConditions = append(whereConditions, fmt.Sprintf("rt.tag_id IN (%s)", strings.Join(placeholders, ",")))
	}

	if query != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("(r.name ILIKE $%d OR r.memo ILIKE $%d)", paramIndex, paramIndex))
		args = append(args, "%"+query+"%")
		paramIndex++
	}

	sqlQuery += " WHERE " + strings.Join(whereConditions, " AND ")
	sqlQuery += " ORDER BY r.created_at DESC"

	rows, err := r.db.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search recipes: %w", err)
	}
	defer rows.Close()

	var recipeEntities []entity.Recipe
	for rows.Next() {
		var rec entity.Recipe
		if err := rows.Scan(&rec.ID, &rec.UserID, &rec.Name, &rec.URL, &rec.Memo, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan recipe: %w", err)
		}
		recipeEntities = append(recipeEntities, rec)
	}
	return r.toRecipeModels(recipeEntities), nil
}

func (r *RecipeRepository) Update(ctx context.Context, recipe *model.Recipe) error {
	e := entity.NewRecipe(recipe)
	query := `
		UPDATE recipes
		SET name = $1, url = $2, memo = $3, updated_at = NOW()
		WHERE id = $4
	`
	_, err := r.db.Execute(ctx, query, e.Name, e.URL, e.Memo, e.ID)
	if err != nil {
		return fmt.Errorf("failed to update recipe: %w", err)
	}
	return nil
}

func (r *RecipeRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM recipes WHERE id = $1`
	_, err := r.db.Execute(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete recipe: %w", err)
	}
	return nil
}

func (r *RecipeRepository) GetTags(ctx context.Context, recipeID int64) ([]model.Tag, error) {
	query := `
		SELECT t.id, t.name
		FROM tags t
		INNER JOIN recipe_tags rt ON t.id = rt.tag_id
		WHERE rt.recipe_id = $1
		ORDER BY t.name
	`
	rows, err := r.db.Query(ctx, query, recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe tags: %w", err)
	}
	defer rows.Close()

	var tagEntities []entity.Tag
	for rows.Next() {
		var tag entity.Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tagEntities = append(tagEntities, tag)
	}
	return r.toTagModels(tagEntities), nil
}

func (r *RecipeRepository) AddTags(ctx context.Context, recipeID int64, tagIDs []int64) error {
	for _, tagID := range tagIDs {
		query := `
			INSERT INTO recipe_tags (recipe_id, tag_id)
			VALUES ($1, $2)
			ON CONFLICT (recipe_id, tag_id) DO NOTHING
		`
		_, err := r.db.Execute(ctx, query, recipeID, tagID)
		if err != nil {
			return fmt.Errorf("failed to add tag to recipe: %w", err)
		}
	}
	return nil
}

func (r *RecipeRepository) RemoveTags(ctx context.Context, recipeID int64, tagIDs []int64) error {
	for _, tagID := range tagIDs {
		query := `DELETE FROM recipe_tags WHERE recipe_id = $1 AND tag_id = $2`
		_, err := r.db.Execute(ctx, query, recipeID, tagID)
		if err != nil {
			return fmt.Errorf("failed to remove tag from recipe: %w", err)
		}
	}
	return nil
}

// Helpers

func (r *RecipeRepository) toRecipeModels(entities []entity.Recipe) []model.Recipe {
	models := make([]model.Recipe, len(entities))
	for i, e := range entities {
		models[i] = *e.ToModel()
	}
	return models
}

func (r *RecipeRepository) toTagModels(entities []entity.Tag) []model.Tag {
	models := make([]model.Tag, len(entities))
	for i, e := range entities {
		models[i] = *e.ToModel()
	}
	return models
}

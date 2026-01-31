package postgres

import (
	"context"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/infrastructure/entity"
)

type TagRepository struct {
	db *DB
}

func NewTagRepository(db *DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(ctx context.Context, tag *model.Tag) error {
	query := `
		INSERT INTO tags (name)
		VALUES ($1)
		RETURNING id
	`
	err := r.db.Pool.QueryRow(ctx, query, tag.Name).Scan(&tag.ID)
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}
	return nil
}

func (r *TagRepository) GetByID(ctx context.Context, id int64) (*model.Tag, error) {
	query := `SELECT id, name FROM tags WHERE id = $1`
	var tagEntity entity.Tag
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&tagEntity.ID, &tagEntity.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag by id: %w", err)
	}
	return tagEntityToModel(&tagEntity), nil
}

func (r *TagRepository) GetByName(ctx context.Context, name string) (*model.Tag, error) {
	query := `SELECT id, name FROM tags WHERE name = $1`
	var tagEntity entity.Tag
	err := r.db.Pool.QueryRow(ctx, query, name).Scan(&tagEntity.ID, &tagEntity.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag by name: %w", err)
	}
	return tagEntityToModel(&tagEntity), nil
}

func (r *TagRepository) GetAll(ctx context.Context) ([]model.Tag, error) {
	query := `SELECT id, name FROM tags ORDER BY name`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all tags: %w", err)
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
	return tagEntitiesToModels(tagEntities), nil
}

func (r *TagRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tags WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return nil
}

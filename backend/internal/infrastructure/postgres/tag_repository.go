package postgres

import (
	"context"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/infrastructure/entity"
)

type TagRepository struct {
	db Database
}

func NewTagRepository(db Database) repository.TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(ctx context.Context, tag *model.Tag) error {
	e := entity.NewTag(tag)
	query := `
		INSERT INTO tags (name)
		VALUES ($1)
		RETURNING id
	`
	rows, err := r.db.Query(ctx, query, e.Name)
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return fmt.Errorf("failed to create tag: no rows returned")
	}

	err = rows.Scan(&e.ID)
	if err != nil {
		return fmt.Errorf("failed to scan tag: %w", err)
	}
	tag.ID = e.ID
	return nil
}

func (r *TagRepository) GetByID(ctx context.Context, id int64) (*model.Tag, error) {
	query := `SELECT id, name FROM tags WHERE id = $1`
	var e entity.Tag
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag by id: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("tag not found")
	}

	err = rows.Scan(&e.ID, &e.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to scan tag: %w", err)
	}
	return e.ToModel(), nil
}

func (r *TagRepository) GetByName(ctx context.Context, name string) (*model.Tag, error) {
	query := `SELECT id, name FROM tags WHERE name = $1`
	var e entity.Tag
	rows, err := r.db.Query(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag by name: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("tag not found")
	}

	err = rows.Scan(&e.ID, &e.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to scan tag: %w", err)
	}
	return e.ToModel(), nil
}

func (r *TagRepository) GetAll(ctx context.Context) ([]model.Tag, error) {
	query := `SELECT id, name FROM tags ORDER BY name`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all tags: %w", err)
	}
	defer rows.Close()

	var entities []entity.Tag
	for rows.Next() {
		var e entity.Tag
		if err := rows.Scan(&e.ID, &e.Name); err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		entities = append(entities, e)
	}

	return r.toModels(entities), nil
}

func (r *TagRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tags WHERE id = $1`
	_, err := r.db.Execute(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return nil
}

func (r *TagRepository) toModels(entities []entity.Tag) []model.Tag {
	models := make([]model.Tag, len(entities))
	for i, e := range entities {
		models[i] = *e.ToModel()
	}
	return models
}

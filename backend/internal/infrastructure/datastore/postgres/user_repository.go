package postgres

import (
	"context"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/infrastructure/entity"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	e := entity.NewUserFromModel(user)
	query := `
		INSERT INTO users (name, created_at, updated_at)
		VALUES ($1, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.Pool.QueryRow(ctx, query, e.Name).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = e.ID
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var e entity.User // Use shared entity struct

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&e.ID,
		&e.Name,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return e.ToModel(), nil
}

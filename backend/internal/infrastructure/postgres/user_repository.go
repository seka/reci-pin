package postgres

import (
	"context"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/infrastructure/entity"
)

type UserRepository struct {
	db Database
}

func NewUserRepository(db Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	e := entity.NewUser(user)
	query := `
		INSERT INTO users (name, created_at, updated_at)
		VALUES ($1, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	rows, err := r.db.Query(ctx, query, e.Name)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return fmt.Errorf("failed to create user: no rows returned")
	}

	err = rows.Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to scan user: %w", err)
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
	var e entity.User

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("user not found")
	}

	err = rows.Scan(
		&e.ID,
		&e.Name,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	return e.ToModel(), nil
}

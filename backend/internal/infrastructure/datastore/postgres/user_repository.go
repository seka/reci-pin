package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/model"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	e := userModelToEntity(user)
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
	var entityUser struct {
		ID        int64     `db:"id"`
		Name      string    `db:"name"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&entityUser.ID,
		&entityUser.Name,
		&entityUser.CreatedAt,
		&entityUser.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	// 簡易変換（converter.goのuserEntityToModelを使うにはentity構造体が必要だが、ここでは内部構造体で受けているため手動マッピング）
	return &model.User{
		ID:   entityUser.ID,
		Name: entityUser.Name,
	}, nil
}

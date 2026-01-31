package postgres

import (
	"context"
	"fmt"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/infrastructure/entity"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User, passwordHash string) error {
	query := `
		INSERT INTO users (email, password_hash, name, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`
	err := r.db.Pool.QueryRow(ctx, query, user.Email, passwordHash, user.Name).
		Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT id, email, name
		FROM users
		WHERE id = $1
	`
	var userEntity entity.User
	err := r.db.Pool.QueryRow(ctx, query, id).
		Scan(&userEntity.ID, &userEntity.Email, &userEntity.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return userEntityToModel(&userEntity), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, string, error) {
	query := `
		SELECT id, email, password_hash, name
		FROM users
		WHERE email = $1
	`
	var userEntity entity.User
	err := r.db.Pool.QueryRow(ctx, query, email).
		Scan(&userEntity.ID, &userEntity.Email, &userEntity.PasswordHash, &userEntity.Name)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user by email: %w", err)
	}
	return userEntityToModel(&userEntity), userEntity.PasswordHash, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET email = $1, name = $2, updated_at = NOW()
		WHERE id = $3
	`
	_, err := r.db.Pool.Exec(ctx, query, user.Email, user.Name, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// HashPassword はパスワードをハッシュ化します
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash はパスワードとハッシュを比較します
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	"github.com/seka/reci-pin/backend/internal/infrastructure/entity"
	"golang.org/x/crypto/bcrypt"
)

type UserEmailCredentialRepository struct {
	db *DB
}

func NewUserEmailCredentialRepository(db *DB) repository.UserEmailCredentialRepository {
	return &UserEmailCredentialRepository{db: db}
}

func (r *UserEmailCredentialRepository) Create(ctx context.Context, credential *model.UserEmailCredential) error {
	e := entity.NewUserEmailCredential(credential)
	query := `
		INSERT INTO user_email_credentials (
			user_id, email, password_hash, email_verified_at, 
			verification_token, verification_token_expires_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`
	_, err := r.db.Pool.Exec(ctx, query,
		e.UserID, e.Email, e.PasswordHash, e.EmailVerifiedAt,
		e.VerificationToken, e.VerificationTokenExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user email credential: %w", err)
	}
	return nil
}

func (r *UserEmailCredentialRepository) GetByEmail(ctx context.Context, email string) (*model.UserEmailCredential, error) {
	query := `
		SELECT user_id, email, password_hash, email_verified_at, 
		       verification_token, verification_token_expires_at, updated_at
		FROM user_email_credentials
		WHERE email = $1
	`
	var e entity.UserEmailCredential
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&e.UserID, &e.Email, &e.PasswordHash, &e.EmailVerifiedAt,
		&e.VerificationToken, &e.VerificationTokenExpiresAt, &e.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to get credential by email: %w", err)
	}
	return e.ToModel(), nil
}

func (r *UserEmailCredentialRepository) GetByUserID(ctx context.Context, userID int64) (*model.UserEmailCredential, error) {
	query := `
		SELECT user_id, email, password_hash, email_verified_at, 
		       verification_token, verification_token_expires_at, updated_at
		FROM user_email_credentials
		WHERE user_id = $1
	`
	var e entity.UserEmailCredential
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&e.UserID, &e.Email, &e.PasswordHash, &e.EmailVerifiedAt,
		&e.VerificationToken, &e.VerificationTokenExpiresAt, &e.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get credential by user_id: %w", err)
	}
	return e.ToModel(), nil
}

func (r *UserEmailCredentialRepository) GetByToken(ctx context.Context, token string) (*model.UserEmailCredential, error) {
	query := `
		SELECT user_id, email, password_hash, email_verified_at, 
		       verification_token, verification_token_expires_at, updated_at
		FROM user_email_credentials
		WHERE verification_token = $1
	`
	var e entity.UserEmailCredential
	err := r.db.Pool.QueryRow(ctx, query, token).Scan(
		&e.UserID, &e.Email, &e.PasswordHash, &e.EmailVerifiedAt,
		&e.VerificationToken, &e.VerificationTokenExpiresAt, &e.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get credential by token: %w", err)
	}
	return e.ToModel(), nil
}

func (r *UserEmailCredentialRepository) Update(ctx context.Context, credential *model.UserEmailCredential) error {
	e := entity.NewUserEmailCredential(credential)
	query := `
		UPDATE user_email_credentials
		SET email = $2, password_hash = $3, email_verified_at = $4,
		    verification_token = $5, verification_token_expires_at = $6, updated_at = NOW()
		WHERE user_id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query,
		e.UserID, e.Email, e.PasswordHash, e.EmailVerifiedAt,
		e.VerificationToken, e.VerificationTokenExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update credential: %w", err)
	}
	return nil
}

// Password utility (Keep this here or move to a service?)
// Keeping here as part of infrastructure implementation details for now.

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

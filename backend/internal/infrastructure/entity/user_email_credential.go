package entity

import (
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/model"
)

type UserEmailCredential struct {
	UserID                     int64      `db:"user_id"`
	Email                      string     `db:"email"`
	PasswordHash               string     `db:"password_hash"`
	EmailVerifiedAt            *time.Time `db:"email_verified_at"`
	VerificationToken          string     `db:"verification_token"`
	VerificationTokenExpiresAt *time.Time `db:"verification_token_expires_at"`
	UpdatedAt                  time.Time  `db:"updated_at"`
}

func (e *UserEmailCredential) ToModel() *model.UserEmailCredential {
	if e == nil {
		return nil
	}
	return &model.UserEmailCredential{
		UserID:                     e.UserID,
		Email:                      e.Email,
		PasswordHash:               e.PasswordHash,
		EmailVerifiedAt:            e.EmailVerifiedAt,
		VerificationToken:          e.VerificationToken,
		VerificationTokenExpiresAt: e.VerificationTokenExpiresAt,
	}
}

func NewUserEmailCredential(m *model.UserEmailCredential) *UserEmailCredential {
	if m == nil {
		return nil
	}
	return &UserEmailCredential{
		UserID:                     m.UserID,
		Email:                      m.Email,
		PasswordHash:               m.PasswordHash,
		EmailVerifiedAt:            m.EmailVerifiedAt,
		VerificationToken:          m.VerificationToken,
		VerificationTokenExpiresAt: m.VerificationTokenExpiresAt,
	}
}

package entity

import "time"

type UserEmailCredential struct {
	UserID                     int64      `db:"user_id"`
	Email                      string     `db:"email"`
	PasswordHash               string     `db:"password_hash"`
	EmailVerifiedAt            *time.Time `db:"email_verified_at"`
	VerificationToken          string     `db:"verification_token"`
	VerificationTokenExpiresAt *time.Time `db:"verification_token_expires_at"`
	UpdatedAt                  time.Time  `db:"updated_at"`
}

package model

import "time"

// UserEmailCredential represents email/password authentication data
type UserEmailCredential struct {
	UserID                     int64
	Email                      string
	PasswordHash               string
	EmailVerifiedAt            *time.Time
	VerificationToken          string
	VerificationTokenExpiresAt *time.Time
}

func (c *UserEmailCredential) IsVerified() bool {
	return c.EmailVerifiedAt != nil
}

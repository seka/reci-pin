package model

import "time"

type AuthToken struct {
	Token     string
	ExpiresAt time.Time
}

type TokenResult struct {
	AccessToken  AuthToken
	RefreshToken AuthToken
}

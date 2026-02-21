package entity

import (
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/model"
)

type RefreshToken struct {
	ID        int64      `db:"id"`
	UserID    int64      `db:"user_id"`
	TokenHash string     `db:"token_hash"`
	ExpiresAt time.Time  `db:"expires_at"`
	CreatedAt time.Time  `db:"created_at"`
	RevokedAt *time.Time `db:"revoked_at"`
	UserAgent string     `db:"user_agent"`
	IPAddress string     `db:"ip_address"`
}

func NewRefreshToken(m *model.RefreshToken) *RefreshToken {
	return &RefreshToken{
		ID:        m.ID,
		UserID:    m.UserID,
		TokenHash: m.TokenHash,
		ExpiresAt: m.ExpiresAt,
		CreatedAt: m.CreatedAt,
		RevokedAt: m.RevokedAt,
		UserAgent: m.UserAgent,
		IPAddress: m.IPAddress,
	}
}

func (e *RefreshToken) ToModel() *model.RefreshToken {
	return &model.RefreshToken{
		ID:        e.ID,
		UserID:    e.UserID,
		TokenHash: e.TokenHash,
		ExpiresAt: e.ExpiresAt,
		CreatedAt: e.CreatedAt,
		RevokedAt: e.RevokedAt,
		UserAgent: e.UserAgent,
		IPAddress: e.IPAddress,
	}
}

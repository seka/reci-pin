package entity

import (
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/model"
)

type User struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (e *User) ToModel() *model.User {
	if e == nil {
		return nil
	}
	return &model.User{
		ID:   e.ID,
		Name: e.Name,
	}
}

func NewUser(m *model.User) *User {
	if m == nil {
		return nil
	}
	return &User{
		ID:   m.ID,
		Name: m.Name,
	}
}

package entity

import (
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/model"
)

type Recipe struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Name      string    `db:"name"`
	URL       string    `db:"url"`
	Memo      string    `db:"memo"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (e *Recipe) ToModel() *model.Recipe {
	if e == nil {
		return nil
	}
	return &model.Recipe{
		ID:     e.ID,
		UserID: e.UserID,
		Name:   e.Name,
		URL:    e.URL,
		Memo:   e.Memo,
	}
}

func NewRecipeFromModel(m *model.Recipe) *Recipe {
	if m == nil {
		return nil
	}
	return &Recipe{
		ID:     m.ID,
		UserID: m.UserID,
		Name:   m.Name,
		URL:    m.URL,
		Memo:   m.Memo,
	}
}

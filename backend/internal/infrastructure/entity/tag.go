package entity

import "github.com/seka/reci-pin/backend/internal/domain/model"

type Tag struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func (e *Tag) ToModel() *model.Tag {
	if e == nil {
		return nil
	}
	return &model.Tag{
		ID:   e.ID,
		Name: e.Name,
	}
}

func NewTag(m *model.Tag) *Tag {
	if m == nil {
		return nil
	}
	return &Tag{
		ID:   m.ID,
		Name: m.Name,
	}
}

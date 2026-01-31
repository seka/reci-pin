package entity

import "time"

type Recipe struct {
	ID        int64         `json:"id"`
	UserID    int64         `json:"user_id"`
	Name      string        `json:"name"`
	URL       string        `json:"url"`
	Memo      string        `json:"memo"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Tags      []Tag         `json:"tags,omitempty"`
	Images    []RecipeImage `json:"images,omitempty"`
}

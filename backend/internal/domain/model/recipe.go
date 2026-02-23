package model

import (
	"time"
)

// Recipe represents a recipe in business logic
type Recipe struct {
	ID        int64
	UserID    int64
	Name      string
	URL       string
	Memo      string
	Tags      []Tag
	Images    []PublicRecipeImage
	CreatedAt time.Time
	UpdatedAt time.Time
}

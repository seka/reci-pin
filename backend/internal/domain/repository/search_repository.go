package repository

import (
	"context"

	"github.com/seka/reci-pin/backend/internal/domain/model"
)

// SearchCriteria defines the criteria for searching recipes
//
//go:generate mockgen -source=$GOFILE -destination=./mock/search_repository_mock.go -package=mock
type SearchCriteria struct {
	UserID   int64
	Keyword  string
	TagIDs   []int64
	Page     int
	PageSize int
}

// RecipeSearchRepository defines the interface for searching recipes
type RecipeSearchRepository interface {
	// Index adds or updates a recipe in the search index
	Index(ctx context.Context, recipe *model.Recipe) error
	// Delete removes a recipe from the search index
	Delete(ctx context.Context, id int64) error
	// Search searches for recipes based on the criteria
	Search(ctx context.Context, criteria SearchCriteria) ([]int64, int64, error)
}

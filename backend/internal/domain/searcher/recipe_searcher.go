package searcher

import (
	"context"

	"github.com/seka/reci-pin/backend/internal/domain/model"
)

// RecipeSearcher defines the interface for searching recipes
//
//go:generate mockgen -source=$GOFILE -destination=./mock/recipe_searcher_mock.go -package=mock
type RecipeSearcher interface {
	// Index adds or updates a recipe in the search index
	Index(ctx context.Context, recipe *model.Recipe) error
	// Delete removes a recipe from the search index
	Delete(ctx context.Context, id int64) error
	// Search searches for recipes based on the criteria
	Search(ctx context.Context, criteria model.RecipeSearchCriteria) ([]int64, int64, error)
}

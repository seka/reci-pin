package model

// RecipeImage represents a recipe image in business logic
type RecipeImage struct {
	ID        int64
	RecipeID  int64
	ImagePath string
}

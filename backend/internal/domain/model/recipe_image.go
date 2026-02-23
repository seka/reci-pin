package model

import "net/url"

// RecipeImage represents a recipe image in business logic
type RecipeImage struct {
	ID        int64
	RecipeID  int64
	ImagePath string
}

// PublicRecipeImage represents a recipe image with its full public URL
type PublicRecipeImage struct {
	RecipeImage
	ImageURL url.URL
}

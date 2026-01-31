package model

// Recipe represents a recipe in business logic
type Recipe struct {
	ID     int64
	UserID int64
	Name   string
	URL    string
	Memo   string
	Tags   []Tag
	Images []RecipeImage
}

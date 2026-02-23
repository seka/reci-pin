package model

// RecipeSearchCriteria defines the criteria for searching recipes
type RecipeSearchCriteria struct {
	UserID   int64
	Keyword  string
	TagIDs   []int64
	Page     int
	PageSize int
}

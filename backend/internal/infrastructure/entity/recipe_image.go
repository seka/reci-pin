package entity

import "time"

type RecipeImage struct {
	ID        int64     `json:"id"`
	RecipeID  int64     `json:"recipe_id"`
	ImagePath string    `json:"image_path"`
	CreatedAt time.Time `json:"created_at"`
}

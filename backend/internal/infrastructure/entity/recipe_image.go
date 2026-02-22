package entity

import (
	"net/url"
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/model"
)

type RecipeImage struct {
	ID        int64     `db:"id"`
	RecipeID  int64     `db:"recipe_id"`
	ImagePath string    `db:"image_path"`
	CreatedAt time.Time `db:"created_at"`
}

func (e *RecipeImage) ToModel() *model.RecipeImage {
	if e == nil {
		return nil
	}
	parsedURL, _ := url.Parse(e.ImagePath)
	return &model.RecipeImage{
		ID:       e.ID,
		RecipeID: e.RecipeID,
		ImageURL: *parsedURL,
	}
}

func NewRecipeImage(m *model.RecipeImage) *RecipeImage {
	if m == nil {
		return nil
	}
	return &RecipeImage{
		ID:        m.ID,
		RecipeID:  m.RecipeID,
		ImagePath: m.ImageURL.Path,
	}
}

package response

import (
	"github.com/seka/reci-pin/backend/internal/domain/model"
)

type RecipeResponse struct {
	ID     int64                 `json:"id"`
	UserID int64                 `json:"user_id"`
	Name   string                `json:"name"`
	URL    string                `json:"url"`
	Memo   string                `json:"memo"`
	Tags   []TagResponse         `json:"tags"`
	Images []RecipeImageResponse `json:"images"`
}

type TagResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type RecipeImageResponse struct {
	ID        int64  `json:"id"`
	RecipeID  int64  `json:"recipe_id"`
	ImagePath string `json:"image_path"`
}

type CreateRecipeImageResponse struct {
	Image     RecipeImageResponse `json:"image"`
	UploadURL string              `json:"upload_url"`
}

func NewRecipe(m *model.Recipe) *RecipeResponse {
	tags := make([]TagResponse, 0, len(m.Tags))
	for _, t := range m.Tags {
		tags = append(tags, TagResponse{
			ID:   t.ID,
			Name: t.Name,
		})
	}

	images := make([]RecipeImageResponse, 0, len(m.Images))
	for _, i := range m.Images {
		images = append(images, RecipeImageResponse{
			ID:        i.ID,
			RecipeID:  i.RecipeID,
			ImagePath: i.ImagePath,
		})
	}

	return &RecipeResponse{
		ID:     m.ID,
		UserID: m.UserID,
		Name:   m.Name,
		URL:    m.URL,
		Memo:   m.Memo,
		Tags:   tags,
		Images: images,
	}
}

func NewRecipes(recipes []model.Recipe) []RecipeResponse {
	responses := make([]RecipeResponse, 0, len(recipes))
	for _, r := range recipes {
		responses = append(responses, *NewRecipe(&r))
	}
	return responses
}

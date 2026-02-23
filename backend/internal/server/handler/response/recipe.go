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
	ID       int64  `json:"id"`
	RecipeID int64  `json:"recipe_id"`
	ImageURL string `json:"image_url"`
}

type CreateRecipeImageResponse struct {
	Image     RecipeImageResponse `json:"image"`
	UploadURL string              `json:"upload_url"`
}

func NewRecipeImageResponse(img model.PublicRecipeImage) RecipeImageResponse {
	return RecipeImageResponse{
		ID:       img.ID,
		RecipeID: img.RecipeID,
		ImageURL: img.ImageURL.String(),
	}
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
		images = append(images, NewRecipeImageResponse(i))
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

func NewRecipeResponse(recipe *model.Recipe) *RecipeResponse {
	// This function seems to be a duplicate of NewRecipe, but the instruction explicitly asks for it.
	// The image handling here is also slightly different (slice of pointers vs slice of values).
	// I will implement it as provided in the snippet, assuming `recipe.Images` contains `model.PublicRecipeImage`.
	images := make([]RecipeImageResponse, len(recipe.Images))
	for i := range recipe.Images {
		images[i] = NewRecipeImageResponse(recipe.Images[i])
	}

	tags := make([]TagResponse, 0, len(recipe.Tags))
	for _, t := range recipe.Tags {
		tags = append(tags, TagResponse{
			ID:   t.ID,
			Name: t.Name,
		})
	}

	return &RecipeResponse{
		ID:     recipe.ID,
		UserID: recipe.UserID,
		Name:   recipe.Name,
		URL:    recipe.URL,
		Memo:   recipe.Memo,
		Tags:   tags,
		Images: images,
	}
}

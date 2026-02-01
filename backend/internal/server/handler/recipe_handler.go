package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/seka/reci-pin/backend/internal/server/middleware"
	"github.com/seka/reci-pin/backend/internal/usecase/recipe"
	"github.com/seka/reci-pin/backend/internal/usecase/recipe_image"
	"github.com/seka/reci-pin/backend/internal/usecase/recipe_tag"
	"github.com/seka/reci-pin/backend/internal/usecase/tag"
)

type RecipeHandler struct {
	createRecipeUseCase   recipe.CreateRecipeUseCase
	getRecipeUseCase      recipe.GetRecipeUseCase
	getUserRecipesUseCase recipe.GetUserRecipesUseCase
	updateRecipeUseCase   recipe.UpdateRecipeUseCase
	deleteRecipeUseCase   recipe.DeleteRecipeUseCase
	searchRecipesUseCase  recipe.SearchRecipesUseCase
	addTagsUseCase        recipe_tag.AddTagsUseCase
	removeTagsUseCase     recipe_tag.RemoveTagsUseCase
	addImageUseCase       recipe_image.AddImageUseCase
	createTagUseCase      tag.CreateTagUseCase
	getAllTagsUseCase     tag.GetAllTagsUseCase
	deleteTagUseCase      tag.DeleteTagUseCase
}

func NewRecipeHandler(
	createRecipeUseCase recipe.CreateRecipeUseCase,
	getRecipeUseCase recipe.GetRecipeUseCase,
	getUserRecipesUseCase recipe.GetUserRecipesUseCase,
	updateRecipeUseCase recipe.UpdateRecipeUseCase,
	deleteRecipeUseCase recipe.DeleteRecipeUseCase,
	searchRecipesUseCase recipe.SearchRecipesUseCase,
	addTagsUseCase recipe_tag.AddTagsUseCase,
	removeTagsUseCase recipe_tag.RemoveTagsUseCase,
	addImageUseCase recipe_image.AddImageUseCase,
	createTagUseCase tag.CreateTagUseCase,
	getAllTagsUseCase tag.GetAllTagsUseCase,
	deleteTagUseCase tag.DeleteTagUseCase,
) *RecipeHandler {
	return &RecipeHandler{
		createRecipeUseCase:   createRecipeUseCase,
		getRecipeUseCase:      getRecipeUseCase,
		getUserRecipesUseCase: getUserRecipesUseCase,
		updateRecipeUseCase:   updateRecipeUseCase,
		deleteRecipeUseCase:   deleteRecipeUseCase,
		searchRecipesUseCase:  searchRecipesUseCase,
		addTagsUseCase:        addTagsUseCase,
		removeTagsUseCase:     removeTagsUseCase,
		addImageUseCase:       addImageUseCase,
		createTagUseCase:      createTagUseCase,
		getAllTagsUseCase:     getAllTagsUseCase,
		deleteTagUseCase:      deleteTagUseCase,
	}
}

type CreateRecipeRequest struct {
	Name   string  `json:"name"`
	URL    string  `json:"url"`
	Memo   string  `json:"memo"`
	TagIDs []int64 `json:"tag_ids"`
}

type UpdateRecipeRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Memo string `json:"memo"`
}

type SearchRecipeRequest struct {
	Query  string  `json:"query"`
	TagIDs []int64 `json:"tag_ids"`
}

type AddTagsRequest struct {
	TagIDs []int64 `json:"tag_ids"`
}

type CreateTagRequest struct {
	Name string `json:"name"`
}

type AddImageRequest struct {
	ImagePath string `json:"image_path"`
}

func (h *RecipeHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	input := recipe.CreateRecipeInput{
		UserID: userID,
		Name:   req.Name,
		URL:    req.URL,
		Memo:   req.Memo,
		TagIDs: req.TagIDs,
	}

	result, err := h.createRecipeUseCase.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func (h *RecipeHandler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid recipe id", http.StatusBadRequest)
		return
	}

	result, err := h.getRecipeUseCase.Execute(r.Context(), id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *RecipeHandler) GetUserRecipes(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	recipes, err := h.getUserRecipesUseCase.Execute(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

func (h *RecipeHandler) SearchRecipes(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req SearchRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	input := recipe.SearchRecipesInput{
		UserID: userID,
		Query:  req.Query,
		TagIDs: req.TagIDs,
	}

	recipes, err := h.searchRecipesUseCase.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

func (h *RecipeHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid recipe id", http.StatusBadRequest)
		return
	}

	var req UpdateRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	input := recipe.UpdateRecipeInput{
		ID:     id,
		UserID: userID,
		Name:   req.Name,
		URL:    req.URL,
		Memo:   req.Memo,
	}

	result, err := h.updateRecipeUseCase.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *RecipeHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid recipe id", http.StatusBadRequest)
		return
	}

	if err := h.deleteRecipeUseCase.Execute(r.Context(), id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RecipeHandler) AddTags(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid recipe id", http.StatusBadRequest)
		return
	}

	var req AddTagsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.addTagsUseCase.Execute(r.Context(), id, userID, req.TagIDs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RecipeHandler) RemoveTags(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid recipe id", http.StatusBadRequest)
		return
	}

	var req AddTagsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.removeTagsUseCase.Execute(r.Context(), id, userID, req.TagIDs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RecipeHandler) AddImage(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid recipe id", http.StatusBadRequest)
		return
	}

	var req AddImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	image, err := h.addImageUseCase.Execute(r.Context(), id, userID, req.ImagePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(image)
}

// Tag endpoints

func (h *RecipeHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	var req CreateTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	tag, err := h.createTagUseCase.Execute(r.Context(), req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tag)
}

func (h *RecipeHandler) GetAllTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.getAllTagsUseCase.Execute(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

func (h *RecipeHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid tag id", http.StatusBadRequest)
		return
	}

	if err := h.deleteTagUseCase.Execute(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

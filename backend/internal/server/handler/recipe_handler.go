package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/seka/reci-pin/backend/internal/domain/validation"
	"github.com/seka/reci-pin/backend/internal/server/handler/request"
	"github.com/seka/reci-pin/backend/internal/server/handler/response"
	"github.com/seka/reci-pin/backend/internal/server/middleware"
	"github.com/seka/reci-pin/backend/internal/usecase/recipe"
	"github.com/seka/reci-pin/backend/internal/usecase/recipe_image"
	"github.com/seka/reci-pin/backend/internal/usecase/recipe_tag"
	"github.com/seka/reci-pin/backend/internal/usecase/tag"
)

type RecipeHandler struct {
	createRecipeUseCase      recipe.CreateRecipeUseCase
	getRecipeUseCase         recipe.GetRecipeUseCase
	getUserRecipesUseCase    recipe.GetUserRecipesUseCase
	updateRecipeUseCase      recipe.UpdateRecipeUseCase
	deleteRecipeUseCase      recipe.DeleteRecipeUseCase
	searchRecipesUseCase     recipe.SearchRecipesUseCase
	addTagsUseCase           recipe_tag.AddTagsUseCase
	removeTagsUseCase        recipe_tag.RemoveTagsUseCase
	createRecipeImageUseCase recipe_image.CreateRecipeImageUseCase
	createTagUseCase         tag.CreateTagUseCase
	getAllTagsUseCase        tag.GetAllTagsUseCase
	deleteTagUseCase         tag.DeleteTagUseCase
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
	createRecipeImageUseCase recipe_image.CreateRecipeImageUseCase,
	createTagUseCase tag.CreateTagUseCase,
	getAllTagsUseCase tag.GetAllTagsUseCase,
	deleteTagUseCase tag.DeleteTagUseCase,
) (*RecipeHandler, error) {
	return &RecipeHandler{
		createRecipeUseCase:      createRecipeUseCase,
		getRecipeUseCase:         getRecipeUseCase,
		getUserRecipesUseCase:    getUserRecipesUseCase,
		updateRecipeUseCase:      updateRecipeUseCase,
		deleteRecipeUseCase:      deleteRecipeUseCase,
		searchRecipesUseCase:     searchRecipesUseCase,
		addTagsUseCase:           addTagsUseCase,
		removeTagsUseCase:        removeTagsUseCase,
		createRecipeImageUseCase: createRecipeImageUseCase,
		createTagUseCase:         createTagUseCase,
		getAllTagsUseCase:        getAllTagsUseCase,
		deleteTagUseCase:         deleteTagUseCase,
	}, nil
}

func (h *RecipeHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req request.CreateRecipeRequest
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
		var validationErrors validation.ValidationErrors
		if errors.As(err, &validationErrors) {
			details := make(map[string][]response.ErrorDetail)
			for _, ve := range validationErrors {
				details[ve.Field] = append(details[ve.Field], response.ErrorDetail{
					Code:   ve.Code,
					Params: ve.Params,
				})
			}
			respondError(w, http.StatusBadRequest, "VALIDATION_FAILED", details)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusCreated, response.NewRecipe(result))
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

	respondJSON(w, http.StatusOK, response.NewRecipe(result))
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

	respondJSON(w, http.StatusOK, response.NewRecipes(recipes))
}

func (h *RecipeHandler) SearchRecipes(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req request.SearchRecipeRequest
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

	respondJSON(w, http.StatusOK, response.NewRecipes(recipes))
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

	var req request.UpdateRecipeRequest
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

	respondJSON(w, http.StatusOK, response.NewRecipe(result))
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

	var req request.AddTagsRequest
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

	var req request.AddTagsRequest // Reusing AddTagsRequest for Remove as struct is same
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

	var req request.CreateRecipeImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	input := recipe_image.CreateRecipeImageInput{
		RecipeID:    id,
		UserID:      userID,
		Filename:    req.Filename,
		ContentType: req.ContentType,
		Size:        req.Size,
	}

	image, uploadURL, err := h.createRecipeImageUseCase.Execute(r.Context(), input)
	if err != nil {
		var validationErrors validation.ValidationErrors
		if errors.As(err, &validationErrors) {
			details := make(map[string][]response.ErrorDetail)
			for _, ve := range validationErrors {
				details[ve.Field] = append(details[ve.Field], response.ErrorDetail{
					Code:   ve.Code,
					Params: ve.Params,
				})
			}
			respondError(w, http.StatusBadRequest, "VALIDATION_FAILED", details)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := response.CreateRecipeImageResponse{
		Image:     response.NewRecipeImageResponse(*image),
		UploadURL: uploadURL,
	}

	respondJSON(w, http.StatusCreated, resp)
}

// Tag endpoints

func (h *RecipeHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	var req request.CreateTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	tag, err := h.createTagUseCase.Execute(r.Context(), req.Name)
	if err != nil {
		var validationErrors validation.ValidationErrors
		if errors.As(err, &validationErrors) {
			details := make(map[string][]response.ErrorDetail)
			for _, ve := range validationErrors {
				details[ve.Field] = append(details[ve.Field], response.ErrorDetail{
					Code:   ve.Code,
					Params: ve.Params,
				})
			}
			respondError(w, http.StatusBadRequest, "VALIDATION_FAILED", details)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusCreated, tag)
}

func (h *RecipeHandler) GetAllTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.getAllTagsUseCase.Execute(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, tags)
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

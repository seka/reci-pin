package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/server/handler"
	"github.com/seka/reci-pin/backend/internal/server/middleware"
	usecasemock "github.com/seka/reci-pin/backend/internal/usecase/mock"
	"github.com/seka/reci-pin/backend/internal/usecase/recipe_image"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRecipeHandler_CreateRecipe(t *testing.T) {
	type args struct {
		body   map[string]any
		userID int64
	}
	type mocks struct {
		setup func(m *usecasemock.MockCreateRecipeUseCase)
	}
	tests := []struct {
		name       string
		args       args
		mocks      mocks
		wantStatus int
	}{
		{
			name: "Success",
			args: args{
				body:   map[string]any{"name": "Pancakes", "url": "http://example.com"},
				userID: 1,
			},
			mocks: mocks{
				setup: func(m *usecasemock.MockCreateRecipeUseCase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any()).
						Return(&model.Recipe{ID: 100, Name: "Pancakes"}, nil)
				},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Unauthorized",
			args: args{
				body:   map[string]any{"name": "Pancakes"},
				userID: 0,
			},
			mocks: mocks{
				setup: func(m *usecasemock.MockCreateRecipeUseCase) {},
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid JSON",
			args: args{
				body:   nil, // simulated by sending invalid bytes manually
				userID: 1,
			},
			mocks: mocks{
				setup: func(m *usecasemock.MockCreateRecipeUseCase) {},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "UseCase Error",
			args: args{
				body:   map[string]any{"name": "Pancakes"},
				userID: 1,
			},
			mocks: mocks{
				setup: func(m *usecasemock.MockCreateRecipeUseCase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any()).
						Return(nil, errors.New("db error"))
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCreateRecipe := usecasemock.NewMockCreateRecipeUseCase(ctrl)
			tt.mocks.setup(mockCreateRecipe)

			h, _ := handler.NewRecipeHandler(
				mockCreateRecipe,
				usecasemock.NewMockGetRecipeUseCase(ctrl),
				usecasemock.NewMockGetUserRecipesUseCase(ctrl),
				usecasemock.NewMockUpdateRecipeUseCase(ctrl),
				usecasemock.NewMockDeleteRecipeUseCase(ctrl),
				usecasemock.NewMockSearchRecipesUseCase(ctrl),
				usecasemock.NewMockAddTagsUseCase(ctrl),
				usecasemock.NewMockRemoveTagsUseCase(ctrl),
				usecasemock.NewMockCreateRecipeImageUseCase(ctrl),
				usecasemock.NewMockCreateTagUseCase(ctrl),
				usecasemock.NewMockGetAllTagsUseCase(ctrl),
				usecasemock.NewMockDeleteTagUseCase(ctrl),
			)

			var req *http.Request
			if tt.name == "Invalid JSON" {
				req = httptest.NewRequest(http.MethodPost, "/api/recipes", bytes.NewReader([]byte("invalid")))
			} else {
				body, _ := json.Marshal(tt.args.body)
				req = httptest.NewRequest(http.MethodPost, "/api/recipes", bytes.NewReader(body))
			}
			req.Header.Set("Content-Type", "application/json")

			if tt.args.userID != 0 {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.args.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			h.CreateRecipe(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestRecipeHandler_GetRecipe(t *testing.T) {
	type args struct {
		recipeID string
		userID   int64
	}
	type mocks struct {
		setup func(m *usecasemock.MockGetRecipeUseCase)
	}
	tests := []struct {
		name       string
		args       args
		mocks      mocks
		wantStatus int
	}{
		{
			name: "Success",
			args: args{recipeID: "100", userID: 1},
			mocks: mocks{
				setup: func(m *usecasemock.MockGetRecipeUseCase) {
					m.EXPECT().Execute(gomock.Any(), int64(100), int64(1)).
						Return(&model.Recipe{ID: 100}, nil)
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Invalid ID",
			args: args{recipeID: "abc", userID: 1},
			mocks: mocks{
				setup: func(m *usecasemock.MockGetRecipeUseCase) {},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Not Found",
			args: args{recipeID: "999", userID: 1},
			mocks: mocks{
				setup: func(m *usecasemock.MockGetRecipeUseCase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil, errors.New("not found"))
				},
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockGetRecipe := usecasemock.NewMockGetRecipeUseCase(ctrl)
			tt.mocks.setup(mockGetRecipe)

			h, _ := handler.NewRecipeHandler(
				usecasemock.NewMockCreateRecipeUseCase(ctrl),
				mockGetRecipe,
				usecasemock.NewMockGetUserRecipesUseCase(ctrl),
				usecasemock.NewMockUpdateRecipeUseCase(ctrl),
				usecasemock.NewMockDeleteRecipeUseCase(ctrl),
				usecasemock.NewMockSearchRecipesUseCase(ctrl),
				usecasemock.NewMockAddTagsUseCase(ctrl),
				usecasemock.NewMockRemoveTagsUseCase(ctrl),
				usecasemock.NewMockCreateRecipeImageUseCase(ctrl),
				usecasemock.NewMockCreateTagUseCase(ctrl),
				usecasemock.NewMockGetAllTagsUseCase(ctrl),
				usecasemock.NewMockDeleteTagUseCase(ctrl),
			)

			req := httptest.NewRequest(http.MethodGet, "/api/recipes/"+tt.args.recipeID, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.args.recipeID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			if tt.args.userID != 0 {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.args.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			h.GetRecipe(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestRecipeHandler_UpdateRecipe(t *testing.T) {
	type args struct {
		recipeID string
		body     map[string]any
		userID   int64
	}
	type mocks struct {
		setup func(m *usecasemock.MockUpdateRecipeUseCase)
	}
	tests := []struct {
		name       string
		args       args
		mocks      mocks
		wantStatus int
	}{
		{
			name: "Success",
			args: args{
				recipeID: "100",
				body:     map[string]any{"name": "Updated"},
				userID:   1,
			},
			mocks: mocks{
				setup: func(m *usecasemock.MockUpdateRecipeUseCase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any()).
						Return(&model.Recipe{ID: 100}, nil)
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Invalid ID",
			args: args{recipeID: "abc", userID: 1},
			mocks: mocks{
				setup: func(m *usecasemock.MockUpdateRecipeUseCase) {},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid JSON",
			args: args{recipeID: "100", userID: 1, body: nil},
			mocks: mocks{
				setup: func(m *usecasemock.MockUpdateRecipeUseCase) {},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Not Found",
			args: args{recipeID: "999", userID: 1, body: map[string]any{"name": "Updated"}},
			mocks: mocks{
				setup: func(m *usecasemock.MockUpdateRecipeUseCase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any()).
						Return(nil, errors.New("not found"))
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUpdate := usecasemock.NewMockUpdateRecipeUseCase(ctrl)
			tt.mocks.setup(mockUpdate)

			h, _ := handler.NewRecipeHandler(
				usecasemock.NewMockCreateRecipeUseCase(ctrl),
				usecasemock.NewMockGetRecipeUseCase(ctrl),
				usecasemock.NewMockGetUserRecipesUseCase(ctrl),
				mockUpdate,
				usecasemock.NewMockDeleteRecipeUseCase(ctrl),
				usecasemock.NewMockSearchRecipesUseCase(ctrl),
				usecasemock.NewMockAddTagsUseCase(ctrl),
				usecasemock.NewMockRemoveTagsUseCase(ctrl),
				usecasemock.NewMockCreateRecipeImageUseCase(ctrl),
				usecasemock.NewMockCreateTagUseCase(ctrl),
				usecasemock.NewMockGetAllTagsUseCase(ctrl),
				usecasemock.NewMockDeleteTagUseCase(ctrl),
			)

			var req *http.Request
			if tt.name == "Invalid JSON" {
				req = httptest.NewRequest(http.MethodPut, "/api/recipes/"+tt.args.recipeID, bytes.NewReader([]byte("invalid")))
			} else {
				body, _ := json.Marshal(tt.args.body)
				req = httptest.NewRequest(http.MethodPut, "/api/recipes/"+tt.args.recipeID, bytes.NewReader(body))
			}
			req.Header.Set("Content-Type", "application/json")

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.args.recipeID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			if tt.args.userID != 0 {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.args.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			h.UpdateRecipe(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestRecipeHandler_DeleteRecipe(t *testing.T) {
	type args struct {
		recipeID string
		userID   int64
	}
	type mocks struct {
		setup func(m *usecasemock.MockDeleteRecipeUseCase)
	}
	tests := []struct {
		name       string
		args       args
		mocks      mocks
		wantStatus int
	}{
		{
			name: "Success",
			args: args{recipeID: "100", userID: 1},
			mocks: mocks{
				setup: func(m *usecasemock.MockDeleteRecipeUseCase) {
					m.EXPECT().Execute(gomock.Any(), int64(100), int64(1)).Return(nil)
				},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "Invalid ID",
			args: args{recipeID: "abc", userID: 1},
			mocks: mocks{
				setup: func(m *usecasemock.MockDeleteRecipeUseCase) {},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Error",
			args: args{recipeID: "100", userID: 1},
			mocks: mocks{
				setup: func(m *usecasemock.MockDeleteRecipeUseCase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(errors.New("failed"))
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDelete := usecasemock.NewMockDeleteRecipeUseCase(ctrl)
			tt.mocks.setup(mockDelete)

			h, _ := handler.NewRecipeHandler(
				usecasemock.NewMockCreateRecipeUseCase(ctrl),
				usecasemock.NewMockGetRecipeUseCase(ctrl),
				usecasemock.NewMockGetUserRecipesUseCase(ctrl),
				usecasemock.NewMockUpdateRecipeUseCase(ctrl),
				mockDelete,
				usecasemock.NewMockSearchRecipesUseCase(ctrl),
				usecasemock.NewMockAddTagsUseCase(ctrl),
				usecasemock.NewMockRemoveTagsUseCase(ctrl),
				usecasemock.NewMockCreateRecipeImageUseCase(ctrl),
				usecasemock.NewMockCreateTagUseCase(ctrl),
				usecasemock.NewMockGetAllTagsUseCase(ctrl),
				usecasemock.NewMockDeleteTagUseCase(ctrl),
			)

			req := httptest.NewRequest(http.MethodDelete, "/api/recipes/"+tt.args.recipeID, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.args.recipeID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			if tt.args.userID != 0 {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.args.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			h.DeleteRecipe(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestRecipeHandler_SearchRecipes(t *testing.T) {
	type args struct {
		body   map[string]any
		userID int64
	}
	type mocks struct {
		setup func(m *usecasemock.MockSearchRecipesUseCase)
	}
	tests := []struct {
		name       string
		args       args
		mocks      mocks
		wantStatus int
	}{
		{
			name: "Success",
			args: args{
				body:   map[string]any{"query": "test"},
				userID: 1,
			},
			mocks: mocks{
				setup: func(m *usecasemock.MockSearchRecipesUseCase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any()).
						Return([]model.Recipe{}, nil)
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Invalid JSON",
			args: args{body: nil, userID: 1},
			mocks: mocks{
				setup: func(m *usecasemock.MockSearchRecipesUseCase) {},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Error",
			args: args{
				body:   map[string]any{"query": "test"},
				userID: 1,
			},
			mocks: mocks{
				setup: func(m *usecasemock.MockSearchRecipesUseCase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any()).
						Return(nil, errors.New("search failed"))
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSearch := usecasemock.NewMockSearchRecipesUseCase(ctrl)
			tt.mocks.setup(mockSearch)

			h, _ := handler.NewRecipeHandler(
				usecasemock.NewMockCreateRecipeUseCase(ctrl),
				usecasemock.NewMockGetRecipeUseCase(ctrl),
				usecasemock.NewMockGetUserRecipesUseCase(ctrl),
				usecasemock.NewMockUpdateRecipeUseCase(ctrl),
				usecasemock.NewMockDeleteRecipeUseCase(ctrl),
				mockSearch,
				usecasemock.NewMockAddTagsUseCase(ctrl),
				usecasemock.NewMockRemoveTagsUseCase(ctrl),
				usecasemock.NewMockCreateRecipeImageUseCase(ctrl),
				usecasemock.NewMockCreateTagUseCase(ctrl),
				usecasemock.NewMockGetAllTagsUseCase(ctrl),
				usecasemock.NewMockDeleteTagUseCase(ctrl),
			)

			var req *http.Request
			if tt.name == "Invalid JSON" {
				req = httptest.NewRequest(http.MethodPost, "/api/recipes/search", bytes.NewReader([]byte("invalid")))
			} else {
				body, _ := json.Marshal(tt.args.body)
				req = httptest.NewRequest(http.MethodPost, "/api/recipes/search", bytes.NewReader(body))
			}
			req.Header.Set("Content-Type", "application/json")

			if tt.args.userID != 0 {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.args.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			h.SearchRecipes(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestRecipeHandler_GetUserRecipes(t *testing.T) {
	type args struct {
		requestUserID string
		authUserID    int64
	}
	type mocks struct {
		setup func(m *usecasemock.MockGetUserRecipesUseCase)
	}
	tests := []struct {
		name       string
		args       args
		mocks      mocks
		wantStatus int
	}{
		{
			name: "Success",
			args: args{requestUserID: "1", authUserID: 1},
			mocks: mocks{
				setup: func(m *usecasemock.MockGetUserRecipesUseCase) {
					m.EXPECT().Execute(gomock.Any(), int64(1)).
						Return([]model.Recipe{}, nil)
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Unauthorized",
			args: args{requestUserID: "1", authUserID: 0},
			mocks: mocks{
				setup: func(m *usecasemock.MockGetUserRecipesUseCase) {},
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Error",
			args: args{requestUserID: "1", authUserID: 1},
			mocks: mocks{
				setup: func(m *usecasemock.MockGetUserRecipesUseCase) {
					m.EXPECT().Execute(gomock.Any(), int64(1)).
						Return(nil, errors.New("failed"))
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockGetUserRecipes := usecasemock.NewMockGetUserRecipesUseCase(ctrl)
			tt.mocks.setup(mockGetUserRecipes)

			h, _ := handler.NewRecipeHandler(
				usecasemock.NewMockCreateRecipeUseCase(ctrl),
				usecasemock.NewMockGetRecipeUseCase(ctrl),
				mockGetUserRecipes,
				usecasemock.NewMockUpdateRecipeUseCase(ctrl),
				usecasemock.NewMockDeleteRecipeUseCase(ctrl),
				usecasemock.NewMockSearchRecipesUseCase(ctrl),
				usecasemock.NewMockAddTagsUseCase(ctrl),
				usecasemock.NewMockRemoveTagsUseCase(ctrl),
				usecasemock.NewMockCreateRecipeImageUseCase(ctrl),
				usecasemock.NewMockCreateTagUseCase(ctrl),
				usecasemock.NewMockGetAllTagsUseCase(ctrl),
				usecasemock.NewMockDeleteTagUseCase(ctrl),
			)

			req := httptest.NewRequest(http.MethodGet, "/api/users/"+tt.args.requestUserID+"/recipes", nil)
			if tt.args.authUserID != 0 {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.args.authUserID)
				req = req.WithContext(ctx)
			}
			// In GetUserRecipes handler, it uses middleware.GetUserIDFromContext only.
			// It doesn't seem to use URL param for user ID in logic, but let's confirm handler impl.
			// Ah, handler uses middleware Auth UserID directly. It ignores {user_id} in path for authorization?
			// Let's check impl: `recipes, err := h.getUserRecipesUseCase.Execute(r.Context(), userID)`
			// `userID` comes from `middleware.GetUserIDFromContext`. So path param is ignored or assumed to match?
			// The handler implementation indeed uses the context userID.

			w := httptest.NewRecorder()
			h.GetUserRecipes(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestRecipeHandler_AddTags(t *testing.T) {
	type args struct {
		recipeID string
		body     map[string]any
		userID   int64
	}
	type mocks struct {
		setup func(m *usecasemock.MockAddTagsUseCase)
	}
	tests := []struct {
		name       string
		args       args
		mocks      mocks
		wantStatus int
	}{
		{
			name: "Success",
			args: args{recipeID: "100", body: map[string]any{"tag_ids": []int64{10}}, userID: 1},
			mocks: mocks{
				setup: func(m *usecasemock.MockAddTagsUseCase) {
					m.EXPECT().Execute(gomock.Any(), int64(100), int64(1), []int64{10}).
						Return(nil)
				},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "Invalid ID",
			args:       args{recipeID: "abc", userID: 1},
			mocks:      mocks{setup: func(m *usecasemock.MockAddTagsUseCase) {}},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Invalid JSON",
			args:       args{recipeID: "100", userID: 1, body: nil},
			mocks:      mocks{setup: func(m *usecasemock.MockAddTagsUseCase) {}},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Error",
			args: args{recipeID: "100", body: map[string]any{"tag_ids": []int64{10}}, userID: 1},
			mocks: mocks{
				setup: func(m *usecasemock.MockAddTagsUseCase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(errors.New("failed"))
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAddTags := usecasemock.NewMockAddTagsUseCase(ctrl)
			tt.mocks.setup(mockAddTags)

			h, _ := handler.NewRecipeHandler(
				usecasemock.NewMockCreateRecipeUseCase(ctrl),
				usecasemock.NewMockGetRecipeUseCase(ctrl),
				usecasemock.NewMockGetUserRecipesUseCase(ctrl),
				usecasemock.NewMockUpdateRecipeUseCase(ctrl),
				usecasemock.NewMockDeleteRecipeUseCase(ctrl),
				usecasemock.NewMockSearchRecipesUseCase(ctrl),
				mockAddTags,
				usecasemock.NewMockRemoveTagsUseCase(ctrl),
				usecasemock.NewMockCreateRecipeImageUseCase(ctrl),
				usecasemock.NewMockCreateTagUseCase(ctrl),
				usecasemock.NewMockGetAllTagsUseCase(ctrl),
				usecasemock.NewMockDeleteTagUseCase(ctrl),
			)

			var req *http.Request
			if tt.name == "Invalid JSON" {
				req = httptest.NewRequest(http.MethodPost, "/api/recipes/"+tt.args.recipeID+"/tags", bytes.NewReader([]byte("invalid")))
			} else {
				body, _ := json.Marshal(tt.args.body)
				req = httptest.NewRequest(http.MethodPost, "/api/recipes/"+tt.args.recipeID+"/tags", bytes.NewReader(body))
			}
			req.Header.Set("Content-Type", "application/json")

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.args.recipeID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			if tt.args.userID != 0 {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.args.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			h.AddTags(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestRecipeHandler_RemoveTags(t *testing.T) {
	type args struct {
		recipeID string
		body     map[string]any
		userID   int64
	}
	type mocks struct {
		setup func(m *usecasemock.MockRemoveTagsUseCase)
	}
	tests := []struct {
		name       string
		args       args
		mocks      mocks
		wantStatus int
	}{
		{
			name: "Success",
			args: args{recipeID: "100", body: map[string]any{"tag_ids": []int64{10}}, userID: 1},
			mocks: mocks{
				setup: func(m *usecasemock.MockRemoveTagsUseCase) {
					m.EXPECT().Execute(gomock.Any(), int64(100), int64(1), []int64{10}).
						Return(nil)
				},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "Invalid ID",
			args:       args{recipeID: "abc", userID: 1},
			mocks:      mocks{setup: func(m *usecasemock.MockRemoveTagsUseCase) {}},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Invalid JSON",
			args:       args{recipeID: "100", userID: 1, body: nil},
			mocks:      mocks{setup: func(m *usecasemock.MockRemoveTagsUseCase) {}},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Error",
			args: args{recipeID: "100", body: map[string]any{"tag_ids": []int64{10}}, userID: 1},
			mocks: mocks{
				setup: func(m *usecasemock.MockRemoveTagsUseCase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(errors.New("failed"))
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRemoveTags := usecasemock.NewMockRemoveTagsUseCase(ctrl)
			tt.mocks.setup(mockRemoveTags)

			h, _ := handler.NewRecipeHandler(
				usecasemock.NewMockCreateRecipeUseCase(ctrl),
				usecasemock.NewMockGetRecipeUseCase(ctrl),
				usecasemock.NewMockGetUserRecipesUseCase(ctrl),
				usecasemock.NewMockUpdateRecipeUseCase(ctrl),
				usecasemock.NewMockDeleteRecipeUseCase(ctrl),
				usecasemock.NewMockSearchRecipesUseCase(ctrl),
				usecasemock.NewMockAddTagsUseCase(ctrl),
				mockRemoveTags,
				usecasemock.NewMockCreateRecipeImageUseCase(ctrl),
				usecasemock.NewMockCreateTagUseCase(ctrl),
				usecasemock.NewMockGetAllTagsUseCase(ctrl),
				usecasemock.NewMockDeleteTagUseCase(ctrl),
			)

			var req *http.Request
			if tt.name == "Invalid JSON" {
				req = httptest.NewRequest(http.MethodDelete, "/api/recipes/"+tt.args.recipeID+"/tags", bytes.NewReader([]byte("invalid")))
			} else {
				body, _ := json.Marshal(tt.args.body)
				req = httptest.NewRequest(http.MethodDelete, "/api/recipes/"+tt.args.recipeID+"/tags", bytes.NewReader(body))
			}
			req.Header.Set("Content-Type", "application/json")

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.args.recipeID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			if tt.args.userID != 0 {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.args.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			h.RemoveTags(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestRecipeHandler_AddImage(t *testing.T) {
	type args struct {
		recipeID string
		body     map[string]any
		userID   int64
	}
	type mocks struct {
		setup func(m *usecasemock.MockCreateRecipeImageUseCase)
	}
	tests := []struct {
		name       string
		args       args
		mocks      mocks
		wantStatus int
	}{
		{
			name: "Success",
			args: args{
				recipeID: "100",
				body: map[string]any{
					"filename":     "test.jpg",
					"content_type": "image/jpeg",
					"size":         int64(10240),
				},
				userID: 1,
			},
			mocks: mocks{
				setup: func(m *usecasemock.MockCreateRecipeImageUseCase) {
					input := recipe_image.CreateRecipeImageInput{
						RecipeID:    100,
						UserID:      1,
						Filename:    "test.jpg",
						ContentType: "image/jpeg",
						Size:        10240,
					}
					m.EXPECT().Execute(gomock.Any(), input).
						Return(&model.PublicRecipeImage{
							RecipeImage: model.RecipeImage{ID: 1, ImagePath: "test.jpg"},
							ImageURL:    url.URL{Path: "http://localhost/test.jpg"},
						}, "upload-url", nil)
				},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Invalid ID",
			args: args{recipeID: "abc", userID: 1},
			mocks: mocks{
				setup: func(m *usecasemock.MockCreateRecipeImageUseCase) {},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid JSON",
			args: args{
				recipeID: "100",
				userID:   1,
				body:     nil,
			},
			mocks: mocks{
				setup: func(m *usecasemock.MockCreateRecipeImageUseCase) {},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "UseCase Error",
			args: args{
				recipeID: "100",
				body: map[string]any{
					"filename":     "test.jpg",
					"content_type": "image/jpeg",
					"size":         int64(10240),
				},
				userID: 1,
			},
			mocks: mocks{
				setup: func(m *usecasemock.MockCreateRecipeImageUseCase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any()).
						Return(nil, "", errors.New("failed"))
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCreateImage := usecasemock.NewMockCreateRecipeImageUseCase(ctrl)
			tt.mocks.setup(mockCreateImage)

			h, _ := handler.NewRecipeHandler(
				usecasemock.NewMockCreateRecipeUseCase(ctrl),
				usecasemock.NewMockGetRecipeUseCase(ctrl),
				usecasemock.NewMockGetUserRecipesUseCase(ctrl),
				usecasemock.NewMockUpdateRecipeUseCase(ctrl),
				usecasemock.NewMockDeleteRecipeUseCase(ctrl),
				usecasemock.NewMockSearchRecipesUseCase(ctrl),
				usecasemock.NewMockAddTagsUseCase(ctrl),
				usecasemock.NewMockRemoveTagsUseCase(ctrl),
				mockCreateImage,
				usecasemock.NewMockCreateTagUseCase(ctrl),
				usecasemock.NewMockGetAllTagsUseCase(ctrl),
				usecasemock.NewMockDeleteTagUseCase(ctrl),
			)

			var req *http.Request
			if tt.name == "Invalid JSON" {
				req = httptest.NewRequest(http.MethodPost, "/api/recipes/"+tt.args.recipeID+"/images", bytes.NewReader([]byte("invalid")))
			} else {
				body, _ := json.Marshal(tt.args.body)
				req = httptest.NewRequest(http.MethodPost, "/api/recipes/"+tt.args.recipeID+"/images", bytes.NewReader(body))
			}
			req.Header.Set("Content-Type", "application/json")

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.args.recipeID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			if tt.args.userID != 0 {
				ctx := context.WithValue(req.Context(), middleware.UserIDKey, tt.args.userID)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			h.AddImage(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

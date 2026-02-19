package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/server/handler"
	usecasemock "github.com/seka/reci-pin/backend/internal/usecase/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRecipeHandler_CreateTag(t *testing.T) {
	type args struct {
		body   map[string]string
		userID int64
	}
	type mocks struct {
		setup func(m *usecasemock.MockCreateTagUseCase)
	}
	tests := []struct {
		name       string
		args       args
		mocks      mocks
		wantStatus int
	}{
		{
			name: "Success",
			args: args{body: map[string]string{"name": "Vegan"}, userID: 1},
			mocks: mocks{
				setup: func(m *usecasemock.MockCreateTagUseCase) {
					m.EXPECT().Execute(gomock.Any(), "Vegan").
						Return(&model.Tag{ID: 10, Name: "Vegan"}, nil)
				},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Invalid JSON",
			args: args{body: nil, userID: 1},
			mocks: mocks{
				setup: func(m *usecasemock.MockCreateTagUseCase) {},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Error",
			args: args{body: map[string]string{"name": "Vegan"}, userID: 1},
			mocks: mocks{
				setup: func(m *usecasemock.MockCreateTagUseCase) {
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

			mockCreateTag := usecasemock.NewMockCreateTagUseCase(ctrl)
			tt.mocks.setup(mockCreateTag)

			h := handler.NewRecipeHandler(
				usecasemock.NewMockCreateRecipeUseCase(ctrl),
				usecasemock.NewMockGetRecipeUseCase(ctrl),
				usecasemock.NewMockGetUserRecipesUseCase(ctrl),
				usecasemock.NewMockUpdateRecipeUseCase(ctrl),
				usecasemock.NewMockDeleteRecipeUseCase(ctrl),
				usecasemock.NewMockSearchRecipesUseCase(ctrl),
				usecasemock.NewMockAddTagsUseCase(ctrl),
				usecasemock.NewMockRemoveTagsUseCase(ctrl),
				usecasemock.NewMockCreateRecipeImageUseCase(ctrl),
				mockCreateTag,
				usecasemock.NewMockGetAllTagsUseCase(ctrl),
				usecasemock.NewMockDeleteTagUseCase(ctrl),
				"",
			)

			var req *http.Request
			if tt.name == "Invalid JSON" {
				req = httptest.NewRequest(http.MethodPost, "/api/tags", bytes.NewReader([]byte("invalid")))
			} else {
				body, _ := json.Marshal(tt.args.body)
				req = httptest.NewRequest(http.MethodPost, "/api/tags", bytes.NewReader(body))
			}
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			h.CreateTag(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestRecipeHandler_GetAllTags(t *testing.T) {
	type mocks struct {
		setup func(m *usecasemock.MockGetAllTagsUseCase)
	}
	tests := []struct {
		name       string
		mocks      mocks
		wantStatus int
	}{
		{
			name: "Success",
			mocks: mocks{
				setup: func(m *usecasemock.MockGetAllTagsUseCase) {
					m.EXPECT().Execute(gomock.Any()).
						Return([]model.Tag{{ID: 1, Name: "Vegan"}}, nil)
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Error",
			mocks: mocks{
				setup: func(m *usecasemock.MockGetAllTagsUseCase) {
					m.EXPECT().Execute(gomock.Any()).
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

			mockGetAll := usecasemock.NewMockGetAllTagsUseCase(ctrl)
			tt.mocks.setup(mockGetAll)

			h := handler.NewRecipeHandler(
				usecasemock.NewMockCreateRecipeUseCase(ctrl),
				usecasemock.NewMockGetRecipeUseCase(ctrl),
				usecasemock.NewMockGetUserRecipesUseCase(ctrl),
				usecasemock.NewMockUpdateRecipeUseCase(ctrl),
				usecasemock.NewMockDeleteRecipeUseCase(ctrl),
				usecasemock.NewMockSearchRecipesUseCase(ctrl),
				usecasemock.NewMockAddTagsUseCase(ctrl),
				usecasemock.NewMockRemoveTagsUseCase(ctrl),
				usecasemock.NewMockCreateRecipeImageUseCase(ctrl),
				usecasemock.NewMockCreateTagUseCase(ctrl),
				mockGetAll,
				usecasemock.NewMockDeleteTagUseCase(ctrl),
				"",
			)

			req := httptest.NewRequest(http.MethodGet, "/api/tags", nil)
			w := httptest.NewRecorder()

			h.GetAllTags(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestRecipeHandler_DeleteTag(t *testing.T) {
	type args struct {
		tagID string
	}
	type mocks struct {
		setup func(m *usecasemock.MockDeleteTagUseCase)
	}
	tests := []struct {
		name       string
		args       args
		mocks      mocks
		wantStatus int
	}{
		{
			name: "Success",
			args: args{tagID: "10"},
			mocks: mocks{
				setup: func(m *usecasemock.MockDeleteTagUseCase) {
					m.EXPECT().Execute(gomock.Any(), int64(10)).Return(nil)
				},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "Invalid ID",
			args: args{tagID: "abc"},
			mocks: mocks{
				setup: func(m *usecasemock.MockDeleteTagUseCase) {},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Error",
			args: args{tagID: "10"},
			mocks: mocks{
				setup: func(m *usecasemock.MockDeleteTagUseCase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any()).
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

			mockDelete := usecasemock.NewMockDeleteTagUseCase(ctrl)
			tt.mocks.setup(mockDelete)

			h := handler.NewRecipeHandler(
				usecasemock.NewMockCreateRecipeUseCase(ctrl),
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
				mockDelete,
				"",
			)

			req := httptest.NewRequest(http.MethodDelete, "/api/tags/"+tt.args.tagID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.args.tagID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			h.DeleteTag(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

package recipe_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository/mock"
	mock_storage "github.com/seka/reci-pin/backend/internal/domain/storage/mock"
	"github.com/seka/reci-pin/backend/internal/usecase/recipe"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUserRecipesUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		userID  int64
		setup   func(*mock.MockRecipeRepository, *mock.MockRecipeImageRepository, *mock_storage.MockClient)
		wantErr bool
		wantLen int
	}{
		{
			name:   "正常系_複数レシピ取得",
			userID: 1,
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository, ms *mock_storage.MockClient) {
				mr.EXPECT().
					GetByUserID(gomock.Any(), int64(1)).
					Return([]model.Recipe{
						{ID: 1, UserID: 1, Name: "Recipe 1"},
						{ID: 2, UserID: 1, Name: "Recipe 2"},
					}, nil)

				// 各レシピのタグ取得
				mr.EXPECT().
					GetTags(gomock.Any(), int64(1)).
					Return([]model.Tag{{ID: 1, Name: "Tag1"}}, nil)
				mr.EXPECT().
					GetTags(gomock.Any(), int64(2)).
					Return([]model.Tag{}, nil)

				// 各レシピの画像取得
				mi.EXPECT().
					GetByRecipeID(gomock.Any(), int64(1)).
					Return([]model.RecipeImage{{ID: 1, RecipeID: 1, ImagePath: "1.jpg"}}, nil)
				mi.EXPECT().
					GetByRecipeID(gomock.Any(), int64(2)).
					Return([]model.RecipeImage{}, nil)

				ms.EXPECT().GetPublicURL().AnyTimes().Return(&url.URL{Scheme: "http", Host: "localhost"})
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name:   "正常系_レシピなし",
			userID: 999,
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository, ms *mock_storage.MockClient) {
				mr.EXPECT().
					GetByUserID(gomock.Any(), int64(999)).
					Return([]model.Recipe{}, nil)
			},
			wantErr: false,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRecipeRepo := mock.NewMockRecipeRepository(ctrl)
			mockImageRepo := mock.NewMockRecipeImageRepository(ctrl)
			mockStorage := mock_storage.NewMockClient(ctrl)
			tt.setup(mockRecipeRepo, mockImageRepo, mockStorage)

			uc := recipe.NewGetUserRecipesUseCase(mockRecipeRepo, mockImageRepo, mockStorage)
			result, err := uc.Execute(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.wantLen)
			}
		})
	}
}

func TestSearchRecipesUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		input   recipe.SearchRecipesInput
		setup   func(*mock.MockRecipeRepository, *mock.MockRecipeImageRepository, *mock.MockRecipeSearchRepository, *mock_storage.MockClient)
		wantErr bool
		wantLen int
	}{
		{
			name: "正常系_クエリ検索",
			input: recipe.SearchRecipesInput{
				UserID: 1,
				Query:  "pasta",
				TagIDs: []int64{},
			},
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository, ms *mock.MockRecipeSearchRepository, storage *mock_storage.MockClient) {
				ms.EXPECT().
					Search(gomock.Any(), gomock.Any()).
					Return([]int64{1}, int64(1), nil)
				mr.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{
						ID: 1, UserID: 1, Name: "Pasta Recipe",
					}, nil)
				mr.EXPECT().
					GetTags(gomock.Any(), int64(1)).
					Return([]model.Tag{}, nil)
				mi.EXPECT().
					GetByRecipeID(gomock.Any(), int64(1)).
					Return([]model.RecipeImage{}, nil)
				storage.EXPECT().GetPublicURL().AnyTimes().Return(&url.URL{Scheme: "http", Host: "localhost"})
			},
			wantErr: false,
			wantLen: 1,
		},
		{
			name: "正常系_タグ検索",
			input: recipe.SearchRecipesInput{
				UserID: 1,
				Query:  "",
				TagIDs: []int64{1, 2},
			},
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository, ms *mock.MockRecipeSearchRepository, storage *mock_storage.MockClient) {
				ms.EXPECT().
					Search(gomock.Any(), gomock.Any()).
					Return([]int64{1, 2}, int64(2), nil)
				mr.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{ID: 1, UserID: 1, Name: "Recipe 1"}, nil)
				mr.EXPECT().
					GetByID(gomock.Any(), int64(2)).
					Return(&model.Recipe{ID: 2, UserID: 1, Name: "Recipe 2"}, nil)
				mr.EXPECT().
					GetTags(gomock.Any(), int64(1)).
					Return([]model.Tag{{ID: 1, Name: "Tag1"}}, nil)
				mr.EXPECT().
					GetTags(gomock.Any(), int64(2)).
					Return([]model.Tag{{ID: 2, Name: "Tag2"}}, nil)
				mi.EXPECT().
					GetByRecipeID(gomock.Any(), int64(1)).
					Return([]model.RecipeImage{}, nil)
				mi.EXPECT().
					GetByRecipeID(gomock.Any(), int64(2)).
					Return([]model.RecipeImage{}, nil)
				storage.EXPECT().GetPublicURL().AnyTimes().Return(&url.URL{Scheme: "http", Host: "localhost"})
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name: "正常系_検索結果なし",
			input: recipe.SearchRecipesInput{
				UserID: 1,
				Query:  "nonexistent",
				TagIDs: []int64{},
			},
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository, ms *mock.MockRecipeSearchRepository, storage *mock_storage.MockClient) {
				ms.EXPECT().
					Search(gomock.Any(), gomock.Any()).
					Return([]int64{}, int64(0), nil)
			},
			wantErr: false,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRecipeRepo := mock.NewMockRecipeRepository(ctrl)
			mockImageRepo := mock.NewMockRecipeImageRepository(ctrl)
			mockSearchRepo := mock.NewMockRecipeSearchRepository(ctrl)
			mockStorage := mock_storage.NewMockClient(ctrl)
			tt.setup(mockRecipeRepo, mockImageRepo, mockSearchRepo, mockStorage)

			uc := recipe.NewSearchRecipesUseCase(mockRecipeRepo, mockImageRepo, mockSearchRepo, mockStorage)
			result, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.wantLen)
			}
		})
	}
}

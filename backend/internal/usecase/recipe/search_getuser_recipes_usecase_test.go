package recipe_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	mockRepo "github.com/seka/reci-pin/backend/internal/domain/repository/mock"
	mockSearcher "github.com/seka/reci-pin/backend/internal/domain/searcher/mock"
	mockStorage "github.com/seka/reci-pin/backend/internal/domain/storage/mock"
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
		setup   func(*mockRepo.MockRecipeRepository, *mockRepo.MockRecipeImageRepository, *mockStorage.MockClient)
		wantErr bool
		wantLen int
	}{
		{
			name:   "正常系_複数レシピ取得",
			userID: 1,
			setup: func(mr *mockRepo.MockRecipeRepository, mi *mockRepo.MockRecipeImageRepository, ms *mockStorage.MockClient) {
				mr.EXPECT().
					GetByUserID(gomock.Any(), int64(1)).
					Return([]model.Recipe{
						{ID: 1, UserID: 1, Name: "Recipe 1"},
						{ID: 2, UserID: 1, Name: "Recipe 2"},
					}, nil)

				// バッチ取得
				mr.EXPECT().
					BulkGetTags(gomock.Any(), gomock.Any()).
					Return(map[int64][]model.Tag{
						1: {{ID: 1, Name: "Tag1"}},
					}, nil)

				mi.EXPECT().
					BulkGetByRecipeIDs(gomock.Any(), gomock.Any()).
					Return(map[int64][]model.RecipeImage{
						1: {{ID: 1, RecipeID: 1, ImagePath: "1.jpg"}},
					}, nil)

				ms.EXPECT().GetPublicURL().AnyTimes().Return(&url.URL{Scheme: "http", Host: "localhost"})
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name:   "正常系_レシピなし",
			userID: 999,
			setup: func(mr *mockRepo.MockRecipeRepository, mi *mockRepo.MockRecipeImageRepository, ms *mockStorage.MockClient) {
				mr.EXPECT().
					GetByUserID(gomock.Any(), int64(999)).
					Return([]model.Recipe{}, nil)
				// 空のリストの場合もバッチ呼び出しが行われる実装になっているため、それらを期待する
				mr.EXPECT().BulkGetTags(gomock.Any(), []int64{}).Return(map[int64][]model.Tag{}, nil)
				mi.EXPECT().BulkGetByRecipeIDs(gomock.Any(), []int64{}).Return(map[int64][]model.RecipeImage{}, nil)
				ms.EXPECT().GetPublicURL().Return(&url.URL{Scheme: "http", Host: "localhost"})
			},
			wantErr: false,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRecipeRepo := mockRepo.NewMockRecipeRepository(ctrl)
			mockImageRepo := mockRepo.NewMockRecipeImageRepository(ctrl)
			mockStorage := mockStorage.NewMockClient(ctrl)
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
		setup   func(*mockRepo.MockRecipeRepository, *mockRepo.MockRecipeImageRepository, *mockSearcher.MockRecipeSearcher, *mockStorage.MockClient)
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
			setup: func(mr *mockRepo.MockRecipeRepository, mi *mockRepo.MockRecipeImageRepository, ms *mockSearcher.MockRecipeSearcher, storage *mockStorage.MockClient) {
				ms.EXPECT().
					Search(gomock.Any(), gomock.Any()).
					Return([]int64{1}, int64(1), nil)
				mr.EXPECT().
					BulkGetByIDs(gomock.Any(), []int64{1}).
					Return([]model.Recipe{
						{ID: 1, UserID: 1, Name: "Pasta Recipe"},
					}, nil)
				mr.EXPECT().
					BulkGetTags(gomock.Any(), []int64{1}).
					Return(map[int64][]model.Tag{}, nil)
				mi.EXPECT().
					BulkGetByRecipeIDs(gomock.Any(), []int64{1}).
					Return(map[int64][]model.RecipeImage{}, nil)
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
			setup: func(mr *mockRepo.MockRecipeRepository, mi *mockRepo.MockRecipeImageRepository, ms *mockSearcher.MockRecipeSearcher, storage *mockStorage.MockClient) {
				ms.EXPECT().
					Search(gomock.Any(), gomock.Any()).
					Return([]int64{1, 2}, int64(2), nil)
				mr.EXPECT().
					BulkGetByIDs(gomock.Any(), []int64{1, 2}).
					Return([]model.Recipe{
						{ID: 1, UserID: 1, Name: "Recipe 1"},
						{ID: 2, UserID: 1, Name: "Recipe 2"},
					}, nil)
				mr.EXPECT().
					BulkGetTags(gomock.Any(), []int64{1, 2}).
					Return(map[int64][]model.Tag{
						1: {{ID: 1, Name: "Tag1"}},
						2: {{ID: 2, Name: "Tag2"}},
					}, nil)
				mi.EXPECT().
					BulkGetByRecipeIDs(gomock.Any(), []int64{1, 2}).
					Return(map[int64][]model.RecipeImage{}, nil)
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
			setup: func(mr *mockRepo.MockRecipeRepository, mi *mockRepo.MockRecipeImageRepository, ms *mockSearcher.MockRecipeSearcher, storage *mockStorage.MockClient) {
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
			mockRecipeRepo := mockRepo.NewMockRecipeRepository(ctrl)
			mockImageRepo := mockRepo.NewMockRecipeImageRepository(ctrl)
			mockSearchRepo := mockSearcher.NewMockRecipeSearcher(ctrl)
			mockStorage := mockStorage.NewMockClient(ctrl)
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

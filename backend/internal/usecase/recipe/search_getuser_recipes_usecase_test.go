package recipe_test

import (
	"context"
	"testing"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/notification/mock"
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
		setup   func(*mock.MockRecipeRepository, *mock.MockRecipeImageRepository)
		wantErr bool
		wantLen int
	}{
		{
			name:   "正常系_複数レシピ取得",
			userID: 1,
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository) {
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
					Return([]model.RecipeImage{{ID: 1, RecipeID: 1}}, nil)
				mi.EXPECT().
					GetByRecipeID(gomock.Any(), int64(2)).
					Return([]model.RecipeImage{}, nil)
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name:   "正常系_レシピなし",
			userID: 999,
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository) {
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
			tt.setup(mockRecipeRepo, mockImageRepo)

			uc := recipe.NewGetUserRecipesUseCase(mockRecipeRepo, mockImageRepo)
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
		setup   func(*mock.MockRecipeRepository, *mock.MockRecipeImageRepository)
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
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository) {
				mr.EXPECT().
					Search(gomock.Any(), int64(1), "pasta", []int64{}).
					Return([]model.Recipe{
						{ID: 1, UserID: 1, Name: "Pasta Recipe"},
					}, nil)
				mr.EXPECT().
					GetTags(gomock.Any(), int64(1)).
					Return([]model.Tag{}, nil)
				mi.EXPECT().
					GetByRecipeID(gomock.Any(), int64(1)).
					Return([]model.RecipeImage{}, nil)
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
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository) {
				mr.EXPECT().
					Search(gomock.Any(), int64(1), "", []int64{1, 2}).
					Return([]model.Recipe{
						{ID: 1, UserID: 1, Name: "Recipe 1"},
						{ID: 2, UserID: 1, Name: "Recipe 2"},
					}, nil)
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
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository) {
				mr.EXPECT().
					Search(gomock.Any(), int64(1), "nonexistent", []int64{}).
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
			tt.setup(mockRecipeRepo, mockImageRepo)

			uc := recipe.NewSearchRecipesUseCase(mockRecipeRepo, mockImageRepo)
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

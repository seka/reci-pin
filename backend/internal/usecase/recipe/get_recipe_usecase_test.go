package recipe_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository/mock"
	"github.com/seka/reci-pin/backend/internal/usecase/recipe"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetRecipeUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		id      int64
		userID  int64
		setup   func(*mock.MockRecipeRepository, *mock.MockRecipeImageRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name:   "正常系_レシピ取得成功",
			id:     1,
			userID: 1,
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository) {
				mr.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{
						ID:     1,
						UserID: 1,
						Name:   "Test Recipe",
						URL:    "https://example.com",
						Memo:   "Test memo",
					}, nil)
				mr.EXPECT().
					GetTags(gomock.Any(), int64(1)).
					Return([]model.Tag{
						{ID: 1, Name: "Tag1"},
					}, nil)
				mi.EXPECT().
					GetByRecipeID(gomock.Any(), int64(1)).
					Return([]model.RecipeImage{
						{ID: 1, RecipeID: 1, ImagePath: "/images/1.jpg"},
					}, nil)
			},
			wantErr: false,
		},
		{
			name:   "異常系_レシピ不在",
			id:     999,
			userID: 1,
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository) {
				mr.EXPECT().
					GetByID(gomock.Any(), int64(999)).
					Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "failed to get recipe",
		},
		{
			name:   "異常系_権限エラー",
			id:     1,
			userID: 2, // 異なるユーザー
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository) {
				mr.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{
						ID:     1,
						UserID: 1, // オーナーは別ユーザー
						Name:   "Test Recipe",
					}, nil)
			},
			wantErr: true,
			errMsg:  "unauthorized access to recipe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRecipeRepo := mock.NewMockRecipeRepository(ctrl)
			mockImageRepo := mock.NewMockRecipeImageRepository(ctrl)
			tt.setup(mockRecipeRepo, mockImageRepo)

			uc := recipe.NewGetRecipeUseCase(mockRecipeRepo, mockImageRepo)
			result, err := uc.Execute(context.Background(), tt.id, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.id, result.ID)
				assert.Equal(t, tt.userID, result.UserID)
				assert.NotEmpty(t, result.Tags)
				assert.NotEmpty(t, result.Images)
			}
		})
	}
}

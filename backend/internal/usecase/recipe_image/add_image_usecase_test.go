package recipe_image_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository/mock"
	"github.com/seka/reci-pin/backend/internal/usecase/recipe_image"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAddImageUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name      string
		recipeID  int64
		userID    int64
		imagePath string
		setup     func(*mock.MockRecipeRepository, *mock.MockRecipeImageRepository)
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "正常系_画像追加成功",
			recipeID:  1,
			userID:    1,
			imagePath: "/images/recipe1.jpg",
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository) {
				mr.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{ID: 1, UserID: 1}, nil)
				mi.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, img *model.RecipeImage) error {
						img.ID = 1
						return nil
					})
			},
			wantErr: false,
		},
		{
			name:      "異常系_レシピ不在",
			recipeID:  999,
			userID:    1,
			imagePath: "/images/recipe.jpg",
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository) {
				mr.EXPECT().
					GetByID(gomock.Any(), int64(999)).
					Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "failed to get recipe",
		},
		{
			name:      "異常系_権限エラー",
			recipeID:  1,
			userID:    2,
			imagePath: "/images/recipe.jpg",
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository) {
				mr.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{ID: 1, UserID: 1}, nil)
			},
			wantErr: true,
			errMsg:  "unauthorized access to recipe",
		},
		{
			name:      "異常系_画像作成失敗",
			recipeID:  1,
			userID:    1,
			imagePath: "/images/recipe.jpg",
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository) {
				mr.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{ID: 1, UserID: 1}, nil)
				mi.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "failed to add image to recipe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRecipeRepo := mock.NewMockRecipeRepository(ctrl)
			mockImageRepo := mock.NewMockRecipeImageRepository(ctrl)
			tt.setup(mockRecipeRepo, mockImageRepo)

			uc := recipe_image.NewAddImageUseCase(mockRecipeRepo, mockImageRepo)
			result, err := uc.Execute(context.Background(), tt.recipeID, tt.userID, tt.imagePath)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.recipeID, result.RecipeID)
				assert.Equal(t, tt.imagePath, result.ImagePath)
			}
		})
	}
}

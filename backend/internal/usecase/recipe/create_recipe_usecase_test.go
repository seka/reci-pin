package recipe_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/notification/mock"
	"github.com/seka/reci-pin/backend/internal/usecase/recipe"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateRecipeUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		input   recipe.CreateRecipeInput
		setup   func(*mock.MockRecipeRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "正常系_レシピ作成成功（タグなし）",
			input: recipe.CreateRecipeInput{
				UserID: 1,
				Name:   "Test Recipe",
				URL:    "https://example.com/recipe",
				Memo:   "Test memo",
				TagIDs: []int64{},
			},
			setup: func(m *mock.MockRecipeRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, r *model.Recipe) error {
						r.ID = 1
						return nil
					})
			},
			wantErr: false,
		},
		{
			name: "正常系_レシピ作成成功（タグあり）",
			input: recipe.CreateRecipeInput{
				UserID: 1,
				Name:   "Test Recipe with Tags",
				URL:    "https://example.com/recipe",
				Memo:   "Test memo",
				TagIDs: []int64{1, 2},
			},
			setup: func(m *mock.MockRecipeRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, r *model.Recipe) error {
						r.ID = 1
						return nil
					})
				m.EXPECT().
					AddTags(gomock.Any(), int64(1), []int64{1, 2}).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "異常系_レシピ作成失敗",
			input: recipe.CreateRecipeInput{
				UserID: 1,
				Name:   "Test Recipe",
				URL:    "https://example.com/recipe",
				Memo:   "Test memo",
				TagIDs: []int64{},
			},
			setup: func(m *mock.MockRecipeRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "failed to create recipe",
		},
		{
			name: "異常系_タグ追加失敗",
			input: recipe.CreateRecipeInput{
				UserID: 1,
				Name:   "Test Recipe",
				URL:    "https://example.com/recipe",
				Memo:   "Test memo",
				TagIDs: []int64{1},
			},
			setup: func(m *mock.MockRecipeRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, r *model.Recipe) error {
						r.ID = 1
						return nil
					})
				m.EXPECT().
					AddTags(gomock.Any(), int64(1), []int64{1}).
					Return(errors.New("tag add error"))
			},
			wantErr: true,
			errMsg:  "failed to add tags to recipe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.NewMockRecipeRepository(ctrl)
			tt.setup(mockRepo)

			uc := recipe.NewCreateRecipeUseCase(mockRepo)
			result, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.input.Name, result.Name)
				assert.Equal(t, tt.input.UserID, result.UserID)
			}
		})
	}
}

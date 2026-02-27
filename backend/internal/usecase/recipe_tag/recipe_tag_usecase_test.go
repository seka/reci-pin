package recipe_tag_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository/mock"
	"github.com/seka/reci-pin/backend/internal/usecase/recipe_tag"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAddTagsUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		recipeID int64
		userID   int64
		tagIDs   []int64
		setup    func(*mock.MockRecipeRepository, *mock.MockTransactionManager)
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "正常系_タグ追加成功",
			recipeID: 1,
			userID:   1,
			tagIDs:   []int64{1, 2},
			setup: func(m *mock.MockRecipeRepository, mtm *mock.MockTransactionManager) {
				mtm.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					},
				)
				m.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{ID: 1, UserID: 1}, nil)
				m.EXPECT().
					AddTags(gomock.Any(), int64(1), []int64{1, 2}).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "異常系_レシピ不在",
			recipeID: 999,
			userID:   1,
			tagIDs:   []int64{1},
			setup: func(m *mock.MockRecipeRepository, mtm *mock.MockTransactionManager) {
				mtm.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					},
				)
				m.EXPECT().
					GetByID(gomock.Any(), int64(999)).
					Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "failed to get recipe",
		},
		{
			name:     "異常系_権限エラー",
			recipeID: 1,
			userID:   2,
			tagIDs:   []int64{1},
			setup: func(m *mock.MockRecipeRepository, mtm *mock.MockTransactionManager) {
				mtm.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					},
				)
				m.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{ID: 1, UserID: 1}, nil)
			},
			wantErr: true,
			errMsg:  "unauthorized access to recipe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mRepo := mock.NewMockRecipeRepository(ctrl)
			mTxMgr := mock.NewMockTransactionManager(ctrl)
			tt.setup(mRepo, mTxMgr)

			uc := recipe_tag.NewAddTagsUseCase(mRepo, mTxMgr)
			err := uc.Execute(context.Background(), tt.recipeID, tt.userID, tt.tagIDs)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRemoveTagsUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		recipeID int64
		userID   int64
		tagIDs   []int64
		setup    func(*mock.MockRecipeRepository, *mock.MockTransactionManager)
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "正常系_タグ削除成功",
			recipeID: 1,
			userID:   1,
			tagIDs:   []int64{1, 2},
			setup: func(m *mock.MockRecipeRepository, mtm *mock.MockTransactionManager) {
				mtm.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					},
				)
				m.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{ID: 1, UserID: 1}, nil)
				m.EXPECT().
					RemoveTags(gomock.Any(), int64(1), []int64{1, 2}).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "異常系_レシピ不在",
			recipeID: 999,
			userID:   1,
			tagIDs:   []int64{1},
			setup: func(m *mock.MockRecipeRepository, mtm *mock.MockTransactionManager) {
				mtm.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					},
				)
				m.EXPECT().
					GetByID(gomock.Any(), int64(999)).
					Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "failed to get recipe",
		},
		{
			name:     "異常系_権限エラー",
			recipeID: 1,
			userID:   2,
			tagIDs:   []int64{1},
			setup: func(m *mock.MockRecipeRepository, mtm *mock.MockTransactionManager) {
				mtm.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					},
				)
				m.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{ID: 1, UserID: 1}, nil)
			},
			wantErr: true,
			errMsg:  "unauthorized access to recipe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mRepo := mock.NewMockRecipeRepository(ctrl)
			mTxMgr := mock.NewMockTransactionManager(ctrl)
			tt.setup(mRepo, mTxMgr)

			uc := recipe_tag.NewRemoveTagsUseCase(mRepo, mTxMgr)
			err := uc.Execute(context.Background(), tt.recipeID, tt.userID, tt.tagIDs)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

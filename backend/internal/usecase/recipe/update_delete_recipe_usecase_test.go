package recipe_test

import (
	"context"
	"errors"
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

func TestUpdateRecipeUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		input   recipe.UpdateRecipeInput
		setup   func(*mockRepo.MockRecipeRepository, *mockRepo.MockRecipeImageRepository, *mockSearcher.MockRecipeSearcher, *mockStorage.MockClient)
		wantErr bool
		errMsg  string
	}{
		{
			name: "正常系_レシピ更新成功",
			input: recipe.UpdateRecipeInput{
				ID:     1,
				UserID: 1,
				Name:   "Updated Recipe",
				URL:    "https://example.com/updated",
				Memo:   "Updated memo",
			},
			setup: func(m *mockRepo.MockRecipeRepository, mi *mockRepo.MockRecipeImageRepository, ms *mockSearcher.MockRecipeSearcher, s *mockStorage.MockClient) {
				m.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{
						ID:     1,
						UserID: 1,
						Name:   "Old Recipe",
					}, nil)
				m.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)

				// GetTags for indexing and response
				m.EXPECT().
					GetTags(gomock.Any(), int64(1)).
					Return([]model.Tag{}, nil)

				// GetImages for response
				mi.EXPECT().
					GetByRecipeID(gomock.Any(), int64(1)).
					Return([]model.RecipeImage{}, nil)

				// Storage URL
				s.EXPECT().
					GetPublicURL().
					Return(&url.URL{Scheme: "https", Host: "localhost", Path: "/storage"})

				// Index
				ms.EXPECT().
					Index(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "異常系_レシピ不在",
			input: recipe.UpdateRecipeInput{
				ID:     999,
				UserID: 1,
				Name:   "Updated Recipe",
			},
			setup: func(m *mockRepo.MockRecipeRepository, mi *mockRepo.MockRecipeImageRepository, ms *mockSearcher.MockRecipeSearcher, s *mockStorage.MockClient) {
				m.EXPECT().
					GetByID(gomock.Any(), int64(999)).
					Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "failed to get recipe",
		},
		{
			name: "異常系_権限エラー",
			input: recipe.UpdateRecipeInput{
				ID:     1,
				UserID: 2,
				Name:   "Updated Recipe",
			},
			setup: func(m *mockRepo.MockRecipeRepository, mi *mockRepo.MockRecipeImageRepository, ms *mockSearcher.MockRecipeSearcher, s *mockStorage.MockClient) {
				m.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{
						ID:     1,
						UserID: 1, // 別ユーザーが所有
					}, nil)
			},
			wantErr: true,
			errMsg:  "unauthorized access to recipe",
		},
		{
			name: "異常系_更新失敗",
			input: recipe.UpdateRecipeInput{
				ID:     1,
				UserID: 1,
				Name:   "Updated Recipe",
			},
			setup: func(m *mockRepo.MockRecipeRepository, mi *mockRepo.MockRecipeImageRepository, ms *mockSearcher.MockRecipeSearcher, s *mockStorage.MockClient) {
				m.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{
						ID:     1,
						UserID: 1,
					}, nil)
				m.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "failed to update recipe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepoInstance := mockRepo.NewMockRecipeRepository(ctrl)
			mockImageRepo := mockRepo.NewMockRecipeImageRepository(ctrl)
			mockSearchRepo := mockSearcher.NewMockRecipeSearcher(ctrl)
			mockStorage := mockStorage.NewMockClient(ctrl)
			tt.setup(mockRepoInstance, mockImageRepo, mockSearchRepo, mockStorage)

			uc := recipe.NewUpdateRecipeUseCase(mockRepoInstance, mockImageRepo, mockSearchRepo, mockStorage)
			result, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.input.Name, result.Name)
			}
		})
	}
}

func TestDeleteRecipeUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		id      int64
		userID  int64
		setup   func(*mockRepo.MockRecipeRepository, *mockSearcher.MockRecipeSearcher)
		wantErr bool
		errMsg  string
	}{
		{
			name:   "正常系_レシピ削除成功",
			id:     1,
			userID: 1,
			setup: func(m *mockRepo.MockRecipeRepository, ms *mockSearcher.MockRecipeSearcher) {
				m.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{
						ID:     1,
						UserID: 1,
					}, nil)
				m.EXPECT().
					Delete(gomock.Any(), int64(1)).
					Return(nil)

				// Delete from index
				ms.EXPECT().
					Delete(gomock.Any(), int64(1)).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "異常系_レシピ不在",
			id:     999,
			userID: 1,
			setup: func(m *mockRepo.MockRecipeRepository, ms *mockSearcher.MockRecipeSearcher) {
				m.EXPECT().
					GetByID(gomock.Any(), int64(999)).
					Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "failed to get recipe",
		},
		{
			name:   "異常系_権限エラー",
			id:     1,
			userID: 2,
			setup: func(m *mockRepo.MockRecipeRepository, ms *mockSearcher.MockRecipeSearcher) {
				m.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{
						ID:     1,
						UserID: 1,
					}, nil)
			},
			wantErr: true,
			errMsg:  "unauthorized access to recipe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepoInstance := mockRepo.NewMockRecipeRepository(ctrl)
			mockSearchRepo := mockSearcher.NewMockRecipeSearcher(ctrl)
			tt.setup(mockRepoInstance, mockSearchRepo)

			uc := recipe.NewDeleteRecipeUseCase(mockRepoInstance, mockSearchRepo)
			err := uc.Execute(context.Background(), tt.id, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

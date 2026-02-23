package recipe_image_test

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository/mock"
	mock_storage "github.com/seka/reci-pin/backend/internal/domain/storage/mock"
	"github.com/seka/reci-pin/backend/internal/usecase/recipe_image"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateRecipeImageUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		input   recipe_image.CreateRecipeImageInput
		setup   func(*mock.MockRecipeRepository, *mock.MockRecipeImageRepository, *mock_storage.MockClient)
		wantErr bool
		errMsg  string
	}{
		{
			name: "正常系_Presigned URL発行成功",
			input: recipe_image.CreateRecipeImageInput{
				RecipeID:    1,
				UserID:      1,
				Filename:    "recipe1.jpg",
				ContentType: "image/jpeg",
				Size:        1024,
			},
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository, ms *mock_storage.MockClient) {
				mr.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{ID: 1, UserID: 1}, nil)
				ms.EXPECT().
					GeneratePresignedURL(gomock.Any(), gomock.Any(), "image/jpeg", int64(1024), gomock.Any()).
					Return("https://s3.example.com/presigned-url", nil)
				mi.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, img *model.RecipeImage) error {
						img.ID = 1
						return nil
					})
				ms.EXPECT().GetPublicURL().Return(&url.URL{Scheme: "http", Host: "localhost"})
			},
			wantErr: false,
		},
		{
			name: "異常系_レシピ不在",
			input: recipe_image.CreateRecipeImageInput{
				RecipeID:    999,
				UserID:      1,
				Filename:    "recipe.jpg",
				ContentType: "image/jpeg",
				Size:        1024,
			},
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository, ms *mock_storage.MockClient) {
				mr.EXPECT().
					GetByID(gomock.Any(), int64(999)).
					Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "failed to get recipe",
		},
		{
			name: "異常系_権限エラー",
			input: recipe_image.CreateRecipeImageInput{
				RecipeID:    1,
				UserID:      2,
				Filename:    "recipe.jpg",
				ContentType: "image/jpeg",
				Size:        1024,
			},
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository, ms *mock_storage.MockClient) {
				mr.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{ID: 1, UserID: 1}, nil)
			},
			wantErr: true,
			errMsg:  "unauthorized access to recipe",
		},
		{
			name: "異常系_ファイルサイズ超過",
			input: recipe_image.CreateRecipeImageInput{
				RecipeID:    1,
				UserID:      1,
				Filename:    "recipe.jpg",
				ContentType: "image/jpeg",
				Size:        51 * 1024 * 1024, // 51MB
			},
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository, ms *mock_storage.MockClient) {
				mr.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{ID: 1, UserID: 1}, nil)
			},
			wantErr: true,
			errMsg:  "file size exceeds limit",
		},
		{
			name: "異常系_不正なContentType",
			input: recipe_image.CreateRecipeImageInput{
				RecipeID:    1,
				UserID:      1,
				Filename:    "recipe.gif",
				ContentType: "image/gif",
				Size:        1024,
			},
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository, ms *mock_storage.MockClient) {
				mr.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{ID: 1, UserID: 1}, nil)
			},
			wantErr: true,
			errMsg:  "unsupported content type",
		},
		{
			name: "異常系_不正な拡張子",
			input: recipe_image.CreateRecipeImageInput{
				RecipeID:    1,
				UserID:      1,
				Filename:    "recipe.bmp",
				ContentType: "image/jpeg",
				Size:        1024,
			},
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository, ms *mock_storage.MockClient) {
				mr.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{ID: 1, UserID: 1}, nil)
			},
			wantErr: true,
			errMsg:  "unsupported file extension",
		},
		{
			name: "異常系_PresignedURL生成失敗",
			input: recipe_image.CreateRecipeImageInput{
				RecipeID:    1,
				UserID:      1,
				Filename:    "recipe.jpg",
				ContentType: "image/jpeg",
				Size:        1024,
			},
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository, ms *mock_storage.MockClient) {
				mr.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{ID: 1, UserID: 1}, nil)
				ms.EXPECT().
					GeneratePresignedURL(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", errors.New("presign failed"))
			},
			wantErr: true,
			errMsg:  "failed to generate presigned URL",
		},
		{
			name: "異常系_DB保存失敗",
			input: recipe_image.CreateRecipeImageInput{
				RecipeID:    1,
				UserID:      1,
				Filename:    "recipe.jpg",
				ContentType: "image/jpeg",
				Size:        1024,
			},
			setup: func(mr *mock.MockRecipeRepository, mi *mock.MockRecipeImageRepository, ms *mock_storage.MockClient) {
				mr.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&model.Recipe{ID: 1, UserID: 1}, nil)
				ms.EXPECT().
					GeneratePresignedURL(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("https://s3.example.com/presigned-url", nil)
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
			mockStorageService := mock_storage.NewMockClient(ctrl)
			tt.setup(mockRecipeRepo, mockImageRepo, mockStorageService)

			uc := recipe_image.NewCreateRecipeImageUseCase(mockRecipeRepo, mockImageRepo, mockStorageService)
			image, url, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, image)
				assert.Empty(t, url)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, image)
				assert.NotEmpty(t, url)
				assert.Equal(t, tt.input.RecipeID, image.RecipeImage.RecipeID)
			}
		})
	}
}

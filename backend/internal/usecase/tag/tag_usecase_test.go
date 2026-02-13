package tag_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/notification/mock"
	"github.com/seka/reci-pin/backend/internal/usecase/tag"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateTagUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		tagName string
		setup   func(*mock.MockTagRepository)
		wantErr bool
		wantID  int64
	}{
		{
			name:    "正常系_新規タグ作成",
			tagName: "NewTag",
			setup: func(m *mock.MockTagRepository) {
				// タグが存在しないことを確認
				m.EXPECT().
					GetByName(gomock.Any(), "NewTag").
					Return(nil, errors.New("not found"))

				// タグを作成
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, tag *model.Tag) error {
						tag.ID = 1
						return nil
					})
			},
			wantErr: false,
			wantID:  1,
		},
		{
			name:    "正常系_既存タグを返す",
			tagName: "ExistingTag",
			setup: func(m *mock.MockTagRepository) {
				// 既存タグを返す
				m.EXPECT().
					GetByName(gomock.Any(), "ExistingTag").
					Return(&model.Tag{ID: 5, Name: "ExistingTag"}, nil)
			},
			wantErr: false,
			wantID:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.NewMockTagRepository(ctrl)
			tt.setup(mockRepo)

			uc := tag.NewCreateTagUseCase(mockRepo)
			result, err := uc.Execute(context.Background(), tt.tagName)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.wantID, result.ID)
				assert.Equal(t, tt.tagName, result.Name)
			}
		})
	}
}

func TestGetAllTagsUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		setup   func(*mock.MockTagRepository)
		wantErr bool
		wantLen int
	}{
		{
			name: "正常系_複数タグ取得",
			setup: func(m *mock.MockTagRepository) {
				m.EXPECT().
					GetAll(gomock.Any()).
					Return([]model.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
						{ID: 3, Name: "Tag3"},
					}, nil)
			},
			wantErr: false,
			wantLen: 3,
		},
		{
			name: "正常系_タグなし",
			setup: func(m *mock.MockTagRepository) {
				m.EXPECT().
					GetAll(gomock.Any()).
					Return([]model.Tag{}, nil)
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name: "異常系_DB エラー",
			setup: func(m *mock.MockTagRepository) {
				m.EXPECT().
					GetAll(gomock.Any()).
					Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.NewMockTagRepository(ctrl)
			tt.setup(mockRepo)

			uc := tag.NewGetAllTagsUseCase(mockRepo)
			result, err := uc.Execute(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.wantLen)
			}
		})
	}
}

func TestDeleteTagUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		tagID   int64
		setup   func(*mock.MockTagRepository)
		wantErr bool
	}{
		{
			name:  "正常系_タグ削除成功",
			tagID: 1,
			setup: func(m *mock.MockTagRepository) {
				m.EXPECT().
					Delete(gomock.Any(), int64(1)).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "異常系_削除失敗",
			tagID: 999,
			setup: func(m *mock.MockTagRepository) {
				m.EXPECT().
					Delete(gomock.Any(), int64(999)).
					Return(errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.NewMockTagRepository(ctrl)
			tt.setup(mockRepo)

			uc := tag.NewDeleteTagUseCase(mockRepo)
			err := uc.Execute(context.Background(), tt.tagID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

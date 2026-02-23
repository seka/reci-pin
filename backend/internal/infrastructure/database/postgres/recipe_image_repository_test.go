package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	mockPostgres "github.com/seka/reci-pin/backend/internal/infrastructure/database/mock"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database/postgres"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRecipeImageRepository_Create(t *testing.T) {
	type args struct {
		image *model.RecipeImage
	}
	type mocks struct {
		setup func(m *mockPostgres.MockDatabase, r *mockPostgres.MockRows)
	}
	tests := []struct {
		name    string
		args    args
		mocks   mocks
		wantErr bool
	}{
		{
			name: "Success",
			args: args{image: &model.RecipeImage{RecipeID: 100, ImagePath: "/path/to/image.jpg"}},
			mocks: mocks{
				setup: func(m *mockPostgres.MockDatabase, r *mockPostgres.MockRows) {
					r.EXPECT().Next().Return(true)
					r.EXPECT().Scan(gomock.Any(), gomock.Any()).DoAndReturn(
						func(dest ...any) error {
							*dest[0].(*int64) = 1
							*dest[1].(*time.Time) = time.Now()
							return nil
						},
					)
					r.EXPECT().Close()
					m.EXPECT().Query(gomock.Any(), gomock.Any(), int64(100), "/path/to/image.jpg").Return(r, nil)
				},
			},
			wantErr: false,
		},
		{
			name: "DB Error",
			args: args{image: &model.RecipeImage{RecipeID: 100, ImagePath: "/path/to/image.jpg"}},
			mocks: mocks{
				setup: func(m *mockPostgres.MockDatabase, r *mockPostgres.MockRows) {
					m.EXPECT().Query(gomock.Any(), gomock.Any(), int64(100), "/path/to/image.jpg").
						Return(nil, errors.New("db error"))
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mockPostgres.NewMockDatabase(ctrl)
			mockRows := mockPostgres.NewMockRows(ctrl)
			tt.mocks.setup(mockDB, mockRows)

			repo := postgres.NewRecipeImageRepository(mockDB)
			err := repo.Create(context.Background(), tt.args.image)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int64(1), tt.args.image.ID)
			}
		})
	}
}

func TestRecipeImageRepository_GetByRecipeID(t *testing.T) {
	type args struct {
		recipeID int64
	}
	type mocks struct {
		setup func(m *mockPostgres.MockDatabase, r *mockPostgres.MockRows)
	}
	tests := []struct {
		name    string
		args    args
		mocks   mocks
		wantLen int
		wantErr bool
	}{
		{
			name: "Success",
			args: args{recipeID: 100},
			mocks: mocks{
				setup: func(m *mockPostgres.MockDatabase, r *mockPostgres.MockRows) {
					r.EXPECT().Next().Return(true)
					r.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
						func(dest ...any) error {
							*dest[0].(*int64) = 1
							*dest[1].(*int64) = 100
							*dest[2].(*string) = "/path/to/image.jpg"
							*dest[3].(*time.Time) = time.Now()
							return nil
						},
					)
					r.EXPECT().Next().Return(false)
					r.EXPECT().Close()

					m.EXPECT().Query(gomock.Any(), gomock.Any(), int64(100)).Return(r, nil)
				},
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "Error",
			args: args{recipeID: 100},
			mocks: mocks{
				setup: func(m *mockPostgres.MockDatabase, r *mockPostgres.MockRows) {
					m.EXPECT().Query(gomock.Any(), gomock.Any(), int64(100)).Return(nil, errors.New("db error"))
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mockPostgres.NewMockDatabase(ctrl)
			mockRows := mockPostgres.NewMockRows(ctrl)
			tt.mocks.setup(mockDB, mockRows)

			repo := postgres.NewRecipeImageRepository(mockDB)
			got, err := repo.GetByRecipeID(context.Background(), tt.args.recipeID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, got, tt.wantLen)
			}
		})
	}
}

func TestRecipeImageRepository_Delete(t *testing.T) {
	type args struct {
		id int64
	}
	type mocks struct {
		setup func(m *mockPostgres.MockDatabase)
	}
	tests := []struct {
		name    string
		args    args
		mocks   mocks
		wantErr bool
	}{
		{
			name: "Success",
			args: args{id: 1},
			mocks: mocks{
				setup: func(m *mockPostgres.MockDatabase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any(), int64(1)).Return(int64(1), nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Error",
			args: args{id: 1},
			mocks: mocks{
				setup: func(m *mockPostgres.MockDatabase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any(), int64(1)).Return(int64(0), errors.New("failed"))
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mockPostgres.NewMockDatabase(ctrl)
			tt.mocks.setup(mockDB)

			repo := postgres.NewRecipeImageRepository(mockDB)
			err := repo.Delete(context.Background(), tt.args.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

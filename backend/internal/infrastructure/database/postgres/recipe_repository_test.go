package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database/postgres"
	mock_postgres "github.com/seka/reci-pin/backend/internal/infrastructure/database/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRecipeRepository_Create(t *testing.T) {
	type args struct {
		recipe *model.Recipe
	}
	type mocks struct {
		setup func(m *mock_postgres.MockDatabase, r *mock_postgres.MockRows)
	}
	tests := []struct {
		name    string
		args    args
		mocks   mocks
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				recipe: &model.Recipe{
					UserID: 1,
					Name:   "Pancakes",
					URL:    "http://example.com",
					Memo:   "Yummy",
				},
			},
			mocks: mocks{
				setup: func(m *mock_postgres.MockDatabase, r *mock_postgres.MockRows) {
					r.EXPECT().Next().Return(true)
					r.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
						func(dest ...interface{}) error {
							*dest[0].(*int64) = 100
							*dest[1].(*time.Time) = time.Now()
							*dest[2].(*time.Time) = time.Now()
							return nil
						},
					)
					r.EXPECT().Close()

					m.EXPECT().Query(gomock.Any(), gomock.Any(), int64(1), "Pancakes", "http://example.com", "Yummy").Return(r, nil)
				},
			},
			wantErr: false,
		},
		{
			name: "DB Error",
			args: args{
				recipe: &model.Recipe{UserID: 1, Name: "Pancakes"},
			},
			mocks: mocks{
				setup: func(m *mock_postgres.MockDatabase, r *mock_postgres.MockRows) {
					m.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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

			mockDB := mock_postgres.NewMockDatabase(ctrl)
			mockRows := mock_postgres.NewMockRows(ctrl)
			tt.mocks.setup(mockDB, mockRows)

			repo := postgres.NewRecipeRepository(mockDB)
			err := repo.Create(context.Background(), tt.args.recipe)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int64(100), tt.args.recipe.ID)
			}
		})
	}
}

func TestRecipeRepository_GetByID(t *testing.T) {
	type args struct {
		id int64
	}
	type mocks struct {
		setup func(m *mock_postgres.MockDatabase, r *mock_postgres.MockRows)
	}
	tests := []struct {
		name    string
		args    args
		mocks   mocks
		want    *model.Recipe
		wantErr bool
	}{
		{
			name: "Success",
			args: args{id: 100},
			mocks: mocks{
				setup: func(m *mock_postgres.MockDatabase, r *mock_postgres.MockRows) {
					// Mock GetByID Query
					r.EXPECT().Next().Return(true)
					r.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
						func(dest ...interface{}) error {
							*dest[0].(*int64) = 100
							*dest[1].(*int64) = 1
							*dest[2].(*string) = "Pancakes"
							*dest[3].(*string) = "http://example.com"
							*dest[4].(*string) = "Delicious"
							*dest[5].(*time.Time) = time.Now()
							*dest[6].(*time.Time) = time.Now()
							return nil
						},
					)
					r.EXPECT().Close()

					m.EXPECT().Query(gomock.Any(), gomock.Any(), int64(100)).Return(r, nil)
				},
			},
			want: &model.Recipe{
				ID:     100,
				UserID: 1,
				Name:   "Pancakes",
				URL:    "http://example.com",
				Memo:   "Delicious",
			},
			wantErr: false,
		},
		{
			name: "Not Found",
			args: args{id: 999},
			mocks: mocks{
				setup: func(m *mock_postgres.MockDatabase, r *mock_postgres.MockRows) {
					r.EXPECT().Next().Return(false)
					r.EXPECT().Close()

					m.EXPECT().Query(gomock.Any(), gomock.Any(), int64(999)).Return(r, nil)
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mock_postgres.NewMockDatabase(ctrl)
			mockRows := mock_postgres.NewMockRows(ctrl)
			tt.mocks.setup(mockDB, mockRows)

			repo := postgres.NewRecipeRepository(mockDB)
			got, err := repo.GetByID(context.Background(), tt.args.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.Name, got.Name)
			}
		})
	}
}

func TestRecipeRepository_Search(t *testing.T) {
	type args struct {
		userID int64
		query  string
		tagIDs []int64
	}
	type mocks struct {
		setup func(m *mock_postgres.MockDatabase, r *mock_postgres.MockRows)
	}
	tests := []struct {
		name    string
		args    args
		mocks   mocks
		deferFn func()
		wantLen int
		wantErr bool
	}{
		{
			name: "Search by Query",
			args: args{userID: 1, query: "Pan", tagIDs: nil},
			mocks: mocks{
				setup: func(m *mock_postgres.MockDatabase, r *mock_postgres.MockRows) {
					r.EXPECT().Next().Return(true)
					r.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
						func(dest ...interface{}) error {
							*dest[0].(*int64) = 1
							*dest[1].(*int64) = 1
							*dest[2].(*string) = "Pancakes"
							*dest[3].(*string) = ""
							*dest[4].(*string) = ""
							*dest[5].(*time.Time) = time.Now()
							*dest[6].(*time.Time) = time.Now()
							return nil
						},
					)
					r.EXPECT().Next().Return(false)
					r.EXPECT().Close()

					// Query arg construction verification handled by checking implementation manually
					// or expecting specific args. Here we relax to Any() for SQL and exact args.
					m.EXPECT().Query(gomock.Any(), gomock.Any(), int64(1), "%Pan%").Return(r, nil)
				},
			},
			wantLen: 1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mock_postgres.NewMockDatabase(ctrl)
			mockRows := mock_postgres.NewMockRows(ctrl)
			tt.mocks.setup(mockDB, mockRows)

			repo := postgres.NewRecipeRepository(mockDB)
			got, err := repo.Search(context.Background(), tt.args.userID, tt.args.query, tt.args.tagIDs)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, got, tt.wantLen)
			}
		})
	}
}

func TestRecipeRepository_Update(t *testing.T) {
	type args struct {
		recipe *model.Recipe
	}
	type mocks struct {
		setup func(m *mock_postgres.MockDatabase)
	}
	tests := []struct {
		name    string
		args    args
		mocks   mocks
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				recipe: &model.Recipe{ID: 100, Name: "New Name", URL: "url", Memo: "memo"},
			},
			mocks: mocks{
				setup: func(m *mock_postgres.MockDatabase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any(), "New Name", "url", "memo", int64(100)).Return(int64(1), nil)
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mock_postgres.NewMockDatabase(ctrl)
			tt.mocks.setup(mockDB)

			repo := postgres.NewRecipeRepository(mockDB)
			err := repo.Update(context.Background(), tt.args.recipe)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

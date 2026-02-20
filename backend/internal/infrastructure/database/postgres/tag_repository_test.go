package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	mock_postgres "github.com/seka/reci-pin/backend/internal/infrastructure/database/mock"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database/postgres"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTagRepository_Create(t *testing.T) {
	type args struct {
		tag *model.Tag
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
			args: args{tag: &model.Tag{Name: "Vegan"}},
			mocks: mocks{
				setup: func(m *mock_postgres.MockDatabase, r *mock_postgres.MockRows) {
					r.EXPECT().Next().Return(true)
					r.EXPECT().Scan(gomock.Any()).DoAndReturn(
						func(dest ...any) error {
							*dest[0].(*int64) = 1
							return nil
						},
					)
					r.EXPECT().Close()
					m.EXPECT().Query(gomock.Any(), gomock.Any(), "Vegan").Return(r, nil)
				},
			},
			wantErr: false,
		},
		{
			name: "DB Error",
			args: args{tag: &model.Tag{Name: "Vegan"}},
			mocks: mocks{
				setup: func(m *mock_postgres.MockDatabase, r *mock_postgres.MockRows) {
					m.EXPECT().Query(gomock.Any(), gomock.Any(), "Vegan").Return(nil, errors.New("db error"))
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

			repo := postgres.NewTagRepository(mockDB)
			err := repo.Create(context.Background(), tt.args.tag)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int64(1), tt.args.tag.ID)
			}
		})
	}
}

func TestTagRepository_GetAll(t *testing.T) {
	type mocks struct {
		setup func(m *mock_postgres.MockDatabase, r *mock_postgres.MockRows)
	}
	tests := []struct {
		name    string
		mocks   mocks
		wantLen int
		wantErr bool
	}{
		{
			name: "Success",
			mocks: mocks{
				setup: func(m *mock_postgres.MockDatabase, r *mock_postgres.MockRows) {
					r.EXPECT().Next().Return(true)
					r.EXPECT().Scan(gomock.Any(), gomock.Any()).DoAndReturn(
						func(dest ...any) error {
							*dest[0].(*int64) = 1
							*dest[1].(*string) = "Vegan"
							return nil
						},
					)
					r.EXPECT().Next().Return(false)
					r.EXPECT().Close()

					m.EXPECT().Query(gomock.Any(), gomock.Any()).Return(r, nil)
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

			repo := postgres.NewTagRepository(mockDB)
			got, err := repo.GetAll(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, got, tt.wantLen)
			}
		})
	}
}

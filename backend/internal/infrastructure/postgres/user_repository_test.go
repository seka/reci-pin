package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/infrastructure/postgres"
	mock_postgres "github.com/seka/reci-pin/backend/internal/infrastructure/postgres/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserRepository_Create(t *testing.T) {
	type args struct {
		user *model.User
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
				user: &model.User{Name: "Test User"},
			},
			mocks: mocks{
				setup: func(m *mock_postgres.MockDatabase, r *mock_postgres.MockRows) {
					// Mock Rows for RETURNING clause
					r.EXPECT().Next().Return(true)
					r.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
						func(dest ...interface{}) error {
							*dest[0].(*int64) = 1
							*dest[1].(*time.Time) = time.Now()
							*dest[2].(*time.Time) = time.Now()
							return nil
						},
					)
					r.EXPECT().Close()

					m.EXPECT().Query(gomock.Any(), gomock.Any(), "Test User").Return(r, nil)
				},
			},
			wantErr: false,
		},
		{
			name: "DB Error",
			args: args{
				user: &model.User{Name: "Test User"},
			},
			mocks: mocks{
				setup: func(m *mock_postgres.MockDatabase, r *mock_postgres.MockRows) {
					m.EXPECT().Query(gomock.Any(), gomock.Any(), "Test User").Return(nil, errors.New("db error"))
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

			repo := postgres.NewUserRepository(mockDB)
			err := repo.Create(context.Background(), tt.args.user)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.args.user.ID)
			}
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
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
		want    *model.User
		wantErr bool
	}{
		{
			name: "Success",
			args: args{id: 1},
			mocks: mocks{
				setup: func(m *mock_postgres.MockDatabase, r *mock_postgres.MockRows) {
					r.EXPECT().Next().Return(true)
					r.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
						func(dest ...interface{}) error {
							*dest[0].(*int64) = 1
							*dest[1].(*string) = "Test User"
							*dest[2].(*time.Time) = time.Now()
							*dest[3].(*time.Time) = time.Now()
							return nil
						},
					)
					r.EXPECT().Close()

					m.EXPECT().Query(gomock.Any(), gomock.Any(), int64(1)).Return(r, nil)
				},
			},
			want:    &model.User{ID: 1, Name: "Test User"},
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
			want:    nil,
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

			repo := postgres.NewUserRepository(mockDB)
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

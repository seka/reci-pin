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

func TestRefreshTokenRepository_Save(t *testing.T) {
	type args struct {
		token *model.RefreshToken
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
			args: args{
				token: &model.RefreshToken{
					UserID:    1,
					TokenHash: "hash",
					ExpiresAt: time.Now().Add(24 * time.Hour),
					UserAgent: "agent",
					IPAddress: "1.1.1.1",
				},
			},
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

					m.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(r, nil)
				},
			},
			wantErr: false,
		},
		{
			name: "DB Error",
			args: args{
				token: &model.RefreshToken{UserID: 1},
			},
			mocks: mocks{
				setup: func(m *mockPostgres.MockDatabase, r *mockPostgres.MockRows) {
					m.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))
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

			repo := postgres.NewRefreshTokenRepository(mockDB)
			err := repo.Save(context.Background(), tt.args.token)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.args.token.ID)
			}
		})
	}
}

func TestRefreshTokenRepository_GetByHash(t *testing.T) {
	type args struct {
		hash string
	}
	type mocks struct {
		setup func(m *mockPostgres.MockDatabase, r *mockPostgres.MockRows)
	}
	tests := []struct {
		name    string
		args    args
		mocks   mocks
		want    *model.RefreshToken
		wantErr bool
	}{
		{
			name: "Success",
			args: args{hash: "hash"},
			mocks: mocks{
				setup: func(m *mockPostgres.MockDatabase, r *mockPostgres.MockRows) {
					r.EXPECT().Next().Return(true)
					r.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
						func(dest ...any) error {
							*dest[0].(*int64) = 1
							*dest[1].(*int64) = 1
							*dest[2].(*string) = "hash"
							*dest[3].(*time.Time) = time.Now()
							*dest[4].(*time.Time) = time.Now()
							*dest[5].(**time.Time) = nil
							*dest[6].(*string) = "agent"
							*dest[7].(*string) = "1.1.1.1"
							return nil
						},
					)
					r.EXPECT().Close()

					m.EXPECT().Query(gomock.Any(), gomock.Any(), "hash").Return(r, nil)
				},
			},
			want:    &model.RefreshToken{ID: 1, UserID: 1, TokenHash: "hash"},
			wantErr: false,
		},
		{
			name: "Not Found",
			args: args{hash: "notfound"},
			mocks: mocks{
				setup: func(m *mockPostgres.MockDatabase, r *mockPostgres.MockRows) {
					r.EXPECT().Next().Return(false)
					r.EXPECT().Close()

					m.EXPECT().Query(gomock.Any(), gomock.Any(), "notfound").Return(r, nil)
				},
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mockPostgres.NewMockDatabase(ctrl)
			mockRows := mockPostgres.NewMockRows(ctrl)
			tt.mocks.setup(mockDB, mockRows)

			repo := postgres.NewRefreshTokenRepository(mockDB)
			got, err := repo.GetByHash(context.Background(), tt.args.hash)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.want == nil {
					assert.Nil(t, got)
				} else {
					assert.NotNil(t, got)
					assert.Equal(t, tt.want.ID, got.ID)
					assert.Equal(t, tt.want.TokenHash, got.TokenHash)
				}
			}
		})
	}
}

func TestRefreshTokenRepository_Revoke(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockPostgres.NewMockDatabase(ctrl)
	mockDB.EXPECT().Execute(gomock.Any(), gomock.Any(), int64(1)).Return(int64(1), nil)

	repo := postgres.NewRefreshTokenRepository(mockDB)
	err := repo.Revoke(context.Background(), 1)
	assert.NoError(t, err)
}

func TestRefreshTokenRepository_RevokeAllByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockPostgres.NewMockDatabase(ctrl)
	mockDB.EXPECT().Execute(gomock.Any(), gomock.Any(), int64(1)).Return(int64(1), nil)

	repo := postgres.NewRefreshTokenRepository(mockDB)
	err := repo.RevokeAllByUserID(context.Background(), 1)
	assert.NoError(t, err)
}

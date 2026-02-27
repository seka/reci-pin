package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	mockPostgres "github.com/seka/reci-pin/backend/internal/infrastructure/database/mock"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database/postgres"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserEmailCredentialRepository_Create(t *testing.T) {
	type args struct {
		credential *model.UserEmailCredential
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
			args: args{
				credential: &model.UserEmailCredential{
					UserID:       1,
					Email:        "test@example.com",
					PasswordHash: "hash",
				},
			},
			mocks: mocks{
				setup: func(m *mockPostgres.MockDatabase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any(), int64(1), "test@example.com", "hash", gomock.Any(), gomock.Any(), gomock.Any()).
						Return(int64(1), nil)
				},
			},
			wantErr: false,
		},
		{
			name: "DB Error",
			args: args{credential: &model.UserEmailCredential{UserID: 1, Email: "test@example.com"}},
			mocks: mocks{
				setup: func(m *mockPostgres.MockDatabase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(int64(0), errors.New("db error"))
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

			repo := postgres.NewUserEmailCredentialRepository(mockDB)
			err := repo.Create(context.Background(), tt.args.credential)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserEmailCredentialRepository_GetByEmail(t *testing.T) {
	type args struct {
		email string
	}
	type mocks struct {
		setup func(m *mockPostgres.MockDatabase, r *mockPostgres.MockRows)
	}
	tests := []struct {
		name    string
		args    args
		mocks   mocks
		want    *model.UserEmailCredential
		wantErr bool
		errIs   error
	}{
		{
			name: "Success",
			args: args{email: "test@example.com"},
			mocks: mocks{
				setup: func(m *mockPostgres.MockDatabase, r *mockPostgres.MockRows) {
					r.EXPECT().Next().Return(true)
					parsedTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
					r.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
						func(dest ...any) error {
							*dest[0].(*int64) = 1
							*dest[1].(*string) = "test@example.com"
							*dest[2].(*string) = "hash"
							if p, ok := dest[3].(**time.Time); ok {
								*p = &parsedTime
							}
							*dest[4].(*string) = "token"
							if p, ok := dest[5].(**time.Time); ok {
								*p = &parsedTime
							}
							*dest[6].(*time.Time) = parsedTime
							return nil
						},
					)
					r.EXPECT().Close()
					m.EXPECT().Query(gomock.Any(), gomock.Any(), "test@example.com").Return(r, nil)
				},
			},
			want: &model.UserEmailCredential{
				UserID:       1,
				Email:        "test@example.com",
				PasswordHash: "hash",
				VerificationToken: "token",
			},
			wantErr: false,
		},
		{
			name: "Not Found",
			args: args{email: "notfound@example.com"},
			mocks: mocks{
				setup: func(m *mockPostgres.MockDatabase, r *mockPostgres.MockRows) {
					r.EXPECT().Next().Return(false)
					r.EXPECT().Close()
					m.EXPECT().Query(gomock.Any(), gomock.Any(), "notfound@example.com").Return(r, nil)
				},
			},
			want:    nil,
			wantErr: true,
			errIs:   repository.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mockPostgres.NewMockDatabase(ctrl)
			mockRows := mockPostgres.NewMockRows(ctrl)
			tt.mocks.setup(mockDB, mockRows)

			repo := postgres.NewUserEmailCredentialRepository(mockDB)
			got, err := repo.GetByEmail(context.Background(), tt.args.email)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errIs != nil {
					assert.True(t, errors.Is(err, tt.errIs))
				}
			} else {
				assert.NoError(t, err)
				if tt.want != nil {
					assert.Equal(t, tt.want.UserID, got.UserID)
					assert.Equal(t, tt.want.Email, got.Email)
				} else {
					assert.Nil(t, got)
				}
			}
		})
	}
}

func TestUserEmailCredentialRepository_Update(t *testing.T) {
	type args struct {
		credential *model.UserEmailCredential
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
			args: args{
				credential: &model.UserEmailCredential{UserID: 1, Email: "new@example.com"},
			},
			mocks: mocks{
				setup: func(m *mockPostgres.MockDatabase) {
					m.EXPECT().Execute(gomock.Any(), gomock.Any(), int64(1), "new@example.com", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(int64(1), nil)
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mockPostgres.NewMockDatabase(ctrl)
			tt.mocks.setup(mockDB)

			repo := postgres.NewUserEmailCredentialRepository(mockDB)
			err := repo.Update(context.Background(), tt.args.credential)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

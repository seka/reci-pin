package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	repositorymock "github.com/seka/reci-pin/backend/internal/domain/repository/mock"
	"github.com/seka/reci-pin/backend/internal/usecase/auth"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUserUseCase_Execute(t *testing.T) {
	type mocks struct {
		setup func(m *repositorymock.MockUserRepository)
	}
	type args struct {
		userID int64
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
			args: args{userID: 1},
			mocks: mocks{
				setup: func(m *repositorymock.MockUserRepository) {
					m.EXPECT().GetByID(gomock.Any(), int64(1)).Return(&model.User{ID: 1, Name: "Test User"}, nil)
				},
			},
			want:    &model.User{ID: 1, Name: "Test User"},
			wantErr: false,
		},
		{
			name: "Error",
			args: args{userID: 999},
			mocks: mocks{
				setup: func(m *repositorymock.MockUserRepository) {
					m.EXPECT().GetByID(gomock.Any(), int64(999)).Return(nil, errors.New("db error"))
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

			repo := repositorymock.NewMockUserRepository(ctrl)
			tt.mocks.setup(repo)

			uc := auth.NewGetUserUseCase(repo)
			got, err := uc.Execute(context.Background(), tt.args.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

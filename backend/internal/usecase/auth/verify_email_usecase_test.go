package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	repositorymock "github.com/seka/reci-pin/backend/internal/domain/repository/mock"
	"github.com/seka/reci-pin/backend/internal/usecase/auth"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestVerifyEmailUseCase_Execute(t *testing.T) {
	type mocks struct {
		setup func(m *repositorymock.MockUserEmailCredentialRepository, mtm *repositorymock.MockTransactionManager)
	}
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		mocks   mocks
		wantErr bool
		errMsg  string
	}{
		{
			name: "Success",
			args: args{token: "valid-token"},
			mocks: mocks{
				setup: func(m *repositorymock.MockUserEmailCredentialRepository, mtm *repositorymock.MockTransactionManager) {
					expiresAt := time.Now().Add(1 * time.Hour)
					cred := &model.UserEmailCredential{
						UserID:                     1,
						Email:                      "test@example.com",
						VerificationToken:          "valid-token",
						VerificationTokenExpiresAt: &expiresAt,
					}
					mtm.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						},
					)
					m.EXPECT().GetByToken(gomock.Any(), "valid-token").Return(cred, nil)
					// Expect Update with cleared token and verified time
					m.EXPECT().Update(gomock.Any(), gomock.AssignableToTypeOf(cred)).DoAndReturn(
						func(_ context.Context, c *model.UserEmailCredential) error {
							assert.Empty(t, c.VerificationToken)
							assert.Nil(t, c.VerificationTokenExpiresAt)
							assert.NotNil(t, c.EmailVerifiedAt)
							return nil
						},
					)
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid Token",
			args: args{token: "invalid-token"},
			mocks: mocks{
				setup: func(m *repositorymock.MockUserEmailCredentialRepository, mtm *repositorymock.MockTransactionManager) {
					mtm.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						},
					)
					m.EXPECT().GetByToken(gomock.Any(), "invalid-token").Return(nil, errors.New("db error"))
				},
			},
			wantErr: true,
			errMsg:  "invalid token",
		},
		{
			name: "Already Verified",
			args: args{token: "verified-token"},
			mocks: mocks{
				setup: func(m *repositorymock.MockUserEmailCredentialRepository, mtm *repositorymock.MockTransactionManager) {
					mtm.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						},
					)
					now := time.Now()
					cred := &model.UserEmailCredential{
						UserID:          1,
						EmailVerifiedAt: &now,
					}
					m.EXPECT().GetByToken(gomock.Any(), "verified-token").Return(cred, nil)
				},
			},
			wantErr: true,
			errMsg:  "email already verified",
		},
		{
			name: "Token Expired",
			args: args{token: "expired-token"},
			mocks: mocks{
				setup: func(m *repositorymock.MockUserEmailCredentialRepository, mtm *repositorymock.MockTransactionManager) {
					mtm.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						},
					)
					expiresAt := time.Now().Add(-1 * time.Hour)
					cred := &model.UserEmailCredential{
						UserID:                     1,
						VerificationToken:          "expired-token",
						VerificationTokenExpiresAt: &expiresAt,
					}
					m.EXPECT().GetByToken(gomock.Any(), "expired-token").Return(cred, nil)
				},
			},
			wantErr: true,
			errMsg:  "token expired",
		},
		{
			name: "Update Error",
			args: args{token: "valid-token"},
			mocks: mocks{
				setup: func(m *repositorymock.MockUserEmailCredentialRepository, mtm *repositorymock.MockTransactionManager) {
					mtm.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
						func(ctx context.Context, fn func(context.Context) error) error {
							return fn(ctx)
						},
					)
					expiresAt := time.Now().Add(1 * time.Hour)
					cred := &model.UserEmailCredential{
						UserID:                     1,
						VerificationToken:          "valid-token",
						VerificationTokenExpiresAt: &expiresAt,
					}
					m.EXPECT().GetByToken(gomock.Any(), "valid-token").Return(cred, nil)
					m.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("db fail"))
				},
			},
			wantErr: true,
			errMsg:  "failed to verify email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := repositorymock.NewMockUserEmailCredentialRepository(ctrl)
			mtm := repositorymock.NewMockTransactionManager(ctrl)
			tt.mocks.setup(repo, mtm)

			uc := auth.NewVerifyEmailUseCase(repo, mtm)
			err := uc.Execute(context.Background(), tt.args.token)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

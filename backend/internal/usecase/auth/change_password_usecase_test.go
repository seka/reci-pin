package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	emailmock "github.com/seka/reci-pin/backend/internal/domain/notification/mock"
	repositorymock "github.com/seka/reci-pin/backend/internal/domain/repository/mock"
	"github.com/seka/reci-pin/backend/internal/infrastructure/database/postgres"
	"github.com/seka/reci-pin/backend/internal/usecase/auth"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestChangePasswordUseCase_Execute(t *testing.T) {
	// Create a dummy hash for testing
	currentPassword := "oldPassword123"
	currentPasswordHash, _ := postgres.HashPassword(currentPassword)
	newPassword := "newPassword123"
	// newPasswordHash will be generated inside the usecase

	tests := []struct {
		name    string
		input   auth.ChangePasswordInput
		setup   func(credRepo *repositorymock.MockUserEmailCredentialRepository, emailSender *emailmock.MockEmailClient)
		wantErr bool
	}{
		{
			name: "Success",
			input: auth.ChangePasswordInput{
				CurrentPassword: currentPassword,
				NewPassword:     newPassword,
			},
			setup: func(credRepo *repositorymock.MockUserEmailCredentialRepository, emailSender *emailmock.MockEmailClient) {
				credRepo.EXPECT().GetByUserID(gomock.Any(), int64(1)).Return(&model.UserEmailCredential{
					UserID:       1,
					Email:        "test@example.com",
					PasswordHash: currentPasswordHash,
				}, nil)
				credRepo.EXPECT().UpdatePassword(gomock.Any(), int64(1), gomock.Any()).Return(nil)
				emailSender.EXPECT().SendPasswordChangeNotification("test@example.com").Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Invalid New Password",
			input: auth.ChangePasswordInput{
				CurrentPassword: currentPassword,
				NewPassword:     "short",
			},
			setup: func(credRepo *repositorymock.MockUserEmailCredentialRepository, emailSender *emailmock.MockEmailClient) {
				// No calls expected
			},
			wantErr: true,
		},
		{
			name: "Incorrect Current Password",
			input: auth.ChangePasswordInput{
				CurrentPassword: "wrongPassword",
				NewPassword:     newPassword,
			},
			setup: func(credRepo *repositorymock.MockUserEmailCredentialRepository, emailSender *emailmock.MockEmailClient) {
				credRepo.EXPECT().GetByUserID(gomock.Any(), int64(1)).Return(&model.UserEmailCredential{
					UserID:       1,
					Email:        "test@example.com",
					PasswordHash: currentPasswordHash,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "User Not Found",
			input: auth.ChangePasswordInput{
				CurrentPassword: currentPassword,
				NewPassword:     newPassword,
			},
			setup: func(credRepo *repositorymock.MockUserEmailCredentialRepository, emailSender *emailmock.MockEmailClient) {
				credRepo.EXPECT().GetByUserID(gomock.Any(), int64(1)).Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name: "Database Error on Get",
			input: auth.ChangePasswordInput{
				CurrentPassword: currentPassword,
				NewPassword:     newPassword,
			},
			setup: func(credRepo *repositorymock.MockUserEmailCredentialRepository, emailSender *emailmock.MockEmailClient) {
				credRepo.EXPECT().GetByUserID(gomock.Any(), int64(1)).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "Database Error on Update",
			input: auth.ChangePasswordInput{
				CurrentPassword: currentPassword,
				NewPassword:     newPassword,
			},
			setup: func(credRepo *repositorymock.MockUserEmailCredentialRepository, emailSender *emailmock.MockEmailClient) {
				credRepo.EXPECT().GetByUserID(gomock.Any(), int64(1)).Return(&model.UserEmailCredential{
					UserID:       1,
					Email:        "test@example.com",
					PasswordHash: currentPasswordHash,
				}, nil)
				credRepo.EXPECT().UpdatePassword(gomock.Any(), int64(1), gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCredRepo := repositorymock.NewMockUserEmailCredentialRepository(ctrl)
			mockEmailClient := emailmock.NewMockEmailClient(ctrl)

			tt.setup(mockCredRepo, mockEmailClient)

			uc := auth.NewChangePasswordUseCase(mockCredRepo, mockEmailClient)
			err := uc.Execute(context.Background(), 1, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

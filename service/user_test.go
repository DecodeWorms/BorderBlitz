package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/DecodeWorms/BorderBlitz/mocks" // assuming the mock is here
	"github.com/DecodeWorms/BorderBlitz/models"
)


func TestUsers_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mocks.NewMockUser(ctrl)
	usersService := &Users{user: mockUser}

	tests := []struct {
		name        string
		input       *models.CreateUserRequest
		mockSetup   func()
		expectError bool
	}{
		{
			name: "successfully created a user",
			input: &models.CreateUserRequest{
				UserType: "email", // Could be email or mobile
			},
			mockSetup: func() {
				mockUser.
					EXPECT().
					Create(gomock.AssignableToTypeOf(&models.Users{})).
					Return(nil)
			},
			expectError: false,
		},
		{
			name: "failed to create a user",
			input: &models.CreateUserRequest{
				UserType: "regular", // This isn neither email or mobile
			},
			mockSetup: func() {
				mockUser.
					EXPECT().
					Create(gomock.AssignableToTypeOf(&models.Users{})).
					Return(errors.New("failed to create user"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := usersService.Create(tt.input)

			if tt.expectError {
				assert.Nil(t, result)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, result)
				assert.NoError(t, err)
				assert.Equal(t, tt.input.UserType, result.UserType)
			}
		})
	}
}

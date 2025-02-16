package authservice_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
	"github.com/vit6556/avito-internship-assignment/internal/database/mock"
	"github.com/vit6556/avito-internship-assignment/internal/entity"
	"github.com/vit6556/avito-internship-assignment/internal/service"
	"github.com/vit6556/avito-internship-assignment/internal/service/auth"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthorizeUser(t *testing.T) {
	ctx := context.Background()
	mockEmployeeRepo := new(mock.MockEmployeeRepository)
	authService := authservice.NewAuthService(mockEmployeeRepo, "secret", time.Hour, 1000)

	tests := []struct {
		name          string
		existingUser  *entity.Employee
		username      string
		password      string
		mockSetup     func()
		expectedError error
	}{
		{
			name:         "Success - Existing User Valid Credentials",
			existingUser: &entity.Employee{ID: 1, Username: "alice", PasswordHash: hashPasswordHelper("password")},
			username:     "alice",
			password:     "password",
			mockSetup: func() {
				mockEmployeeRepo.ExpectedCalls = nil
				mockEmployeeRepo.On("GetEmployeeByUsername", ctx, "alice").
					Return(&entity.Employee{ID: 1, Username: "alice", PasswordHash: hashPasswordHelper("password")}, nil)
			},
			expectedError: nil,
		},
		{
			name:         "Error - Existing User Invalid Credentials",
			existingUser: &entity.Employee{ID: 1, Username: "alice", PasswordHash: hashPasswordHelper("password")},
			username:     "alice",
			password:     "wrongpassword",
			mockSetup: func() {
				mockEmployeeRepo.ExpectedCalls = nil
				mockEmployeeRepo.On("GetEmployeeByUsername", ctx, "alice").
					Return(&entity.Employee{ID: 1, Username: "alice", PasswordHash: hashPasswordHelper("password")}, nil)
			},
			expectedError: service.ErrInvalidCredentials,
		},
		{
			name:         "Success - New User Creation",
			existingUser: nil,
			username:     "bob",
			password:     "newpassword",
			mockSetup: func() {
				mockEmployeeRepo.ExpectedCalls = nil
				mockEmployeeRepo.On("GetEmployeeByUsername", ctx, "bob").
					Return(nil, service.ErrEmployeeNotFound)
				mockEmployeeRepo.On("CreateEmployee", ctx, testifyMock.Anything).
					Return(2, nil)
			},
			expectedError: nil,
		},
		{
			name:         "Error - Failed to Create New User",
			existingUser: nil,
			username:     "bob",
			password:     "newpassword",
			mockSetup: func() {
				mockEmployeeRepo.ExpectedCalls = nil
				mockEmployeeRepo.On("GetEmployeeByUsername", ctx, "bob").
					Return(nil, service.ErrEmployeeNotFound)
				mockEmployeeRepo.On("CreateEmployee", ctx, testifyMock.Anything).
					Return(0, service.ErrEmployeeCreationFailed)
			},
			expectedError: service.ErrEmployeeCreationFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			token, err := authService.AuthorizeUser(ctx, tt.username, tt.password)

			assert.Equal(t, tt.expectedError, err)
			if tt.expectedError == nil {
				assert.NotEmpty(t, token)
			}

			mockEmployeeRepo.AssertExpectations(t)
		})
	}
}

func TestValidateToken(t *testing.T) {
	mockEmployeeRepo := new(mock.MockEmployeeRepository)
	authService := authservice.NewAuthService(mockEmployeeRepo, "secret", time.Hour, 1000)

	tests := []struct {
		name          string
		token         string
		mockSetup     func()
		expectedUser  int
		expectedError error
	}{
		{
			name:          "Success - Valid Token",
			token:         generateToken(1, "secret"),
			mockSetup:     func() {},
			expectedUser:  1,
			expectedError: nil,
		},
		{
			name:          "Error - Invalid Token",
			token:         "invalid.token.string",
			mockSetup:     func() {},
			expectedUser:  0,
			expectedError: service.ErrInvalidToken,
		},
		{
			name:          "Error - Expired Token",
			token:         generateExpiredToken(1, "secret"),
			mockSetup:     func() {},
			expectedUser:  0,
			expectedError: service.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			userID, err := authService.ValidateToken(tt.token)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedUser, userID)
		})
	}
}

func hashPasswordHelper(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func generateToken(userID int, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func generateExpiredToken(userID int, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(-time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

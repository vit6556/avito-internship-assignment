package httphandler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
	"github.com/vit6556/avito-internship-assignment/internal/delivery/http/dto"
	"github.com/vit6556/avito-internship-assignment/internal/delivery/http/handler"
	"github.com/vit6556/avito-internship-assignment/internal/service"
	"github.com/vit6556/avito-internship-assignment/internal/service/mock"
)

func TestGetToken(t *testing.T) {
	e := echo.New()
	mockAuthService := new(mock.MockAuthService)
	handler := httphandler.NewAuthHandler(mockAuthService, time.Hour, false)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Success - Valid credentials",
			requestBody: dto.AuthRequest{
				Username: "testuser",
				Password: "password123",
			},
			mockSetup: func() {
				mockAuthService.On("AuthorizeUser", testifyMock.Anything, "testuser", "password123").
					Return("valid-token", nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "Error - Invalid JSON",
			requestBody:    "invalid-json",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid JSON data",
		},
		{
			name: "Error - Missing fields",
			requestBody: dto.AuthRequest{
				Username: "",
				Password: "",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid request data",
		},
		{
			name: "Error - Invalid credentials",
			requestBody: dto.AuthRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			mockSetup: func() {
				mockAuthService.On("AuthorizeUser", testifyMock.Anything, "testuser", "wrongpassword").
					Return("", service.ErrInvalidCredentials)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid username or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			var reqBody []byte
			if str, ok := tt.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.GetToken(c)

			if tt.expectedError != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedStatus, rec.Code)
			mockAuthService.AssertExpectations(t)
		})
	}
}

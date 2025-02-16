package httphandler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"

	"github.com/vit6556/avito-internship-assignment/internal/delivery/http/dto"
	"github.com/vit6556/avito-internship-assignment/internal/delivery/http/handler"
	"github.com/vit6556/avito-internship-assignment/internal/service"
	"github.com/vit6556/avito-internship-assignment/internal/service/mock"
)

func TestSendCoin(t *testing.T) {
	e := echo.New()
	mockTransactionService := new(mock.MockTransactionService)
	handler := httphandler.NewTransactionHandler(mockTransactionService)

	tests := []struct {
		name           string
		userID         interface{}
		requestBody    dto.SendCoinRequest
		contentType    string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Success - Coins sent successfully",
			userID: 1,
			requestBody: dto.SendCoinRequest{
				ToUser: "bob",
				Amount: 50,
			},
			contentType: "application/json",
			mockSetup: func() {
				mockTransactionService.On("SendCoins", testifyMock.Anything, 1, "bob", 50).
					Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"coins sent successfully"}`,
		},
		{
			name:           "Error - Unauthorized",
			userID:         nil,
			requestBody:    dto.SendCoinRequest{ToUser: "bob", Amount: 50},
			contentType:    "application/json",
			mockSetup:      func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"unauthorized"}`,
		},
		{
			name:           "Error - Invalid Content-Type",
			userID:         1,
			requestBody:    dto.SendCoinRequest{ToUser: "bob", Amount: 50},
			contentType:    "text/plain",
			mockSetup:      func() {},
			expectedStatus: http.StatusUnsupportedMediaType,
			expectedBody:   `{"error":"Content-Type must be application/json"}`,
		},
		{
			name:           "Error - Invalid request data",
			userID:         1,
			requestBody:    dto.SendCoinRequest{},
			contentType:    "application/json",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid request data"}`,
		},
		{
			name:   "Error - Employee Not Found",
			userID: 1,
			requestBody: dto.SendCoinRequest{
				ToUser: "bob",
				Amount: 50,
			},
			contentType: "application/json",
			mockSetup: func() {
				mockTransactionService.On("SendCoins", testifyMock.Anything, 1, "bob", 50).
					Return(service.ErrEmployeeNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"employee not found"}`,
		},
		{
			name:   "Error - Insufficient Funds",
			userID: 1,
			requestBody: dto.SendCoinRequest{
				ToUser: "bob",
				Amount: 50,
			},
			contentType: "application/json",
			mockSetup: func() {
				mockTransactionService.On("SendCoins", testifyMock.Anything, 1, "bob", 50).
					Return(service.ErrInsufficientFunds).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"insufficient funds"}`,
		},
		{
			name:   "Error - Self Transaction",
			userID: 1,
			requestBody: dto.SendCoinRequest{
				ToUser: "alice",
				Amount: 50,
			},
			contentType: "application/json",
			mockSetup: func() {
				mockTransactionService.On("SendCoins", testifyMock.Anything, 1, "alice", 50).
					Return(service.ErrSelfTransaction).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"sender and recipient cannot be the same user"}`,
		},
		{
			name:   "Error - Internal Server Error",
			userID: 1,
			requestBody: dto.SendCoinRequest{
				ToUser: "bob",
				Amount: 50,
			},
			contentType: "application/json",
			mockSetup: func() {
				mockTransactionService.On("SendCoins", testifyMock.Anything, 1, "bob", 50).
					Return(errors.New("unexpected error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"failed to send coins"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
			req.Header.Set("Content-Type", tt.contentType)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.userID != nil {
				c.Set("userID", tt.userID)
			}

			err := handler.SendCoin(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())

			mockTransactionService.AssertExpectations(t)
		})
	}
}

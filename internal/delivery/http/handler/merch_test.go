package httphandler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"

	"github.com/vit6556/avito-internship-assignment/internal/delivery/http/handler"
	"github.com/vit6556/avito-internship-assignment/internal/service"
	"github.com/vit6556/avito-internship-assignment/internal/service/mock"
)

func TestBuyItem(t *testing.T) {
	e := echo.New()
	mockMerchService := new(mock.MockMerchService)
	handler := httphandler.NewMerchHandler(mockMerchService)

	tests := []struct {
		name           string
		userID         interface{}
		itemName       string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "Success - Merch purchased successfully",
			userID:   1,
			itemName: "book",
			mockSetup: func() {
				mockMerchService.On("BuyItem", testifyMock.Anything, 1, "book").
					Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"merch purchased successfully"}`,
		},
		{
			name:           "Error - Unauthorized",
			userID:         nil,
			itemName:       "book",
			mockSetup:      func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"unauthorized"}`,
		},
		{
			name:           "Error - Invalid merch name",
			userID:         1,
			itemName:       "",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid merch name"}`,
		},
		{
			name:     "Error - Merch Not Found",
			userID:   1,
			itemName: "book",
			mockSetup: func() {
				mockMerchService.On("BuyItem", testifyMock.Anything, 1, "book").
					Return(service.ErrMerchNotFound).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"merch not found"}`,
		},
		{
			name:     "Error - Insufficient Funds",
			userID:   1,
			itemName: "book",
			mockSetup: func() {
				mockMerchService.On("BuyItem", testifyMock.Anything, 1, "book").
					Return(service.ErrInsufficientFunds).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"insufficient funds"}`,
		},
		{
			name:     "Error - Internal Server Error",
			userID:   1,
			itemName: "book",
			mockSetup: func() {
				mockMerchService.On("BuyItem", testifyMock.Anything, 1, "book").
					Return(errors.New("unexpected error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"failed to buy merch"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodGet, "/api/buy/"+tt.itemName, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.userID != nil {
				c.Set("userID", tt.userID)
			}
			c.SetParamNames("item")
			c.SetParamValues(tt.itemName)

			err := handler.BuyItem(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())

			mockMerchService.AssertExpectations(t)
		})
	}
}

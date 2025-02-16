package httphandler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
	"github.com/vit6556/avito-internship-assignment/internal/delivery/http/dto"
	httphandler "github.com/vit6556/avito-internship-assignment/internal/delivery/http/handler"
	"github.com/vit6556/avito-internship-assignment/internal/service/mock"
)

func TestGetEmployeeInfo(t *testing.T) {
	e := echo.New()
	mockEmployeeService := new(mock.MockEmployeeService)
	employeeHandler := httphandler.NewEmployeeHandler(mockEmployeeService)

	tests := []struct {
		name           string
		userID         interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "Success - Fetch Employee Info",
			userID: 1,
			mockSetup: func() {
				mockEmployeeService.On("GetEmployeeInfo", testifyMock.Anything, 1).
					Return(&dto.EmployeeInfoResponse{
						Coins: 100,
						Inventory: []*dto.InventoryItem{
							{Type: "book", Quantity: 1},
						},
						CoinHistory: &dto.CoinHistory{
							Received: []dto.CoinTransaction{
								{User: "bob", Amount: 50},
							},
							Sent: []dto.CoinTransaction{
								{User: "alice", Amount: 30},
							},
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"coins": float64(100),
				"inventory": []interface{}{
					map[string]interface{}{"type": "book", "quantity": float64(1)},
				},
				"coinHistory": map[string]interface{}{
					"received": []interface{}{
						map[string]interface{}{"user": "bob", "amount": float64(50)},
					},
					"sent": []interface{}{
						map[string]interface{}{"user": "alice", "amount": float64(30)},
					},
				},
			},
		},
		{
			name:           "Error - Unauthorized User",
			userID:         nil,
			mockSetup:      func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   map[string]interface{}{"error": "unauthorized"},
		},
		{
			name:   "Error - Internal Server Error",
			userID: 2,
			mockSetup: func() {
				mockEmployeeService.On("GetEmployeeInfo", testifyMock.Anything, 2).
					Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "failed to fetch employee info"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("userID", tt.userID)

			err := employeeHandler.GetEmployeeInfo(c)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			var responseBody map[string]interface{}
			err = json.Unmarshal(rec.Body.Bytes(), &responseBody)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, responseBody)

			mockEmployeeService.AssertExpectations(t)
		})
	}
}

package employeeservice_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vit6556/avito-internship-assignment/internal/database/mock"
	"github.com/vit6556/avito-internship-assignment/internal/delivery/http/dto"
	"github.com/vit6556/avito-internship-assignment/internal/entity"
	"github.com/vit6556/avito-internship-assignment/internal/service"
	"github.com/vit6556/avito-internship-assignment/internal/service/employee"
)

func TestGetEmployeeInfo(t *testing.T) {
	ctx := context.Background()
	mockEmployeeRepo := new(mock.MockEmployeeRepository)
	mockMerchRepo := new(mock.MockMerchRepository)
	mockTransactionRepo := new(mock.MockTransactionRepository)
	employeeService := employeeservice.NewEmployeeService(mockEmployeeRepo, mockMerchRepo, mockTransactionRepo)

	tests := []struct {
		name          string
		user          *entity.Employee
		coinHistory   *entity.CoinHistory
		inventory     []*entity.InventoryItem
		mockSetup     func()
		expectedError error
		expectedData  *dto.EmployeeInfoResponse
	}{
		{
			name: "Success - Get Employee Info",
			user: &entity.Employee{ID: 1, Username: "alice", Balance: 100},
			coinHistory: &entity.CoinHistory{
				Received: []entity.CoinTransaction{{User: "bob", Amount: 50}},
				Sent:     []entity.CoinTransaction{{User: "charlie", Amount: 20}},
			},
			inventory: []*entity.InventoryItem{
				{Type: "book", Quantity: 2},
				{Type: "powerbank", Quantity: 1},
			},
			mockSetup: func() {
				mockEmployeeRepo.ExpectedCalls = nil
				mockMerchRepo.ExpectedCalls = nil
				mockTransactionRepo.ExpectedCalls = nil

				mockEmployeeRepo.On("GetEmployeeByID", ctx, 1).
					Return(&entity.Employee{ID: 1, Username: "alice", Balance: 100}, nil)
				mockTransactionRepo.On("GetCoinHistory", ctx, 1).
					Return(&entity.CoinHistory{
						Received: []entity.CoinTransaction{{User: "bob", Amount: 50}},
						Sent:     []entity.CoinTransaction{{User: "charlie", Amount: 20}},
					}, nil)
				mockMerchRepo.On("GetUserPurchases", ctx, 1).
					Return([]*entity.InventoryItem{
						{Type: "book", Quantity: 2},
						{Type: "powerbank", Quantity: 1},
					}, nil)
			},
			expectedError: nil,
			expectedData: &dto.EmployeeInfoResponse{
				Coins: 100,
				Inventory: []*dto.InventoryItem{
					{Type: "book", Quantity: 2},
					{Type: "powerbank", Quantity: 1},
				},
				CoinHistory: &dto.CoinHistory{
					Received: []dto.CoinTransaction{{User: "bob", Amount: 50}},
					Sent:     []dto.CoinTransaction{{User: "charlie", Amount: 20}},
				},
			},
		},
		{
			name:        "Error - Employee Not Found",
			user:        nil,
			coinHistory: nil,
			inventory:   nil,
			mockSetup: func() {
				mockEmployeeRepo.ExpectedCalls = nil
				mockMerchRepo.ExpectedCalls = nil
				mockTransactionRepo.ExpectedCalls = nil

				mockEmployeeRepo.On("GetEmployeeByID", ctx, 1).
					Return(nil, service.ErrEmployeeNotFound)
			},
			expectedError: service.ErrEmployeeNotFound,
			expectedData:  nil,
		},
		{
			name:        "Error - Failed to Fetch Coin History",
			user:        &entity.Employee{ID: 1, Username: "alice", Balance: 100},
			coinHistory: nil,
			inventory:   []*entity.InventoryItem{},
			mockSetup: func() {
				mockEmployeeRepo.ExpectedCalls = nil
				mockMerchRepo.ExpectedCalls = nil
				mockTransactionRepo.ExpectedCalls = nil

				mockEmployeeRepo.On("GetEmployeeByID", ctx, 1).
					Return(&entity.Employee{ID: 1, Username: "alice", Balance: 100}, nil)
				mockTransactionRepo.On("GetCoinHistory", ctx, 1).
					Return(nil, service.ErrDatabaseError)
			},
			expectedError: service.ErrDatabaseError,
			expectedData:  nil,
		},
		{
			name: "Error - Failed to Fetch Inventory",
			user: &entity.Employee{ID: 1, Username: "alice", Balance: 100},
			coinHistory: &entity.CoinHistory{
				Received: []entity.CoinTransaction{{User: "bob", Amount: 50}},
				Sent:     []entity.CoinTransaction{{User: "charlie", Amount: 20}},
			},
			inventory: nil,
			mockSetup: func() {
				mockEmployeeRepo.ExpectedCalls = nil
				mockMerchRepo.ExpectedCalls = nil
				mockTransactionRepo.ExpectedCalls = nil

				mockEmployeeRepo.On("GetEmployeeByID", ctx, 1).
					Return(&entity.Employee{ID: 1, Username: "alice", Balance: 100}, nil)
				mockTransactionRepo.On("GetCoinHistory", ctx, 1).
					Return(&entity.CoinHistory{
						Received: []entity.CoinTransaction{{User: "bob", Amount: 50}},
						Sent:     []entity.CoinTransaction{{User: "charlie", Amount: 20}},
					}, nil)
				mockMerchRepo.On("GetUserPurchases", ctx, 1).
					Return(nil, service.ErrDatabaseError)
			},
			expectedError: service.ErrDatabaseError,
			expectedData:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			var userID int
			if tt.user != nil {
				userID = tt.user.ID
			} else {
				userID = 1
			}

			result, err := employeeService.GetEmployeeInfo(ctx, userID)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedData, result)

			mockEmployeeRepo.AssertExpectations(t)
			mockMerchRepo.AssertExpectations(t)
			mockTransactionRepo.AssertExpectations(t)
		})
	}
}

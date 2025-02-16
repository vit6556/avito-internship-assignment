package merchservice_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vit6556/avito-internship-assignment/internal/database"
	"github.com/vit6556/avito-internship-assignment/internal/database/mock"
	"github.com/vit6556/avito-internship-assignment/internal/entity"
	"github.com/vit6556/avito-internship-assignment/internal/service"
	"github.com/vit6556/avito-internship-assignment/internal/service/merch"
)

func TestBuyItem(t *testing.T) {
	ctx := context.Background()
	mockEmployeeRepo := new(mock.MockEmployeeRepository)
	mockMerchRepo := new(mock.MockMerchRepository)
	merchService := merchservice.NewMerchService(mockEmployeeRepo, mockMerchRepo)

	tests := []struct {
		name          string
		user          *entity.Employee
		item          *entity.MerchItem
		mockSetup     func()
		expectedError error
	}{
		{
			name: "Success - Purchase Item",
			user: &entity.Employee{ID: 1, Username: "alice", Balance: 100},
			item: &entity.MerchItem{ID: 1, Name: "book", Price: 50},
			mockSetup: func() {
				mockEmployeeRepo.ExpectedCalls = nil
				mockMerchRepo.ExpectedCalls = nil

				mockMerchRepo.On("GetItemByName", ctx, "book").
					Return(&entity.MerchItem{ID: 1, Name: "book", Price: 50}, nil)
				mockEmployeeRepo.On("GetEmployeeByID", ctx, 1).
					Return(&entity.Employee{ID: 1, Username: "alice", Balance: 100}, nil)
				mockMerchRepo.On("BuyItem", ctx, 1, 1).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Error - Item Not Found",
			user: &entity.Employee{ID: 1, Username: "alice", Balance: 100},
			item: nil,
			mockSetup: func() {
				mockEmployeeRepo.ExpectedCalls = nil
				mockMerchRepo.ExpectedCalls = nil

				mockMerchRepo.On("GetItemByName", ctx, "book").
					Return(nil, database.ErrMerchNotFound)
			},
			expectedError: service.ErrMerchNotFound,
		},
		{
			name: "Error - Employee Not Found",
			user: nil,
			item: &entity.MerchItem{ID: 1, Name: "book", Price: 50},
			mockSetup: func() {
				mockEmployeeRepo.ExpectedCalls = nil
				mockMerchRepo.ExpectedCalls = nil

				mockMerchRepo.On("GetItemByName", ctx, "book").
					Return(&entity.MerchItem{ID: 1, Name: "book", Price: 50}, nil)
				mockEmployeeRepo.On("GetEmployeeByID", ctx, 1).
					Return(nil, database.ErrEmployeeNotFound)
			},
			expectedError: service.ErrEmployeeNotFound,
		},
		{
			name: "Error - Insufficient Funds",
			user: &entity.Employee{ID: 1, Username: "alice", Balance: 20},
			item: &entity.MerchItem{ID: 1, Name: "book", Price: 50},
			mockSetup: func() {
				mockEmployeeRepo.ExpectedCalls = nil
				mockMerchRepo.ExpectedCalls = nil

				mockMerchRepo.On("GetItemByName", ctx, "book").
					Return(&entity.MerchItem{ID: 1, Name: "book", Price: 50}, nil)
				mockEmployeeRepo.On("GetEmployeeByID", ctx, 1).
					Return(&entity.Employee{ID: 1, Username: "alice", Balance: 20}, nil)
			},
			expectedError: service.ErrInsufficientFunds,
		},
		{
			name: "Error - Database Error on Purchase",
			user: &entity.Employee{ID: 1, Username: "alice", Balance: 100},
			item: &entity.MerchItem{ID: 1, Name: "book", Price: 50},
			mockSetup: func() {
				mockEmployeeRepo.ExpectedCalls = nil
				mockMerchRepo.ExpectedCalls = nil

				mockMerchRepo.On("GetItemByName", ctx, "book").
					Return(&entity.MerchItem{ID: 1, Name: "book", Price: 50}, nil)
				mockEmployeeRepo.On("GetEmployeeByID", ctx, 1).
					Return(&entity.Employee{ID: 1, Username: "alice", Balance: 100}, nil)
				mockMerchRepo.On("BuyItem", ctx, 1, 1).
					Return(database.ErrDatabaseTransaction)
			},
			expectedError: service.ErrDatabaseError,
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

			var itemName string
			if tt.item != nil {
				itemName = tt.item.Name
			} else {
				itemName = "book"
			}

			err := merchService.BuyItem(ctx, userID, itemName)

			assert.Equal(t, tt.expectedError, err)

			mockEmployeeRepo.AssertExpectations(t)
			mockMerchRepo.AssertExpectations(t)
		})
	}
}

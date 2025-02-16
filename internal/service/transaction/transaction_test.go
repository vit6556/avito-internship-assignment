package transactionservice_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vit6556/avito-internship-assignment/internal/database"
	"github.com/vit6556/avito-internship-assignment/internal/database/mock"
	"github.com/vit6556/avito-internship-assignment/internal/entity"
	"github.com/vit6556/avito-internship-assignment/internal/service"
	"github.com/vit6556/avito-internship-assignment/internal/service/transaction"
)

func TestSendCoins(t *testing.T) {
	ctx := context.Background()
	mockEmployeeRepo := new(mock.MockEmployeeRepository)
	mockTransactionRepo := new(mock.MockTransactionRepository)
	transactionService := transactionservice.NewTransactionService(mockEmployeeRepo, mockTransactionRepo)

	tests := []struct {
		name          string
		sender        *entity.Employee
		receiver      *entity.Employee
		amount        int
		mockSetup     func()
		expectedError error
	}{
		{
			name:     "Success - Valid transaction",
			sender:   &entity.Employee{ID: 1, Username: "alice"},
			receiver: &entity.Employee{ID: 2, Username: "bob"},
			amount:   50,
			mockSetup: func() {
				mockEmployeeRepo.ExpectedCalls = nil
				mockTransactionRepo.ExpectedCalls = nil

				mockEmployeeRepo.On("GetEmployeeByUsername", ctx, "bob").
					Return(&entity.Employee{ID: 2, Username: "bob"}, nil)
				mockEmployeeRepo.On("GetEmployeeByID", ctx, 1).
					Return(&entity.Employee{ID: 1, Username: "alice", Balance: 100}, nil)
				mockTransactionRepo.On("SendCoins", ctx, 1, 2, 50).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "Error - Receiver Not Found",
			sender:   &entity.Employee{ID: 1, Username: "alice", Balance: 100},
			receiver: nil,
			amount:   50,
			mockSetup: func() {
				mockEmployeeRepo.ExpectedCalls = nil
				mockTransactionRepo.ExpectedCalls = nil

				mockEmployeeRepo.On("GetEmployeeByUsername", ctx, "bob").
					Return(nil, database.ErrEmployeeNotFound)
			},
			expectedError: service.ErrEmployeeNotFound,
		},
		{
			name:     "Error - Sender Not Found",
			sender:   nil,
			receiver: &entity.Employee{ID: 2, Username: "bob"},
			amount:   50,
			mockSetup: func() {
				mockEmployeeRepo.ExpectedCalls = nil
				mockTransactionRepo.ExpectedCalls = nil

				mockEmployeeRepo.On("GetEmployeeByUsername", ctx, "bob").
					Return(&entity.Employee{ID: 2, Username: "bob"}, nil)
				mockEmployeeRepo.On("GetEmployeeByID", ctx, 1).
					Return(nil, database.ErrEmployeeNotFound)
			},
			expectedError: service.ErrEmployeeNotFound,
		},
		{
			name:     "Error - Sender and Receiver are the Same",
			sender:   &entity.Employee{ID: 2, Username: "bob"},
			receiver: &entity.Employee{ID: 2, Username: "bob"},
			amount:   50,
			mockSetup: func() {
				mockEmployeeRepo.ExpectedCalls = nil
				mockTransactionRepo.ExpectedCalls = nil

				mockEmployeeRepo.On("GetEmployeeByUsername", ctx, "bob").
					Return(&entity.Employee{ID: 2, Username: "bob"}, nil)
				mockEmployeeRepo.On("GetEmployeeByID", ctx, 2).
					Return(&entity.Employee{ID: 2, Username: "bob"}, nil)
			},
			expectedError: service.ErrSelfTransaction,
		},
		{
			name:     "Error - Insufficient Funds",
			sender:   &entity.Employee{ID: 1, Username: "alice", Balance: 30},
			receiver: &entity.Employee{ID: 2, Username: "bob"},
			amount:   50,
			mockSetup: func() {
				mockEmployeeRepo.ExpectedCalls = nil
				mockTransactionRepo.ExpectedCalls = nil

				mockEmployeeRepo.On("GetEmployeeByUsername", ctx, "bob").
					Return(&entity.Employee{ID: 2, Username: "bob"}, nil)
				mockEmployeeRepo.On("GetEmployeeByID", ctx, 1).
					Return(&entity.Employee{ID: 1, Username: "alice", Balance: 30}, nil)
			},
			expectedError: service.ErrInsufficientFunds,
		},
		{
			name:     "Error - Transaction Failed",
			sender:   &entity.Employee{ID: 1, Username: "alice", Balance: 100},
			receiver: &entity.Employee{ID: 2, Username: "bob"},
			amount:   50,
			mockSetup: func() {
				mockEmployeeRepo.ExpectedCalls = nil
				mockTransactionRepo.ExpectedCalls = nil

				mockEmployeeRepo.On("GetEmployeeByUsername", ctx, "bob").
					Return(&entity.Employee{ID: 2, Username: "bob"}, nil)
				mockEmployeeRepo.On("GetEmployeeByID", ctx, 1).
					Return(&entity.Employee{ID: 1, Username: "alice", Balance: 100}, nil)
				mockTransactionRepo.On("SendCoins", ctx, 1, 2, 50).
					Return(database.ErrDatabaseTransaction)
			},
			expectedError: service.ErrDatabaseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			var senderID int
			if tt.sender != nil {
				senderID = tt.sender.ID
			} else {
				senderID = 1
			}

			var receiverUsername string
			if tt.receiver != nil {
				receiverUsername = tt.receiver.Username
			} else {
				receiverUsername = "bob"
			}

			err := transactionService.SendCoins(ctx, senderID, receiverUsername, tt.amount)

			assert.Equal(t, tt.expectedError, err)

			mockEmployeeRepo.AssertExpectations(t)
			mockTransactionRepo.AssertExpectations(t)
		})
	}
}

package transactionservice

import (
	"context"
	"log"

	"github.com/vit6556/avito-internship-assignment/internal/database"
	"github.com/vit6556/avito-internship-assignment/internal/service"
)

type TransactionService struct {
	employeeRepo    database.EmployeeRepository
	transactionRepo database.TransactionRepository
}

func NewTransactionService(employeeRepo database.EmployeeRepository, transactionRepo database.TransactionRepository) *TransactionService {
	return &TransactionService{
		employeeRepo:    employeeRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *TransactionService) SendCoins(ctx context.Context, senderID int, toUser string, amount int) error {
	receiver, err := s.employeeRepo.GetEmployeeByUsername(ctx, toUser)
	if err != nil {
		log.Printf("recipient %q not found: %v", toUser, err)
		return service.ErrEmployeeNotFound
	}

	sender, err := s.employeeRepo.GetEmployeeByID(ctx, senderID)
	if err != nil {
		log.Printf("sender %q not found: %v", toUser, err)
		return service.ErrEmployeeNotFound
	}

	if sender.ID == receiver.ID {
		return service.ErrSelfTransaction
	}

	if sender.Balance < amount {
		return service.ErrInsufficientFunds
	}

	err = s.transactionRepo.SendCoins(ctx, senderID, receiver.ID, amount)
	if err != nil {
		log.Printf("transaction failed: user %d -> %d, amount: %d, error: %v", senderID, receiver.ID, amount, err)
		switch err {
		case database.ErrEmployeeNotFound:
			return service.ErrEmployeeNotFound
		case database.ErrInsufficientFunds:
			return service.ErrInsufficientFunds
		default:
			return service.ErrDatabaseError
		}
	}

	return nil
}

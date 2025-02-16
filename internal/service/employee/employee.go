package employeeservice

import (
	"context"
	"log"

	"github.com/vit6556/avito-internship-assignment/internal/database"
	"github.com/vit6556/avito-internship-assignment/internal/delivery/http/dto"
	"github.com/vit6556/avito-internship-assignment/internal/entity"
	"github.com/vit6556/avito-internship-assignment/internal/service"
)

type EmployeeService struct {
	employeeRepo    database.EmployeeRepository
	merchRepo       database.MerchRepository
	transactionRepo database.TransactionRepository
}

func NewEmployeeService(employeeRepo database.EmployeeRepository, merchRepo database.MerchRepository, transactionRepo database.TransactionRepository) *EmployeeService {
	return &EmployeeService{
		employeeRepo:    employeeRepo,
		merchRepo:       merchRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *EmployeeService) GetEmployeeInfo(ctx context.Context, id int) (*dto.EmployeeInfoResponse, error) {
	employee, err := s.employeeRepo.GetEmployeeByID(ctx, id)
	if err != nil {
		log.Printf("failed to get employee by ID %q: %v", id, err)
		return nil, service.ErrEmployeeNotFound
	}

	coinHistory, err := s.transactionRepo.GetCoinHistory(ctx, id)
	if err != nil {
		log.Printf("failed to get coin history for user %d: %v", id, err)
		return nil, service.ErrDatabaseError
	}

	purchases, err := s.merchRepo.GetUserPurchases(ctx, id)
	if err != nil {
		log.Printf("failed to get purchases for user %d: %v", id, err)
		return nil, service.ErrDatabaseError
	}

	return &dto.EmployeeInfoResponse{
		Coins:       employee.Balance,
		Inventory:   mapInventoryToDTO(purchases),
		CoinHistory: mapCoinHistoryToDTO(coinHistory),
	}, nil
}

func mapInventoryToDTO(items []*entity.InventoryItem) []*dto.InventoryItem {
	purchases := make([]*dto.InventoryItem, len(items))

	for i, item := range items {
		purchases[i] = &dto.InventoryItem{
			Type:     item.Type,
			Quantity: item.Quantity,
		}
	}

	return purchases
}

func mapCoinHistoryToDTO(history *entity.CoinHistory) *dto.CoinHistory {
	if history == nil {
		return nil
	}

	received := make([]dto.CoinTransaction, len(history.Received))
	sent := make([]dto.CoinTransaction, len(history.Sent))

	for i := range history.Received {
		received[i] = dto.CoinTransaction{
			User:   history.Received[i].User,
			Amount: history.Received[i].Amount,
		}
	}

	for i := range history.Sent {
		sent[i] = dto.CoinTransaction{
			User:   history.Sent[i].User,
			Amount: history.Sent[i].Amount,
		}
	}

	return &dto.CoinHistory{
		Received: received,
		Sent:     sent,
	}
}

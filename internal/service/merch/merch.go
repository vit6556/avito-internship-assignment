package merchservice

import (
	"context"
	"log"

	"github.com/vit6556/avito-internship-assignment/internal/database"
	"github.com/vit6556/avito-internship-assignment/internal/service"
)

type MerchService struct {
	employeeRepo database.EmployeeRepository
	merchRepo    database.MerchRepository
}

func NewMerchService(employeeRepo database.EmployeeRepository, merchRepo database.MerchRepository) *MerchService {
	return &MerchService{
		employeeRepo: employeeRepo,
		merchRepo:    merchRepo,
	}
}

func (s *MerchService) BuyItem(ctx context.Context, userID int, itemName string) error {
	item, err := s.merchRepo.GetItemByName(ctx, itemName)
	if err != nil {
		log.Printf("merch %q not found: %v", itemName, err)
		return service.ErrMerchNotFound
	}

	user, err := s.employeeRepo.GetEmployeeByID(ctx, userID)
	if err != nil {
		log.Printf("employee %d not found: %v", userID, err)
		return service.ErrEmployeeNotFound
	}

	if user.Balance < item.Price {
		return service.ErrInsufficientFunds
	}

	err = s.merchRepo.BuyItem(ctx, userID, item.ID)
	if err != nil {
		log.Printf("failed to process purchase for employee %d and item %q: %v", userID, itemName, err)
		switch err {
		case database.ErrMerchNotFound:
			return service.ErrMerchNotFound
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

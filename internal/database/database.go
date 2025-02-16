package database

import (
	"context"
	"errors"

	"github.com/vit6556/avito-internship-assignment/internal/entity"
)

var (
	ErrEmployeeNotFound       = errors.New("employee not found")
	ErrEmployeeCreationFailed = errors.New("failed to create Employee")

	ErrMerchNotFound     = errors.New("merch not found")
	ErrInsufficientFunds = errors.New("insufficient funds")

	ErrDatabaseQueryFailed  = errors.New("database query failed")
	ErrDatabaseScanFailed   = errors.New("failed to scan database row")
	ErrDatabaseTransaction  = errors.New("database transaction failed")
	ErrDatabaseInsertFailed = errors.New("failed to insert data into database")
	ErrDatabaseUpdateFailed = errors.New("failed to update database record")
)

type EmployeeRepository interface {
	CreateEmployee(ctx context.Context, employee entity.Employee) (int, error)
	GetEmployeeByID(ctx context.Context, userID int) (*entity.Employee, error)
	GetEmployeeByUsername(ctx context.Context, username string) (*entity.Employee, error)
}

type MerchRepository interface {
	BuyItem(ctx context.Context, userID int, itemID int) error
	GetItemByID(ctx context.Context, itemID int) (*entity.MerchItem, error)
	GetItemByName(ctx context.Context, name string) (*entity.MerchItem, error)
	GetUserPurchases(ctx context.Context, userID int) ([]*entity.InventoryItem, error)
}

type TransactionRepository interface {
	GetCoinHistory(ctx context.Context, userID int) (*entity.CoinHistory, error)
	SendCoins(ctx context.Context, senderID, receiverID, amount int) error
}

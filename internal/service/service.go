package service

import (
	"context"
	"errors"

	"github.com/vit6556/avito-internship-assignment/internal/delivery/http/dto"
)

var (
	ErrInvalidCredentials   = errors.New("invalid username or password")
	ErrInvalidToken         = errors.New("invalid or expired token")
	ErrAuthenticationFailed = errors.New("authentication failed")

	ErrDatabaseError = errors.New("database operation failed")

	ErrMerchNotFound = errors.New("merch not found")

	ErrEmployeeCreationFailed = errors.New("failed to create employee")
	ErrEmployeeNotFound       = errors.New("employee not found")

	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrSelfTransaction   = errors.New("sender and receiver cannot be the same user")
)

type AuthService interface {
	AuthorizeUser(ctx context.Context, username, password string) (string, error)
	ValidateToken(tokenString string) (int, error)
}

type EmployeeService interface {
	GetEmployeeInfo(ctx context.Context, id int) (*dto.EmployeeInfoResponse, error)
}

type MerchService interface {
	BuyItem(ctx context.Context, userID int, itemName string) error
}

type TransactionService interface {
	SendCoins(ctx context.Context, senderID int, toUser string, amount int) error
}

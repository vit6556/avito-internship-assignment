package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/vit6556/avito-internship-assignment/internal/entity"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) GetCoinHistory(ctx context.Context, userID int) (*entity.CoinHistory, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.CoinHistory), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTransactionRepository) SendCoins(ctx context.Context, senderID, receiverID, amount int) error {
	args := m.Called(ctx, senderID, receiverID, amount)
	return args.Error(0)
}

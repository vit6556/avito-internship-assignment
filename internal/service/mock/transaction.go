package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) SendCoins(ctx context.Context, senderID int, toUser string, amount int) error {
	args := m.Called(ctx, senderID, toUser, amount)
	return args.Error(0)
}

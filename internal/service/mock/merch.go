package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockMerchService struct {
	mock.Mock
}

func (m *MockMerchService) BuyItem(ctx context.Context, userID int, itemName string) error {
	args := m.Called(ctx, userID, itemName)
	return args.Error(0)
}

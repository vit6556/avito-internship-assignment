package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/vit6556/avito-internship-assignment/internal/entity"
)

type MockMerchRepository struct {
	mock.Mock
}

func (m *MockMerchRepository) BuyItem(ctx context.Context, userID int, itemID int) error {
	args := m.Called(ctx, userID, itemID)
	return args.Error(0)
}

func (m *MockMerchRepository) GetItemByID(ctx context.Context, itemID int) (*entity.MerchItem, error) {
	args := m.Called(ctx, itemID)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.MerchItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMerchRepository) GetItemByName(ctx context.Context, name string) (*entity.MerchItem, error) {
	args := m.Called(ctx, name)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.MerchItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMerchRepository) GetUserPurchases(ctx context.Context, userID int) ([]*entity.InventoryItem, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) != nil {
		return args.Get(0).([]*entity.InventoryItem), args.Error(1)
	}
	return nil, args.Error(1)
}

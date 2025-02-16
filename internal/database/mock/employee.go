package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/vit6556/avito-internship-assignment/internal/entity"
)

type MockEmployeeRepository struct {
	mock.Mock
}

func (m *MockEmployeeRepository) CreateEmployee(ctx context.Context, employee entity.Employee) (int, error) {
	args := m.Called(ctx, employee)
	return args.Int(0), args.Error(1)
}

func (m *MockEmployeeRepository) GetEmployeeByID(ctx context.Context, userID int) (*entity.Employee, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Employee), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockEmployeeRepository) GetEmployeeByUsername(ctx context.Context, username string) (*entity.Employee, error) {
	args := m.Called(ctx, username)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Employee), args.Error(1)
	}
	return nil, args.Error(1)
}

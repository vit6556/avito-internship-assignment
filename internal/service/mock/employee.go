package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/vit6556/avito-internship-assignment/internal/delivery/http/dto"
)

type MockEmployeeService struct {
	mock.Mock
}

func (m *MockEmployeeService) GetEmployeeInfo(ctx context.Context, id int) (*dto.EmployeeInfoResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.EmployeeInfoResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

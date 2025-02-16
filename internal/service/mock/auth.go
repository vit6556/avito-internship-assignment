package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) AuthorizeUser(ctx context.Context, username, password string) (string, error) {
	args := m.Called(ctx, username, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidateToken(tokenString string) (int, error) {
	args := m.Called(tokenString)
	return args.Int(0), args.Error(1)
}

package service

import (
	"context"

	"github.com/haqer0002/avito-shop/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository мок для репозитория пользователей
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateCoins(ctx context.Context, userID int64, amount int64) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}

// MockTransactionRepository мок для репозитория транзакций
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, transaction *models.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetUserTransactions(ctx context.Context, userID int64) ([]models.Transaction, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

// MockUserMerchRepository мок для репозитория купленного мерча
type MockUserMerchRepository struct {
	mock.Mock
}

func (m *MockUserMerchRepository) Create(ctx context.Context, userMerch *models.UserMerch) error {
	args := m.Called(ctx, userMerch)
	return args.Error(0)
}

func (m *MockUserMerchRepository) GetUserMerch(ctx context.Context, userID int64) ([]models.UserMerch, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.UserMerch), args.Error(1)
}

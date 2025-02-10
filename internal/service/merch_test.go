package service

import (
	"context"
	"errors"
	"testing"

	"github.com/haqer0002/avito-shop/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMerchRepository мок для репозитория мерча
type MockMerchRepository struct {
	mock.Mock
}

func (m *MockMerchRepository) GetByName(ctx context.Context, name string) (*models.MerchItem, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MerchItem), args.Error(1)
}

func (m *MockMerchRepository) GetAll(ctx context.Context) ([]models.MerchItem, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.MerchItem), args.Error(1)
}

func TestMerchService_BuyMerch(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockMerchRepo := new(MockMerchRepository)
	mockUserMerchRepo := new(MockUserMerchRepository)

	service := NewMerchService(mockUserRepo, mockMerchRepo, mockUserMerchRepo)

	ctx := context.Background()
	userID := int64(1)
	merchName := "t-shirt"

	// Создаем тестовый мерч
	testMerch := &models.MerchItem{
		ID:    1,
		Name:  merchName,
		Price: 80,
	}

	// Настраиваем моки
	mockMerchRepo.On("GetByName", ctx, merchName).Return(testMerch, nil)
	mockUserRepo.On("UpdateCoins", ctx, userID, -testMerch.Price).Return(nil)
	mockUserMerchRepo.On("Create", ctx, mock.AnythingOfType("*models.UserMerch")).Return(nil)

	// Вызываем тестируемый метод
	err := service.BuyMerch(ctx, userID, merchName)

	// Проверяем результаты
	assert.NoError(t, err)
	mockMerchRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockUserMerchRepo.AssertExpectations(t)
}

func TestMerchService_BuyMerch_InsufficientFunds(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockMerchRepo := new(MockMerchRepository)
	mockUserMerchRepo := new(MockUserMerchRepository)

	service := NewMerchService(mockUserRepo, mockMerchRepo, mockUserMerchRepo)

	ctx := context.Background()
	userID := int64(1)
	merchName := "t-shirt"

	// Создаем тестовый мерч
	testMerch := &models.MerchItem{
		ID:    1,
		Name:  merchName,
		Price: 80,
	}

	// Настраиваем моки
	mockMerchRepo.On("GetByName", ctx, merchName).Return(testMerch, nil)
	mockUserRepo.On("UpdateCoins", ctx, userID, -testMerch.Price).Return(errors.New("insufficient funds"))

	// Вызываем тестируемый метод
	err := service.BuyMerch(ctx, userID, merchName)

	// Проверяем результаты
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to deduct coins")
	mockMerchRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

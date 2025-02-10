package service

import (
	"context"
	"errors"
	"testing"

	"github.com/haqer0002/avito-shop/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_SendCoins(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTransactionRepo := new(MockTransactionRepository)
	mockUserMerchRepo := new(MockUserMerchRepository)

	service := NewUserService(mockUserRepo, mockTransactionRepo, mockUserMerchRepo)

	ctx := context.Background()
	fromUserID := int64(1)
	toUsername := "recipient"
	amount := int64(100)

	// Создаем тестового получателя
	recipient := &models.User{
		ID:       2,
		Username: toUsername,
		Coins:    1000,
	}

	// Настраиваем моки
	mockUserRepo.On("GetByUsername", ctx, toUsername).Return(recipient, nil)
	mockUserRepo.On("UpdateCoins", ctx, fromUserID, -amount).Return(nil)
	mockUserRepo.On("UpdateCoins", ctx, recipient.ID, amount).Return(nil)
	mockTransactionRepo.On("Create", ctx, mock.AnythingOfType("*models.Transaction")).Return(nil)

	// Вызываем тестируемый метод
	err := service.SendCoins(ctx, fromUserID, toUsername, amount)

	// Проверяем результаты
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestUserService_SendCoins_InsufficientFunds(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTransactionRepo := new(MockTransactionRepository)
	mockUserMerchRepo := new(MockUserMerchRepository)

	service := NewUserService(mockUserRepo, mockTransactionRepo, mockUserMerchRepo)

	ctx := context.Background()
	fromUserID := int64(1)
	toUsername := "recipient"
	amount := int64(2000) // Больше, чем начальный баланс

	// Создаем тестового получателя
	recipient := &models.User{
		ID:       2,
		Username: toUsername,
		Coins:    1000,
	}

	// Настраиваем моки
	mockUserRepo.On("GetByUsername", ctx, toUsername).Return(recipient, nil)
	mockUserRepo.On("UpdateCoins", ctx, fromUserID, -amount).Return(errors.New("insufficient funds"))

	// Вызываем тестируемый метод
	err := service.SendCoins(ctx, fromUserID, toUsername, amount)

	// Проверяем результаты
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to deduct coins from sender")
	mockUserRepo.AssertExpectations(t)
}

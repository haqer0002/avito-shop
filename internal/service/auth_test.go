package service

import (
	"context"
	"testing"

	"github.com/haqer0002/avito-shop/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)

	ctx := context.Background()
	username := "testuser"
	password := "testpass"

	// Настраиваем мок
	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.User")).Return(nil)

	// Вызываем тестируемый метод
	err := service.CreateUser(ctx, username, password)

	// Проверяем результаты
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_GenerateToken(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo)
	service := authService.(*authServiceImpl)

	ctx := context.Background()
	username := "testuser"
	password := "testpass"

	// Создаем тестового пользователя
	testUser := &models.User{
		ID:       1,
		Username: username,
		Password: service.generatePasswordHash(password),
	}

	// Настраиваем мок
	mockRepo.On("GetByUsername", ctx, username).Return(testUser, nil)

	// Вызываем тестируемый метод
	token, err := service.GenerateToken(ctx, username, password)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_ParseToken(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo)
	service := authService.(*authServiceImpl)

	ctx := context.Background()
	username := "testuser"
	password := "testpass"

	// Создаем тестового пользователя
	testUser := &models.User{
		ID:       1,
		Username: username,
		Password: service.generatePasswordHash(password),
	}

	// Настраиваем мок
	mockRepo.On("GetByUsername", ctx, username).Return(testUser, nil)

	// Генерируем токен
	token, err := service.GenerateToken(ctx, username, password)
	assert.NoError(t, err)

	// Парсим токен
	userID, err := service.ParseToken(token)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, userID)
	mockRepo.AssertExpectations(t)
}

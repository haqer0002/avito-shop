package service

import (
	"context"

	"github.com/haqer0002/avito-shop/internal/models"
	"github.com/haqer0002/avito-shop/internal/repository"
)

// AuthService представляет интерфейс сервиса аутентификации
type AuthService interface {
	CreateUser(ctx context.Context, username, password string) error
	GenerateToken(ctx context.Context, username, password string) (string, error)
	ParseToken(token string) (int64, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

// UserService представляет интерфейс сервиса пользователей
type UserService interface {
	GetUserInfo(ctx context.Context, userID int64) (*models.InfoResponse, error)
	SendCoins(ctx context.Context, fromUserID int64, toUsername string, amount int64) error
}

// MerchService представляет интерфейс сервиса мерча
type MerchService interface {
	BuyMerch(ctx context.Context, userID int64, merchName string) error
	GetAllMerch(ctx context.Context) ([]models.MerchItem, error)
}

// Service представляет все сервисы приложения
type Service struct {
	Auth  AuthService
	User  UserService
	Merch MerchService
}

// NewService создает новый экземпляр Service
func NewService(repos *repository.Repository) *Service {
	return &Service{
		Auth:  NewAuthService(repos.Users),
		User:  NewUserService(repos.Users, repos.Transactions, repos.UserMerch),
		Merch: NewMerchService(repos.Users, repos.Merch, repos.UserMerch),
	}
}

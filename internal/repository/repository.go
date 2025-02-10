package repository

import (
	"context"

	"github.com/haqer0002/avito-shop/internal/models"
)

// UserRepository определяет методы для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
	UpdateCoins(ctx context.Context, userID int64, amount int64) error
}

// TransactionRepository определяет методы для работы с транзакциями
type TransactionRepository interface {
	Create(ctx context.Context, transaction *models.Transaction) error
	GetUserTransactions(ctx context.Context, userID int64) ([]models.Transaction, error)
}

// MerchRepository определяет методы для работы с мерчем
type MerchRepository interface {
	GetByName(ctx context.Context, name string) (*models.MerchItem, error)
	GetAll(ctx context.Context) ([]models.MerchItem, error)
}

// UserMerchRepository определяет методы для работы с купленным мерчем
type UserMerchRepository interface {
	Create(ctx context.Context, userMerch *models.UserMerch) error
	GetUserMerch(ctx context.Context, userID int64) ([]models.UserMerch, error)
}

// Repository объединяет все репозитории
type Repository struct {
	Users        UserRepository
	Transactions TransactionRepository
	Merch        MerchRepository
	UserMerch    UserMerchRepository
}

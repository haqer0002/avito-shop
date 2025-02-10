package postgres

import (
	"fmt"
	"log"

	"github.com/haqer0002/avito-shop/internal/repository"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewPostgresDB создает новое подключение к базе данных
func NewPostgresDB(connStr string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	log.Println("Successfully connected to database")
	return db, nil
}

// NewRepository создает новый экземпляр репозитория
func NewRepository(db *sqlx.DB) *repository.Repository {
	return &repository.Repository{
		Users:        NewUserRepository(db),
		Transactions: NewTransactionRepository(db),
		Merch:        NewMerchRepository(db),
		UserMerch:    NewUserMerchRepository(db),
	}
}

// Repository объединяет все репозитории
type Repository struct {
	Users        *UserRepository
	Transactions *TransactionRepository
	Merch        *MerchRepository
	UserMerch    *UserMerchRepository
}

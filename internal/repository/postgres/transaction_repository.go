package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/haqer0002/avito-shop/internal/models"
)

// TransactionRepository реализует интерфейс repository.TransactionRepository
type TransactionRepository struct {
	db *sqlx.DB
}

// NewTransactionRepository создает новый экземпляр TransactionRepository
func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

// Create создает новую транзакцию
func (r *TransactionRepository) Create(ctx context.Context, transaction *models.Transaction) error {
	query := `
		INSERT INTO transactions (from_user_id, to_user_id, amount, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		transaction.FromUserID,
		transaction.ToUserID,
		transaction.Amount,
		transaction.Description,
	).Scan(&transaction.ID, &transaction.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

// GetUserTransactions получает все транзакции пользователя
func (r *TransactionRepository) GetUserTransactions(ctx context.Context, userID int64) ([]models.Transaction, error) {
	query := `
		SELECT id, from_user_id, to_user_id, amount, created_at, description
		FROM transactions
		WHERE from_user_id = $1 OR to_user_id = $1
		ORDER BY created_at DESC`

	var transactions []models.Transaction
	err := r.db.SelectContext(ctx, &transactions, query, userID)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

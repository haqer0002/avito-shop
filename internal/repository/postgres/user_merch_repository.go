package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/haqer0002/avito-shop/internal/models"
)

// UserMerchRepository реализует интерфейс repository.UserMerchRepository
type UserMerchRepository struct {
	db *sqlx.DB
}

// NewUserMerchRepository создает новый экземпляр UserMerchRepository
func NewUserMerchRepository(db *sqlx.DB) *UserMerchRepository {
	return &UserMerchRepository{
		db: db,
	}
}

// Create создает запись о купленном мерче
func (r *UserMerchRepository) Create(ctx context.Context, userMerch *models.UserMerch) error {
	query := `
		INSERT INTO user_merch (user_id, merch_id)
		VALUES ($1, $2)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		userMerch.UserID,
		userMerch.MerchID,
	).Scan(&userMerch.ID, &userMerch.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

// GetUserMerch получает весь мерч пользователя
func (r *UserMerchRepository) GetUserMerch(ctx context.Context, userID int64) ([]models.UserMerch, error) {
	query := `
		SELECT um.id, um.user_id, um.merch_id, um.created_at
		FROM user_merch um
		WHERE um.user_id = $1
		ORDER BY um.created_at DESC`

	var userMerch []models.UserMerch
	err := r.db.SelectContext(ctx, &userMerch, query, userID)
	if err != nil {
		return nil, err
	}

	return userMerch, nil
}

package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/haqer0002/avito-shop/internal/models"
)

// MerchRepository реализует интерфейс repository.MerchRepository
type MerchRepository struct {
	db *sqlx.DB
}

// NewMerchRepository создает новый экземпляр MerchRepository
func NewMerchRepository(db *sqlx.DB) *MerchRepository {
	return &MerchRepository{
		db: db,
	}
}

// GetByName получает мерч по названию
func (r *MerchRepository) GetByName(ctx context.Context, name string) (*models.MerchItem, error) {
	item := &models.MerchItem{}
	query := `
		SELECT id, name, price
		FROM merch_items
		WHERE name = $1`

	err := r.db.GetContext(ctx, item, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("merch item not found")
		}
		return nil, err
	}

	return item, nil
}

// GetAll получает все доступные товары
func (r *MerchRepository) GetAll(ctx context.Context) ([]models.MerchItem, error) {
	var items []models.MerchItem
	query := `
		SELECT id, name, price
		FROM merch_items
		ORDER BY price ASC`

	err := r.db.SelectContext(ctx, &items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

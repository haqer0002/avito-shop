package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/haqer0002/avito-shop/internal/models"
	"github.com/jmoiron/sqlx"
)

// UserRepository реализует интерфейс repository.UserRepository
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create создает нового пользователя
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (username, password, coins)
		VALUES ($1, $2, $3)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		user.Username,
		user.Password,
		user.Coins,
	).Scan(&user.ID)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}

	log.Printf("Successfully created user with ID: %d", user.ID)
	return nil
}

// GetByUsername получает пользователя по имени пользователя
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, password, coins
		FROM users
		WHERE username = $1`

	err := r.db.GetContext(ctx, user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// UpdateCoins обновляет количество монет пользователя
func (r *UserRepository) UpdateCoins(ctx context.Context, userID int64, amount int64) error {
	query := `
		UPDATE users
		SET coins = coins + $1
		WHERE id = $2 AND coins + $1 >= 0`

	result, err := r.db.ExecContext(ctx, query, amount, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("insufficient funds")
	}

	return nil
}

// GetByID получает пользователя по ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	log.Printf("Getting user by ID: %d (type: %T)", id, id)
	user := &models.User{}
	query := `
		SELECT id, username, password, coins
		FROM users
		WHERE id = $1`

	log.Printf("Executing query: %s with ID: %d", query, id)
	err := r.db.GetContext(ctx, user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User with ID %d not found", id)
			return nil, errors.New("user not found")
		}
		log.Printf("Error getting user by ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	log.Printf("Successfully got user by ID %d: %+v", id, user)
	return user, nil
}

func (r *UserRepository) GetDB() *sqlx.DB {
	return r.db
}

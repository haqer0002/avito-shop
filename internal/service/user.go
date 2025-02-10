package service

import (
	"context"
	"fmt"
	"log"

	"github.com/haqer0002/avito-shop/internal/models"
	"github.com/haqer0002/avito-shop/internal/repository"
)

type userServiceImpl struct {
	userRepo        repository.UserRepository
	transactionRepo repository.TransactionRepository
	userMerchRepo   repository.UserMerchRepository
}

func NewUserService(userRepo repository.UserRepository, transactionRepo repository.TransactionRepository, userMerchRepo repository.UserMerchRepository) UserService {
	return &userServiceImpl{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
		userMerchRepo:   userMerchRepo,
	}
}

func (s *userServiceImpl) GetUserInfo(ctx context.Context, userID int64) (*models.InfoResponse, error) {
	log.Printf("Starting GetUserInfo for userID: %d (type: %T)", userID, userID)

	// Получаем информацию о пользователе
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		log.Printf("Error getting user by ID %d: %v", userID, err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	log.Printf("Successfully got user info: %+v", user)

	// Получаем транзакции пользователя
	transactions, err := s.transactionRepo.GetUserTransactions(ctx, userID)
	if err != nil {
		log.Printf("Error getting transactions for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	log.Printf("Successfully got transactions: %+v", transactions)

	// Получаем купленный мерч
	userMerch, err := s.userMerchRepo.GetUserMerch(ctx, userID)
	if err != nil {
		log.Printf("Error getting user merch for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to get user merch: %w", err)
	}
	log.Printf("Successfully got user merch: %+v", userMerch)

	// Формируем историю транзакций
	coinHistory := models.CoinTransactionHistory{
		Received: make([]models.CoinTransaction, 0),
		Sent:     make([]models.CoinTransaction, 0),
	}

	for _, t := range transactions {
		if t.ToUserID == userID {
			coinHistory.Received = append(coinHistory.Received, models.CoinTransaction{
				FromUser: fmt.Sprint(t.FromUserID),
				Amount:   t.Amount,
			})
		} else {
			coinHistory.Sent = append(coinHistory.Sent, models.CoinTransaction{
				ToUser: fmt.Sprint(t.ToUserID),
				Amount: t.Amount,
			})
		}
	}

	// Формируем инвентарь
	inventory := make(map[int64]int)
	for _, m := range userMerch {
		inventory[m.MerchID]++
	}

	inventoryItems := make([]models.InventoryItem, 0, len(inventory))
	for merchID, quantity := range inventory {
		inventoryItems = append(inventoryItems, models.InventoryItem{
			Type:     fmt.Sprint(merchID),
			Quantity: quantity,
		})
	}

	response := &models.InfoResponse{
		Coins:       user.Coins,
		Inventory:   inventoryItems,
		CoinHistory: coinHistory,
	}
	log.Printf("Successfully prepared response: %+v", response)
	return response, nil
}

func (s *userServiceImpl) SendCoins(ctx context.Context, fromUserID int64, toUsername string, amount int64) error {
	// Получаем пользователя-получателя
	toUser, err := s.userRepo.GetByUsername(ctx, toUsername)
	if err != nil {
		return fmt.Errorf("recipient not found: %w", err)
	}

	// Проверяем достаточность средств и списываем монеты у отправителя
	err = s.userRepo.UpdateCoins(ctx, fromUserID, -amount)
	if err != nil {
		return fmt.Errorf("failed to deduct coins from sender: %w", err)
	}

	// Начисляем монеты получателю
	err = s.userRepo.UpdateCoins(ctx, toUser.ID, amount)
	if err != nil {
		// В случае ошибки возвращаем монеты отправителю
		_ = s.userRepo.UpdateCoins(ctx, fromUserID, amount)
		return fmt.Errorf("failed to add coins to recipient: %w", err)
	}

	// Создаем запись о транзакции
	transaction := &models.Transaction{
		FromUserID:  fromUserID,
		ToUserID:    toUser.ID,
		Amount:      amount,
		Description: fmt.Sprintf("Transfer from user %d to user %s", fromUserID, toUsername),
	}

	err = s.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return fmt.Errorf("failed to create transaction record: %w", err)
	}

	return nil
}

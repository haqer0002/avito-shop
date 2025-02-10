package service

import (
	"context"
	"fmt"

	"github.com/haqer0002/avito-shop/internal/models"
	"github.com/haqer0002/avito-shop/internal/repository"
)

type merchServiceImpl struct {
	userRepo      repository.UserRepository
	merchRepo     repository.MerchRepository
	userMerchRepo repository.UserMerchRepository
}

func NewMerchService(userRepo repository.UserRepository, merchRepo repository.MerchRepository, userMerchRepo repository.UserMerchRepository) MerchService {
	return &merchServiceImpl{
		userRepo:      userRepo,
		merchRepo:     merchRepo,
		userMerchRepo: userMerchRepo,
	}
}

func (s *merchServiceImpl) BuyMerch(ctx context.Context, userID int64, merchName string) error {
	// Получаем информацию о мерче
	merch, err := s.merchRepo.GetByName(ctx, merchName)
	if err != nil {
		return fmt.Errorf("merch not found: %w", err)
	}

	// Списываем монеты у пользователя
	err = s.userRepo.UpdateCoins(ctx, userID, -merch.Price)
	if err != nil {
		return fmt.Errorf("failed to deduct coins: %w", err)
	}

	// Создаем запись о покупке
	userMerch := &models.UserMerch{
		UserID:  userID,
		MerchID: merch.ID,
	}

	err = s.userMerchRepo.Create(ctx, userMerch)
	if err != nil {
		// В случае ошибки возвращаем монеты пользователю
		_ = s.userRepo.UpdateCoins(ctx, userID, merch.Price)
		return fmt.Errorf("failed to record purchase: %w", err)
	}

	return nil
}

func (s *merchServiceImpl) GetAllMerch(ctx context.Context) ([]models.MerchItem, error) {
	return s.merchRepo.GetAll(ctx)
}

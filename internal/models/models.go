package models

import "time"

// User представляет пользователя системы
type User struct {
	ID       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"-" db:"password"`
	Coins    int64  `json:"coins" db:"coins"`
}

// Transaction представляет транзакцию между пользователями
type Transaction struct {
	ID          int64     `json:"id" db:"id"`
	FromUserID  int64     `json:"from_user_id" db:"from_user_id"`
	ToUserID    int64     `json:"to_user_id" db:"to_user_id"`
	Amount      int64     `json:"amount" db:"amount"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Description string    `json:"description" db:"description"`
}

// MerchItem представляет товар в магазине
type MerchItem struct {
	ID    int64  `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Price int64  `json:"price" db:"price"`
}

// UserMerch представляет купленный пользователем мерч
type UserMerch struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	MerchID   int64     `json:"merch_id" db:"merch_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// InfoResponse представляет ответ на запрос информации о пользователе
type InfoResponse struct {
	Coins       int64                  `json:"coins"`
	Inventory   []InventoryItem        `json:"inventory"`
	CoinHistory CoinTransactionHistory `json:"coinHistory"`
}

// InventoryItem представляет предмет в инвентаре пользователя
type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

// CoinTransactionHistory представляет историю транзакций пользователя
type CoinTransactionHistory struct {
	Received []CoinTransaction `json:"received"`
	Sent     []CoinTransaction `json:"sent"`
}

// CoinTransaction представляет отдельную транзакцию
type CoinTransaction struct {
	FromUser string `json:"fromUser,omitempty"`
	ToUser   string `json:"toUser,omitempty"`
	Amount   int64  `json:"amount"`
}

// AuthRequest представляет запрос на аутентификацию
type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse представляет ответ на аутентификацию
type AuthResponse struct {
	Token string `json:"token"`
}

// SendCoinRequest представляет запрос на отправку монет
type SendCoinRequest struct {
	ToUser string `json:"toUser" binding:"required"`
	Amount int64  `json:"amount" binding:"required,gt=0"`
}

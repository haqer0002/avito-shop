package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/haqer0002/avito-shop/internal/models"
	"github.com/haqer0002/avito-shop/internal/repository"
)

const (
	salt       = "hjqrhjqw124617ajfhajs"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"user_id"`
}

type authServiceImpl struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authServiceImpl{repo: repo}
}

func (s *authServiceImpl) CreateUser(ctx context.Context, username, password string) error {
	// Валидация входных данных
	if username == "" || password == "" {
		log.Printf("Validation error: username and password cannot be empty")
		return errors.New("username and password cannot be empty")
	}

	if len(username) < 3 {
		log.Printf("Validation error: username must be at least 3 characters long")
		return errors.New("username must be at least 3 characters long")
	}

	if len(password) < 6 {
		log.Printf("Validation error: password must be at least 6 characters long")
		return errors.New("password must be at least 6 characters long")
	}

	hashedPassword := s.generatePasswordHash(password)
	log.Printf("Generated password hash for user %s: %s", username, hashedPassword)

	user := &models.User{
		Username: username,
		Password: hashedPassword,
		Coins:    1000, // Начальное количество монет
	}

	log.Printf("Attempting to create user: %s with coins: %d", username, user.Coins)
	if err := s.repo.Create(ctx, user); err != nil {
		log.Printf("Error creating user in repository: %v", err)
		return err
	}

	log.Printf("Successfully created user %s with ID: %d", username, user.ID)
	return nil
}

func (s *authServiceImpl) GenerateToken(ctx context.Context, username, password string) (string, error) {
	log.Printf("Attempting to get user: %s", username)
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return "", err
	}

	hashedPassword := s.generatePasswordHash(password)
	log.Printf("Comparing passwords for user %s. Input hash: %s, Stored hash: %s", username, hashedPassword, user.Password)
	if user.Password != hashedPassword {
		log.Printf("Password mismatch for user %s", username)
		return "", errors.New("incorrect password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: user.ID,
	})

	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return "", err
	}

	log.Printf("Successfully generated token for user %s", username)
	return tokenString, nil
}

func (s *authServiceImpl) ParseToken(accessToken string) (int64, error) {
	log.Printf("Parsing token: %s", accessToken)
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Invalid signing method: %v", token.Method)
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		log.Printf("Token claims are not of type *tokenClaims")
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	log.Printf("Successfully parsed token. User ID: %d (type: %T)", claims.UserID, claims.UserID)
	return claims.UserID, nil
}

func (s *authServiceImpl) generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	hash.Write([]byte(salt))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (s *authServiceImpl) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	log.Printf("Getting user by username: %s", username)
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		log.Printf("Error getting user by username: %v", err)
		return nil, err
	}
	log.Printf("Successfully got user by username: %+v", user)
	return user, nil
}

func (s *authServiceImpl) GetRepo() repository.UserRepository {
	return s.repo
}

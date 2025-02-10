package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/haqer0002/avito-shop/internal/config"
	"github.com/haqer0002/avito-shop/internal/handlers"
	"github.com/haqer0002/avito-shop/internal/models"
	"github.com/haqer0002/avito-shop/internal/repository/postgres"
	"github.com/haqer0002/avito-shop/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testDB    *sqlx.DB
	handler   *handlers.Handler
	testToken string
)

func TestMain(m *testing.M) {
	// Загружаем тестовую конфигурацию
	os.Setenv("DB_NAME", "postgres")
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Подключаемся к тестовой базе данных
	testDB, err = postgres.NewPostgresDB(cfg.GetDBConnString())
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}

	repos := postgres.NewRepository(testDB)

	services := service.NewService(repos)

	handler = handlers.NewHandler(services)

	code := m.Run()

	testDB.Close()

	os.Exit(code)
}

func TestIntegration_Auth(t *testing.T) {
	// Создаем тестовый запрос на аутентификацию
	authData := models.AuthRequest{
		Username: "testuser",
		Password: "testpass",
	}

	body, err := json.Marshal(authData)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/auth/sign-up", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder для записи ответа
	w := httptest.NewRecorder()

	// Выполняем запрос
	handler.InitRoutes().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.AuthResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	testToken = response.Token
	assert.NotEmpty(t, testToken)
}

func TestIntegration_BuyMerch(t *testing.T) {
	// Пропускаем тест, если нет токена
	if testToken == "" {
		t.Skip("No auth token available")
	}

	// Создаем тестовый запрос на покупку мерча
	req := httptest.NewRequest("POST", "/api/merch/buy/t-shirt", nil)
	req.Header.Set("Authorization", "Bearer "+testToken)

	w := httptest.NewRecorder()

	// Выполняем запрос
	handler.InitRoutes().ServeHTTP(w, req)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_SendCoins(t *testing.T) {
	// Пропускаем тест, если нет токена
	if testToken == "" {
		t.Skip("No auth token available")
	}

	// Создаем тестовый запрос на отправку монет
	sendData := models.SendCoinRequest{
		ToUser: "anotheruser",
		Amount: 100,
	}

	body, err := json.Marshal(sendData)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/user/send", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+testToken)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.InitRoutes().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_GetInfo(t *testing.T) {
	// Пропускаем тест, если нет токена
	if testToken == "" {
		t.Skip("No auth token available")
	}

	// Создаем тестовый запрос на получение информации
	req := httptest.NewRequest("GET", "/api/user/info", nil)
	req.Header.Set("Authorization", "Bearer "+testToken)

	w := httptest.NewRecorder()

	handler.InitRoutes().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.InfoResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.NotNil(t, response.Coins)
	assert.NotNil(t, response.Inventory)
	assert.NotNil(t, response.CoinHistory)
}

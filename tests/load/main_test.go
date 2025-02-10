package load

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/haqer0002/avito-shop/internal/config"
	"github.com/haqer0002/avito-shop/internal/handlers"
	"github.com/haqer0002/avito-shop/internal/models"
	"github.com/haqer0002/avito-shop/internal/repository/postgres"
	"github.com/haqer0002/avito-shop/internal/service"
)

const (
	concurrentUsers = 100
	requestsPerUser = 10
)

func TestLoadBuyMerch(t *testing.T) {
	// Загрузка конфигурации
	os.Setenv("DB_NAME", "postgres")
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	// Подключение к базе данных
	db, err := postgres.NewPostgresDB(cfg.GetDBConnString())
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Инициализация репозиториев и сервисов
	repos := postgres.NewRepository(db)
	services := service.NewService(repos)
	handler := handlers.NewHandler(services)

	// Создание тестового сервера
	server := httptest.NewServer(handler.InitRoutes())
	defer server.Close()

	var wg sync.WaitGroup
	errorsChan := make(chan error, concurrentUsers*requestsPerUser)
	successCount := 0
	var mu sync.Mutex

	startTime := time.Now()

	// Запуск конкурентных пользователей
	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()

			// Аутентификация пользователя
			token, err := authenticateUser(server.URL, fmt.Sprintf("user%d", userID))
			if err != nil {
				errorsChan <- fmt.Errorf("authentication failed for user%d: %v", userID, err)
				return
			}

			// Выполнение запросов на покупку мерча
			for j := 0; j < requestsPerUser; j++ {
				err := buyMerch(server.URL, token)
				if err != nil {
					errorsChan <- fmt.Errorf("buy merch failed for user%d, request%d: %v", userID, j, err)
				} else {
					mu.Lock()
					successCount++
					mu.Unlock()
				}
				time.Sleep(time.Millisecond * 100)
			}
		}(i)
	}

	wg.Wait()
	close(errorsChan)

	duration := time.Since(startTime)
	totalRequests := concurrentUsers * requestsPerUser

	var errors []error
	for err := range errorsChan {
		errors = append(errors, err)
	}

	t.Logf("Load test completed in %v", duration)
	t.Logf("Total requests: %d", totalRequests)
	t.Logf("Successful requests: %d", successCount)
	t.Logf("Failed requests: %d", len(errors))
	t.Logf("Success rate: %.2f%%", float64(successCount)/float64(totalRequests)*100)
	t.Logf("Requests per second: %.2f", float64(totalRequests)/duration.Seconds())

	if len(errors) > 0 {
		t.Logf("Sample of errors:")
		for i := 0; i < min(5, len(errors)); i++ {
			t.Logf("  %v", errors[i])
		}
	}
}

func authenticateUser(baseURL, username string) (string, error) {
	authData := models.AuthRequest{
		Username: username,
		Password: "testpass",
	}
	body, err := json.Marshal(authData)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(baseURL+"/auth/sign-up", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authentication failed with status: %d", resp.StatusCode)
	}

	var authResponse models.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return "", err
	}

	return authResponse.Token, nil
}

func buyMerch(baseURL, token string) error {
	req, err := http.NewRequest("POST", baseURL+"/api/merch/buy/t-shirt", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("buy merch failed with status: %d", resp.StatusCode)
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

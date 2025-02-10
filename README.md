# Avito Shop

## Описание проекта

Сервис внутреннего магазина мерча Авито, где сотрудники могут приобретать товары за внутреннюю валюту (coins).

### Функциональность

- Регистрация и авторизация пользователей через JWT
- Управление балансом пользователей
- Система транзакций между пользователями
- Каталог мерча и система покупок
- История транзакций и инвентарь пользователя

### Технологический стек

- **Go 1.21+**
- **PostgreSQL** - основная база данных
- **Docker** и **Docker Compose** - контейнеризация
- **JWT** - аутентификация
- **Gin** - веб-фреймворк
- **SQLx** - работа с базой данных
- **Testify** - тестирование

### Установка и запуск

#### Предварительные требования

- Go 1.21 или выше
- Docker и Docker Compose
- PostgreSQL (если запуск без Docker)

#### Запуск через Docker

```bash
# Клонирование репозитория
git clone https://github.com/your-username/avito-shop.git
cd avito-shop

# Запуск через Docker Compose
docker-compose up -d
```

#### Локальный запуск

1. Клонируйте репозиторий
```bash
git clone https://github.com/your-username/avito-shop.git
cd avito-shop
```

2. Создайте файл .env на основе .env.example
```bash
cp .env.example .env
```

3. Установите зависимости
```bash
go mod download
```

4. Запустите PostgreSQL и примените миграции

5. Запустите приложение
```bash
go run cmd/app/main.go
```

### API Endpoints

#### Аутентификация

##### POST /auth/sign-up
Регистрация и авторизация пользователя
```json
{
    "username": "user123",
    "password": "password123"
}
```

#### Пользователь

##### GET /api/user/info
Получение информации о балансе и транзакциях (требует авторизации)

##### POST /api/user/send
Отправка монет другому пользователю (требует авторизации)
```json
{
    "toUser": "recipient123",
    "amount": 100
}
```

#### Мерч

##### GET /api/merch/list
Получение списка доступного мерча (требует авторизации)

##### POST /api/merch/buy/:item
Покупка мерча (требует авторизации)

### Тестирование

```bash
# Unit-тесты
go test ./internal/...

# Интеграционные тесты
go test ./tests/integration

# Нагрузочные тесты
go test ./tests/load
```

### Безопасность

- Все пароли хешируются перед сохранением
- Используется JWT для аутентификации
- Реализована защита от отрицательного баланса
- Транзакции выполняются атомарно

### Масштабирование

Проект подготовлен к горизонтальному масштабированию:
- Stateless архитектура
- Docker контейнеризация
- Возможность запуска нескольких инстансов

---

# Avito Shop (English)

## Project Description

Internal merchandise shop service for Avito, where employees can purchase items using internal currency (coins).

### Features

- User registration and authentication via JWT
- User balance management
- User-to-user transaction system
- Merchandise catalog and purchase system
- Transaction history and user inventory

### Tech Stack

- **Go 1.21+**
- **PostgreSQL** - main database
- **Docker** and **Docker Compose** - containerization
- **JWT** - authentication
- **Gin** - web framework
- **SQLx** - database operations
- **Testify** - testing

### Installation and Running

#### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL (if running without Docker)

#### Running with Docker

```bash
# Clone repository
git clone https://github.com/your-username/avito-shop.git
cd avito-shop

# Run with Docker Compose
docker-compose up -d
```

#### Local Setup

1. Clone the repository
```bash
git clone https://github.com/your-username/avito-shop.git
cd avito-shop
```

2. Create .env file from .env.example
```bash
cp .env.example .env
```

3. Install dependencies
```bash
go mod download
```

4. Start PostgreSQL and apply migrations

5. Run the application
```bash
go run cmd/app/main.go
```

### API Endpoints

#### Authentication

##### POST /auth/sign-up
User registration and authentication
```json
{
    "username": "user123",
    "password": "password123"
}
```

#### User

##### GET /api/user/info
Get balance and transaction information (requires authentication)

##### POST /api/user/send
Send coins to another user (requires authentication)
```json
{
    "toUser": "recipient123",
    "amount": 100
}
```

#### Merchandise

##### GET /api/merch/list
Get available merchandise list (requires authentication)

##### POST /api/merch/buy/:item
Purchase merchandise (requires authentication)

### Testing

```bash
# Unit tests
go test ./internal/...

# Integration tests
go test ./tests/integration

# Load tests
go test ./tests/load
```

### Security

- All passwords are hashed before storage
- JWT authentication
- Protection against negative balance
- Atomic transactions

### Scalability

The project is prepared for horizontal scaling:
- Stateless architecture
- Docker containerization
- Multiple instance capability

## Project Structure

```
.
├── cmd/                    # Application entry point
├── internal/              
│   ├── config/            # Configuration
│   ├── handlers/          # HTTP handlers
│   ├── middleware/        # Middleware components
│   ├── models/            # Data models
│   ├── repository/        # Database layer
│   └── service/           # Business logic
├── migrations/            # SQL migrations
└── tests/                 # Tests
    ├── integration/       # Integration tests
    └── load/              # Load tests
```

## Автор

- [Ваше имя](https://github.com/your-username)

## Лицензия

MIT License
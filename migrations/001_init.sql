-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    coins BIGINT NOT NULL DEFAULT 1000
);

-- Создание таблицы транзакций
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    from_user_id BIGINT REFERENCES users(id),
    to_user_id BIGINT REFERENCES users(id),
    amount BIGINT NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы мерча
CREATE TABLE IF NOT EXISTS merch_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    price BIGINT NOT NULL
);

-- Создание таблицы купленного мерча
CREATE TABLE IF NOT EXISTS user_merch (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    merch_id BIGINT REFERENCES merch_items(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Добавляем начальные товары
INSERT INTO merch_items (name, price) VALUES
    ('t-shirt', 50),
    ('hoodie', 100),
    ('cap', 30),
    ('stickers', 10); 
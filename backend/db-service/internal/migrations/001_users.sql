-- +goose Up
CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    crypted_password TEXT NOT NULL,
    telegram_id TEXT NOT NULL,
    crypted_refresh TEXT,
    room TEXt NOT NULL,
    dorm_number TEXT,
    role TEXT
);

-- +goose Down
DROP TABLE users;
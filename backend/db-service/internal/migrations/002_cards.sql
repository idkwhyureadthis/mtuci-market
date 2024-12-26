-- +goose Up
CREATE TABLE cards(
    id SERIAL PRIMARY KEY,
    created_by INTEGER REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    card_status TEXT,
    price FLOAT
);

-- +goose Down
DROP TABLE cards;
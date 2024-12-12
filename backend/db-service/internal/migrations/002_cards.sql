--+goose Up
CREATE TYPE card_state AS ENUM('accepted', 'on_moderation', 'rejected', 'sold');

CREATE TABLE cards(
    id INTEGER PRIMARY KEY,
    created_by INTEGER REFERENCES users(id),
    name TEXT NOT NULL,
    card_status card_state,
    price FLOAT
);

-- +goose Down
DROP TABLE cards;
DROP TYPE card_state;
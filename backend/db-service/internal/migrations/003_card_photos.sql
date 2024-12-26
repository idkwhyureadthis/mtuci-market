-- +goose Up
CREATE TABLE card_photos(
    card_id INTEGER REFERENCES cards(id) ON DELETE CASCADE,
    photo_link TEXT NOT NULL
);

-- +goose Down
DROP TABLE card_photos;
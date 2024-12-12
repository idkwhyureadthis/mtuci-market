-- +goose Up
CREATE TABLE purchases(
    card_id INTEGER REFERENCES cards(id) PRIMARY KEY,
    bought_by INTEGER REFERENCES 
)


-- +goose Down
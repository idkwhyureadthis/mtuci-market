-- +goose Up
CREATE TYPE user_role AS ENUM('user', 'moderator', 'admin');

CREATE TABLE users(
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    about TEXT,
    crypted_refresh TEXT,
    room INTEGER NOT NULL,
    dorm_number INTEGER
    role user_role,
);

-- +goose Down
DROP TABLE users;
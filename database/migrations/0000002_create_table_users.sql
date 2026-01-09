-- +goose up
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL CHECK (LENGTH(password) > 8)
);

-- +goose down
DROP TABLE users;

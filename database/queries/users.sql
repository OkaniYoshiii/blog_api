-- name: ListUsers :many
SELECT id, email, password FROM users;

-- name: GetUserByEmail :one
SELECT id, email, password FROM users WHERE users.email = ?;

-- name: CreateUser :one
INSERT INTO users (email, password) VALUES (?, ?) RETURNING *;

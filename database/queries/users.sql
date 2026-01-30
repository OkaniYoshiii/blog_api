-- name: ListUsers :many
SELECT id, email, password FROM users;

-- name: GetUserById :one
SELECT id, email FROM users WHERE users.id = ?;

-- name: GetUserByEmail :one
SELECT id, email, password FROM users WHERE users.email = ?;

-- name: CreateUser :one
INSERT INTO users (email, password) VALUES (?, ?) RETURNING *;

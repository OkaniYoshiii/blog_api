-- name: ListApiKeys :many
SELECT id, value, created_at, updated_at, revoked_at FROM api_keys;

-- name: GetApiKeyByValue :one
SELECT id FROM api_keys WHERE api_keys.value = ?;

-- name: CreateApiKey :exec
INSERT INTO api_keys (
    value, application
) VALUES (
    ?, ?
);

-- +goose up
CREATE TABLE api_keys (
    id INTEGER PRIMARY KEY,
    value TEXT NOT NULL UNIQUE,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    revoked_at TEXT,
    application TEXT NOT NULL CHECK(application IN ('web_backend', 'web_frontend'))
);

-- +goose down
DROP TABLE api_keys;

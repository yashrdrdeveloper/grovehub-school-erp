-- +goose Up
ALTER TABLE identity.password_tokens ADD COLUMN session_generation bigint NOT NULL DEFAULT 0 CHECK(session_generation >= 0);

-- +goose Down
ALTER TABLE identity.password_tokens DROP COLUMN session_generation;

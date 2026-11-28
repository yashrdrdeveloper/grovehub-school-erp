-- +goose Up
ALTER TABLE integration.oauth_states ADD COLUMN session_generation bigint NOT NULL DEFAULT 0 CHECK(session_generation >= 0);

-- +goose Down
ALTER TABLE integration.oauth_states DROP COLUMN session_generation;

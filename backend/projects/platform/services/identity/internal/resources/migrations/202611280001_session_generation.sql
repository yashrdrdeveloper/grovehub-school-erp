-- +goose Up
ALTER TABLE identity.sessions ADD COLUMN session_generation bigint NOT NULL DEFAULT 0 CHECK(session_generation >= 0);

-- Existing unbound sessions must authenticate again when the registry policy
-- is enabled. Zero is not a valid registry-issued generation.

-- +goose Down
ALTER TABLE identity.sessions DROP COLUMN session_generation;

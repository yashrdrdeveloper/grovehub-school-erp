-- +goose Up
CREATE SCHEMA IF NOT EXISTS accounting;

-- +goose Down
DROP SCHEMA IF EXISTS accounting;

-- +goose Up
CREATE SCHEMA IF NOT EXISTS foundation;

-- +goose Down
DROP SCHEMA IF EXISTS foundation;

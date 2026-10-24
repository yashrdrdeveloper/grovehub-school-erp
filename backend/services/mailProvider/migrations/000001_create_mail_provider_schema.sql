-- +goose Up
CREATE SCHEMA IF NOT EXISTS mail_provider;

-- +goose Down
DROP SCHEMA IF EXISTS mail_provider CASCADE;

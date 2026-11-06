-- +goose Up
-- Prior versions persisted message references and consumer cursors in Integration.
-- Preserve data for an eventual Communication-owned migration; never discard it.
CREATE SCHEMA IF NOT EXISTS integration_legacy;
ALTER TABLE integration.mail_references SET SCHEMA integration_legacy;
ALTER TABLE integration.sync_states SET SCHEMA integration_legacy;

-- +goose Down
ALTER TABLE integration_legacy.sync_states SET SCHEMA integration;
ALTER TABLE integration_legacy.mail_references SET SCHEMA integration;
DROP SCHEMA IF EXISTS integration_legacy;

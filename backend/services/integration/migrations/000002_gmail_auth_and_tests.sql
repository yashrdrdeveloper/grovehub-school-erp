-- +goose Up
ALTER TABLE integration.connections
 ADD COLUMN external_id text,
 ADD COLUMN last_tested_at timestamptz,
 ADD COLUMN last_test_status text;
CREATE UNIQUE INDEX uq_integration_connection_provider_external
 ON integration.connections(tenant_id,provider,external_id)
 WHERE external_id IS NOT NULL AND status = 'ACTIVE';
CREATE TABLE integration.oauth_states (
 hash text PRIMARY KEY,
 connection_id uuid NOT NULL REFERENCES integration.connections(id) ON DELETE CASCADE,
 expires_at timestamptz NOT NULL
);
CREATE INDEX idx_integration_oauth_states_expiry ON integration.oauth_states(expires_at);
CREATE TABLE integration.oauth_credentials (
 connection_id uuid PRIMARY KEY REFERENCES integration.connections(id) ON DELETE CASCADE,
 ciphertext bytea NOT NULL,
 updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE integration.connection_tests (
 id uuid PRIMARY KEY,
 connection_id uuid NOT NULL REFERENCES integration.connections(id) ON DELETE CASCADE,
 direction text NOT NULL CHECK(direction IN ('INBOUND','OUTBOUND')),
 status text NOT NULL CHECK(status IN ('SUCCESS','FAILED')),
 detail text NOT NULL,
 started_at timestamptz NOT NULL,
 completed_at timestamptz NOT NULL
);
CREATE INDEX idx_integration_tests_recent ON integration.connection_tests(connection_id,completed_at DESC,id DESC);
-- +goose Down
DROP TABLE IF EXISTS integration.connection_tests;
DROP TABLE IF EXISTS integration.oauth_credentials;
DROP TABLE IF EXISTS integration.oauth_states;
DROP INDEX IF EXISTS integration.uq_integration_connection_provider_external;
ALTER TABLE integration.connections DROP COLUMN last_test_status,DROP COLUMN last_tested_at,DROP COLUMN external_id;

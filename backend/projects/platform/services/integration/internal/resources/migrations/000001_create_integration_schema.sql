-- +goose Up
-- Development baseline for Integration. This schema owns provider connectivity only.
-- IMPORTANT: Existing databases that already ran historical 000002..000005
-- require a deliberate data-preserving migration/rebaseline procedure.
CREATE SCHEMA IF NOT EXISTS integration;

CREATE TABLE IF NOT EXISTS integration.connections (
    id uuid PRIMARY KEY,
    tenant_id uuid NOT NULL,
    owner_gid uuid NOT NULL,
    provider text NOT NULL,
    credential_ref text,
    external_id text,
    status text NOT NULL DEFAULT 'PENDING_SETUP',
    last_tested_at timestamptz,
    last_test_status text,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_integration_provider CHECK (length(btrim(provider)) > 0),
    CONSTRAINT chk_integration_status CHECK (status IN ('PENDING_SETUP','ACTIVE','DISCONNECTED'))
);

CREATE INDEX IF NOT EXISTS idx_integration_connections_tenant_owner
    ON integration.connections(tenant_id, owner_gid, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_integration_active_mailbox
    ON integration.connections(provider, lower(external_id))
    WHERE status = 'ACTIVE' AND external_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uq_integration_connection_provider_external
    ON integration.connections(tenant_id, provider, external_id)
    WHERE external_id IS NOT NULL AND status = 'ACTIVE';

CREATE TABLE IF NOT EXISTS integration.oauth_states (
    hash text PRIMARY KEY,
    connection_id uuid NOT NULL REFERENCES integration.connections(id) ON DELETE CASCADE,
    expires_at timestamptz NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_integration_oauth_states_expiry
    ON integration.oauth_states(expires_at);

CREATE TABLE IF NOT EXISTS integration.oauth_credentials (
    connection_id uuid PRIMARY KEY REFERENCES integration.connections(id) ON DELETE CASCADE,
    ciphertext bytea NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS integration.connection_tests (
    id uuid PRIMARY KEY,
    connection_id uuid NOT NULL REFERENCES integration.connections(id) ON DELETE CASCADE,
    direction text NOT NULL CHECK(direction IN ('INBOUND','OUTBOUND')),
    status text NOT NULL CHECK(status IN ('SUCCESS','FAILED')),
    detail text NOT NULL,
    started_at timestamptz NOT NULL,
    completed_at timestamptz NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_integration_tests_recent
    ON integration.connection_tests(connection_id, completed_at DESC, id DESC);

CREATE TABLE IF NOT EXISTS integration.webhook_subscriptions (
    id uuid PRIMARY KEY,
    connection_id uuid NOT NULL REFERENCES integration.connections(id) ON DELETE CASCADE,
    resource text NOT NULL,
    external_id text NOT NULL,
    status text NOT NULL,
    expires_at timestamptz,
    CONSTRAINT uq_integration_webhook_resource UNIQUE(connection_id, resource),
    CONSTRAINT uq_integration_webhook_external UNIQUE(connection_id, external_id)
);

CREATE INDEX IF NOT EXISTS idx_integration_watch_expiry
    ON integration.webhook_subscriptions(expires_at)
    WHERE status = 'ACTIVE';

-- +goose Down
DROP SCHEMA IF EXISTS integration CASCADE;

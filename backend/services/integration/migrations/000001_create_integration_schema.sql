-- +goose Up
CREATE SCHEMA IF NOT EXISTS integration;

CREATE TABLE IF NOT EXISTS integration.connections (
 id uuid PRIMARY KEY,
 tenant_id uuid NOT NULL,
 owner_gid uuid NOT NULL,
 provider text NOT NULL,
 credential_ref text,
 status text NOT NULL DEFAULT 'PENDING_SETUP',
 created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
 updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
 CONSTRAINT chk_integration_provider CHECK (length(btrim(provider)) > 0),
 CONSTRAINT chk_integration_status CHECK (status IN ('PENDING_SETUP','ACTIVE','DISCONNECTED'))
);
CREATE INDEX IF NOT EXISTS idx_integration_connections_tenant_owner
 ON integration.connections(tenant_id,owner_gid,created_at DESC);

CREATE TABLE IF NOT EXISTS integration.sync_states (
 id uuid PRIMARY KEY,
 connection_id uuid NOT NULL REFERENCES integration.connections(id) ON DELETE CASCADE,
 resource text NOT NULL,
 cursor text NOT NULL,
 last_synced_at timestamptz,
 CONSTRAINT uq_integration_sync_resource UNIQUE(connection_id,resource)
);

CREATE TABLE IF NOT EXISTS integration.webhook_subscriptions (
 id uuid PRIMARY KEY,
 connection_id uuid NOT NULL REFERENCES integration.connections(id) ON DELETE CASCADE,
 resource text NOT NULL,
 external_id text NOT NULL,
 status text NOT NULL,
 expires_at timestamptz,
 CONSTRAINT uq_integration_webhook_external UNIQUE(connection_id,external_id)
);
-- +goose Down
DROP SCHEMA IF EXISTS integration CASCADE;

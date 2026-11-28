-- +goose Up
CREATE INDEX idx_integration_connections_owner_order ON integration.connections(tenant_id,owner_gid,created_at DESC,id DESC);
CREATE INDEX idx_integration_active_mailbox_exact ON integration.connections(provider,external_id) WHERE status='ACTIVE' AND external_id IS NOT NULL;
CREATE INDEX idx_integration_oauth_states_connection ON integration.oauth_states(connection_id);
-- +goose Down
DROP INDEX integration.idx_integration_oauth_states_connection;
DROP INDEX integration.idx_integration_active_mailbox_exact;
DROP INDEX integration.idx_integration_connections_owner_order;

-- +goose Up
CREATE UNIQUE INDEX uq_integration_webhook_resource
 ON integration.webhook_subscriptions(connection_id,resource);
CREATE INDEX idx_integration_active_gmail_mailbox
 ON integration.connections(provider,lower(external_id))
 WHERE status = 'ACTIVE';
CREATE INDEX idx_integration_watch_expiry
 ON integration.webhook_subscriptions(expires_at)
 WHERE status = 'ACTIVE';
-- +goose Down
DROP INDEX IF EXISTS integration.idx_integration_watch_expiry;
DROP INDEX IF EXISTS integration.idx_integration_active_gmail_mailbox;
DROP INDEX IF EXISTS integration.uq_integration_webhook_resource;

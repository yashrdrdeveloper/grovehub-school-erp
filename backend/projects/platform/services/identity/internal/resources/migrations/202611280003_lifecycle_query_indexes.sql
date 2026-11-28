-- +goose Up
CREATE INDEX idx_identity_accounts_tenant_account ON identity.accounts(tenant_id, account_id) WHERE tenant_id IS NOT NULL;
CREATE INDEX idx_identity_sessions_active_generation ON identity.sessions(account_id, session_generation, created_at DESC) INCLUDE (refresh_expires_at) WHERE revoked_at IS NULL;
CREATE INDEX idx_identity_invitations_ready ON identity.invitations(next_attempt_at, id) WHERE sent_at IS NULL AND suppressed_at IS NULL;
CREATE INDEX idx_identity_tokens_live_account ON identity.password_tokens(account_id, purpose) WHERE used_at IS NULL;
-- +goose Down
DROP INDEX identity.idx_identity_tokens_live_account;
DROP INDEX identity.idx_identity_invitations_ready;
DROP INDEX identity.idx_identity_sessions_active_generation;
DROP INDEX identity.idx_identity_accounts_tenant_account;

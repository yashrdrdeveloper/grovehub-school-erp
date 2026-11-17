-- +goose Up
CREATE TABLE IF NOT EXISTS identity.password_setup_tokens (
    token_hash CHAR(64) PRIMARY KEY,
    account_id UUID NOT NULL REFERENCES identity.accounts(account_id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_password_setup_tokens_account ON identity.password_setup_tokens (account_id);
CREATE INDEX IF NOT EXISTS idx_password_setup_tokens_expiry ON identity.password_setup_tokens (expires_at) WHERE used_at IS NULL;

CREATE TABLE IF NOT EXISTS identity.password_setup_mail (
    account_id UUID PRIMARY KEY REFERENCES identity.accounts(account_id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    ciphertext BYTEA NOT NULL,
    nonce BYTEA NOT NULL,
    attempts INTEGER NOT NULL DEFAULT 0,
    next_attempt_at TIMESTAMPTZ NOT NULL,
    lease_until TIMESTAMPTZ,
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_password_setup_mail_pending ON identity.password_setup_mail (next_attempt_at) WHERE sent_at IS NULL;

-- IMPORTANT: Existing account status CHECK/ENUM constraints must be inspected
-- and updated explicitly to permit PENDING_SETUP before rollout. Do not drop
-- unknown constraints blindly; migration verification is required.

-- +goose Down
DROP TABLE IF EXISTS identity.password_setup_mail;
DROP TABLE IF EXISTS identity.password_setup_tokens;

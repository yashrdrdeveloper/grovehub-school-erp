-- +goose Up
CREATE TABLE IF NOT EXISTS identity.sessions (
    session_id UUID PRIMARY KEY,
    account_id UUID NOT NULL REFERENCES identity.accounts(account_id) ON DELETE CASCADE,
    access_hash CHAR(64) NOT NULL UNIQUE,
    refresh_hash CHAR(64) NOT NULL UNIQUE,
    access_expires_at TIMESTAMPTZ NOT NULL,
    refresh_expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_identity_sessions_account ON identity.sessions (account_id);
CREATE INDEX IF NOT EXISTS idx_identity_sessions_refresh_expiry ON identity.sessions (refresh_expires_at) WHERE revoked_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS identity.sessions;

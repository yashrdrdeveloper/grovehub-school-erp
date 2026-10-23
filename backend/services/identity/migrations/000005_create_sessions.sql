-- +goose Up
CREATE TABLE identity.sessions (
    session_id uuid NOT NULL,
    account_id uuid NOT NULL,
    refresh_token_hash text NOT NULL,
    status text NOT NULL DEFAULT 'ACTIVE',
    expires_at timestamptz NOT NULL,
    last_seen_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at timestamptz NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT pk_sessions PRIMARY KEY (session_id),
    CONSTRAINT fk_sessions_account FOREIGN KEY (account_id)
        REFERENCES identity.accounts(account_id) ON DELETE CASCADE,
    CONSTRAINT uq_sessions_refresh_token_hash UNIQUE (refresh_token_hash),
    CONSTRAINT chk_sessions_status CHECK (status IN ('ACTIVE', 'REVOKED', 'EXPIRED')),
    CONSTRAINT chk_sessions_refresh_token_hash CHECK (length(btrim(refresh_token_hash)) > 0),
    CONSTRAINT chk_sessions_expiry CHECK (expires_at > created_at),
    CONSTRAINT chk_sessions_revocation CHECK (
        (status = 'REVOKED' AND revoked_at IS NOT NULL)
        OR (status <> 'REVOKED' AND revoked_at IS NULL)
    )
);

CREATE INDEX idx_sessions_account_status
    ON identity.sessions(account_id, status);

CREATE INDEX idx_sessions_expires_at
    ON identity.sessions(expires_at);

-- +goose Down
DROP TABLE IF EXISTS identity.sessions;

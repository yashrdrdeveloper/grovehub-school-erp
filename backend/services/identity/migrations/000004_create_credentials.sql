-- +goose Up
CREATE TABLE identity.credentials (
    credential_id uuid NOT NULL,
    account_id uuid NOT NULL,
    type text NOT NULL DEFAULT 'PASSWORD',
    secret_hash text NOT NULL,
    algorithm text NOT NULL DEFAULT 'ARGON2ID',
    status text NOT NULL DEFAULT 'ACTIVE',
    last_used_at timestamptz NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT pk_credentials PRIMARY KEY (credential_id),
    CONSTRAINT fk_credentials_account FOREIGN KEY (account_id)
        REFERENCES identity.accounts(account_id) ON DELETE CASCADE,
    CONSTRAINT chk_credentials_type CHECK (type IN ('PASSWORD')),
    CONSTRAINT chk_credentials_algorithm CHECK (algorithm IN ('ARGON2ID')),
    CONSTRAINT chk_credentials_status CHECK (status IN ('ACTIVE', 'REVOKED')),
    CONSTRAINT chk_credentials_secret_hash CHECK (length(btrim(secret_hash)) > 0)
);

CREATE INDEX idx_credentials_account_id
    ON identity.credentials(account_id);

CREATE UNIQUE INDEX uq_credentials_account_active_type
    ON identity.credentials(account_id, type)
    WHERE status = 'ACTIVE';

-- +goose Down
DROP TABLE IF EXISTS identity.credentials;

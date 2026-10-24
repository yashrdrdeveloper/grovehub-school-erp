-- +goose Up
CREATE TABLE identity.password_credentials (
    account_id uuid NOT NULL,
    password_hash text NOT NULL,
    password_changed_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version bigint NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT pk_password_credentials PRIMARY KEY (account_id),
    CONSTRAINT fk_password_credentials_account
        FOREIGN KEY (account_id)
        REFERENCES identity.accounts (account_id)
        ON DELETE CASCADE,
    CONSTRAINT chk_password_credentials_hash CHECK (length(btrim(password_hash)) > 0),
    CONSTRAINT chk_password_credentials_version CHECK (version >= 1)
);

-- Password credentials are internal Identity state. They are not exposed as an HTTP resource.
-- Account deletion cascades to the credential row.

-- +goose Down
DROP TABLE IF EXISTS identity.password_credentials;

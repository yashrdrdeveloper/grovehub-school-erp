-- +goose Up
CREATE TABLE identity.login_identifiers (
    identifier_id uuid NOT NULL,
    account_id uuid NOT NULL,
    type text NOT NULL,
    value text NOT NULL,
    normalized_value text NOT NULL,
    is_primary boolean NOT NULL DEFAULT false,
    is_verified boolean NOT NULL DEFAULT false,
    status text NOT NULL DEFAULT 'ACTIVE',
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT pk_login_identifiers PRIMARY KEY (identifier_id),
    CONSTRAINT fk_login_identifiers_account FOREIGN KEY (account_id)
        REFERENCES identity.accounts(account_id) ON DELETE CASCADE,
    CONSTRAINT uq_login_identifiers_type_normalized_value UNIQUE (type, normalized_value),
    CONSTRAINT chk_login_identifiers_type CHECK (type IN ('EMAIL', 'PHONE')),
    CONSTRAINT chk_login_identifiers_status CHECK (status IN ('ACTIVE', 'DISABLED')),
    CONSTRAINT chk_login_identifiers_value CHECK (length(btrim(value)) > 0),
    CONSTRAINT chk_login_identifiers_normalized_value CHECK (length(btrim(normalized_value)) > 0)
);

CREATE INDEX idx_login_identifiers_account_id
    ON identity.login_identifiers(account_id);

CREATE UNIQUE INDEX uq_login_identifiers_account_primary_type
    ON identity.login_identifiers(account_id, type)
    WHERE is_primary = true AND status = 'ACTIVE';

-- +goose Down
DROP TABLE IF EXISTS identity.login_identifiers;

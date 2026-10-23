-- +goose Up
CREATE TABLE identity.accounts (
    account_id uuid NOT NULL,
    person_gid uuid NOT NULL,
    status text NOT NULL DEFAULT 'ACTIVE',
    version bigint NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT pk_accounts PRIMARY KEY (account_id),
    CONSTRAINT uq_accounts_person_gid UNIQUE (person_gid),
    CONSTRAINT chk_accounts_status CHECK (
        status IN ('PENDING_VERIFICATION', 'ACTIVE', 'LOCKED', 'SUSPENDED', 'DEACTIVATED')
    ),
    CONSTRAINT chk_accounts_version CHECK (version >= 1)
);

-- Identity intentionally does not declare a foreign key to foundation.persons.
-- Cross-domain Person existence is validated through the Foundation client.

-- +goose Down
DROP TABLE IF EXISTS identity.accounts;

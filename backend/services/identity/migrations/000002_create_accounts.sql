-- +goose Up
CREATE TABLE identity.accounts (
    account_id uuid NOT NULL,
    person_gid uuid NOT NULL,
    email text NOT NULL,
    phone text,
    status text NOT NULL DEFAULT 'ACTIVE',
    version bigint NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT pk_accounts PRIMARY KEY (account_id),
    CONSTRAINT uq_accounts_person_gid UNIQUE (person_gid),
    CONSTRAINT uq_accounts_email UNIQUE (email),
    CONSTRAINT chk_accounts_email CHECK (length(btrim(email)) > 0),
    CONSTRAINT chk_accounts_phone CHECK (phone IS NULL OR length(btrim(phone)) > 0),
    CONSTRAINT chk_accounts_status CHECK (status IN ('ACTIVE', 'SUSPENDED', 'DEACTIVATED')),
    CONSTRAINT chk_accounts_version CHECK (version >= 1)
);

CREATE UNIQUE INDEX uq_accounts_phone
    ON identity.accounts (phone)
    WHERE phone IS NOT NULL;

-- Identity intentionally does not declare a foreign key to foundation.persons.
-- Person existence is owned by Foundation and will be resolved through its public client.

-- +goose Down
DROP TABLE IF EXISTS identity.accounts;

-- +goose Up
DROP SCHEMA IF EXISTS identity CASCADE;
DROP TABLE IF EXISTS identity_schema_migrations;

CREATE SCHEMA identity;

CREATE TABLE identity.accounts (
    account_id uuid NOT NULL,
    registration_idempotency_key text,
    registration_request_hash text,
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
    CONSTRAINT chk_accounts_registration_idempotency CHECK (
        (registration_idempotency_key IS NULL AND registration_request_hash IS NULL)
        OR (
            length(btrim(registration_idempotency_key)) > 0
            AND length(registration_request_hash) = 64
        )
    ),
    CONSTRAINT chk_accounts_email CHECK (length(btrim(email)) > 0),
    CONSTRAINT chk_accounts_phone CHECK (phone IS NULL OR length(btrim(phone)) > 0),
    CONSTRAINT chk_accounts_status CHECK (status IN ('ACTIVE', 'SUSPENDED', 'DEACTIVATED')),
    CONSTRAINT chk_accounts_version CHECK (version >= 1)
);

CREATE UNIQUE INDEX uq_accounts_registration_idempotency_key
    ON identity.accounts (registration_idempotency_key)
    WHERE registration_idempotency_key IS NOT NULL;

CREATE UNIQUE INDEX uq_accounts_phone
    ON identity.accounts (phone)
    WHERE phone IS NOT NULL;

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

-- Identity intentionally does not declare a foreign key from accounts.person_gid
-- to foundation.persons. Person existence is owned by Foundation and resolved
-- through the Foundation client boundary.

-- +goose Down
DROP SCHEMA IF EXISTS identity CASCADE;

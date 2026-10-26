-- +goose Up
DROP SCHEMA IF EXISTS foundation CASCADE;
DROP TABLE IF EXISTS foundation_schema_migrations;

CREATE SCHEMA foundation;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE foundation.persons (
    gid uuid NOT NULL,
    first_name text NOT NULL,
    middle_name text,
    last_name text NOT NULL,
    status text NOT NULL DEFAULT 'ACTIVE',
    version bigint NOT NULL DEFAULT 1,
    search_details text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT pk_persons PRIMARY KEY (gid),
    CONSTRAINT chk_persons_status CHECK (status IN ('ACTIVE', 'DEACTIVATED')),
    CONSTRAINT chk_persons_version CHECK (version >= 1)
);

CREATE TABLE foundation.person_contacts (
    person_gid uuid NOT NULL,
    type text NOT NULL,
    value text NOT NULL,
    normalized_value text NOT NULL,
    is_primary boolean NOT NULL DEFAULT FALSE,
    is_verified boolean NOT NULL DEFAULT FALSE,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT pk_person_contacts PRIMARY KEY (person_gid, type, normalized_value),
    CONSTRAINT fk_person_contacts_person
        FOREIGN KEY (person_gid)
        REFERENCES foundation.persons (gid)
        ON DELETE CASCADE,
    CONSTRAINT chk_person_contacts_type CHECK (type IN ('EMAIL', 'PHONE')),
    CONSTRAINT chk_person_contacts_value CHECK (length(btrim(value)) > 0),
    CONSTRAINT chk_person_contacts_normalized_value CHECK (length(btrim(normalized_value)) > 0)
);

CREATE UNIQUE INDEX uq_person_contacts_type_normalized_value
    ON foundation.person_contacts (type, normalized_value);

CREATE UNIQUE INDEX uq_person_contacts_primary_type
    ON foundation.person_contacts (person_gid, type)
    WHERE is_primary;

CREATE INDEX idx_persons_search_details_trgm
    ON foundation.persons
    USING GIN (search_details gin_trgm_ops);

-- +goose Down
DROP SCHEMA IF EXISTS foundation CASCADE;

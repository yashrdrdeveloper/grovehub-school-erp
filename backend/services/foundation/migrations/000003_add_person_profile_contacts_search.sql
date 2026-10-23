-- +goose Up
CREATE EXTENSION IF NOT EXISTS pg_trgm;

ALTER TABLE foundation.persons
    ADD COLUMN first_name text NOT NULL DEFAULT '',
    ADD COLUMN middle_name text,
    ADD COLUMN last_name text NOT NULL DEFAULT '',
    ADD COLUMN search_details text NOT NULL DEFAULT '';

ALTER TABLE foundation.persons
    ALTER COLUMN first_name DROP DEFAULT,
    ALTER COLUMN last_name DROP DEFAULT;

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

UPDATE foundation.persons
SET search_details = LOWER(
    TRIM(
        CONCAT_WS(
            ' ',
            gid::text,
            first_name,
            middle_name,
            last_name
        )
    )
);

CREATE INDEX idx_persons_search_details_trgm
    ON foundation.persons
    USING GIN (search_details gin_trgm_ops);

-- +goose Down
DROP INDEX IF EXISTS foundation.idx_persons_search_details_trgm;
DROP INDEX IF EXISTS foundation.uq_person_contacts_primary_type;
DROP INDEX IF EXISTS foundation.uq_person_contacts_type_normalized_value;
DROP TABLE IF EXISTS foundation.person_contacts;

ALTER TABLE foundation.persons
    DROP COLUMN IF EXISTS search_details,
    DROP COLUMN IF EXISTS last_name,
    DROP COLUMN IF EXISTS middle_name,
    DROP COLUMN IF EXISTS first_name;

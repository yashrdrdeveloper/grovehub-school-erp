-- +goose Up
CREATE TABLE foundation.persons (
    gid uuid NOT NULL,
    status text NOT NULL DEFAULT 'ACTIVE',
    version bigint NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT pk_persons PRIMARY KEY (gid),
    CONSTRAINT chk_persons_status CHECK (status IN ('ACTIVE', 'DEACTIVATED')),
    CONSTRAINT chk_persons_version CHECK (version >= 1)
);

-- +goose Down
DROP TABLE IF EXISTS foundation.persons;

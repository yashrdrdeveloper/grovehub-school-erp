-- +goose Up
CREATE TABLE mail_provider.domains (
    domain_id uuid NOT NULL,
    name text NOT NULL,
    verification_token text NOT NULL,
    status text NOT NULL DEFAULT 'PENDING_VERIFICATION',
    version bigint NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT pk_mail_provider_domains PRIMARY KEY (domain_id),
    CONSTRAINT uq_mail_provider_domains_name UNIQUE (name),
    CONSTRAINT chk_mail_provider_domains_name CHECK (length(btrim(name)) > 0),
    CONSTRAINT chk_mail_provider_domains_status CHECK (status IN ('PENDING_VERIFICATION', 'VERIFIED', 'DEACTIVATED')),
    CONSTRAINT chk_mail_provider_domains_version CHECK (version >= 1)
);

-- +goose Down
DROP TABLE IF EXISTS mail_provider.domains;

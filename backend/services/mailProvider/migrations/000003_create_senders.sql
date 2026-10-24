-- +goose Up
CREATE TABLE mail_provider.senders (
    sender_id uuid NOT NULL,
    domain_id uuid NOT NULL,
    email text NOT NULL,
    display_name text NOT NULL,
    status text NOT NULL DEFAULT 'ACTIVE',
    version bigint NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT pk_mail_provider_senders PRIMARY KEY (sender_id),
    CONSTRAINT fk_mail_provider_senders_domain FOREIGN KEY (domain_id) REFERENCES mail_provider.domains(domain_id),
    CONSTRAINT uq_mail_provider_senders_email UNIQUE (email),
    CONSTRAINT chk_mail_provider_senders_email CHECK (length(btrim(email)) > 0),
    CONSTRAINT chk_mail_provider_senders_display_name CHECK (length(btrim(display_name)) > 0),
    CONSTRAINT chk_mail_provider_senders_status CHECK (status IN ('ACTIVE', 'DEACTIVATED')),
    CONSTRAINT chk_mail_provider_senders_version CHECK (version >= 1)
);

CREATE INDEX idx_mail_provider_senders_domain_id ON mail_provider.senders(domain_id);

-- +goose Down
DROP TABLE IF EXISTS mail_provider.senders;

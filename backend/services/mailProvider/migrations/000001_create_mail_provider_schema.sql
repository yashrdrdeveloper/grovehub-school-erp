-- +goose Up
DROP SCHEMA IF EXISTS mail_provider CASCADE;
DROP TABLE IF EXISTS mail_provider_schema_migrations;

CREATE SCHEMA mail_provider;

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
    CONSTRAINT fk_mail_provider_senders_domain
        FOREIGN KEY (domain_id)
        REFERENCES mail_provider.domains(domain_id),
    CONSTRAINT uq_mail_provider_senders_email UNIQUE (email),
    CONSTRAINT chk_mail_provider_senders_email CHECK (length(btrim(email)) > 0),
    CONSTRAINT chk_mail_provider_senders_display_name CHECK (length(btrim(display_name)) > 0),
    CONSTRAINT chk_mail_provider_senders_status CHECK (status IN ('ACTIVE', 'DEACTIVATED')),
    CONSTRAINT chk_mail_provider_senders_version CHECK (version >= 1)
);

CREATE INDEX idx_mail_provider_senders_domain_id
    ON mail_provider.senders(domain_id);

CREATE TABLE mail_provider.messages (
    message_id uuid NOT NULL,
    sender_id uuid NOT NULL,
    to_addresses text[] NOT NULL,
    cc_addresses text[] NOT NULL DEFAULT '{}',
    bcc_addresses text[] NOT NULL DEFAULT '{}',
    subject text NOT NULL,
    text_body text NOT NULL DEFAULT '',
    html_body text NOT NULL DEFAULT '',
    status text NOT NULL DEFAULT 'QUEUED',
    failure_message text,
    version bigint NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT pk_mail_provider_messages PRIMARY KEY (message_id),
    CONSTRAINT fk_mail_provider_messages_sender
        FOREIGN KEY (sender_id)
        REFERENCES mail_provider.senders(sender_id),
    CONSTRAINT chk_mail_provider_messages_to CHECK (cardinality(to_addresses) > 0),
    CONSTRAINT chk_mail_provider_messages_subject CHECK (length(btrim(subject)) > 0),
    CONSTRAINT chk_mail_provider_messages_body CHECK (length(btrim(text_body)) > 0 OR length(btrim(html_body)) > 0),
    CONSTRAINT chk_mail_provider_messages_status CHECK (status IN ('QUEUED', 'SUBMITTED', 'FAILED')),
    CONSTRAINT chk_mail_provider_messages_version CHECK (version >= 1)
);

CREATE INDEX idx_mail_provider_messages_sender_id
    ON mail_provider.messages(sender_id);

CREATE INDEX idx_mail_provider_messages_created_at
    ON mail_provider.messages(created_at DESC);

-- +goose Down
DROP SCHEMA IF EXISTS mail_provider CASCADE;

-- +goose Up
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
    CONSTRAINT fk_mail_provider_messages_sender FOREIGN KEY (sender_id) REFERENCES mail_provider.senders(sender_id),
    CONSTRAINT chk_mail_provider_messages_to CHECK (cardinality(to_addresses) > 0),
    CONSTRAINT chk_mail_provider_messages_subject CHECK (length(btrim(subject)) > 0),
    CONSTRAINT chk_mail_provider_messages_body CHECK (length(btrim(text_body)) > 0 OR length(btrim(html_body)) > 0),
    CONSTRAINT chk_mail_provider_messages_status CHECK (status IN ('QUEUED', 'SUBMITTED', 'FAILED')),
    CONSTRAINT chk_mail_provider_messages_version CHECK (version >= 1)
);

CREATE INDEX idx_mail_provider_messages_sender_id ON mail_provider.messages(sender_id);
CREATE INDEX idx_mail_provider_messages_created_at ON mail_provider.messages(created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS mail_provider.messages;

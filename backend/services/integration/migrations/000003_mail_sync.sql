-- +goose Up
CREATE TABLE integration.mail_references (
 id uuid PRIMARY KEY,
 connection_id uuid NOT NULL REFERENCES integration.connections(id) ON DELETE CASCADE,
 provider_message_id text NOT NULL,
 thread_id text NOT NULL,
 direction text NOT NULL CHECK (direction IN ('INBOUND','OUTBOUND')),
 history_id text NOT NULL DEFAULT '',
 labels text NOT NULL DEFAULT '',
 deleted boolean NOT NULL DEFAULT false,
 updated_at timestamptz NOT NULL DEFAULT now(),
 CONSTRAINT uq_mail_ref UNIQUE(connection_id,provider_message_id)
);
CREATE INDEX idx_mail_ref_connection_updated ON integration.mail_references(connection_id,updated_at DESC);
-- +goose Down
DROP TABLE IF EXISTS integration.mail_references;

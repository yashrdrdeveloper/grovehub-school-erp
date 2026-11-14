-- +goose Up
-- Run after the existing Foundation history: old queued events are copied
-- before the legacy outbox is retired. The same PostgreSQL database is required
-- so business writes and the Infrastructure message record commit together.
CREATE SCHEMA IF NOT EXISTS infrastructure;

CREATE TABLE infrastructure.messaging (
    event_id text PRIMARY KEY,
    topic text NOT NULL,
    payload bytea NOT NULL,
    metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
    status text NOT NULL DEFAULT 'PENDING',
    attempts integer NOT NULL DEFAULT 0,
    next_attempt_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    lease_until timestamptz,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    published_at timestamptz,
    CONSTRAINT chk_infrastructure_messaging_status CHECK (status IN ('PENDING', 'SENDING', 'SENT')),
    CONSTRAINT chk_infrastructure_messaging_attempts CHECK (attempts >= 0)
);

CREATE INDEX idx_infrastructure_messaging_pending
    ON infrastructure.messaging (next_attempt_at, created_at)
    WHERE status = 'PENDING';

CREATE INDEX idx_infrastructure_messaging_leases
    ON infrastructure.messaging (lease_until)
    WHERE status = 'SENDING';

CREATE INDEX idx_infrastructure_messaging_sent_retention
    ON infrastructure.messaging (published_at, event_id)
    WHERE status = 'SENT';

-- +goose StatementBegin
DO $$
BEGIN
    IF to_regclass('foundation.event_outbox') IS NOT NULL THEN
        INSERT INTO infrastructure.messaging
            (event_id, topic, payload, metadata, status, attempts,
             next_attempt_at, lease_until, created_at, published_at)
        SELECT event_id, topic, payload,
               COALESCE(convert_from(metadata, 'UTF8')::jsonb, '{}'::jsonb),
               status, attempts, next_attempt_at, lease_until, created_at,
               published_at
        FROM foundation.event_outbox
        ON CONFLICT (event_id) DO NOTHING;

        -- This transactional migration only drops the legacy table after the
        -- copy succeeds. An active old dispatcher must be stopped first.
        DROP TABLE foundation.event_outbox;
    END IF;
END
$$;
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS infrastructure.messaging;

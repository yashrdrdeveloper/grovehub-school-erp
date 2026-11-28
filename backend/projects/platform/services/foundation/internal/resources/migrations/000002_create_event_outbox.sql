-- +goose Up
-- Each aggregate mutation and its event enqueue are committed atomically.
CREATE TABLE foundation.event_outbox (
    event_id text PRIMARY KEY,
    topic text NOT NULL,
    payload bytea NOT NULL,
    metadata bytea NOT NULL,
    status text NOT NULL DEFAULT 'PENDING',
    attempts integer NOT NULL DEFAULT 0,
    next_attempt_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    lease_until timestamptz,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    published_at timestamptz,
    CONSTRAINT chk_foundation_event_outbox_status CHECK (status IN ('PENDING', 'SENDING', 'SENT')),
    CONSTRAINT chk_foundation_event_outbox_attempts CHECK (attempts >= 0)
);

CREATE INDEX idx_foundation_event_outbox_pending
    ON foundation.event_outbox (next_attempt_at, created_at)
    WHERE status = 'PENDING';

CREATE INDEX idx_foundation_event_outbox_leases
    ON foundation.event_outbox (lease_until)
    WHERE status = 'SENDING';

-- +goose Down
DROP TABLE IF EXISTS foundation.event_outbox;

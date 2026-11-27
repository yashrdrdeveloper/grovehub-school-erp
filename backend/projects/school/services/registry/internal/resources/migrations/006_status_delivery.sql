CREATE TABLE school_registry.status_outbox (
 execution_id uuid PRIMARY KEY REFERENCES school_registry.status_events(execution_id),
 attempts bigint NOT NULL DEFAULT 0 CHECK(attempts>=0),
 next_attempt_at timestamptz NOT NULL DEFAULT now(),
 delivered_at timestamptz,
 blocked_at timestamptz,
 CHECK(delivered_at IS NULL OR blocked_at IS NULL)
);
CREATE INDEX status_outbox_pending ON school_registry.status_outbox(next_attempt_at,execution_id) WHERE delivered_at IS NULL AND blocked_at IS NULL;
CREATE FUNCTION school_registry.enqueue_status_receipt() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
 INSERT INTO school_registry.status_outbox(execution_id) VALUES(NEW.execution_id);
 RETURN NEW;
END;
$$;
CREATE TRIGGER status_outbox_enqueue AFTER INSERT ON school_registry.status_events
 FOR EACH ROW EXECUTE FUNCTION school_registry.enqueue_status_receipt();
INSERT INTO school_registry.status_outbox(execution_id)
 SELECT execution_id FROM school_registry.status_events ON CONFLICT DO NOTHING;

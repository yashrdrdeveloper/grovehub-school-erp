ALTER TABLE school_registry.rename_events ADD COLUMN correlation_id uuid;
CREATE TABLE school_registry.rename_outbox (
 execution_id uuid PRIMARY KEY REFERENCES school_registry.rename_events(execution_id),
 attempts bigint NOT NULL DEFAULT 0 CHECK(attempts>=0),
 next_attempt_at timestamptz NOT NULL DEFAULT now(),
 delivered_at timestamptz,
 blocked_at timestamptz,
 CHECK(delivered_at IS NULL OR blocked_at IS NULL)
);
CREATE INDEX rename_outbox_pending ON school_registry.rename_outbox(next_attempt_at,execution_id) WHERE delivered_at IS NULL AND blocked_at IS NULL;
CREATE FUNCTION school_registry.enqueue_rename_receipt() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
 IF NEW.previous_version IS NOT NULL THEN
  INSERT INTO school_registry.rename_outbox(execution_id) VALUES(NEW.execution_id);
 END IF;
 RETURN NEW;
END;
$$;
CREATE TRIGGER rename_outbox_enqueue AFTER INSERT ON school_registry.rename_events
 FOR EACH ROW EXECUTE FUNCTION school_registry.enqueue_rename_receipt();
INSERT INTO school_registry.rename_outbox(execution_id)
 SELECT execution_id FROM school_registry.rename_events WHERE previous_version IS NOT NULL ON CONFLICT DO NOTHING;

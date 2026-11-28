CREATE TABLE internal_tool.audit_outbox (
 workflow_event_id bigint PRIMARY KEY REFERENCES internal_tool.workflow_events(id),
 event_id uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
 attempts integer NOT NULL DEFAULT 0 CHECK(attempts>=0),
 next_attempt_at timestamptz NOT NULL DEFAULT now(),
 delivered_at timestamptz,
 blocked_at timestamptz,
 CHECK(delivered_at IS NULL OR blocked_at IS NULL)
);
CREATE INDEX audit_outbox_pending ON internal_tool.audit_outbox(next_attempt_at,workflow_event_id) WHERE delivered_at IS NULL AND blocked_at IS NULL;
CREATE FUNCTION internal_tool.enqueue_audit_event() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
 INSERT INTO internal_tool.audit_outbox(workflow_event_id) VALUES(NEW.id);
 RETURN NEW;
END;
$$;
CREATE TRIGGER workflow_audit_outbox AFTER INSERT ON internal_tool.workflow_events FOR EACH ROW EXECUTE FUNCTION internal_tool.enqueue_audit_event();
INSERT INTO internal_tool.audit_outbox(workflow_event_id) SELECT id FROM internal_tool.workflow_events ON CONFLICT DO NOTHING;

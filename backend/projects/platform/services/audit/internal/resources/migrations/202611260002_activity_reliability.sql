-- +goose Up
ALTER TABLE audit.activities ADD COLUMN subject_account_id uuid;
CREATE INDEX idx_audit_subject_account_time ON audit.activities(tenant_id, subject_account_id, occurred_at DESC, audit_id DESC);
CREATE TABLE audit.rejected_events (
    fingerprint char(64) PRIMARY KEY,
    event_type text NOT NULL,
    reason text NOT NULL CHECK (reason IN ('INVALID_EVENT','IDENTITY_CONFLICT')),
    payload_bytes bigint NOT NULL CHECK (payload_bytes >= 0),
    observed_at timestamptz NOT NULL DEFAULT clock_timestamp()
);
CREATE TRIGGER rejected_events_immutable BEFORE UPDATE OR DELETE ON audit.rejected_events
FOR EACH ROW EXECUTE FUNCTION audit.reject_mutation();
CREATE TRIGGER rejected_events_no_truncate BEFORE TRUNCATE ON audit.rejected_events
FOR EACH STATEMENT EXECUTE FUNCTION audit.reject_mutation();

-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM audit.rejected_events) OR EXISTS (SELECT 1 FROM audit.activities WHERE subject_account_id IS NOT NULL) THEN
        RAISE EXCEPTION 'audit diagnostic or subject history must be preserved';
    END IF;
END;
$$;
-- +goose StatementEnd
DROP TABLE audit.rejected_events;
DROP INDEX audit.idx_audit_subject_account_time;
ALTER TABLE audit.activities DROP COLUMN subject_account_id;

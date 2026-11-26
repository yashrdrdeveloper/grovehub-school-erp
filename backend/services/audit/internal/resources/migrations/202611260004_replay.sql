-- +goose Up
CREATE TABLE audit.event_receipts (
    source_event_id uuid PRIMARY KEY,
    audit_id uuid NOT NULL REFERENCES audit.activities(audit_id),
    received_at timestamptz NOT NULL DEFAULT clock_timestamp()
);
INSERT INTO audit.event_receipts(source_event_id,audit_id)
SELECT source_event_id,audit_id FROM audit.activities;
CREATE INDEX idx_audit_receipts_activity ON audit.event_receipts(audit_id);
CREATE TRIGGER event_receipts_immutable BEFORE UPDATE OR DELETE ON audit.event_receipts
FOR EACH ROW EXECUTE FUNCTION audit.reject_mutation();
CREATE TRIGGER event_receipts_no_truncate BEFORE TRUNCATE ON audit.event_receipts
FOR EACH STATEMENT EXECUTE FUNCTION audit.reject_mutation();

CREATE TABLE audit.replay_runs (
    replay_id uuid PRIMARY KEY,
    tenant_id text NOT NULL CHECK (tenant_id <> ''),
    source_event_id uuid NOT NULL REFERENCES audit.event_receipts(source_event_id),
    audit_id uuid NOT NULL REFERENCES audit.activities(audit_id),
    operator_label text NOT NULL CHECK (length(operator_label) BETWEEN 1 AND 128),
    replayed_at timestamptz NOT NULL DEFAULT clock_timestamp()
);
CREATE INDEX idx_audit_replay_tenant_time ON audit.replay_runs(tenant_id,replayed_at DESC);
CREATE TRIGGER replay_runs_immutable BEFORE UPDATE OR DELETE ON audit.replay_runs
FOR EACH ROW EXECUTE FUNCTION audit.reject_mutation();
CREATE TRIGGER replay_runs_no_truncate BEFORE TRUNCATE ON audit.replay_runs
FOR EACH STATEMENT EXECUTE FUNCTION audit.reject_mutation();

-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM audit.replay_runs)
       OR EXISTS (SELECT 1 FROM audit.event_receipts r JOIN audit.activities a USING (audit_id) WHERE r.source_event_id <> a.source_event_id) THEN
        RAISE EXCEPTION 'audit replay and receipt history must be preserved';
    END IF;
END;
$$;
-- +goose StatementEnd
DROP TABLE audit.replay_runs;
DROP TABLE audit.event_receipts;

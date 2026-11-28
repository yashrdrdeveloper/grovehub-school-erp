-- +goose Up
CREATE TABLE audit.operator_activities (
 event_id uuid PRIMARY KEY,
 source text NOT NULL CHECK(source='internal-tool'),
 actor_id uuid NOT NULL,
 resource_id uuid NOT NULL,
 operation text NOT NULL,
 occurred_at timestamptz NOT NULL,
 ingested_at timestamptz NOT NULL DEFAULT clock_timestamp()
);
CREATE INDEX idx_operator_activity_actor ON audit.operator_activities(actor_id,occurred_at DESC);
CREATE INDEX idx_operator_activity_resource ON audit.operator_activities(resource_id,occurred_at DESC);
CREATE TRIGGER operator_activities_immutable BEFORE UPDATE OR DELETE ON audit.operator_activities
FOR EACH ROW EXECUTE FUNCTION audit.reject_mutation();
CREATE TRIGGER operator_activities_no_truncate BEFORE TRUNCATE ON audit.operator_activities
FOR EACH STATEMENT EXECUTE FUNCTION audit.reject_mutation();
-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
 RAISE EXCEPTION 'operator audit history must be preserved';
END;
$$;
-- +goose StatementEnd

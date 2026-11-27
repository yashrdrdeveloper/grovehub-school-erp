-- +goose Up
CREATE TABLE audit.school_mutation_receipts (
 execution_id uuid PRIMARY KEY,
 school_id uuid NOT NULL,
 previous_version bigint NOT NULL CHECK(previous_version>0),
 correlation_id uuid,
 occurred_at timestamptz NOT NULL,
 ingested_at timestamptz NOT NULL DEFAULT clock_timestamp()
);
CREATE INDEX school_mutation_receipts_school ON audit.school_mutation_receipts(school_id,occurred_at DESC);
CREATE INDEX school_mutation_receipts_correlation ON audit.school_mutation_receipts(correlation_id) WHERE correlation_id IS NOT NULL;
CREATE TRIGGER school_mutation_receipts_immutable BEFORE UPDATE OR DELETE ON audit.school_mutation_receipts FOR EACH ROW EXECUTE FUNCTION audit.reject_mutation();
CREATE TRIGGER school_mutation_receipts_no_truncate BEFORE TRUNCATE ON audit.school_mutation_receipts FOR EACH STATEMENT EXECUTE FUNCTION audit.reject_mutation();
-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
 RAISE EXCEPTION 'school mutation receipts must be preserved';
END;
$$;
-- +goose StatementEnd

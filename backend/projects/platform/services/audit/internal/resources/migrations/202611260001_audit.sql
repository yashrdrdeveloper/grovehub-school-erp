-- +goose Up
CREATE SCHEMA IF NOT EXISTS audit;
CREATE TABLE audit.activities (
    audit_id uuid PRIMARY KEY,
    source_event_id uuid NOT NULL UNIQUE,
    tenant_id text NOT NULL CHECK (tenant_id <> ''),
    asset_type text NOT NULL,
    asset_id uuid NOT NULL,
    audit_type text NOT NULL,
    operation text NOT NULL,
    source text NOT NULL,
    asset_version bigint NOT NULL CHECK (asset_version > 0),
    actor_id uuid,
    person_id uuid,
    data jsonb NOT NULL CHECK (jsonb_typeof(data) = 'object'),
    semantic_hash char(64) NOT NULL,
    occurred_at timestamptz NOT NULL,
    ingested_at timestamptz NOT NULL DEFAULT clock_timestamp(),
    UNIQUE (tenant_id, asset_type, asset_id, audit_type, operation, asset_version)
);
CREATE INDEX idx_audit_tenant_time ON audit.activities(tenant_id, occurred_at DESC, audit_id DESC);
CREATE INDEX idx_audit_asset_time ON audit.activities(tenant_id, asset_type, asset_id, occurred_at DESC, audit_id DESC);
CREATE INDEX idx_audit_actor_time ON audit.activities(tenant_id, actor_id, occurred_at DESC, audit_id DESC);
CREATE INDEX idx_audit_person_time ON audit.activities(tenant_id, person_id, occurred_at DESC, audit_id DESC);
-- +goose StatementBegin
CREATE FUNCTION audit.reject_mutation() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    RAISE EXCEPTION 'audit activities are immutable';
END;
$$;
-- +goose StatementEnd
CREATE TRIGGER activities_immutable BEFORE UPDATE OR DELETE ON audit.activities
FOR EACH ROW EXECUTE FUNCTION audit.reject_mutation();
CREATE TRIGGER activities_no_truncate BEFORE TRUNCATE ON audit.activities
FOR EACH STATEMENT EXECUTE FUNCTION audit.reject_mutation();

-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM audit.activities) THEN
        RAISE EXCEPTION 'audit history must be preserved';
    END IF;
END;
$$;
-- +goose StatementEnd
DROP TABLE audit.activities;
DROP FUNCTION audit.reject_mutation();
DROP SCHEMA audit;

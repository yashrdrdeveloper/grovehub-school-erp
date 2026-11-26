-- +goose Up
ALTER TABLE audit.activities ADD COLUMN occurrence_id uuid;
ALTER TABLE audit.activities DROP CONSTRAINT activities_asset_version_check;
-- The original unnamed UNIQUE constraint may have a truncated generated name.
-- +goose StatementBegin
DO $$
DECLARE constraint_name text;
BEGIN
    SELECT conname INTO STRICT constraint_name
    FROM pg_constraint
    WHERE conrelid = 'audit.activities'::regclass AND contype = 'u'
      AND pg_get_constraintdef(oid) = 'UNIQUE (tenant_id, asset_type, asset_id, audit_type, operation, asset_version)';
    EXECUTE format('ALTER TABLE audit.activities DROP CONSTRAINT %I', constraint_name);
END;
$$;
-- +goose StatementEnd
ALTER TABLE audit.activities ADD CONSTRAINT activities_identity_check CHECK (
    (asset_version > 0 AND occurrence_id IS NULL) OR (asset_version = 0 AND occurrence_id IS NOT NULL)
);
CREATE UNIQUE INDEX idx_audit_version_identity ON audit.activities(tenant_id,asset_type,asset_id,audit_type,operation,asset_version) WHERE occurrence_id IS NULL;
CREATE UNIQUE INDEX idx_audit_occurrence_identity ON audit.activities(tenant_id,source,asset_type,asset_id,audit_type,operation,occurrence_id) WHERE occurrence_id IS NOT NULL;

-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM audit.activities WHERE occurrence_id IS NOT NULL) THEN
        RAISE EXCEPTION 'occurrence history must be preserved';
    END IF;
END;
$$;
-- +goose StatementEnd
DROP INDEX audit.idx_audit_occurrence_identity;
DROP INDEX audit.idx_audit_version_identity;
ALTER TABLE audit.activities DROP CONSTRAINT activities_identity_check;
ALTER TABLE audit.activities DROP COLUMN occurrence_id;
ALTER TABLE audit.activities ADD CONSTRAINT activities_asset_version_check CHECK (asset_version > 0);
ALTER TABLE audit.activities ADD UNIQUE (tenant_id,asset_type,asset_id,audit_type,operation,asset_version);

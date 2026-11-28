-- +goose Up
SET LOCAL lock_timeout = '5s';
CREATE INDEX idx_audit_tenant_ingested ON audit.activities(tenant_id, ingested_at DESC, audit_id DESC);
CREATE INDEX idx_audit_tenant_operation_time ON audit.activities(tenant_id, operation, occurred_at DESC, audit_id DESC);

-- +goose Down
SET LOCAL lock_timeout = '5s';
DROP INDEX audit.idx_audit_tenant_operation_time;
DROP INDEX audit.idx_audit_tenant_ingested;

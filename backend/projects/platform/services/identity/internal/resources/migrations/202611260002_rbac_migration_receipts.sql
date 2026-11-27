-- +goose Up
CREATE TABLE identity.rbac_migration_receipts (
    receipt_id uuid PRIMARY KEY,
    tenant_id text NOT NULL CHECK (length(btrim(tenant_id)) > 0),
    operator_id uuid NOT NULL,
    operation text NOT NULL CHECK (operation IN ('BACKFILL', 'BOOTSTRAP')),
    input_sha256 text NOT NULL CHECK (input_sha256 ~ '^[0-9a-f]{64}$'),
    result jsonb NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, operation)
);

-- +goose Down
DROP TABLE identity.rbac_migration_receipts;

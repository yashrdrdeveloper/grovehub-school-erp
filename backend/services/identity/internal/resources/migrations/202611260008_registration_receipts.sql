-- +goose Up
CREATE TABLE identity.registration_receipts (
    tenant_id text NOT NULL,
    actor_id uuid NOT NULL,
    idempotency_key text NOT NULL CHECK (length(idempotency_key) BETWEEN 1 AND 128 AND btrim(idempotency_key) = idempotency_key),
    request_sha256 text NOT NULL CHECK (request_sha256 ~ '^[0-9a-f]{64}$'),
    account_id uuid NOT NULL,
    event_id uuid NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, actor_id, idempotency_key),
    UNIQUE (tenant_id, account_id),
    FOREIGN KEY (tenant_id, actor_id) REFERENCES identity.accounts(tenant_id, account_id) ON DELETE RESTRICT,
    FOREIGN KEY (tenant_id, account_id) REFERENCES identity.accounts(tenant_id, account_id) ON DELETE RESTRICT
);
CREATE TRIGGER registration_receipts_immutable BEFORE UPDATE OR DELETE ON identity.registration_receipts
FOR EACH ROW EXECUTE FUNCTION identity.reject_resource_audit_mutation();

-- +goose Down
DROP TABLE identity.registration_receipts;

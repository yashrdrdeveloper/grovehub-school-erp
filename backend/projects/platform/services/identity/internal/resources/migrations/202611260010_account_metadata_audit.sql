-- +goose Up
CREATE TABLE identity.account_metadata_audit (
    event_id uuid PRIMARY KEY,
    tenant_id text NOT NULL,
    actor_id uuid NOT NULL,
    account_id uuid NOT NULL,
    account_version bigint NOT NULL CHECK (account_version > 1),
    operation text NOT NULL CHECK (operation = 'UPDATED'),
    event_payload jsonb NOT NULL CHECK (jsonb_typeof(event_payload) = 'object'),
    occurred_at timestamptz NOT NULL,
    UNIQUE (tenant_id, account_id, account_version),
    FOREIGN KEY (tenant_id, account_id) REFERENCES identity.accounts(tenant_id, account_id) ON DELETE RESTRICT,
    FOREIGN KEY (tenant_id, actor_id) REFERENCES identity.accounts(tenant_id, account_id) ON DELETE RESTRICT
);
CREATE INDEX idx_identity_account_metadata_audit_time ON identity.account_metadata_audit(tenant_id, occurred_at DESC, event_id);
CREATE TRIGGER account_metadata_audit_immutable BEFORE UPDATE OR DELETE ON identity.account_metadata_audit
FOR EACH ROW EXECUTE FUNCTION identity.reject_resource_audit_mutation();

-- +goose Down
DROP TABLE identity.account_metadata_audit;

-- +goose Up
CREATE TABLE identity.approval_audit (
    event_id uuid PRIMARY KEY,
    tenant_id text NOT NULL,
    request_id uuid NOT NULL,
    actor_id uuid NOT NULL,
    request_version bigint NOT NULL CHECK (request_version >= 1),
    operation text NOT NULL CHECK (operation IN ('CREATED','APPROVE','REJECT','CANCELLED','APPLIED','STALE','EXPIRED')),
    event_payload jsonb NOT NULL CHECK (jsonb_typeof(event_payload) = 'object'),
    occurred_at timestamptz NOT NULL,
    UNIQUE (tenant_id, request_id, request_version),
    FOREIGN KEY (tenant_id, request_id) REFERENCES identity.approval_requests(tenant_id, request_id) ON DELETE RESTRICT,
    FOREIGN KEY (tenant_id, actor_id) REFERENCES identity.accounts(tenant_id, account_id) ON DELETE RESTRICT
);

CREATE INDEX idx_identity_approval_audit_tenant_time ON identity.approval_audit(tenant_id, occurred_at, event_id);

-- +goose StatementBegin
CREATE FUNCTION identity.reject_approval_audit_mutation() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    RAISE EXCEPTION 'approval audit is immutable';
END;
$$;
-- +goose StatementEnd

CREATE TRIGGER approval_audit_immutable BEFORE UPDATE OR DELETE ON identity.approval_audit
FOR EACH ROW EXECUTE FUNCTION identity.reject_approval_audit_mutation();

-- +goose Down
DROP TABLE identity.approval_audit;
DROP FUNCTION identity.reject_approval_audit_mutation();

-- +goose Up
CREATE TABLE identity.resource_audit (
    event_id uuid PRIMARY KEY,
    tenant_id text NOT NULL,
    request_id uuid NOT NULL,
    actor_id uuid NOT NULL,
    resource_type text NOT NULL CHECK (resource_type IN ('ROLE','ACCOUNT')),
    resource_id uuid NOT NULL,
    resource_version bigint NOT NULL CHECK (resource_version > 1),
    operation text NOT NULL CHECK (operation IN ('ROLE_LIFECYCLE_CHANGE','ROLE_PERMISSION_CHANGE','ROLE_ASSIGNMENT_CHANGE','ACCOUNT_LIFECYCLE_CHANGE')),
    event_payload jsonb NOT NULL CHECK (jsonb_typeof(event_payload) = 'object'),
    occurred_at timestamptz NOT NULL,
    UNIQUE (tenant_id, request_id),
    UNIQUE (tenant_id, resource_type, resource_id, resource_version),
    FOREIGN KEY (tenant_id, request_id) REFERENCES identity.approval_requests(tenant_id, request_id) ON DELETE RESTRICT,
    FOREIGN KEY (tenant_id, actor_id) REFERENCES identity.accounts(tenant_id, account_id) ON DELETE RESTRICT
);
CREATE INDEX idx_identity_resource_audit_tenant_time ON identity.resource_audit(tenant_id, occurred_at DESC, event_id);
CREATE INDEX idx_identity_resource_audit_resource ON identity.resource_audit(tenant_id, resource_type, resource_id, occurred_at DESC, event_id);

-- +goose StatementBegin
CREATE FUNCTION identity.reject_resource_audit_mutation() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    RAISE EXCEPTION 'resource audit is immutable';
END;
$$;
-- +goose StatementEnd
CREATE TRIGGER resource_audit_immutable BEFORE UPDATE OR DELETE ON identity.resource_audit
FOR EACH ROW EXECUTE FUNCTION identity.reject_resource_audit_mutation();

-- +goose Down
DROP TABLE identity.resource_audit;
DROP FUNCTION identity.reject_resource_audit_mutation();

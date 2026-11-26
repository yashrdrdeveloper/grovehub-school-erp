-- +goose Up
CREATE TABLE identity.role_metadata_audit (
    event_id uuid PRIMARY KEY,
    tenant_id text NOT NULL,
    actor_id uuid NOT NULL,
    role_id uuid NOT NULL,
    role_version bigint NOT NULL CHECK (role_version >= 1),
    operation text NOT NULL CHECK (operation IN ('CREATED','UPDATED')),
    event_payload jsonb NOT NULL CHECK (jsonb_typeof(event_payload) = 'object'),
    occurred_at timestamptz NOT NULL,
    CHECK ((operation = 'CREATED' AND role_version = 1) OR (operation = 'UPDATED' AND role_version > 1)),
    UNIQUE (tenant_id, role_id, role_version),
    FOREIGN KEY (tenant_id, role_id) REFERENCES identity.roles(tenant_id, role_id) ON DELETE RESTRICT,
    FOREIGN KEY (tenant_id, actor_id) REFERENCES identity.accounts(tenant_id, account_id) ON DELETE RESTRICT
);
CREATE INDEX idx_identity_role_metadata_audit_time ON identity.role_metadata_audit(tenant_id, occurred_at DESC, event_id);
CREATE TRIGGER role_metadata_audit_immutable BEFORE UPDATE OR DELETE ON identity.role_metadata_audit
FOR EACH ROW EXECUTE FUNCTION identity.reject_resource_audit_mutation();

-- +goose Down
DROP TABLE identity.role_metadata_audit;

-- +goose Up
-- School-local approval data. No workflow is seeded or enabled by this migration.
CREATE TABLE identity.approval_policies (
    policy_id uuid PRIMARY KEY,
    tenant_id text NOT NULL CHECK (length(btrim(tenant_id)) > 0),
    request_type text NOT NULL CHECK (request_type IN ('ROLE_ACCESS','PERMISSION_ACCESS','ROLE_ASSIGNMENT_CHANGE','ROLE_PERMISSION_CHANGE','FEATURE_FLAG_CHANGE','APPROVAL_POLICY_CHANGE','ROLE_LIFECYCLE_CHANGE','ACCOUNT_LIFECYCLE_CHANGE')),
    version bigint NOT NULL CHECK (version >= 1),
    active boolean NOT NULL DEFAULT false,
    created_by uuid NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, policy_id),
    UNIQUE (tenant_id, request_type, version),
    FOREIGN KEY (tenant_id, created_by) REFERENCES identity.accounts(tenant_id, account_id) ON DELETE RESTRICT
);
CREATE UNIQUE INDEX uq_identity_approval_active_policy ON identity.approval_policies(tenant_id, request_type) WHERE active;

CREATE TABLE identity.approval_policy_steps (
    tenant_id text NOT NULL,
    policy_id uuid NOT NULL,
    step_number integer NOT NULL CHECK (step_number BETWEEN 1 AND 20),
    required_permission_code text NOT NULL REFERENCES identity.permissions(code) ON DELETE RESTRICT,
    required_role_id uuid,
    PRIMARY KEY (tenant_id, policy_id, step_number),
    FOREIGN KEY (tenant_id, policy_id) REFERENCES identity.approval_policies(tenant_id, policy_id) ON DELETE RESTRICT,
    FOREIGN KEY (tenant_id, required_role_id) REFERENCES identity.roles(tenant_id, role_id) ON DELETE RESTRICT
);

CREATE TABLE identity.approval_requests (
    request_id uuid PRIMARY KEY,
    tenant_id text NOT NULL,
    policy_id uuid NOT NULL,
    requester_id uuid NOT NULL,
    request_type text NOT NULL CHECK (request_type IN ('ROLE_ACCESS','PERMISSION_ACCESS','ROLE_ASSIGNMENT_CHANGE','ROLE_PERMISSION_CHANGE','FEATURE_FLAG_CHANGE','APPROVAL_POLICY_CHANGE','ROLE_LIFECYCLE_CHANGE','ACCOUNT_LIFECYCLE_CHANGE')),
    payload jsonb NOT NULL CHECK (jsonb_typeof(payload) = 'object'),
    payload_sha256 text NOT NULL CHECK (payload_sha256 ~ '^[0-9a-f]{64}$'),
    reason text NOT NULL CHECK (length(btrim(reason)) BETWEEN 1 AND 2000),
    idempotency_key text NOT NULL CHECK (length(btrim(idempotency_key)) BETWEEN 1 AND 128),
    status text NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING','APPROVED','REJECTED','CANCELLED','EXPIRED','APPLIED','FAILED','STALE')),
    current_step integer NOT NULL DEFAULT 1 CHECK (current_step BETWEEN 1 AND 20),
    step_count integer NOT NULL CHECK (step_count BETWEEN 1 AND 20),
    version bigint NOT NULL DEFAULT 1 CHECK (version >= 1),
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CHECK (current_step <= step_count),
    CHECK (expires_at > created_at),
    UNIQUE (tenant_id, request_id),
    UNIQUE (tenant_id, request_id, requester_id),
    UNIQUE (tenant_id, requester_id, idempotency_key),
    FOREIGN KEY (tenant_id, policy_id) REFERENCES identity.approval_policies(tenant_id, policy_id) ON DELETE RESTRICT,
    FOREIGN KEY (tenant_id, requester_id) REFERENCES identity.accounts(tenant_id, account_id) ON DELETE RESTRICT
);
CREATE INDEX idx_identity_approval_requests_inbox ON identity.approval_requests(tenant_id, status, expires_at, request_id);

-- Snapshot the resolved policy criteria on request creation. Later policy edits
-- must not reinterpret a pending request.
CREATE TABLE identity.approval_request_steps (
    tenant_id text NOT NULL,
    request_id uuid NOT NULL,
    step_number integer NOT NULL CHECK (step_number BETWEEN 1 AND 20),
    required_permission_code text NOT NULL REFERENCES identity.permissions(code) ON DELETE RESTRICT,
    required_role_id uuid,
    PRIMARY KEY (tenant_id, request_id, step_number),
    FOREIGN KEY (tenant_id, request_id) REFERENCES identity.approval_requests(tenant_id, request_id) ON DELETE RESTRICT,
    FOREIGN KEY (tenant_id, required_role_id) REFERENCES identity.roles(tenant_id, role_id) ON DELETE RESTRICT
);

CREATE TABLE identity.approval_decisions (
    tenant_id text NOT NULL,
    request_id uuid NOT NULL,
    requester_id uuid NOT NULL,
    step_number integer NOT NULL,
    reviewer_id uuid NOT NULL,
    decision text NOT NULL CHECK (decision IN ('APPROVE','REJECT')),
    comment text NOT NULL DEFAULT '' CHECK (length(comment) <= 2000),
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, request_id, step_number),
    CHECK (requester_id <> reviewer_id),
    CHECK (decision <> 'REJECT' OR length(btrim(comment)) > 0),
    FOREIGN KEY (tenant_id, request_id, requester_id) REFERENCES identity.approval_requests(tenant_id, request_id, requester_id) ON DELETE RESTRICT,
    FOREIGN KEY (tenant_id, request_id, step_number) REFERENCES identity.approval_request_steps(tenant_id, request_id, step_number) ON DELETE RESTRICT,
    FOREIGN KEY (tenant_id, reviewer_id) REFERENCES identity.accounts(tenant_id, account_id) ON DELETE RESTRICT
);

-- +goose StatementBegin
CREATE FUNCTION identity.reject_approval_decision_mutation() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    RAISE EXCEPTION 'approval decisions are immutable';
END;
$$;
-- +goose StatementEnd
CREATE TRIGGER approval_decisions_immutable BEFORE UPDATE OR DELETE ON identity.approval_decisions
FOR EACH ROW EXECUTE FUNCTION identity.reject_approval_decision_mutation();

-- +goose Down
DROP TABLE identity.approval_decisions;
DROP FUNCTION identity.reject_approval_decision_mutation();
DROP TABLE identity.approval_request_steps;
DROP TABLE identity.approval_requests;
DROP TABLE identity.approval_policy_steps;
DROP TABLE identity.approval_policies;

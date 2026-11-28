-- +goose Up
-- Operational history, separate from user-facing RBAC approval/resource audit.
-- No recipient address, proof, token hash, encrypted content or provider error is copied.
CREATE TABLE identity.invitation_delivery_audit (
    audit_id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id text,
    account_id uuid NOT NULL,
    invitation_id uuid NOT NULL,
    purpose text NOT NULL CHECK (purpose IN ('ACCOUNT_SETUP', 'PASSWORD_RESET', 'UNKNOWN')),
    event_type text NOT NULL CHECK (event_type IN ('QUEUED', 'CLAIMED', 'RETRY_SCHEDULED', 'DELIVERED', 'SUPPRESSED')),
    attempt integer NOT NULL CHECK (attempt >= 0),
    reason text,
    next_attempt_at timestamptz,
    occurred_at timestamptz NOT NULL DEFAULT clock_timestamp(),
    UNIQUE (invitation_id, attempt, event_type),
    CHECK (
        (event_type = 'SUPPRESSED' AND reason IS NOT NULL AND reason IN ('INVALID_PROOF', 'ACCOUNT_INELIGIBLE') AND next_attempt_at IS NULL)
        OR (event_type = 'RETRY_SCHEDULED' AND reason IS NOT NULL AND reason = 'DELIVERY_FAILED' AND next_attempt_at IS NOT NULL)
        OR (event_type IN ('QUEUED', 'CLAIMED', 'DELIVERED') AND reason IS NULL AND next_attempt_at IS NULL)
    )
);
CREATE INDEX idx_invitation_delivery_audit_tenant ON identity.invitation_delivery_audit(tenant_id, occurred_at DESC, audit_id DESC);
CREATE INDEX idx_invitation_delivery_audit_account ON identity.invitation_delivery_audit(tenant_id, account_id, audit_id);
CREATE TRIGGER invitation_delivery_audit_immutable BEFORE UPDATE OR DELETE ON identity.invitation_delivery_audit
FOR EACH ROW EXECUTE FUNCTION identity.reject_resource_audit_mutation();
CREATE TRIGGER invitation_delivery_audit_no_truncate BEFORE TRUNCATE ON identity.invitation_delivery_audit
FOR EACH STATEMENT EXECUTE FUNCTION identity.reject_resource_audit_mutation();

-- +goose StatementBegin
CREATE FUNCTION identity.record_invitation_delivery_audit() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    audit_event text;
    audit_reason text;
    retry_at timestamptz;
    account_tenant text;
BEGIN
    IF TG_OP = 'INSERT' THEN
        audit_event := 'QUEUED';
    ELSIF NEW.suppressed_at IS NOT NULL AND OLD.suppressed_at IS NULL THEN
        audit_event := 'SUPPRESSED';
        audit_reason := NEW.suppression_reason;
    ELSIF NEW.sent_at IS NOT NULL AND OLD.sent_at IS NULL THEN
        audit_event := 'DELIVERED';
    ELSIF NEW.attempts > OLD.attempts AND NEW.lease_until IS NOT NULL THEN
        audit_event := 'CLAIMED';
    ELSIF OLD.lease_until IS NOT NULL AND NEW.lease_until IS NULL
        AND NEW.sent_at IS NULL AND NEW.suppressed_at IS NULL THEN
        audit_event := 'RETRY_SCHEDULED';
        audit_reason := 'DELIVERY_FAILED';
        retry_at := NEW.next_attempt_at;
    ELSE
        RETURN NEW;
    END IF;

    SELECT NULLIF(tenant_id, '') INTO account_tenant
    FROM identity.accounts WHERE account_id = NEW.account_id;

    INSERT INTO identity.invitation_delivery_audit
        (tenant_id, account_id, invitation_id, purpose, event_type, attempt, reason, next_attempt_at)
    VALUES
        (account_tenant, NEW.account_id, NEW.id,
         CASE WHEN NEW.purpose IN ('ACCOUNT_SETUP', 'PASSWORD_RESET') THEN NEW.purpose ELSE 'UNKNOWN' END,
         audit_event, NEW.attempts, audit_reason, retry_at);
    RETURN NEW;
END;
$$;
-- +goose StatementEnd
CREATE TRIGGER invitation_delivery_audit_record AFTER INSERT OR UPDATE ON identity.invitations
FOR EACH ROW EXECUTE FUNCTION identity.record_invitation_delivery_audit();

-- +goose Down
DROP TRIGGER invitation_delivery_audit_record ON identity.invitations;
DROP FUNCTION identity.record_invitation_delivery_audit();
DROP TABLE identity.invitation_delivery_audit;

-- +goose Up
ALTER TABLE identity.invitations DROP CONSTRAINT invitation_suppression_state;
ALTER TABLE identity.invitations ADD CONSTRAINT invitation_suppression_state CHECK (
    (suppressed_at IS NULL AND suppression_reason IS NULL)
    OR (suppressed_at IS NOT NULL AND suppression_reason IS NOT NULL
        AND suppression_reason IN ('INVALID_PROOF', 'ACCOUNT_INELIGIBLE', 'DELIVERY_UNCERTAIN') AND sent_at IS NULL)
);
-- The previous migration used an unnamed table-level transition check.
ALTER TABLE identity.invitation_delivery_audit DROP CONSTRAINT invitation_delivery_audit_check;
ALTER TABLE identity.invitation_delivery_audit ADD CONSTRAINT invitation_delivery_audit_transition CHECK (
    (event_type = 'SUPPRESSED' AND reason IS NOT NULL AND reason IN ('INVALID_PROOF', 'ACCOUNT_INELIGIBLE', 'DELIVERY_UNCERTAIN') AND next_attempt_at IS NULL)
    OR (event_type = 'RETRY_SCHEDULED' AND reason IS NOT NULL AND reason = 'DELIVERY_FAILED' AND next_attempt_at IS NOT NULL)
    OR (event_type IN ('QUEUED', 'CLAIMED', 'DELIVERED') AND reason IS NULL AND next_attempt_at IS NULL)
);
-- +goose Down
-- The original constraints intentionally reject rollback while uncertain rows exist.
ALTER TABLE identity.invitation_delivery_audit DROP CONSTRAINT invitation_delivery_audit_transition;
ALTER TABLE identity.invitation_delivery_audit ADD CONSTRAINT invitation_delivery_audit_check CHECK (
    (event_type = 'SUPPRESSED' AND reason IS NOT NULL AND reason IN ('INVALID_PROOF', 'ACCOUNT_INELIGIBLE') AND next_attempt_at IS NULL)
    OR (event_type = 'RETRY_SCHEDULED' AND reason IS NOT NULL AND reason = 'DELIVERY_FAILED' AND next_attempt_at IS NOT NULL)
    OR (event_type IN ('QUEUED', 'CLAIMED', 'DELIVERED') AND reason IS NULL AND next_attempt_at IS NULL)
);
ALTER TABLE identity.invitations DROP CONSTRAINT invitation_suppression_state;
ALTER TABLE identity.invitations ADD CONSTRAINT invitation_suppression_state CHECK (
    (suppressed_at IS NULL AND suppression_reason IS NULL)
    OR (suppressed_at IS NOT NULL AND suppression_reason IS NOT NULL
        AND suppression_reason IN ('INVALID_PROOF', 'ACCOUNT_INELIGIBLE') AND sent_at IS NULL)
);

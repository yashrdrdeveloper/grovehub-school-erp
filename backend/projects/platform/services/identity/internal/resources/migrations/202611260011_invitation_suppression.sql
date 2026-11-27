-- +goose Up
ALTER TABLE identity.invitations
    ADD COLUMN suppressed_at timestamptz,
    ADD COLUMN suppression_reason text,
    ADD CONSTRAINT invitation_suppression_state CHECK (
        (suppressed_at IS NULL AND suppression_reason IS NULL)
        OR (suppressed_at IS NOT NULL AND suppression_reason IS NOT NULL
            AND suppression_reason IN ('INVALID_PROOF', 'ACCOUNT_INELIGIBLE') AND sent_at IS NULL)
    );

-- +goose Down
-- Retain suppressed jobs as terminal before removing the new state.
UPDATE identity.invitations SET sent_at = suppressed_at WHERE suppressed_at IS NOT NULL;
ALTER TABLE identity.invitations DROP CONSTRAINT invitation_suppression_state,
    DROP COLUMN suppression_reason, DROP COLUMN suppressed_at;

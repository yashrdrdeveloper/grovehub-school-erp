-- +goose Up
DROP TRIGGER invitation_delivery_audit_record ON identity.invitations;
DROP FUNCTION identity.record_invitation_delivery_audit();
DROP TABLE identity.invitation_delivery_audit;
-- Registration receipts still use this immutability guard. Preserve it.

-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
    RAISE EXCEPTION 'retired invitation audit journal cannot be restored';
END;
$$;
-- +goose StatementEnd

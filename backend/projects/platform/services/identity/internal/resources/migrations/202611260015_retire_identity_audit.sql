-- +goose Up
-- Pre-release replacement: Audit owns searchable activity history.
-- Approval decisions, registration receipts and policy/migration receipts are
-- business state and must not be dropped.
DROP TABLE identity.approval_audit;
DROP TABLE identity.resource_audit;
DROP TABLE identity.role_metadata_audit;
DROP TABLE identity.account_metadata_audit;
DROP FUNCTION identity.reject_approval_audit_mutation();

-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
    RAISE EXCEPTION 'retired Identity audit journals cannot be restored; restore an approved database backup instead';
END;
$$;
-- +goose StatementEnd

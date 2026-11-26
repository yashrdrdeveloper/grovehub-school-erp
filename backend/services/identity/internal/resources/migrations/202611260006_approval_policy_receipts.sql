-- +goose Up
CREATE TABLE identity.approval_policy_receipts (
    receipt_id uuid PRIMARY KEY,
    tenant_id text NOT NULL CHECK (length(btrim(tenant_id)) > 0),
    operator_id uuid NOT NULL,
    author_id uuid NOT NULL,
    input_sha256 text NOT NULL CHECK (input_sha256 ~ '^[0-9a-f]{64}$'),
    result jsonb NOT NULL CHECK (jsonb_typeof(result) = 'object'),
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, input_sha256),
    FOREIGN KEY (tenant_id, author_id) REFERENCES identity.accounts(tenant_id, account_id) ON DELETE RESTRICT
);
-- +goose StatementBegin
CREATE FUNCTION identity.reject_approval_policy_receipt_mutation() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    RAISE EXCEPTION 'approval policy receipts are immutable';
END;
$$;
-- +goose StatementEnd
CREATE TRIGGER approval_policy_receipts_immutable BEFORE UPDATE OR DELETE ON identity.approval_policy_receipts
FOR EACH ROW EXECUTE FUNCTION identity.reject_approval_policy_receipt_mutation();

-- +goose Down
DROP TABLE identity.approval_policy_receipts;
DROP FUNCTION identity.reject_approval_policy_receipt_mutation();

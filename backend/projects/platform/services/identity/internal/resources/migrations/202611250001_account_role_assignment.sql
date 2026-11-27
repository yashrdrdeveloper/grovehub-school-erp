-- +goose Up
ALTER TABLE identity.accounts
    ADD COLUMN primary_role_id text,
    ADD COLUMN secondary_role_ids jsonb NOT NULL DEFAULT '[]'::jsonb;

ALTER TABLE identity.accounts
    ADD CONSTRAINT ck_identity_account_roles_array
    CHECK (jsonb_typeof(secondary_role_ids) = 'array'),
    ADD CONSTRAINT ck_identity_account_roles_primary
    CHECK (primary_role_id IS NOT NULL OR secondary_role_ids = '[]'::jsonb),
    ADD CONSTRAINT ck_identity_account_roles_distinct
    CHECK (primary_role_id IS NULL OR NOT (secondary_role_ids ? primary_role_id));

CREATE INDEX idx_identity_account_tenant_primary_role
    ON identity.accounts (tenant_id, primary_role_id);

-- +goose Down
DROP INDEX IF EXISTS identity.idx_identity_account_tenant_primary_role;
ALTER TABLE identity.accounts
    DROP CONSTRAINT IF EXISTS ck_identity_account_roles_distinct,
    DROP CONSTRAINT IF EXISTS ck_identity_account_roles_primary,
    DROP CONSTRAINT IF EXISTS ck_identity_account_roles_array,
    DROP COLUMN IF EXISTS secondary_role_ids,
    DROP COLUMN IF EXISTS primary_role_id;

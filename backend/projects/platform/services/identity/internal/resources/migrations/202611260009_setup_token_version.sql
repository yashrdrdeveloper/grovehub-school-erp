-- +goose Up
ALTER TABLE identity.password_tokens ADD COLUMN account_version bigint;
ALTER TABLE identity.password_tokens ADD CONSTRAINT ck_setup_account_version
CHECK (account_version IS NULL OR (account_version >= 1 AND purpose IN ('ACCOUNT_SETUP','SETUP_GRANT')));

-- +goose Down
ALTER TABLE identity.password_tokens DROP CONSTRAINT ck_setup_account_version;
ALTER TABLE identity.password_tokens DROP COLUMN account_version;

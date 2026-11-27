-- +goose Up
ALTER TABLE identity.password_tokens DROP CONSTRAINT ck_setup_account_version;
ALTER TABLE identity.password_tokens ADD CONSTRAINT ck_token_account_version
CHECK (account_version IS NULL OR (account_version >= 1 AND purpose IN ('ACCOUNT_SETUP','SETUP_GRANT','PASSWORD_RESET','RESET_GRANT')));

-- +goose Down
-- Refuse to discard recovery version bindings on databases with recovery history.
ALTER TABLE identity.password_tokens DROP CONSTRAINT ck_token_account_version;
ALTER TABLE identity.password_tokens ADD CONSTRAINT ck_setup_account_version
CHECK (account_version IS NULL OR (account_version >= 1 AND purpose IN ('ACCOUNT_SETUP','SETUP_GRANT')));

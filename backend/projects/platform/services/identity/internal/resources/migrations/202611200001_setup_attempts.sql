-- +goose Up
-- Setup attempt budget is persisted on the short-lived grant token, not in browser state.
ALTER TABLE identity.password_tokens ADD COLUMN attempts integer NOT NULL DEFAULT 0;
ALTER TABLE identity.password_tokens ADD CONSTRAINT ck_identity_password_token_attempts CHECK (attempts >= 0);

-- +goose Down
ALTER TABLE identity.password_tokens DROP CONSTRAINT IF EXISTS ck_identity_password_token_attempts;
ALTER TABLE identity.password_tokens DROP COLUMN IF EXISTS attempts;

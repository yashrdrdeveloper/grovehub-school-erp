-- +goose Up
ALTER TABLE identity.accounts ADD COLUMN IF NOT EXISTS tenant_id text;
ALTER TABLE identity.accounts ADD COLUMN IF NOT EXISTS username text;
ALTER TABLE identity.accounts ADD COLUMN IF NOT EXISTS search_details text NOT NULL DEFAULT '';
CREATE UNIQUE INDEX IF NOT EXISTS uq_identity_account_tenant_username ON identity.accounts (tenant_id, lower(username)) WHERE username IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_identity_account_tenant_status ON identity.accounts (tenant_id, status, created_at, account_id);
CREATE INDEX IF NOT EXISTS idx_identity_account_person_gid ON identity.accounts (person_gid);
CREATE TABLE IF NOT EXISTS identity.password_tokens (
 token_hash char(64) PRIMARY KEY,
 account_id uuid NOT NULL REFERENCES identity.accounts(account_id) ON DELETE CASCADE,
 purpose text NOT NULL CHECK (purpose IN ('ACCOUNT_SETUP','PASSWORD_RESET','SETUP_GRANT','RESET_GRANT')),
 expires_at timestamptz NOT NULL,
 used_at timestamptz,
 created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_identity_password_tokens_account ON identity.password_tokens (account_id, purpose);
CREATE TABLE IF NOT EXISTS identity.invitation_mail (
 id uuid PRIMARY KEY,
 account_id uuid NOT NULL REFERENCES identity.accounts(account_id) ON DELETE CASCADE,
 email text NOT NULL,
 purpose text NOT NULL,
 ciphertext bytea NOT NULL,
 nonce bytea NOT NULL,
 attempts integer NOT NULL DEFAULT 0,
 next_attempt_at timestamptz NOT NULL,
 lease_until timestamptz,
 sent_at timestamptz,
 created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_identity_invitation_mail_pending ON identity.invitation_mail (next_attempt_at) WHERE sent_at IS NULL;
ALTER TABLE identity.sessions ADD COLUMN IF NOT EXISTS reauthenticated_at timestamptz;
ALTER TABLE identity.sessions ADD COLUMN IF NOT EXISTS user_agent text;
ALTER TABLE identity.sessions ADD COLUMN IF NOT EXISTS last_active_at timestamptz;
-- +goose Down
DROP TABLE IF EXISTS identity.invitation_mail;
DROP TABLE IF EXISTS identity.password_tokens;
ALTER TABLE identity.sessions DROP COLUMN IF EXISTS last_active_at;
ALTER TABLE identity.sessions DROP COLUMN IF EXISTS user_agent;
ALTER TABLE identity.sessions DROP COLUMN IF EXISTS reauthenticated_at;
DROP INDEX IF EXISTS identity.uq_identity_account_tenant_username;
DROP INDEX IF EXISTS identity.idx_identity_account_tenant_status;
DROP INDEX IF EXISTS identity.idx_identity_account_person_gid;
ALTER TABLE identity.accounts DROP COLUMN IF EXISTS search_details;
ALTER TABLE identity.accounts DROP COLUMN IF EXISTS username;
ALTER TABLE identity.accounts DROP COLUMN IF EXISTS tenant_id;

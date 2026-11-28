CREATE SCHEMA IF NOT EXISTS internal_tool;
CREATE TABLE internal_tool.operators (
 id uuid PRIMARY KEY,
 email text NOT NULL UNIQUE CHECK (email=lower(email) AND length(email)<=254),
 password_salt bytea NOT NULL CHECK (octet_length(password_salt)=16),
 password_hash bytea NOT NULL CHECK (octet_length(password_hash)=32),
 persona text NOT NULL CHECK (persona IN ('L1','L2')),
 status text NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE','DISABLED')),
 created_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE internal_tool.sessions (
 token_hash bytea PRIMARY KEY CHECK (octet_length(token_hash)=32),
 operator_id uuid NOT NULL REFERENCES internal_tool.operators(id),
 csrf_hash bytea NOT NULL CHECK (octet_length(csrf_hash)=32),
 created_at timestamptz NOT NULL DEFAULT now(),
 expires_at timestamptz NOT NULL,
 revoked_at timestamptz,
 CHECK (expires_at>created_at)
);
CREATE INDEX internal_tool_sessions_operator ON internal_tool.sessions(operator_id);
CREATE INDEX internal_tool_sessions_expiry ON internal_tool.sessions(expires_at);
CREATE TABLE internal_tool.login_attempts (
 email_hash bytea PRIMARY KEY CHECK (octet_length(email_hash)=32),
 window_start timestamptz NOT NULL,
 attempts integer NOT NULL CHECK (attempts>0)
);

-- +goose Up
-- Additive RBAC persistence. No existing role assignment is rewritten and no
-- permissions or entitlements are granted by this migration.
ALTER TABLE identity.accounts
    ADD CONSTRAINT uq_identity_accounts_tenant_account UNIQUE (tenant_id, account_id);

CREATE TABLE identity.features (
    feature_id uuid PRIMARY KEY,
    code text NOT NULL UNIQUE CHECK (length(btrim(code)) > 0),
    module text NOT NULL CHECK (length(btrim(module)) > 0),
    name text NOT NULL CHECK (length(btrim(name)) > 0),
    description text NOT NULL DEFAULT '',
    active boolean NOT NULL DEFAULT false,
    version bigint NOT NULL DEFAULT 1 CHECK (version >= 1),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE identity.permissions (
    permission_id uuid PRIMARY KEY,
    code text NOT NULL UNIQUE CHECK (length(btrim(code)) > 0),
    module text NOT NULL CHECK (length(btrim(module)) > 0),
    name text NOT NULL CHECK (length(btrim(name)) > 0),
    description text NOT NULL DEFAULT '',
    feature_id uuid REFERENCES identity.features(feature_id) ON DELETE RESTRICT,
    active boolean NOT NULL DEFAULT false,
    version bigint NOT NULL DEFAULT 1 CHECK (version >= 1),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_identity_permissions_feature ON identity.permissions(feature_id);

CREATE TABLE identity.tenant_features (
    tenant_id text NOT NULL CHECK (length(btrim(tenant_id)) > 0),
    feature_id uuid NOT NULL REFERENCES identity.features(feature_id) ON DELETE RESTRICT,
    enabled boolean NOT NULL DEFAULT false,
    version bigint NOT NULL DEFAULT 1 CHECK (version >= 1),
    updated_by uuid,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, feature_id)
);

CREATE TABLE identity.roles (
    role_id uuid PRIMARY KEY,
    tenant_id text NOT NULL CHECK (length(btrim(tenant_id)) > 0),
    name text NOT NULL CHECK (length(btrim(name)) > 0),
    normalized_name text NOT NULL CHECK (length(btrim(normalized_name)) > 0),
    description text NOT NULL DEFAULT '',
    status text NOT NULL DEFAULT 'INACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE')),
    search_details text NOT NULL DEFAULT '',
    version bigint NOT NULL DEFAULT 1 CHECK (version >= 1),
    created_by uuid,
    updated_by uuid,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, role_id),
    UNIQUE (tenant_id, normalized_name)
);

CREATE INDEX idx_identity_roles_tenant_status ON identity.roles(tenant_id, status, created_at, role_id);

CREATE TABLE identity.role_personas (
    tenant_id text NOT NULL,
    role_id uuid NOT NULL,
    persona_code text NOT NULL CHECK (persona_code IN (
        'SCHOOL_ADMINISTRATION', 'TEACHER', 'COORDINATOR',
        'ACCOUNTANT', 'LIBRARIAN', 'RECEPTIONIST'
    )),
    PRIMARY KEY (tenant_id, role_id, persona_code),
    FOREIGN KEY (tenant_id, role_id) REFERENCES identity.roles(tenant_id, role_id) ON DELETE RESTRICT
);

CREATE TABLE identity.role_permissions (
    tenant_id text NOT NULL,
    role_id uuid NOT NULL,
    permission_id uuid NOT NULL REFERENCES identity.permissions(permission_id) ON DELETE RESTRICT,
    granted_by uuid,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, role_id, permission_id),
    FOREIGN KEY (tenant_id, role_id) REFERENCES identity.roles(tenant_id, role_id) ON DELETE RESTRICT
);

CREATE TABLE identity.account_roles (
    tenant_id text NOT NULL CHECK (length(btrim(tenant_id)) > 0),
    account_id uuid NOT NULL,
    role_id uuid NOT NULL,
    assignment_type text NOT NULL CHECK (assignment_type IN ('PRIMARY', 'SECONDARY')),
    granted_by uuid,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, account_id, role_id),
    FOREIGN KEY (tenant_id, account_id) REFERENCES identity.accounts(tenant_id, account_id) ON DELETE RESTRICT,
    FOREIGN KEY (tenant_id, role_id) REFERENCES identity.roles(tenant_id, role_id) ON DELETE RESTRICT
);

CREATE UNIQUE INDEX uq_identity_account_roles_primary
    ON identity.account_roles(tenant_id, account_id) WHERE assignment_type = 'PRIMARY';

CREATE INDEX idx_identity_account_roles_role ON identity.account_roles(tenant_id, role_id);

-- +goose Down
DROP TABLE identity.account_roles;
DROP TABLE identity.role_permissions;
DROP TABLE identity.role_personas;
DROP TABLE identity.roles;
DROP TABLE identity.tenant_features;
DROP TABLE identity.permissions;
DROP TABLE identity.features;
ALTER TABLE identity.accounts DROP CONSTRAINT uq_identity_accounts_tenant_account;

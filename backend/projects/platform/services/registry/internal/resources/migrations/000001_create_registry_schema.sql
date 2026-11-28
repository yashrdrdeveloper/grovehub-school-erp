-- +goose Up
CREATE SCHEMA registry;

CREATE TABLE registry.projects (
    project_id uuid PRIMARY KEY,
    code text NOT NULL UNIQUE CHECK (code ~ '^[a-z][a-z0-9_-]{1,63}$'),
    name text NOT NULL CHECK (length(btrim(name)) BETWEEN 1 AND 128),
    status text NOT NULL DEFAULT 'DISABLED' CHECK (status IN ('ACTIVE','DISABLED')),
    version bigint NOT NULL DEFAULT 1 CHECK (version > 0),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE registry.tenants (
    tenant_id uuid PRIMARY KEY,
    project_id uuid NOT NULL REFERENCES registry.projects(project_id),
    external_reference text NOT NULL CHECK (length(btrim(external_reference)) BETWEEN 1 AND 128),
    status text NOT NULL DEFAULT 'DISABLED' CHECK (status IN ('ACTIVE','DISABLED')),
    session_generation bigint NOT NULL DEFAULT 1 CHECK (session_generation > 0),
    version bigint NOT NULL DEFAULT 1 CHECK (version > 0),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (project_id, tenant_id),
    UNIQUE (project_id, external_reference)
);

CREATE TABLE registry.workloads (
    workload_id uuid PRIMARY KEY,
    project_id uuid REFERENCES registry.projects(project_id),
    subject text NOT NULL CHECK (length(btrim(subject)) BETWEEN 1 AND 128),
    issuer text NOT NULL CHECK (length(btrim(issuer)) BETWEEN 1 AND 128),
    status text NOT NULL DEFAULT 'DISABLED' CHECK (status IN ('ACTIVE','DISABLED')),
    version bigint NOT NULL DEFAULT 1 CHECK (version > 0),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (issuer, subject)
);

CREATE TABLE registry.workload_keys (
    key_id text PRIMARY KEY CHECK (length(btrim(key_id)) BETWEEN 1 AND 128),
    workload_id uuid NOT NULL REFERENCES registry.workloads(workload_id),
    public_key bytea NOT NULL UNIQUE CHECK (octet_length(public_key) = 32),
    valid_from timestamptz NOT NULL,
    valid_until timestamptz NOT NULL CHECK (valid_until > valid_from),
    revoked_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX workload_keys_owner ON registry.workload_keys(workload_id);

-- These capabilities are provider-owned, never defined by untrusted request data.
CREATE TABLE registry.capabilities (
    code text NOT NULL CHECK (length(btrim(code)) BETWEEN 1 AND 128),
    scope text NOT NULL CHECK (scope IN ('PLATFORM','PROJECT','TENANT')),
    PRIMARY KEY (code, scope)
);

CREATE TABLE registry.workload_grants (
    grant_id uuid PRIMARY KEY,
    workload_id uuid NOT NULL REFERENCES registry.workloads(workload_id),
    target_project_id uuid REFERENCES registry.projects(project_id),
    target_tenant_id uuid,
    scope text NOT NULL CHECK (scope IN ('PLATFORM','PROJECT','TENANT')),
    capability text NOT NULL,
    allow_delegation boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (capability, scope) REFERENCES registry.capabilities(code, scope),
    FOREIGN KEY (target_project_id, target_tenant_id) REFERENCES registry.tenants(project_id, tenant_id),
    CHECK (
        (scope = 'PLATFORM' AND target_project_id IS NULL AND target_tenant_id IS NULL)
        OR (scope = 'PROJECT' AND target_project_id IS NOT NULL AND target_tenant_id IS NULL)
        OR (scope = 'TENANT' AND target_project_id IS NOT NULL AND target_tenant_id IS NOT NULL)
    ),
    UNIQUE NULLS NOT DISTINCT (workload_id, target_project_id, target_tenant_id, scope, capability)
);

-- +goose Down
DROP TABLE registry.workload_grants;
DROP TABLE registry.capabilities;
DROP TABLE registry.workload_keys;
DROP TABLE registry.workloads;
DROP TABLE registry.tenants;
DROP TABLE registry.projects;
DROP SCHEMA registry;

# School backend

Owns browser-facing School and Identity use-case APIs. Platform calls use the restricted authenticated generated SDK clients. No direct Platform database access.

Startup requires SCHOOL_DATABASE_URL for a separately provisioned PostgreSQL database and least-privilege role. No fallback to DATABASE_URL or other project credentials. GET /health/database checks only the owned database; /health/ready checks the owned database and Platform /ready after successful local configuration validation; workload grant and browser verification are separate. School registry migration version 2 is required before startup.

## School registry

The School project owns the authoritative UUID, code, name and lifecycle status. Run `make migrate` before startup, then `go run -mod=mod ./cmd/school register-school UUID CODE "School name"` against the owned database. Registration is idempotent for the same UUID/code/name; conflicting records are rejected. Use the existing canonical tenant UUID when provisioning; there is no automatic Identity tenant creation or translation. Only the deployment operator can provision records in this phase. No rename or disable command is exposed yet.

The School-owned OpenAPI contract is in services/registry/internal/resources/api/openapi/v1/openapi.yaml. POST /school/restricted/registry/lookup is a private, signed workload API. Configure SCHOOL_REGISTRY_READ_GRANTS as a JSON array of objects containing key (WorkloadKeyConfig with public key only) and schoolIds (allowed UUID array). Each key has issuer internal-tool-school-directory, subject internal-tool-school-directory:<reader UUID>, one scope internal.school.registry.read and no delegation; audience school. Every requested school must be authorized before the query runs. Missing rows are omitted. No public enumeration, cross-project database access, browser credential forwarding or mutation API is available.

## School name mutation

Migration 2 adds immutable school_registry.rename_events. Configure SCHOOL_REGISTRY_RENAME_KEYS as a JSON array of independent WorkloadKeyConfig public keys with issuer internal-tool-school-rename, subject internal-tool-school-rename:<school UUID>, only internal.school.registry.rename and no delegation. PATCH /school/restricted/registry/schools/{schoolId}/name accepts expectedName and name plus the execution UUID in Idempotency-Key. New writes require ACTIVE status and the exact current name; identical committed retries return 204 without modifying the current name. Conflicting retries return 409. The local audit receipt is written in the same transaction. This is name-only mutation; no school status, deletion, tenant creation or code changes are enabled. Local receipts do not yet replicate to centralized Audit.

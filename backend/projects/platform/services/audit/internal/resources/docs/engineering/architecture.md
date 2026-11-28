# Audit architecture

Audit owns event ingestion, activity persistence, storage diagnostics and replay. Its bootstrap registers `audit-storage` and `audit-replay` and the Audit database migration. Storage inspection lives in `internal/app/services/storage`; it reads PostgreSQL statistics in a read-only transaction and requires the exact database name. Estimates are not exact row counts. Replay reads retained published events from Infrastructure; preserve tenant/event identity checks and dry-run/apply safeguards.

Migrations are embedded from `internal/resources/migrations`. Run Platform command and Audit unit/integration tests before changes; integration tests require disposable PostgreSQL.

## Test resource layout

Service-owned test directories live under `internal/resources/tests/`, including unit and integration suites, helpers and mocks. Generated suites mirror the resource paths under `.generated/`.

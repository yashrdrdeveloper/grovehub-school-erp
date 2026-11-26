# Offline Identity RBAC migration runbook

Status: staged operator tooling; NOT authorization cutover. Do not use against school data until integration CI is green, a verified backup exists and cutover is reviewed.

## Preconditions

- Stop every application instance, worker and other database writer for the maintenance window. The confirmation flag is an operator attestation, not an automatic process check.
- Verify target connection, database name, exact tenant, operator UUID and restore procedure. The operator UUID is an audit attribution supplied by the trusted database operator, not application authentication.
- Apply the additive Identity schema migrations including the migration receipts table using the established migration procedure; the CLI never migrates schema automatically.
- Store plans/manifests securely; they contain account IDs. Use absolute paths for input files. Do not publish them or commit them to source control.

## Read-only plan

```bash
umask 077
bash scripts/identity/rbac-migrate.sh plan --tenant-id "$TENANT_ID" --database "$DB_NAME" > "$PLAN_FILE"
shasum -a 256 "$PLAN_FILE"
```

Review all assignments and unassigned accounts; retain the reviewed digest. The output contains JSON only on success (connection logging, if configured, should go to stderr; verify before review).

## Backfill (offline only)

```bash
bash scripts/identity/rbac-migrate.sh backfill --tenant-id "$TENANT_ID" --database "$DB_NAME" \
  --operator-id "$OPERATOR_ID" --input "$PLAN_FILE" --sha256 "$REVIEWED_PLAN_SHA" \
  --confirm-offline-maintenance
```

Creates inactive roles and inert normalized assignments only. The transaction includes a unique migration receipt. Reapplication fails rather than merging. A failed operation leaves no receipt or partial writes.

## Bootstrap (offline only)

Create a manifest JSON with the exact original plan under `plan`, the nominated active account UUID under `administratorId`, and an explicit `grants` array of objects with `roleId` and `permissionCodes`. Review all affected accounts. Recovery grants require exclusive role ownership. Register/activate capabilities and enable the relevant tenant features through a separately reviewed operator process; this command does not enable them. Hash the reviewed manifest file separately.

```bash
bash scripts/identity/rbac-migrate.sh bootstrap --tenant-id "$TENANT_ID" --database "$DB_NAME" \
  --operator-id "$OPERATOR_ID" --input "$MANIFEST_FILE" --sha256 "$REVIEWED_MANIFEST_SHA" \
  --confirm-offline-maintenance
```

Success returns a receipt UUID; query identity.rbac_migration_receipts by tenant and operation to reconcile uncertain client output. Do not retry blindly. Receipts record actor, reviewed input digest, result and timestamp, not raw manifests.

## Cutover and recovery

Legacy fields and runtime permission policy remain unchanged. Do not enable normal school access based on completion of these commands. Live endpoint permission enforcement, compatibility, approval, audit/outbox, and explicit migration readiness checks remain release blockers. Keep writers stopped if a failure requires restore; verify backup and reconciliation with the database operator. There is no destructive automatic rollback after a committed bootstrap.

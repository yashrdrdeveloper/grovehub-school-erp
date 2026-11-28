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

## Approval policy provisioning (offline only)

After backfill and administrator bootstrap receipts exist, review a manifest with exact tenantId, authorId (an active account in that tenant, used for policy attribution), and policies. Include exactly one revision for each of ACCOUNT_LIFECYCLE_CHANGE, ROLE_LIFECYCLE_CHANGE, ROLE_PERMISSION_CHANGE, ROLE_ASSIGNMENT_CHANGE. Each revision has requestType, expectedVersion (0 only for a new type, otherwise the highest stored revision), and ordered steps (1–20 objects with permissionCode and optional roleId). Choose school-delegable permissions and role constraints explicitly; there are no default grants or automatic approvals. Each step must currently have at least two distinct eligible active accounts. Assignment beneficiary exclusion may require a third; review actual workflow liveness separately. Policy author is not a substitute for the trusted offline operator's authorization.

```bash
bash scripts/identity/rbac-migrate.sh policies --tenant-id "$TENANT_ID" --database "$DB_NAME" \
  --operator-id "$OPERATOR_ID" --input "$POLICY_FILE" --sha256 "$REVIEWED_POLICY_SHA" \
  --confirm-offline-maintenance
```

The batch atomically deactivates previous policy revisions and inserts new ones with immutable provisioning receipt in identity.approval_policy_receipts. Old requests retain their snapshotted criteria. Stale revisions or unavailable reviewers abort the entire batch. Review the receipt before retrying an uncertain outcome. Do not use online; this command uses table locks and requires all other writers stopped. This does not activate execution routes or establish full cutover readiness.

## Authority lock protocol

Cooperating online tenant writers must acquire shared global authority lock before the tenant lock. Future global permission/feature/entitlement mutation entrypoints must use the fresh-transaction exclusive GlobalGuard and preserve every protected school's effective administrator. Offline migration and policy commands acquire exclusive global lock before table locks, but must still run with all writers stopped. Do not perform ad hoc SQL capability changes: database owners can bypass application safeguards. Composite execution readiness now checks live administrators, all four active policy revision IDs against the latest receipt, and live eligible reviewers. This is not a production cutover approval.

## Controlled execution opt-in

Leave IDENTITY_APPROVAL_EXECUTION_TENANT unset until the rollout is independently reviewed and CI is green. For a disposable staging database only, set it to one exact tenant identifier in strict mode, restart every runtime consistently, and verify migration/policy receipts, independent reviewers, audit/outbox delivery and session revocation. No wildcard or legacy combination is accepted. The requester must have authenticated within five minutes; the executor still revalidates the approved request, requester/reviewer authority and composite readiness inside its transaction. Unset the variable and restart all instances to disable execution; this does not undo committed effects. Legacy-development requires OTEL_DEPLOYMENT_ENVIRONMENT=development and HTTP loopback allowed origin, and refuses a database containing migration receipts; never use it behind a public tunnel or against real school data. No production cutover is authorized by this runbook alone.

## Invitation delivery reconciliation (read-only operator diagnostic)

Use a database credential authorized for read-only access to the two mail schemas, with both services' migrations applied. Confirm the exact database and tenant before running from the repository root:

```bash
umask 077
cd backend
go run ./cmd identity-mail-reconcile --tenant-id "$TENANT_ID" --database "$DB_NAME" --limit 100
```

If nextCursor is present, pass it as --after for the next page. Pages are live snapshots, not a frozen export; repeat review after concurrent work settles. Output contains account and delivery UUIDs, not email addresses or secrets. Keep it private.

- REVIEW_UNCERTAIN or STALE_DISPATCH: provider may have accepted the mail. Inspect provider-side records using authorized operator access. Do not replay the same ID, delete reservations, or edit audit history.
- ACCEPTED_WITHOUT_SENT: provider acceptance is recorded; it does not prove inbox delivery. A still-valid queued job may naturally reconcile through the existing no-send replay path. A suppressed job must remain suppressed; review manually.
- SENT_WITHOUT_ACCEPTANCE or RECEIPT_WITHOUT_JOB: investigate discrepancy; do not infer that no email was sent.
- If a pending setup account needs another link, use the existing authorized resend endpoint with recent authentication and current delegation. A new proof and delivery UUID are issued; the old proof is invalidated. For an active account, use the normal forgot-password workflow. Never manually restore the old proof. No blanket retry or force-delivery operation is provided.

Independent Integration reservations and receipts require an available pool connection while Identity holds its transaction. A pool capped at one connection cannot support this path. Two connections suffice for one isolated dispatch, not for simultaneous tenants or production load. Size and measure the shared pool against competing requests, connection wait duration and database capacity; review before enabling real school traffic.

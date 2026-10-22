# Database Guidelines

> **Status:** Mandatory Phase 0 engineering standard  
> **Location:** `docs/engineering/database/DATABASE_GUIDELINES.md`  
> **Applies to:** PostgreSQL schema design, migrations, repositories, transactions, indexes, constraints, JSONB, pagination, locking, data ownership, backups, test databases, and database evolution  
> **Primary database:** PostgreSQL  
> **Core principles:** explicit ownership, strong constraints, transactional correctness, application-generated IDs, safe migrations, service isolation, query observability, real integration tests, and extraction-friendly boundaries

---

## 1. Objective

PostgreSQL is Perihellion's primary transactional system of record.

The database must support:

```text
correct business invariants
predictable migrations
safe concurrent writes
fast indexed access
clear service ownership
future service extraction
backup/restore
observability
```

The database is not a shared-memory shortcut between services.

---

## 2. Initial Topology

Initially:

```text
one PostgreSQL deployment
        ↓
one Perihellion product database
        ↓
logical ownership by backend service/module
```

Sharing the same PostgreSQL instance does **not** allow cross-service table access.

---

## 3. Logical Service Ownership

Default recommendation:

```text
one PostgreSQL schema per backend service/module
```

Example:

```text
catalog.applications
catalog.application_revisions

execution.executions
execution.execution_steps

governance.approvals
```

This gives strong naming/ownership without requiring separate database servers on day one.

If a service intentionally uses another layout, ownership rules still apply.

---

## 4. Mandatory Rules

1. PostgreSQL is the default transactional database.
2. Each table has exactly one owning service/module.
3. A service never queries another service's tables.
4. Cross-service data access goes through public service clients/contracts.
5. Cross-service foreign keys are forbidden.
6. Foreign keys inside one owning service are encouraged where they protect invariants.
7. Every schema change has a migration.
8. Migrations are committed with application code.
9. Production schema is never changed manually.
10. IDs are generated before repository insertion unless database generation is a deliberate contract.
11. Repositories contain persistence logic only.
12. Business transaction boundaries are owned by the business service/use case.
13. Important invariants use database constraints where practical.
14. `NOT NULL` is used when absence is not a valid state.
15. Unique business invariants use unique constraints/indexes.
16. Check constraints are preferred for simple stable invariants.
17. PostgreSQL ENUM is avoided unless the value set is exceptionally stable and database-level enum behavior is specifically needed.
18. Application enums normally use text/small integer plus application validation and optional check constraints.
19. `timestamptz` is the default for instants.
20. `timestamp without time zone` is used only for true timezone-less local values.
21. JSONB is for flexible/provider metadata, not an excuse to avoid relational modeling.
22. Queries are parameterized.
23. Dynamic identifiers/sort fields use whitelists.
24. Unbounded table scans in request paths are forbidden.
25. Large lists paginate.
26. Every important query has an indexing plan.
27. Migration rollback/forward recovery is considered before merge.
28. Destructive migrations use expand/migrate/contract strategy.
29. Long table locks are avoided during live upgrades.
30. Database errors are mapped to semantic errors before leaving the repository/service boundary.
31. Real PostgreSQL is used for integration testing.
32. Backups are not considered valid until restore is tested.
33. Secrets/credentials are never stored in ordinary configuration tables when a secret reference is sufficient.
34. Audit/security data follows explicit retention and access policy.
35. Query observability is mandatory.

---

## 5. Naming

Use `snake_case`.

Recommended:

```text
schemas:              catalog
tables:               applications
columns:              application_id
primary keys:         pk_applications
foreign keys:         fk_application_revisions_application
unique constraints:   uq_applications_org_name
indexes:              idx_applications_org_created_at
check constraints:    chk_applications_status
```

Names must communicate intent.

---

## 6. Table Names

Use plural nouns for entity collections unless the repository already has a stronger consistent convention.

Examples:

```text
applications
executions
provider_bindings
```

Avoid abbreviations that only the original author understands.

---

## 7. Primary Keys

The default is a single stable application-generated identifier.

Recommended storage:

```text
uuid
```

when the application ID strategy emits UUID-compatible values.

If another opaque ID representation is selected platform-wide, use it consistently.

Repositories must not expose assumptions about the textual format.

---

## 8. ID Generation

Generate identifiers before persistence.

Conceptually:

```text
business/service layer
      ↓
ID generator abstraction
      ↓
entity with ID
      ↓
repository insert
```

This supports:

```text
events before commit preparation
testing
idempotency
non-database implementations
```

---

## 9. Natural Keys

Natural/business keys may have unique constraints.

They should normally not replace stable primary IDs when values can change.

Example:

```text
primary ID: application_id
unique business key: (organization_id, name)
```

---

## 10. Tenant / Organization Scoping

Multi-tenant/domain-owned tables that require organization/project isolation should carry the required scope explicitly.

Indexes and uniqueness must include the appropriate scope.

Example:

```text
UNIQUE (organization_id, name)
```

Do not rely only on application filtering for critical isolation invariants.

---

## 11. Cross-Service References

A table may store another service's opaque ID as data when business semantics require it.

But:

```text
no cross-service DB join
no cross-service FK
no direct lookup into foreign schema
```

Resolution belongs through the public service contract when current data is required.

---

## 12. Foreign Keys

Inside one service:

```text
use FK when it protects a real invariant
```

Choose deletion behavior deliberately:

```text
RESTRICT
CASCADE
SET NULL
```

Do not default to CASCADE without understanding domain lifecycle.

---

## 13. `NOT NULL`

If a field is required in all valid rows, use `NOT NULL`.

Do not encode required business fields as nullable just to make migrations easier.

Use staged migration when introducing a required field.

---

## 14. Defaults

Database defaults are appropriate for storage-level invariants such as:

```text
created_at
simple boolean flags
```

Business defaults that may vary by policy belong in application logic/configuration.

---

## 15. Timestamps

Recommended common columns where applicable:

```text
created_at timestamptz NOT NULL
updated_at timestamptz NOT NULL
```

Use application/DB conventions consistently.

Do not rely on local server timezone.

---

## 16. Soft Delete

Soft delete is not the default.

Use it only when domain semantics require recoverability/history/retention.

If used, define:

```text
deleted_at
query behavior
unique constraint behavior
restore semantics
retention/purge
```

---

## 17. Status Fields

Stable workflow status is usually stored as text or small integer with a constraint where useful.

Avoid PostgreSQL ENUM for fast-evolving product workflows.

Example:

```sql
status text NOT NULL
CHECK (status IN ('CREATED', 'RUNNING', 'SUCCEEDED', 'FAILED'))
```

For very rapidly evolving state, application validation may be preferable to frequently changing DB checks.

---

## 18. JSONB

Good uses:

```text
provider-specific metadata
extension payloads
non-core optional attributes
captured evidence metadata
```

Bad uses:

```text
entire relational domain model
fields constantly filtered/joined but never modeled
avoiding migrations
```

---

## 19. JSONB Indexes

Add GIN/expression indexes only for known query patterns.

Do not index every JSONB column automatically.

Measure write/storage cost.

---

## 20. Arrays

PostgreSQL arrays are allowed when the value is genuinely a bounded atomic collection and relational querying is simple.

If array members have identity/lifecycle/query semantics, model a child table.

---

## 21. Money / Decimal

Use precision-safe storage.

Typical:

```text
numeric
integer minor units + currency
```

The owning business domain defines the representation.

Never use floating point for accounting semantics.

---

## 22. Binary Data

Large artifacts/files belong in object storage when appropriate.

PostgreSQL stores:

```text
metadata
checksum
object reference
ownership
state
```

Do not turn the primary database into an unbounded artifact store.

---

## 23. Secrets

Store references/metadata rather than raw secret values when possible.

If secret material must be persisted, it requires security review, encryption strategy, rotation, and access controls.

---

## 24. Repository Boundary

Repositories expose domain-focused persistence operations.

Good:

```go
GetByID
Create
UpdateStatus
ListByProject
```

Avoid leaking raw SQL/query builder details into services.

---

## 25. Repository Interfaces

Interfaces live near the business consumer/owning service contract.

Concrete PostgreSQL implementation lives in infrastructure.

Business services depend on interfaces, not `pgx/sql.DB` directly.

---

## 26. Query Parameterization

Never concatenate untrusted values into SQL.

Good:

```sql
WHERE application_id = $1
```

Dynamic fields such as sort keys come from an allowlist/strategy map.

---

## 27. `SELECT *`

Avoid `SELECT *` in production repository queries.

Select the columns intentionally needed.

This improves compatibility, clarity, and scan cost.

---

## 28. N+1 Queries

Request/business workflows should be reviewed for N+1 behavior.

Use:

```text
set-based queries
batch lookup
joins within same service ownership
preloading
```

when appropriate.

---

## 29. Joins

Joins are allowed within one service-owned schema/model.

Cross-service joins are forbidden.

---

## 30. Pagination

Large queries paginate.

Prefer cursor/keyset pagination for large mutable collections.

A cursor must be opaque to API consumers.

---

## 31. Stable Ordering

Paginated queries must have deterministic order.

Example:

```text
ORDER BY created_at DESC, application_id DESC
```

The cursor contains/represents the stable ordering key.

---

## 32. Index Design

Index for actual query patterns.

Consider:

```text
WHERE columns
JOIN keys
ORDER BY
uniqueness
partial filters
tenant scope
```

Avoid speculative index explosion.

---

## 33. Composite Indexes

Column order matters.

Design based on real filters/order.

Document non-obvious indexes close to the migration or schema documentation.

---

## 34. Partial Indexes

Use when a stable high-value subset is queried frequently.

Example:

```text
active/non-deleted resources
pending executions
```

Do not use partial indexes with volatile/unclear predicates.

---

## 35. Unique Constraints

Use unique constraints for real business uniqueness.

Map violations to semantic `Conflict` errors.

Do not implement uniqueness through "check then insert" only.

---

## 36. Check Constraints

Protect simple invariants.

Examples:

```text
non-negative count
valid bounded state
end time >= start time
```

Complex policy remains in business code.

---

## 37. Transactions

Transaction ownership belongs to the business use case.

```text
Business Service
   ↓
transaction boundary
   ↓
one/more repositories from same service
```

Repositories must not independently open hidden transactions that prevent use-case atomicity unless explicitly designed.

---

## 38. Transaction Scope

Keep transactions as short as practical.

Do not:

```text
call external network service while holding DB transaction
wait for human approval in transaction
run long computation inside transaction
```

Persist state, commit, then continue durable workflow.

---

## 39. Isolation

Use PostgreSQL default isolation unless a use case proves stronger isolation is required.

When using:

```text
repeatable read
serializable
explicit locks
```

document the invariant/race being protected.

---

## 40. Row Locks

Use `FOR UPDATE`/locking deliberately.

Lock in consistent order to reduce deadlock risk.

Never lock rows while waiting on remote calls.

---

## 41. Optimistic Concurrency

For editable resources, a revision/version may be preferable to pessimistic locking.

Update:

```text
WHERE id = ? AND version = ?
```

and increment version atomically.

Return semantic conflict on stale update.

---

## 42. Deadlocks

Deadlocks are possible.

Applications should:

```text
keep transactions short
use consistent lock order
surface/retry only safe operations
instrument failures
```

Do not add broad retry loops around all transactions.

---

## 43. Connection Pools

One bounded pool per database target, not per request.

Pool configuration must respect a global connection budget.

Expose:

```text
open
in-use
idle
waits
timeouts
```

metrics.

---

## 44. Prepared Statements / Driver Features

Use driver features deliberately.

Do not let per-tenant/provider dynamic behavior create unbounded prepared-statement or pool growth.

---

## 45. Query Timeouts

Every request path has a context/deadline.

Long analytical/reporting queries should use separate workload controls when needed.

---

## 46. Slow Queries

Slow-query thresholds should be observable.

Investigate:

```text
missing index
bad cardinality estimate
N+1
lock wait
large payload
unbounded filter
```

before increasing timeout.

---

## 47. Migrations

Canonical layout may be:

```text
migrations/
├── 000001_initial.up.sql
├── 000001_initial.down.sql
├── 000002_add_application_status.up.sql
└── 000002_add_application_status.down.sql
```

or the repository's selected migration tool convention.

Ordering must be deterministic.

---

## 48. Migration Ownership

Migrations live with the owning service/module or a clearly centralized migration package that preserves ownership.

Do not let one service modify another service's tables.

---

## 49. Up / Down

Provide down migrations where safe and repository policy requires them.

Some irreversible data migrations cannot truly be undone.

In those cases document forward recovery rather than inventing unsafe rollback SQL.

---

## 50. Expand / Migrate / Contract

Breaking live-schema evolution uses:

```text
1. Expand
   add compatible schema

2. Migrate
   backfill/use both

3. Contract
   remove old schema after old code is gone
```

Do not combine all three into one zero-downtime release.

---

## 51. Adding `NOT NULL`

Safe staged pattern:

```text
add nullable column
deploy code writing it
backfill
validate
add NOT NULL
remove compatibility path later
```

---

## 52. Renaming Columns

For live compatible upgrades:

```text
add new
dual read/write or copy
migrate
stop old usage
remove old later
```

Direct rename may break old binaries during rolling upgrade.

---

## 53. Dropping Data

Destructive migrations require explicit review.

Confirm:

```text
retention
backup
old-version compatibility
rollback strategy
```

---

## 54. Backfills

Large backfills should be bounded/batched.

Avoid one migration transaction that rewrites millions of rows and blocks startup.

Long backfills may belong in a controlled job.

---

## 55. Migration Execution

Development may auto-run migrations through approved scripts.

Production migration execution is controlled by deployment/release workflow.

Application replicas should not race uncontrolled schema migration.

---

## 56. Schema Drift

Production schema must be reproducible from migrations.

Manual drift is a defect.

Provide tooling/checks to detect it where practical.

---

## 57. Seed Data

Migrations create schema, not arbitrary mutable demo state.

Stable system reference data may be migrated if it is truly part of the product contract.

Development/test fixtures belong in test/bootstrap tooling.

---

## 58. Test Databases

Backend integration tests use real PostgreSQL.

Use:

```text
isolated database/schema/container
deterministic migrations
test-owned data
cleanup/disposable environment
```

Never run integration tests against a developer's production-like shared DB accidentally.

---

## 59. SQL Mocking

SQL mocks may test a narrow adapter behavior, but they do **not** replace PostgreSQL integration tests.

Features such as:

```text
constraints
JSONB
transactions
locking
indexes
PostgreSQL syntax
```

must be validated against PostgreSQL when relevant.

---

## 60. Database Fixtures

Prefer builders/factories through repositories/services.

Direct SQL fixtures are allowed for repository integration tests when the SQL itself is the boundary under test.

Do not use direct DB seeding in Playwright product regression.

---

## 61. Backup

Backups must include:

```text
database data
required encryption metadata
migration/version information
```

Store backups outside the primary failure domain.

---

## 62. Restore

A backup is not proven until restore succeeds.

Automate periodic restore drills into an isolated environment.

Verify application-level usability after restore.

---

## 63. Point-in-Time Recovery

Enterprise/production profiles should support PITR where operational requirements demand it.

Retention and RPO/RTO are deployment policy, not hardcoded in application logic.

---

## 64. Encryption

Use TLS for database connections where required.

Disk/storage encryption should be provided by deployment/storage layer.

Sensitive fields needing application-level encryption require explicit design and key management.

---

## 65. Row-Level Security

PostgreSQL RLS is not the default substitute for service authorization.

It may be used as defense in depth where a domain deliberately designs and tests it.

Authorization remains explicit in the application/service architecture.

---

## 66. Audit Data

Audit records should be append-oriented where practical.

Define:

```text
actor
action
target
result
timestamp
request/correlation ID
organization/project context
```

Retention and tamper requirements belong to security/governance policy.

---

## 67. Outbox Tables

Outbox tables are owned by the service emitting the event.

They participate in the same transaction as the state change.

Publisher state/retry metadata must be indexed for polling/claiming.

---

## 68. Job/Execution Claiming

Workers claiming DB-backed work should use safe concurrency patterns such as:

```text
SELECT ... FOR UPDATE SKIP LOCKED
```

when appropriate.

The exact pattern must preserve business semantics and retry/idempotency.

---

## 69. Data Retention

Every high-volume dataset should eventually define:

```text
retention
archive
purge
legal/audit constraints
```

Do not let logs/evidence/temporary execution data grow without bound.

---

## 70. Database Observability

Record:

```text
query latency
query errors
pool usage/waits
transaction failures
deadlocks
lock waits
migration duration
database availability
```

Do not label metrics with raw query/user IDs causing cardinality explosion.

---

## 71. Query Logging

Do not log full query parameters when they may contain secrets/PII.

Use normalized query/operation names.

---

## 72. Repository Error Mapping

Map PostgreSQL/driver errors.

Examples:

```text
no rows -> semantic not found when expected
unique violation -> conflict
serialization failure -> controlled retry candidate
timeout -> timeout
connection failure -> unavailable
```

Never return `pq`/`pgx` errors directly to clients.

---

## 73. Data Access Security

Database credentials are secret references/configured securely.

Use least privilege.

Separate migration credentials from runtime credentials when production security posture requires it.

---

## 74. Read Replicas

Do not add read replicas until workload needs them.

If introduced, define consistency requirements per query.

Never send read-after-write critical behavior to eventually consistent replica accidentally.

---

## 75. Search

Start with PostgreSQL full-text search / `pg_trgm` where suitable.

Move to external search only when scale/functionality justifies it.

Search projection is not transactional source of truth.

---

## 76. Graph Projection

Engineering graph source data may start in relational entity/edge tables in PostgreSQL.

A future graph database is a projection/derived system unless architecture deliberately changes source-of-truth ownership.

---

## 77. Vector Data

Use `pgvector` only when a concrete AI/search use case requires it.

Do not introduce vector infrastructure in Phase 0 without a use case.

---

## 78. Database PR Checklist

```text
[ ] Owning service/schema identified
[ ] No cross-service table access
[ ] Migration added
[ ] Up/down/forward recovery considered
[ ] Constraints protect important invariants
[ ] Nullability/default semantics correct
[ ] Indexes match query paths
[ ] No unnecessary PostgreSQL ENUM
[ ] JSONB use justified
[ ] Transaction ownership correct
[ ] Query parameterized
[ ] Pagination bounded
[ ] Semantic DB errors mapped
[ ] Integration test uses real PostgreSQL
[ ] Migration tested from previous schema
[ ] Observability impact considered
[ ] Destructive/locking risk reviewed
```

---

## 79. Definition of Done

```text
[ ] Schema is reproducible from migrations
[ ] Service ownership is explicit
[ ] No cross-service DB dependency introduced
[ ] Repository API remains persistence-only
[ ] Critical constraints exist
[ ] Query/index plan is reasonable
[ ] Concurrent behavior tested where relevant
[ ] Migration succeeds on clean database
[ ] Migration succeeds from supported previous state
[ ] Real PostgreSQL integration tests pass
[ ] Errors are normalized
[ ] Backup/restore implications understood
```

---

## 80. Phase 0 Deliverables

Phase 0 should establish:

```text
PostgreSQL local Compose service
migration tool/convention
one service-owned PostgreSQL schema
application-generated ID abstraction
repository interface + PostgreSQL implementation
transaction helper
normalized DB error mapping
real PostgreSQL integration-test harness
migration test
pool metrics
backup/restore development proof or documented follow-up
```

---

## 81. Agent Rules

Before a schema/repository change inspect:

```text
docs/AGENT.md
this file
owning entity
repository interface
PostgreSQL adapter
migrations
relevant queries/indexes
integration tests
public contract if persistence affects it
```

Do not design schema from a frontend request alone.

---

## 82. Final Mandatory Rules

1. PostgreSQL is source of truth for transactional data.
2. Every table has one owner.
3. Cross-service SQL/FKs/joins are forbidden.
4. Repositories are persistence-only.
5. Business services own transactions.
6. IDs are stable and generated before persistence by default.
7. Required invariants use constraints.
8. JSONB is targeted, not default modeling.
9. Queries are parameterized and bounded.
10. Indexes follow real query patterns.
11. Every schema change is migrated.
12. Destructive changes are staged.
13. Production drift is not accepted.
14. Real PostgreSQL tests are mandatory for persistence behavior.
15. Database errors become semantic errors.
16. Backups require restore proof.
17. Pool/query/migration behavior is observable.
18. Data ownership supports future extraction.

---

## 83. Guiding Principle

```text
Service Ownership
      +
Relational Integrity
      +
Safe Migration
      +
Transactional Correctness
      +
Real PostgreSQL Tests
      =
Durable System of Record
```

> **Use PostgreSQL strongly, but never let a shared database erase Perihellion's service boundaries.**

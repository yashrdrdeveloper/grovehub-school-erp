# Persistence, Transactions, ACID and Migrations

## 1. Responsibilities

```text
modules/<resource>/entity/                  # persisted fields and GORM table metadata
modules/<resource>/repository/              # typed repository contract
modules/<resource>/repository/impl/         # SQL/GORM implementation
modules/<resource>/repository/helpers/      # resource-only queries/constraints
modules/<resource>/services/operations/     # transaction + business invariants
libs/core/database/                         # generic context-bound UnitOfWork mechanics
internal/resources/migrations/               # embedded, versioned, owned DDL
```

Generic DB mechanics belong to the core database capability. Queries, unique constraint names, aggregate rules, and search projections remain resource-owned.

## 2. ACID

**Atomicity:** one aggregate command commits all rows and required outbox intent together, or rolls back entirely. **Consistency:** constraints and business invariants are preserved, not just Go validators. **Isolation:** choose a level/locking policy that withstands concurrent writes, including optimistic version checks or `SELECT ... FOR UPDATE` when appropriate. **Durability:** acknowledged commits are persistent according to actual DB configuration, WAL/backup policy, not inferred from a successful HTTP response.

ACID applies to one database transaction. It cannot make a remote provider API call part of a PostgreSQL transaction. Use outbox (for events) or explicit long-running workflow/compensation for remote effects.

## 3. Transaction ownership

The **business operation** defines the aggregate transaction boundary. Pass the same context to all participating repositories. The reusable library's `UnitOfWork` and `transaction.Database(ctx, db)` allow repository calls to participate without receiving raw transaction handles in public DTOs.

```go
err := unitOfWork.Run(ctx, func(txCtx context.Context) error {
    if err := firstRepository.Save(txCtx, first); err != nil {
        return err
    }
    if err := secondRepository.Save(txCtx, second); err != nil {
        return err
    }
    return outbox.Enqueue(txCtx, event)
})
```

This is illustrative. Actual signatures are defined by the installed library. **Never** open a new independent DB transaction halfway through an existing operation. Do not publish events before commit. `AfterCommit` callbacks are not equivalent to a durable outbox and may be lost on process crash.

## 4. Invariants and optimistic concurrency

- Persist foreign keys, unique normalized values and controlled lifecycle status with SQL constraints.
- Require optimistic version predicates when concurrent updates can overwrite one another; report conflicts as domain errors.
- Check `RowsAffected` and distinguish not-found from stale version.
- Never retry a non-idempotent command blindly.
- DB unique constraint errors are mapped to stable domain error codes by the owning repository/helper.
- Search projection updates are part of the aggregate transaction, not best-effort cache maintenance if exposed as committed state.

## 5. Idempotency

For retriable command APIs: validate a scoped idempotency key (actor/tenant/operation), compute a stable fingerprint of normalized input, record result/owner and preserve one transition. Same key + same request returns the same result; same key + different request conflicts. Concurrency cases must rely on DB constraints/locks, not an unprotected `Find` followed by `Create`. Expiration and tombstone policy is explicit.

## 6. GORM rules

- Repositories use `transaction.Database(ctx, db)` consistently.
- Select only needed columns for read queries and batch relation hydration.
- Keep trusted SQL identifiers/expressions in code; parameterize values.
- Avoid implicit `Save` for a write when insert/update and version semantics differ.
- Do not expose `*gorm.DB` in service or client contracts.
- Avoid hidden cascading deletes/association mutations unless documented and tested.
- Pagination uses stable cursor/keyset with resource-owned configuration.
- Never log raw SQL bindings containing credentials or PII.

## 7. Migration contract

1. Every module has its own ordered immutable migrations and Goose history.
2. Migration files are embedded in `internal/resources/`, and `Migrate` delegates to them.
3. An applied migration **must never be edited, deleted, or rewritten** even if a later migration replaces its tables.
4. Use a new forward-only migration for indexes, constraints, new tables, backfills and data transfers.
5. For destructive changes: expand -> backfill -> validate -> switch readers/writers -> contract, with safe rollout.
6. Check migration ordering when schemas are shared; stop dependent workers before relocating live queues.
7. Schema changes include both existing-data and clean-install test paths, plus rollback/restore plan.

## 8. Backup, recovery and retention

Backup strategy must define RPO, RTO, offsite durability, encryption, retention, periodic restore tests and schema compatibility. Table retention is not broker retention, and outbox row cleanup is not an audit or compliance archive. Verify `SENT` rows are eligible for cleanup while pending/retrying messages are protected.

## 9. Verification

Test one transaction rollback across two repositories, event enqueue atomicity, duplicate concurrent commands, idempotency hash mismatch, version conflict, constraint error mapping, migration from real previous schema, and crash-window delivery. Use an actual PostgreSQL instance for ACID and SQL behavior; an in-memory fake does not prove transaction isolation or SQL constraints.

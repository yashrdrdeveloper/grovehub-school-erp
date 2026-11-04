# Identity registration transaction

Route: `POST /auth/register` with no request GID.

1. Normalize and validate registration and idempotency key.
2. Start one shared GORM/PostgreSQL transaction using `core/database/transaction.Run`.
3. Check idempotency, resolve Foundation Person by normalized email, update it or create it.
4. Create linked Identity Account and password credential using the same context-scoped transaction.
5. On failure roll back **all** Person, Account and credential database writes; discard queued domain events.
6. After successful commit, run queued callbacks for Foundation/Identity events, and return the Account with the Person GID.

## Code ownership

- Identity orchestrator: `backend/services/identity/services/impl/auth.go`.
- Transaction capability: `backend/libs/core/database/transaction/transaction.go`.
- Repositories use `transaction.Database(ctx, repository.db)` for all database access. Foundation's nested contact transactions use GORM savepoints.
- Event decorators use `transaction.AfterCommit(ctx, callback)`; no publish before commit.

## Limitations

Database atomicity requires the same PostgreSQL connection; this is not a distributed transaction. NATS publication is still best effort after commit, not an outbox. Password setup/reset by email, email token and mail outbox remain unimplemented. Existing-Person name changes by email alone require ownership verification before public use. PostgreSQL integration tests (rollback, commit, idempotency and concurrency) remain necessary.

After-commit callbacks use `transaction.DetachedContext(ctx)` to retain trace metadata without retaining the old transaction marker; merely calling `context.WithoutCancel` would re-enqueue the publish attempt.

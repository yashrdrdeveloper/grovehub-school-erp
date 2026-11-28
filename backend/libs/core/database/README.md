# Core Database

The only consumer import is `github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database`. Database provisioning, migrations, business repositories and tenant authorization stay in their owning applications.

```text
database/
├── database.go
└── internal/
    ├── src/
    │   ├── postgres/
    │   │   ├── gorm.go
    │   │   ├── pool.go
    │   │   ├── constants/
    │   │   ├── helpers/
    │   │   └── validators/
    │   └── transaction/
    │       ├── transaction.go
    │       ├── unit_of_work.go        # interface
    │       ├── helpers/context.go
    │       └── impl/unit_of_work.go
    └── resources/tests/unit-test/
```

- `NewDB`: opens a GORM PostgreSQL connection.
- `OpenProjectPool`: validates an explicit database URI, rejects identity overrides and verifies the connected role/database. Connection failures do not disclose credentials.
- `NewUnitOfWork` and `RunTransaction`: coordinate context-bound transactions, including nested calls.
- `TransactionDB` and `TransactionActive`: resolve and inspect the active transaction.
- `AfterCommit`: registers callbacks that run only after successful commit.
- `DetachedContext`: preserves context values while removing transaction state and cancellation.

The UnitOfWork interface is defined in transaction/unit_of_work.go beside impl/unit_of_work.go. The public facade embeds the interface without exposing the implementation. The implementation constructor returns the interface and includes a compile-time conformance check. No artificial interfaces are required for pure validation and mapping functions. A transaction context is intended for sequential repository operations, not concurrent use.

## Testing

Public unit tests: `tests/unit-test/database/` in Core. Private unit tests: `database/internal/resources/tests/unit-test/`. Unit tests use Ginkgo, Gomega and sqlmock; they do not connect to PostgreSQL.

Real database isolation checks: `tests/sanity/database/postgres/` in Core. CI provisions disposable project databases and passes the three project test DSNs. If all are absent, the suite skips; partial configuration fails.

From `backend/libs/core`:

```bash
go test ./database/... ./tests/unit-test/database/...
go test -count=1 -v ./tests/sanity/database/postgres
```

# GroveHub GORM Repository Guidelines

> **Status:** Mandatory GroveHub backend persistence standard  
> **Applies to:** every Go service under `backend/services/*` that persists data in PostgreSQL  
> **ORM:** GORM v2 with `gorm.io/driver/postgres`  
> **Boundary:** GORM is private to `repository/impl` and shared database bootstrap code.

---

## 1. Objective

GroveHub uses GORM to remove repetitive PostgreSQL CRUD/scanning code without allowing persistence concerns to leak into the domain.

The architecture remains:

```text
controller
   ↓
service
   ↓
repository interface
   ↓
repository/impl
   ↓
GORM
   ↓
PostgreSQL
```

GORM does not replace repository boundaries, domain entities, service validation, or SQL migrations.

---

## 2. Canonical Placement

```text
backend/services/<service>/
├── entity/
├── repository/
│   ├── <resource>.go
│   └── impl/
│       ├── <resource>.go
│       ├── <resource>_model.go
│       └── <resource>_mapper.go
└── migrations/
```

Responsibilities:

```text
repository/<resource>.go
    public repository interface
    repository sentinel errors
    domain/search contracts only

repository/impl/<resource>.go
    GORM queries
    transactions local to one aggregate persistence operation
    constraint/error translation
    validated filter/sort/cursor execution

repository/impl/<resource>_model.go
    persistence-only structs
    GORM tags
    table/column mapping

repository/impl/<resource>_mapper.go
    explicit persistence model <-> domain entity mapping
```

Do not create a second `postgres/` repository package under a service. The concrete PostgreSQL/GORM adapter is `repository/impl`.

---

## 3. Hard Boundary Rules

The following packages may import `gorm.io/gorm`:

```text
backend/libs/core/database/postgres
backend/services/<service>/repository/impl
```

The following must not import GORM:

```text
entity
repository interface package
services
controller
routes
client
DTOs
public mappers
```

Domain entities must not embed `gorm.Model`, contain GORM tags, expose database column names, contain persistence hooks, or know schema/table names.

Repository interfaces must use domain/shared contract types, not `*gorm.DB`, GORM clauses, or persistence models.

---

## 4. Persistence Models Are Not Domain Entities

Use separate persistence models inside `repository/impl`.

```go
type personModel struct {
    GID       string    `gorm:"column:gid;primaryKey"`
    FirstName string    `gorm:"column:first_name"`
    CreatedAt time.Time `gorm:"column:created_at"`
}
```

The domain `entity.Person` remains persistence-agnostic.

Use explicit mapper functions between the two representations. Do not use reflection-based generic copying for domain persistence mapping.

This is the Go equivalent of keeping a deliberate MapStruct-style mapping boundary: the mapping is compile-time-visible, reviewable, and field-by-field.

---

## 5. Database Bootstrap

Create one GORM database handle for the PostgreSQL target at application startup.

```text
launcher/bootstrap
      ↓
libs/core/database/postgres
      ↓
*gorm.DB
      ↓
service composition roots
      ↓
repository/impl constructors
```

Do not open a new database handle per request or repository call. Use the underlying `database/sql` pool managed by GORM for connection pooling and shutdown.

Every request-path repository operation must use `db.WithContext(ctx)`.

---

## 6. Migrations Remain the Source of Truth

GORM `AutoMigrate` is forbidden for GroveHub production/runtime schema evolution.

Canonical schema ownership remains:

```text
backend/services/<service>/migrations/
```

Use migrations for tables, columns, constraints, indexes, extensions, and deliberately required database functions/triggers.

Persistence model tags describe application mapping; they are not the authoritative production schema definition.

---

## 7. CRUD Rules

Prefer GORM's normal model/query APIs for ordinary persistence:

```text
Create
First / Take
Find
Updates
Delete
Pluck
```

Do not use raw SQL merely because the previous implementation used PGX directly.

Raw SQL or expressions are allowed only where PostgreSQL-specific behavior is materially clearer or required, such as search projection refresh, special PostgreSQL expressions, complex keyset predicates, or advanced aggregates.

Keep those fragments inside `repository/impl` and parameterize all values.

---

## 8. Updates

For domain updates, prefer explicit columns or maps over `Save`.

This prevents accidental full-row writes, makes zero-value behavior explicit, avoids unintended association persistence, and keeps changed columns reviewable.

Do not let GORM automatically persist aggregate associations unless the repository intentionally owns that behavior and it is tested.

---

## 9. Aggregate Transactions

A single repository operation may open a transaction when one aggregate spans multiple tables and must be stored atomically.

For example, Person + `person_contacts` + search projection refresh must either all succeed or all roll back.

Use `db.WithContext(ctx).Transaction(...)`.

For a future use case that needs one transaction across multiple repositories, promote the boundary to an explicit service/unit-of-work abstraction rather than nesting hidden independent transactions.

Keep transactions short and never perform external network calls while holding one.

---

## 10. Constraints and Semantic Errors

Database constraints remain authoritative for concurrency-safe uniqueness and integrity.

Repository implementations translate known PostgreSQL constraint failures into repository sentinel errors.

Example:

```text
uq_person_contacts_type_normalized_value
        ↓
ErrPersonContactAlreadyExists
```

Do not implement uniqueness only as `SELECT` followed by `INSERT`, because concurrent requests can race.

Do not expose raw PostgreSQL/GORM errors through service/API contracts when a semantic repository error exists.

---

## 11. Search

Shared search request/query/result types remain generic and entity-neutral in `libs/core/contracts/search` and `libs/core/search`.

The owning service validates public fields/operators. The repository implementation maps approved public semantics to database expressions.

```text
firstName      -> LOWER(p.first_name)
createdAt      -> p.created_at
contacts.email -> EXISTS person_contacts(...)
```

Never accept a client field name and concatenate it directly into `ORDER BY`, filters, table names, or column expressions.

Sort expressions and relation mappings must come from repository-owned allowlists/registries. Complex GORM `Where` expressions are acceptable when SQL fragments are fixed by code and user values remain bound parameters.

---

## 12. Cursor Pagination

Cursor/keyset pagination remains deterministic. For non-unique sort fields, use the resource ID as the stable tie-breaker.

```text
ORDER BY created_at DESC, gid DESC
```

The repository applies the matching keyset condition and fetches `limit + 1` rows to determine `HasNext`. GORM does not change the public cursor contract.

---

## 13. Associations and N+1

Do not automatically `Preload` every relation.

Choose deliberately between `Preload`, same-service joins, batch `IN` lookup, or a separate aggregate query.

For list/search paths, prefer bounded set-based/batch loading over one relation query per row.

---

## 14. Naming and Schemas

Service-owned PostgreSQL schemas remain explicit, for example `foundation.persons` and `foundation.person_contacts`.

Persistence models should define explicit column mapping where it improves safety and clarity. Do not rely on GORM naming conventions when an existing migration/schema contract already defines the name.

---

## 15. Delete Behavior

GroveHub does not get soft delete merely because GORM supports it.

Do not add `gorm.DeletedAt` unless the domain explicitly requires soft-delete semantics and migrations/API behavior are designed for it.

Existing database foreign-key delete behavior remains controlled by migrations.

---

## 16. GORM Hooks

Avoid GORM lifecycle hooks for business behavior. Business rules belong in services/domain helpers.

Hooks may be used only for narrowly scoped persistence concerns when there is a clear reason and tests cover the behavior. No hidden workflow orchestration belongs in `BeforeCreate`, `AfterCreate`, `BeforeUpdate`, or `AfterUpdate`.

---

## 17. Generic Repository Rule

Do not create a generic `BaseRepository[T]` merely because GORM supports generics.

Repositories remain domain/use-case focused, such as `PersonRepository`, `StudentRepository`, and `AdmissionRepository`.

Promote only truly entity-neutral technical helpers to `libs/core`. Avoid abstractions that erase domain semantics or turn every query into a generic CRUD API.

---

## 18. Testing

Repository integration tests should use real PostgreSQL behavior when testing constraints, transactions, PostgreSQL casts, `ILIKE`/trigram behavior, composite keys, partial unique indexes, and cursor ordering.

Service unit tests continue to mock/fake repository interfaces rather than GORM. Tests outside `repository/impl` should not need to know that GORM is used.

---

## 19. Current Foundation Pattern

Foundation is the reference implementation:

```text
repository/person.go
    PersonRepository contract

repository/impl/person.go
    GORM repository implementation

repository/impl/person_model.go
    GORM-only Person/contact persistence structs

repository/impl/person_mapper.go
    explicit entity <-> persistence mapping
```

Person creation/update remains atomic across `foundation.persons`, `foundation.person_contacts`, and `search_details` refresh.

The public service/controller/client contracts do not change because of the ORM migration.

---

## 20. Review Checklist

- [ ] GORM imports are contained to allowed persistence packages.
- [ ] Domain entities have no GORM tags or `gorm.Model`.
- [ ] Repository interfaces contain no GORM types.
- [ ] SQL migrations remain the schema source of truth.
- [ ] Every request-path query uses `WithContext(ctx)`.
- [ ] Dynamic identifiers come only from allowlists.
- [ ] User values are parameterized.
- [ ] Multi-table aggregate writes are atomic.
- [ ] Known constraints map to semantic repository errors.
- [ ] List/search relation loading is not N+1.
- [ ] Updates explicitly control changed columns.
- [ ] No accidental soft-delete behavior was introduced.
- [ ] No business rules were moved into GORM hooks.
- [ ] Repository tests cover PostgreSQL-specific behavior where relevant.

---

## 21. Final Rule

> **Use GORM to simplify persistence implementation, not to collapse GroveHub's domain boundaries. Repository contracts remain stable; GORM stays behind `repository/impl`.**

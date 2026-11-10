# GroveHub Core Library Guidelines

> **Status:** Mandatory GroveHub backend standard  
> **Applies to:** `backend/libs/core` and every backend module that consumes core libraries  
> **Goal:** capability-owned, encapsulated shared libraries with one stable public entry point per capability

---

## 1. Capability-First Rule

`backend/libs/core` is organized by **capability**, never by generic technical role.

Correct:

```text
backend/libs/core/
├── search/
├── diff/
├── events/
├── telemetry/
├── logging/
├── errors/
├── normalize/
├── id/
├── http/
└── database/
```

Incorrect:

```text
backend/libs/core/
├── contracts/
├── helpers/
├── dto/
├── mappers/
├── validators/
├── utils/
└── common/
```

A contract, helper, mapper, validator, DTO, provider adapter, or implementation exists **inside the capability that owns it**.

For example, search contracts belong to `search/contracts`, not `core/contracts/search`.

---

## 2. Canonical Capability Structure

A capability may use the following internal layout when those responsibilities exist:

```text
backend/libs/core/<capability>/
├── <capability>.go        # public façade / stable API
├── contracts/             # capability-owned data/contracts
├── helpers/               # implementation mechanics / pure helpers
├── mappers/               # representation/config/name mapping
├── validators/            # capability-level technical validation
└── <provider>/             # optional provider-specific adapter, when justified
```

Example:

```text
backend/libs/core/search/
├── search.go
├── contracts/
│   └── search.go
├── helpers/
│   └── query.go
├── mappers/
│   └── cursor.go
└── validators/
    └── query.go
```

Do **not** create empty folders merely to match the template. `contracts`, `helpers`, `mappers`, and `validators` exist only when the capability has that responsibility.

---

## 3. Root Façade Rule

The root file is the supported consumer API.

Examples:

```text
search/search.go
telemetry/telemetry.go
logging/logging.go
errors/errors.go
normalize/normalize.go
id/id.go
```

Service/module code imports the capability root:

```go
import coresearch "grovehub/libs/core/search"
import coreerrors "grovehub/libs/core/errors"
import corelogging "grovehub/libs/core/logging"
import coretelemetry "grovehub/libs/core/telemetry"
import corenormalize "grovehub/libs/core/normalize"
import coreid "grovehub/libs/core/id"
```

Service/module code must not import implementation internals such as:

```text
grovehub/libs/core/search/helpers
grovehub/libs/core/search/contracts
grovehub/libs/core/telemetry/helpers
grovehub/libs/core/errors/contracts
```

Those packages are implementation details of their owning capability.

The capability root may use aliases/wrappers to expose stable types/functions while allowing the internal layout to change without forcing repository-wide import churn.

---

## 4. Internal Folder Responsibilities

### `contracts/`

Contains capability-owned data structures, value types, configuration contracts, and generic result/input types.

Contracts must not perform I/O or own orchestration.

### `helpers/`

Contains reusable implementation mechanics owned by that capability.

Examples:

```text
search predicate/query construction
logging handler implementation
telemetry runtime/metric/trace mechanics
UUID generation implementation
normalization algorithms
```

Helpers are not a dumping ground. If a helper belongs to another capability, move it there.

### `mappers/`

Contains deterministic representation mappings where a real mapping responsibility exists.

Examples:

```text
opaque cursor serialization/deserialization
telemetry service-name -> service metadata mapping
```

Do not create a mapper package for identity functions or simple pass-throughs.

### `validators/`

Contains domain-neutral technical validation owned by the capability.

Examples:

```text
search query configuration validation
telemetry sample-ratio/export-interval validation
```

Business/domain validation remains in the owning service validator.

---

## 5. Current Core Layout

The intended current shape is:

```text
backend/libs/core/
├── database/
│   └── postgres/
├── errors/
│   ├── contracts/
│   ├── helpers/
│   └── errors.go
├── http/
│   ├── contracts/
│   ├── headers/
│   ├── logging/
│   ├── request/
│   ├── response/
│   └── tracing/
├── id/
│   ├── helpers/
│   └── id.go
├── logging/
│   ├── contracts/
│   ├── helpers/
│   └── logging.go
├── normalize/
│   ├── helpers/
│   └── normalize.go
├── search/
│   ├── contracts/
│   ├── helpers/
│   ├── mappers/
│   ├── validators/
│   └── search.go
└── telemetry/
    ├── contracts/
    ├── helpers/
    ├── mappers/
    ├── validators/
    └── telemetry.go
```

`database/postgres` is valid because PostgreSQL is a provider implementation **inside** the database capability. HTTP request/response/tracing packages are valid sub-capabilities because they remain owned by `http` rather than becoming repository-wide buckets.

---

## 6. Search Example

Consumers import only:

```go
import coresearch "grovehub/libs/core/search"
```

The root façade exposes both HTTP-facing search contracts and the validated repository search model with distinct names:

```text
Request
Filter
Sort
Result[T]

QueryFilter
QuerySort
Cursor[TId]
Query[TId]
QueryResult[T]
```

This avoids two unrelated packages both being imported as `search` and makes the boundary explicit:

```text
search.Request
    ↓ resource validator
search.Query[TId]
    ↓ repository
search.QueryResult[Entity]
    ↓ service
search.Result[Response]
```

Generic query mechanics remain inside the search capability.

---

## 7. Errors, Logging, and Telemetry

The same façade rule applies:

```text
core/errors
    -> public Code, Definition, Error and constructors/wrappers
    -> internal contracts/helpers remain hidden from services

core/logging
    -> public Levels, parsing and handler construction
    -> internal contracts/helpers own implementation

core/telemetry
    -> public service metadata, initialization, tracing and metric functions
    -> internal contracts/helpers/mappers/validators own implementation
```

Do not make services know how telemetry runtime state, logger filtering, or application-error representation is internally implemented.

---

## 8. ID and Normalization

These are independent capabilities, not members of a generic `helpers` package.

Use:

```go
coreid "grovehub/libs/core/id"
corenormalize "grovehub/libs/core/normalize"
```

Do not use:

```text
grovehub/libs/core/helpers/id
grovehub/libs/core/helpers/normalize
```

The implementation may live under each capability's `helpers/`, but consumers use the capability root.

---

## 9. HTTP Contracts

Transport response contracts belong to the HTTP capability:

```text
backend/libs/core/http/contracts/
```

They must not live in a repository-wide `core/dto` bucket.

The public response writer remains:

```text
backend/libs/core/http/response
```

Controllers use HTTP-facing packages; business services do not depend on HTTP contracts.

---

## 10. Tests

Core tests should mirror capability ownership rather than obsolete technical-role paths.

Preferred:

```text
backend/libs/core/tests/unit-test/diff/
backend/libs/core/tests/unit-test/id/
backend/libs/core/tests/unit-test/logging/
backend/libs/core/tests/unit-test/telemetry/
backend/libs/core/tests/unit-test/http/
```

Do not preserve paths such as `tests/unit-test/helpers/id` after `helpers/id` no longer exists as a public capability.

---

## 11. Adding New Core Code

Before adding code under `libs/core`, answer these questions:

1. Is it truly domain-neutral and shared?
2. Which capability owns it?
3. Can it live behind an existing capability root?
4. Is a new internal `contracts/helpers/mappers/validators` package actually necessary?
5. Can consumers avoid importing the internal package directly?

If there is no clear capability owner, the code is probably not ready for `libs/core`.

---

## 11.1 Generic JSON Diff Capability

`backend/libs/core/diff` is the domain-neutral public comparison API. Consumers use:

```go
changes, err := diff.Compare(old, next)
```

Both values have the same inferred generic Go type; no Foundation dependency or manual datatype string is required. `Compare[T]` returns `([]diff.Change, error)`, where each change has `field`, `oldValue`, and `newValue`. Paths are flat and JSON-named, such as `contacts.emails[0].value`. Change ordering is deterministic; arrays compare by index, not by domain identifier. Missing properties use JSON null in the result, and numeric values decode as `json.Number`. Consumers should check errors.

Core diff owns comparison mechanics only. Source types remain in their domain (for example, Foundation `dto/response.Person`). Do not create Foundation-specific types inside core, automatically publish a diff as part of an event, or place diff inside messaging infrastructure.

## 11.2 Database Transaction Capability

`backend/libs/core/database/transaction` provides technical GORM transaction scoping for operations spanning modules in the same PostgreSQL database:

- `transaction.Run(ctx, db, operation)` starts a transaction; a nested `Run` joins the active one.
- `transaction.Database(ctx, repository.db)` selects the active transaction for repository reads/writes or the ordinary DB outside a transaction.
- `transaction.AfterCommit(ctx, callback)` queues work to run only after successful commit. No callbacks execute on rollback or failed commit.
- `transaction.DetachedContext(ctx)` retains trace metadata while masking the committed transaction; publishers use it to prevent re-queueing events.

GORM remains an implementation detail of core database and repository/orchestration construction; do not expose a transaction handle inside Foundation/Identity API DTOs or event contracts. This capability does not coordinate remote services or provide guaranteed event delivery. See `docs/engineering/backend/REGISTRATION_TRANSACTION.md`.

## 12. Hard Rules

```text
NO top-level libs/core/contracts
NO top-level libs/core/helpers
NO top-level libs/core/dto
NO generic libs/core/utils or libs/core/common dumping ground
NO service importing a core capability's internal contracts/helpers/mappers/validators
NO duplicate public entry points for the same capability
NO domain/service-specific behavior inside libs/core
NO empty internal folders created only for symmetry
NO provider name at the repository-wide root when it belongs to a capability
YES one stable root façade per capability
YES internal implementation grouped under the owning capability
YES capability tests aligned with capability ownership
```

> **Guiding principle:** organize shared code by what it **does**, keep its implementation details inside that capability, and expose one stable root package to the rest of GroveHub.

## Shared generic domain-event contract

`backend/libs/core/events/envelope.go` defines `Envelope[T any]`, a transport-agnostic `Metadata` struct, and generic `NewEnvelope(ctx, metadata, payload)`. The stable top-level JSON keys are `eventId`, `eventType`, `version` (event **schema** version), `aggregateId`, `occurredAt`, `correlationId`, `causationId`, and `data` (typed payload). Correlation ID may be explicitly propagated using `WithCorrelationID(ctx, id)` or derived from the current trace; causation ID is explicitly propagated with `WithCausationID(ctx, id)`. Unknown optional IDs must not be fabricated. Domains keep their event names, operation identifiers and payload structs under their own `services/<domain>/events/` packages. See `docs/engineering/backend/EVENT_GUIDELINES.md`.

---

## Entity Builder (No Code Generation)

`backend/libs/core/entity` offers an optional, domain-neutral generic builder directly at runtime. Go cannot define a new `Entity<Table> Name { ... }` syntax or automatically synthesize typed `Get<Field>` and `Set<Field>` methods from a struct without code generation; GroveHub does not require generators for this capability.

Define entities as ordinary passive Go structs, preserving existing GORM annotations and `TableName()` mappings where required by the service's schema. Construct values without creating a per-entity builder:

```go
credential := coreentity.New[entity.PasswordCredential]().
    With(func(v *entity.PasswordCredential) {
        v.AccountId = accountId
        v.PasswordHash = passwordHash
        v.PasswordChangedAt = now
        v.Version = 1
        v.CreatedAt = now
        v.UpdatedAt = now
    }).
    Build()
```

`New[T]()`, `From(value)`, `With(func(*T))`, and `Build()` are the complete stable core API. `With` acts on a copy of a value struct; it does not implement domain validation, permissions, clock/ID generation, persistence, or lifecycle rules. Struct fields must remain accessible to their callers, and nested pointer/slice/map fields are not deep-copied. Service orchestration continues to own all business decisions. No `go generate` command or generated files are needed. Do not use generic reflection-based field-name setters that lose compile-time type checking.

For compatibility with current integrations, legacy entity constructors/accessors may remain until their consumers are deliberately migrated. For new entities, do not add trivial constructors or accessors solely to implement generic construction; prefer `coreentity.New[T]()` or a direct struct literal when more readable.

Run the dedicated Ginkgo suite with `cd backend && go test ./libs/core/tests/unit-test/entity/...`.

# GroveHub Current Backend Structure

> **Status:** Mandatory GroveHub architecture standard  
> **Applies to:** the current Go modular monolith and all new backend services/modules  
> **Precedence:** this file is authoritative for GroveHub backend placement and layer ownership.

Also read when relevant:

```text
docs/engineering/backend/EVENT_GUIDELINES.md
docs/engineering/backend/CORE_LIBRARY_GUIDELINES.md
docs/engineering/backend/ANNOTATION_FRAMEWORK_MASTER_PLAN.md  # Phase 1 definitions; Phase 2 partially migrated
docs/engineering/backend/ANNOTATION_PHASE1_REFERENCE.md  # complete directive contracts
docs/engineering/backend/ANNOTATION_PHASE2_MIGRATION.md  # implementation status and acceptance gates
docs/engineering/backend/MAPPER_AND_DTO_GUIDELINES.md
docs/engineering/backend/VALIDATOR_GUIDELINES.md
docs/engineering/backend/NORMALIZATION_AND_DRY_GUIDELINES.md
docs/engineering/backend/SEARCH_API_GUIDELINES.md
```

---

## Integration-backed service design rules (mandatory for all backend services)

These rules govern new services and refactors, not a particular Integration implementation phase.

### Public contracts and implementation boundaries
- Expose stable **interfaces, DTOs, domain contracts and reusable types** at package roots (`providers/providers.go`, `services/connection.go`, `repository/watch.go`, `controller/connection.go`, `config/config.go`).
- Put stateful/concrete implementations under the corresponding `impl/` folder. Provider adapters use `providers/<provider>/impl/` and an intentionally small public `providers/<provider>/` facade. Module wiring belongs in `<service>.go`, not in public contracts.
- A provider registry is itself an interface. Registry map/sort/dispatch, vendor routing and constructor logic belong under `providers/impl/`. Declare **only implemented capabilities**; use optional action or OAuth interfaces rather than forcing all providers to implement mail, webhook, send, or read functions.
- Put pure reusable helpers in `helpers/` of the owning domain, grouped by resource where appropriate (e.g., `services/helpers/connections/`). Keep concrete `services/impl/connections.go` focused on orchestration and persistence decisions; move stable ID validation and entity-to-response transformations into `services/helpers/connections/`. Reuse `core/normalize` for shared UUID/email/phone normalization instead of reimplementing parsing. Do not add helpers solely to rearrange one-line calls. If a data structure is reused across multiple layers, put its stable public definition in the **owning package** or `types/` when there is no clearer owner; do not create duplicate copies in each layer.
- All action/capability/provider keys and status literals should be defined once in capability-specific `constants/` subpackages (for example `constants/actions/gmail.go`, `constants/providers/providers.go`, `constants/connections/status.go`). Preserve type-safe aliases in a public contract only if needed for existing consumers. Keep typed error definitions/sentinels in `constants/error/`, with transport-agnostic errors mapped to HTTP only in controllers.

### Controller and route organization
- Split controllers by the owning **resource**, rather than by each action: put connection CRUD/OAuth/tests/watch/actions methods in `controller/connections.go` and `controller/impl/connections.go`; use `controller/provider.go` and `controller/impl/provider.go` for provider discovery. Compose public contracts when needed.
- Put reusable controller authorization/error/response mapping into an owning common/helper file; do not duplicate it in every controller.
- Keep `routes/http.go` as the **composition entry point**. Annotate the owning public controller interfaces with `@Routes` and HTTP verb directives; `libs/core/annotation/generation/routes` generates service-local route registrations into `.generated/routes/`. Do not recreate the removed manual `routes/connections.go`, `routes/providers.go` or `routes/helpers/http.go` registrations. Controllers retain all request handling, security and response logic; routing annotations define method/path registration only.

### Provider strategy and internal clients
- Separate provider-discovery service contracts (`services/providers.go`) from connection lifecycle contracts (`services/connections.go`); provide independent concrete service implementations.
- For a multi-vendor capability, use `providers/<capability>/<capability>.go` for a stable strategy/factory contract, `providers/<capability>/impl/` for generic factory implementation, and `providers/<capability>/providers/<vendor>/` for each vendor adapter. Prefer provider-ID/capability lookups to vendor switches.
- Every module consumed by sibling services should expose an in-process interface via `client/client.go` and a constructor with injected dependencies in `client/impl/client.go`. For supported simple forwarding, `@Client`/`@Proxy` generate the delegates into `.generated/client/impl/`; business rules must remain in their service layer. Do not loop back through HTTP inside the monolith.

### Configuration
- `config/config.go` declares the module's configuration structures and a `Loader` interface. `config/impl/config.go` loads the environment, normalizes and validates optional provider settings. Pure parsing helpers may live in `config/helpers/` when reused.
- Never load secrets in controller/request code. Load at composition startup; inject configured provider clients; when optional secrets are absent, disable those capabilities without exposing secret values.

### GORM ORM-first repository convention
- Use GORM models, `Model`, `Where(&typedStruct)`, associations (`Preload`), `Find/First`, `OnConflict`, and `Transaction` for ordinary CRUD/upsert/joins. A GORM `Where` is an ORM query builder, **not** equivalent to hand-writing a full SQL query.
- Do **not** hand-build `Table(...).Select(...).Joins(...).Scan(...)` when an entity relation and `Preload` express the same operation clearly.
- A small SQL condition inside `Where` for non-equality operators (such as `expires_at < ?`) is appropriate. Raw SQL is allowed only when a genuine database-specific feature cannot be expressed cleanly in GORM; document why.
- Define foreign-key relations on entity structs only if they express the owning service's persistence relation, not to duplicate foreign-domain data.
- Scope all tenant/owner queries correctly. Do not list another tenant's data and filter in application code if the database can scope the query.

### Events
- Define a **payload contract** alongside an alias of `coreevents.Envelope[Payload]` for each versioned service event. Use `core/events.NewEnvelope` via pure event helpers for event identity/correlation metadata.
- Put publication/serialization in `events/<resource>/impl/`, rather than in normal business-service implementation. Delivery and acknowledgment semantics must be explicit; failure to publish a provider notification must not be reported as successful delivery.
- Provider notification events contain references/hints, not copied messages or secret data. Each consuming domain owns its own processing checkpoint and deduplication.

### Schema and development commands
- A **brand-new service on an unapplied development branch** may use one `000001` schema baseline. After migrations have been applied/shared, **never squash or renumber them** without a documented upgrade/migration plan; resetting a database deletes data.
- `make dev-start` must start the platform's infrastructure, apply every registered service's migrations, load module environment settings and run the backend/Swagger. Update `scripts/database/migrate.sh`, `scripts/config/services.sh`, and `scripts/lib/env.sh` when adding a service or its new required runtime variables.
- Document service boundaries and shared rules here, not in a new phase document for every change. Update per-service documentation only for genuine API/schema/operational contracts.

## 1. Canonical Service Structure

GroveHub uses flat service roots. Do not add `src/` or service-local launchers.

```text
backend/services/<service>/
├── <service>.go
├── api/
├── config/
├── constants/
│   ├── error/
│   └── events/
│       └── operation.go
├── dto/
│   ├── request/
│   └── response/
├── entity/
├── events/
│   └── <resource>/
│       ├── <resource>.go
│       ├── impl/
│       │   └── <resource>.go
│       └── helpers/
│           └── <resource>.go
├── mappers/
│   ├── <resource>.go
│   └── impl/
│       └── <resource>.go
├── validators/
│   ├── <resource>.go
│   ├── impl/
│   │   └── <resource>.go
│   └── helpers/
│       └── <resource>.go
├── repository/
│   ├── <resource>.go
│   └── impl/
│       └── <resource>.go
├── services/
│   ├── <resource>.go
│   ├── impl/
│   │   └── <resource>.go
│   └── helpers/
│       └── <resource>.go
├── controller/
│   ├── <resource>.go
│   ├── impl/
│   │   └── <resource>.go
│   └── helpers/
│       └── <resource>.go
├── routes/
├── client/
│   ├── <resource>.go
│   └── impl/
│       └── <resource>.go
├── migrations/
└── tests/
```

Create only folders/files that have a current responsibility. Do not create empty event, mapper, validator, helper, or implementation packages merely for symmetry.

---

## 2. Contract / Implementation Rule

Interfaces and public contracts belong at the package boundary. Concrete implementations belong under `impl/`.

```text
constants/events/operation.go
    -> service-owned lifecycle operation type and constants (when needed)
events/<resource>/<resource>.go
    -> public domain-event contracts + stable topic/type/source identifiers
events/<resource>/impl/<resource>.go
    -> event publication/decorator implementation
events/<resource>/helpers/<resource>.go
    -> pure event envelope/payload/metadata construction

mappers/<resource>.go
    -> mapper interface
mappers/impl/<resource>.go
    -> concrete representation conversion

validators/<resource>.go
    -> validator interface
validators/impl/<resource>.go
    -> validation orchestration
validators/helpers/<resource>.go
    -> pure resource-specific validation/normalization policy

repository/<resource>.go
    -> repository interface + repository sentinel errors
repository/impl/<resource>.go
    -> database query/persistence implementation

services/<resource>.go
    -> service interface + service-owned input/result contracts
services/impl/<resource>.go
    -> use-case orchestration
services/helpers/<resource>.go
    -> pure helpers used by service orchestration

controller/<resource>.go
    -> controller interface
controller/impl/<resource>.go
    -> HTTP adapter
controller/helpers/<resource>.go
    -> request/error transport helpers

client/<resource>.go
    -> public in-process client interface
client/impl/<resource>.go
    -> adapter to provider service
```

Consumers depend on public interfaces/contracts, not concrete `impl` packages. Concrete construction belongs in the owning module composition root.

---

## 3. Layer Flow

HTTP flow:

```text
Route
  ↓
Controller interface
  ↓
Controller implementation
  ↓
Request DTO / shared contract
  ↓
Service interface
  ↓
Service implementation
  ├──► Validator interface
  │       ↓
  │    Validator implementation
  │       ↓
  │    normalized existing contract
  ↓
Repository interface
  ↓
Repository implementation
  ↓
Database
```

In-process service-to-service flow:

```text
Consuming service
  ↓
Provider client interface
  ↓
Provider client implementation
  ↓
Provider service interface
  ↓
Provider service implementation
  ↓
Provider repository
```

Domain-event publication flow:

```text
Owning service operation
  ↓
Event implementation/decorator
  ↓
Infrastructure MessagingClient
  ↓
Messaging provider strategy
  ↓
NATS JetStream
```

Do not bypass the provider service layer by importing its repository or implementation packages.

Do not import NATS/Watermill/provider mechanics from domain event packages. Domain events depend only on the public Infrastructure messaging client.

---

## 4. Entity Rule

An entity represents the persisted resource/domain record used by the owning service.

Foundation identity fields use **plain Go `string`**: `Person.GID string` and `Contact.PersonGID string`. There is no `entity.GID` type or `entity/gid.go`; services and repositories accept strings while validators enforce UUIDv7 and GORM retains PostgreSQL `uuid` columns.

For GORM-backed modules, the entity is also the GORM database struct. Do not introduce a parallel persistence model and mapper when the entity already represents the same database record.

Foundation Person is the canonical example; entity structs contain only persisted columns:

```text
entity/person.go
    -> Person passive GORM/database struct and field tags
    -> @Entity table name mapping only (generated outside source)

entity/contact.go
    -> Contact passive GORM/database struct and field tags
    -> @Entity table name mapping only (generated outside source)


services/helpers/person.go
    -> NewPerson and contact construction (independent values)
    -> OptionalString, NullableString, CopyContacts and association rules

repository/impl/person.go
    -> persists Person columns and their SearchDetails projection
    -> uses @Searchable-generated Search with validated read-only contact filters

repository/impl/contact.go
    -> persists Contact rows only; queries PersonGID foreign keys

repository/impl/transaction.go
    -> coordinates both repositories through a shared GORM transaction

services/impl/person.go
    -> coordinates both entity values through Person's lifecycle

mappers/impl/person.go
    -> joins Person and []Contact into the existing response DTO
```

Do not create:

```text
person_model.go
person_mapper.go
```

solely to copy identical fields between a persistence model and `entity.Person`.

Entities must **not** own:

```text
repository queries
search projection calculation
HTTP behavior
authorization
workflow orchestration
cross-aggregate business decisions
external I/O
```

Persistence-derived behavior that exists only to store/query database data belongs in the repository implementation.

For example, Foundation's denormalized `search_details` projection is calculated by `services/helpers/person.go` and written through `PersonRepository`, not by `entity.Person`. The `Person.IdempotencyKey` and `Person.RequestHash` Go fields continue mapping to the existing `create_idempotency_key` and `create_request_hash` PostgreSQL columns; no schema migration is needed.

---

## 5. Repository Rule

The repository boundary is persistence/query ownership.

`repository/<resource>.go` contains the interface and repository-specific errors.

`repository/impl/<resource>.go` owns:

```text
GORM/PGX database access
query composition
transactions
constraint mapping
aggregate persistence
contact/child persistence
validated search query execution
persistence-derived projection values
```

For GORM-backed repositories, prefer GORM APIs for normal CRUD/query behavior. Raw SQL is an exception for database behavior that cannot be expressed cleanly/maintainably through the ORM.

Do not move persistence calculations into entities merely to avoid repository code.

The repository public package path must not encode database technology. Consumers depend on `repository.PersonRepository`, not `repository/postgres` or a concrete GORM type.

---

## 6. Mapper Rule

Create a mapper only when there is a real representation boundary.

Examples:

```text
entity.Person -> dto/response.Person
request.GID (string) -> UUIDv7 validation -> repository string argument
entity metadata -> response error metadata
```

Do not create a mapper for identical database/entity structures.

Do not create separate resource error mapper files when the metadata belongs only to that resource.

---

## 7. Validator Rule

Validators own request/service-input validation and canonicalization.

They may:

```text
validate required/optional fields
normalize email/phone/UUID through core capabilities
validate resource search fields/operators
validate cursor values
return normalized existing DTO/service/core contracts
```

They must not:

```text
query repositories
perform authorization
call external services
own business workflows
generate IDs/read clocks
create duplicate ValidatedX wrappers without new information
```

If validation leaves the shape unchanged, return the existing request/service-input contract.

---

## 8. Service Rule

Service implementations own use-case/business orchestration.

Typical responsibilities:

```text
call validators
perform repository-backed existence/uniqueness workflows
construct/update entities
manage lifecycle transitions
call other module clients
map repository errors to domain/application errors
```

Services must not duplicate request normalization already owned by validators or generic query mechanics already owned by core capabilities.

Domain-event publication orchestration belongs under the owning service's `events/<resource>/impl` package, not mixed into the normal `services/impl` implementation. For registration, Identity owns one `core/database/transaction.Run(ctx, db, operation)` transaction; all Foundation/Identity repositories use `transaction.Database(ctx, repository.db)` and events use `transaction.AfterCommit`. Domain contracts do not expose GORM.

---

## 9. Controller Rule

Controllers are HTTP adapters.

They own:

```text
request decode
path/query extraction
calling the service interface
HTTP status selection
HTTP response writing
```

They do not own repository/business logic.

Routes depend on the controller interface, not controller implementations.

---

## 10. Client Rule

A service client is an in-process adapter to the provider's service interface.

It must:

```text
forward context.Context
reuse provider request/response contracts
call provider service methods
avoid duplicate business logic
avoid internal HTTP when both modules are in-process
```

A consuming module imports only the provider's public `client` package.

---

## 10.1 Event Rule

A service owns the semantic meaning and public contract of the events caused by changes to resources it owns.

Examples:

```text
Foundation -> Person/GID events
Identity   -> Account/Credential/Session events
```

Event contracts live under:

```text
services/<service>/events/<resource>/
```

Event publication implementations live under:

```text
services/<service>/events/<resource>/impl/
```

Pure mapping/envelope helpers live under:

```text
services/<service>/events/<resource>/helpers/
```

Stable subjects/topics must be reused across message instances. Message/event IDs provide uniqueness; random subjects must not be created for every message or sanity-test run.

Created events are emitted only when a create operation genuinely returns `created == true`; idempotent replays returning `created == false` must not emit duplicate created events.

Event payloads should contain the minimum required contract and should not copy PII/secrets by default.

All events use the generic envelope from `libs/core/events` and publish through Infrastructure's public `MessagingClient`.

See `EVENT_GUIDELINES.md` for the complete event standard.

---

## 11. Core Library Structure

`backend/libs/core` is capability-first.

Canonical shape:

```text
backend/libs/core/
├── database/
│   └── postgres/
├── errors/
│   ├── contracts/
│   ├── helpers/
│   └── errors.go
├── events/
│   └── envelope.go
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

There must not be repository-wide role buckets such as:

```text
libs/core/contracts
libs/core/helpers
libs/core/dto
libs/core/utils
libs/core/common
```

A contract/helper/mapper/validator belongs under the capability that owns it.

Service/module code imports only capability roots:

```go
coreerrors "grovehub/libs/core/errors"
coreevents "grovehub/libs/core/events"
coreid "grovehub/libs/core/id"
corelogging "grovehub/libs/core/logging"
corenormalize "grovehub/libs/core/normalize"
coresearch "grovehub/libs/core/search"
coretelemetry "grovehub/libs/core/telemetry"
```

Do not import a core capability's internal `contracts/helpers/mappers/validators` package from a service.

See `CORE_LIBRARY_GUIDELINES.md` for the full rule.

---

## 12. Shared Search Infrastructure

Search is one encapsulated core capability:

```text
backend/libs/core/search/
├── search.go
├── contracts/
├── helpers/
├── mappers/
└── validators/
```

Consumers import only `grovehub/libs/core/search`.

External request/response names:

```text
search.Request
search.Filter
search.Sort
search.Result[T]
```

Validated repository names:

```text
search.QueryFilter
search.QuerySort
search.Cursor[TId]
search.Query[TId]
search.QueryResult[T]
```

Resource validators own allowed field/operator policy and resource-specific normalization.

Repositories own safe persistence expression mapping.

The search capability owns generic operator/cursor/query mechanics.

Do not create resource-specific search wrappers that merely duplicate these generic contracts.

---

## 13. Normalization and ID Capabilities

Use:

```go
coreid "grovehub/libs/core/id"
corenormalize "grovehub/libs/core/normalize"
```

`core/id` owns generation such as `NewUUIDv7`.

`core/normalize` owns domain-neutral canonicalization such as:

```text
Email
Phone
PhoneCharacters
UUIDv7
```

Domain error mapping remains in the owning service validator/helper.

---

## 14. Foundation Person Aggregate

Foundation exposes one Person aggregate:

```text
Person
├── GID
├── FirstName
├── MiddleName
├── LastName
├── Contacts[]       # Person-owned child rows (NOT a field in entity.Person)
├── Status
├── Version
├── SearchDetails   # persisted search projection; repository maintained
├── CreatedAt
└── UpdatedAt
```

Contacts are nested API values, not separately addressable public resources. `entity.Person` contains **no** `Contacts` field. `PersonRepository` returns Person rows, `ContactRepository` returns separate Contact rows, `PersonService` owns and transactionally coordinates contact changes, and the mapper combines them into the existing DTO. No `PersonWithContacts` model or new `@Model` annotation is needed.

Search remains a Person use case; do not create `person_search.go` merely because Person supports search.

Foundation owns the Person/GID lifecycle event contract. Its current unified event is `foundation.person.changed.v1` on the stable `foundation_person_changed` topic. Its `data` contains `operation`, `old`, and `new` using complete `dto/response.Person` snapshots; supported operations are CREATED, UPDATED and DEACTIVATED. A successful deactivation publishes only DEACTIVATED, never a duplicate UPDATED event; no-op changes do not publish. Consumers can calculate a flat field-level diff with `libs/core/diff`, which never imports Foundation. See `EVENT_GUIDELINES.md` and `backend/services/foundation/events/persons/README.md` for contract and caveats.

---

## 15. SRP / OCP Rules

Prefer registries/factories when behavior is selected by stable keys such as:

```text
search operator
search field
sort field
error code
contact type
```

A new generic search operator should normally be implemented once in `core/search` and enabled by resource policy, rather than copied into every repository.

Do not force abstractions where simple linear logic is clearer. The goal is clear ownership and low duplication.

---

## 16. Hard Rules

```text
NO src/ wrapper
NO service-local launcher
NO consuming service importing another service's impl package
NO controller/repository/network I/O from validators
NO business workflow inside validators/controllers/mappers
NO persistence query logic inside entities
NO parallel persistence model/mapper when entity already is the DB model
NO raw SQL for ordinary ORM CRUD when GORM expresses it cleanly
NO resource-specific wrappers duplicating generic search contracts
NO arbitrary client field names becoming SQL expressions
NO top-level libs/core/contracts
NO top-level libs/core/helpers
NO top-level libs/core/dto
NO generic libs/core/utils/common dumping ground
NO service importing core capability internals
NO NATS/Watermill/provider imports from domain event packages
NO random event topics/subjects per message or test run
YES repository/impl owns database queries and persistence-derived projections
YES entities contain only passive persisted state, types and database field tags; service helpers own constructors, optional values and collection copies
YES one stable façade per core capability
YES resource policy stays outside generic core mechanics
YES domain services own event meaning/contracts
YES Infrastructure owns event transport mechanics
YES events are grouped under events/<resource>/{impl,helpers}
```

> **Guiding principle:** services own use cases and domain-event meaning, repositories own persistence/query behavior, entities represent persisted resource state without workflow/persistence-query logic, Infrastructure owns transport mechanics, and shared core code is encapsulated by capability behind stable root packages.

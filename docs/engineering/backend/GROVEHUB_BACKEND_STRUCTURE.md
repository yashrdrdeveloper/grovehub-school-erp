# GroveHub Current Backend Structure

> **Status:** Mandatory GroveHub architecture standard  
> **Applies to:** the current Go modular monolith and all new backend services/modules  
> **Precedence:** this file is authoritative for GroveHub backend placement and layer ownership.

Also read when relevant:

```text
docs/engineering/backend/CORE_LIBRARY_GUIDELINES.md
docs/engineering/backend/MAPPER_AND_DTO_GUIDELINES.md
docs/engineering/backend/VALIDATOR_GUIDELINES.md
docs/engineering/backend/NORMALIZATION_AND_DRY_GUIDELINES.md
docs/engineering/backend/SEARCH_API_GUIDELINES.md
```

---

## 1. Canonical Service Structure

GroveHub uses flat service roots. Do not add `src/` or service-local launchers.

```text
backend/services/<service>/
├── <service>.go
├── api/
├── config/
├── constants/
│   └── error/
├── dto/
│   ├── request/
│   └── response/
├── entity/
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

Create only folders/files that have a current responsibility. Do not create empty mapper, validator, helper, or implementation packages merely for symmetry.

---

## 2. Contract / Implementation Rule

Interfaces belong at the package boundary. Concrete implementations belong under `impl/`.

```text
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

Consumers depend on public interfaces, not concrete `impl` packages. Concrete construction belongs in the owning module composition root.

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

Do not bypass the provider service layer by importing its repository or implementation packages.

---

## 4. Entity Rule

An entity represents the persisted resource/domain record used by the owning service.

For GORM-backed modules, the entity is also the GORM database struct. Do not introduce a parallel persistence model and mapper when the entity already represents the same database record.

Foundation Person is the canonical example:

```text
entity/person.go
    -> Person GORM/database struct
    -> TableName
    -> constructor
    -> simple getters/setters/builders related to entity state

entity/contact.go
    -> Contact GORM/database struct
    -> TableName
    -> constructor/entity-state helpers
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

For example, Foundation's denormalized `search_details` projection is calculated by `repository/impl/person.go`, not by `entity.Person`.

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
request GID string -> entity.GID
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
├── Contacts[]
├── Status
├── Version
├── SearchDetails   # persisted search projection; repository maintained
├── CreatedAt
└── UpdatedAt
```

Contacts are nested values, not separately addressable public resources.

Search remains a Person use case; do not create `person_search.go` merely because Person supports search.

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
NO duplicate service interfaces in controller/client packages
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
YES repository/impl owns database queries and persistence-derived projections
YES entity owns DB/entity state and entity-level constructors/helpers
YES one stable façade per core capability
YES resource policy stays outside generic core mechanics
```

> **Guiding principle:** services own use cases, repositories own persistence/query behavior, entities represent persisted resource state without workflow/persistence-query logic, and shared core code is encapsulated by capability behind stable root packages.

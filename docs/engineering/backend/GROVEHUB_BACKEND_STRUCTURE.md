# GroveHub Current Backend Structure

> **Status:** Mandatory GroveHub architecture standard  
> **Applies to:** the current Go modular monolith and all new backend services/modules  
> **Precedence:** this file is authoritative for GroveHub backend placement and layer ownership.

Also read when relevant:

```text
docs/engineering/backend/MAPPER_AND_DTO_GUIDELINES.md
docs/engineering/backend/SEARCH_API_GUIDELINES.md
docs/engineering/backend/GORM_REPOSITORY_GUIDELINES.md
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
├── repository/
│   ├── <resource>.go
│   └── impl/
│       ├── <resource>.go
│       ├── <resource>_model.go
│       └── <resource>_mapper.go
├── services/
│   ├── <resource>.go
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
    └── unit-test/
```

Create only files/folders with a current responsibility. Persistence model/mapper files are created only for resources that actually persist data.

---

## 2. Contract / Implementation Rule

Interfaces belong at the package boundary. Concrete implementations belong under `impl/`.

### Controller

```text
controller/person.go
    -> PersonController interface
    -> PersonService dependency interface

controller/impl/person.go
    -> HTTP implementation of PersonController

controller/helpers/person.go
    -> Person-specific request/error transport helpers
```

Routes depend on `controller.PersonController`, never on `controller/impl.PersonController`.

### Repository

```text
repository/person.go
    -> PersonRepository interface
    -> persistence query/result types when genuinely repository-specific
    -> repository sentinel errors

repository/impl/person.go
    -> GORM/PostgreSQL implementation
    -> aggregate persistence, transactions, validated query execution, constraint mapping

repository/impl/person_model.go
    -> persistence-only models and GORM tags

repository/impl/person_mapper.go
    -> explicit persistence model <-> domain entity mapping
```

Do not encode the database technology into the public repository package path. `impl/` is the concrete adapter boundary.

GORM types, persistence models, table names, and PostgreSQL driver types must not leak into `repository/<resource>.go`.

### Client

```text
client/person.go
    -> public client interface only

client/impl/person.go
    -> concrete in-process client implementation
    -> service delegation
```

Consuming services depend on `client.Client` and import only the public `client` package.

A consuming service must never import another service's `client/impl` package. Concrete client construction belongs to the owning service composition root.

### Construction

Only the service root composition file imports its concrete `impl` packages:

```text
foundation.go
    *gorm.DB
       ↓
    repository/impl
       ↓
    services
       ↓
    controller/impl
       ↓
    controller interface

    services
       ↓
    client/impl
       ↓
    client interface
```

The launcher/bootstrap owns creation and shutdown of the shared PostgreSQL GORM handle. Service composition roots receive that handle and pass it only to repository implementations.

The launcher may receive the owning module's public client interface and inject it into another module. It must not construct or expose another service's client implementation directly.

---

## 3. Layer Responsibilities

```text
Route
  ↓
Controller interface
  ↓
Controller implementation
  ↓
Request DTO / shared contract
  ↓
Service
  ↓
Repository interface
  ↓
Repository implementation
  ↓
GORM
  ↓
PostgreSQL
```

For in-process service-to-service calls:

```text
Consuming service
  ↓
Provider client interface
  ↓
Provider client implementation
  ↓
Provider service
  ↓
Provider repository
```

### Entity

Entities own encapsulated state only. Stateful fields remain private; external code reads through getters and changes state through copy-returning setters.

Entities must not own business decisions, repositories, HTTP, clocks, ID generation, authorization, or workflow orchestration.

Entities must also remain persistence-agnostic: no `gorm.Model`, no GORM tags, no table/column names, and no persistence lifecycle hooks.

### DTO

DTOs contain data only. Resource-owned response metadata belongs with that resource DTO when it has no independent lifecycle.

Example:

```text
dto/response/person.go
    Person
    ContactValue
    Contacts
    PersonErrorMetadata
```

Do not create `dto/response/error.go` solely for Person metadata.

### Mapper

A resource mapper owns representation changes for that resource.

```text
mappers/person.go
    ToPersonGID
    ToPersonResponse
    ToPersonErrorMetadata
```

Do not create `mappers/error.go` when the mapper only serves Person.

Persistence mapping is a separate private concern and lives under `repository/impl/<resource>_mapper.go`.

Do not use reflection-based generic copy libraries to hide domain/persistence mapping. Persistence mapping should remain explicit and compile-time-visible.

### Service

`services/<resource>.go` owns use-case orchestration.

`services/helpers/<resource>.go` may own pure, resource-specific helper behavior such as:

```text
normalization
validation
contact construction
search field registry
search cursor field registry
resource comparison helpers
```

Helpers must not call repositories or perform workflow orchestration.

Keep service methods readable. Extract reusable resource logic instead of growing deeply nested conditionals.

### Repository implementation

`repository/impl/<resource>.go` owns:

```text
GORM query execution
aggregate-local persistence transactions
constraint/error mapping
validated search query execution
PostgreSQL-specific SQL expressions when required
```

`repository/impl/<resource>_model.go` owns persistence structs and GORM tags.

`repository/impl/<resource>_mapper.go` owns explicit domain entity <-> persistence model mapping.

Use ordinary GORM APIs for ordinary CRUD. Raw SQL is allowed only when PostgreSQL-specific behavior is materially clearer or required. All user values remain parameterized.

Generic search contracts/mechanics must not be copied into every repository. Shared search types remain in `backend/libs/core/contracts/search` and entity-neutral mechanics remain in `backend/libs/core/search` where they are genuinely reusable.

SQL migrations remain the authoritative schema definition. Runtime `AutoMigrate` is forbidden for GroveHub schema evolution.

See `GORM_REPOSITORY_GUIDELINES.md` for the full persistence standard.

### Client implementation

`client/impl/<resource>.go` owns only the concrete adapter from the public client contract to the provider service.

It must:

```text
forward context.Context unchanged
reuse provider request/response contracts
call provider service methods
avoid duplicate business logic
avoid internal HTTP when both modules are in-process
```

It must not import repository implementations or bypass the provider service layer.

---

## 4. SRP / OCP Rules

Prefer registries/factories over repeated `switch` trees when behavior is selected by a stable key such as:

```text
search field
search operator
sort field
error code
contact type
```

Example:

```text
operator -> shared factory
field -> field specification
error code -> HTTP status
sort field -> sort specification
```

Adding one supported field/operator should normally extend a registry rather than require editing multiple branches across service and repository layers.

Do not force a registry where ordinary linear control flow is clearer. The objective is lower cyclomatic complexity, not abstraction for its own sake.

---

## 5. Foundation Person Aggregate

Foundation exposes one Person aggregate:

```text
Person
├── GID
├── FirstName
├── MiddleName
├── LastName
├── Contacts
│   ├── Emails[]
│   └── Phones[]
├── Status
├── Version
├── CreatedAt
└── UpdatedAt
```

Contacts are nested values, not separately addressable public resources.

Current Foundation layout:

```text
backend/services/foundation/
├── dto/
│   ├── request/person.go
│   └── response/person.go
├── entity/
│   ├── gid.go
│   ├── contact.go
│   └── person.go
├── mappers/
│   └── person.go
├── repository/
│   ├── person.go
│   └── impl/
│       ├── person.go
│       ├── person_model.go
│       └── person_mapper.go
├── services/
│   ├── person.go
│   └── helpers/
│       └── person.go
├── controller/
│   ├── person.go
│   ├── impl/
│   │   └── person.go
│   └── helpers/
│       └── person.go
├── client/
│   ├── person.go
│   └── impl/
│       └── person.go
├── docs/
│   ├── contracts.md
│   ├── search.md
│   └── persistence.md
├── routes/http.go
└── migrations/
```

The Person entity remains free of GORM. `person_model.go` owns storage mapping and `person_mapper.go` explicitly maps between persistence and domain representations.

Do not create `person_search.go` merely because Person has a search use case. Search remains part of the Person responsibility unless it becomes an independently owned subsystem.

---

## 6. Shared Search Infrastructure

Shared request/result contracts:

```text
backend/libs/core/contracts/search/
```

Reusable search mechanics:

```text
backend/libs/core/search/
```

`libs/core/search` may own generic technical behavior such as:

```text
opaque cursor encode/decode
LIKE escaping
entity-neutral search contracts/helpers
```

It must not know Foundation tables, Person fields, GORM models, or business rules.

Foundation owns only mappings such as:

```text
firstName -> LOWER(p.first_name)
contacts.email -> person_contacts EMAIL relation
createdAt -> p.created_at
```

The owning service validates public search fields/operators. The repository implementation translates only those approved semantics to fixed GORM/SQL expressions.

This keeps shared mechanics OCP-friendly while service ownership stays explicit.

---

## 7. GORM / Database Rules

One GORM PostgreSQL handle is created by launcher bootstrap and reused by service repository implementations.

Every request-path repository operation must use:

```go
db.WithContext(ctx)
```

A single repository operation may use `db.Transaction(...)` when one aggregate spans multiple owned tables and must persist atomically.

If a future business use case requires one transaction across multiple repositories, promote the transaction boundary to an explicit service/unit-of-work contract rather than nesting hidden independent transactions.

Do not add GORM association auto-save, soft-delete behavior, hooks, or generic `BaseRepository[T]` abstractions by default.

---

## 8. Hard Rules

```text
NO src/ wrapper
NO service-local launcher/
NO repository/postgres package for current GroveHub repository implementations; use repository/impl
NO GORM imports outside libs/core/database/postgres and services/*/repository/impl
NO gorm.Model or GORM tags in domain entities
NO GORM types/models/clauses in public repository interfaces
NO runtime AutoMigrate as GroveHub schema evolution
NO concrete controller dependency in routes
NO controller implementation mixed into controller contract file
NO client implementation mixed into client contract file
NO consuming service importing another service's client/impl package
NO client implementation bypassing the provider service layer
NO resource-specific helper dumped into generic common/misc files
NO separate Person search implementation file when it is the same Person responsibility
NO business logic in controller/repository/entity/mapper/client adapter
NO duplicated entity-specific search contract wrappers when shared generic search types already express the requirement
NO arbitrary client field names becoming SQL expressions
NO reflection-based generic domain/persistence mapping
NO direct mutable entity field access outside entity package
NO separate error DTO/mapper files when metadata belongs only to one resource
NO generic BaseRepository[T] that erases domain semantics
```

> **Guiding principle:** interfaces define boundaries, `impl` packages implement them, GORM simplifies persistence only behind that boundary, resource helpers reduce local complexity, and truly generic mechanics are promoted to `libs/core` only when domain-specific knowledge remains outside the shared package.

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
    -> repository sentinel errors
    -> domain/shared query/result contracts only

repository/impl/person.go
    -> GORM/PostgreSQL implementation
    -> aggregate transaction/query/error handling

repository/impl/person_model.go
    -> persistence-only structs and GORM tags

repository/impl/person_mapper.go
    -> explicit persistence model <-> entity mapping
```

Do not encode database technology into the public repository package path. `impl/` is the concrete adapter boundary.

The public repository package must not expose `*gorm.DB`, GORM clauses, table models, driver rows, or PostgreSQL-specific types.

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

Only the service root composition file imports its concrete implementation packages:

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

The launcher/bootstrap creates the shared PostgreSQL GORM handle and injects it into service composition roots. Repositories must not open their own application-level pools.

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

Entities must not own business decisions, repositories, HTTP, clocks, ID generation, authorization, workflow orchestration, GORM tags, table names, or persistence hooks.

### DTO

DTOs contain data only. Resource-owned response metadata belongs with that resource DTO when it has no independent lifecycle.

### Mapper

Public/resource mappers own API/service representation changes.

```text
mappers/person.go
    ToPersonGID
    ToPersonResponse
    ToPersonErrorMetadata
```

Persistence mapping is separate and private:

```text
repository/impl/person_mapper.go
    entity.Person <-> personModel
```

Do not use reflection-based generic copy libraries to hide domain/persistence mapping.

### Service

`services/<resource>.go` owns use-case orchestration.

`services/helpers/<resource>.go` may own pure, resource-specific helper behavior such as normalization, validation, contact construction, search field registries, cursor registries, and resource comparisons.

Helpers must not call repositories or perform workflow orchestration.

### Repository implementation

`repository/impl` owns persistence implementation only:

```text
GORM queries
persistence models
explicit entity/model mapping
aggregate-local transactions
constraint/error translation
validated search query execution
PostgreSQL-specific expressions when required
```

GORM stays inside `repository/impl` and `libs/core/database/postgres`.

Use ordinary GORM APIs for ordinary CRUD. Raw SQL is allowed only when PostgreSQL-specific behavior is clearer or required, and all user values must remain parameterized.

SQL migrations remain the authoritative schema definition. Do not use runtime `AutoMigrate` for GroveHub schema evolution.

See `GORM_REPOSITORY_GUIDELINES.md` for the complete persistence rules.

### Client implementation

`client/impl/<resource>.go` owns only the concrete adapter from the public client contract to the provider service.

It must forward `context.Context`, reuse provider contracts, call provider service methods, avoid duplicate business logic, and avoid internal HTTP when both modules are in-process.

---

## 4. SRP / OCP Rules

Prefer registries/factories over repeated `switch` trees when behavior is selected by a stable key such as search field, search operator, sort field, error code, or contact type.

Adding one supported field/operator should normally extend an allowlist/registry rather than require editing multiple branches across service and repository layers.

Do not force a registry where ordinary linear control flow is clearer. The objective is lower complexity, not abstraction for its own sake.

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

Current Foundation persistence layout:

```text
backend/services/foundation/
├── entity/
│   ├── gid.go
│   ├── contact.go
│   └── person.go
├── repository/
│   ├── person.go
│   └── impl/
│       ├── person.go
│       ├── person_model.go
│       └── person_mapper.go
├── services/
├── controller/
├── client/
├── docs/
├── routes/
└── migrations/
```

The domain entity remains free of GORM. `person_model.go` owns database mapping, and `person_mapper.go` explicitly maps between persistence and domain representations.

Do not create `person_search.go` merely because Person has a search use case. Search remains part of the Person repository responsibility unless it becomes an independently owned subsystem.

---

## 6. Shared Search Infrastructure

Shared request/result contracts:

```text
backend/libs/core/contracts/search/
```

Reusable search contracts/mechanics:

```text
backend/libs/core/search/
```

Shared search code must remain entity-neutral. It must not know Foundation tables, Person fields, GORM models, or business rules.

Foundation owns mappings such as:

```text
firstName      -> LOWER(p.first_name)
contacts.email -> foundation.person_contacts EMAIL relation
createdAt      -> p.created_at
```

The service validates client-visible fields/operators. The repository implementation translates only approved values into fixed GORM/SQL expressions.

No client-provided identifier may become a raw SQL table/column/order expression.

---

## 7. Database and Transaction Rules

One PostgreSQL GORM handle is created during launcher bootstrap and reused by repositories.

Every request-path query uses:

```go
db.WithContext(ctx)
```

A repository may use a transaction to persist one aggregate atomically across multiple owned tables. If a future use case requires one transaction across multiple repositories, promote that boundary to an explicit service/unit-of-work abstraction.

Do not perform network calls inside a database transaction.

Do not add GORM soft delete, lifecycle hooks, association auto-save, or generic base repositories unless the domain/use case explicitly requires them.

---

## 8. Hard Rules

```text
NO src/ wrapper
NO service-local launcher/
NO repository/postgres package for service implementations; use repository/impl
NO GORM import in entity/service/controller/client/public repository packages
NO gorm.Model or gorm tags in domain entities
NO GORM types in repository interfaces
NO runtime AutoMigrate as the schema source of truth
NO concrete controller dependency in routes
NO consuming service importing another service's client/impl package
NO client implementation bypassing the provider service layer
NO resource-specific helper dumped into generic common/misc files
NO separate Person search implementation file when it is the same Person responsibility
NO business logic in controller/repository/entity/mapper/client adapter
NO reflection-based generic domain/persistence mapping
NO arbitrary client field names becoming SQL expressions
NO direct mutable entity field access outside entity package
NO generic BaseRepository[T] that erases domain semantics
```

> **Guiding principle:** interfaces define boundaries, `impl` packages implement them, GORM simplifies persistence behind the repository boundary, and domain/service contracts remain independent of the ORM.

# GroveHub Current Backend Structure

> **Status:** Mandatory GroveHub architecture standard  
> **Applies to:** the current Go modular monolith and all new backend services/modules  
> **Precedence:** this file is authoritative for GroveHub backend placement and layer ownership.

Also read when relevant:

```text
docs/engineering/backend/MAPPER_AND_DTO_GUIDELINES.md
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
├── repository/
│   ├── <resource>.go
│   └── impl/
│       └── <resource>.go
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
├── migrations/
└── tests/
    └── unit-test/
```

Create only files/folders with a current responsibility.

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
    -> persistence query/result types
    -> repository sentinel errors

repository/impl/person.go
    -> PostgreSQL implementation
```

Do not encode the database technology into the public repository package path. `impl/` is the concrete adapter boundary.

### Construction

Only the service root composition file imports concrete `impl` packages:

```text
foundation.go
    repository/impl
       ↓
    services
       ↓
    controller/impl
       ↓
    controller interface
       ↓
    routes
```

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
PostgreSQL
```

### Entity

Entities own encapsulated state only. Stateful fields remain private; external code reads through getters and changes state through copy-returning setters.

Entities must not own business decisions, repositories, HTTP, clocks, ID generation, authorization, or workflow orchestration.

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
SQL
transactions
row scanning
constraint mapping
aggregate persistence
validated search query execution
```

Generic SQL operator behavior must not be copied into every repository. Reusable search SQL mechanics belong in `backend/libs/core/search`.

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
│       └── person.go
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
├── routes/http.go
└── migrations/
```

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
parameter placeholder builder
operator factories
scalar/relation filter builders
LIKE escaping
cursor comparison construction
opaque cursor encode/decode
```

It must not know Foundation tables, Person fields, or business rules.

Foundation owns only mappings such as:

```text
firstName -> LOWER(p.first_name)
contacts.email -> person_contacts EMAIL relation
createdAt -> p.created_at
```

This keeps shared mechanics OCP-friendly while service ownership stays explicit.

---

## 7. Hard Rules

```text
NO src/ wrapper
NO service-local launcher/
NO repository/postgres package for current GroveHub repository implementations; use repository/impl
NO concrete controller dependency in routes
NO controller implementation mixed into controller contract file
NO resource-specific helper dumped into generic common/misc files
NO separate Person search implementation file when it is the same Person responsibility
NO business logic in controller/repository/entity/mapper
NO duplicated generic search SQL operator switches in services
NO duplicated generic search SQL operator switches in repositories
NO arbitrary client field names becoming SQL expressions
NO direct mutable entity field access outside entity package
NO separate error DTO/mapper files when metadata belongs only to one resource
```

> **Guiding principle:** interfaces define boundaries, `impl` packages implement them, resource helpers reduce local complexity, and truly generic mechanics are promoted to `libs/core` only when the domain-specific knowledge remains outside the shared package.

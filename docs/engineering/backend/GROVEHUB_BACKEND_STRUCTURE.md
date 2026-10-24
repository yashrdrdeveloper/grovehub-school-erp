# GroveHub Current Backend Structure

> **Status:** Mandatory GroveHub architecture standard  
> **Applies to:** the current Go modular monolith and all new backend services/modules  
> **Precedence:** this file is authoritative for GroveHub backend placement and layer ownership.

Also read when relevant:

```text
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
    └── unit-test/
```

Create only files/folders with a current responsibility.

Do not create an empty `mappers/` package for a resource that has no real representation conversion yet.

Do not create an empty `validators/` package when there is no request/service-input validation responsibility yet.

---

## 2. Contract / Implementation Rule

Interfaces belong at the package boundary. Concrete implementations belong under `impl/`.

### Mapper

```text
mappers/person.go
    -> PersonMapper interface only

mappers/impl/person.go
    -> concrete PersonMapper implementation
```

Consumers depend on `mappers.PersonMapper`, not on the concrete implementation. Concrete mapper construction belongs to the owning service composition root.

### Validator

```text
validators/person.go
    -> PersonValidator interface
    -> reuses existing request/service-input/core contracts

validators/impl/person.go
    -> concrete request/service-input validation

validators/helpers/person.go
    -> pure Person-specific validation/normalization helpers
```

Consumers depend on `validators.PersonValidator`, never directly on `validators/impl.PersonValidator`.

Validator implementations must not query repositories or own business workflow decisions. Their job is payload validation and canonicalization.

Validation does not justify a parallel data model. If validation/canonicalization leaves the shape unchanged, validators return the existing request DTO or service-input contract. If a generic internal contract already exists, such as `coresearch.Query[TId]` or `coresearch.CursorPayload`, validators return that contract directly.

Do not create `ValidatedCreateX`, `ValidatedUpdateX`, or `ValidatedSearchX` wrappers merely to mark data as validated.

### Service

```text
services/person.go
    -> PersonService interface only
    -> service-owned input/result contracts when they are not transport DTOs

services/impl/person.go
    -> concrete use-case orchestration

services/helpers/person.go
    -> pure helpers actually used by service orchestration
```

Controller and client adapters depend on `services.PersonService`. Do not redeclare duplicate service interfaces in controller or client packages.

Validation/normalization helpers used only by validators belong under `validators/helpers`, not under `services/helpers`.

### Controller

```text
controller/person.go
    -> PersonController interface only

controller/impl/person.go
    -> HTTP implementation of PersonController
    -> depends on services.PersonService

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
    -> PostgreSQL implementation
```

Do not encode the database technology into the public repository package path. `impl/` is the concrete adapter boundary.

### Client

```text
client/person.go
    -> public client interface only

client/impl/person.go
    -> concrete in-process client implementation
    -> depends on services.PersonService
    -> service delegation only
```

Consuming services depend on `client.Client` and import only the public `client` package.

A consuming service must never import another service's `client/impl` package. Concrete client construction belongs to the owning service composition root.

### Construction

Only the service root composition file imports its concrete `impl` packages:

```text
foundation.go
    repository/impl
         │
    mappers/impl
         │
    validators/impl
         │
         ▼
    services/impl
         │
         ├──────────────► client/impl
         │                    │
         │                    ▼
         │               client interface
         │
         ▼
    controller/impl
         │
         ▼
    controller interface
```

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
Provider service interface
  ↓
Provider service implementation
  ↓
Provider repository
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

`mappers/<resource>.go` defines the mapper interface.

`mappers/impl/<resource>.go` owns real representation changes for that resource.

Example:

```text
mappers/person.go
    -> PersonMapper interface

mappers/impl/person.go
    -> ToPersonGID
    -> ToPersonResponse
    -> ToPersonErrorMetadata
```

Do not create a mapper when there is no actual representation boundary yet. Do not create `mappers/error.go` when the mapper only serves Person.

### Validator

`validators/<resource>.go` defines the validator interface and reuses existing request/service-input/generic core contracts whenever their shape already represents the validated data.

`validators/impl/<resource>.go` owns request/service-input validation orchestration.

`validators/helpers/<resource>.go` owns pure resource-specific validation helpers such as:

```text
required/optional field validation
request normalization
contact payload validation/canonicalization
search field/operator validation
search cursor validation
resource-specific domain error mapping for invalid payloads
```

Validators may reuse shared technical normalizers from `backend/libs/core/helpers/normalize`.

Validators must not query repositories, call external services, authorize users, generate IDs, read clocks, or choose business workflows.

A validator-specific result type is allowed only when validation produces genuinely new information that cannot be represented safely by an existing DTO, service input, entity value object, or generic core contract. A type whose only meaning is “the same request, but validated” is not allowed.

### Service

`services/<resource>.go` defines the service interface.

`services/impl/<resource>.go` owns use-case orchestration.

`services/helpers/<resource>.go` may own pure, resource-specific helper behavior actually used by service orchestration, such as:

```text
entity/value construction from already validated request data
entity comparison for update orchestration
next-page cursor value extraction after repository results return
service-result construction helpers
service-owned technical operations that do not perform I/O or workflow decisions
```

A concrete technical operation with one current implementation may remain a service helper instead of introducing an interface solely for abstraction. Identity password hashing is the canonical example: `services/helpers/password_hasher.go` owns hashing and verification while the password credential service owns orchestration. Introduce a separate interface only when there is a real interchangeable implementation requirement.

Payload validation, normalization, and validator-only search parsing do not belong in `services/helpers`.

Helpers must not call repositories or perform workflow orchestration.

Keep service methods readable. Validate/canonicalize through the validator boundary, then perform business/repository work in the service.

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
│   ├── person.go
│   └── impl/
│       └── person.go
├── validators/
│   ├── person.go
│   ├── impl/
│   │   └── person.go
│   └── helpers/
│       └── person.go
├── repository/
│   ├── person.go
│   └── impl/
│       └── person.go
├── services/
│   ├── person.go
│   ├── impl/
│   │   └── person.go
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
│   └── search.md
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

Foundation validators own request-level mappings and validation such as:

```text
allowed filter fields/operators
email/phone search value normalization
sort field validation
cursor value validation
```

The validator returns the existing generic search contracts directly:

```text
coresearch.Query[entity.GID]
coresearch.CursorPayload
```

Do not introduce `ValidatedPersonSearch` or another wrapper around those generic contracts.

Foundation repository code owns SQL mappings such as:

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
NO concrete mapper dependency outside the owning service composition root
NO concrete validator dependency outside the owning service composition root
NO concrete service dependency outside the owning service composition root
NO duplicate service interfaces in controller/client packages
NO concrete controller dependency in routes
NO controller implementation mixed into controller contract file
NO client implementation mixed into client contract file
NO consuming service importing another service's client/impl package
NO client implementation bypassing the provider service layer
NO resource-specific helper dumped into generic common/misc files
NO validator-only helper left under services/helpers
NO request-payload normalization duplicated inside service implementations
NO repository/network access from validators
NO authorization/business workflow in validators
NO ValidatedCreateX/ValidatedUpdateX wrappers when an existing request/service-input contract has the same shape
NO ValidatedSearchX wrapper when existing generic search contracts already represent the result
NO separate Person search implementation file when it is the same Person responsibility
NO business logic in controller/repository/entity/mapper/client adapter
NO duplicated generic search SQL operator switches in services
NO duplicated generic search SQL operator switches in repositories
NO arbitrary client field names becoming SQL expressions
NO direct mutable entity field access outside entity package
NO separate error DTO/mapper files when metadata belongs only to one resource
NO mapper package for a resource with no real representation change
NO interface introduced for a single concrete service helper without an actual substitution requirement
```

> **Guiding principle:** interfaces define boundaries, `impl` packages implement them, validators canonicalize payloads without inventing duplicate contracts, resource helpers stay with their actual consumer, and truly generic mechanics are promoted to `libs/core` only when domain-specific knowledge remains outside the shared package.

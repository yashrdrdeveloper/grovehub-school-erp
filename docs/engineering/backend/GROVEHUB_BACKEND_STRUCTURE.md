# GroveHub Current Backend Structure

> **Status:** Mandatory GroveHub architecture standard  
> **Applies to:** the current Go modular monolith and all new backend services/modules  
> **Precedence:** this file is authoritative for current GroveHub service folder structure and layer responsibilities. When older examples conflict with this file, this file wins.

For DTO/mapping work, also read:

```text
docs/engineering/backend/MAPPER_AND_DTO_GUIDELINES.md
```

---

## 1. Current Runtime

GroveHub runs as one Go platform process with the root `backend/launcher/` as the platform edge/composition boundary.

Service-local `launcher/` folders are not used.

---

## 2. Canonical Service Structure

GroveHub services use flat responsibility folders. Do not add a `src/` wrapper.

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
│   └── postgres/
├── services/
├── controller/
├── routes/
├── client/
├── migrations/
└── tests/
    └── unit-test/
```

Create only folders with a current responsibility.

Forbidden legacy wrappers:

```text
backend/services/<service>/src/
backend/services/<service>/launcher/
```

---

## 3. Mandatory Layer Responsibilities

```text
Route
  ↓
Controller
  ↓
Request DTO
  ↓
Service
  ↓
Entity + Repository Contract
  ↓
PostgreSQL Repository
  ↓
PostgreSQL

Repository result Entity
  ↓
Service
  ↓
Mapper
  ↓
Response DTO
  ↓
Controller
```

### 3.1 Entity

`entity/` owns encapsulated domain data types.

For stateful entities, prefer private fields with explicit getter methods.

Example:

```go
type Person struct {
    gid       GID
    status    PersonStatus
    version   int64
    createdAt time.Time
    updatedAt time.Time
}
```

A pure constructor may assign values supplied by the caller:

```go
entity.NewPerson(gid, status, version, createdAt, updatedAt)
```

A pure constructor must not:

```text
generate IDs
read the current clock
choose business defaults
perform validation workflows
query repositories
perform lifecycle/business decisions
```

Reads happen through getters when the entity exposes them:

```text
GetGID()
GetStatus()
GetVersion()
GetCreatedAt()
GetUpdatedAt()
```

Mutable entity state must not be changed by external direct field assignment.

For copy-style entities, setters use **value receivers** and return a new copied entity:

```go
func (person Person) SetStatus(status PersonStatus) Person {
    person.status = status
    return person
}
```

The setter only assigns a caller-supplied value. It must not decide whether the transition is allowed or what value should be used.

Business decisions remain in `services/`.

For immutable identity/history fields such as GID and creation time, do not expose setters unless a real business requirement exists.

Entities must not contain repository/database access, HTTP/gRPC behavior, ID generation, clock reads, authorization, or workflow orchestration.

### 3.2 DTO

Current DTO structure:

```text
dto/request/
dto/response/
```

DTOs contain data only.

Request DTOs represent application/transport input. Response DTOs represent intentional application/public output and error metadata shapes.

Rules:

```text
NO business logic in DTOs
NO repository/database access in DTOs
NO UUID/time generation in DTOs
NO entity serialized directly as a public response
NO repository returning response DTOs
```

### 3.3 Mappers

`mappers/` owns real representation transformations only.

Examples:

```text
request DTO -> entity/value
entity getters -> response DTO
service context -> error metadata DTO
persistence model -> entity, only when a separate persistence model is required
```

Mappers do not own entity lifecycle changes.

Mappers must not query databases/networks, generate IDs, read clocks, perform authorization, validate business policy, choose status transitions, choose version increments, or choose HTTP status codes.

When entity getters exist, mappers read through those getters instead of field access.

Builders are not the default pattern. Do not create builders for simple field assignment.

### 3.4 Repository Contract

`repository/` defines persistence operations only.

Prefer names such as:

```text
Create
FindByID
FindByGID
Update
Delete
Exists
List
```

Avoid business workflow names such as `Approve`, `Deactivate`, `Promote`, `Convert`, or `Validate`.

Repository-level sentinel errors may represent technical persistence outcomes. Services translate them into stable service errors.

### 3.5 PostgreSQL Repository

`repository/postgres/` owns SQL, query parameters, row scanning, PostgreSQL errors, and technical row-count handling.

If entity fields are private, scan database columns into local variables and construct the entity through its pure constructor.

Correct:

```text
SQL row
  ↓ scan local values
entity.NewPerson(...)
```

Do not mutate private entity state from the repository.

When writing an entity, use getters.

Do not create generic `repository/postgres/entity.go` by default. Add a resource-specific persistence model only when the stored shape genuinely differs from the domain entity.

### 3.6 Services

`services/` owns business behavior and use-case orchestration:

```text
validation rules
business defaults
ID generation
status/lifecycle decisions
version changes
timestamp decisions
authorization decisions belonging to the domain
transaction/use-case orchestration
repository coordination
cross-service client orchestration
translation of repository errors into stable service errors
```

The service chooses the values of a state transition, then uses entity copy-setters to produce a new entity value.

Example:

```go
updatedPerson := person.
    SetStatus(entity.PersonStatusDeactivated).
    SetVersion(person.GetVersion() + 1).
    SetUpdatedAt(now)
```

The original `person` value remains unchanged.

The service depends on repository contracts, not concrete PostgreSQL implementations.

### 3.7 Controller

`controller/` is a thin transport adapter:

```text
read route/body/query/header input
bind request DTO
call service
select HTTP status
write response DTO or normalized service error
```

Controllers must not access repositories/databases directly, own business rules, mutate entities, or build response DTO fields manually when a mapper exists.

### 3.8 Routes

`routes/` only connects routes to controllers/middleware.

No business logic and no database access.

---

## 4. Dependency Direction

Required:

```text
routes
  ↓
controller
  ↓
services
  ↓
entity + mappers + repository contracts
                   ↑
          repository/postgres
```

Forbidden:

```text
entity -> repository/services/controller
mapper -> repository/database/network
repository -> services/controller
controller -> repository/database
```

---

## 5. Foundation Example

```text
backend/services/foundation/
├── constants/
│   └── error/
│       └── person.go
├── dto/
│   ├── request/
│   │   └── person.go
│   └── response/
│       ├── person.go
│       └── error.go
├── entity/
│   ├── gid.go
│   └── person.go
├── mappers/
│   ├── person.go
│   └── error.go
├── repository/
│   ├── person.go
│   └── postgres/
│       └── person.go
├── services/
│   └── person.go
├── migrations/
└── tests/
```

Current Foundation `Person` rules:

```text
fields are private
reads use getters
NewPerson only assigns supplied values
SetStatus/SetVersion/SetUpdatedAt return copies
GID/CreatedAt have no setters
repository reads/writes through constructor/getters
service chooses all business transition values
mappers convert representations only
```

Foundation tests are added together after the Foundation flow is complete for the current implementation phase.

---

## 6. Error Contract

Service semantic error definitions live under:

```text
constants/error/<resource>.go
```

Business services return typed `coreerrors.Error` values. Raw PostgreSQL/driver errors must not become API contracts.

Error metadata is converted by `mappers/error.go` and attached to the typed application error.

HTTP status remains a controller decision.

Do not duplicate the shared HTTP error response envelope inside each service.

---

## 7. Shared Technical Libraries

Cross-service technical capabilities belong under:

```text
backend/libs/core/<capability>/
```

Small stateless helpers use focused packages:

```text
backend/libs/core/helpers/<helper>/
```

Examples:

```text
libs/core/errors
libs/core/dto
libs/core/http
libs/core/database/postgres
libs/core/helpers/id
libs/core/telemetry
```

---

## 8. Agent Hard Rules

```text
NO src/ wrapper in GroveHub services
NO service-local launcher/
NO exported mutable entity fields for stateful entities when encapsulation is established
NO direct entity field mutation outside the entity package
NO pointer-receiver setters for copy-style entities
NO business decisions inside entity setters/constructors
NO repository/database access in entities
NO business logic in DTOs
NO business logic in mappers
NO business logic in repositories
NO repository/database access from controllers
NO direct entity serialization as public API response
NO repository returning response DTOs
NO builder for simple field assignment
NO raw PostgreSQL/driver errors as API response contracts
NO duplicate per-service HTTP error response schema
```

If an older generic document conflicts with these rules, this GroveHub standard and `MAPPER_AND_DTO_GUIDELINES.md` win.

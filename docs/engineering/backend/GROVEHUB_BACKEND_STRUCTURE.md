# GroveHub Current Backend Structure

> **Status:** Mandatory GroveHub architecture standard  
> **Applies to:** the current Go modular monolith and all new backend services/modules  
> **Precedence:** This file is authoritative for current GroveHub service folder structure and layer responsibilities. When older generic backend/database examples conflict with this file, this file wins.

For DTO/mapping work, also read:

```text
docs/engineering/backend/MAPPER_AND_DTO_GUIDELINES.md
```

---

## 1. Current Runtime

GroveHub runs as one Go platform process with a root launcher acting as the HTTP edge/composition boundary.

```text
cmd/platform
   ↓
launcher (edge/composition)
   ↓
services
```

The root launcher is allowed and required.

Service-local `launcher/` folders are not part of the current GroveHub service structure.

---

## 2. Canonical Service Structure

GroveHub services use a flat responsibility-based structure. Do not add a `src/` wrapper.

```text
backend/services/<service>/
├── <service>.go
├── api/                         # public OpenAPI/proto contracts when required
├── config/                      # when required
├── constants/
│   └── error/                   # stable service error definitions grouped by resource
├── dto/
│   ├── request/                 # application/transport input DTOs
│   └── response/                # output DTOs + error metadata DTOs
├── entity/                      # data-only domain entities/types
├── mappers/                     # all object-to-object transformations/builders
├── repository/                  # persistence contracts
│   └── postgres/                # PostgreSQL implementation
├── services/                    # business use cases and business rules
├── controller/                  # transport binding + service invocation
├── routes/                      # route registration
├── client/                      # when a public cross-service client exists
├── migrations/                  # service-owned database migrations
└── tests/
    └── unit-test/               # added according to the current implementation phase
```

Use only folders with a current responsibility.

Do not add:

```text
backend/services/<service>/src/
backend/services/<service>/launcher/
```

unless an explicit architecture decision changes this standard first.

---

## 3. Mandatory Layer Responsibilities

GroveHub uses this flow:

```text
Route
  ↓
Controller
  ↓
Request DTO
  ↓
Service
  ↓
Mapper -> Entity/Value
  ↓
Repository interface
  ↓
PostgreSQL repository
  ↓
PostgreSQL

PostgreSQL repository
  ↓
Entity
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

`entity/` contains data only:

```text
structs
types
value aliases
entity-owned constants
```

Entities must not contain:

```text
business-rule constructors
validation workflows
lifecycle transitions
ID generation
clock/time generation
repository calls
SQL/database mapping
HTTP/gRPC behavior
business branching
```

If logic answers a business question, it belongs in `services/`.

### 3.2 DTO

Current GroveHub DTO structure is flat:

```text
dto/request/
dto/response/
```

DTOs contain data only.

Request DTOs represent application/transport input.

Response DTOs represent intentional public/application output and error metadata shapes.

Rules:

```text
NO business logic in DTOs
NO repository/database access in DTOs
NO UUID/time generation in DTOs
NO entity returned directly as a public response
NO repository returning response DTOs
```

### 3.3 Mappers

All object-to-object transformation belongs in `mappers/`.

Examples:

```text
request DTO -> entity/value
entity -> response DTO
existing entity -> updated entity construction
service context -> error metadata DTO
persistence model -> entity, only when a persistence model is actually required
```

Mappers are pure transformation code.

Mappers must not:

```text
query repositories/databases
call networks
read environment variables
generate IDs
read current time
perform authorization
validate business policy
choose status/lifecycle transitions
choose version increments
choose HTTP status codes
```

When an entity/DTO has several fields or is built in multiple use cases, use a focused builder in its mapper file.

Example:

```go
person := mappers.NewPersonBuilder().
    WithGID(gid).
    WithStatus(entity.PersonStatusActive).
    WithVersion(1).
    WithCreatedAt(now).
    WithUpdatedAt(now).
    Build()
```

The service decides `ACTIVE`, version `1`, and `now`. The builder only assigns supplied values.

For updates:

```go
person = mappers.NewPersonBuilderFrom(person).
    WithStatus(entity.PersonStatusDeactivated).
    WithVersion(person.Version + 1).
    WithUpdatedAt(now).
    Build()
```

### 3.4 Repository Contract

`repository/` defines persistence operations only.

Prefer persistence-oriented names:

```text
Create
FindByID
FindByGID
Update
Delete
Exists
List
```

Avoid business workflow names:

```text
Approve
Deactivate
Promote
Convert
Validate
CalculateReadiness
```

Repository-level sentinel errors may represent technical persistence outcomes such as a missing row. Services translate them into stable service errors.

### 3.5 PostgreSQL Repository

`repository/postgres/` owns PostgreSQL-specific behavior:

```text
SQL
row scanning
query parameters
PostgreSQL error handling
technical row-count checks
```

When the PostgreSQL row shape matches the domain entity, scan directly into the domain entity/resource fields.

Do not create a generic `repository/postgres/entity.go` by default.

Only introduce a separate persistence model when the stored shape actually differs from the domain shape. Name it for the resource, for example:

```text
repository/postgres/person_entity.go
```

Repositories may handle technical conditions such as `ErrNoRows`, constraint mapping, or `RowsAffected() == 0`. They must not contain business-condition branching.

### 3.6 Services

`services/` owns all business behavior and use-case orchestration:

```text
validation rules
business defaults
ID generation
status/lifecycle transitions
version changes
timestamp decisions
authorization decisions belonging to the domain
transaction/use-case orchestration
repository coordination
cross-service client orchestration
translation of repository errors into stable service errors
```

Current GroveHub application services may accept request DTOs and return response DTOs.

Inside a service:

```text
request DTO
   ↓ mapper
entity/value
   ↓ business decisions
repository
   ↓ entity
mapper
   ↓
response DTO
```

Services must use mappers for object conversion instead of assembling DTO/entity fields ad hoc.

The service depends on repository contracts, not concrete PostgreSQL implementations.

### 3.7 Controller

`controller/` is a thin transport adapter:

```text
read route/body/query/header transport input
bind into request DTO
call service
select HTTP status
write service response DTO or normalized service error
```

Controllers must not:

```text
query repositories directly
generate IDs
apply business validation
perform lifecycle transitions
own domain defaults
convert entities into response DTOs field-by-field
build error metadata maps manually
```

HTTP status selection belongs in the controller/transport layer.

### 3.8 Routes

`routes/` only connects routes to controllers/middleware.

Routes contain no business logic and no database access.

---

## 4. Dependency Direction

Required direction:

```text
routes
  ↓
controller
  ↓
services
  ↓
mappers + repository contracts
  ↓                ↑
entity          repository/postgres
```

Forbidden:

```text
entity -> services
entity -> repository
entity -> controller
mapper -> repository/database/network
repository -> services
repository -> controller
controller -> repository
controller -> database
```

Entities and DTOs remain passive data structures. Mappers remain pure transformations.

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
└── tests/                       # add Foundation tests after the Foundation flow is complete
```

Do not introduce `src/`.

---

## 6. Error Contract

Service-specific semantic error definitions live under:

```text
constants/error/<resource>.go
```

A definition owns stable `Code` + `Message`.

Business services return typed `coreerrors.Error` values. Repository/driver errors may be retained as wrapped causes for diagnostics, but raw PostgreSQL errors must not become API contracts.

Optional service error metadata is constructed by a mapper and attached to `coreerrors.Error`.

Example:

```text
GID/entity/service context
      ↓
mappers.ToPersonErrorMetadata(...)
      ↓
coreerrors.Error.Metadata()
      ↓
controller
      ↓
shared HTTP error response envelope
```

Do not build arbitrary error metadata maps in controllers/services.

HTTP status remains a controller decision.

The canonical shared HTTP error body remains:

```json
{
  "success": false,
  "errorCode": "FOUNDATION_PERSON_NOT_FOUND",
  "errorMessage": "Person was not found",
  "httpErrorCode": 404,
  "errorMetaData": {
    "gid": "..."
  }
}
```

Do not duplicate the error response envelope inside each service.

---

## 7. Shared Technical Libraries

Cross-service technical capabilities belong under:

```text
backend/libs/core/<capability>/
```

Small stateless reusable helpers use focused helper packages:

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

`helpers/` must not become a dumping ground.

UUID generation remains in:

```text
backend/libs/core/helpers/id/
```

Shared PostgreSQL pool construction remains in:

```text
backend/libs/core/database/postgres/
```

---

## 8. Unit Tests

Backend unit tests live inside their owner.

For the current Foundation implementation phase, complete the Foundation flow first and add Foundation tests together afterward.

Pure data-only entities/DTOs do not need meaningless tests for fields alone.

---

## 9. Agent Rule

Before creating backend code, classify it by responsibility:

```text
data/type/constant                 -> entity/
request/application input          -> dto/request/
response/application output        -> dto/response/
object conversion/builder          -> mappers/<resource>.go
error metadata conversion          -> mappers/error.go
service semantic error             -> constants/error/<resource>.go
persistence contract               -> repository/
PostgreSQL persistence             -> repository/postgres/<resource>.go
business rule/use case             -> services/
HTTP/gRPC transport                -> controller/
route registration                 -> routes/
service-owned schema               -> migrations/
shared technical capability        -> libs/core/<capability>/
small stateless technical helper   -> libs/core/helpers/<helper>/
root platform composition          -> backend/launcher/
```

Hard rules:

```text
NO src/ wrapper in GroveHub services
NO business logic in entity/
NO business logic in DTOs
NO business logic in mappers/
NO business logic in repository/
NO repository/database/network calls from mappers/
NO repository/database access from controller/
NO direct entity serialization as public API response
NO repository returning response DTOs
NO ad-hoc field-by-field entity/DTO conversion outside mappers/
NO ad-hoc error metadata maps outside mappers/
NO ID/time/default generation inside builders
NO business logic in routes/
NO service-local launcher/
NO generic repository/postgres/entity.go unless a real persistence-shape divergence requires a resource-specific model
NO raw PostgreSQL/driver errors as API response contracts
NO duplicate per-service HTTP error response schema
NO generic files directly inside libs/core/helpers/
```

If an older generic document shows a conflicting structure, this GroveHub standard and `MAPPER_AND_DTO_GUIDELINES.md` win.

# GroveHub Current Backend Structure

> **Status:** Mandatory GroveHub architecture standard  
> **Applies to:** the current Go modular monolith and all new backend services/modules  
> **Precedence:** This file is authoritative for current GroveHub service folder structure **and layer responsibilities**. When older generic backend/database examples conflict with this file, this file wins.

---

## 1. Current Runtime

GroveHub currently runs as one Go platform process with a root launcher acting as the HTTP edge/composition boundary.

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

GroveHub services use a flat responsibility-based structure. Do **not** add a `src/` wrapper.

```text
backend/services/<service>/
├── <service>.go
├── api/                         # public OpenAPI/proto contracts when required
├── config/                      # when required
├── constants/                   # service-wide technical/semantic constants when required
│   └── error/                   # stable service error definitions grouped by resource
├── entity/                      # data-only domain entities/types
├── repository/                  # persistence contracts
│   └── postgres/                # PostgreSQL implementation
├── services/                    # business use cases and business rules
├── controller/                  # transport/controller adapters
├── routes/                      # route registration
├── dto/
│   └── vN/
│       ├── request/
│       └── response/
├── client/                      # when a public cross-service client exists
├── migrations/                  # service-owned database migrations
└── tests/
    └── unit-test/               # mirrors production responsibility
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

GroveHub follows this responsibility split:

```text
Route
  ↓
Controller
  ↓
Service
  ↓
Repository interface
  ↓
PostgreSQL repository
  ↓
PostgreSQL
```

`entity/` is a shared data model used by the service and repository layers. It does not call or depend on those layers.

### 3.1 Entity

`entity/` contains **data only**:

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

Example:

```go
type GID string

type PersonStatus string

const (
    PersonStatusActive      PersonStatus = "ACTIVE"
    PersonStatusDeactivated PersonStatus = "DEACTIVATED"
)

type Person struct {
    GID       GID
    Status    PersonStatus
    Version   int64
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

If logic answers a business question, it belongs in `services/`, not `entity/`.

### 3.2 Repository Contract

`repository/` defines persistence operations only.

Prefer persistence-oriented names such as:

```text
Create
FindByID
FindByGID
Update
Delete
Exists
List
```

Avoid business workflow names such as:

```text
Approve
Deactivate
Promote
Convert
Validate
CalculateReadiness
```

Repository interfaces must not contain business decisions or workflow orchestration.

Repository-level sentinel errors may represent persistence outcomes such as a missing row. Those errors are translated by the service layer into stable service errors.

### 3.3 PostgreSQL Repository

`repository/postgres/` owns PostgreSQL-specific behavior:

```text
SQL
row scanning
query parameters
PostgreSQL error handling
technical row-count checks
mapping into domain entities when shapes match
```

Do **not** create a generic `repository/postgres/entity.go` by default.

When the PostgreSQL row shape matches the domain entity, keep the persistence code in the resource repository file, for example:

```text
repository/postgres/person.go
```

Only introduce a separate persistence model when the stored shape actually diverges from the domain shape. When needed, name it for the resource, for example:

```text
repository/postgres/person_entity.go
```

Repositories may contain technical condition handling required by the database driver, for example `ErrNoRows`, constraint mapping, or `RowsAffected() == 0`. They must not contain business-condition branching.

### 3.4 Services

`services/` owns all business behavior and use-case orchestration, including:

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

Example creation flow:

```text
service.CreatePerson
      ↓
validate business input
      ↓
generate UUIDv7 through libs/core/helpers/id
      ↓
construct entity.Person data
      ↓
repository.Create
```

The service depends on repository contracts, not concrete PostgreSQL implementations.

### 3.5 Controller

`controller/` is a thin transport adapter:

```text
read request
map request DTO/input
call service
map service result/error
write response
```

Controllers must not:

```text
query repositories directly
generate IDs
apply business validation
perform lifecycle transitions
own domain defaults
```

HTTP status selection belongs in the controller/transport layer, not in the service.

### 3.6 Routes

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
repository contract
  ↑
repository/postgres
```

The following are forbidden:

```text
entity -> services
entity -> repository
entity -> controller
repository -> services
repository -> controller
controller -> repository
controller -> database
```

Domain entities remain passive data structures.

---

## 5. Identity Example

Current Identity shape may remain smaller until more responsibilities exist:

```text
backend/services/identity/
├── identity.go
├── api/
│   └── openapi/v1/openapi.yaml
├── constants/
│   └── errorcode/               # legacy/current Identity path; new service code uses constants/error/
├── controller/
├── routes/
├── dto/
│   └── v1/
│       └── response/
└── tests/
    └── unit-test/
```

Do not create empty folders merely to match the full blueprint.

---

## 6. Foundation Example

Foundation uses the flat structure:

```text
backend/services/foundation/
├── constants/
│   └── error/
│       └── person.go
├── entity/
│   ├── gid.go
│   └── person.go
├── repository/
│   ├── person.go
│   └── postgres/
│       └── person.go
├── services/
│   └── person.go
├── migrations/
└── tests/                       # add Foundation tests after the Foundation flow is complete
```

As Foundation grows, add responsibilities directly under the service root:

```text
controller/
routes/
dto/
```

Do not introduce `src/`.

---

## 7. Root Launcher

The root launcher owns platform-wide concerns such as:

```text
HTTP server lifecycle
edge routing/fallback
cross-cutting middleware
OpenTelemetry server instrumentation
service module registration
platform bootstrap
```

It must not absorb service business logic.

The launcher is the current edge/API-gateway boundary and starts the top-level HTTP trace span.

---

## 8. Shared Technical Libraries and Helpers

Cross-service technical capabilities belong under:

```text
backend/libs/core/<capability>/
```

Small stateless reusable helpers belong under a named helper package:

```text
backend/libs/core/helpers/<helper>/
```

`helpers/` must never become a dumping ground. Every helper gets a focused subpackage named for one technical responsibility.

Examples:

```text
libs/core/errors
libs/core/dto
libs/core/http
libs/core/database/postgres
libs/core/helpers/id
libs/core/telemetry
```

Shared business/domain concepts do not belong in `libs/core` merely because more than one service might eventually use them.

UUID generation is a stateless technical helper and belongs in:

```text
backend/libs/core/helpers/id/
```

The GroveHub helper should expose a small stable wrapper API and delegate UUID algorithm details to the selected UUID library. Application code must not reimplement UUID bit layout, timestamp encoding, randomness layout, or text formatting.

Business services may use the ID helper to generate IDs before persistence.

Shared PostgreSQL pool construction belongs in:

```text
backend/libs/core/database/postgres/
```

It receives configuration from composition/bootstrap code; it does not read service-specific environment variables or own business logic.

---

## 9. Error Contract

Service-specific semantic error definitions live with the owning service and are grouped by resource, for example:

```text
backend/services/foundation/constants/error/person.go
```

A semantic error definition owns stable:

```text
Code
Message
```

Business services return typed `coreerrors.Error` values using those definitions. Repository/driver errors may be retained as wrapped causes for diagnostics, but raw PostgreSQL errors must not become API contracts.

HTTP status remains a transport/controller decision.

Optional metadata is supplied at the request/transport boundary when needed.

The canonical HTTP error body is the existing shared schema:

```json
{
  "success": false,
  "errorCode": "FOUNDATION_PERSON_NOT_FOUND",
  "errorMessage": "Person was not found",
  "httpErrorCode": 404,
  "errorMetaData": {}
}
```

Flow:

```text
PostgreSQL/Repository error
        ↓
Service error definition + optional wrapped cause
        ↓
Controller chooses HTTP status/metadata
        ↓
libs/core/http/response writes canonical error JSON
```

Do not duplicate this response schema inside each service.

---

## 10. Trace Context

The root launcher starts/continues the edge trace.

Service route boundaries add child spans while preserving the same trace ID.

All nested calls must pass `context.Context`.

See:

```text
docs/engineering/observability/TRACE_CONTEXT_GUIDELINES.md
```

---

## 11. Unit Tests

Backend unit tests live inside their owner, not under `backend/testing`.

Tests mirror the flat production responsibility:

```text
entity/              -> tests/unit-test/entity/          # only when entity behavior exists
services/            -> tests/unit-test/services/
controller/          -> tests/unit-test/controller/
repository/          -> tests/unit-test/repository/
repository/postgres/ -> tests/unit-test/repository/postgres/
helpers/id/          -> tests/unit-test/helpers/id/
```

Pure data-only entities do not require meaningless tests merely to test struct fields.

For the current Foundation implementation phase, complete the Foundation flow first and add its tests together afterward, as explicitly directed for this phase.

See:

```text
docs/engineering/testing/BACKEND_UNIT_TEST_STRUCTURE.md
```

---

## 12. Agent Rule

Before creating backend code, classify it by responsibility:

```text
data/type/constant                 -> entity/
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
owner-local unit tests             -> tests/unit-test/<mirrored path>/
```

Hard rules:

```text
NO src/ wrapper in GroveHub services
NO business logic in entity/
NO business logic in repository/
NO repository/database access from controller/
NO business logic in routes/
NO service-local launcher/
NO generic repository/postgres/entity.go unless a real persistence-shape divergence requires a resource-specific model
NO raw PostgreSQL/driver errors as API response contracts
NO duplicate per-service HTTP error response schema
NO generic files directly inside libs/core/helpers/; use a focused helper subpackage
```

If an older generic document shows a conflicting structure or places business behavior in entities/repositories, this GroveHub standard wins.

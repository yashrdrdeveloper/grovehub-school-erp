# GroveHub Current Backend Structure

> **Status:** Mandatory GroveHub architecture standard  
> **Applies to:** the current Go modular monolith and all new backend services/modules  
> **Precedence:** this file is authoritative for current GroveHub service folder structure and layer responsibilities. When older examples conflict with this file, this file wins.

Also read when relevant:

```text
docs/engineering/backend/MAPPER_AND_DTO_GUIDELINES.md
docs/engineering/backend/SEARCH_API_GUIDELINES.md
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
Request DTO / shared contract
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
Response DTO / shared result contract
  ↓
Controller
```

### 3.1 Entity

`entity/` owns encapsulated domain data types.

For stateful entities, prefer private fields with explicit getter methods.

A pure constructor may assign values supplied by the caller, but it must not generate IDs, read the current clock, choose business defaults, query repositories, perform validation workflows, or make lifecycle/business decisions.

Mutable entity state must not be changed by external direct field assignment.

For copy-style entities, setters use value receivers and return a new copied entity:

```go
func (person Person) SetStatus(status PersonStatus) Person {
    person.status = status
    return person
}
```

The setter only assigns the caller-supplied value. Business decisions remain in `services/`.

Collection getters and setters must preserve encapsulation. A getter such as `GetContacts()` returns a copy rather than exposing the entity's internal slice for mutation.

For immutable identity/history fields such as GID and creation time, do not expose setters unless a real business requirement exists.

Entities must not contain repository/database access, HTTP/gRPC behavior, ID generation, clock reads, authorization, or workflow orchestration.

### 3.2 DTO and Shared Contract

Current service DTO structure:

```text
dto/request/
dto/response/
```

DTOs contain transport/application data only.

Cross-service generic contracts may live under:

```text
backend/libs/core/contracts/<contract>/
```

Example:

```text
backend/libs/core/contracts/search/
```

Shared contracts define reusable boundary shapes only. They do not own service-specific field names, SQL, tables, business validation, or repository behavior.

Rules:

```text
NO business logic in DTOs/contracts
NO repository/database access in DTOs/contracts
NO UUID/time generation in DTOs/contracts
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

`repository/` defines persistence operations and persistence query contracts only.

Prefer names such as:

```text
Create
FindByID
FindByGID
FindGIDsByContacts
Search
Update
Delete
Exists
List
```

Avoid business workflow names such as `Approve`, `Deactivate`, `Promote`, `Convert`, or `Validate`.

Repository-level sentinel errors may represent technical persistence outcomes. Services translate them into stable service errors.

A repository may define a service-specific validated search query type. It must not accept arbitrary database field names from the public request and concatenate them into SQL.

### 3.5 PostgreSQL Repository

`repository/postgres/` owns SQL, query parameters, row scanning, PostgreSQL errors, technical row-count handling, transactions required to persist one aggregate consistently, and service-owned query/index behavior.

If entity fields are private, scan database columns into local variables and construct the entity through its pure constructor.

Correct:

```text
SQL row
  ↓ scan local values
entity.NewPerson(...)
```

When writing an entity, use getters.

Dynamic search/filter/sort SQL must use a server-owned allowlist and parameterized values. Raw client field strings must never become SQL identifiers or expressions.

Do not create generic `repository/postgres/entity.go` by default. Add a persistence model only when the stored shape genuinely differs from the domain entity.

### 3.6 Services

`services/` owns business behavior and use-case orchestration:

```text
validation rules
business defaults
ID generation
normalization
status/lifecycle decisions
version changes
timestamp decisions
authorization decisions belonging to the domain
transaction/use-case orchestration
repository coordination
cross-service client orchestration
translation of repository errors into stable service errors
search field/operator validation
cursor validation/encoding semantics
```

The service chooses values for a state transition, then uses entity copy-setters to produce a new entity value.

The service depends on repository contracts, not concrete PostgreSQL implementations.

For search, the service converts the shared contract into a validated service-specific repository query before persistence code sees it.

### 3.7 Controller

`controller/` is a thin transport adapter:

```text
read route/body/query/header input
bind request DTO/shared contract
call service
select HTTP status
write response DTO/shared result or normalized service error
```

Controllers must not access repositories/databases directly, own business rules, mutate entities, or build response DTO fields manually when a mapper exists.

Resource-specific controller helpers stay private in the resource controller file. Do not create `controller/error.go`, `controller/helpers.go`, or `controller/common.go` for behavior used by only one resource. Extract shared controller infrastructure only after multiple real controllers share the concern.

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

Shared technical/application contracts may be imported by service layers without reversing ownership:

```text
libs/core/contracts/search
        ↓
controller + services + repository query types as needed
```

Forbidden:

```text
entity -> repository/services/controller
mapper -> repository/database/network
repository -> services/controller
controller -> repository/database
shared contract -> service-specific SQL/business rules
```

---

## 5. Foundation Person Aggregate

Foundation exposes Person as one aggregate.

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

Contacts are nested values, not separately addressable API resources.

Public rules:

```text
NO public contactId
NO /foundation/persons/{gid}/contacts endpoint
GET Person returns names + all contacts
PATCH Person updates names and/or replaces the complete contact set when contacts is supplied
Deactivate preserves the aggregate and contact history
isVerified is Foundation-controlled, not client-controlled
```

Persistence may normalize contacts into `foundation.person_contacts` for indexing, uniqueness, and query performance. The storage table does not create a separate public Contact resource.

Current Foundation layout includes:

```text
backend/services/foundation/
├── api/openapi/v1/openapi.yaml
├── constants/error/person.go
├── dto/
│   ├── request/person.go
│   └── response/
│       ├── person.go
│       └── error.go
├── entity/
│   ├── gid.go
│   ├── contact.go
│   └── person.go
├── mappers/
│   ├── person.go
│   └── error.go
├── repository/
│   ├── person.go
│   └── postgres/
│       ├── person.go
│       └── person_search.go
├── services/
│   ├── person.go
│   └── person_search.go
├── controller/
│   └── person.go
├── routes/http.go
├── migrations/
└── tests/
```

Current Foundation Person rules:

```text
fields are private
reads use getters
collection getters return copies
NewPerson/NewContact only assign supplied values
copy setters mutate Person state without deciding business policy
repository reads/writes through constructors/getters
service chooses business transitions, normalization, verification defaults, and version changes
mappers convert representations only
```

Foundation tests are added together after the Foundation flow is complete for the current implementation phase.

---

## 6. Foundation Contact Uniqueness

Foundation owns canonical contact values for the human.

Exact normalized email/phone values are structured persistence invariants:

```text
(type, normalized_value) is globally unique inside Foundation
```

Create behavior:

```text
normalize submitted email/phone
        ↓
exact indexed lookup
        ↓
one existing GID -> return existing Person, do not create a new GID
multiple GIDs -> stable conflict error
no GID -> create new Person
```

The database unique index remains the final race-safety guard.

Do not use `search_details` or generic `/search` for duplicate prevention.

---

## 7. Search API

All GroveHub search work must follow:

```text
docs/engineering/backend/SEARCH_API_GUIDELINES.md
```

Shared contract:

```text
backend/libs/core/contracts/search/
```

Foundation exposes:

```text
POST /foundation/persons/search
```

Foundation uses:

```text
structured exact contact indexes -> uniqueness / duplicate lookup
persons.search_details + indexed partial matching -> human-friendly searchText
service-owned filter/sort allowlists -> safe structured search
opaque keyset cursor -> pagination
```

`search_details` is derived data and must be refreshed transactionally whenever names or contacts change.

---

## 8. Error Contract

Service semantic error definitions live under:

```text
constants/error/<resource>.go
```

Business services return typed `coreerrors.Error` values. Raw PostgreSQL/driver errors must not become API contracts.

Error metadata is converted by `mappers/error.go` and attached to typed application errors when needed.

HTTP status remains a controller decision.

Do not duplicate the shared HTTP error response envelope inside each service.

---

## 9. Shared Technical Libraries

Cross-service technical capabilities belong under:

```text
backend/libs/core/<capability>/
```

Reusable boundary contracts belong under:

```text
backend/libs/core/contracts/<contract>/
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
libs/core/contracts/search
libs/core/helpers/id
libs/core/telemetry
```

---

## 10. Agent Hard Rules

```text
NO src/ wrapper in GroveHub services
NO service-local launcher/
NO exported mutable entity fields for stateful entities when encapsulation is established
NO direct entity field mutation outside the entity package
NO leaking mutable entity collection slices
NO pointer-receiver setters for copy-style entities
NO business decisions inside entity setters/constructors
NO repository/database access in entities
NO business logic in DTOs/shared contracts
NO business logic in mappers
NO business logic in repositories
NO repository/database access from controllers
NO direct entity serialization as public API response
NO repository returning response DTOs
NO builder for simple field assignment
NO raw PostgreSQL/driver errors as API response contracts
NO duplicate per-service HTTP error response schema
NO raw client field/sort strings concatenated into SQL
NO offset/page-number pagination in the shared search contract
NO search projection as source of truth for uniqueness
NO separate public Contact resource for current Foundation Person aggregate
```

If an older generic document conflicts with these rules, this GroveHub standard, `MAPPER_AND_DTO_GUIDELINES.md`, and `SEARCH_API_GUIDELINES.md` win.

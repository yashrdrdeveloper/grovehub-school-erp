# GroveHub Search API Guidelines

> **Status:** Mandatory GroveHub search standard  
> **Applies to:** shared search contracts, filters, sorting, cursor pagination, validators, repository query construction, and search projection data.

Also read:

```text
docs/engineering/backend/CORE_LIBRARY_GUIDELINES.md
```

---

## 1. One Public Search Package

All consumers import:

```go
import coresearch "grovehub/libs/core/search"
```

Do not import `core/search/contracts`, `core/search/helpers`, `core/search/mappers`, or `core/search/validators` from services/modules. Those are implementation details of the search capability.

The internal layout is:

```text
backend/libs/core/search/
├── search.go
├── contracts/
│   └── search.go
├── helpers/
│   └── query.go
├── mappers/
│   └── cursor.go
└── validators/
    └── query.go
```

`search/search.go` is the supported façade.

---

## 2. External Search Contract

The root search package exposes the transport-safe request/response contract:

```go
search.Request
search.Filter
search.Sort
search.PageRequest
search.PageResponse
search.Result[T]
```

Canonical request:

```json
{
  "searchText": "arun",
  "filters": [
    {"field": "status", "operator": "IN", "values": ["ACTIVE"]}
  ],
  "sort": [
    {"field": "createdAt", "direction": "DESC"}
  ],
  "pageInfo": {
    "rows": 25,
    "cursor": "opaque-cursor"
  }
}
```

Cursor pagination is mandatory. Do not add offset/page-number semantics to this shared contract.

A raw `search.Request` must not be passed directly to a repository.

---

## 3. Validated Repository Search Model

After resource validation/normalization, use:

```go
search.QueryFilter
search.QuerySort
search.Cursor[TId]
search.Query[TId]
search.QueryResult[T]
search.CursorPayload
```

The `Query` prefix intentionally distinguishes validated repository input/output from HTTP request/response types.

Example:

```text
search.Request
      ↓
PersonValidator
      ↓
search.Query[entity.GID]
      ↓
PersonRepository
      ↓
search.QueryResult[entity.Person]
      ↓
PersonService
      ↓
search.Result[response.Person]
```

Do not create resource wrappers such as:

```text
PersonSearchFilter
PersonSearchSort
PersonSearchCursor
PersonSearchQuery
PersonSearchResult
```

when they only duplicate the generic search capability.

---

## 4. Search Capability Responsibilities

### `search/contracts`

Owns shared search types and configuration contracts.

### `search/helpers`

Owns generic query mechanics:

```text
parameter argument collection
operator factories
scalar predicates
relation/EXISTS predicates
LIKE escaping
keyset cursor predicates
```

### `search/mappers`

Owns generic representation conversion such as opaque cursor encode/decode.

### `search/validators`

Owns domain-neutral query/config validation such as valid limit, sort mapping, and direction.

### `search/search.go`

Is the public façade and orchestration entry point. `Apply` composes validation and query-building behavior and is the only package path repositories need.

The shared search capability must never know Foundation tables, Person fields, or service-specific business policy.

---

## 5. Resource Validator Owns Public Search Policy

Each resource validator owns:

```text
allowed public fields
allowed operators per field
value normalization
sort-field policy
cursor value validation
resource-specific cursor ID validation
```

Example:

```text
firstName
    operators: EQUALS, IN, CONTAINS, STARTS_WITH
    normalizer: lower-case trimmed string

createdAt
    operators: EQUALS, GT, GTE, LT, LTE
    normalizer: RFC3339 UTC

contacts.email
    operators: EQUALS, IN, EXISTS, CONTAINS
    normalizer: normalized email
```

The flow is:

```text
search.Request
      ↓
validators/<resource>.go
      ↓
validators/impl/<resource>.go
      ↓
validators/helpers/<resource>.go
      ↓
search.Query[TId]
```

Unknown fields/operators must be rejected before repository execution.

A service must not repeat this validation after receiving a validated query.

---

## 6. Repository Owns Safe Field Mapping

Repositories receive only validated `search.Query[TId]` values.

The repository implementation owns mappings from public field names to safe persistence expressions:

```text
firstName      -> LOWER(p.first_name)
createdAt      -> p.created_at ::timestamptz
contacts.email -> person_contacts relation + EMAIL fixed predicate
```

Raw client field names must never be concatenated directly into SQL.

The repository calls `search.Apply` using a capability `search.Config`.

Repository search methods return:

```go
search.QueryResult[Entity]
```

not transport-facing `search.Result`.

---

## 7. Generic Operators

Generic operator behavior is implemented once inside the search capability:

```text
EQUALS / NOT_EQUALS
IN / NIN
EXISTS / NOT_EXISTS
GT / GTE / LT / LTE
CONTAINS / STARTS_WITH
```

Do not copy generic operator switch trees into every validator or repository.

Use:

```text
operator -> shared search helper/factory
field -> resource validator/repository specification
```

A new generic operator is added to `core/search` once and explicitly enabled in each resource policy that supports it.

---

## 8. Cursor Rules

Cursor values are opaque to clients.

The search capability owns:

```text
CursorPayload
encoding/decoding
keyset predicate construction
```

The resource validator owns validation of:

```text
requested sort field
requested direction
sort value type
resource-specific cursor ID type
```

The repository owns the actual stable cursor expression.

Example:

```text
Foundation Person:
    CursorId         = Person GID
    CursorExpression = p.gid
```

The generic search capability must not assume a UUID/GID cursor ID or hardcode a cursor-ID cast.

Repositories use keyset pagination, not `OFFSET`.

---

## 9. Search Projection vs Source of Truth

`search_details` is denormalized projection data for free-text search. It is not authoritative and must not be used for uniqueness/deduplication.

Foundation example:

```text
person_contacts unique/indexed normalized values -> uniqueness/source of truth
persons.search_details + trigram index            -> free-text search projection
```

When a Person or its contacts change, the repository rebuilds and persists `search_details` inside the same transaction. This projection calculation belongs to persistence/repository code, not the entity.

Do not perform a second SQL aggregation query when the repository already has the complete aggregate in memory.

---

## 10. Hard Rules

```text
NO imports of core/search internals from services/modules
NO core/contracts/search package
NO raw search.Request passed to repositories
NO search request validation duplicated in services
NO resource wrappers duplicating Query/QueryResult
NO raw client field -> SQL identifier mapping
NO unparameterized filter values
NO copied generic operator switch tree per repository
NO copied cursor codec per service
NO offset/page-number pagination
NO search projection used as uniqueness/source of truth
NO Foundation/Person-specific knowledge in core/search
NO UUID/GID-specific cursor assumptions in core/search
YES search/search.go as the supported façade
YES QueryFilter/QuerySort/QueryResult for repository-side models
YES Filter/Sort/Result for external request/response models
```

> **Guiding principle:** one search capability, one public import, explicit transport-vs-query names, resource-owned policy, repository-owned field mapping, and generic mechanics implemented once.

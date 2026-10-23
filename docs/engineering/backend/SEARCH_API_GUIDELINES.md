# GroveHub Search API Guidelines

> **Status:** Mandatory GroveHub search standard  
> **Applies to:** shared search contracts, filters, sorting, cursor pagination, service validation, repository query construction, and search indexing.

---

## 1. Shared Boundary

Every GroveHub search endpoint uses the external HTTP contract from:

```text
backend/libs/core/contracts/search/
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

Cursor pagination is mandatory. Do not mix page-number/offset semantics into this contract.

The raw `contracts/search.Request` must not be passed directly to a repository.

---

## 2. Shared Validated Search Model

After service validation and normalization, every service uses the generic internal search model from:

```text
backend/libs/core/search/
```

Canonical internal types:

```go
search.Filter
search.Sort
search.Cursor[TId]
search.Query[TId]
search.Result[T]
```

Example for Foundation Person:

```go
search.Query[entity.GID]
search.Result[entity.Person]
```

Example for another resource:

```go
search.Query[entity.LeadId]
search.Result[entity.Lead]
```

Do not create resource-specific wrappers such as:

```text
PersonSearchFilter
PersonSearchSort
PersonSearchCursor
PersonSearchQuery
PersonSearchResult
```

when they only duplicate the shared search model.

Resource-specific types are justified only when they contain genuinely resource-specific state or behavior that cannot be represented by the generic search model.

---

## 3. Shared Technical Search Engine

Reusable search mechanics live in:

```text
backend/libs/core/search/
```

The shared engine owns generic mechanics only:

```text
validated Filter / Sort / Cursor / Query / Result models
query arguments/placeholders
EQUALS / NOT_EQUALS
IN / NIN
EXISTS / NOT_EXISTS
GT / GTE / LT / LTE
CONTAINS / STARTS_WITH
LIKE escaping
scalar field predicates
relation/EXISTS predicates
keyset cursor predicates
opaque cursor payload encoding/decoding
generic invalid-query errors
```

Operator behavior is registered through factories. Services/repositories must not recreate large operator switch blocks.

---

## 4. Service-Owned Search Policy

Each service owns a registry of supported public fields and operators.

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

The service translates:

```text
contracts/search.Request
        ↓
service validation + normalization
        ↓
search.Query[TId]
```

Unknown fields/operators are rejected before repository execution.

The service also validates its concrete cursor identity type. For example, Foundation validates that `search.Cursor[entity.GID].CursorId` is a valid Person UUIDv7.

---

## 5. Repository-Owned Field Mapping

Repositories receive only validated generic `search.Query[TId]` values.

The repository implementation owns a field registry mapping public field names to safe SQL definitions:

```text
firstName -> LOWER(p.first_name)
createdAt -> p.created_at ::timestamptz
contacts.email -> person_contacts relation + EMAIL fixed predicate
```

Raw client field names must never be concatenated into SQL.

Repository implementations call the shared `libs/core/search` factories for generic operator behavior.

Repository results use:

```go
search.Result[Entity]
```

rather than resource-specific search result wrappers.

---

## 6. No Service-Specific Operator Switch Trees

Do not implement patterns such as:

```go
switch operator {
case EQUALS:
case IN:
case EXISTS:
case CONTAINS:
...
}
```

inside every service or repository.

Use:

```text
operator -> shared factory
field -> service/repository specification
```

A new generic operator should be implemented once in `libs/core/search` and then explicitly enabled by each service field registry that supports it.

---

## 7. Cursor Rules

Cursor values are opaque to clients.

The shared search package owns serialization/encoding and generic keyset predicate construction. The owning service/repository supplies the stable cursor identity for its resource.

The shared core must never assume that the stable cursor ID is a GroveHub GID, UUID, or any specific domain identifier.

Examples of valid service-owned cursor IDs include:

```text
Person          -> GID
Lead            -> leadId
Communication   -> communicationId
Admission       -> admissionId
```

The validated internal cursor is typed:

```go
search.Cursor[TId]
```

and contains:

```text
Value
CursorId
```

The opaque encoded cursor payload contains only transport-safe values and is converted to the typed cursor after service validation.

The generic cursor predicate receives:

```text
sort expression
sort direction
sort value
cursorId
cursorExpression
```

`cursorExpression` is owned by the repository mapping and may point to any stable unique field appropriate for that resource.

Example:

```text
Foundation Person:
    cursorId         = person GID value
    cursorExpression = p.gid

Admissions Lead:
    cursorId         = leadId value
    cursorExpression = lead.lead_id
```

The shared core must not hardcode casts such as `::uuid` for the cursor ID. PostgreSQL type resolution comes from the owning repository's cursor expression. If a resource requires an explicit expression/cast, the repository owns that expression.

The owning service validates that the decoded cursor matches:

```text
requested sort field
requested direction
sort field value type
service-specific cursorId type
```

Repositories use keyset predicates, never `OFFSET` for the shared search contract.

---

## 8. Search Projection vs Uniqueness

`search_details` is derived search projection data for human-friendly free-text matching.

It is not authoritative and must not be used for uniqueness/deduplication.

Foundation example:

```text
exact normalized person_contacts indexes -> duplicate/GID protection
persons.search_details + trigram index -> searchText
```

Whenever names/contacts change, refresh `search_details` in the same persistence transaction.

---

## 9. Complexity Rules

Search code must remain SRP/OCP friendly:

```text
external contract       -> libs/core/contracts/search
validated generic model -> libs/core/search
shared SQL mechanics    -> libs/core/search
service policy          -> services/helpers/<resource>.go
repository mapping      -> repository/impl/<resource>.go
use-case orchestration  -> services/<resource>.go
```

Do not create `<resource>_search.go` by default. Split only when search has become a genuinely independent responsibility rather than a use case of the resource.

Do not duplicate generic search DTOs inside a resource repository contract.

---

## 10. Hard Rules

```text
NO raw contracts/search.Request passed directly to repositories
NO resource-specific wrappers that merely duplicate search.Filter/Sort/Cursor/Query/Result
NO raw client field -> SQL identifier mapping
NO unparameterized filter values
NO copied operator switch tree per repository
NO copied cursor codec per service
NO offset/page-number pagination in the shared contract
NO search projection as source of truth
NO generic core package containing Person/Foundation-specific knowledge
NO GID/UUID-specific cursor logic in libs/core/search
NO hardcoded cursor-id cast such as ::uuid in generic cursor mechanics
```

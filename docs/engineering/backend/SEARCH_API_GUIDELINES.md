# GroveHub Search API Guidelines

> **Status:** Mandatory GroveHub search standard  
> **Applies to:** shared search contracts, filters, sorting, cursor pagination, service validation, repository query construction, and search indexing.

---

## 1. Shared Boundary

Every GroveHub search endpoint uses:

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

---

## 2. Shared Technical Search Engine

Reusable search mechanics live in:

```text
backend/libs/core/search/
```

The shared engine owns generic mechanics only:

```text
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
opaque cursor encoding/decoding
```

Operator behavior is registered through factories. Services/repositories must not recreate large operator switch blocks.

---

## 3. Service-Owned Search Policy

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

The service translates `search.Request` into a validated repository query.

Unknown fields/operators are rejected before repository execution.

---

## 4. Repository-Owned Field Mapping

Repositories receive only validated service-specific search queries.

The repository implementation owns a field registry mapping public field names to safe SQL definitions:

```text
firstName -> LOWER(p.first_name)
createdAt -> p.created_at ::timestamptz
contacts.email -> person_contacts relation + EMAIL fixed predicate
```

Raw client field names must never be concatenated into SQL.

Repository implementations call the shared `libs/core/search` factories for generic operator behavior.

---

## 5. No Service-Specific Operator Switch Trees

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

## 6. Cursor Rules

Cursor values are opaque to clients.

The shared search package owns serialization/encoding. The owning service validates that the cursor matches:

```text
requested sort field
requested direction
field value type
stable tie-breaker ID
```

Repositories use keyset predicates, never `OFFSET` for the shared search contract.

---

## 7. Search Projection vs Uniqueness

`search_details` is derived search projection data for human-friendly free-text matching.

It is not authoritative and must not be used for uniqueness/deduplication.

Foundation example:

```text
exact normalized person_contacts indexes -> duplicate/GID protection
persons.search_details + trigram index -> searchText
```

Whenever names/contacts change, refresh `search_details` in the same persistence transaction.

---

## 8. Complexity Rules

Search code must remain SRP/OCP friendly:

```text
shared contract      -> libs/core/contracts/search
shared SQL mechanics -> libs/core/search
service policy       -> services/helpers/<resource>.go
repository mapping   -> repository/impl/<resource>.go
use-case orchestration -> services/<resource>.go
```

Do not create `<resource>_search.go` by default. Split only when search has become a genuinely independent responsibility rather than a use case of the resource.

---

## 9. Hard Rules

```text
NO raw client field -> SQL identifier mapping
NO unparameterized filter values
NO copied operator switch tree per repository
NO copied cursor codec per service
NO offset/page-number pagination in the shared contract
NO search projection as source of truth
NO generic core package containing Person/Foundation-specific knowledge
```

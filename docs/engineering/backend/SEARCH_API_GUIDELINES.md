# GroveHub Search API Guidelines

> **Status:** Mandatory GroveHub search contract standard  
> **Applies to:** backend REST search endpoints, shared search DTOs, service validation, cursor pagination, repository query building, PostgreSQL search projections, filters, sorting, and indexing  
> **Precedence:** this document is authoritative for GroveHub search APIs when older generic pagination/search examples conflict with it.

---

## 1. Objective

GroveHub search APIs must expose one reusable request/response shape while allowing each owning service to decide what can actually be searched.

The shared contract lives at:

```text
backend/libs/core/contracts/search/
```

The shared package owns transport/application contract types only. It does not own SQL, service-specific fields, table names, or repository behavior.

---

## 2. Canonical Endpoint Shape

Search is an explicit POST operation on the owning collection:

```text
POST /<service>/<resources>/search
```

Example:

```text
POST /foundation/persons/search
```

Use POST because the request may contain structured filters, sort information, and an opaque cursor that do not belong in an unbounded query string.

---

## 3. Shared Request Contract

Canonical request:

```json
{
  "searchText": "arun",
  "filters": [
    {
      "field": "status",
      "operator": "IN",
      "values": ["ACTIVE"]
    }
  ],
  "sort": [
    {
      "field": "createdAt",
      "direction": "DESC"
    }
  ],
  "pageInfo": {
    "rows": 25,
    "cursor": "opaque-value-from-previous-response"
  }
}
```

The shared Go contract is represented by:

```text
search.Request
search.Filter
search.Operator
search.Sort
search.Direction
search.PageRequest
```

Do not create service-specific copies of these generic transport structures.

---

## 4. Canonical Operators

The initial shared operator set is:

```text
EQUALS
NOT_EQUALS
IN
NIN
EXISTS
NOT_EXISTS
GT
GTE
LT
LTE
CONTAINS
STARTS_WITH
```

Semantics:

```text
GT   -> >
GTE  -> >=
LT   -> <
LTE  -> <=
```

Do not add duplicate aliases such as `GTE_EQUALS` for the same comparison.

A shared operator existing does not mean every service field accepts it.

---

## 5. Service-Owned Search Schema

Each service owns an explicit allowlist mapping:

```text
public search field
        ↓
allowed operators
        ↓
normalization / value type
        ↓
repository-safe SQL expression
```

Clients must never send database column names that are passed directly into SQL.

Correct:

```text
firstName -> LOWER(persons.first_name)
createdAt -> persons.created_at
contacts.email -> EXISTS(...) over service-owned contact rows
```

Forbidden:

```text
request.field -> concatenate directly into SQL
```

Unknown fields or unsupported field/operator combinations return a stable `400` service error.

---

## 6. Filter Values

The shared contract uses string values:

```json
{
  "field": "createdAt",
  "operator": "GTE",
  "values": ["2026-10-01T00:00:00Z"]
}
```

The owning service parses and normalizes values according to the field type before the repository receives them.

Examples:

```text
UUID field      -> validate UUID contract
status          -> normalize/validate service status
instant         -> parse RFC3339 and normalize timezone
email           -> service-defined canonical email normalization
phone           -> service-defined canonical phone normalization
name            -> trim/case-normalize for search semantics
```

`EXISTS` and `NOT_EXISTS` take no values.

`IN` and `NIN` require one or more values.

Single-value operators require exactly one value.

---

## 7. searchText

`searchText` is human-friendly discovery, not a substitute for structured invariants.

The owning service may maintain a denormalized search projection such as:

```text
search_details TEXT
```

The projection may include normalized copies of searchable information from multiple service-owned tables.

Example Person projection:

```text
gid
first name
middle name
last name
email values
phone values
```

Rules:

```text
search_details is derived/search-only data
search_details is not the source of truth
structured columns/tables remain authoritative
search_details is refreshed in the same transaction as authoritative changes
search_details must have an indexing plan
```

Do not use a search projection for exact uniqueness or duplicate-prevention decisions.

---

## 8. Exact Matching vs Free-Text Search

These are different operations.

```text
Exact normalized lookup
    -> business invariant / duplicate detection
    -> structured indexed columns

searchText
    -> user discovery
    -> indexed search projection
```

For Foundation Person creation:

```text
normalized email/phone exact match
        ↓
existing GID found
        ↓
return existing Person
```

The generic `/search` endpoint must not be used as the uniqueness guard for creation.

Database unique constraints/indexes must still protect the invariant against races.

---

## 9. Cursor Pagination

GroveHub search uses cursor/keyset pagination.

Canonical page request:

```json
{
  "rows": 25,
  "cursor": "opaque-value"
}
```

Do not mix database page numbers/offset pagination with cursor pagination in the shared contract.

Canonical response:

```json
{
  "items": [],
  "pageInfo": {
    "rows": 25,
    "nextCursor": "opaque-next-value",
    "hasNext": true
  }
}
```

Rules:

```text
cursor is opaque to clients
cursor encodes only the state required to continue the validated sort
cursor must include a stable tie-breaker
cursor must be validated against the current sort
cursor must not contain secrets
invalid cursor -> stable 400 error
```

A typical stable ordering is:

```text
created_at DESC, id DESC
```

The owning service chooses its default sort and tie-breaker.

---

## 10. Sorting

The shared request allows sort descriptors, but each service explicitly allowlists sortable fields.

Direction values are:

```text
ASC
DESC
```

Dynamic SQL sort expressions must come only from server-owned mappings.

Never interpolate arbitrary client field strings into `ORDER BY`.

A service may initially support fewer sort descriptors than the shared contract can represent. That restriction must be documented in its OpenAPI contract.

---

## 11. PostgreSQL Rules

Search repository implementations must follow the database standard:

```text
parameterized values
whitelisted dynamic identifiers/expressions
bounded rows
stable ordering
appropriate indexes
no unbounded request-path scan
```

For partial text discovery, PostgreSQL trigram/GIN or another deliberately selected indexed strategy may be used when appropriate to the data shape.

For exact values such as normalized email/phone, use ordinary structured indexes/unique indexes rather than the free-text projection.

---

## 12. Shared vs Service-Specific Ownership

Shared:

```text
backend/libs/core/contracts/search/
    request/response contract
    operators
    sort direction
    page request/response
    default/max generic rows
```

Service-specific:

```text
searchable fields
allowed operators per field
normalization
business limits stricter than shared limits
cursor payload contents
repository query model
SQL mapping
search projection contents
indexes
stable semantic errors
```

Do not move service-specific field names or SQL into `libs/core/contracts/search`.

---

## 13. Foundation Person Search

Foundation currently exposes:

```text
POST /foundation/persons/search
```

Searchable fields:

```text
gid
firstName
middleName
lastName
status
createdAt
updatedAt
contacts.email
contacts.phone
contacts.primaryEmail
contacts.primaryPhone
contacts.verifiedEmail
contacts.verifiedPhone
```

The contact variants are intentional compound semantics. For example:

```text
contacts.primaryEmail EQUALS x
```

must match the same contact row where both:

```text
normalized email = x
is_primary = true
```

Do not emulate this with independent filters that could match different contact rows.

Foundation limits the initial request to:

```text
maximum 100 rows
maximum 20 filters
maximum 50 values per filter
maximum 1 explicit sort descriptor
maximum 256 searchText characters
```

Default sort:

```text
createdAt DESC, gid DESC
```

---

## 14. Foundation Person Aggregate Rule

Contacts are nested Person values in the public contract:

```text
Person
└── contacts
    ├── emails[]
    └── phones[]
```

There is no public contact ID and no `/persons/{gid}/contacts` API.

The database may normalize contacts into a service-owned table for indexing, uniqueness, and query performance. That persistence shape does not make Contact a separately addressable API resource.

---

## 15. Mandatory Agent Rules

When creating or changing a GroveHub search API:

1. Reuse `backend/libs/core/contracts/search`; do not clone the generic request contract.
2. Define a service-owned field allowlist.
3. Define allowed operators per field.
4. Normalize/validate values before repository query construction.
5. Keep all SQL values parameterized.
6. Never pass raw client field/sort strings into SQL.
7. Use cursor/keyset pagination; do not add offset/page-number behavior to the shared contract.
8. Include a stable tie-breaker in cursor ordering.
9. Keep cursors opaque and non-secret.
10. Use structured exact indexes for uniqueness/deduplication.
11. Treat `search_details` or equivalent as derived data only.
12. Refresh search projections transactionally with their source data.
13. Document searchable fields/operators in the service OpenAPI contract.
14. Add an indexing plan for every request-path search.
15. Return stable service errors for invalid search fields/operators/cursors.

---

## 16. Guiding Principle

> **The shared contract standardizes how GroveHub asks for search; the owning service controls what may be searched, how values are normalized, and how PostgreSQL executes it safely and efficiently.**

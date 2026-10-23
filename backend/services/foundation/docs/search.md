# Foundation Person Search

> **Status:** Authoritative Foundation Person search contract  
> **Generic standard:** `docs/engineering/backend/SEARCH_API_GUIDELINES.md`

## 1. What this document owns

This document defines Foundation-specific Person search behavior:

- Person search endpoint
- Foundation client search method
- supported Person fields
- allowed operators per field
- normalization rules
- sorting
- cursor behavior
- `search_details`
- Person-specific validation/error behavior

Generic search mechanics belong to `libs/core/contracts/search` and `libs/core/search` and must not be copied into Foundation-specific types.

## 2. Entry points

HTTP:

```text
POST /foundation/persons/search
```

In-process:

```go
foundationClient.SearchPersons(ctx, request)
```

Both entry points execute the same Foundation Person service behavior.

## 3. Request lifecycle

When another service searches Foundation during an existing request, it must pass the same context:

```go
result, err := foundationClient.SearchPersons(request.Context(), searchRequest)
```

The client does not create an HTTP request and does not replace the context.

## 4. Public search request

Foundation accepts the shared HTTP contract:

```json
{
  "searchText": "arun",
  "filters": [
    {
      "field": "status",
      "operator": "EQUALS",
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
    "rows": 20,
    "cursor": null
  }
}
```

The raw request is validated and normalized by Foundation before repository execution.

```text
contracts/search.Request
        ↓
Foundation Person field/operator validation
        ↓
core/search.Query[entity.GID]
        ↓
Person repository
        ↓
core/search.Result[entity.Person]
```

Do not recreate `PersonSearchFilter`, `PersonSearchSort`, `PersonSearchCursor`, `PersonSearchQuery`, or `PersonSearchResult` wrappers unless Foundation later gains genuinely Person-specific search state that cannot be represented by the generic model.

## 5. Operators

Shared operators currently available to Foundation fields include:

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

An operator being available in shared core does not automatically mean every Person field supports it.

## 6. Person field matrix

| Field | Supported operators |
|---|---|
| `gid` | `EQUALS`, `NOT_EQUALS`, `IN`, `NIN` |
| `firstName` | `EQUALS`, `NOT_EQUALS`, `IN`, `NIN`, `EXISTS`, `NOT_EXISTS`, `CONTAINS`, `STARTS_WITH` |
| `middleName` | `EQUALS`, `NOT_EQUALS`, `IN`, `NIN`, `EXISTS`, `NOT_EXISTS`, `CONTAINS`, `STARTS_WITH` |
| `lastName` | `EQUALS`, `NOT_EQUALS`, `IN`, `NIN`, `EXISTS`, `NOT_EXISTS`, `CONTAINS`, `STARTS_WITH` |
| `status` | `EQUALS`, `NOT_EQUALS`, `IN`, `NIN` |
| `createdAt` | `EQUALS`, `GT`, `GTE`, `LT`, `LTE` |
| `updatedAt` | `EQUALS`, `GT`, `GTE`, `LT`, `LTE` |
| `contacts.email` | `EQUALS`, `NOT_EQUALS`, `IN`, `NIN`, `EXISTS`, `NOT_EXISTS`, `CONTAINS`, `STARTS_WITH` |
| `contacts.phone` | same contact operators |
| `contacts.primaryEmail` | same contact operators, restricted to primary email contacts |
| `contacts.primaryPhone` | same contact operators, restricted to primary phone contacts |
| `contacts.verifiedEmail` | same contact operators, restricted to verified email contacts |
| `contacts.verifiedPhone` | same contact operators, restricted to verified phone contacts |

Unknown fields or unsupported field/operator combinations are rejected before repository execution.

## 7. Normalization

Foundation owns Person-specific normalization.

```text
gid        -> canonical UUIDv7 string
name       -> trim; comparison/search value normalized to lower case
status     -> upper case and validate against Person statuses
time       -> RFC3339 input normalized to UTC RFC3339Nano
email      -> Foundation email normalization
phone      -> Foundation phone normalization
```

Partial email/phone operators such as `CONTAINS` and `STARTS_WITH` use partial-safe normalization instead of requiring a complete email/phone value.

## 8. Free-text search

`searchText` searches the derived `foundation.persons.search_details` projection.

The projection includes searchable Person identity/profile/contact text such as:

```text
gid
first name
middle name
last name
email values
normalized email values
phone values
normalized phone values
```

`search_details` is for discovery only.

It must never be used for Person uniqueness or duplicate detection.

```text
exact normalized person_contacts index -> duplicate protection
search_details + trigram index          -> free-text search
```

## 9. Sorting

Supported sort fields:

```text
gid
firstName
lastName
createdAt
updatedAt
```

Default:

```text
createdAt DESC
```

Only one public sort field is currently accepted per request.

## 10. Cursor pagination

Foundation uses keyset/cursor pagination, never offset pagination for this search contract.

The external cursor is opaque.

Internally the generic core cursor model is:

```go
type Cursor[TId any] struct {
    Value    string
    CursorId TId
}
```

For Person:

```text
TId      = entity.GID
CursorId = Person GID
```

This does not make GID a generic search-core concept. Other services can use `leadId`, `communicationId`, or another stable identifier type.

The encoded payload also records the requested sort field and direction. A cursor created for one sort cannot be reused with a different sort.

Example first page:

```json
{
  "sort": [
    {
      "field": "createdAt",
      "direction": "DESC"
    }
  ],
  "pageInfo": {
    "rows": 20
  }
}
```

Example next page:

```json
{
  "sort": [
    {
      "field": "createdAt",
      "direction": "DESC"
    }
  ],
  "pageInfo": {
    "rows": 20,
    "cursor": "<opaque nextCursor>"
  }
}
```

## 11. Response

```json
{
  "items": [
    {
      "gid": "019...",
      "firstName": "Arun",
      "lastName": "Sharma",
      "contacts": {
        "emails": [],
        "phones": []
      },
      "status": "ACTIVE",
      "version": 1,
      "createdAt": "2026-10-24T00:00:00Z",
      "updatedAt": "2026-10-24T00:00:00Z"
    }
  ],
  "pageInfo": {
    "rows": 20,
    "nextCursor": "<opaque cursor>",
    "hasNext": true
  }
}
```

The HTTP response is wrapped by GroveHub's normal success envelope.

## 12. Examples

Exact email:

```json
{
  "filters": [
    {
      "field": "contacts.email",
      "operator": "EQUALS",
      "values": ["arun@gmail.com"]
    }
  ],
  "pageInfo": {
    "rows": 20
  }
}
```

Active Persons created after a timestamp:

```json
{
  "filters": [
    {
      "field": "status",
      "operator": "EQUALS",
      "values": ["ACTIVE"]
    },
    {
      "field": "createdAt",
      "operator": "GTE",
      "values": ["2026-10-01T00:00:00Z"]
    }
  ],
  "sort": [
    {
      "field": "createdAt",
      "direction": "DESC"
    }
  ],
  "pageInfo": {
    "rows": 25
  }
}
```

Persons without an email:

```json
{
  "filters": [
    {
      "field": "contacts.email",
      "operator": "NOT_EXISTS",
      "values": []
    }
  ],
  "pageInfo": {
    "rows": 20
  }
}
```

## 13. Hard rules

```text
NO raw search Request directly to repository
NO Foundation-specific copies of generic Query/Filter/Sort/Cursor/Result
NO arbitrary client field -> SQL identifier interpolation
NO offset pagination
NO GID assumptions in libs/core/search
NO duplicate detection through search_details
NO service-to-service HTTP call when Foundation client is available in-process
NO context.Background() replacing the consuming request context
```

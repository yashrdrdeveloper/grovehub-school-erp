# Foundation Contracts

> **Status:** Authoritative Foundation service contract  
> **Audience:** Foundation maintainers and consuming services such as Identity  
> **Scope:** Person ownership, public request/response DTOs, HTTP surface, in-process client surface, lifecycle behavior, and service boundaries.

## 1. What this document owns

This document defines **what Foundation exposes** and **how other GroveHub services consume it**.

Use this document for:

- Person create/read/update/deactivate behavior
- Person request and response shapes
- contact semantics
- duplicate-person behavior
- public HTTP endpoints
- in-process Foundation client usage
- service-to-service dependency rules

Search-specific rules and examples live in [`search.md`](./search.md).

The cross-service generic search standard lives in:

```text
docs/engineering/backend/SEARCH_API_GUIDELINES.md
```

## 2. Foundation ownership

Foundation owns the canonical GroveHub Person.

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

Foundation is the source of truth for human/profile/contact data.

Identity owns authentication/account/security concerns. Identity may reference a Foundation Person by GID, but Identity must not duplicate or mutate Foundation Person state directly.

## 3. Public Foundation surfaces

Foundation exposes the same business behavior through two entry points.

### HTTP

Use HTTP for external callers and requests entering GroveHub through the edge.

```text
POST   /foundation/persons
GET    /foundation/persons/{gid}
PATCH  /foundation/persons/{gid}
POST   /foundation/persons/search
PATCH  /foundation/persons/{gid}/deactivate
```

### In-process client

Other services inside the GroveHub process must use:

```text
backend/services/foundation/client
```

Do not import Foundation repositories, services, controllers, or database implementation packages from another service.

The client exposes:

```go
CreatePerson(...)
GetPersonByGID(...)
UpdatePerson(...)
SearchPersons(...)
DeactivatePerson(...)
```

The client delegates to the same Foundation service used by the HTTP controller. It does **not** call Foundation over HTTP.

## 4. Request lifecycle rule

Every client method accepts `context.Context`.

A consuming service must forward its current request context:

```go
person, err := foundationClient.GetPersonByGID(
    request.Context(),
    foundationrequest.PersonByGID{GID: gid},
)
```

Do not replace the request context with `context.Background()` inside service-to-service calls.

This keeps cancellation, deadlines, tracing, request correlation, and other context values on the same request lifecycle.

```text
incoming request
    ↓
Identity controller/service
    ↓ same context.Context
Foundation client
    ↓ same context.Context
Foundation service
    ↓
Foundation repository
```

## 5. Create Person

Request:

```json
{
  "firstName": "Arun",
  "middleName": "Kumar",
  "lastName": "Sharma",
  "contacts": {
    "emails": [
      {
        "value": "arun@gmail.com",
        "isPrimary": true
      }
    ],
    "phones": [
      {
        "value": "+919876543210",
        "isPrimary": true
      }
    ]
  }
}
```

`isVerified` is not writable through the public create/update DTO. Verification is Foundation-controlled state.

### Duplicate behavior

Foundation normalizes submitted email/phone contacts before creation.

```text
normalize contacts
    ↓
lookup exact normalized contacts
    ↓
0 Person GIDs -> create Person
1 Person GID  -> return existing Person
>1 Person GID -> contact conflict
```

HTTP behavior:

```text
new Person      -> 201 Created
existing Person -> 200 OK
```

The client returns the same distinction through the `created bool` result from `CreatePerson`.

## 6. Person response

Canonical response DTO:

```json
{
  "gid": "019...",
  "firstName": "Arun",
  "middleName": "Kumar",
  "lastName": "Sharma",
  "contacts": {
    "emails": [
      {
        "value": "arun@gmail.com",
        "isPrimary": true,
        "isVerified": false
      }
    ],
    "phones": [
      {
        "value": "+919876543210",
        "isPrimary": true,
        "isVerified": false
      }
    ]
  },
  "status": "ACTIVE",
  "version": 1,
  "createdAt": "2026-10-24T00:00:00Z",
  "updatedAt": "2026-10-24T00:00:00Z"
}
```

Contacts do not expose a public contact ID.

## 7. Get Person

HTTP:

```text
GET /foundation/persons/{gid}
```

In-process:

```go
foundationClient.GetPersonByGID(ctx, foundationrequest.PersonByGID{GID: gid})
```

The GID must be a valid Foundation Person GID. A valid GID that does not exist returns the stable Person-not-found application error.

## 8. Update Person

HTTP:

```text
PATCH /foundation/persons/{gid}
```

Only supplied top-level fields are updated.

Example:

```json
{
  "firstName": "Arun Updated"
}
```

### Contact replacement semantics

`contacts` is optional.

```text
contacts omitted -> contacts unchanged
contacts present -> supplied contacts become desired aggregate contact state
```

Example:

```json
{
  "contacts": {
    "emails": [],
    "phones": [
      {
        "value": "+919876543210",
        "isPrimary": true
      }
    ]
  }
}
```

This clears email contacts and retains only the supplied phone contacts.

Existing verification state is preserved for unchanged normalized contacts.

## 9. Deactivate Person

HTTP:

```text
PATCH /foundation/persons/{gid}/deactivate
```

In-process:

```go
foundationClient.DeactivatePerson(ctx, foundationrequest.PersonByGID{GID: gid})
```

Deactivation is idempotent. It changes Person status; it does not delete the Person or contacts.

## 10. Contact rules

- zero or many emails are allowed
- zero or many phones are allowed
- at most one primary contact per contact type
- email normalization is case-insensitive
- phone normalization produces the canonical international representation used by Foundation
- normalized `(type, value)` is globally unique for Person duplicate protection
- contact verification is controlled by Foundation

## 11. Consumer rules

A consuming service such as Identity must:

```text
DO import foundation/client
DO use Foundation request/response contracts
DO forward the current context.Context
DO store/reference the Person GID when identity linkage is required

DO NOT import foundation/repository
DO NOT import foundation/repository/impl
DO NOT call Foundation tables directly
DO NOT copy Person DTOs into Identity merely to avoid the client
DO NOT make an internal HTTP request when the in-process client is available
DO NOT change Foundation contact data as a side effect of Identity login changes
```

## 12. Composition

The launcher owns dependency wiring.

```text
launcher
├── Foundation module
│   ├── Person repository
│   ├── Person service
│   ├── HTTP controller
│   └── Foundation client
│
└── Identity module
    └── receives Foundation client
```

Foundation does not depend on Identity.

Identity may depend on the Foundation public client.

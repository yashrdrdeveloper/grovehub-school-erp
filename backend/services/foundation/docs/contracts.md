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

## 3. Foundation contracts and HTTP surfaces

Foundation exposes the same business behavior through two entry points.

### HTTP

Use HTTP for requests entering GroveHub through the edge. The `protected` path segment declares the intended **future gateway policy**; it does not enforce authentication today. Do not expose these person-data routes to untrusted clients until Identity/gateway enforcement is implemented.

Foundation application errors declare HTTP statuses alongside their code and
message in `constants/error/person.go`, via `coreerrors.NewDefinition`.
The controller reads `Definition().HTTPStatus()` directly; no second status
map is maintained in the controller. Unexpected non-application errors use
HTTP 500 and the configured safe fallback response.

```text
POST   /foundation/protected/persons
GET    /foundation/protected/persons/{gid}
PATCH  /foundation/protected/persons/{gid}
POST   /foundation/protected/persons/search
PATCH  /foundation/protected/persons/{gid}/deactivate
```

### In-process client

Other services inside the GroveHub process must depend on the public client contract:

```text
backend/services/foundation/client/person/
├── person.go
└── impl/
    └── person.go
```

Ownership:

```text
client/person/person.go
    -> public Client interface only

client/person/impl/person.go
    -> concrete in-process implementation
    -> stores an internal HTTP router registered with Foundation's controller
    -> generated @RequestLine methods dispatch through the REST pipeline
```

Consumers such as Identity import only:

```text
grovehub/services/foundation/client/person
```

They must never import:

```text
grovehub/services/foundation/client/person/person/impl
```

Only Foundation composition code may construct the concrete implementation.

Do not import Foundation repositories, services, controllers, or database implementation packages from another service.

The client exposes:

```go
CreatePerson(...)
GetPersonByGID(...)
UpdatePerson(...)
SearchPersons(...)
DeactivatePerson(...)
```

The `@RequestLine` client uses one concrete request DTO per method,
embedding `corerequest.Request`, and **only** declares the needed sections.
The reusable core introspector maps a declared `Params` struct into route
placeholders, `Query` into URL-encoded query parameters, and
`RequestBody` into JSON. Missing sections are absent automatically;
no generic `Empty` placeholders are required. The controller receives
and validates the exact same struct, then forwards it to the business
service unchanged. For example:

```go
type UpdatePersonRequest struct {
    corerequest.Request
    Params PersonByGID `validate:"valid"`
    RequestBody UpdatePerson `validate:"valid"`
}
```

A `GetPersonByGIDRequest` declares just `Params`; a
`SearchPersonsRequest` declares just `RequestBody`. Routes are
service-prefix-relative (e.g. `@RequestLine POST /protected/persons`).
The code generator verifies the registered controller route and dispatches
in-process without TCP. HTTP 201 maps to `created=true`; 200 maps to
`false`. Safe error codes, HTTP status and public metadata are preserved.

The internal router does **not** run the launcher's outer IAM/gateway middleware; URL tiers are not access controls today. Revisit this trust boundary before enforcing protected/restricted policies or exposing the handler to untrusted callers.

## 4. Request lifecycle rule

Every client method accepts `context.Context`.

A consuming service must forward its current request context:

```go
person, err := foundationClient.GetPersonByGID(
    request.Context(),
    foundationrequest.GetPersonByGIDRequest{
        Params: foundationrequest.PersonByGID{GID: gid},
    },
)
```

Do not replace the request context with `context.Background()` inside service-to-service calls.

This keeps cancellation, deadlines, tracing, request correlation, and other context values on the same request lifecycle.

```text
incoming request
    ↓
Identity controller/service
    ↓ same context.Context
Foundation client interface
    ↓
Foundation client implementation
    ↓ same context.Context
Foundation service
    ↓
Foundation repository
```

## DTO layout

Foundation DTOs are resource-based and use a single Go package per resource:

```text
dto/person/operations.go              # CreatePersonRequest, UpdatePersonRequest, etc.
dto/person/request.go                 # CreatePerson, UpdatePerson, PersonNames, Contacts
dto/person/params.go                  # PersonByGID
dto/person/queryparams.go             # Query DTOs when URL query endpoints are introduced
dto/person/response.go                # Person, ResponseContacts, PersonErrorMetadata
libs/core/annotation/dto/patch.go     # shared presence-aware Patch[T] type
```

Go consumers import `grovehub/services/foundation/dto/person` (typically `persondto`). All request, params, operation and response structs share the `person` package, so response-only contact types are named `ResponseContacts` and `ResponseContactValue`. `@PatchRequestDto` generates its JSON serializer under `.generated/dto/person` and preserves absent/null/supplied semantics.

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
GET /foundation/protected/persons/{gid}
```

In-process:

```go
foundationClient.GetPersonByGID(ctx, foundationrequest.PersonByGID{GID: gid})
```

The GID must be a valid Foundation Person GID. A valid GID that does not exist returns the stable Person-not-found application error.

## Repository search implementation

The repository interface in `repository/person.go` exposes a plain `Search` method.
The concrete `repository/person/impl/person.go` owns its `@Searchable config=search.PersonConfig alias=p` annotation and an explicit `Search(ctx, query)` method that delegates to `repository.searchGenerated(ctx, query)`.

The generator validates the interface signature, `@Repository entity=Person` ownership, approved search configuration and stable cursor, then emits the private query helper in `.generated/repository/person/impl/zz_search.gen.go`. Existing search filters and field expressions remain in `repository/person/helpers/search/person.go`.

## HTTP request execution pipeline

Foundation Person controllers use `rest.RunStruct(w,r,...)` with the
operation-specific struct and shared request rules. The binder reads only
declared sections and the service receives the same validated struct. Services
also validate relevant fields when invoked directly. Invalid path/body/query
inputs return semantic 4xx responses, and domain errors retain status and
safe metadata. Create returns 201 for new and 200 for existing Persons.

## 8. Update Person

HTTP:

```text
PATCH /foundation/protected/persons/{gid}
```

Only supplied top-level fields are updated.

Example:

```json
{
  "firstName": "Arun Updated"
}
```

### PATCH presence and clearing semantics

```text
firstName / lastName omitted -> unchanged
firstName / lastName null    -> 400 (required names cannot be cleared)
firstName / lastName string  -> trim, validate, update
middleName omitted          -> unchanged
middleName null             -> clear (persist NULL)
middleName string           -> trim, validate, update (empty string also clears)
contacts omitted            -> unchanged
contacts null               -> remove every email and phone contact
contacts object             -> replace all contacts with the supplied arrays
```

The DTO uses `annotation/dto.Patch[T]` with `Set` for JSON presence and `Value` for null/value; its `@PatchRequestDto` annotation generates JSON serialization, without resource-specific helpers; nested `emails` or `phones` omitted from a supplied contacts object means an empty group.

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
PATCH /foundation/protected/persons/{gid}/deactivate
```

In-process:

```go
foundationClient.DeactivatePerson(ctx, foundationrequest.PersonByGID{GID: gid})
```

Deactivation is idempotent. It changes Person status; it does not delete the Person or contacts. PersonRepository intentionally has no public hard-Delete method.

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
DO import foundation/client/person
DO depend on client.Client
DO use Foundation request/response contracts
DO forward the current context.Context
DO store/reference the Person GID when identity linkage is required

DO NOT import foundation/client/person/person/impl
DO NOT import foundation/repository/<resource>
DO NOT import foundation/repository/<resource>/<resource>/impl
DO NOT import foundation/services/person directly
DO NOT call Foundation tables directly
DO NOT copy Person DTOs into Identity merely to avoid the client
DO NOT create an external HTTP/TCP loopback request; use the generated in-process endpoint client
DO NOT change Foundation contact data as a side effect of Identity login changes
```

## 12. Composition

The service root owns Foundation concrete construction, while the launcher passes only the public client contract to consumers.

```text
Foundation module
├── PersonRepository (persons table, 5 methods)
├── ContactRepository (person_contacts table, 4 methods)
├── core/transaction.UnitOfWork (reusable shared transaction)
├── repository/constants/errors/ (sentinel persistence errors)
├── repository/<resource>/helpers/errors/ (database error translation)
├── repository/constants/{person,contact}.go (schema constraint names)
├── repository/person/helpers/search/person.go (trusted Search Config)
├── PersonService (sole owner of contact lifecycle)
├── controller/person/protected/person.go (controller interface)
│   └── controller/person/protected/impl/person.go (HTTP adapter)
└── client interface
    └── client/impl

launcher
├── constructs Foundation module
└── passes foundationModule.Client() to Identity
```

Foundation does not depend on Identity.

Identity may depend on the Foundation public client contract.

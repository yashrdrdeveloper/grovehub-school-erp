# API Contract Guidelines

> **Status:** Mandatory Phase 0 engineering standard  
> **Location:** `docs/engineering/api/API_CONTRACT_GUIDELINES.md`  
> **Applies to:** REST/JSON, OpenAPI, gRPC/Protobuf, public clients, events, WebSocket/SSE payloads, request/response DTOs, error contracts, pagination, idempotency, compatibility, generated code, and cross-service/public API evolution  
> **Primary transports:** REST + JSON + OpenAPI for public/browser APIs; gRPC + Protobuf for internal/provider RPC where appropriate; WebSocket/SSE for realtime streams  
> **Core principles:** explicit contracts, additive evolution, stable semantics, opaque identifiers, normalized errors, idempotent mutation where required, generated-code discipline, transport/domain separation, and backward compatibility

---

## 1. Objective

Every Perihellion boundary must be understandable without reading the implementation behind it.

A contract answers:

```text
What operation exists?
What input does it accept?
What output does it return?
What errors are possible?
What authorization is required?
What compatibility guarantee exists?
How is the operation observed/retried?
```

The contract is not the business implementation.

---

## 2. Contract Hierarchy

```text
Domain / Business Semantics
          ↓
Stable DTO / Message Contract
          ↓
Transport Mapping
   ┌──────┼────────┐
   ↓      ↓        ↓
 REST    gRPC    Event/Stream
```

Domain entities must not become transport models merely for convenience.

---

## 3. Mandatory Rules

1. OpenAPI is the source of truth for HTTP contracts.
2. Protobuf is the source of truth for gRPC contracts.
3. Generated code is never hand-edited.
4. Public contracts are versioned by compatibility boundary, not by implementation change.
5. Additive compatible changes stay in the current major version.
6. Breaking changes require explicit compatibility review and normally a new major contract version.
7. Internal business services are not duplicated by API version.
8. DTOs are transport contracts, not persistence entities.
9. Database models are never returned directly.
10. IDs crossing contracts are opaque.
11. Consumers must not infer business meaning from ID format.
12. Timestamps use RFC 3339/ISO-8601 semantics with timezone/offset.
13. Public errors are normalized and never expose SQL, stack traces, secrets, or provider internals.
14. Pagination is explicit.
15. Mutations that may be retried must define idempotency semantics.
16. Request and correlation identifiers propagate across service/provider boundaries.
17. Unknown fields must be tolerated where the serialization protocol allows it.
18. Removal/rename of public fields is a breaking change unless the field was explicitly experimental.
19. Protobuf field numbers are never reused.
20. Deprecated Protobuf fields/numbers are reserved.
21. Enum values must define unknown/unspecified behavior.
22. Boolean fields must not encode multi-state business concepts.
23. Public request validation is structural; business authorization remains in business services.
24. Cross-service calls use the target service's `client/vN`.
25. Event consumers depend on versioned event contracts, not producer internals.
26. Realtime payloads are versioned and resumable where loss/reconnect matters.
27. Contract examples and tests must stay synchronized with the source contract.
28. Compatibility checks run in CI.
29. Documentation does not override a machine-readable contract.
30. Contract design must remain provider-neutral unless the endpoint is explicitly provider administration.

---

## 4. Backend Contract Layout

Recommended service boundary:

```text
backend/services/<service>/
├── api/
│   ├── openapi/
│   │   └── v1/
│   │       └── openapi.yaml
│   └── proto/
│       └── v1/
│           └── service.proto
│
├── dto/
│   └── v1/
│
├── client/
│   └── v1/
│       ├── client.go
│       ├── inprocess/
│       ├── grpc/
│       └── http/
│
├── src/
└── ...
```

The contract folders define what a consumer may depend on.

---

## 5. HTTP Contract Design

Public/browser APIs use:

```text
HTTP
JSON
OpenAPI
standard status codes
normalized errors
```

Avoid custom transport semantics when ordinary HTTP already expresses them.

---

## 6. URL Design

URLs describe resources and stable operations.

Prefer:

```text
/applications
/applications/{applicationId}
/applications/{applicationId}/revisions
/executions/{executionId}/cancel
```

Avoid:

```text
/getApplications
/doCreateApplication
/applicationAction?action=delete
```

Actions that are not natural CRUD operations may use an explicit action subresource/command endpoint.

---

## 7. Version Placement

Contract versioning is major-version based.

Conceptually:

```text
api/openapi/v1
dto/v1
client/v1
```

HTTP routing may expose a major version in the route according to the service's routing convention.

Do not create `v2` because internal code was refactored.

---

## 8. HTTP Methods

Use semantic methods:

```text
GET     read
POST    create / command
PUT     full idempotent replacement when truly supported
PATCH   partial update
DELETE  delete/revoke
```

Do not use `GET` for mutations.

---

## 9. HTTP Status Codes

Typical mappings:

```text
200 OK
201 Created
202 Accepted
204 No Content
400 Bad Request
401 Unauthorized
403 Forbidden
404 Not Found
409 Conflict
412 Precondition Failed
422 Unprocessable Entity where structural syntax is valid but request semantics cannot be accepted
429 Too Many Requests
500 Internal Server Error
503 Service Unavailable
504 Gateway Timeout
```

Use one repository-wide semantic mapping.

---

## 10. Async Operations

Long-running operations should not hold a request open unnecessarily.

Preferred:

```text
POST operation
    ↓
202 Accepted
    ↓
operation/execution ID
    ↓
Mission Runtime / durable execution
    ↓
status/events
```

Response example:

```json
{
  "executionId": "opaque-id",
  "status": "QUEUED"
}
```

---

## 11. Request DTOs

Request DTOs contain only public input.

They must not contain:

```text
database-only columns
server-computed ownership fields
trusted actor identity supplied by the caller
internal provider secrets unless the endpoint explicitly accepts a secret reference/config input
```

Server-derived context stays server-derived.

---

## 12. Response DTOs

Responses should be designed for the contract, not as raw entity dumps.

Use stable fields.

Avoid exposing internal implementation details such as:

```text
table IDs
provider connection strings
internal enum/storage names
stack traces
raw secret values
```

---

## 13. Field Naming

JSON fields use one consistent convention:

```text
camelCase
```

Examples:

```text
applicationId
createdAt
correlationId
providerType
```

Protobuf uses idiomatic `snake_case` field names; generated language bindings map appropriately.

---

## 14. Opaque IDs

External contracts represent identifiers as opaque strings unless a protocol requires another representation.

Example:

```json
{
  "applicationId": "01J..."
}
```

Consumers must not:

```text
parse timestamps from IDs
infer service ownership from prefixes
sort by undocumented ID structure
```

---

## 15. Timestamps

Use RFC 3339 with timezone/offset.

Prefer UTC for emitted server timestamps:

```text
2026-10-01T10:20:30Z
```

Do not emit locale-formatted timestamps in APIs.

---

## 16. Durations

Use explicit units.

Good:

```json
{
  "timeoutSeconds": 30
}
```

or an explicitly documented duration string.

Avoid ambiguous:

```json
{
  "timeout": 30
}
```

---

## 17. Money and Precise Numeric Values

Never assume binary floating point is acceptable for money/cost/accounting semantics.

Contracts should carry:

```text
amount + currency
```

using a precision-safe representation chosen by the owning domain.

---

## 18. Boolean Design

Do not model states such as:

```text
pending
approved
rejected
cancelled
```

with multiple booleans.

Use a typed status/state field.

---

## 19. Enums

Every public enum must have an explicit unknown/unspecified strategy.

Protobuf:

```proto
enum ExecutionStatus {
  EXECUTION_STATUS_UNSPECIFIED = 0;
  EXECUTION_STATUS_CREATED = 1;
  EXECUTION_STATUS_RUNNING = 2;
  EXECUTION_STATUS_SUCCEEDED = 3;
}
```

Never reuse numeric values.

---

## 20. Pagination

List endpoints that can grow unbounded must paginate.

Preferred scalable contract:

```text
cursor
pageSize
nextCursor
```

Example:

```json
{
  "items": [],
  "nextCursor": "opaque-cursor"
}
```

Use offset pagination only when the use case genuinely needs stable positional pages and scale/consistency characteristics are understood.

---

## 21. Filtering

Filters are explicit query parameters.

Example:

```text
GET /applications?environmentId=...&status=ACTIVE
```

Do not overload free-form search with hidden filter syntax unless a documented query language exists.

---

## 22. Sorting

Sorting must use approved public field identifiers.

Example:

```text
sort=createdAt
order=desc
```

Reject unsupported sort keys.

Never interpolate arbitrary sort strings into SQL.

---

## 23. Search

Search semantics must be documented:

```text
case sensitivity
tokenization
fields searched
ranking
pagination
```

Do not promise exact-match semantics when the implementation is full-text search.

---

## 24. Partial Updates

PATCH semantics must be explicit.

Distinguish:

```text
field omitted
field supplied as null
field supplied with value
```

when the contract cares.

Do not rely on Go zero values to infer caller intent.

---

## 25. Optimistic Concurrency

Resources with concurrent edits may expose a version/revision/ETag.

Example:

```text
If-Match: "<revision>"
```

or explicit revision field.

Conflicts return a semantic conflict/precondition error.

---

## 26. Idempotency

Create/execute endpoints that are likely to be retried should support an idempotency key where duplicate execution would be unsafe.

Example:

```text
Idempotency-Key: <opaque-key>
```

The contract must define:

```text
scope
retention window
same-key/same-payload behavior
same-key/different-payload behavior
```

---

## 27. Request ID vs Correlation ID

Use both concepts distinctly.

```text
requestId
  unique to one inbound request/attempt

correlationId
  stable across related calls/workflow when available
```

Services generate a request ID when missing.

Trusted correlation context is propagated, normalized, and bounded in length.

---

## 28. Standard Headers

Approved platform headers should be centralized.

Examples may include:

```text
request/correlation
idempotency
trace context
conditional request headers
```

Do not invent service-specific variants for the same concept.

---

## 29. Authentication Context

A caller may present credentials.

The server resolves trusted identity/context.

Never trust arbitrary caller-supplied:

```text
userId
organizationId
role
permission
tenantId
```

as authorization truth without validating it through trusted authentication context.

---

## 30. Authorization Documentation

OpenAPI should indicate security requirements.

Business-level permission names/semantics should be documented close to the endpoint.

Frontend visibility never substitutes for backend authorization.

---

## 31. Error Envelope

Use one normalized public shape.

Example:

```json
{
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "message": "Application was not found.",
    "requestId": "opaque-request-id",
    "details": {}
  }
}
```

`message` is safe for the intended consumer.

`details` is structured and bounded.

---

## 32. Error Codes

Error codes are stable machine contracts.

Examples:

```text
NOT_FOUND
UNAUTHORIZED
FORBIDDEN
CONFLICT
INVALID_CONFIGURATION
UNSUPPORTED_CAPABILITY
QUOTA_EXCEEDED
TIMEOUT
UNAVAILABLE
```

Do not make consumers parse human text.

---

## 33. Validation Errors

Field validation may return structured field problems.

Example:

```json
{
  "error": {
    "code": "VALIDATION_FAILED",
    "details": {
      "fields": [
        {
          "field": "name",
          "code": "REQUIRED"
        }
      ]
    }
  }
}
```

Do not expose validator internals.

---

## 34. Internal Error Mapping

Internal errors map to public semantics at the transport boundary.

```text
Postgres unique violation
        ↓
semantic conflict
        ↓
HTTP 409 / equivalent gRPC status
```

Raw driver errors never cross the API.

---

## 35. gRPC Status Mapping

Use standard gRPC status semantics.

Examples:

```text
InvalidArgument
Unauthenticated
PermissionDenied
NotFound
AlreadyExists
FailedPrecondition
ResourceExhausted
Unavailable
DeadlineExceeded
Internal
```

Keep semantics aligned with HTTP errors.

---

## 36. Protobuf Package Naming

Use versioned packages.

Example:

```proto
package perihellion.catalog.v1;
```

Generated package options should follow repository conventions.

---

## 37. Protobuf Field Evolution

Allowed:

```text
add optional/new fields
add new messages
add new RPCs
```

Breaking:

```text
change field meaning
change field number
reuse removed field number
change scalar type incompatibly
remove required consumer behavior
```

Reserve removed fields/numbers.

---

## 38. Protobuf `oneof`

Use `oneof` for mutually exclusive variants.

Do not encode variants through several optional booleans/fields that can become contradictory.

---

## 39. Protobuf Maps

Use map fields when the key/value semantics are genuinely open.

Do not use `map<string,string>` to avoid designing a stable contract.

Prefer typed fields for stable business concepts.

---

## 40. Generated Code

Generated output directories are read-only to engineers.

Workflow:

```text
edit OpenAPI/proto
    ↓
validate
    ↓
generate
    ↓
review generated diff
    ↓
run compatibility tests
```

---

## 41. OpenAPI Requirements

Every operation should define:

```text
operationId
summary
request schema
response schema
error responses
security
path/query parameters
examples for non-obvious contracts
```

Operation IDs must remain stable where generated clients depend on them.

---

## 42. OpenAPI Component Reuse

Reusable schemas belong under components.

Do not create slightly different copies of the same stable error/pagination/request metadata contract.

DRY applies to shared semantics, not merely matching syntax.

---

## 43. OpenAPI Validation

CI validates:

```text
syntax
references
operation ID uniqueness
breaking compatibility where tooling supports it
generated-code freshness
```

---

## 44. Public Client Contract

Cross-service consumers use:

```text
client/vN
```

The client exposes business-facing operations, not transport details.

Example:

```go
type Client interface {
    GetApplication(
        ctx context.Context,
        id string,
    ) (Application, error)
}
```

The consumer should not know whether the implementation is in-process, gRPC, or HTTP.

---

## 45. In-Process Client

The in-process client calls the same business boundary that a remote transport reaches.

It must not access:

```text
repository directly
controller
database
private implementation from the caller
```

---

## 46. Client Timeouts

Remote client methods respect caller context.

Do not create hidden unbounded timeouts.

Transport defaults may enforce upper limits, but must not ignore earlier caller deadlines.

---

## 47. Client Retries

Retries belong in a controlled client/resilience layer.

Retry only when:

```text
operation is safe/idempotent
error is retryable
attempts are bounded
backoff is bounded
deadline permits
```

Never blindly retry non-idempotent mutations.

---

## 48. Events

Events are contracts.

Recommended envelope:

```json
{
  "id": "opaque-event-id",
  "type": "application.created",
  "version": 1,
  "timestamp": "2026-10-01T10:20:30Z",
  "correlationId": "opaque-correlation-id",
  "organizationId": "opaque-org-id",
  "projectId": "opaque-project-id",
  "principal": {},
  "data": {}
}
```

Not every field is mandatory for every event, but envelope semantics must be standardized.

---

## 49. Event Naming

Use stable semantic names:

```text
resource.action
```

Examples:

```text
application.created
execution.started
execution.completed
approval.required
```

Do not put implementation/vendor names into generic domain event names.

---

## 50. Event Versioning

Event schema evolution follows the same compatibility discipline as APIs.

Consumers must not depend on producer-internal types.

Breaking event changes require a new event version/type contract.

---

## 51. Event Delivery Semantics

Every event channel documents:

```text
at-most-once / at-least-once expectations
ordering scope
deduplication key
retry behavior
dead-letter behavior where applicable
retention
```

Consumers must assume duplicates when at-least-once delivery is used.

---

## 52. Outbox Boundary

When business persistence and event publication must be atomic, prefer an outbox pattern.

The business service owns the transaction.

The transport producer publishes from the outbox.

---

## 53. Realtime Streams

WebSocket/SSE contracts must define:

```text
event type
payload version
ordering
heartbeat
reconnect behavior
resume cursor/last event ID when supported
authorization
backpressure/drop behavior
```

---

## 54. Execution Event Streams

Mission Runtime events should be structured, not raw log text only.

Example categories:

```text
execution.started
step.started
step.log
step.progress
step.output
approval.required
execution.failed
rollback.started
execution.completed
```

---

## 55. Nullability

Nullability is semantic.

Do not use `null` indiscriminately.

Contract documentation must distinguish:

```text
not provided
unknown
not applicable
explicitly cleared
```

where relevant.

---

## 56. Defaults

Server-side defaults are contract behavior.

Document them.

Do not require clients to duplicate every server default.

---

## 57. Experimental Fields

Experimental contract surface must be visibly marked and isolated.

Do not silently treat production public contracts as experimental to avoid compatibility responsibility.

---

## 58. Deprecation

Deprecation requires:

```text
deprecated marker
replacement path
timeline/policy
consumer communication
compatibility tests until removal
```

---

## 59. Breaking Change Examples

Breaking:

```text
removing field
renaming field
changing meaning
making optional input required
changing error semantics incompatibly
changing enum values in a way existing clients cannot handle
changing endpoint authorization semantics without migration planning
```

---

## 60. Compatible Change Examples

Normally compatible:

```text
new optional response field
new endpoint
new optional request field with server default
new enum value when clients follow unknown-value guidance
new event field ignored by older consumers
```

Compatibility must still be tested.

---

## 61. API Boundary Validation

Validate untrusted request structures at the transport boundary.

Business invariants remain in business services/entities.

Avoid:

```text
controller querying repository to validate ownership
validator calling another service
```

---

## 62. Contract Tests

Every public client/transport needs contract tests for meaningful behavior.

Examples:

```text
OpenAPI endpoint status/body semantics
gRPC method/status semantics
client/vN in-process and remote equivalence
event serialization compatibility
```

---

## 63. Golden Contract Fixtures

Small stable fixtures may be used to detect accidental wire-format changes.

Do not use giant snapshots as the only compatibility mechanism.

---

## 64. Consumer-Driven Assumptions

If a consumer relies on behavior not expressed in the contract, either:

```text
make the behavior explicit
or remove the hidden dependency
```

Do not allow tribal knowledge to become an API.

---

## 65. Frontend Contract Consumption

Frontend services consume the public backend API contract.

Frontend readers/mappers may adapt DTOs to view models.

Frontend components must not derive business semantics from undocumented backend fields.

---

## 66. Provider APIs

Generic product APIs expose capabilities, not vendor-specific branches.

Provider administration endpoints may expose provider metadata/configuration explicitly.

Business APIs should not require callers to pass `"redis"` or `"jenkins"` when the binding can be resolved from context.

---

## 67. Capability Discovery

When consumers need capability discovery, expose normalized capability identifiers.

Example:

```json
{
  "capabilities": [
    "CACHE_KV",
    "CI_RUN",
    "OBSERVABILITY_LOGS"
  ]
}
```

Do not force consumers to infer capabilities from provider brand.

---

## 68. Security-Sensitive Contracts

Authentication/authorization/token/secret contracts require threat review.

Never return raw secret material when a reference/handle is sufficient.

Sensitive values must be write-only or redacted where possible.

---

## 69. File Upload Contracts

Define:

```text
maximum size
allowed content types
streaming behavior
checksum if required
virus/malware scanning responsibility
object ownership
retention
```

Never trust filename or content type alone.

---

## 70. Download Contracts

Downloads should use explicit content type/disposition and authorization.

Pre-signed URLs, when used, must be short-lived and scoped.

---

## 71. Bulk Operations

Bulk APIs must define whether behavior is:

```text
atomic
best-effort
partial success
```

If partial, return per-item outcomes.

---

## 72. Rate Limits

When rate limits are part of the public contract, document:

```text
scope
limit
window
retry guidance
```

Do not expose infrastructure-specific limiter internals.

---

## 73. API Docs Are Generated from Source Contracts

Generated/reference docs are derived from OpenAPI/Protobuf.

Do not maintain a second contradictory manual API specification.

---

## 74. Examples

Examples must:

```text
compile/validate where practical
avoid real secrets
use opaque sample IDs
match current schema
```

---

## 75. Contract PR Checklist

```text
[ ] Owning service/domain is clear
[ ] Change is additive or explicitly breaking
[ ] Correct contract version used
[ ] OpenAPI/proto source updated
[ ] DTO/client updated
[ ] Generated files regenerated
[ ] No domain entity leaked directly
[ ] IDs remain opaque
[ ] Timestamps/durations explicit
[ ] Errors normalized
[ ] Pagination/filter/sort semantics explicit
[ ] Idempotency considered for mutation
[ ] Auth/authz documented
[ ] Contract tests added/updated
[ ] Compatibility check run
[ ] Frontend/other consumers reviewed
[ ] Examples/docs updated
```

---

## 76. Definition of Done

```text
[ ] Machine-readable contract is authoritative
[ ] Implementation matches contract
[ ] Public client matches contract
[ ] Errors/status mappings are tested
[ ] Breaking compatibility was not introduced accidentally
[ ] Generated code is fresh
[ ] Cross-service consumers use client/vN
[ ] Events/streams have explicit schema/version semantics
[ ] Security-sensitive fields are protected
[ ] Actual contract checks were run and reported
```

---

## 77. Phase 0 Deliverables

Phase 0 must prove:

```text
one OpenAPI v1 contract
one gRPC/protobuf v1 contract or a documented reason it is not yet needed
one normalized error schema
one pagination schema
one request/correlation convention
one public client/v1
one in-process client implementation
one contract compatibility check
one generated-code freshness check
one versioned event envelope
```

---

## 78. Agent Rules

Before changing a public API, an agent must inspect:

```text
docs/AGENT.md
this file
owning service public contract
DTOs
client/vN
controllers/mappers
contract tests
known consumers
```

Never change a public contract from only the implementation file.

---

## 79. Final Mandatory Rules

1. Contract first for externally consumed behavior.
2. OpenAPI owns HTTP shape.
3. Protobuf owns gRPC shape.
4. Generated code is read-only.
5. DTOs are not persistence/domain entities.
6. Contracts evolve additively inside a major version.
7. Breaking changes require deliberate versioning/migration.
8. IDs are opaque.
9. Time and duration units are explicit.
10. Errors are stable machine contracts.
11. Pagination is mandatory for unbounded lists.
12. Retryable mutations define idempotency.
13. Request/correlation context propagates.
14. Public clients hide transport choice.
15. Events and realtime payloads are versioned contracts.
16. Provider-neutral product APIs expose capabilities, not vendor branching.
17. Compatibility is enforced by tests/CI.
18. Security-sensitive contracts get explicit review.

---

## 80. Guiding Principle

```text
Stable Semantics
      +
Machine-Readable Contract
      +
Additive Evolution
      +
Normalized Errors
      +
Transport Independence
      =
Reliable Public Boundary
```

> **A Perihellion consumer should be able to integrate safely from the contract alone, without reading private implementation code.**

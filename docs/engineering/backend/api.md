# API, Controller, DTO, Validation and Mapper Standard

## 1. Transport boundary blueprint

```text
client/<resource>/<resource>.go                # stable typed @Client facade
internal/app/routes/http.go                    # service prefix + generated route registration
internal/app/controller/<resource>/<tier>/     # @Routes interface + implementation
internal/app/dto/<resource>/{request,response,params,queryparams,operations}.go
internal/app/validators/<resource>/{<resource>.go,impl/,helpers/}
internal/app/mappers/<resource>/{<resource>.go,impl/}
internal/app/services/<resource>/              # authorized domain use-case
internal/resources/api/openapi/v1/openapi.yaml # authoritative machine-readable contract
```

Create only DTO category files that contain real types. Group DTOs by **resource**, not separate top-level `dto/request`/`dto/response` packages. External consumers may use public aliases from `client/<resource>`; private DTOs are not imported from another module's `internal/app`.

## 2. HTTP responsibilities

- **Routes:** compose the generated registration functions; HTTP method and relative path are declared once on the controller interface with annotations.
- **Launcher gateway:** classify matched route, verify principal/tenant/workload, check Origin/CSRF and rate limits. Domain modules must not implement duplicate copies of these generic controls.
- **Shared middleware:** request/correlation/trace IDs, body limits and recovery.
- **Controller:** binding/structured HTTP result, invoke validator and use-case, verify authorization context, translate stable errors.
- **Validator:** pure technical checks and canonicalization, not DB lookups or response writing.
- **Service:** ownership/business invariant and change orchestration.
- **Repository:** GORM SQL and transaction-aware persistence.
- **Mapper:** transform real representation differences, enforce safe exported fields.

`public`, `protected`, and `restricted` may describe route tiers. Directory names and paths **do not enforce access control**. A temporary development principal must not become production authorization.

## 3. Request lifecycle

Use the common request contract with optional typed sections: `Params`, `Query`, `RequestBody`. Normalize deterministically, strict-decode unknown JSON fields, validate resource rules, obtain verified identity, authorize against trusted ownership, invoke operation, map and encode success/error response.

A generated `@RequestLine` client may dispatch in-process through the same REST path. `@Proxy` calls a supplied service method directly. The latter is a trusted internal boundary; do not mistakenly assume it passes through HTTP middleware. Document the security implications.

## 4. DTO rules

- Path/query/request/response models only; never expose GORM entities/transactions to clients.
- Use explicit JSON names, allowed fields and null handling. DTO validation should not copy every domain invariant into tags.
- PATCH requires **absent**, **null**, **value** states (including `null` to clear optional data) and explicit semantics for collection replacement vs merge.
- Response DTOs exclude secret tokens, password hashes, internal request fingerprints, unsafe PII and provider SDK objects.
- Use field names/casing that match OpenAPI and generated clients.
- No meaningless wrapper/duplicate type: a new type requires a different contract or policy.

## 5. Validation versus business policy

Validation covers syntax, length, normalization, known enum/operator and shape. Authorization checks a trusted actor, capability and resource ownership. Business rules cover allowed lifecycle transitions, duplicate entities, state preconditions and cross-aggregate consistency. Each has a distinct error class and location.

Do not repeat the same normalization in controller, service and repository. Direct service clients must receive equivalent invariant enforcement; an HTTP-only validator is not enough for every internal call.

## 6. Mapper rule

Use typed mappers for entity -> response, external DTO -> internal representation, or meaningful request -> command conversions. Do **not** use them for identical type copies, persistence, network calls or validation. Annotated mapping may generate straightforward transformations, but complex ownership/secret rules remain explicit, reviewed, and tested.

## 7. Status and errors

Use standard HTTP methods and idempotency semantics; map errors into the established structured envelope, preserving correlation metadata without leaking internals. Explicitly document created vs existing result (201 vs 200 if implemented), invalid input (400), denied access (401/403), missing resource (404), conflict (409), rate limit (429), unsupported capability (501 when appropriate), unavailable dependency (503). Follow the central API contract guideline for detailed governance.

## 8. API evolution

OpenAPI, route annotations, client method signature, DTO shape and product walkthroughs must be updated in the same commit. Breaking changes require versioning/deprecation strategy; never rename a route or event only to make folders symmetric. Include handling for error cases and access tiers.

## 9. HTTP test matrix

Test malformed/unknown JSON properties, absent vs null PATCH, invalid IDs, unauthorized/forbidden actors, resource isolation across tenants, correct status+body envelope, idempotent retries, conflict, large body, cancellation, and generated `@RequestLine` compatibility with actual routes. A route test that bypasses IAM is not a security test.

## 10. Annotated route limits

Route-specific HTTP throttling is owned by the controller contract, for example `@POST /login` followed by `@RouteLimit 12`. The shared generator validates and translates that declaration into an exact method-and-route policy; each service exposes it from bootstrap to Launcher. Hardcoded endpoint switches in the Launcher gateway are prohibited. An absent annotation does not override other applicable platform security controls. See `annotation.md` and `gateway.md`.

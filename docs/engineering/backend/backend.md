# Backend Engineering — Master Standard

> Scope: every backend service and shared library. This file describes architecture, responsibilities, boundaries and engineering principles. It deliberately contains **no concrete product service examples**, API paths, domain entities or business workflows.

## 1. Authority and documentation hierarchy

1. This document governs cross-cutting backend engineering rules.
2. [searchDetails.md](searchDetails.md) governs searching, filtering, sorting and cursor pagination.
3. [annotation.md](annotation.md) governs code-generation annotations.
4. Each module's `docs/engineering/` records implementation-specific decisions, deviations and extension points, linking here instead of copying the rules.
5. Each service's `docs/product/` describes business behavior, API contracts, examples, permissions, workflows and reproducible verification steps. Libraries have engineering docs, not product docs unless they expose an actual product surface.
6. Real code, tests, migrations and API contracts must agree with documentation. A proposal is not a shipped feature; mark unimplemented capabilities clearly.

## 2. Design philosophy

- **SRP (Single Responsibility):** every package, type and function has one cohesive reason to change.
- **DRY:** share stable knowledge, not coincidentally similar code; avoid prematurely extracting one-off abstractions.
- **KISS:** choose the simplest implementation that satisfies verified requirements.
- **YAGNI:** do not build speculative capabilities; create extension seams only where useful.
- **SOLID:** SRP, open/closed through composition, substitutability of contracts, small interfaces and dependency inversion at boundaries.
- **Separation of concerns:** transport, orchestration, domain behavior and persistence must not leak into each other.
- **High cohesion, low coupling:** keep related behavior together; dependencies point to smaller stable contracts.
- **Composition over inheritance:** use interfaces and composed dependencies rather than pseudo-inheritance.
- **Explicit over implicit:** make ownership, transactions, side effects, defaults and errors visible.
- **Fail fast at boundaries:** validate untrusted data early; do not silently repair malformed input.
- **Least privilege & defense in depth:** authorization at the use-case boundary, never solely at routing.
- **Idempotency:** repeated identical commands should not duplicate externally observable effects.
- **Backward compatibility:** treat APIs, schemas and published events as versioned contracts.

## 3. Architectural layers and dependency direction

```text
Transport -> Controller -> Application Service -> Repository Interface
                  |                |                  |
                  v                v                  v
             Validation       Domain rules      Repository impl -> Database
                                   |
                                   v
                          Event publisher / Outbox
```

The graph is conceptual, not a mandate for an extra package per box. Dependencies may point from outer layers to stable shared contracts. No repository imports a controller, and no domain rule depends on HTTP or message broker payload types.

## 4. Package/folder responsibility matrix

| Folder / concept | Owns | Must not own |
|---|---|---|
| `routes/` | Route registration, grouping, transport middleware wiring, endpoint-to-handler mapping | Validation rules, SQL, business decisions |
| `controller/` | Decode/bind transport requests, invoke application operation, map result to response and status | Domain workflows, transactions, direct DB access |
| `dto/request/` | Input contract and field optionality/patch semantics | Persistence annotations and business logic |
| `dto/response/` | Output contract, compatibility and serialization | Entity persistence or internal secrets |
| `dto/operations/` | Internal operation inputs/outputs when they reduce coupling | Duplicate request DTOs with no reason |
| `entity/` | Domain state, invariants and persistence model where intentionally combined | HTTP concerns and message transport |
| `validation/ or validators/` | Pure input validation, normalization contracts, domain-independent constraints | Side effects, DB writes, authorization decisions |
| `mappers/` | Explicit shape conversion between DTOs, entities and events | IO, business rules, silent data loss |
| `services/` | Use-case orchestration, authorization decisions, transactions, domain operations and coordination | HTTP status codes, SQL syntax, broker protocol |
| `repository/` | Persistence interfaces and aggregate-oriented operations | HTTP concerns and external API calls |
| `repository/impl/` | GORM/SQL query execution, constraints, transaction participation | Transport objects, business workflows |
| `repository/helpers/` | Private reusable query/scan/persistence helpers | Public API or cross-module policy |
| `events/producer/` | Event construction boundary, publishing integration and delivery intent | Request handling, database business decisions |
| `events/consumer/` | Subscription setup, deserialize/validate, idempotent dispatch, retry/ack policy | Arbitrary cross-service DB mutation |
| `events/contracts/` | Schemas, versions, event names, compatibility rules | Direct persistence implementation |
| `externalService/` | Incoming/outgoing external provider integration boundary | Core domain ownership |
| `adapters/` | Provider-specific protocol translation | Provider-specific rules leaking into service logic |
| `middleware/` | Cross-cutting HTTP/gRPC transport controls | Domain authorization as the only check |
| `pipeline/` | Deterministic boundary processing stages | Business use cases |
| `constants/` | Truly shared stable typed constants / keys | Unrelated dumping ground, magic configuration |
| `errorcode/` | Stable typed error catalog and safe error mapping | Raw HTTP response construction |
| `config/` | Typed configuration parsing and validation | Business behavior |
| `helpers/` | Narrow private reusable functions; prefer placement near owner | Hidden global utility grab bag |
| `utils/` | Small generic deterministic functions only | Service-coupled behavior |
| `db/ or migrations/` | Versioned schema evolution, rollback policy, indexes and constraints | Runtime business orchestration |
| `tests/` | Unit, integration, contract and end-to-end evidence | Production-only logic |
| `docs/engineering/` | Module design, dependencies, invariants, extension/configuration, operational behavior | Repetition of global standards |
| `docs/product/` | User journeys, workflows, APIs, examples and reproduction steps | General Go style or infrastructure architecture |

**Do not create every folder automatically.** Only create a folder when the module has that responsibility. Resource-specific subfolders are appropriate when they improve navigation; avoid mirrored empty folder trees.

## 5. Request lifecycle

1. Edge middleware establishes request identity, trace/correlation context, size limits and authentication.
2. Pipeline performs deterministic normalization and safe decode in the agreed order; parsing/validation errors are reported consistently.
3. Validator checks shape and semantic constraints that do not require side effects.
4. Controller adapts transport into a use-case input.
5. Application service applies authorization, business invariants and orchestration.
6. Repository implementation performs persistence under an explicit transaction boundary.
7. Domain changes and any required outbox record commit atomically.
8. Controller maps outcome to versioned response and stable errors.
9. Async publisher delivers committed events; consumer validates, deduplicates, handles and acknowledges.

Specify whether normalization happens before decoding for each content type; do not mutate arbitrary raw payloads in a way that changes signature validation or canonical value semantics. **A controller never substitutes for the use-case layer.**

## 6. Transaction and persistence guarantees

**ACID:** Atomicity (all-or-none transaction), Consistency (declared invariants and constraints survive), Isolation (concurrent transactions follow the chosen isolation semantics), Durability (committed state persists according to DB guarantees).

- A transaction has one owner: usually an application operation or unit of work.
- Repository methods accept/participate in caller context and transaction; they do not silently commit independent parts of one use case.
- Keep transactions short; do not block on network calls or event delivery while holding DB locks.
- Use DB constraints to enforce uniqueness and referential rules; app validation alone is insufficient under concurrency.
- Explain isolation level, lock order, deadlock retry and optimistic concurrency when pertinent.
- Use **transactional outbox** when local data and event-intent require atomicity; broker publish is after commit.
- Cross-service distributed workflows need explicit saga/compensation and eventual consistency, not implied cross-DB ACID.
- Migrations must be versioned, ordered, reviewable and compatible with rolling deployments.

## 7. Events, messaging and delivery semantics

Event contracts require an event ID, type, version, aggregate identifier, occurrence time, correlation/causation metadata where relevant and a typed payload. Keep producer and consumer responsibilities separate.

- **Producer:** builds a versioned fact, schedules durable delivery (outbox when needed), emits telemetry and handles publisher retries.
- **Consumer:** validates contract/version, enforces idempotency by event ID or business key, performs one bounded operation, acknowledges only on success, and classifies transient/permanent errors.
- Define stream, topic/subject, retention, replay, ordering key, delivery policy, backoff, DLQ/quarantine and poison-message handling per integration.
- Expect at-least-once delivery unless a stronger mechanism is demonstrated; **exactly once** is never assumed from a broker name.
- Schema evolution is backward/forward compatible within policy; breaking changes require a new version.
- Event names describe completed facts, not HTTP commands.
- Never publish PII/secrets by default, and preserve tenant boundaries and authorization for event-triggered operations.

## 8. API contracts and security

- Version externally visible contracts; use consistent naming, status codes, pagination and stable machine-readable errors.
- Request bodies distinguish omitted, null and concrete values for PATCH; document field-removal rules.
- Use explicit idempotency keys and request hash checks for retry-sensitive operations.
- Authentication proves identity; authorization evaluates actor, tenant, resource and action at the use-case boundary.
- Sanitize inputs, parameterize SQL, constrain query complexity, avoid logging secrets and personally identifiable information.
- Define token propagation, request timeouts, cancellation, bounded retries, circuit breakers and rate limits.
- Integrations receive minimum necessary credentials; store secrets outside source and rotate safely.

## 9. Reliability and observability

Use structured logs with correlation/trace identifiers, actor/tenant identifiers only where safe, operation, duration and error code. Add metrics for latency, error rates, retries, queue lag, DB pool exhaustion and outbox backlog. Propagate distributed traces. Define SLOs, alert thresholds, graceful shutdown, liveness/readiness checks and overload behavior. Do not expose sensitive values in logs or metrics labels.

## 10. Concurrency, performance and correctness

- Assume concurrent requests and event redelivery; specify locking and race behavior.
- Respect context cancellation and deadlines across DB/provider calls.
- Bound worker pools, batch sizes, query limits and memory usage.
- Prevent N+1 queries, unbounded scans and unindexed filtering; explain query plans and indexing decisions.
- Cache only with explicit TTL, key ownership, invalidation and consistency requirements.
- Prefer correctness and measurable bottlenecks to speculative optimization.

## 11. Patterns and when to apply them

| Pattern | Use | Avoid |
|---|---|---|
| Repository | Isolate storage/query mechanics and enable contract tests | A trivial wrapper for every ORM call |
| Unit of Work | Coordinate atomic writes across repositories | Transactions hidden deep in helpers |
| Adapter | Translate external provider SDK/shape into internal contract | Polluting core types with vendor details |
| Strategy | Select among interchangeable policies/providers | One-off if/else unnecessarily abstracted |
| Factory | Complex validated construction | Boilerplate construction of simple structs |
| Builder | Many optional construction parameters with invariants | Simple required-field structs |
| Dependency injection | Wire dependencies explicitly for isolation | Service locator/global mutable state |
| Observer / Pub-sub | Decoupled asynchronous event notification | Work requiring synchronous consistency |
| Outbox | Atomic state change + durable delivery intent | Assuming a publish API is transactional with DB |
| Saga | Multi-service long-running coordination | Claiming full atomicity across independent DBs |
| CQRS | Read/write workload and models diverge significantly | Splitting straightforward CRUD by habit |
| Decorator | Add cross-cutting behavior to a stable interface | Hiding ordering-sensitive business rules |
| Middleware / Pipeline | Shared boundary behavior in deliberate order | Domain logic dependent on HTTP middleware |
| Circuit breaker | Protect calls to unstable remote dependencies | Local pure in-process functions |
| Bulkhead | Isolate resource pools/failure domains | Unlimited goroutines or shared mutable bottlenecks |

Describe tradeoffs and failure modes whenever adopting a pattern.

## 12. Go code conventions

Keep function and interface method signatures on one source line when practical, as a house preference; avoid sacrificing readability for exceptionally long signatures. Use small interfaces defined by their consumers, explicit errors, typed enums, context as first argument for I/O operations, Go naming idioms and direct composition. Generated files must be deterministic, clearly marked and kept in designated generated-code locations. Do not edit generated output manually. Do not treat annotations as runtime frameworks.

## 13. Required quality gates

Every change should address: formatting/lint, tests appropriate to the layer, migration compatibility, API/event contract versioning, concurrency/idempotency cases, security review, observable failure modes and documentation updates. New packages require a documented owner and boundary. Product features must have reproducible API examples under the owning service's `docs/product/`.

## 14. Documentation anti-patterns and maintenance

Never duplicate this master standard in a service. Avoid stale snapshot dumps, promises phrased as shipped behavior, speculative directory scaffolding, service-specific examples in this file and undocumented deviations. Prefer links, tables of ownership and tested examples. Changes to shared conventions update this file first; implementation-specific decisions stay next to the module.

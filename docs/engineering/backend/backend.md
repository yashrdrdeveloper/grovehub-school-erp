# Backend Engineering Standard

**Status:** Normative • **Scope:** every Go backend module and shared library • **Owner:** Backend engineering

This file is the **generic source of truth** for designing and maintaining backend packages. It deliberately contains **no specific service, vendor, event, endpoint, table, or product workflow**. Module-specific implementation belongs in the module's `docs/engineering/`; business API execution and expected results belong in `docs/product/`.

The words **MUST**, **MUST NOT**, and **SHOULD** express enforceable rules. Existing external contracts, already-applied database migrations, and production data take precedence over a purely cosmetic folder move; document exceptions in the owning module.

## 1. Information architecture and ownership

| Location | Responsibility | Forbidden |
|---|---|---|
| `docs/engineering/backend/backend.md` | This architecture/SRP/quality master | Service-specific details |
| `docs/engineering/backend/<topic>.md` | Reusable technical guidance | Individual product workflows |
| `backend/services/<module>/docs/engineering/` | Module-specific implementation, deviations, contracts, dependencies, runbooks | General platform-wide rules duplicated from here |
| `backend/services/<module>/docs/product/` | Use cases, actor journeys, API sequences, request/response examples, acceptance criteria | Internal algorithm blueprints |
| `backend/libs/<library>/docs/engineering/` | Shared capability contracts, use conditions, invariant/security guarantees | Domain- or vendor-specific product rules |
| `backend/services/<module>/internal/resources/api/openapi/` | Machine-readable canonical API schema | Prose-only undocumented contract |
| `backend/services/<module>/internal/resources/migrations/` | Immutable, ordered schema changes | Unreviewed generated schema mutation |

When documentation and code differ, first establish what is actually deployed, then **correct the owner-local documentation** and any generic rule only if the latter is wrong. A plan or annotation inventory is not proof that code generation, authorization, transactions, or other runtime behavior exists.

## 2. Canonical module blueprint

Use a **resource-first layout inside a stable technical layer**. Place only folders that carry real responsibilities.

```text
backend/
├── launcher/                                 # process composition + gateway access policy
├── libs/
│   ├── core/                                 # domain-neutral reusable capabilities
│   └── <technical-library>/                  # separately owned capability contracts
└── services/
    └── <module>/
        ├── <module>.go                       # exported New, Migrate, Module only as needed
        ├── client/
        │   └── <resource>/<resource>.go      # stable intermodule facade, aliases, @Client
        ├── internal/
        │   ├── app/
        │   │   ├── bootstrap/
        │   │   │   ├── bootstrap.go          # builder/runtime interfaces
        │   │   │   └── impl/
        │   │   │       ├── build.go          # minimal dependency composition
        │   │   │       ├── runtime.go        # module interface implementation, Close
        │   │   │       ├── client.go         # internal clients
        │   │   │       ├── routes.go         # HTTP wiring, only for HTTP modules
        │   │   │       ├── repository.go     # DB wiring
        │   │   │       ├── services.go       # use-case composition
        │   │   │       ├── events.go         # producer/consumer composition
        │   │   │       ├── external_services.go # adapter registration
        │   │   │       └── workers.go        # background lifecycle, if needed
        │   │   ├── config/{config.go,impl/}  # typed private config and loading
        │   │   ├── constants/{error,keys,events}/
        │   │   ├── dto/<resource>/
        │   │   │   ├── request.go
        │   │   │   ├── response.go
        │   │   │   ├── params.go             # if actual path params exist
        │   │   │   ├── queryparams.go        # if actual query params exist
        │   │   │   └── operations.go         # if distinct operation values exist
        │   │   ├── entity/<entity>.go
        │   │   ├── controller/<resource>/
        │   │   │   ├── <access>/<resource>.go # @Routes / method contract
        │   │   │   └── <access>/impl/<resource>.go
        │   │   ├── routes/http.go            # compose GENERATED annotated route groups
        │   │   ├── validators/<resource>/
        │   │   │   ├── <resource>.go          # typed rule interface
        │   │   │   ├── impl/<resource>.go
        │   │   │   └── helpers/               # validation-only helpers
        │   │   ├── mappers/<resource>/
        │   │   │   ├── <resource>.go          # representation conversions
        │   │   │   └── impl/<resource>.go
        │   │   ├── repository/<resource>/
        │   │   │   ├── <resource>.go          # @Repository interface
        │   │   │   ├── impl/<resource>.go     # GORM and SQL
        │   │   │   └── helpers/search/        # trusted search config as needed
        │   │   ├── services/<resource>/
        │   │   │   ├── <resource>.go          # use-case interface
        │   │   │   ├── impl/<resource>.go     # @Service dependency composition
        │   │   │   ├── helpers/               # used by >=2 operations
        │   │   │   └── operations/
        │   │   │       └── <operation>/
        │   │   │           ├── <operation>.go # single operation entry/chain
        │   │   │           ├── flows/          # typed chain state, if multistep
        │   │   │           ├── helpers/        # operation-private pure logic
        │   │   │           └── strategies/
        │   │   │               ├── strategy.go
        │   │   │               └── impl/      # actual operation algorithms
        │   │   ├── events/
        │   │   │   ├── producers/<aggregate>/
        │   │   │   └── consumers/<aggregate>/
        │   │   ├── providers/                 # generic provider engine if needed
        │   │   ├── externalService/<capability>/<vendor>/
        │   │   └── client/<resource>/impl/     # generated client/proxy adapter
        │   ├── resources/
        │   │   ├── resources.go               # embed resources
        │   │   ├── migrations/
        │   │   └── api/openapi/v1/
        │   └── tests/
        │       ├── unit-test/                 # mirrors source ownership
        │       ├── mocks/
        │       └── fixtures/
        └── docs/
            ├── engineering/{README.md,architecture.md}
            └── product/{README.md,api_workflows.md}
```

**Not a scaffolding checklist.** Absence of an HTTP API means no controller/routes folder. No persistence means no repository/entity/migration. A single-step operation does not need `flows/` or a chain. Never create `impl/` merely to wrap a one-line pure function without a genuine boundary.

## 3. Folder SRP: what to do and what not to do

| Folder / file | Exactly one reason to change | Do | MUST NOT |
|---|---|---|---|
| Module root | Stable module lifecycle interface | `New`, `Migrate`, public aliases | SQL, provider SDK calls, HTTP method logic |
| `client/<resource>/` | Cross-module public API | Typed contracts, public DTO aliases, `@Client` | Import another module's `internal/`, leak GORM/entity/SDK |
| `bootstrap/` | Object graph lifecycle | Define `Bootstrap`/`Runtime` | Business policy |
| `bootstrap/impl/build.go` | Minimal dependency assembly | Validate dependencies and delegate typed wiring | Controllers, operations, workers, or SQL logic inline |
| `bootstrap/impl/runtime.go` | Module Runtime implementation | RegisterHTTP, expose typed clients, VerifyAccess and Close through injected interfaces | Business policy, direct service calls or concrete controller assembly |
| `bootstrap/impl/client.go` | Client composition | Supply generated @RequestLine client and in-process router dependencies | Direct service proxies for HTTP-exposed methods, business validations |
| `bootstrap/impl/routes.go` | Controller + route composition | Register actual handlers | Inline handlers |
| `bootstrap/impl/repository.go` | Persistence DI | Construct repository, mapper, validator | SQL query bodies |
| `bootstrap/impl/services.go` | Use-case DI | Compose services from typed dependencies | Vendor switches in operations |
| `bootstrap/impl/events.go` | Event subscriptions/producers | Bind producer and consumer contracts | Keep raw broker SDK logic |
| `bootstrap/impl/external_services.go` | Vendor adapters | Register concrete SDK implementations | Domain rules |
| `bootstrap/impl/workers.go` | Worker lifecycle | Cancel, join, expose failure/health | Detached, untracked goroutines |
| `services/<capability>/` | Identity/security and other reusable module-owned use cases | Typed capability contracts and impl packages | Parallel feature-root `security/` implementing business use cases |
| `config/` | Private runtime config | Load/validate env, redact secrets | Hardcoded production credentials/endpoints |
| `constants/error/` | Error identity | Stable codes, metadata | HTTP implementation details |
| `constants/keys/` | Stable persistence/header identifiers | Table constraints, agreed keys | Business decisions |
| `dto/<resource>/` | API or command representation | Requests, responses, PATCH presence | SQL, side effects, credential decryption |
| `entity/` | Persisted state schema | GORM mapping, timestamps, keys | HTTP, business workflows, auto-fetch |
| `controller/` | HTTP transport | Decode, invoke rules, auth, invoke service, encode | GORM, event dispatch, vendor calls |
| `routes/` | HTTP routing | Register method/path/tier | Authorization by folder naming alone |
| `validators/` | Input policy | Normalize, reject bad input, pure rule checks | Queries, transactions, response writing |
| `mappers/` | Representation change | Explicit field mapping, null/PII policy | Fetching, persistence, validation duplication |
| `repository/` | Data access | GORM, trusted queries, optimistic version, constraint mapping | HTTP, NATS, remote API calls |
| `services/<resource>/` | Business use-case contract | Expose operation signatures | Vendor-specific transport types |
| `services/<resource>/impl/` | Operation composition | Embed/construct operations | Monolithic logic, dispatch by giant switch |
| `services/<resource>/operations/<op>/` | One business operation | Orchestrate steps and transaction | Competing resource APIs |
| `operations/<op>/flows/` | Typed multi-step state | Request, accumulated state, handoff | Persistence entity dumping ground |
| `operations/<op>/strategies/` | Replaceable operation behavior | Typed interface, `impl/` algorithm | Unbounded string/reflection dispatch |
| `operations/<op>/helpers/` | Pure helper needed only in one operation | Parse, compare, derive values | Reused cross-operation logic |
| `services/<resource>/helpers/` | Helpers shared by multiple operations | Shared computation, typed domain helpers | Generic global utility dump |
| `events/producers/` | Outbound event meaning | Stable schema, envelope, emit rules | Broker implementation |
| `events/consumers/` | Inbound event response | Decode, deduplicate, invoke authorized use case | Owning the upstream event schema |
| `externalService/` | Third-party integration | HTTP clients/SDK, rate limits, auth tokens | Application/domain policy |
| `providers/` | Generic provider discovery/execution | Capability-neutral contracts, registry | `mail.*`, payment-specific, vendor specifics |
| `internal/resources/` | Immutable deployable assets | Embed versioned SQL/OpenAPI | Runtime policy decisions |
| `internal/tests/` | Confidence and contracts | Mocks, scenario tests, exact regression cases | Shipping test-only shims in production |
| `docs/engineering/` | How this module works | Architecture, diagrams, invariants, exceptions | Duplicating this entire master |
| `docs/product/` | What user or caller can do | API journey, preconditions, examples, expected results | Hidden internal implementation dependency |

**Placement test:** Determine the *owner* (resource/capability), the *change trigger* (reason to modify), and the *side-effect boundary*. Put code in the lowest owner-specific folder satisfying all three.

## 4. Responsibility flow

### HTTP command

```text
request -> middleware [trace, identity, limits]
       -> routes [method/path only]
       -> controller [decode, authorization, REST result]
       -> validator [canonical input]
       -> service / operation [business rules and coordination]
       -> repository [transaction-aware GORM]
       -> mapper [if representation actually changes]
       -> response DTO / common HTTP envelope
```

A validator is not authorization. A mapper is not a business operation. A controller cannot validate authorization by trusting a user-supplied owner ID. Operation output must not expose stored credentials or an ORM entity accidentally.

### External dependency

```text
service operation -> typed capability contract -> provider registry
                  -> externalService/<capability>/<vendor>/adapter
                  -> external API
```

Transport DTOs live at their owner boundary. Never import a concrete vendor adapter from unrelated business modules.

### Business event

```text
successful domain mutation -> aggregate-owned producer
     -> (transactional outbox if necessary) -> Infrastructure dispatcher
     -> broker -> aggregate-owned consumer -> deduplicated use case
```

A producer owns business meaning and eligibility; Infrastructure owns broker delivery/retries, not domain semantics. A consumer may update its own state, never another module's private tables.

## 5. Domain modeling and cohesion

- **Entity:** passive persisted state, primary/foreign keys, schema mapping, locking version, lifecycle timestamps. Do not embed service or controller dependencies.
- **Aggregate:** transactional consistency boundary. It may span multiple tables without requiring the same Go struct; enforce invariant transitions in its owning operation.
- **Value object:** immutable meaning (validated identifiers, money, dates) where semantics justify a type. Do not create wrappers only to rename primitives.
- **DTO:** transport or operation boundary, not database storage. Resource folder groups `request`, `response`, `params`, `queryparams`, and genuinely distinct operation types.
- **Repository:** one data-owner interface at the resource boundary. Search/relation queries remain trusted and parameterized.
- **Service:** only meaningful business use cases. One public service per resource/capability, with operation packages when more than one operation or flow exists.
- **Provider:** capability-neutral interface + specialized capability contract + concrete external adapter. Optional capabilities must not be forced into all providers.

## 6. SOLID and package-level architecture

**SRP:** each package/component has one owner and change trigger; constructors do only wiring. **OCP:** introduce a new operation strategy or vendor adapter without editing a central switch in domain code. **LSP:** substituting an implementation must preserve authorization, errors, side effects, and timing guarantees. **ISP:** split broad interfaces by consumer needs, e.g. discovery vs executable actions. **DIP:** business operations depend on typed local interfaces; concrete DB, broker, and vendor SDKs are wired at bootstrap.

Use the **dependency rule**: internal algorithm packages may depend inward on domain DTOs and stable shared capabilities, never outward on launcher/controller or concrete vendor packages. Import cycles indicate misplaced ownership. Keep interfaces small; a Go interface is warranted when there is a real consumer-facing boundary, testing seam, or multiple implementations—not because Java has one.

## 7. Additional quality principles

- **DRY:** unify *identical stable knowledge*, not coincidentally similar code. Reuse cross-service mechanics under a named shared capability; keep domain policy separate. Avoid generic `helpers/` or global `utils/` buckets.
- **KISS:** choose the smallest explicit solution that enforces real invariants; single-step operations do not need artificial chains/factories.
- **YAGNI:** future providers/modules are not an excuse to implement empty interfaces, routes, or tables today.
- **Separation of concerns:** HTTP, business decisions, persistence, event delivery, provider I/O, and observability have distinct owners.
- **High cohesion / low coupling:** organize resources inside layers; public contracts are smaller and more stable than private implementations.
- **Explicitness over magic:** compile-time annotations are metadata/generation only. Runtime transactions, authorization, retries, idempotency, and dependency registration must be visible and testable.
- **Fail fast:** reject missing config and invalid adapters during startup; fail closed for authorization. Do not hide errors in log-only fallbacks.
- **Prefer composition:** Go embedding of operation implementations can form a service facade; avoid inheritance-like hierarchies and unnecessary indirection.
- **Consistency:** same naming, nesting, import boundaries, and documentation scope for every *real* capability; not identical sets of empty directories.

## 8. Data and transaction rules (ACID)

- **Atomicity:** all changes to one aggregate, associated relations, idempotency state, and required outbox events commit together or roll back together.
- **Consistency:** enforce invariants at both valid input boundaries and DB constraints. Duplicate normalization or broken references must not silently pass.
- **Isolation:** choose correct transaction level, locks, constraint handling and optimistic versions for contention; a transaction alone does not prevent lost updates.
- **Durability:** successful commit must survive process failure under the configured database durability. A broker publish performed *after* commit is a separate guarantee unless a durable outbox bridges it.
- **Unit of Work:** a single transaction-bound context is passed through all repositories. Use explicit `transaction.UnitOfWork` for multi-repository operations.
- **Idempotency:** the key is scoped, validated, stored with a request fingerprint, and tied to a stable result. Replays do not repeat side effects; concurrent conflicts are handled.
- **Outbox:** insert in the same PostgreSQL transaction as domain writes when reliable delivery is required. Asynchronous delivery is at-least-once; consumers persist deduplication by event ID. A dispatcher's delivery-state cleanup is **not** audit history.
- **Migrations:** forward-only when deployed, versioned, reviewed, rollback/data-backfill strategy documented. Never edit an applied migration to make CI green.

See [persistence.md](persistence.md) and [events.md](events.md) for concrete generic mechanisms.

## 9. Patterns: when and when NOT to use

| Pattern | Use for | Avoid when |
|---|---|---|
| **Strategy** | Different interchangeable algorithms/provider actions; typed interface selected in bootstrap | Only one stable inline pure function |
| **Factory** | Construct a valid group of related dependencies, enforce nontrivial initialization | One direct constructor suffices |
| **Abstract Factory** | Entire families of adapters with shared lifecycle/configuration | No meaningful family distinction |
| **Adapter** | Translate vendor SDK/HTTP types into application contracts | Merely renaming identical types |
| **Registry** | Find registered implementations by stable identity/capability | Static one-provider dependency |
| **Facade** | Expose a small public client contract over private implementation | Exporting every internal helper |
| **Command** | Typed capability operation execution, validation and result | Giant string switch without typed policy |
| **Chain of Responsibility / Pipeline** | Multi-step ordered prepare/load/persist flow with clear short-circuiting | Single straightforward operation |
| **Unit of Work** | Atomic multi-repository database changes | Remote distributed transactions claimed by local DB |
| **Repository** | Resource persistence policy and query boundary | Generic `Save(any)` for all domains |
| **Transactional Outbox** | Atomic business state and durable event intent | Best-effort ephemeral telemetry |
| **Decorator** | Cross-cutting behavior with stable contract and no hidden ordering surprises | Obscuring business invariants |
| **Specification / Policy** | Typed reusable predicates with independently tested semantics | Arbitrary SQL supplied by callers |
| **CQRS** | Measured need for different read/write projections and latency models | Default CRUD complexity |
| **Saga** | Long-lived distributed business coordination with compensation | Local PostgreSQL transaction |
| **Observer / Pub-sub** | Decoupled asynchronous domain notifications | Synchronous commands that require immediate consistency |

A pattern is acceptable only when its owner, failure model, tests, and benefit are stated. See [patterns.md](patterns.md).

## 10. Input, validation, mapping, and search

Transport pipeline: normalize -> decode into typed request -> validate technical rules -> authorize with verified context -> business invariant checks -> operation. When normalization is performed before decoding, it must not silently hide malformed JSON or unknown fields; both representations require a consistent policy.

PATCH distinguishes **absent**, **explicit null**, and **supplied value**. Never treat a Go pointer alone as full presence tracking. Preserve nullability, field ownership, idempotency and optimistic concurrency throughout. The common `libs/core/request` contract includes optional `Params`, `Query`, and `RequestBody`, only when needed.

Mappers transform between meaningful representations. Do not create mappers that simply pass the same struct through. Search uses validated public requests -> resource-owned allowlist/config -> generic cursor/keyset query -> hydrated domain response. No raw field names or untrusted SQL expressions.

See [api.md](api.md) and [searchDetails.md](searchDetails.md).

## 11. Security, privacy, and multi-tenancy

Use trusted principal context, authorization policy, and scoped repository predicates; never treat route naming (`public`/`protected`/`restricted`) or annotations as an enforcement mechanism. Distinguish authentication, authorization, tenancy, and ownership. Reject unknown/unauthorized resources without leaking their existence as required by policy.

Use credential encryption at rest, least-privilege provider scopes, explicit secret rotation, constant-time comparison where relevant, sensitive-field redaction, and strict rate/size limits on external operations. Do not log raw message payloads, credentials, personal data, or token values. External callbacks must validate signatures/audience/state/replay before mutation; never trust headers that clients can spoof. Follow platform security standards for deployment, not DIY crypto.

## 12. Messaging and resilience

Stable event type, topic and schema version; domain envelope includes event ID, aggregate ID, occurred-at and correlation metadata. The aggregate version is distinct from schema version. Consumers are idempotent; duplicate/redelivery, out-of-order events and poison messages must have tested behavior. Lease/retry/backoff/DLQ policies are owned by Infrastructure. Provider-side HTTP retries need timeouts, cancelable contexts, bounded concurrency, and error classification. No infinite retry loop inside a user-facing request.

Broker storage and database outbox are independent. `Persist:false` is direct publication; `Persist:true` enqueues durable delivery intent when called in a transaction. It does **not** imply permanent audit retention. See [events.md](events.md).

## 13. Error model and logging

Declare stable, namespaced errors at the owning module/capability; map DB/remote errors once at their boundaries. Operations return typed errors, controllers translate them to approved HTTP envelopes/codes. Unexpected errors produce safe server errors with correlation metadata. No swallowed failures, panic for expected validation, or leakage of SQL/provider details to clients.

Use structured, context-aware logging, W3C trace propagation, bounded metric labels, and stable service identity. Logs must help identify *which operation* failed without exposing contents.

## 14. Configuration and lifecycle

Load private configuration from environment/secret sources; validate required values before server readiness. Development defaults belong in explicit local startup scripts, not hardcoded production endpoints or credentials. Bootstrap constructs repositories -> providers -> operations -> clients -> handlers -> workers. Track all background goroutines, cancel contexts and join work on shutdown before closing dependent brokers or databases. A non-nil dependency contract should fail at construction, not be silently replaced with a production no-op.

## 15. Go source and naming style

1. **One-line Go function/method signatures** when readable and required by the module style; split only truly unmanageable signatures deliberately. Format every file with `gofmt`. CI MUST reject unformatted Go. No dense one-line function bodies or packed switch cases.
2. **Blank line between every top-level type, const, var, and func declaration.** Also separate successive method definitions. Attach GoDoc/annotation comments directly to the following declaration.
3. Naming: `<resource>.go`, `operations/<verb>/<verb>.go`, `strategies/strategy.go`, `strategies/impl/<verb>.go`, `helpers/<meaning>.go`; no repetitive `zz_thing` in source.
4. Prefer domain-significant names over generic `manager`, `processor`, `common`, `handler2` or `misc`.
5. Don't create `index.go` barrels for DTO, constants, helpers, or types. Keep modules importable through coherent facades.
6. Use Go interface conformance assertions when useful; return errors and context explicitly; avoid shadowed context, global mutable registries and unnecessary reflection.
7. Generated files stay in service/library `.generated/` and are applied by the approved Go overlay. Do not commit generated noise or add direct imports of generated internals.
8. No commented-out code, stale compatibility shim without owner/deletion plan, or misleading "implemented" comment.

## 16. Tests, docs and delivery gates

- Place unit tests in mirrored `internal/tests/unit-test/` or approved owner-local packages where needed for Go `internal` visibility; prefer Ginkgo/Gomega for repository conventions; Go table tests are valid when clearer.
- Test success, invalid input, nil dependencies, authorization denial, transaction rollback, idempotent replay, duplicate delivery, and cancellation.
- Mock at the boundary; database repositories need realistic query tests; provider contracts need adapter tests; full broker/database behavior needs integration or sanity tests.
- Mocking a transaction is not proof of ACID; mocking a publisher is not proof of outbox delivery.
- Update owner engineering and product docs in the **same commit** as changes to behavior or API.
- Regenerate and validate annotations, build with overlay, run unit, integration/sanity, schema and API checks as appropriate. Fix warning/compile failures before declaring completion.
- Roll out behind compatibility safeguards where consumers exist. Verify branch head, test status, and linked commit. Never claim CI passed if it was not checked.

## 17. Architecture review checklist

A change is complete only when a reviewer can answer **yes** to all applicable questions:

- Is ownership explicit and unchanged service/API/event behavior preserved?
- Are generic utilities in their named technical library, and domain rules local?
- Can this resource evolve without editing unrelated modules?
- Do controllers avoid DB/vendor logic, services avoid HTTP, and repositories avoid brokers?
- Is a new directory justified by a real responsibility?
- Are input, authorization and tenant boundaries independently enforced?
- Are all write steps atomic where needed; is idempotency concurrency-safe?
- Is the event's source of truth, delivery guarantee and retention accurately documented?
- Are annotations implemented (rather than merely proposed), and overlays reproducible?
- Are secrets and PII protected in logs, events, storage and tests?
- Do generated clients, OpenAPI schema and controller routes agree?
- Are migration history and backwards compatibility safe?
- Are Go declarations formatted, tests run and module/product docs updated?

## 18. Focused generic references

| Topic | Authoritative detail |
|---|---|
| Annotations and processors | [annotation.md](annotation.md) |
| Search projection, filters, cursor/keyset | [searchDetails.md](searchDetails.md) |
| Domain events, consumer, producer, outbox | [events.md](events.md) |
| ACID, UnitOfWork, migrations, idempotency | [persistence.md](persistence.md) |
| HTTP routes, DTOs, validation, mappers | [api.md](api.md) |
| Providers, OAuth, external adapters | [providers.md](providers.md) |
| Patterns and operation implementation | [patterns.md](patterns.md) |
| Testing / release gates | [testing.md](testing.md) |
| Core and technical library design | [libraries.md](libraries.md) |

A service's or library's `docs/engineering/README.md` must link here and list its **actual deviations/implementations**. A service's `docs/product/README.md` must map user-visible use cases to reproducible API steps. There is no reason to put product examples in this generic master.

## Gateway boundary and route annotations

The process launcher owns one access-policy gateway. Controller packages declare `@GET`, `@POST`, `@PATCH`, and other method/path annotations. Each service's `routes/http.go` MUST compose generated typed registrations, not duplicate the controller's paths using handwritten `mux.Handle` calls. The route generator determines its tier from `controller/<resource>/{public,protected,restricted}`.

The external gateway classifies the **matched ServeMux route pattern**, then enforces the public/protected/restricted policy before asking ServeMux to dispatch. Raw request paths, untrusted forwarding headers, and directory names are NOT authentication evidence. Paths missing a recognized tier fail closed, aside from explicitly allowlisted health endpoints. The service must still enforce resource permissions and tenant ownership inside operations.

The in-process gateway has a separate typed invocation interface. A direct `@Proxy` service client does not pass through HTTP middleware and therefore MUST enforce authorization at its own boundary. Restricted external routes require a workload credential verifier, not a user/browser cookie. Maintain independent verifiers for external providers such as signed webhooks.

Gateway implementation folders follow the same SRP convention: `launcher/services/gateway/operations/{external,internal}`. Shared principal contracts belong in `libs/core/gateway`; domain-specific credential verification remains owned by the authentication service.

Use `gofmt` on all changed Go source files before a change is merged. Keep imports grouped, consistent tab indentation, one-line function signatures, and multiline bodies for control flow. Do not accept dense semicolon-separated statements as a coding convention.


## Strict module consistency gates (mandatory)

1. Executable commands MUST live in `backend/cmd/<command>/` or an operator script under `scripts/`, never in `backend/services/<module>/cmd/`. Module-owned implementation may be exposed through a narrowly scoped, public operator facade when Go's `internal` import rules require it.
2. Each bootstrap MUST separate `build.go` and `runtime.go` from `routes.go`, `client.go`, `repository.go`, `services.go`, `events.go` and `workers.go` as needed. The generic `bootstrap/impl/bootstrap.go` catch-all file is prohibited. The runtime must not impersonate an HTTP client by implementing domain service methods.
3. For a business operation exposed as a generated HTTP controller, the public typed client MUST use `@RequestLine` with the same route, request DTO and response contract. Direct `@Proxy` is allowed only for deliberate trusted internal capabilities lacking an HTTP route, with authorization and bypass implications documented.
4. Transport constants (headers, directives, cookie paths and expiry durations), domain enums and resource search-field keys MUST have one owning declaration. Use entity-owned enum constants for domain states and resource-owned search-field identifiers for filters, sorts and SQL allowlists. Add tests confirming JSON/OpenAPI field-name parity; Go field tags cannot reference constants.
5. Validator implementations MUST delegate substantial, pure resource-specific normalization and search policy to `validators/<resource>/helpers/<purpose>.go`. SQL expressions stay in `repository/<resource>/helpers/search/`. Neither layer may construct SQL identifiers from untrusted fields.
6. Authentication and authorization use cases belong in `services/<capability>/`, not an unrelated top-level technical folder. Keep transport-only cookies in `transport/` and gateway cross-module policy in Launcher.
7. Persistently configurable business selections (including the active Integration mail sender for a tenant) MUST NOT be frozen in service environment variables. Environment settings are reserved for deployment secrets, trusted origins, runtime endpoints and other deployment-owned invariants. Validate dynamic selections against tenant scope, active connection state and provider capability every use.
8. Public routes for account activation, invitation grants or password recovery do not imply public account registration. Protected account creation is mandatory; grant-bound activation may be public-tier only if token consumption and business invariants are independently enforced.
9. Changed source MUST pass `gofmt`, generation/build, authorization/tenant tests and documentation parity before merging. Avoid string-matching error messages when a typed/sentinel error can preserve the same safe API mapping.

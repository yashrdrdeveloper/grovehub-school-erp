# GrooveHub | Backend Annotation Framework - Master Plan

**Version:** 1.3 Phase 2 migration verified in Go 1.27 CI  |  **Date:** 10 November 2026  |  **Branch:** `phase-2.1-identity-account`  |  **Change type:** Generator implementation, selected backend migrations and CI

**Scope:** Complete Go modular-monolith annotation catalog, all backend folders and cross-cutting concerns, 23 GoF patterns, phased execution, acceptance criteria and risks.

---

## Phase 1 / Phase 2 implementation agreement (supersedes P0-P8 as release stages)

**Phase 1 — annotation contract definitions (this change).** Register every named annotation centrally under `backend/libs/core/annotation`, with name/syntax, purpose, category, intended mechanism (GEN/WRAP/REG/CHECK/TAG), current implementation status, roadmap workstream and structural target. Include all 23 GoF pattern decisions, a read-only AST scanner, typed contracts and a complete reference document. **Do not apply new annotations to any service in Phase 1.** The pre-existing Foundation `@Entity` workflow remains unchanged.

**Phase 2 — implementation and migration (partially started).** Design the per-annotation argument schema, implement the approved generator or runtime wrapper, prove behavior in tests, then replace repetitive code layer by layer. The original P0-P8 roadmap below is now the **set of Phase 2 workstreams**, not a claim that all those processors are working. The current entity, route, client, mapper, configuration, validation and event migrations are tracked in [ANNOTATION_PHASE2_MIGRATION.md](ANNOTATION_PHASE2_MIGRATION.md). In particular, `@Service` is registered for completeness **but not applied to service code** under the current Phase 1 request.

The exhaustive Phase 1 catalog and status contract are published in [ANNOTATION_PHASE1_REFERENCE.md](ANNOTATION_PHASE1_REFERENCE.md). The code catalog is under `backend/libs/core/annotation/{annotation.go,inventory.go,contract/,entity/,validation/,event/,utilities/,patterns/, ...}`.

**Runtime-safety rule:** An annotation that is recognized by the Phase 1 parser is **not automatically enforced**. Do not use unimplemented `@NonNull`, `@Authorize`, `@Transactional`, `@Retry`, `@Consumer` or `@Service` as a substitute for working code.

---

### Verified Phase 2 scope (November 10, 2026)

The clean-checkout `make backend-build` has compiled on GitHub Actions with Go 1.27.3. Generated source is injected by a Go overlay without committing generated files. Active migrations are `@Entity` (GORM table names only for Foundation, Identity and Integration), `@Routes`, `@Client`/`@Proxy`, strict `@Mapper`/`@Mapping` DTO projections, `@Event` topic/type/version constants, typed Integration configuration and Foundation/Identity request validation. Annotation-specific test cases were removed at request; CI continues to verify compilation, package loading and unrelated backend unit suites. Live database, NATS, HTTP authorization and full Ginkgo/Gomega regression coverage are separate acceptance gates.

`@Entity schema.table` now means **GORM table mapping only for every entity**. `mode=table` is still accepted for backward compatibility but is not required. Sensitive fields never receive generated setters or builders. For example:

```go
// @Entity identity.password_credentials mode=table
type PasswordCredential struct {
    PasswordHash string `gorm:"column:password_hash;not null"`
}
```

This adds only GORM's `TableName()`. Default `@Entity foundation.person_contacts` does **not** generate builders, constructors or accessors.

The original upstream Gombok CLI was removed from the unified generation pipeline. `annotation/generation/entities` supports explicit builder/getter/setter directives only on non-entity types; entity annotations exclusively generate table mappings. Standalone Gombok-only directives that are cataloged but not supported by the new generator remain **unimplemented**. See [ANNOTATION_PHASE2_MIGRATION.md](ANNOTATION_PHASE2_MIGRATION.md) for current file-level scope.

---

## 01. Executive decision and scope

> **Document status:** Phase 1 definition registry created in `backend/libs/core/annotation` (all 226 use cases, 232 unique names including GoF-only candidates, all 23 GoF pattern decisions). This means **metadata/catalog and AST discovery, not executable annotations**. The Phase 2 implementation adds Go AST-based entity, route, client, mapper, table-name and event generators, plus typed configuration and validation. Service business workflows remain handwritten.

**Folder ownership (November 2026).** Category definitions are located in dedicated annotation packages (`entity/`, `validation/`, `event/`, `utilities/`, `patterns/` and other responsibilities). The acyclic metadata contract is in `annotation/contract/`; the root registry preserves existing consumers. Annotation-specific test cases have been removed, but the backend CI build and unrelated unit suites remain enabled.

**Current entity ownership:** Foundation Person/Contact and Identity Account/PasswordCredential are passive structs; construction, nil conversions, contact copying and lifecycle updates moved to owning services. All other existing Integration entity models were already passive. `@Entity` only generates GORM table names.

**Purpose.** Replace mechanical boilerplate with typed, auditable generation while retaining readable Go and GrooveHub's explicit layer ownership. Cover all backend folders, middleware, server bootstrap, database and repositories, mappers, validation, external providers, events, messaging, observability, configuration, security, tests, scheduling, utilities and the 23 Gang of Four (GoF) patterns.

**Non-goal.** A Spring container clone, arbitrary aspect weaving, invisible business rules, auto-generated SQL from guessed method names, implicit network calls, and indiscriminate use of annotations on every Go symbol.

**Core choices**

1. Adopt `// @Name arguments` comment directives on declarations, with ordinary struct tags for field constraints. These are build-tool metadata, **not Go language features**.
2. Keep typed, deterministic generated source in `services/<service>/.generated/<package>/` and `libs/core/.generated/<package>/`; map virtual files to the maintained package using Go `-overlay` where necessary.
3. Run source discovery, syntax/type checking, annotation expansion, code generation, manifest creation and compilation through one documented developer/CI entry point. Fail on unknown annotations, collisions or unsupported combinations.
4. Prefer **existing libraries** and explicit Go implementations: the in-repository Gombok-compatible AST generator for builder/accessor generation, strict typed mapping generation (Goverter remains optional for later expansion), evaluate go-playground/validator for struct validation, GORM for persistence, standard `net/http` for routing, existing core telemetry (OpenTelemetry), NATS JetStream + Watermill for transport, `slog` for logging and Ginkgo/Gomega for tests.
5. Generated code may remove technical glue, but **business logic, security, ownership checks, lifecycle, transaction boundaries and event semantics remain explicit**.

**Acceptance standard.** No feature is considered adopted until it works on clean checkout, `make backend-build`, Air, direct build/test wrappers, CI, editor navigation (or documented fallback), deterministic regeneration, migration compatibility and existing behavior regressions.


## 02. Current repository baseline - inspected

**Repository:** `yashrdrdeveloper/grovehub-school-erp`, branch `phase-2.1-identity-account`, inspected 10 November 2026. Repository tree contained **246 tracked Go files**; module counts are approximate per folder, not lines of code.

| Area | Observed code and architecture | Planning implication |
| --- | --- | --- |
| Foundation | `entity/{person,contact,gid}.go`, GORM repository, `services`, `validators`, `mappers`, routes and Person event decorator | First vertical reference; do not move contact association rules back into entities |
| Identity | Account and credential entities; auth/account/password validators and services; cross-module Foundation client | Strong validation, secret handling, transactions and client-adapter cases |
| Integration | Gmail/Outlook provider strategies, actions/OAuth/webhooks, config loaders, repository, service and controller | Provider and config/adapter descriptors; no provider selection switches |
| Infrastructure | MessagingClient -> MessagingService -> strategy registry -> NATS/Watermill adapter | Separate domain event annotations from broker transport registration |
| Launcher | `cmd/platform`, `launcher/bootstrap`, `launcher/middleware`, `launcher/routing`, `http.Server` | Explicit startup and security wiring, carefully generated route registration |
| Core libraries | annotation, database, diff, errors, events, http, id, idempotency, logging, normalize, search, telemetry | Annotation owns generator metadata; technical capabilities stay in their existing core packages |

**Pre-existing implementation baseline (before Phase 1 registry):** `// @Entity foundation.person_contacts`, `// @Entity foundation.persons`; Gombok directives `@Builder`, `@Getter`, `@Setter` and related supported CLI directives; `libs/core/annotation/generation/{expand,gormtable,overlay}`; module-local `.generated`; `-overlay` integration. The complete build and IDE behavior are **not assumed verified**.

**Observed pain points:** repeated field conversion in `foundation/mappers/impl/person.go`; long normalization/search checks in `foundation/validators/helpers/person.go`; route registration duplication across modules; in-process client forwarding wrappers; manual provider catalog/strategy registration; config `os.Getenv` parsing; repeated logging and tracing context; extensive service construction in `foundation.go`/`identity.go`/`integration.go`.

**Current policy constraints:** `docs/AGENT.md` and backend guidelines require passive entities, service-owned business logic, mapper-only representation changes, validator-only structural validation, repository-only persistence, route-only registration, service-owned events, Infrastructure-owned transport, interface/implementation separation and capability-first core packages.


## 03. Framework semantics: syntax, ownership and output

### Syntax and parsing

- Declaration directive: `// @Entity foundation.person_contacts` immediately preceding `type Contact struct {...}`.
- Field rules should normally be tags, e.g. ``validate:"required,max=100" normalize:"trim"``. Do not invent a second expression language unless existing tags cannot express a verified requirement.
- Method directive: `// @POST /` immediately preceding an interface method or declared route contract. Parameters, path, auth and response schema must be typed and validated.
- Annotation identifiers resolve through a versioned registry; unknown, misspelled, duplicated or misplaced directives are **compile-time diagnostics**, not silently ignored comments.
- Every proposal in this document is **subject to capability-specific implementation and contract review**. The examples are illustrative syntax, not working code.

### Generation engine: proposed pipeline

`source discovery -> parse Go AST -> evaluate annotation registry -> validate types/constraints -> build intermediate graph -> resolve cross-package references -> generate typed files -> deterministic Go format -> source hashes/manifests -> overlay mapping -> go build/go test/gopls`.

- Single orchestration entry point, backed by small capability generators under `libs/core/annotation/generation/`.
- Stable naming, source-position diagnostics, preserved type parameters/import aliases, duplicate detection, explicit warning/error codes, and per-module caches.
- Generate only affected packages, but regenerate dependents when public contracts change. Atomic writes and file locks prevent Air/CI races. Delete stale outputs safely.
- Do not parse generated files as maintained input; do not publish generated files to Git. Include an inspectable generation report (`sources -> outputs -> rules`) and a `--dry-run`/`--explain` mode.
- Use exact toolchain versions and lock generator dependencies; record an annotation language/schema version and generator version in metadata.
- Go overlays are a **build input**; direct `go test` and third-party tooling without the manifest will not see virtual code. Gopls support must be demonstrated, not presumed.

### Modes and statuses

**GEN** = generate typed Go code; **WRAP** = explicit generated decorator/wrapper; **REG** = generate descriptor/registry/wiring; **CHECK** = compile-time checker or linter only; **TAG** = struct tag interpreted by a known validator; **AVOID** = intentionally do not automate that behavior. 

**EXISTING** = observed on branch, **PROPOSED** = not implemented, **OPTIONAL** = consider only with proven recurring need, **DISCOURAGED** = do not introduce as a default.

### Safety boundaries

- A comment must never silently make a service transactional, authorized, idempotent or reliable; a generated wrapper must be visible and wired at an explicit boundary.
- Do not generate setters for immutable IDs, password hashes, tenant ownership or security-sensitive fields by default.
- Preserve PATCH distinction: omitted vs explicit JSON `null` vs provided value; pointers alone may be insufficient.
- Do not implicitly log request bodies, access tokens, PII, secrets or complete Foundation Person event payloads.
- For side effects, enforce delivery/duplicate semantics, explicit retries and idempotency. Never imply exactly-once processing from JetStream acknowledgments.


## 04. Full annotation catalog by backend responsibility

## Annotation inventory: complete proposed register

This catalog names the planned semantics for the backend. **Only the observed Foundation entity/Gombok capability is currently available; the rest are proposals.** Annotation names may be amended by design review to avoid duplicates and align with third-party library syntax.

### 01. Entity, DTO and model lifecycle

| Directive | Expected effect | Mode | Status / phase |
| --- | --- | --- | --- |
| @Entity schema.table [mode=table] | Generate GORM TableName only; all construction and mutation belongs to owning service | GEN | EXISTING / P0 |
| @Builder | Fluent typed builder; preserve field types and optional semantics | GEN | EXISTING / P0 |
| @Getter | Typed getter; suppress conflicting custom methods | GEN | EXISTING / P0 |
| @Setter | Mutable typed setter; never default for sensitive fields | GEN | EXISTING / P0 |
| @AllArgsConstructor | All-fields constructor when appropriate | GEN | GOMBOK / P0 |
| @NoArgsConstructor | Zero-argument constructor; use sparingly in Go | GEN | GOMBOK / P0 |
| @RequiredArgsConstructor | Required field constructor; evaluate Go-specific meaning | GEN | GOMBOK / P0 |
| @ToString | String formatting; redact sensitive fields | GEN | GOMBOK / P0 |
| @Equals | Equality; explicit semantics for pointers/slices/PII | GEN | GOMBOK / P0 |
| @Data | Bundle DTO convenience only, NOT blanket entity setter generation | GEN | PROPOSED / P2 |
| @Immutable | Exclude mutators; enforce construction-only fields | CHECK/GEN | PROPOSED / P2 |
| @Ignore | Exclude selected fields from a specified generator | CHECK | PROPOSED / P2 |
| @Default | Typed default on builder/config input; avoid hidden persistence defaults | GEN | OPTIONAL / P2 |
| @Sensitive | Mark secret/PII for logging and export policy checks | CHECK | PROPOSED / P7 |
| @Redact | Generate safe string/log representation or redaction metadata | GEN/CHECK | PROPOSED / P7 |
| @Version | Optional mapping to explicit optimistic-lock field | CHECK | OPTIONAL / P4 |
| @SoftDelete | Optional persistence policy; only where schema supports it | CHECK | OPTIONAL / P4 |

**Boundary:** GORM tags continue to own column definitions; do not generate schema migrations from struct annotations by default.

### 02. Mappers and representation conversion

| Directive | Expected effect | Mode | Status / phase |
| --- | --- | --- | --- |
| @Mapper | Generate typed mapping implementation for explicitly annotated interface methods | GEN | PILOT / P2 |
| @Mapping | Explicit source-to-target projection; require every destination field | GEN | PILOT / P2 |
| @MapIgnore | Prevent source/destination field copy; required for secrets | CHECK/GEN | PROPOSED / P2 |
| @MappingIgnore | Deprecated alias for @MapIgnore; no separate implementation | CHECK | DISCOURAGED / P2 |
| @ConvertWith | Call declared pure conversion function | GEN | PROPOSED / P2 |
| @EnumMapping | Enum/name conversion with exhaustive unknown handling | GEN | PROPOSED / P2 |
| @AfterMapping | Invoke explicit pure post-mapper, not business behavior | GEN | OPTIONAL / P2 |
| @CollectionMapping | Typed collection conversion with nil/empty policy | GEN | OPTIONAL / P2 |
| @Flatten | Nested field mapping for proven DTO boundary | GEN | OPTIONAL / P2 |
| @MappingPolicy | Strict unmapped-field/null/overflow behavior | CHECK | PROPOSED / P2 |

**Boundary:** First library candidate: Goverter. Must reject unintended secret export, unmapped response fields, UUID/string mistakes, and JSON contract drift.

### 03. Validation, normalization and structured errors

| Directive | Expected effect | Mode | Status / phase |
| --- | --- | --- | --- |
| @Validator | Typed validator registration / wrapper for DTO type | GEN/REG | PROPOSED / P2 |
| @Valid | Nested/request validation at explicit service/controller boundary | WRAP/TAG | PROPOSED / P2 |
| @NonNull | Reject nil/absent nullable input; distinguish from blank/zero | TAG | PROPOSED / P2 |
| @NotBlank | Trim-aware non-whitespace text validation | TAG | PROPOSED / P2 |
| @NotEmpty | Non-empty string/list/map | TAG | PROPOSED / P2 |
| @Length | Rune-based or byte-based length with declared unit | TAG | PROPOSED / P2 |
| @Min | Minimum numeric value or length, explicit semantics | TAG | PROPOSED / P2 |
| @Max | Maximum numeric value or length, explicit semantics | TAG | PROPOSED / P2 |
| @Range | Bounded numeric/rune length range | TAG | PROPOSED / P2 |
| @Email | Email syntax; coordinate with normalized canonical value | TAG | PROPOSED / P2 |
| @Phone | Phone technical format and supported region policy | TAG | PROPOSED / P2 |
| @UUIDv7 | UUIDv7 requirement, not generic UUID-only | TAG | PROPOSED / P2 |
| @UUID | Generic UUID syntax when version is unrestricted; prefer @UUIDv7 for domain GIDs | TAG | OPTIONAL / P2 |
| @Pattern | Constrained regex/pattern with bounded compile cost | TAG | OPTIONAL / P2 |
| @OneOf | Allow-listed enum values and unknown handling | TAG | PROPOSED / P2 |
| @UniqueItems | Duplicate detection with explicit key/normalization | TAG | PROPOSED / P2 |
| @Normalize | Named normalization pipeline: trim, email, phone, UUID | TAG/GEN | PROPOSED / P2 |
| @RequiredIf | Conditional field requiredness | TAG | OPTIONAL / P2 |
| @MutuallyExclusive | Cross-field XOR invariant | TAG | OPTIONAL / P2 |
| @CrossField | Register custom typed cross-field rule | REG | PROPOSED / P2 |
| @ValidatePatch | Presence-aware PATCH normalization/validation | GEN | PROPOSED / P2 |
| @Error | Resolve domain error definition from constants/error; not stringly typed | CHECK/WRAP | PROPOSED / P2 |
| @StrictJSON | Reject unknown properties at transport boundary | WRAP | OPTIONAL / P3 |

**Boundary:** Use go-playground/validator as a candidate for technical tags; keep existing `core/normalize` for canonicalization. Do not move authorization, uniqueness queries or business state transitions into validators.

### 04. Persistence, repository, transactions and search

| Directive | Expected effect | Mode | Status / phase |
| --- | --- | --- | --- |
| @Repository | Validate repository ownership declaration (`entity=Type`), without wiring or transactions | CHECK | PILOT / P4 |
| @CRUD | Opt-in GORM basic CRUD only; no guessed aggregate persistence | GEN | OPTIONAL / P4 |
| @Query | Explicit typed query declaration; parameterized only | GEN/CHECK | OPTIONAL / P4 |
| @FindBy | Allowlisted trivial finder generation; fail on ambiguity | GEN | OPTIONAL / P4 |
| @Transactional | Generated wrapper using core/database/transaction | WRAP | PROPOSED / P4 |
| @ReadOnly | Read-only transaction intent; enforce at DB where possible | WRAP | OPTIONAL / P4 |
| @OptimisticLock | Version predicate plus conflict mapping | GEN | OPTIONAL / P4 |
| @TenantScoped | Require trusted tenant context in repository predicates | CHECK/WRAP | PROPOSED / P7 |
| @SoftDelete | Optional GORM soft-delete integration | GEN | OPTIONAL / P4 |
| @Searchable | Generate typed Search method delegating to shared `core/search.Execute` with maintained config | GEN | PILOT / P4 |
| @Filterable | Verify configured filter field keys match declared allowlist | CHECK | PILOT / P4 |
| @Sortable | Verify configured sort field keys match declared allowlist | CHECK | PILOT / P4 |
| @Cursor | Verify stable cursor field matches configured sort cursor expression | CHECK | PILOT / P4 |
| @Migration | Only mark migration files for linting; never auto-apply DDL | CHECK | DISCOURAGED / P4 |

**Boundary:** GORM remains the ORM. Foundation contact upsert, search projection, optimistic concurrency conflict resolution, and database error mapping remain explicit repository code. No repository may invent service business decisions.

### 05. Services, composition and dependency injection

| Directive | Expected effect | Mode | Status / phase |
| --- | --- | --- | --- |
| @Service | Register service contract/constructor for generated composition | REG | PROPOSED / P5 |
| @Component | Generic registered component when no narrower role | REG | OPTIONAL / P5 |
| @Inject | Constructor parameter binding, not reflection field injection | REG | OPTIONAL / P5 |
| @Provides | Explicit provider function for external dependencies | REG | PROPOSED / P5 |
| @Named | Qualify dependencies by stable identifier | REG | OPTIONAL / P5 |
| @Scope | Lifetime metadata; default service-instance scope | REG | OPTIONAL / P5 |
| @PostConstruct | Explicit startup hook; ordered and failure-aware | REG | OPTIONAL / P5 |
| @PreDestroy | Explicit close/shutdown hook with deadlines | REG | OPTIONAL / P5 |
| @FeatureFlag | Register typed feature gate; not authorization | WRAP/CHECK | OPTIONAL / P7 |
| @Idempotent | Apply declared key scope/result/retention rules at service boundary | WRAP | PROPOSED / P4 |

**Boundary:** Keep use-case methods handwritten. Compile-time wiring is preferable to global locator/reflection; detect cycles and ambiguity before runtime.

### 06. HTTP controllers, routes, request/response and API docs

| Directive | Expected effect | Mode | Status / phase |
| --- | --- | --- | --- |
| @Controller | Register typed HTTP controller/adapter | REG | PROPOSED / P3 |
| @Routes | Controller/resource path prefix and grouping | REG | PROPOSED / P3 |
| @GET | GET route with typed path and handler contract | GEN | PROPOSED / P3 |
| @POST | POST route with typed handler contract | GEN | PROPOSED / P3 |
| @PUT | PUT route with replacement semantics | GEN | PROPOSED / P3 |
| @PATCH | PATCH route preserving field presence | GEN | PROPOSED / P3 |
| @DELETE | DELETE route | GEN | PROPOSED / P3 |
| @PathParam | Typed path extraction with decoding/validation | GEN | OPTIONAL / P3 |
| @QueryParam | Typed query extraction, limits and defaults | GEN | OPTIONAL / P3 |
| @HeaderParam | Typed header extraction; do not trust user/tenant IDs | GEN | OPTIONAL / P3 |
| @RequestBody | Strict typed decoding/size bound | GEN | PROPOSED / P3 |
| @Response | Typed response serializer using shared envelope | GEN | PROPOSED / P3 |
| @Status | Declared success status and error mapping | GEN/CHECK | OPTIONAL / P3 |
| @Consumes | Request media-type constraints | CHECK | OPTIONAL / P3 |
| @Produces | Response content-type constraints | CHECK | OPTIONAL / P3 |
| @OpenAPI | Generate or validate OpenAPI from contract source | CHECK/GEN | PROPOSED / P3 |
| @Schema | DTO/OpenAPI schema metadata checked against authoritative public schema | CHECK/GEN | PROPOSED / P3 |
| @OperationID | Stable API operation identity | CHECK | OPTIONAL / P3 |
| @Example | API examples with contract validation | CHECK | OPTIONAL / P3 |
| @Deprecated | Deprecation metadata and compatibility warning | CHECK | OPTIONAL / P3 |
| @Public | Explicit public endpoint exemption, security-reviewed | CHECK | PROPOSED / P7 |

**Boundary:** Keep `routes/http.go` as composition point grouped by resource; never embed business logic or direct repo calls in route generation. OpenAPI remains the canonical public contract under API guidelines.

### 07. HTTP middleware and cross-cutting web concerns

| Directive | Expected effect | Mode | Status / phase |
| --- | --- | --- | --- |
| @UseMiddleware | Ordered middleware bindings from approved registry | REG | PROPOSED / P3 |
| @RequestID | Generate/propagate unique request ID | WRAP | PROPOSED / P3 |
| @CorrelationID | Propagate correlation ID distinct from request ID | WRAP | PROPOSED / P3 |
| @Trace | OpenTelemetry instrumentation and context propagation | WRAP | PROPOSED / P6 |
| @Log | Structured request/operation log policy; redact by default | WRAP | PROPOSED / P6 |
| @Recover | Panic-to-safe-HTTP response, no suppressed observability | WRAP | PROPOSED / P3 |
| @Timeout | Deadline enforcement with cancellation-aware code | WRAP | PROPOSED / P3 |
| @RateLimit | Tenant/actor-scoped limits and bounded state | WRAP | PROPOSED / P7 |
| @CORS | Approved-origin policy, not permissive wildcard default | WRAP | OPTIONAL / P7 |
| @MaxBody | Request size limit before decoding | WRAP | PROPOSED / P3 |
| @SecureHeaders | Security headers / transport policy | WRAP | OPTIONAL / P7 |
| @CacheControl | HTTP caching semantics (private/no-store as needed) | WRAP | OPTIONAL / P7 |
| @CSRF | Browser session/unsafe-request protection where applicable | WRAP | OPTIONAL / P7 |
| @Compress | Compression with content-type and side-channel checks | WRAP | OPTIONAL / P7 |

**Boundary:** The order must be explicit: recovery, tracing/correlation, authentication, authorization, limits, decoding and response. The exact sequence is an architecture decision with integration tests.

### 08. Security, secrets, tenancy and policy

| Directive | Expected effect | Mode | Status / phase |
| --- | --- | --- | --- |
| @Authenticate | Require verified identity; never trust caller-filled principal | WRAP/CHECK | PROPOSED / P7 |
| @Authorize | Evaluate business permission against trusted actor/context | WRAP | PROPOSED / P7 |
| @RolesAllowed | Role-based rule using central IAM policy | WRAP | OPTIONAL / P7 |
| @RequiresPermission | Named policy requirement with resource context | WRAP | PROPOSED / P7 |
| @OwnerOnly | Ownership check must resolve resource via service safely | WRAP | OPTIONAL / P7 |
| @TenantScoped | Enforce tenant context across route/service/repository | CHECK/WRAP | PROPOSED / P7 |
| @Sensitive | Data classification for logs/contracts/events | CHECK | PROPOSED / P7 |
| @Redact | Redaction of diagnostics/log output | WRAP/GEN | PROPOSED / P7 |
| @Secret | Config/credential never exposed in diagnostics | CHECK | PROPOSED / P5 |
| @Encrypt | Encryption policy for stored fields, provider-owned crypto | WRAP | OPTIONAL / P7 |
| @Audit | Emit bounded, access-controlled audit record | WRAP | PROPOSED / P6 |

**Boundary:** Declare policies only if the runtime checks them; generation must reject routes that declare permissions without a configured policy implementation. Never derive trust from annotations alone.

### 09. Typed HTTP clients, provider adapters and external services

| Directive | Expected effect | Mode | Status / phase |
| --- | --- | --- | --- |
| @Client | Typed in-process or HTTP client descriptor / adapter | GEN/REG | PROPOSED / P3 |
| @ExternalService | Explicit third-party boundary and provider ownership | REG | PROPOSED / P5 |
| @Provider | Provider identity, capabilities and factory registration | REG | PROPOSED / P5 |
| @Adapter | Bind vendor implementation to stable provider interface | REG | PROPOSED / P5 |
| @Endpoint | Typed base URL/route metadata from configuration | REG | PROPOSED / P5 |
| @HTTPMethod | Outgoing request verb/path mapping | GEN | OPTIONAL / P5 |
| @OAuth2 | OAuth flow configuration metadata, no hidden token grant | REG | OPTIONAL / P5 |
| @Credential | Scoped credential reference/injection; redact | REG | PROPOSED / P5 |
| @Retry | Bounded idempotent retry policy, explicit status set | WRAP | PROPOSED / P5 |
| @Backoff | Deterministic backoff with jitter/cap | WRAP | PROPOSED / P5 |
| @CircuitBreaker | Bounded failure protection per provider/tenant | WRAP | OPTIONAL / P5 |
| @Bulkhead | Limit concurrent provider calls | WRAP | OPTIONAL / P5 |
| @Fallback | Typed fallback function for safe error types | WRAP | OPTIONAL / P5 |
| @Webhook | Webhook route/handler registration metadata | REG | PROPOSED / P5 |
| @VerifySignature | Invoke provider-specific authenticity verifier | WRAP | PROPOSED / P5 |
| @HealthProbe | Explicit provider health check with timeout | REG | OPTIONAL / P5 |

**Boundary:** No generic annotation should bypass Gmail/Outlook OAuth capability checks, provider verification, secret stores, refresh rules, or vendor rate limits. Integration owns provider connections; Communication will own messages and sync cursors.

### 10. Events, envelopes and event lifecycle

| Directive | Expected effect | Mode | Status / phase |
| --- | --- | --- | --- |
| @Event | Generate stable topic/type/schema version identifiers without automatic publication | GEN/CHECK | PILOT / P6 |
| @EventVersion | Stable event schema version and compatibility checks | CHECK | PROPOSED / P6 |
| @EventHandler | Typed event handler registration | REG | PROPOSED / P6 |
| @PublishEvent | Explicit event publisher decorator metadata only | REG/WRAP | OPTIONAL / P6 |
| @AfterCommit | Bind to existing after-commit callback semantics | WRAP | PROPOSED / P6 |
| @Outbox | Opt-in durable outbox persistence/relay contract | REG | OPTIONAL / P6 |
| @EventSchema | Contract/schema compatibility checks | CHECK | OPTIONAL / P6 |
| @EventSource | Domain source metadata, not envelope field | REG | OPTIONAL / P6 |
| @Causation | Propagate existing command/parent event ID | WRAP | PROPOSED / P6 |
| @Correlation | Propagate existing correlation/trace context | WRAP | PROPOSED / P6 |
| @EventPII | Classify full-snapshot events and retention policy | CHECK | PROPOSED / P7 |

**Boundary:** Keep core `Envelope[T]`. Foundation Person changed event uses `data: { operation, old, new }`; Identity Account event and Integration mailbox changed retain stable event types/topics. Never automatically emit events for idempotent no-op or rolled-back operations.

### 11. Messaging producer, consumer, broker topology and delivery

| Directive | Expected effect | Mode | Status / phase |
| --- | --- | --- | --- |
| @Producer | Typed publisher binding to MessagingClient, domain topic | REG/GEN | PROPOSED / P6 |
| @Consumer | Typed consumer registration using MessagingClient | REG | PROPOSED / P6 |
| @Subscribe | Stable topic/subject and handler association | REG | PROPOSED / P6 |
| @Subject | Typed subject metadata; no per-event randomized subjects | REG/CHECK | PROPOSED / P6 |
| @Stream | Stream configuration owned by Infrastructure provisioning | REG/CHECK | OPTIONAL / P6 |
| @Durable | Stable durable consumer identifier | REG | PROPOSED / P6 |
| @QueueGroup | Consumer work-sharing group metadata | REG | OPTIONAL / P6 |
| @Ack | Explicit ack policy and success handling | WRAP | PROPOSED / P6 |
| @Nack | Explicit failure retry or reject behavior | WRAP | OPTIONAL / P6 |
| @RetryPolicy | Consumer delivery attempt/backoff constraints | REG | PROPOSED / P6 |
| @MaxDeliver | Upper bound on attempts; poison advisory handling | REG | PROPOSED / P6 |
| @DeadLetter | Explicit dead-letter/recovery workflow; not assumed native DLQ | REG | OPTIONAL / P6 |
| @Concurrency | Consumer worker limits/backpressure | REG | PROPOSED / P6 |
| @Batch | Batch processing/ack semantics | REG | OPTIONAL / P6 |
| @IdempotentConsumer | Deduplication key and state-store contract | WRAP | PROPOSED / P6 |
| @MessageTrace | W3C trace context publish/extract | WRAP | PROPOSED / P6 |
| @MessageTimeout | Processing timeout / ack wait coordination | WRAP | PROPOSED / P6 |

**Boundary:** JetStream consumer acknowledgments are at-least-once, not exactly-once side effects. MaxDeliver advisory is not an automatic application dead-letter queue; explicit operational handling and recovery are needed.

### 12. Configuration, profiles, typed environment and secrets

| Directive | Expected effect | Mode | Status / phase |
| --- | --- | --- | --- |
| @Configuration | Typed config struct registered for startup loading | REG/GEN | PROPOSED / P5 |
| @ConfigPrefix | Namespace environment properties by service | GEN | PROPOSED / P5 |
| @Env | Bind explicitly named environment variable | GEN | PROPOSED / P5 |
| @Default | Typed default with explicit precedence | GEN | PROPOSED / P5 |
| @Required | Configuration missing error at startup | CHECK | PROPOSED / P5 |
| @ValidateConfig | Run technical config validation before module startup | WRAP | PROPOSED / P5 |
| @Secret | Mark sensitive config redaction and source restrictions | CHECK | PROPOSED / P5 |
| @ConfigSource | Describe allowed source (env, file, secret store) | REG | OPTIONAL / P5 |
| @Profile | Explicit development/test/production selection | REG/CHECK | OPTIONAL / P5 |
| @Reloadable | Controlled reload with atomic switch and validation | WRAP | OPTIONAL / P5 |
| @FeatureFlag | Typed feature selection; not a security boundary | REG | OPTIONAL / P5 |

**Boundary:** Configuration contract stays in each service `config/`; `impl` loads it. Missing optional Gmail/Outlook credentials disable only their optional capabilities, never silently enable insecure defaults.

### 13. Logging, tracing, metrics, auditing and diagnostics

| Directive | Expected effect | Mode | Status / phase |
| --- | --- | --- | --- |
| @Log | Structured log wrapper through `slog`/core logging | WRAP | PROPOSED / P7 |
| @Trace | OpenTelemetry span boundary with context | WRAP | PROPOSED / P7 |
| @Span | Named child span with attributes allowlist | WRAP | OPTIONAL / P7 |
| @Timed | Histogram for operation duration with bounded labels | WRAP | PROPOSED / P7 |
| @Counted | Counter for stable operation/result labels | WRAP | OPTIONAL / P7 |
| @Metric | Register stable metric identity and unit | REG | OPTIONAL / P7 |
| @Audit | Security/business audit without secret exposure | WRAP | PROPOSED / P7 |
| @Redact | Hide fields from logs, traces and errors | CHECK/WRAP | PROPOSED / P7 |
| @NoLog | Suppress payload logging, preserve audit/error metadata | CHECK | PROPOSED / P7 |
| @SLO | Attach objective metadata to monitored endpoint/operation | REG | OPTIONAL / P7 |
| @CorrelationID | Carry correlation field and trace metadata | WRAP | PROPOSED / P7 |

**Boundary:** Existing `libs/core/logging` and `libs/core/telemetry` remain the runtime API. Avoid uncontrolled high-cardinality labels, sensitive payload captures, logger duplication, missing span endings or incorrect trace propagation.

### 14. Server, launcher, startup, shutdown and background jobs

| Directive | Expected effect | Mode | Status / phase |
| --- | --- | --- | --- |
| @Module | Declare module composition entry, dependencies and exports | REG | PROPOSED / P3 |
| @Server | HTTP server policy (bind address, limits, TLS) | REG/CHECK | OPTIONAL / P3 |
| @HealthCheck | Register liveness probe function | REG | PROPOSED / P3 |
| @ReadyCheck | Register dependencies/readiness policy | REG | PROPOSED / P3 |
| @Startup | Register ordered initializer and failure behavior | REG | OPTIONAL / P3 |
| @GracefulShutdown | Stop accepting requests, drain work, close dependencies | REG | PROPOSED / P3 |
| @Shutdown | Module close handler and deadline | REG | PROPOSED / P3 |
| @Scheduled | Schedule registered function with timezone/cadence | REG | OPTIONAL / P8 |
| @Cron | Cron expression with explicit timezone/misfire handling | REG | OPTIONAL / P8 |
| @LeaderOnly | Distributed lock/lease requirement for singleton jobs | WRAP | OPTIONAL / P8 |
| @BackgroundJob | Registered worker managed by lifecycle, not hidden goroutine | REG | OPTIONAL / P8 |

**Boundary:** Preserve the existing launcher/bootstrap and module `New`/`RegisterHTTP` flow. Integration currently has a development-only loopback principal safeguard; annotations must not bypass it or turn it on in production.

### 15. Testing, fixtures, mocks, contract tests and benchmarks

| Directive | Expected effect | Mode | Status / phase |
| --- | --- | --- | --- |
| @UnitTest | Test category metadata for Ginkgo suite grouping | CHECK | OPTIONAL / P8 |
| @IntegrationTest | Service/database/broker integration category | CHECK | OPTIONAL / P8 |
| @E2E | Full HTTP/worker flow category | CHECK | OPTIONAL / P8 |
| @Fixture | Typed deterministic fixture constructor | GEN | PROPOSED / P8 |
| @Mock | Generate interface fake/mock without expectation magic | GEN | PROPOSED / P8 |
| @GenerateMock | Explicit mock target and stable generated location | GEN | PROPOSED / P8 |
| @ContractTest | Cross-module OpenAPI/client/event contract test manifest | REG/GEN | PROPOSED / P8 |
| @PropertyTest | Property/fuzz test input metadata | GEN | OPTIONAL / P8 |
| @TestMatrix | Explicit table-driven variations | GEN | OPTIONAL / P8 |
| @TestTimeout | Timeout for test suite/fixture | REG | OPTIONAL / P8 |
| @RequiresDocker | Skip/fail rationale when containers unavailable | CHECK | OPTIONAL / P8 |
| @Benchmark | Register performance benchmark and baseline | REG | OPTIONAL / P8 |
| @Seed | Deterministic seed/clock and resource cleanup | REG | OPTIONAL / P8 |
| @Golden | Golden-data drift check with explicit update step | CHECK | OPTIONAL / P8 |

**Boundary:** Prefer existing Ginkgo/Gomega suite structure; do not build a parallel runner. Tests must verify generated code behavior, not only whether a file was emitted.

### 16. Utilities, helpers, cache, concurrency and shared contracts

| Directive | Expected effect | Mode | Status / phase |
| --- | --- | --- | --- |
| @Pure | Static advisory for side-effect-free functions; lint where possible | CHECK | OPTIONAL / P8 |
| @Deprecated | Expose deprecation diagnostics across call sites | CHECK | PROPOSED / P8 |
| @Experimental | Mark unstable API and compatibility expectations | CHECK | OPTIONAL / P8 |
| @Memoize | Cache only proven pure, bounded and safe functions | WRAP | DISCOURAGED / P8 |
| @Cacheable | Explicit cache keys, tenant partition, TTL, invalidation | WRAP | OPTIONAL / P8 |
| @CacheEvict | Invalidate declared cache keys after successful mutation | WRAP | OPTIONAL / P8 |
| @CacheKey | Stable typed key fields, no secrets | GEN/CHECK | OPTIONAL / P8 |
| @TTL | Retention/expiry declaration with time unit | REG | OPTIONAL / P8 |
| @ConcurrencySafe | Document verification obligation; not an automatic race fix | CHECK | DISCOURAGED / P8 |
| @Bulkhead | Concurrency semaphore for external work | WRAP | OPTIONAL / P5 |
| @Checksum | Explicit integrity check for files/messages | WRAP | OPTIONAL / P8 |
| @Retention | Data lifecycle metadata; requires actual cleanup process | REG/CHECK | OPTIONAL / P8 |
| @UploadLimit | Enforce file request/storage limit | WRAP | OPTIONAL / P8 |
| @Encrypt | Bounded encryption policy via core/crypto if added | WRAP | OPTIONAL / P8 |
| @GenerateStringer | Typed enum stringer when needed | GEN | OPTIONAL / P8 |

**Boundary:** No `@Utils` annotation: helpers are owned by their existing capability/service and should not be generated just for folder symmetry.


## 05. GoF design patterns - all 23 and annotation decision

These are the **23 classic GoF patterns**, now registered as Phase 1 metadata decisions, not 23 mandatory executable annotations. A pattern is an architectural technique, not automatically a code-generation feature. Prefer `@Strategy`, `@Adapter`, `@FactoryMethod`, `@Decorator`, `@Proxy`, `@Observer` and `@Chain` for meaningful registry/wrapper generation. Most other names would only decorate code without eliminating real boilerplate and should remain plain Go or a linter.

| Family | Pattern | Candidate annotation | Allowed automation / purpose | Decision |
| --- | --- | --- | --- | --- |
| Creational | Abstract Factory | @AbstractFactory | Registry of related provider families; explicit factory contract | Optional registry |
| Creational | Builder | @Builder | Typed fluent construction | Supported |
| Creational | Factory Method | @FactoryMethod | Register constructors for a capability | Useful registry |
| Creational | Prototype | @Prototype | Explicit clone/copy contract only | Avoid default auto-clone |
| Creational | Singleton | @Singleton | Lifetime constraint in composition root | Avoid global instance |
| Structural | Adapter | @Adapter | Bind vendor to stable interface | Useful registry |
| Structural | Bridge | @Bridge | Map abstraction to independent implementor | Optional metadata |
| Structural | Composite | @Composite | Declare child composition, no hidden recursion | Optional metadata |
| Structural | Decorator | @Decorator | Explicit typed wrapper order | Useful generated wiring |
| Structural | Facade | @Facade | Mark stable public capability API | Check only |
| Structural | Flyweight | @Flyweight | Object sharing semantics require manual proof | Avoid generation |
| Structural | Proxy | @Proxy | Typed forwarding for client/resilience boundary | Useful generation |
| Behavioral | Chain of Responsibility | @Chain | Ordered middleware/handler chain | Useful registry |
| Behavioral | Command | @Command | Typed command registry and dispatch | Optional registry |
| Behavioral | Interpreter | @Interpreter | Registry of parsing grammar components | Avoid auto-parser |
| Behavioral | Iterator | @Iterator | Typed iteration interface only when needed | Avoid generation |
| Behavioral | Mediator | @Mediator | Explicit event/notification coordination | Optional registry |
| Behavioral | Memento | @Memento | Snapshot semantics with sensitive-state controls | Avoid automation |
| Behavioral | Observer | @Observer | Listener registration around event contract | Useful registry |
| Behavioral | State | @State | Register state transitions with guards | Check only initially |
| Behavioral | Strategy | @Strategy | Provider capability registration + selection | Useful registry |
| Behavioral | Template Method | @TemplateMethod | Common skeleton remains an explicit Go method | Avoid inheritance emulation |
| Behavioral | Visitor | @Visitor | Typed dispatch with explicit exhaustive cases | Optional, rare |

**Applied examples.** Integration uses provider strategy/factory/resolver; Infrastructure uses messaging strategy with NATS adapter; Foundation/Identity event implementations decorate services; Foundation in-process client is a proxy; launcher middleware is a chain. None of these should switch to runtime reflection solely to imitate Java.

**Disallow** automatic inheritance emulation for Template Method, global singleton mutation, implicit clone of security-sensitive state, and magical State transitions that bypass services.


## 06. Folder-by-folder coverage and ownership map

The guiding rule is **annotation metadata belongs in `libs/core/annotation`; runtime behavior remains owned by its actual capability/service**. The following map covers existing folders and future folders only where justified.

| Folder / layer | Potential annotation or generator output | Never delegate to annotation |
| --- | --- | --- |
| `backend/cmd/platform` | `@Server`, startup/telemetry configuration checks | Actual process exit/shutdown policy |
| `backend/launcher/bootstrap` | `@Configuration`, `@Provides`, startup dependency graph | Security defaults, secret loading and DB failure policy |
| `backend/launcher/middleware` | `@UseMiddleware`, `@RequestID`, `@Trace`, `@Authenticate`, `@Recover` | Runtime request security and authorization guarantees |
| `backend/launcher/routing` | `@HealthCheck`, `@ReadyCheck`, route manifest | Operational health semantics |
| `services/<service>/<service>.go` | `@Module`, `@Service`, `@Repository`, `@Provides` registries | Service composition semantics and startup dependencies |
| `services/<service>/entity` | `@Entity`, `@Builder`, `@Getter`, `@Immutable` | Business rules/ownership/transitions |
| `services/<service>/dto/request` | `@Data`, validation/normalization tags, `@Valid` | Trusted principal, business checks |
| `services/<service>/dto/response` | `@Data`, `@Sensitive`, `@Schema` | Expose raw persistence records |
| `services/<service>/validators` | `@Validator`, `@CrossField`, `@Error` | Repository lookups/authorization/ID creation |
| `services/<service>/validators/helpers` | Typed normalizers and custom validators | Business lifecycle decisions |
| `services/<service>/mappers` | `@Mapper`, `@Mapping`, `@ConvertWith` | Business validation, I/O |
| `services/<service>/repository` | `@Repository`, `@CRUD`, `@Query`, `@Transactional` | Unknown SQL/data lifecycle guesses |
| `services/<service>/services` | `@Service`, wrapper policies, `@Idempotent` | Use-case implementation and decisions |
| `services/<service>/services/helpers` | `@Pure` lint only, selected safe utilities | Service-owned association/business rules |
| `services/<service>/controller` | `@Controller`, `@RequestBody`, `@Response`, security wrappers | Direct repository access |
| `services/<service>/routes` | `@Routes`, HTTP verb directives | Business behavior or transport-private secrets |
| `services/<service>/client` | `@Client`, `@Proxy`, `@Retry` when remote | Cross-module internal repository access |
| `services/<service>/events` | `@Event`, `@EventHandler`, `@AfterCommit` | Blind automatic publication |
| `services/<service>/providers` | `@Strategy`, `@FactoryMethod`, `@Adapter`, `@Provider` | Vendor selection in consumer business code |
| `services/<service>/externalService` | `@ExternalService`, `@VerifySignature`, `@OAuth2` | Domain workflows and trust shortcuts |
| `services/<service>/config` | `@Configuration`, `@Env`, `@Secret` | Reading env variables inside controllers |
| `services/<service>/constants` | `@OneOf`, error-symbol checks, stable event identifiers | Auto-renaming public error/event codes |
| `services/<service>/migrations` | Schema compatibility lint, no DDL generation | Auto migration on startup |
| `services/<service>/tests` | `@Fixture`, `@Mock`, `@ContractTest`, `@Benchmark` | Replacing behavior tests with generator snapshots |
| `libs/core/database` | Transaction wrapper adapter, typed query utilities | Domain-specific GORM models |
| `libs/core/events` | Envelope serialization and trace/correlation helpers | Domain event meaning/topic ownership |
| `libs/core/http` | Decoder/envelope/middleware wrappers | Resource policy decisions |
| `libs/core/logging` | Structured log decorators and redaction metadata | Sensitive data capture by default |
| `libs/core/telemetry` | Span/metric instrumentation adapters | Unbounded cardinality |
| `libs/core/errors` | Typed error-symbol resolution | Leak raw DB/provider errors |
| `libs/core/normalize` | Pure email/phone/UUID adapters | Domain meaning |
| `libs/core/search` | Search field/operator/cursor metadata | Raw user-controlled SQL |
| `libs/core/annotation` | AST/generation/overlays/registry/diagnostics | Duplicating every runtime capability |
| `scripts/backend` | Generator, format, build, Air and direct Go wrappers | Undocumented magical build prerequisites |
| `scripts/testing` | Overlay-aware test commands, reports | Skipping compiled generated behavior |
| `docs/engineering/backend` | Catalog, rules, plan, ADRs, compatibility | Stale/silent runtime conventions |

**Runtime dependency direction:** service packages depend on public core capabilities, not `annotation/generation`. Generated implementation imports the same explicit domain/core interfaces as handwritten code. No generated source is written into maintained `entity/`, `mappers/`, `validators/` or `routes/` folders.


## 07. Target repository shape (lazy creation only)

```text
backend/
  cmd/platform/                          # existing process entry
  launcher/
    bootstrap/ middleware/ routing/      # existing server responsibilities
  libs/core/
    annotation/
      entity/                            # existing generic builder facade
      generation/
        expand/ gormtable/ overlay/     # EXISTING implementation tools
        parser/ registry/ diagnostics/   # PLANNED parser, schema, errors
        mapper/ validation/ routes/     # PLANNED typed code generators
        client/ composition/ config/    # PLANNED wiring generators
        event/ messaging/ policy/       # PLANNED descriptors/wrappers
        fixture/                        # PLANNED test fixture generation
    database/ errors/ events/ http/ id/ idempotency/
    logging/ normalize/ search/ telemetry/ # existing capabilities
    .generated/                         # only when core needs generated sources
  services/
    foundation/
      entity/ dto/ validators/ mappers/ repository/
      services/ controller/ routes/ client/ events/
      config/ migrations/ tests/         # create only when needed
      .generated/
        entity/ mappers/ validators/ routes/ client/
        repository/ services/ events/    # folders appear on demand
    identity/
      ...
      .generated/
    integration/
      ...
      .generated/
    infrastructure/
      ...
      .generated/
  .generated/overlay.json               # GENERATED manifest, Git ignored
scripts/backend/generate-entities.sh    # existing, evolves into orchestration
scripts/backend/generate-annotations.sh # proposed unified command
scripts/testing/                        # overlay-aware tests

docs/engineering/backend/
  ANNOTATION_FRAMEWORK_MASTER_PLAN.md   # THIS document
```

**No empty folders:** The tree above is a target ownership map, not permission to scaffold every directory immediately. Existing capability-first guidelines remain authoritative.


## 08. Concrete before/after design examples - illustrative only

### 1. Existing Foundation entity

```go
// @Entity foundation.person_contacts
type Contact struct {
    PersonGID GID    `gorm:"column:person_gid;type:uuid;primaryKey"`
    Value     string `gorm:"column:value;not null"`
}
```

This now generates only `TableName()` outside `entity/`. Builders/constructors/accessors are not generated for entities; service helpers own creation, optional-value normalization and contact copying.

### 2. DTO validation with named error definition (proposed)

```go
// @Validator
// @Error foundation.InvalidPersonRequest
type CreatePerson struct {
    FirstName string `json:"firstName" validate:"required,max=100" normalize:"trim"`
    LastName  string `json:"lastName"  validate:"required,max=100" normalize:"trim"`
}
```

The exact rule syntax and package symbol resolution require a separate ADR. `max=100` must deliberately preserve the existing **rune-based** name limit; do not silently change Unicode behavior. For PATCH, field presence must be tracked separately.

### 3. Typed mapping (proposed)

```go
// @Mapper
// @Mapping GID -> GID
// @MappingIgnore PasswordHash
// @ConvertWith MapContacts
type PersonMapper interface {
    ToResponse(entity.Person) response.Person
}
```

Actual Goverter grammar must be reused when adopted; this pseudocode is illustrative. Mapper generation should fail closed on source/destination shape drift.

### 4. Routes (proposed)

```go
// @Routes /foundation/persons
// @Authenticate
type PersonController interface {
    // @POST /
    Create(http.ResponseWriter, *http.Request)
    // @GET /{gid}
    GetByGID(http.ResponseWriter, *http.Request)
}
```

Generated registration still calls existing controller implementations and approved middleware. Authoritative OpenAPI policy must remain consistent.

### 5. Provider strategy (proposed)

```go
// @Strategy mail
// @Provider gmail
// @Adapter mailprovider.Strategy
type GmailStrategy struct { /* explicit vendor dependencies */ }
```

Only register/validate the implementation. OAuth, signatures, token refresh, provider capabilities and rate limits remain handwritten and protected.

### 6. Event and consumer (proposed)

```go
// @Event integration.mailbox.changed.v1
// @Subject integration_mailbox_changed
type MailboxChangedPayload struct { /* stable public fields */ }

// @Consumer
// @Subscribe integration_mailbox_changed
// @Durable communication_mailbox_sync
// @IdempotentConsumer eventId
func HandleMailboxChanged(ctx context.Context, event coreevents.Envelope[MailboxChangedPayload]) error {
    /* explicit business processing */
    return nil
}
```

The generator can register handlers and typed decode; it **must not** claim at-most-once or exactly-once side effects. Processing/dedup persistence is explicit.

### 7. Typed config (proposed)

```go
// @Configuration
// @ConfigPrefix GROVEHUB_GMAIL
type GmailConfig struct {
    ClientID     string `env:"CLIENT_ID" validate:"required"`
    ClientSecret string `env:"CLIENT_SECRET" secret:"true"`
    RedirectURI  string `env:"REDIRECT_URI" validate:"required,url"`
}
```

Provider callback-path validation stays custom; secrets never enter logs or generated debug output.


## 09. Roadmap, deliverables and phase gates

Implementation should be incremental and entirely backward compatible until a reviewed migration point.

| Phase | What changes | Deliverables | Go/no-go gate |
| --- | --- | --- | --- |
| P0 - standardize | Catalog & ADRs; audit and stable semantics; no application code changes | Versioned annotation grammar, allowlist, docs and ownership matrix | Review approves single truth per capability and security invariants |
| P1 - engine hardening | Generalized parser, registry, diagnostics, deterministic output; stage/overlay reliability | Unified generator orchestration, clean checkout path, source manifest, explain mode | Build, Air, CI, gopls workaround, concurrency and stale cleanup verified |
| P2 - mapper/validation | Goverter evaluation + typed DTO validators/normalizers, domain errors | Foundation vertical slice, then Identity; strict field mismatch diagnostics | Unicode/PATCH/null/PII and normalization parity tests |
| P3 - API/client/bootstrap | Routes, client proxy generation, middleware metadata, health modules | Explicit route manifest and typed client adapters | OpenAPI/routes/security/trace compatibility matrix passes |
| P4 - data/service wrappers | Repository registration, limited CRUD, transactional/idempotent decorator | GORM-compliant wrappers, DI graph checks | Rollback, idempotency, concurrency, migrations unchanged |
| P5 - config/providers | Configuration, secrets, strategy/factory/adapter wiring, resilience | Gmail/Outlook descriptors, typed config loader (as opt-in) | Missing optional secrets safe, webhook verify and no production dev-auth bypass |
| P6 - events/messaging | Event schema registry; publishers, consumers, ack/retry/trace/outbox planning | Domain event metadata, NATS policy descriptions, consumer contracts | Stable topics/envelopes, no-op rules, duplicates and crash recovery verified |
| P7 - security/observability | Authorize/tenant-aware wrapper, trace/metrics/log, audit/redaction | Approved middleware policies, bounded telemetry labels | Negative security tests, no PII/secret leakage, trace continuity |
| P8 - tests/tooling/pattern polish | Fixtures, contract tests, benchmarking, optional GoF registry helpers | Golden tests, measured build-time baseline, safe migrations and operations playbook | No unmeasured annotation added; rollback strategy ready |

### Repeatable phase acceptance tests

- Clean checkout and `make backend-build` succeed with exact locked tool versions; no generated files in tracked paths.
- Generation runs only for changed packages, emits deterministic hashes and a fully inspectable manifest; fails with precise file/line diagnostics.
- Go build, Go test, Air reload, CI and documented IDE/`gopls` workflows operate consistently; user can inspect generated code.
- No duplicate methods, package import cycles or incorrect transformations; deleted annotations remove stale code.
- Structural request validation and error codes match current behavior; default and nullable handling is explicitly tested.
- GORM schema names, indexes, query filters, search cursor behavior, transaction rollback, version checks and migration state remain stable.
- Authentication and authorization fail closed; no caller-controlled tenant/actor information is trusted, and production dev-auth is disabled.
- Events maintain stable topic/type/payload, correlation/causation, after-commit semantics, and avoid idempotent duplicate publication.
- JetStream tests explicitly cover ack/nack, duplicate deliveries, redelivery, exhaustion advisory and recovery; consumer processing is idempotent.
- Logging/tracing redact secrets and PII and do not introduce high-cardinality telemetry dimensions; no silent background workers.
- Generated runtime performance is benchmarked against handwritten equivalents; codegen/IDE cost is measured independently.
- Existing Ginkgo unit, service sanity, integration and contract suites pass; rollback can disable the new generator per module.


## 10. Risks, deliberately rejected shortcuts and migration policy

| Risk | Impact | Required mitigation |
| --- | --- | --- |
| Magic comments appear to run at runtime | Security/business behavior omitted | Every wrapper generated, registered and testable; warn when annotation is not wired |
| Generated code invisible to direct Go tooling | Failed IDE navigation/tests | Document overlay workflow, build command wrappers and gopls proof |
| Annotation/API contract drift | Backward compatibility break | OpenAPI/schema diff checks; stable operation IDs and event versions |
| Silent mapping of sensitive fields | Data leakage | Strict unmapped-field policy and sensitive-field denylist |
| Mixed validation and business rules | Wrong authorization/uniqueness | Validators do pure structural checks; services own repository-backed decisions |
| Repository magic or unchecked query strings | SQL injection/tenant escape | Use typed GORM/query parameters and explicit search field allowlists |
| Transactional decorator wraps wrong boundary | Partial writes or deadlocks | Explicit core transaction usage and retry/rollback acceptance tests |
| Retry on non-idempotent operations | Duplicate side effects | Require idempotency classification/key and safe retry policy |
| Broker ack interpreted as exactly once | Duplicate downstream writes | Idempotent consumer contract and durable dedup when required |
| Automatic pattern annotation proliferation | Hard-to-read boilerplate and build cost | Only generate when repeated mechanical work is measurable |
| Annotations expand into security bypass | Privilege escalation | Fail closed when policy/tenant context unbound |
| Logging of user, provider or credential data | Privacy/security incident | Classify/redact by default, negative tests |
| Slow generation / Air loops | Developer friction | Per-package cache, atomic writes, watcher exclusions, profiling |
| Incompatible generator upgrades | Reproducibility issue | Pinned versions, stable manifests, versioned schema and rollback command |

### Migration strategy

- **Observe** current implementation before introducing any annotation for a folder.
- **Opt in per package/resource**, not globally. Keep manual implementations as the reference until behavior parity is proven.
- **Do not delete handwritten domain logic.** Remove only boilerplate that is mechanically reproduced and checked.
- Maintain alias/adapter compatibility when public interfaces or exported constructors change. Fail generation on public API mismatch.
- Keep a per-capability `annotations.yaml` policy (proposed, not implemented) documenting allowed annotations, generator version and exclusions; do not place an untyped global rule engine in `utils`.
- Do not roll out a design-pattern annotation unless it reduces real boilerplate (registration/wiring) with a measurable benefit.


## 11. Final recommendations and architectural decisions

**Must implement first:** `@Mapper`, `@Validator`, `@Valid`, required/normalization tags, field-specific error mapping, route registration and simple `@Client` adapters. These address the repetition observed in Foundation/Identity without obscuring workflows.

**Next where useful:** annotation-driven `@Configuration` and provider `@Strategy`/`@Adapter` registration, explicit transaction/idempotency wrappers, stable event/consumer metadata, and approved middleware composition.

**Keep optional or lint-only:** most GoF pattern annotations, repository query generation beyond simple CRUD, cache, scheduling, lifecycle injection and automatic code for advanced resilient transport behavior.

**Never automate blindly:** authorization, ownership, password/secret management, provider verification, data migrations, cross-service business transactions, Foundation Person event publication semantics, state-machine transitions and message dedup/recovery.

**Documentation placement:** commit this plan under `docs/engineering/backend/ANNOTATION_FRAMEWORK_MASTER_PLAN.md`; cross-reference it from `CORE_LIBRARY_GUIDELINES.md`. Follow existing `docs/AGENT.md`, `GROVEHUB_BACKEND_STRUCTURE.md`, `VALIDATOR_GUIDELINES.md`, `MAPPER_AND_DTO_GUIDELINES.md`, `EVENT_GUIDELINES.md`, `INFRASTRUCTURE_MESSAGING_STRATEGY.md`, `INTEGRATION_SERVICE_GUIDELINES.md`, `API_CONTRACT_GUIDELINES.md`, and security/testing standards. A planning document is **not** a new mandatory coding rule until approved and implemented.

**Decision needed before P1:** choose a stable annotation grammar (Go comments + struct tags), whether Goverter is the mapper backend, which validator library/rules normalize existing Unicode and PATCH semantics, the generated wrapper visibility policy, and how `gopls` will discover virtual files across modules.


## 12. Sources and implementation references

**Current repository:** [GrooveHub development branch](https://github.com/yashrdrdeveloper/grovehub-school-erp/tree/phase-2.1-identity-account) and its internal engineering guidelines (reviewed 10 November 2026).

**External technical references** (supporting design constraints, not proof of implementation):

- Go command reference (`-overlay` and limitations): https://go.dev/src/cmd/go/alldocs.go
- Goverter typed mapping generator: https://github.com/jmattheis/goverter
- `go-playground/validator/v10` tags and cross-field validation: https://pkg.go.dev/github.com/go-playground/validator/v10
- Gombok generator: https://github.com/YangTaeyoung/gombok
- GORM: https://gorm.io/docs/
- NATS JetStream consumer delivery/ack documentation: https://docs.nats.io/nats-concepts/jetstream/consumers
- Watermill: https://watermill.io/
- OpenTelemetry Go: https://opentelemetry.io/docs/languages/go/
- Go `slog`: https://pkg.go.dev/log/slog

**Reminder:** This file is a design proposal and backlog. Annotation examples other than existing `@Entity`/Gombok directives are illustrative and must not be assumed to compile before their phase is implemented.

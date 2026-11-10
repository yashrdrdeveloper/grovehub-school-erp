# GrooveHub Annotation Contracts — Phase 1 Reference

**Status:** Phase 1 contract definition only  |  **Last updated:** 10 November 2026  |  **Branch:** `phase-2.1-identity-account`

This is the complete directive inventory implemented as **typed metadata definitions**, not as working backend decorators. The authoritative architecture proposal is [ANNOTATION_FRAMEWORK_MASTER_PLAN.md](ANNOTATION_FRAMEWORK_MASTER_PLAN.md). **No services are annotated or refactored in Phase 1.**

## Two-stage delivery agreement

| Stage | Output | Explicit exclusions |
| --- | --- | --- |
| **Phase 1 — now** | Every documented directive has metadata for canonical spelling, syntax, purpose, effect category, availability, roadmap phase, and structural target; AST scanner and contract tests exist | No replacement of Foundation/Identity/Integration/Infrastructure code. No new runtime middleware, DAO, event publisher, DI container or parser injection |
| **Phase 2 — later** | Implement selected processors and typed wrappers; migrate existing code incrementally with verified parity | Never infer business decisions, trust caller credentials, invent SQL, or auto-publish domain events on repository writes |

## What a Phase 1 annotation *is*

Annotations in Go are currently `// @Name [arguments]` comments. They are not Go syntax, reflection-based decorators, or runtime interceptors. The library at `backend/libs/core/annotation` exports `All`, `Lookup`, `Uses`, `Patterns`, `ParseComment` and `ValidateDirective`. The AST scanner at `backend/libs/core/annotation/generation/parser` recognizes comments without changing files. Validation tags and executable generators/wrappers remain future Phase 2 capabilities.

The **only pre-existing executable custom directive** is Foundation's `@Entity schema.table`, implemented by the existing staging/Gombok/GORM-table/Go-overlay toolchain. Gombok's upstream annotations are also supported by that toolchain. Other entries below are definitions only, even when their intended mode is GEN, WRAP, REG or TAG.

### Contract rules

- Canonical names are case-sensitive and singular (`@Mapper`, not `@Mappers`). `@MapIgnore` is canonical; `@MappingIgnore` is a deprecated spelling and must not acquire a distinct generator.
- The scanner rejects unknown names, and `@Entity` requires a safe `schema.table` argument. Full directive-specific argument schemas and field-tag validation are gated on Phase 2.
- `@NonNull identity.InvalidRegistration` is recognized as metadata but does **not** yet reject nil or resolve the error constant; these require Phase 2 typed validation and error binding.
- Struct tags such as `validate:"required"` are **not** interpreted by this Phase 1 registry. A directive must not be considered enforced until an approved generator or runtime wrapper demonstrably performs the rule.
- Preserve service-owned domain decisions, GORM repository semantics, DTO contract precedence, API security boundaries, and existing NATS/Watermill consumers.
- Discouraged annotations are listed for completeness, not approval for production usage. Pattern metadata does not imply a design pattern is implemented.
- Generated files, when Phase 2 processors are implemented, must live in service-local `.generated/` directories, ignored by Git, and compile using a deterministic overlay.

## Complete directive inventory — 226 use cases

Repeated names indicate **different planned contexts** of a single canonical annotation (for example `@Retry` for external clients and general resilience). No duplicate runtime implementations are implied.

### Entity, DTO and model lifecycle

| Directive / usage | Intended semantics | Mechanism | Status and roadmap |
| --- | --- | --- | --- |
| `@Entity schema.table` | Gombok builder, safe accessors, GORM TableName; existing Foundation directive | GEN | EXISTING / P0 |
| `@Builder` | Fluent typed builder; preserve field types and optional semantics | GEN | EXISTING / P0 |
| `@Getter` | Typed getter; suppress conflicting custom methods | GEN | EXISTING / P0 |
| `@Setter` | Mutable typed setter; never default for sensitive fields | GEN | EXISTING / P0 |
| `@AllArgsConstructor` | All-fields constructor when appropriate | GEN | GOMBOK / P0 |
| `@NoArgsConstructor` | Zero-argument constructor; use sparingly in Go | GEN | GOMBOK / P0 |
| `@RequiredArgsConstructor` | Required field constructor; evaluate Go-specific meaning | GEN | GOMBOK / P0 |
| `@ToString` | String formatting; redact sensitive fields | GEN | GOMBOK / P0 |
| `@Equals` | Equality; explicit semantics for pointers/slices/PII | GEN | GOMBOK / P0 |
| `@Data` | Bundle DTO convenience only, NOT blanket entity setter generation | GEN | PROPOSED / P2 |
| `@Immutable` | Exclude mutators; enforce construction-only fields | CHECK/GEN | PROPOSED / P2 |
| `@Ignore` | Exclude selected fields from a specified generator | CHECK | PROPOSED / P2 |
| `@Default` | Typed default on builder/config input; avoid hidden persistence defaults | GEN | OPTIONAL / P2 |
| `@Sensitive` | Mark secret/PII for logging and export policy checks | CHECK | PROPOSED / P7 |
| `@Redact` | Generate safe string/log representation or redaction metadata | GEN/CHECK | PROPOSED / P7 |
| `@Version` | Optional mapping to explicit optimistic-lock field | CHECK | OPTIONAL / P4 |
| `@SoftDelete` | Optional persistence policy; only where schema supports it | CHECK | OPTIONAL / P4 |

### Mappers and representation conversion

| Directive / usage | Intended semantics | Mechanism | Status and roadmap |
| --- | --- | --- | --- |
| `@Mapper` | Generate typed mapping implementation for interface | GEN | PROPOSED / P2 |
| `@Mapping` | Explicit source-to-target field relation / rename | GEN | PROPOSED / P2 |
| `@MapIgnore` | Prevent source/destination field copy; required for secrets | CHECK/GEN | PROPOSED / P2 |
| `@MappingIgnore` | Deprecated alias for @MapIgnore; no separate implementation | CHECK | DISCOURAGED / P2 |
| `@ConvertWith` | Call declared pure conversion function | GEN | PROPOSED / P2 |
| `@EnumMapping` | Enum/name conversion with exhaustive unknown handling | GEN | PROPOSED / P2 |
| `@AfterMapping` | Invoke explicit pure post-mapper, not business behavior | GEN | OPTIONAL / P2 |
| `@CollectionMapping` | Typed collection conversion with nil/empty policy | GEN | OPTIONAL / P2 |
| `@Flatten` | Nested field mapping for proven DTO boundary | GEN | OPTIONAL / P2 |
| `@MappingPolicy` | Strict unmapped-field/null/overflow behavior | CHECK | PROPOSED / P2 |

### Validation, normalization and structured errors

| Directive / usage | Intended semantics | Mechanism | Status and roadmap |
| --- | --- | --- | --- |
| `@Validator` | Typed validator registration / wrapper for DTO type | GEN/REG | PROPOSED / P2 |
| `@Valid` | Nested/request validation at explicit service/controller boundary | WRAP/TAG | PROPOSED / P2 |
| `@NonNull` | Reject nil/absent nullable input; distinguish from blank/zero | TAG | PROPOSED / P2 |
| `@NotBlank` | Trim-aware non-whitespace text validation | TAG | PROPOSED / P2 |
| `@NotEmpty` | Non-empty string/list/map | TAG | PROPOSED / P2 |
| `@Length` | Rune-based or byte-based length with declared unit | TAG | PROPOSED / P2 |
| `@Min` | Minimum numeric value or length, explicit semantics | TAG | PROPOSED / P2 |
| `@Max` | Maximum numeric value or length, explicit semantics | TAG | PROPOSED / P2 |
| `@Range` | Bounded numeric/rune length range | TAG | PROPOSED / P2 |
| `@Email` | Email syntax; coordinate with normalized canonical value | TAG | PROPOSED / P2 |
| `@Phone` | Phone technical format and supported region policy | TAG | PROPOSED / P2 |
| `@UUIDv7` | UUIDv7 requirement, not generic UUID-only | TAG | PROPOSED / P2 |
| `@UUID` | Generic UUID syntax when version is unrestricted; prefer @UUIDv7 for domain GIDs | TAG | OPTIONAL / P2 |
| `@Pattern` | Constrained regex/pattern with bounded compile cost | TAG | OPTIONAL / P2 |
| `@OneOf` | Allow-listed enum values and unknown handling | TAG | PROPOSED / P2 |
| `@UniqueItems` | Duplicate detection with explicit key/normalization | TAG | PROPOSED / P2 |
| `@Normalize` | Named normalization pipeline: trim, email, phone, UUID | TAG/GEN | PROPOSED / P2 |
| `@RequiredIf` | Conditional field requiredness | TAG | OPTIONAL / P2 |
| `@MutuallyExclusive` | Cross-field XOR invariant | TAG | OPTIONAL / P2 |
| `@CrossField` | Register custom typed cross-field rule | REG | PROPOSED / P2 |
| `@ValidatePatch` | Presence-aware PATCH normalization/validation | GEN | PROPOSED / P2 |
| `@Error` | Resolve domain error definition from constants/error; not stringly typed | CHECK/WRAP | PROPOSED / P2 |
| `@StrictJSON` | Reject unknown properties at transport boundary | WRAP | OPTIONAL / P3 |

### Persistence, repository, transactions and search

| Directive / usage | Intended semantics | Mechanism | Status and roadmap |
| --- | --- | --- | --- |
| `@Repository` | Repository implementation metadata and constructor wiring | REG | PROPOSED / P4 |
| `@CRUD` | Opt-in GORM basic CRUD only; no guessed aggregate persistence | GEN | OPTIONAL / P4 |
| `@Query` | Explicit typed query declaration; parameterized only | GEN/CHECK | OPTIONAL / P4 |
| `@FindBy` | Allowlisted trivial finder generation; fail on ambiguity | GEN | OPTIONAL / P4 |
| `@Transactional` | Generated wrapper using core/database/transaction | WRAP | PROPOSED / P4 |
| `@ReadOnly` | Read-only transaction intent; enforce at DB where possible | WRAP | OPTIONAL / P4 |
| `@OptimisticLock` | Version predicate plus conflict mapping | GEN | OPTIONAL / P4 |
| `@TenantScoped` | Require trusted tenant context in repository predicates | CHECK/WRAP | PROPOSED / P7 |
| `@SoftDelete` | Optional GORM soft-delete integration | GEN | OPTIONAL / P4 |
| `@Searchable` | Declare allowed search fields, not raw SQL expressions | REG/CHECK | PROPOSED / P4 |
| `@Filterable` | Allowed filters/operators with canonicalization | REG/CHECK | OPTIONAL / P4 |
| `@Sortable` | Allowed stable sort keys and null ordering | REG/CHECK | OPTIONAL / P4 |
| `@Cursor` | Typed cursor serialization/validation policy | GEN/CHECK | OPTIONAL / P4 |
| `@Migration` | Only mark migration files for linting; never auto-apply DDL | CHECK | DISCOURAGED / P4 |

### Services, composition and dependency injection

| Directive / usage | Intended semantics | Mechanism | Status and roadmap |
| --- | --- | --- | --- |
| `@Service` | Register service contract/constructor for generated composition | REG | PROPOSED / P5 |
| `@Component` | Generic registered component when no narrower role | REG | OPTIONAL / P5 |
| `@Inject` | Constructor parameter binding, not reflection field injection | REG | OPTIONAL / P5 |
| `@Provides` | Explicit provider function for external dependencies | REG | PROPOSED / P5 |
| `@Named` | Qualify dependencies by stable identifier | REG | OPTIONAL / P5 |
| `@Scope` | Lifetime metadata; default service-instance scope | REG | OPTIONAL / P5 |
| `@PostConstruct` | Explicit startup hook; ordered and failure-aware | REG | OPTIONAL / P5 |
| `@PreDestroy` | Explicit close/shutdown hook with deadlines | REG | OPTIONAL / P5 |
| `@FeatureFlag` | Register typed feature gate; not authorization | WRAP/CHECK | OPTIONAL / P7 |
| `@Idempotent` | Apply declared key scope/result/retention rules at service boundary | WRAP | PROPOSED / P4 |

### HTTP controllers, routes, request/response and API docs

| Directive / usage | Intended semantics | Mechanism | Status and roadmap |
| --- | --- | --- | --- |
| `@Controller` | Register typed HTTP controller/adapter | REG | PROPOSED / P3 |
| `@Routes` | Controller/resource path prefix and grouping | REG | PROPOSED / P3 |
| `@GET` | GET route with typed path and handler contract | GEN | PROPOSED / P3 |
| `@POST` | POST route with typed handler contract | GEN | PROPOSED / P3 |
| `@PUT` | PUT route with replacement semantics | GEN | PROPOSED / P3 |
| `@PATCH` | PATCH route preserving field presence | GEN | PROPOSED / P3 |
| `@DELETE` | DELETE route | GEN | PROPOSED / P3 |
| `@PathParam` | Typed path extraction with decoding/validation | GEN | OPTIONAL / P3 |
| `@QueryParam` | Typed query extraction, limits and defaults | GEN | OPTIONAL / P3 |
| `@HeaderParam` | Typed header extraction; do not trust user/tenant IDs | GEN | OPTIONAL / P3 |
| `@RequestBody` | Strict typed decoding/size bound | GEN | PROPOSED / P3 |
| `@Response` | Typed response serializer using shared envelope | GEN | PROPOSED / P3 |
| `@Status` | Declared success status and error mapping | GEN/CHECK | OPTIONAL / P3 |
| `@Consumes` | Request media-type constraints | CHECK | OPTIONAL / P3 |
| `@Produces` | Response content-type constraints | CHECK | OPTIONAL / P3 |
| `@OpenAPI` | Generate or validate OpenAPI from contract source | CHECK/GEN | PROPOSED / P3 |
| `@Schema` | DTO/OpenAPI schema metadata checked against authoritative public schema | CHECK/GEN | PROPOSED / P3 |
| `@OperationID` | Stable API operation identity | CHECK | OPTIONAL / P3 |
| `@Example` | API examples with contract validation | CHECK | OPTIONAL / P3 |
| `@Deprecated` | Deprecation metadata and compatibility warning | CHECK | OPTIONAL / P3 |
| `@Public` | Explicit public endpoint exemption, security-reviewed | CHECK | PROPOSED / P7 |

### HTTP middleware and cross-cutting web concerns

| Directive / usage | Intended semantics | Mechanism | Status and roadmap |
| --- | --- | --- | --- |
| `@UseMiddleware` | Ordered middleware bindings from approved registry | REG | PROPOSED / P3 |
| `@RequestID` | Generate/propagate unique request ID | WRAP | PROPOSED / P3 |
| `@CorrelationID` | Propagate correlation ID distinct from request ID | WRAP | PROPOSED / P3 |
| `@Trace` | OpenTelemetry instrumentation and context propagation | WRAP | PROPOSED / P6 |
| `@Log` | Structured request/operation log policy; redact by default | WRAP | PROPOSED / P6 |
| `@Recover` | Panic-to-safe-HTTP response, no suppressed observability | WRAP | PROPOSED / P3 |
| `@Timeout` | Deadline enforcement with cancellation-aware code | WRAP | PROPOSED / P3 |
| `@RateLimit` | Tenant/actor-scoped limits and bounded state | WRAP | PROPOSED / P7 |
| `@CORS` | Approved-origin policy, not permissive wildcard default | WRAP | OPTIONAL / P7 |
| `@MaxBody` | Request size limit before decoding | WRAP | PROPOSED / P3 |
| `@SecureHeaders` | Security headers / transport policy | WRAP | OPTIONAL / P7 |
| `@CacheControl` | HTTP caching semantics (private/no-store as needed) | WRAP | OPTIONAL / P7 |
| `@CSRF` | Browser session/unsafe-request protection where applicable | WRAP | OPTIONAL / P7 |
| `@Compress` | Compression with content-type and side-channel checks | WRAP | OPTIONAL / P7 |

### Security, secrets, tenancy and policy

| Directive / usage | Intended semantics | Mechanism | Status and roadmap |
| --- | --- | --- | --- |
| `@Authenticate` | Require verified identity; never trust caller-filled principal | WRAP/CHECK | PROPOSED / P7 |
| `@Authorize` | Evaluate business permission against trusted actor/context | WRAP | PROPOSED / P7 |
| `@RolesAllowed` | Role-based rule using central IAM policy | WRAP | OPTIONAL / P7 |
| `@RequiresPermission` | Named policy requirement with resource context | WRAP | PROPOSED / P7 |
| `@OwnerOnly` | Ownership check must resolve resource via service safely | WRAP | OPTIONAL / P7 |
| `@TenantScoped` | Enforce tenant context across route/service/repository | CHECK/WRAP | PROPOSED / P7 |
| `@Sensitive` | Data classification for logs/contracts/events | CHECK | PROPOSED / P7 |
| `@Redact` | Redaction of diagnostics/log output | WRAP/GEN | PROPOSED / P7 |
| `@Secret` | Config/credential never exposed in diagnostics | CHECK | PROPOSED / P5 |
| `@Encrypt` | Encryption policy for stored fields, provider-owned crypto | WRAP | OPTIONAL / P7 |
| `@Audit` | Emit bounded, access-controlled audit record | WRAP | PROPOSED / P6 |

### Typed HTTP clients, provider adapters and external services

| Directive / usage | Intended semantics | Mechanism | Status and roadmap |
| --- | --- | --- | --- |
| `@Client` | Typed in-process or HTTP client descriptor / adapter | GEN/REG | PROPOSED / P3 |
| `@ExternalService` | Explicit third-party boundary and provider ownership | REG | PROPOSED / P5 |
| `@Provider` | Provider identity, capabilities and factory registration | REG | PROPOSED / P5 |
| `@Adapter` | Bind vendor implementation to stable provider interface | REG | PROPOSED / P5 |
| `@Endpoint` | Typed base URL/route metadata from configuration | REG | PROPOSED / P5 |
| `@HTTPMethod` | Outgoing request verb/path mapping | GEN | OPTIONAL / P5 |
| `@OAuth2` | OAuth flow configuration metadata, no hidden token grant | REG | OPTIONAL / P5 |
| `@Credential` | Scoped credential reference/injection; redact | REG | PROPOSED / P5 |
| `@Retry` | Bounded idempotent retry policy, explicit status set | WRAP | PROPOSED / P5 |
| `@Backoff` | Deterministic backoff with jitter/cap | WRAP | PROPOSED / P5 |
| `@CircuitBreaker` | Bounded failure protection per provider/tenant | WRAP | OPTIONAL / P5 |
| `@Bulkhead` | Limit concurrent provider calls | WRAP | OPTIONAL / P5 |
| `@Fallback` | Typed fallback function for safe error types | WRAP | OPTIONAL / P5 |
| `@Webhook` | Webhook route/handler registration metadata | REG | PROPOSED / P5 |
| `@VerifySignature` | Invoke provider-specific authenticity verifier | WRAP | PROPOSED / P5 |
| `@HealthProbe` | Explicit provider health check with timeout | REG | OPTIONAL / P5 |

### Events, envelopes and event lifecycle

| Directive / usage | Intended semantics | Mechanism | Status and roadmap |
| --- | --- | --- | --- |
| `@Event` | Domain-owned topic, version and payload metadata | REG/CHECK | PROPOSED / P6 |
| `@EventVersion` | Stable event schema version and compatibility checks | CHECK | PROPOSED / P6 |
| `@EventHandler` | Typed event handler registration | REG | PROPOSED / P6 |
| `@PublishEvent` | Explicit event publisher decorator metadata only | REG/WRAP | OPTIONAL / P6 |
| `@AfterCommit` | Bind to existing after-commit callback semantics | WRAP | PROPOSED / P6 |
| `@Outbox` | Opt-in durable outbox persistence/relay contract | REG | OPTIONAL / P6 |
| `@EventSchema` | Contract/schema compatibility checks | CHECK | OPTIONAL / P6 |
| `@EventSource` | Domain source metadata, not envelope field | REG | OPTIONAL / P6 |
| `@Causation` | Propagate existing command/parent event ID | WRAP | PROPOSED / P6 |
| `@Correlation` | Propagate existing correlation/trace context | WRAP | PROPOSED / P6 |
| `@EventPII` | Classify full-snapshot events and retention policy | CHECK | PROPOSED / P7 |

### Messaging producer, consumer, broker topology and delivery

| Directive / usage | Intended semantics | Mechanism | Status and roadmap |
| --- | --- | --- | --- |
| `@Producer` | Typed publisher binding to MessagingClient, domain topic | REG/GEN | PROPOSED / P6 |
| `@Consumer` | Typed consumer registration using MessagingClient | REG | PROPOSED / P6 |
| `@Subscribe` | Stable topic/subject and handler association | REG | PROPOSED / P6 |
| `@Subject` | Typed subject metadata; no per-event randomized subjects | REG/CHECK | PROPOSED / P6 |
| `@Stream` | Stream configuration owned by Infrastructure provisioning | REG/CHECK | OPTIONAL / P6 |
| `@Durable` | Stable durable consumer identifier | REG | PROPOSED / P6 |
| `@QueueGroup` | Consumer work-sharing group metadata | REG | OPTIONAL / P6 |
| `@Ack` | Explicit ack policy and success handling | WRAP | PROPOSED / P6 |
| `@Nack` | Explicit failure retry or reject behavior | WRAP | OPTIONAL / P6 |
| `@RetryPolicy` | Consumer delivery attempt/backoff constraints | REG | PROPOSED / P6 |
| `@MaxDeliver` | Upper bound on attempts; poison advisory handling | REG | PROPOSED / P6 |
| `@DeadLetter` | Explicit dead-letter/recovery workflow; not assumed native DLQ | REG | OPTIONAL / P6 |
| `@Concurrency` | Consumer worker limits/backpressure | REG | PROPOSED / P6 |
| `@Batch` | Batch processing/ack semantics | REG | OPTIONAL / P6 |
| `@IdempotentConsumer` | Deduplication key and state-store contract | WRAP | PROPOSED / P6 |
| `@MessageTrace` | W3C trace context publish/extract | WRAP | PROPOSED / P6 |
| `@MessageTimeout` | Processing timeout / ack wait coordination | WRAP | PROPOSED / P6 |

### Configuration, profiles, typed environment and secrets

| Directive / usage | Intended semantics | Mechanism | Status and roadmap |
| --- | --- | --- | --- |
| `@Configuration` | Typed config struct registered for startup loading | REG/GEN | PROPOSED / P5 |
| `@ConfigPrefix` | Namespace environment properties by service | GEN | PROPOSED / P5 |
| `@Env` | Bind explicitly named environment variable | GEN | PROPOSED / P5 |
| `@Default` | Typed default with explicit precedence | GEN | PROPOSED / P5 |
| `@Required` | Configuration missing error at startup | CHECK | PROPOSED / P5 |
| `@ValidateConfig` | Run technical config validation before module startup | WRAP | PROPOSED / P5 |
| `@Secret` | Mark sensitive config redaction and source restrictions | CHECK | PROPOSED / P5 |
| `@ConfigSource` | Describe allowed source (env, file, secret store) | REG | OPTIONAL / P5 |
| `@Profile` | Explicit development/test/production selection | REG/CHECK | OPTIONAL / P5 |
| `@Reloadable` | Controlled reload with atomic switch and validation | WRAP | OPTIONAL / P5 |
| `@FeatureFlag` | Typed feature selection; not a security boundary | REG | OPTIONAL / P5 |

### Logging, tracing, metrics, auditing and diagnostics

| Directive / usage | Intended semantics | Mechanism | Status and roadmap |
| --- | --- | --- | --- |
| `@Log` | Structured log wrapper through `slog`/core logging | WRAP | PROPOSED / P7 |
| `@Trace` | OpenTelemetry span boundary with context | WRAP | PROPOSED / P7 |
| `@Span` | Named child span with attributes allowlist | WRAP | OPTIONAL / P7 |
| `@Timed` | Histogram for operation duration with bounded labels | WRAP | PROPOSED / P7 |
| `@Counted` | Counter for stable operation/result labels | WRAP | OPTIONAL / P7 |
| `@Metric` | Register stable metric identity and unit | REG | OPTIONAL / P7 |
| `@Audit` | Security/business audit without secret exposure | WRAP | PROPOSED / P7 |
| `@Redact` | Hide fields from logs, traces and errors | CHECK/WRAP | PROPOSED / P7 |
| `@NoLog` | Suppress payload logging, preserve audit/error metadata | CHECK | PROPOSED / P7 |
| `@SLO` | Attach objective metadata to monitored endpoint/operation | REG | OPTIONAL / P7 |
| `@CorrelationID` | Carry correlation field and trace metadata | WRAP | PROPOSED / P7 |

### Server, launcher, startup, shutdown and background jobs

| Directive / usage | Intended semantics | Mechanism | Status and roadmap |
| --- | --- | --- | --- |
| `@Module` | Declare module composition entry, dependencies and exports | REG | PROPOSED / P3 |
| `@Server` | HTTP server policy (bind address, limits, TLS) | REG/CHECK | OPTIONAL / P3 |
| `@HealthCheck` | Register liveness probe function | REG | PROPOSED / P3 |
| `@ReadyCheck` | Register dependencies/readiness policy | REG | PROPOSED / P3 |
| `@Startup` | Register ordered initializer and failure behavior | REG | OPTIONAL / P3 |
| `@GracefulShutdown` | Stop accepting requests, drain work, close dependencies | REG | PROPOSED / P3 |
| `@Shutdown` | Module close handler and deadline | REG | PROPOSED / P3 |
| `@Scheduled` | Schedule registered function with timezone/cadence | REG | OPTIONAL / P8 |
| `@Cron` | Cron expression with explicit timezone/misfire handling | REG | OPTIONAL / P8 |
| `@LeaderOnly` | Distributed lock/lease requirement for singleton jobs | WRAP | OPTIONAL / P8 |
| `@BackgroundJob` | Registered worker managed by lifecycle, not hidden goroutine | REG | OPTIONAL / P8 |

### Testing, fixtures, mocks, contract tests and benchmarks

| Directive / usage | Intended semantics | Mechanism | Status and roadmap |
| --- | --- | --- | --- |
| `@UnitTest` | Test category metadata for Ginkgo suite grouping | CHECK | OPTIONAL / P8 |
| `@IntegrationTest` | Service/database/broker integration category | CHECK | OPTIONAL / P8 |
| `@E2E` | Full HTTP/worker flow category | CHECK | OPTIONAL / P8 |
| `@Fixture` | Typed deterministic fixture constructor | GEN | PROPOSED / P8 |
| `@Mock` | Generate interface fake/mock without expectation magic | GEN | PROPOSED / P8 |
| `@GenerateMock` | Explicit mock target and stable generated location | GEN | PROPOSED / P8 |
| `@ContractTest` | Cross-module OpenAPI/client/event contract test manifest | REG/GEN | PROPOSED / P8 |
| `@PropertyTest` | Property/fuzz test input metadata | GEN | OPTIONAL / P8 |
| `@TestMatrix` | Explicit table-driven variations | GEN | OPTIONAL / P8 |
| `@TestTimeout` | Timeout for test suite/fixture | REG | OPTIONAL / P8 |
| `@RequiresDocker` | Skip/fail rationale when containers unavailable | CHECK | OPTIONAL / P8 |
| `@Benchmark` | Register performance benchmark and baseline | REG | OPTIONAL / P8 |
| `@Seed` | Deterministic seed/clock and resource cleanup | REG | OPTIONAL / P8 |
| `@Golden` | Golden-data drift check with explicit update step | CHECK | OPTIONAL / P8 |

### Utilities, helpers, cache, concurrency and shared contracts

| Directive / usage | Intended semantics | Mechanism | Status and roadmap |
| --- | --- | --- | --- |
| `@Pure` | Static advisory for side-effect-free functions; lint where possible | CHECK | OPTIONAL / P8 |
| `@Deprecated` | Expose deprecation diagnostics across call sites | CHECK | PROPOSED / P8 |
| `@Experimental` | Mark unstable API and compatibility expectations | CHECK | OPTIONAL / P8 |
| `@Memoize` | Cache only proven pure, bounded and safe functions | WRAP | DISCOURAGED / P8 |
| `@Cacheable` | Explicit cache keys, tenant partition, TTL, invalidation | WRAP | OPTIONAL / P8 |
| `@CacheEvict` | Invalidate declared cache keys after successful mutation | WRAP | OPTIONAL / P8 |
| `@CacheKey` | Stable typed key fields, no secrets | GEN/CHECK | OPTIONAL / P8 |
| `@TTL` | Retention/expiry declaration with time unit | REG | OPTIONAL / P8 |
| `@ConcurrencySafe` | Document verification obligation; not an automatic race fix | CHECK | DISCOURAGED / P8 |
| `@Bulkhead` | Concurrency semaphore for external work | WRAP | OPTIONAL / P5 |
| `@Checksum` | Explicit integrity check for files/messages | WRAP | OPTIONAL / P8 |
| `@Retention` | Data lifecycle metadata; requires actual cleanup process | REG/CHECK | OPTIONAL / P8 |
| `@UploadLimit` | Enforce file request/storage limit | WRAP | OPTIONAL / P8 |
| `@Encrypt` | Bounded encryption policy via core/crypto if added | WRAP | OPTIONAL / P8 |
| `@GenerateStringer` | Typed enum stringer when needed | GEN | OPTIONAL / P8 |

## All 23 GoF pattern annotation decisions

| Family | Pattern | Annotation | Permitted use | Decision |
| --- | --- | --- | --- | --- |
| Creational | Abstract Factory | `@AbstractFactory` | Registry of related provider families; explicit factory contract | Optional registry |
| Creational | Builder | `@Builder` | Typed fluent construction | Supported |
| Creational | Factory Method | `@FactoryMethod` | Register constructors for a capability | Useful registry |
| Creational | Prototype | `@Prototype` | Explicit clone/copy contract only | Avoid default auto-clone |
| Creational | Singleton | `@Singleton` | Lifetime constraint in composition root | Avoid global instance |
| Structural | Adapter | `@Adapter` | Bind vendor to stable interface | Useful registry |
| Structural | Bridge | `@Bridge` | Map abstraction to independent implementor | Optional metadata |
| Structural | Composite | `@Composite` | Declare child composition, no hidden recursion | Optional metadata |
| Structural | Decorator | `@Decorator` | Explicit typed wrapper order | Useful generated wiring |
| Structural | Facade | `@Facade` | Mark stable public capability API | Check only |
| Structural | Flyweight | `@Flyweight` | Object sharing semantics require manual proof | Avoid generation |
| Structural | Proxy | `@Proxy` | Typed forwarding for client/resilience boundary | Useful generation |
| Behavioral | Chain of Responsibility | `@Chain` | Ordered middleware/handler chain | Useful registry |
| Behavioral | Command | `@Command` | Typed command registry and dispatch | Optional registry |
| Behavioral | Interpreter | `@Interpreter` | Registry of parsing grammar components | Avoid auto-parser |
| Behavioral | Iterator | `@Iterator` | Typed iteration interface only when needed | Avoid generation |
| Behavioral | Mediator | `@Mediator` | Explicit event/notification coordination | Optional registry |
| Behavioral | Memento | `@Memento` | Snapshot semantics with sensitive-state controls | Avoid automation |
| Behavioral | Observer | `@Observer` | Listener registration around event contract | Useful registry |
| Behavioral | State | `@State` | Register state transitions with guards | Check only initially |
| Behavioral | Strategy | `@Strategy` | Provider capability registration + selection | Useful registry |
| Behavioral | Template Method | `@TemplateMethod` | Common skeleton remains an explicit Go method | Avoid inheritance emulation |
| Behavioral | Visitor | `@Visitor` | Typed dispatch with explicit exhaustive cases | Optional, rare |

No automatic inheritance for Template Method, global Singleton construction, unsafe Prototype clone, Flyweight sharing, or business-state transition inference. `@Strategy`, `@Adapter`, `@Decorator`, `@FactoryMethod`, `@Proxy`, `@Chain`, and `@Observer` are candidates for *explicit* future registries/wrappers.

## Phase 1 files and ownership

```text
backend/libs/core/annotation/
  annotation.go                     # contract types, registry, comment validation
  annotations_{entity,mapping,validation,persistence,composition,http}.go
  annotations_{middleware,security,external,events,messaging,configuration}.go
  annotations_{observability,server,testing,utilities}.go
  inventory.go                      # combines 16 source catalogs
  patterns.go                       # all 23 GoF decisions
  entity/entity.go                  # pre-existing generic builder
  generation/{expand,gormtable,overlay}/ # pre-existing tooling
  generation/parser/parser.go       # read-only AST discovery (new)
backend/libs/core/tests/unit-test/annotation/ # contract/scanner tests
```

These folders are metadata/tool ownership; they do not change the public package boundaries of `core/events`, `core/http`, `core/database`, `core/logging`, `core/telemetry`, or any service.

## Phase 2 migration roadmap (not started)

| Workstream | Required verification before application |
| --- | --- |
| Mapping, validation and normalization | Type-safe source/destination mapping, Unicode rules, PATCH presence, field errors, security exclusions |
| HTTP routes/controllers/middleware | Stable OpenAPI, mux path rules, middleware order, authenticated principal, trace and error envelopes |
| Configuration/DI/external adapters | No reflection-based global service locator, validated secret handling, provider contracts and capability selection |
| Repository/persistence/transactions | Preserve GORM predicates, tenant filters, transaction context, conflict errors and idempotent writes |
| Events/messaging/producers/consumers | Stable envelope/topic/version, no duplicate idempotent publication, after-commit behavior, NATS ack/DLQ safety |
| Observability/security/server/jobs | Secure logging/redaction, bounded telemetry labels, shutdown deadlines, no unsafe job re-entry |
| Testing/utility contracts | Existing Ginkgo tests and deterministic fixtures; benchmark overhead and clean-checkout build |

At every migration step, enable an annotation **only after** its processor exists and tests prove it performs what the directive declares. Source comments with only Phase 1 definitions must not be introduced in service source as if they enforce runtime behavior.

## Completion criteria

- Every documented row has a discoverable typed `UseCase` and all 23 GoF decisions have `Pattern` metadata.
- Parsing is deterministic, unknown directives fail closed, registry callers cannot mutate shared definitions.
- Existing Foundation code generation stays intact; no existing service code has been rewritten.
- Phase 2 is explicitly separate; only implement or migrate it after separate authorization.


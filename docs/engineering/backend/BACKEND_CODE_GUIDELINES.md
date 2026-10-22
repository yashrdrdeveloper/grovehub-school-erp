# Backend Code Guidelines

> **Status:** Mandatory Phase 0 engineering standard  
> **Applies to:** All Go backend code, services, public contracts, clients, persistence, messaging, provider/capability adapters, workers, launchers, tests, generated code, and backend shared libraries  
> **Primary language:** Go  
> **Architecture:** Modular monolith with microservice-grade boundaries  
> **Initial runtime:** One repository, one binary, one process, one shared HTTP/gRPC runtime  
> **Evolution target:** Selected modules may later become independently deployable services without rewriting their business logic  
> **Core principles:** SRP, OCP, DRY, explicit contracts, constructor injection, small interfaces, capability-based infrastructure, testability, observability, security, and backward compatibility

---

# 1. Objective

The backend must be designed so that it is cheap to run as one application today and safe to split into independent services later only when a real need exists.

The architecture must optimize for:

1. **Single Responsibility Principle (SRP)**
2. **Open/Closed Principle (OCP)**
3. **Don't Repeat Yourself (DRY)**
4. Strong module/service ownership
5. Explicit public contracts
6. No cross-service implementation coupling
7. Infrastructure behind capability interfaces
8. Constructor-based dependency injection
9. Small interfaces
10. Pure mapping and validation boundaries
11. Explicit transaction ownership
12. Stable error contracts
13. Observable execution
14. Safe concurrency
15. Backward-compatible APIs and events
16. Fast local development
17. Clone-and-run repository experience
18. Architecture rules that are enforceable by CI and AI coding agents

The backend is not organized to create folders. It is organized to make responsibility, ownership, dependency direction, compatibility, testing, and future extraction obvious.

---

# 2. Core Architecture Philosophy

The system is:

> **A modular monolith with explicit versioned service boundaries and microservice-grade extraction seams.**

Initially:

```text
ONE REPOSITORY
      ↓
ONE GO BINARY
      ↓
ONE PROCESS
      ↓
ONE SHARED HTTP/gRPC RUNTIME
      ↓
ROOT ROUTER / COMPOSITION ROOT
      ↓
IN-PROCESS SERVICES
```

We do **not** introduce localhost HTTP/gRPC calls simply to imitate microservices.

Strong service boundaries exist even while services are in the same process.

If Service A needs Service B, Service A consumes Service B's **public client contract**, not Service B's internal packages.

Later:

```text
Service A
   ↓
Service B client/vN
   ↓
transport strategy
   ├── inprocess
   ├── grpc
   └── http
   ↓
Service B
```

Business consumers should not care whether the target service is local or remote.

---

# 3. Non-Negotiable Backend Rules

1. Every service owns its domain/business logic.
2. Every service exposes one clear root construction boundary such as `New(...)`.
3. The root launcher does not construct service-internal repositories/controllers/business services itself.
4. Public contracts are versioned.
5. Internal business logic is **not** versioned merely because HTTP/gRPC gets a new major version.
6. A service never imports another service's internal implementation (`src/`, internal repository, infrastructure, etc.).
7. A service never reads or writes another service's database directly.
8. Cross-service access goes through the target service's public client.
9. Controllers contain no business logic.
10. Repositories contain no business workflows.
11. Mappers contain no business decisions.
12. Request validators contain no business authorization or repository queries.
13. Generated files are never hand-edited.
14. Protobuf is a transport contract, not a domain model.
15. OpenAPI is the HTTP contract.
16. Shared technical primitives belong in explicit `libs/core/<capability>` packages.
17. Shared domain/business concepts do not go into `libs/core`.
18. Circular service dependencies are forbidden.
19. Constructor injection is preferred.
20. Global mutable dependencies are forbidden.
21. Service locators are forbidden in business code.
22. Small interfaces are preferred over god interfaces.
23. Composition is preferred over inheritance-like/base-service abstractions.
24. New abstractions must solve a real current problem.
25. `common`, `misc`, giant `utils`, and giant `helpers` dumping grounds are forbidden.
26. Infrastructure/provider selection must not leak into business logic.
27. `switch`/vendor branching should not be the mechanism for extensibility.
28. Retry is allowed only when safety/idempotency semantics are understood.
29. Errors crossing public boundaries must be normalized.
30. Every production goroutine has explicit lifecycle/cancellation ownership.
31. Every externally visible contract change requires compatibility review.
32. Every schema change requires a migration.
33. Important domain invariants should be protected at both service and database layers when practical.
34. Observability is part of implementation, not an afterthought.
35. Secrets must never be logged or returned accidentally.

---

# 4. Repository Layout

Recommended backend structure inside the product monorepo:

```text
backend/
│
├── go.mod
├── go.sum
├── go.work                    # only when multiple modules are intentionally used
├── Makefile
├── README.md
├── .env.example
├── buf.yaml
├── buf.gen.yaml
├── golangci.yml
│
├── cmd/
│   └── platform/
│       └── main.go
│
├── launcher/
│   ├── launcher.go
│   ├── bootstrap/
│   │   ├── config.go
│   │   ├── logger.go
│   │   └── telemetry.go
│   ├── composition/
│   │   ├── services.go
│   │   ├── providers.go
│   │   └── dependencies.go
│   ├── routing/
│   │   ├── http.go
│   │   ├── grpc.go
│   │   └── registry.go
│   ├── lifecycle/
│   │   ├── start.go
│   │   ├── stop.go
│   │   └── shutdown.go
│   ├── health/
│   │   ├── health.go
│   │   └── readiness.go
│   └── runtime/
│       ├── http.go
│       └── grpc.go
│
├── services/
│   ├── foundation/
│   ├── identity/
│   ├── catalog/
│   ├── execution/
│   └── ...
│
├── libs/
│   └── core/
│       ├── config/
│       ├── logging/
│       ├── database/
│       ├── cache/
│       ├── http/
│       ├── grpc/
│       ├── messaging/
│       ├── telemetry/
│       ├── health/
│       ├── resilience/
│       ├── validation/
│       ├── pagination/
│       ├── security/
│       ├── id/
│       ├── clock/
│       ├── lifecycle/
│       ├── concurrency/
│       ├── serialization/
│       ├── errors/
│       ├── capability/
│       └── provider/
│
├── providers/
│   ├── mock/
│   ├── memory/
│   ├── filesystem/
│   └── reference/
│
├── testing/
│   └── unit-test/
│       ├── launcher/          # mirrors backend/launcher/
│       ├── services/          # mirrors backend/services/
│       ├── libs/              # mirrors backend/libs/
│       ├── providers/         # mirrors backend/providers/
│       └── common/            # genuinely shared test-only support
│
├── deployments/
├── scripts/
└── docs/
    └── engineering/
```

Do not create folders that have no current responsibility.

---

# 5. Standard Service Blueprint

Each backend service/module follows the same vocabulary.

```text
services/<service>/
│
├── <service>.go
│
├── launcher/
│   ├── launcher.go
│   ├── bootstrap/
│   ├── wiring/
│   │   ├── repositories.go
│   │   ├── services.go
│   │   ├── controllers.go
│   │   ├── messaging.go
│   │   └── clients.go
│   ├── routes/
│   │   ├── http.go
│   │   └── grpc.go
│   ├── lifecycle/
│   │   ├── start.go
│   │   └── stop.go
│   └── health/
│       └── health.go
│
├── config/
│   ├── config.go
│   ├── defaults.go
│   ├── loader.go
│   └── validator.go
│
├── api/
│   ├── proto/v1/
│   └── openapi/v1/openapi.yaml
│
├── dto/
│   └── v1/
│       ├── request/
│       ├── response/
│       └── common/
│
├── client/
│   └── v1/
│       ├── client.go
│       ├── options.go
│       ├── errors.go
│       ├── inprocess/
│       ├── grpc/
│       └── http/
│
├── gen/
│   ├── proto/v1/
│   └── openapi/v1/
│
├── src/
│   ├── entity/
│   ├── services/
│   ├── repository/
│   ├── ports/
│   │   ├── inbound/
│   │   └── outbound/
│   ├── controller/
│   │   ├── http/v1/
│   │   └── grpc/v1/
│   ├── mapper/v1/
│   ├── validator/v1/
│   ├── middleware/
│   ├── messaging/
│   │   ├── events/v1/
│   │   ├── producer/
│   │   ├── consumer/
│   │   └── outbox/
│   ├── infrastructure/
│   │   ├── persistence/
│   │   ├── cache/
│   │   ├── clients/
│   │   ├── storage/
│   │   ├── crypto/
│   │   └── telemetry/
│   ├── errors/
│   ├── constants/
│   ├── helpers/
│   └── utils/
│
├── migrations/
├── scripts/
├── docs/
└── mocks/

# Unit tests for this service do not live here.
# They mirror this service under:
# backend/testing/unit-test/services/<service>/
```

Use only folders required by the service.

Service-local production folders do not contain the canonical backend Go unit-test tree.

Tests mirror the service under:

```text
backend/testing/unit-test/services/<service>/
```

Example:

```text
backend/services/catalog/src/services/application/
        ↓
backend/testing/unit-test/services/catalog/src/services/application/
```

Repository-wide sanity/regression live outside the backend module:

```text
testing/sanity/
testing/regression/
```

---

# 6. Public vs Internal Surface

Public/compatibility-sensitive surface:

```text
api/
dto/
client/
src/messaging/events/
```

Internal implementation:

```text
src/
```

Construction:

```text
launcher/
```

Configuration:

```text
config/
```

Generated output:

```text
gen/
```

Schema evolution:

```text
migrations/
```

Only public surfaces should be imported by other services.

---

# 7. Service Root Facade

Each service exposes a small root facade.

Example:

```go
package catalog

type Module interface {
    RegisterHTTP(router Router)
    RegisterGRPC(server GRPCServer)
    Health(ctx context.Context) Health
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
}

func New(opts Options) (Module, error) {
    // service-owned construction
}
```

Root composition may import:

```go
services/catalog
```

It must not import:

```text
services/catalog/src/services
services/catalog/src/repository
services/catalog/src/infrastructure
```

---

# 8. Service Construction Rule

`New(...)` / service launcher owns internal wiring.

Construction may include:

1. validate service config
2. create infrastructure adapters
3. construct repositories
4. construct outbound clients
5. construct business services
6. construct mappers/validators
7. construct controllers
8. construct producers/consumers
9. construct background workers
10. build health/readiness checks
11. expose route registration
12. return the module facade

The root launcher must not reproduce service internals.

---

# 9. Dependency Direction

Required direction:

```text
HTTP/gRPC Controller
        ↓
Inbound Port / Business Service
        ↓
Domain Entity
        ↓
Repository / Outbound Port Interface
        ↑
Infrastructure Adapter
```

Messaging:

```text
Consumer
   ↓
Inbound Port / Business Service
```

Events:

```text
Business Service
   ↓
EventPublisher interface
   ↑
Producer Adapter
```

Cross-service:

```text
Service A
   ↓
Service B client/vN
   ↓
inprocess OR grpc/http
   ↓
Service B
```

---

# 10. Forbidden Dependency Direction

```text
entity        ✗ controller
entity        ✗ protobuf
entity        ✗ HTTP
entity        ✗ gRPC
entity        ✗ SQL driver
entity        ✗ Redis/Valkey client

service       ✗ concrete PostgreSQL adapter
service       ✗ concrete cache adapter
service       ✗ raw HTTP/gRPC transport implementation
service       ✗ provider brand checks

controller    ✗ repository
controller    ✗ DB
controller    ✗ business rules

repository    ✗ controller
repository    ✗ business workflow

mapper        ✗ repository
mapper        ✗ network
mapper        ✗ business decision

validator     ✗ repository
validator     ✗ business authorization

producer      ✗ business decision
consumer      ✗ business decision

service A     ✗ service B/src/*
service A     ✗ service B/database
```

---

# 11. Entities and Value Objects

Location:

```text
src/entity/
```

Entities/value objects contain domain state and intrinsic domain behavior.

They must not depend on:

- HTTP
- gRPC
- Protobuf
- SQL
- Redis/Valkey
- generated transport code
- another service's internal implementation

Keep entities transport- and persistence-independent.

Value objects should be used when a primitive has important invariants or semantics.

Example:

```go
type ApplicationID string
```

Do not create wrapper types without benefit.

---

# 12. Business Services / Use Cases

Location:

```text
src/services/
```

Business services own:

- business rules
- lifecycle rules
- relationship validation
- hierarchy/cycle checks
- readiness calculations
- authorization decisions belonging to the domain
- orchestration across repositories
- orchestration across outbound ports/public service clients
- transaction boundaries
- event creation decisions
- reconciliation after changes
- semantic errors

Business services do not know concrete infrastructure vendors.

Correct:

```go
type Service struct {
    repo      repository.ApplicationRepository
    clock     clock.Clock
    publisher ports.EventPublisher
}
```

Incorrect:

```go
type Service struct {
    db    *sql.DB
    redis *redis.Client
}
```

unless the package itself is explicitly an infrastructure package.

---

# 13. Interfaces and Ports

Location:

```text
src/ports/inbound/
src/ports/outbound/
```

Use interfaces at meaningful boundaries.

Prefer small interfaces.

Good:

```go
type EventPublisher interface {
    Publish(ctx context.Context, event Event) error
}
```

Avoid giant interfaces that combine unrelated responsibilities.

Do not create an interface for every struct merely to claim SOLID compliance.

In Go, interfaces should generally be owned by the consumer or a neutral contract package when multiple consumers require the same stable capability.

---

# 14. Repository Contracts

Location:

```text
src/repository/
```

Repositories are persistence contracts.

Good:

```go
type ApplicationRepository interface {
    Create(ctx context.Context, application entity.Application) error
    Update(ctx context.Context, application entity.Application) error
    FindByID(ctx context.Context, id entity.ApplicationID) (entity.Application, error)
    Delete(ctx context.Context, id entity.ApplicationID) error
}
```

Repository methods should be persistence-oriented.

Prefer:

```text
Create
Update
Delete
FindByID
FindByKey
ListByProjectID
ExistsByName
```

Avoid business workflow names:

```text
Promote
Approve
CalculateReadiness
ValidateParent
ResolveRequirements
PerformMigration
```

Those belong to business services.

---

# 15. Persistence Implementations

Location:

```text
src/infrastructure/persistence/<technology>/
```

Examples:

```text
postgres/
mysql/
dynamodb/
```

Persistence adapters own:

- query construction
- SQL/database API calls
- scanning
- persistence row models if needed
- transaction mechanics
- DB-specific error translation
- persistence-specific serialization

They do not own business policy.

---

# 16. Database Ownership

Each service/domain owns its logical schema and migrations.

A service must not query another service's tables directly even if everything currently shares one PostgreSQL server/database.

Physical database sharing is an implementation detail.

Logical ownership is mandatory.

```text
Service A → Service A repository → Service A schema
Service B → Service B repository → Service B schema
```

Cross-service data needs go through public clients/contracts.

---

# 17. PostgreSQL Guidelines

PostgreSQL is the initial primary system of record.

Use:

- explicit primary keys
- foreign keys where ownership is known
- unique constraints for domain uniqueness
- check constraints for stable database invariants
- indexes based on actual query requirements
- JSONB for genuinely flexible metadata/configuration
- relational join tables for real relationships

Avoid:

- giant generic `things` tables
- speculative indexes
- native PostgreSQL ENUM by default
- duplicated derived/readiness state that can be safely calculated
- cross-domain tables with unclear ownership

Use Go-defined string constants for controlled values unless the database itself genuinely needs a native enum.

---

# 18. Migrations

Every schema change requires migration files.

Recommended pattern:

```text
migrations/postgres/
├── 000001_create_applications.up.sql
├── 000001_create_applications.down.sql
├── 000002_add_application_key.up.sql
└── 000002_add_application_key.down.sql
```

Rules:

- migrations are ordered and immutable after release
- test migrations from clean database
- test upgrade from supported previous version
- provide rollback/down migration when technically safe
- destructive changes require staged rollout when compatibility matters
- add constraints intentionally
- data migrations must be bounded and observable
- avoid unbounded table rewrites during production startup

For large production changes, prefer expand → migrate/backfill → switch → contract.

---

# 19. Transaction Ownership

Business service/use-case layer owns multi-operation transactions.

Example:

```go
func (s *Service) ReplaceMappings(ctx context.Context, req ReplaceMappingsRequest) error {
    return s.tx.Within(ctx, func(ctx context.Context) error {
        if err := s.mappingRepo.DeleteByModuleID(ctx, req.ModuleID); err != nil {
            return err
        }

        return s.mappingRepo.CreateMany(ctx, req.Mappings)
    })
}
```

Repositories may participate in a transaction but should not decide business transaction scope.

Do not open nested transactions casually.

---

# 20. DTO Rules

Public DTOs are separate from entities and DB models.

```text
dto/v1/
├── request/
├── response/
└── common/
```

Rules:

- DTOs are versioned
- entities are not returned directly
- request DTOs expose only client-controlled fields
- response DTOs expose intentional public data
- IDs/timestamps owned by server are not accepted from clients unless the contract explicitly requires it
- DTOs contain no business logic
- DTO compatibility is preserved within a major version

---

# 21. API Versioning

Version public/compatibility-sensitive surfaces:

```text
api/proto/v1/
api/openapi/v1/
dto/v1/
client/v1/
src/controller/http/v1/
src/controller/grpc/v1/
src/mapper/v1/
src/validator/v1/
src/messaging/events/v1/
```

Do not version internal business logic by API version:

```text
src/entity/
src/services/
src/repository/
src/ports/
src/infrastructure/
launcher/
config/
```

Correct model:

```text
HTTP v1 ─┐
         ├── current business service
HTTP v2 ─┘
```

Do not create duplicate v1/v2 business services merely because the transport changes.

---

# 22. Public Client Boundary

Every consumable service exposes:

```text
client/vN/
```

Example:

```go
type Client interface {
    GetApplication(ctx context.Context, req GetApplicationRequest) (GetApplicationResponse, error)
}
```

Implementations:

```text
client/v1/inprocess/
client/v1/grpc/
client/v1/http/
```

Internal consumer code depends on the client interface, not transport.

---

# 23. Cross-Service Rule

If Service A requires Service B:

```text
Service A
   ↓
Service B client/vN
   ↓
Service B
```

Never:

```text
Service A → Service B/src
Service A → Service B repository
Service A → Service B DB
Service A → Service B infrastructure
```

If Service B does not expose the required operation, extend Service B's public contract first.

Do not bypass the boundary for convenience.

---

# 24. OpenAPI

Source:

```text
api/openapi/vN/openapi.yaml
```

The OpenAPI contract must document:

- paths
- methods
- parameters
- request bodies
- response bodies
- status codes
- authorization requirements
- normalized errors
- pagination
- examples
- version

Do not document endpoints that do not exist.

Run OpenAPI validation in CI.

Generated code derived from OpenAPI is never edited manually.

---

# 25. Protobuf / gRPC

Source:

```text
api/proto/vN/
```

Generated:

```text
gen/proto/vN/
```

Rules:

- transport contract only
- never use protobuf as domain entity
- never hand-edit generated code
- use Buf lint
- use Buf breaking checks
- never reuse removed field numbers
- reserve removed fields/names
- prefer additive changes
- use clear package/go_package declarations

---

# 26. Controller Rules

Controllers/adapters:

```text
src/controller/http/vN/
src/controller/grpc/vN/
```

Responsibilities:

1. receive request
2. extract transport metadata/context
3. structural validation
4. map transport DTO to command/input
5. call inbound/business service
6. map result to response
7. map semantic errors to transport errors

Controllers do not:

- query DB
- perform business decisions
- orchestrate repositories
- construct infrastructure
- parse provider-specific errors

---

# 27. Mapper Rules

Mappers are pure transformations.

Examples:

```text
request DTO → command
protobuf → command
entity/result → response DTO
entity → event contract
persistence row → entity
```

No:

- DB access
- network access
- repository calls
- service calls
- business authorization
- hidden side effects

Non-trivial mappers require unit tests.

---

# 28. Validator Rules

Transport validators handle structural/input validation:

- required fields
- syntax/format
- UUID/ID shape
- ranges
- enums
- mutually exclusive request fields
- payload structure

Validators do not query repositories.

Existence checks, ownership, lifecycle state, permissions, and other business rules belong to business services/policies.

---

# 29. Middleware Rules

Middleware/interceptors own cross-cutting transport concerns:

- authentication extraction
- request ID
- correlation ID
- tracing
- request logging
- recovery/panic protection
- tenant/org/project context extraction
- generic rate limiting
- request-size limits
- security headers where applicable

Middleware must not become a hidden business layer.

---

# 30. Errors

Use semantic errors inside services.

Examples:

```go
var (
    ErrNotFound            = errors.New("not found")
    ErrConflict            = errors.New("conflict")
    ErrPermissionDenied    = errors.New("permission denied")
    ErrInvalidState        = errors.New("invalid state")
    ErrInvalidConfiguration = errors.New("invalid configuration")
)
```

Wrap internal causes when useful:

```go
return fmt.Errorf("load application: %w", err)
```

Do not expose raw SQL, cache, HTTP, cloud-provider, or internal stack errors publicly.

Transport mapping should be centralized.

Use `errors.Is` / `errors.As` for semantic classification.

Do not compare error strings.

---

# 31. Error Taxonomy

Common normalized categories should include:

```text
INVALID_ARGUMENT
UNAUTHENTICATED
PERMISSION_DENIED
NOT_FOUND
CONFLICT
PRECONDITION_FAILED
QUOTA_EXCEEDED
UNSUPPORTED_CAPABILITY
TIMEOUT
UNAVAILABLE
CANCELLED
INTERNAL
```

Domain-specific codes can extend these without leaking provider internals.

Public errors should include a stable machine code and safe human message.

---

# 32. Context Rules

`context.Context` is the first parameter for request-scoped operations.

Good:

```go
func (s *Service) Get(ctx context.Context, id ApplicationID) (Application, error)
```

Rules:

- never store `context.Context` in long-lived structs
- honor cancellation/deadlines in DB/network/provider calls
- do not use context as a bag for arbitrary dependencies
- context values are for request-scoped metadata only
- define typed private context keys
- propagate trace/correlation context across outbound calls

---

# 33. Dependency Injection

Use explicit constructor injection.

```go
func NewService(
    applications repository.ApplicationRepository,
    publisher ports.EventPublisher,
    clock clock.Clock,
) *Service {
    return &Service{
        applications: applications,
        publisher:    publisher,
        clock:        clock,
    }
}
```

Avoid:

- global DB handles in domain packages
- package-level clients
- hidden singletons
- runtime service locator lookups inside business methods
- reflection-heavy DI by default

Manual wiring is preferred while the application remains manageable.

---

# 34. Capability / Provider Architecture

Infrastructure that companies may replace is modeled as a **capability**, not a vendor.

Examples:

```text
CACHE
CI
OBJECT_STORAGE
MESSAGING
SECRETS
RUNTIME
DEPLOYMENT
OBSERVABILITY
PROVISIONING
SEARCH
MODEL_GATEWAY
```

Business/domain code uses a capability contract.

It must not branch on provider brand.

Bad:

```go
if provider == "redis" {
    // ...
} else if provider == "valkey" {
    // ...
}
```

Good:

```go
type Cache interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}
```

---

# 35. Provider Strategy Pattern

Providers implement stable capability contracts.

```text
Cache
  ├── Valkey adapter
  ├── Redis adapter
  ├── Dragonfly adapter
  └── Memory adapter
```

Registration is configuration/composition driven.

Example:

```go
type CacheFactory func(cfg CacheConfig) (Cache, error)

type CacheRegistry struct {
    factories map[ProviderType]CacheFactory
}
```

The registry is an infrastructure/composition mechanism.

It is **not** a general global service locator for business code.

---

# 36. Avoid `switch` for Provider Selection

Do not grow central files such as:

```go
switch providerType {
case "redis":
    ...
case "valkey":
    ...
case "dragonfly":
    ...
}
```

Use registration maps/factories:

```go
registry.Register(ProviderTypeRedis, newRedis)
registry.Register(ProviderTypeValkey, newValkey)
registry.Register(ProviderTypeDragonfly, newDragonfly)
```

Resolution:

```go
factory, ok := registry.Get(providerType)
if !ok {
    return nil, ErrUnsupportedProvider
}

return factory(cfg)
```

A simple guard `if` is acceptable.

The restriction is against repeated type/vendor branching as an extension mechanism.

---

# 37. Dynamic Provider Resolution

Some platform capabilities are chosen dynamically by organization/project/environment.

In that case, business code may receive an injected logical resolver:

```go
type CacheResolver interface {
    Resolve(ctx context.Context, scope Scope) (Cache, error)
}
```

Callers provide logical context, not vendor names.

```go
cache, err := s.cacheResolver.Resolve(ctx, Scope{
    OrganizationID: orgID,
    ProjectID:      projectID,
    EnvironmentID: environmentID,
})
```

The resolver handles binding precedence internally.

Do not expose `redis`, `valkey`, `jenkins`, etc. through business-service APIs.

---

# 38. Provider Capabilities

Provider contracts may expose capabilities when implementations differ.

Example:

```go
type CapabilitySet interface {
    Supports(capability Capability) bool
}
```

Good:

```go
if provider.Supports(CacheCapabilityStreams) {
    // use optional semantic capability
}
```

Bad:

```go
if provider.Name() == "redis" {
    // vendor-dependent feature assumption
}
```

Check semantics, not brand.

---

# 39. Provider Middleware / Decorators

Cross-cutting provider behavior should be implemented through decorators/middleware when practical.

```text
Business Capability
       ↓
Tracing decorator
       ↓
Metrics decorator
       ↓
Audit decorator
       ↓
Timeout decorator
       ↓
Retry decorator
       ↓
Circuit breaker decorator
       ↓
Provider adapter
```

Do not duplicate retry/logging/metrics/error mapping in every provider.

---

# 40. Provider Error Normalization

Provider adapters map vendor errors into stable platform errors.

Example:

```text
Redis NOAUTH
Vault permission denied
Cloud 403
    ↓
normalized semantic/platform errors
```

Business code does not parse vendor strings.

---

# 41. Shared Technical Libraries

Reusable domain-independent technical code belongs under:

```text
libs/core/<capability>/
```

Examples:

```text
libs/core/database/
libs/core/cache/
libs/core/http/
libs/core/grpc/
libs/core/messaging/
libs/core/telemetry/
libs/core/security/
libs/core/resilience/
libs/core/provider/
libs/core/capability/
```

Good shared code:

- database pool helpers
- transaction primitives
- HTTP server/client primitives
- gRPC primitives
- retry policies
- circuit breaker
- OpenTelemetry setup
- structured logger
- IDs
- clocks
- pagination
- generic secret-reference types
- provider registration primitives

Forbidden shared domain code:

```text
User
Role
Permission
Application
Release
Migration
IAMPolicy
business authorization logic
```

---

# 42. Shared Utility Rule

Before adding reusable helper code:

1. decide if it is domain-independent
2. search existing explicit `libs/core/<capability>` packages
3. place it under the appropriate named capability
4. keep service-specific behavior inside the service

Forbidden generic dumping grounds:

```text
libs/core/common/
libs/core/utils/
libs/core/helpers/
libs/core/misc/
```

Prefer explicit packages.

---

# 43. Service-Local Helpers / Utils

Service-local:

```text
src/helpers/
src/utils/
```

Use only for small pure support functions.

Rules:

- no business workflow orchestration
- no DB calls
- no network calls
- no hidden global state
- meaningful file/package names

If a helper grows into a business responsibility, move it into a properly named service/entity package.

---

# 44. Configuration

Configuration should have explicit ownership and validation.

Typical service config:

```go
type Config struct {
    Enabled bool
    Timeout time.Duration
}
```

Rules:

- defaults are explicit
- validation occurs at startup/config-load boundary
- secrets are references, not embedded values where possible
- environment variables are bootstrap/config inputs, not business APIs
- do not read environment variables throughout business code
- config reload semantics must be explicit
- immutable config is preferred unless live reload is a real requirement

---

# 45. Bootstrap Configuration

Some components cannot depend on themselves to bootstrap.

Minimal trusted bootstrap may include:

```text
environment
listen addresses
database connection reference
initial trust configuration
service identity
bootstrap config source
```

Keep bootstrap minimal.

---

# 46. Constants

Service-specific stable constants belong in:

```text
src/constants/
```

Use typed string constants when behavior depends on a controlled set of values.

```go
type EnvironmentTag string

const (
    EnvironmentTagDevelopment EnvironmentTag = "DEVELOPMENT_ENVIRONMENT"
    EnvironmentTagProduction  EnvironmentTag = "PRODUCTION_ENVIRONMENT"
)
```

Constants are not mutable configuration.

Avoid string literals scattered across business logic.

---

# 47. IDs

Domain IDs are generated by service/domain orchestration, not repositories.

Use a shared injected generator when useful:

```go
type IDGenerator interface {
    New() string
}
```

Rules:

- IDs are opaque unless contract explicitly states otherwise
- do not derive authorization/business state from string parsing
- request bodies cannot override path-owned/server-owned IDs without an explicit contract

---

# 48. Time

Use injected clocks for business time decisions that need deterministic tests.

```go
type Clock interface {
    Now() time.Time
}
```

Avoid calling `time.Now()` throughout business services when time affects business behavior.

Persist timestamps in a consistent timezone/format, normally UTC.

---

# 49. Messaging Structure

Messaging belongs under:

```text
src/messaging/
├── events/vN/
├── producer/
├── consumer/
└── outbox/
```

Event contracts that cross service boundaries are versioned.

---

# 50. Event Contracts

Events are contracts, not business processors.

Recommended envelope:

```go
type EventEnvelope[T any] struct {
    ID            string
    Type          string
    Version       string
    OccurredAt    time.Time
    CorrelationID string
    CausationID   string
    OrganizationID string
    ProjectID     string
    Data          T
}
```

Only include fields required by event consumers and platform observability.

Do not leak internal DB models.

---

# 51. Event Producers

Producers own publishing mechanics:

- serialization
- broker headers/metadata
- publish operation
- retry mechanics
- outbox dispatch
- provider-specific broker mapping

Producer does not decide whether a business event should exist.

Business service makes that decision through an `EventPublisher` port.

---

# 52. Event Consumers

Consumers own:

1. deserialize
2. validate envelope
3. establish trace/correlation context
4. invoke business/inbound service
5. ack/retry/dead-letter according to policy

Consumer contains no business policy.

Consumers must be idempotent when the broker can redeliver.

---

# 53. Outbox Pattern

Use an outbox where DB state and external event publication must be coordinated reliably.

Typical flow:

```text
Business transaction
   ├── domain data mutation
   └── outbox row
          ↓
Outbox dispatcher
          ↓
Message provider
```

Do not attempt naive distributed transactions across databases/brokers.

---

# 54. Idempotency

Operations that may be retried must have explicit idempotency semantics.

For externally retried commands, use idempotency keys where appropriate.

Consumers should protect against duplicate delivery.

Do not assume a broker provides exactly-once business semantics.

---

# 55. Retry

Retry only when:

- operation is idempotent or safely deduplicated
- failure is transient
- total deadline permits it
- retry budget is bounded

Use exponential backoff with jitter for remote dependencies when appropriate.

Never retry:

- validation failures
- authorization failures
- deterministic conflicts
- arbitrary non-idempotent writes without protection

---

# 56. Timeouts

Every remote/provider call should have a bounded timeout/deadline.

Do not rely on infinite defaults.

Timeout ownership should be clear:

```text
request deadline
   ↓
use-case budget
   ↓
provider/client timeout
```

Lower layers should not silently extend the caller's deadline.

---

# 57. Circuit Breakers

Use circuit breakers for unstable remote dependencies when they provide real operational value.

Do not add them to local in-process calls.

Circuit-breaker state must be observable.

---

# 58. Concurrency

Every goroutine must have an owner and a termination path.

Rules:

- prefer `errgroup.WithContext` for related concurrent work
- respect context cancellation
- bound worker pools/concurrency
- avoid fire-and-forget goroutines in request handlers
- do not leak tickers/timers
- close channels from the sending/owning side
- avoid channels when a synchronous call is simpler
- protect shared mutable state explicitly
- run race detector in CI for relevant packages/tests

---

# 59. Background Workers

Workers must expose explicit lifecycle:

```go
type Worker interface {
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
}
```

Worker behavior should include:

- bounded concurrency
- backpressure
- cancellation
- health state
- metrics
- retry/dead-letter semantics where appropriate
- graceful shutdown

---

# 60. Panic Policy

Panic is not normal error handling.

Allowed examples:

- impossible programmer invariant
- generated code assumptions
- startup helper that deliberately offers `Must...` semantics

Request/worker boundaries must recover unexpected panics, log safely, emit telemetry, and return/record failure without leaking sensitive stack details to clients.

---

# 61. Go Style

Mandatory:

- `gofmt`
- `go vet`
- `golangci-lint`
- clear package names
- short, readable functions
- early/guard returns
- avoid unnecessary abstraction
- avoid stutter where Go naming conventions make intent obvious

Prefer readable one-line signatures when reasonable.

Split long signatures only when readability improves.

---

# 62. Naming

Prefer responsibility-revealing names:

```text
AuthorizationService
ApplicationRepository
PolicyStrategy
IAMClient
ConfigProvider
TokenSigner
ApplicationMapper
ExecutionPlanner
```

Avoid vague names:

```text
Manager
Processor
Thing
Common
BaseService
GenericService
Data
UtilService
HelperService
Misc
```

If a package/type cannot be named clearly, its responsibility is probably unclear.

---

# 63. Package Size

Split by responsibility, not file count.

Split when:

- responsibilities are unrelated
- dependencies differ significantly
- tests are difficult to isolate
- separate evolution is obvious

Do not create one package per function.

Do not keep unrelated behavior together merely to reduce folder count.

---

# 64. SOLID in Go

## SRP

Each package/type/function has one primary responsibility.

## OCP

Extend through strategies/adapters/registries/small interfaces rather than expanding central conditionals.

## LSP

Implementations preserve interface semantics.

## ISP

Prefer focused interfaces.

## DIP

Business code depends on abstractions; infrastructure adapters implement them.

Do not apply SOLID mechanically when plain Go code is simpler.

---

# 65. Approved Patterns

Use when they solve a real problem:

- Repository
- Strategy
- Adapter
- Facade
- Constructor Dependency Injection
- Decorator
- Chain of Responsibility
- Options pattern
- Outbox
- Circuit breaker
- Retry
- Factory/Registry for runtime implementations
- State machine for explicit lifecycle-heavy workflows

Patterns should simplify reasoning, not demonstrate cleverness.

---

# 66. Patterns to Avoid by Default

Avoid unless clearly justified:

- god service
- god repository
- global mutable state
- service locator
- generic base service
- generic base repository with domain behavior
- inheritance emulation
- reflection-heavy DI
- giant utility packages
- premature CQRS
- premature event sourcing
- premature distributed transactions
- localhost RPC between local modules
- interface per struct
- generic repository abstraction that hides domain queries
- magic framework layers with unclear ownership

---

# 67. Avoid `switch` as Extensibility Mechanism

A `switch` is acceptable for a closed algorithm where the cases are truly part of one stable domain rule.

It is **not** acceptable as the normal mechanism for adding providers/strategies/handlers.

Bad:

```go
switch kind {
case "redis":
case "valkey":
case "dragonfly":
}
```

Better:

```go
strategies := map[ProviderType]Factory{
    ProviderRedis:     newRedis,
    ProviderValkey:    newValkey,
    ProviderDragonfly: newDragonfly,
}
```

This preserves OCP.

---

# 68. Guard Clauses

Use guard clauses to reduce nesting.

Good:

```go
if req.ID == "" {
    return ErrInvalidArgument
}

if !allowed {
    return ErrPermissionDenied
}
```

Avoid deep pyramids of nested `if` blocks.

---

# 69. Security Rules

Never log:

- passwords
- access tokens
- refresh tokens
- API keys
- client secrets
- private keys
- DB credentials
- secret values
- sensitive request payloads without explicit redaction policy

Use references for secrets where possible.

Validate trust boundaries.

Apply authorization server-side regardless of frontend behavior.

---

# 70. Input Safety

Public endpoints should enforce appropriate limits:

- request body size
- collection length
- string length
- upload size
- pagination maximum
- regex complexity when user-controlled
- decompression limits

Avoid unbounded user-controlled memory allocation.

---

# 71. SQL Safety

Use parameterized queries.

Never concatenate untrusted values into SQL.

Dynamic identifiers/order-by fields must come from explicit allowlists/maps.

Do not expose raw DB error details.

---

# 72. Authorization

Authentication extraction may happen in middleware.

Business authorization belongs in the appropriate business/policy layer.

Avoid hidden authorization in repositories or validators.

Security-sensitive operations should have explicit permission/capability checks near the use case.

---

# 73. Observability

Every service should expose/propagate:

- service name
- service version
- environment
- request ID
- correlation ID
- trace ID
- structured logs
- metrics
- traces
- health
- readiness

Use OpenTelemetry as the default instrumentation standard.

---

# 74. Logging

Use structured logs.

Preferred fields:

```text
service
module
operation
request_id
correlation_id
trace_id
organization_id
project_id
resource_id
error_code
duration_ms
```

Do not log entire objects by default.

Log semantic context, not sensitive payload dumps.

Avoid duplicate logging of the same error at every layer.

Log at the layer that can add meaningful ownership/context.

---

# 75. Metrics

Metrics should be bounded-cardinality.

Good dimensions:

```text
service
operation
status
provider_type
```

Avoid high-cardinality labels such as raw user IDs, request IDs, URLs with IDs, arbitrary error messages.

Measure:

- request rate
- error rate
- latency
- queue depth
- worker throughput
- provider failures
- retry counts
- DB pool utilization
- cache hit/miss when meaningful

---

# 76. Tracing

Create spans around meaningful boundaries:

- HTTP/gRPC requests
- business use cases when useful
- DB calls
- provider calls
- message publish/consume
- execution steps

Propagate trace context across messages and remote calls.

Do not create spans for every trivial helper.

---

# 77. Health and Readiness

Typical states:

```text
STARTING
UP
DEGRADED
DOWN
STOPPING
```

Root endpoints:

```text
/health
/ready
```

Health answers whether the process is alive/operational.

Readiness answers whether it can receive intended traffic/work.

Do not make readiness depend on every optional downstream provider if the service can still function without it.

---

# 78. Graceful Shutdown

Shutdown should generally:

1. stop accepting new requests/work
2. mark not-ready
3. cancel root context
4. stop consumers/schedulers
5. drain/finish bounded in-flight work
6. close listeners
7. flush telemetry/logs
8. close owned resources

Use bounded shutdown deadlines.

---

# 79. Startup Order

Typical startup:

```text
1. bootstrap configuration
2. logging
3. telemetry
4. shared runtime primitives
5. required foundational modules/providers
6. domain services
7. HTTP/gRPC route registration
8. background workers
9. listeners
10. readiness
```

Shutdown reverses dependency order where practical.

---

# 80. Cache Rules

Cache is an optimization unless the domain contract explicitly defines otherwise.

Generic connector/client primitives belong in `libs/core/cache`.

Service-specific cache semantics belong in the owning service.

Examples:

```text
src/infrastructure/cache/valkey/
src/infrastructure/cache/redis/
```

Domain cache-key strategy does not belong in generic core.

Rules:

- explicit TTL
- cache invalidation semantics documented
- correctness should not silently depend on stale cached data
- avoid cache stampedes where relevant
- do not store secrets casually

---

# 81. Object Storage

Object/blob storage adapters should expose logical operations instead of vendor SDKs to business code.

Store metadata/references in the database and large blobs in object storage when appropriate.

Use streaming for large objects.

Do not load unbounded blobs into memory.

---

# 82. External HTTP/gRPC Clients

Raw external calls belong under:

```text
src/infrastructure/clients/
```

or a reusable technical client package when genuinely domain-independent.

Client adapters own:

- endpoint resolution
- authentication
- serialization
- timeout
- retry classification
- response/error translation
- telemetry

Business services call an outbound port/public client interface.

---

# 83. Pagination

Use one consistent pagination contract.

For large/changing collections, cursor pagination is preferred where it provides better correctness/scalability.

Offset pagination may be used for bounded administrative lists.

Never allow unbounded list endpoints.

---

# 84. Sorting / Filtering

Public sort/filter fields must be explicit allowlists.

Do not directly interpolate user-supplied field names into SQL.

Map API field names to persistence expressions.

---

# 85. Generated Code

Never manually edit:

```text
gen/proto/
gen/openapi/
generated mocks
```

Change the source contract/interface and regenerate.

CI should verify generated-code freshness.

---

# 86. Testing Philosophy

Perihellion uses only three top-level test suites:

```text
UNIT TEST
    ↓
SANITY
    ↓
REGRESSION
```

`integration`, `contract`, `conformance`, `functional`, `E2E`, `smoke`, and `critical` are not separate top-level Perihellion test suites.

They describe a validation responsibility or execution style that must be placed inside one of the three canonical suites.

Use this rule:

```text
isolated implementation/business behavior
        → UNIT TEST

fast assembled-system / real-infrastructure confidence
        → SANITY

supported end-user workflow protection
        → REGRESSION
```

Do not duplicate the same assertion across all three levels without a concrete reason.

---

# 87. Backend Unit Testing

All backend Go unit tests use:

```text
Go testing
Ginkgo v2
Gomega
```

Canonical imports:

```go
import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)
```

Backend unit tests are centralized under:

```text
backend/testing/unit-test/
```

They mirror the production path.

Example:

```text
backend/launcher/bootstrap/config/
        ↓
backend/testing/unit-test/launcher/bootstrap/config/

backend/libs/core/logging/
        ↓
backend/testing/unit-test/libs/core/logging/

backend/services/catalog/src/services/application/
        ↓
backend/testing/unit-test/services/catalog/src/services/application/
```

Backend unit tests are **not adjacent** to production Go source files.

Prefer external test packages where practical:

```go
package config_test
```

Do not expose production internals only to make tests easier.

---

# 88. Ginkgo Suite Structure and Unit Coverage

Each coherent backend unit-test package should have a `suite_test.go` when it owns a Ginkgo suite.

Example:

```go
package config_test

import (
    "testing"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

func TestLauncherBootstrapConfigUnit(t *testing.T) {
    RegisterFailHandler(Fail)

    RunSpecs(
        t,
        "Launcher Bootstrap Config Unit Tests",
    )
}
```

Actual behavior tests live in focused files such as:

```text
config_test.go
loader_test.go
validator_test.go
```

`suite_test.go` may contain small suite-wide constants and helpers used by multiple files in that exact suite.

If support grows, split it into focused files:

```text
constants_test.go
helpers_test.go
fixtures_test.go
builders_test.go
```

Do not turn `suite_test.go` or `testing/unit-test/common/` into a generic dumping ground.

Environment-modifying tests must clean up deterministically using Ginkgo lifecycle support such as `BeforeEach` / `DeferCleanup`, or an equivalent deterministic pattern.

New executable business code targets very high meaningful coverage. Coverage is evidence, not a substitute for good scenarios.

Cover important:

- success
- validation
- not found
- conflicts
- authorization
- dependency/repository errors
- ID/time behavior when injected
- transaction decisions
- mapping errors
- provider errors
- cancellation/timeouts where important

Do not claim a coverage percentage unless measured by the actual coverage command.

---

# 89. Real Infrastructure Validation

Real infrastructure correctness belongs in **SANITY** when the semantics cannot be proven by an isolated unit test.

For PostgreSQL, sanity validation should exercise as applicable:

- real PostgreSQL
- real migrations
- real constraints
- real JSONB behavior
- real transactions
- real repository implementations
- startup/readiness with the real database

Do not call SQL-mock-only behavior proof of PostgreSQL correctness.

Use disposable/test-container or isolated test infrastructure where practical.

The top-level suite remains:

```text
SANITY
```

not `integration`.

---

# 90. Public Contract Validation

Versioned public API/client/event contracts must be validated, but `contract` is not a separate top-level suite.

Place contract validation according to scope:

```text
pure encoding/mapping/version rule
        → UNIT TEST

running HTTP/gRPC/public-client compatibility
        → SANITY

user workflow depending on the contract
        → REGRESSION
```

Verify as applicable:

- encoding/decoding
- schema compatibility
- normalized errors
- version behavior
- additive evolution
- in-process and remote client semantic equivalence

---

# 91. Provider Conformance Validation

Every provider implementing a capability must preserve the capability semantics.

Conformance is a **validation responsibility**, not a fourth top-level test suite.

Use:

```text
isolated provider strategy behavior
        → UNIT TEST

real/reference provider capability proof
        → SANITY

product workflow through the capability
        → REGRESSION
```

Example capability expectations:

```text
Cache capability
├── set/get
├── delete
├── TTL
├── not-found semantics
├── cancellation
├── normalized errors
└── health
```

Provider conformance must be reusable without copying the full product regression suite for every provider.

---

# 92. Sanity and Regression Boundaries

Sanity is the fast assembled-system confidence suite.

Sanity may validate:

- application startup
- health/readiness
- migrations
- real PostgreSQL connectivity
- important API behavior
- provider capability health
- authentication/authorization baseline
- a minimal important platform path

Regression validates supported product workflows from the user perspective.

Regression remains product-wide and uses Playwright + TypeScript for browser-driven journeys.

There is no separate `smoke`, `critical`, `functional`, or `E2E` top-level suite.

Risk/criticality may be metadata inside sanity/regression when useful.

---

# 93. Canonical Test Locations

Backend Go unit tests:

```text
backend/testing/unit-test/
├── launcher/
├── services/
├── libs/
├── providers/
└── common/
```

The path below `unit-test/` mirrors production responsibility.

Example:

```text
backend/testing/unit-test/
└── launcher/
    └── bootstrap/
        └── config/
            ├── config_test.go
            └── suite_test.go
```

Repository-wide assembled suites:

```text
testing/
├── sanity/
└── regression/
```

Frontend unit tests remain adjacent to the unit being tested:

```text
Something/
├── Something.tsx
└── __tests__/
    └── Something.test.tsx
```

or:

```text
useApplication/
├── useApplication.ts
└── __tests__/
    └── useApplication.test.ts
```

Use only folders that own real tests/support.

---

# 94. Mocking

Mock real boundaries:

- repositories in business-unit tests
- external/public clients
- clock
- ID generator
- event publisher
- provider capability interfaces

Avoid mocking:

- domain entities
- trivial pure functions
- implementation details
- every internal function call

---

# 95. Test Fixtures / Builders

Use reusable builders for large test entities/DTOs.

```go
application := fixtures.NewApplication(
    fixtures.WithApplicationName("CRM"),
)
```

Keep tests focused on the behavior under test.

Avoid giant repeated struct literals when they obscure intent.

---

# 96. Race Testing

Run `go test -race` for packages/flows where concurrency is meaningful.

Do not ignore detected races.

Concurrency correctness is part of Definition of Done.

---

# 97. Benchmarks

Add benchmarks only for code with a meaningful performance concern.

Do not micro-benchmark trivial functions without a reason.

Use representative payloads and document what regression the benchmark protects.

---

# 98. CI Architecture Enforcement

CI should run, as applicable:

```text
gofmt check
go vet
golangci-lint
go test ./...
go test -race (selected/full cadence)
coverage
Buf lint
Buf breaking
OpenAPI validation
generated-code freshness
architecture/forbidden-import checks
provider conformance tests
migration tests
vulnerability scan
secret scan
container build
```

Architecture should be enforced mechanically where practical.

---

# 99. Forbidden Import Checks

CI should reject patterns such as:

```text
services/A importing services/B/src
services/A importing services/B/migrations
business packages importing concrete provider SDKs
entity packages importing transport packages
services importing launcher
```

A lint/architecture test should encode these rules.

---

# 100. Public Compatibility

Prefer additive changes.

HTTP:

- add optional fields
- add new endpoints
- preserve field meaning
- avoid changing types/semantics in-place

Protobuf:

- add fields
- reserve removed fields
- never reuse field numbers
- preserve method meaning

Client:

- do not remove public methods inside a major version
- do not change semantic meaning inside a major version

Events:

- maintain schema compatibility according to declared version policy

---

# 101. Semantic Versioning

Public artifacts/packages may use semantic versions:

```text
PATCH → compatible bug fix
MINOR → backward-compatible capability
MAJOR → breaking public contract
```

API path major version and package/release semantic version serve different purposes.

Internal refactoring alone does not require a new API major version.

---

# 102. Documentation Rules

Every public service/module should document:

- purpose/ownership
- public APIs
- clients
- config
- dependencies
- persistence ownership
- emitted/consumed events
- health/readiness
- migration notes
- operational considerations

Architecture decisions that change cross-cutting rules belong in ADRs.

---

# 103. ADRs

Use ADRs for decisions such as:

- changing persistence technology
- new cross-service contract pattern
- new provider/capability family
- introducing distributed execution
- changing versioning policy
- adopting a major framework

Do not create ADRs for trivial implementation choices.

---

# 104. AI / Coding-Agent Placement Algorithm

Before creating a file, determine responsibility.

## Public HTTP/gRPC contract?

```text
api/openapi/vN/
api/proto/vN/
dto/vN/
client/vN/
```

## External event contract?

```text
src/messaging/events/vN/
```

## Business logic?

```text
src/services/
```

or intrinsic entity behavior:

```text
src/entity/
```

## Persistence contract?

```text
src/repository/
```

## Persistence implementation?

```text
src/infrastructure/persistence/<technology>/
```

## Provider adapter?

```text
src/infrastructure/<capability>/<provider>/
```

or platform reference provider:

```text
providers/<provider>/
```

## Transport controller?

```text
src/controller/http/vN/
src/controller/grpc/vN/
```

## Mapping?

```text
src/mapper/vN/
```

## Structural validation?

```text
src/validator/vN/
```

## Messaging?

```text
src/messaging/
```

## Service construction?

```text
launcher/
```

## Generic technical reusable capability?

Search:

```text
libs/core/
```

first.

---

# 105. AI / Coding-Agent Lookup Rules

Before writing code, an agent must:

1. inspect the affected service/package
2. search for an existing implementation
3. search `libs/core` for reusable technical primitives
4. inspect target service public client for cross-service calls
5. inspect current public API version
6. inspect current DTO/event patterns
7. inspect repository contracts
8. inspect existing provider/strategy registrations
9. inspect tests for expected semantics
10. avoid duplicate helpers/abstractions
11. preserve existing architecture unless the task explicitly changes it

---

# 106. HTTP API Change Workflow

```text
1. inspect current OpenAPI contract
2. update api/openapi/vN/openapi.yaml
3. update request/response DTOs if needed
4. regenerate generated code if applicable
5. add/update mapper
6. add/update structural validator
7. add/update business service
8. add/update controller
9. update client/vN
10. add/update unit-test validation for pure mapping/business behavior
11. add/update sanity validation for the running public boundary when required
12. add/update regression only when supported user behavior changes
13. update docs
14. run compatibility checks
```

---

# 107. gRPC Change Workflow

```text
1. inspect current proto version
2. update api/proto/vN
3. generate protobuf
4. add/update mapper
5. add/update business service
6. add/update gRPC controller
7. update client/vN
8. run Buf lint
9. run Buf breaking
10. add unit-test and/or sanity contract validation at the appropriate scope
```

---

# 108. Cross-Service Change Workflow

If Service A needs Service B:

```text
1. inspect Service B client/vN
2. use existing supported operation if available
3. otherwise extend Service B public contract
4. inject Service B client into Service A
5. use inprocess implementation while local
6. use grpc/http implementation when remote
7. never import Service B internal packages
8. never query Service B DB
```

---

# 109. Provider Change Workflow

When adding a new infrastructure provider:

```text
1. identify existing logical capability
2. do not create vendor-specific business APIs
3. inspect capability interface
4. extend interface only if semantic capability is genuinely required
5. implement provider adapter
6. register factory/strategy
7. normalize errors
8. add health check
9. add provider telemetry
10. run capability conformance validation in unit-test or sanity as appropriate
11. document configuration and supported capabilities
12. prove existing business modules require no vendor-specific changes
```

---

# 110. Shared Utility Workflow

Before adding shared code:

```text
1. decide whether behavior is truly domain-independent
2. search explicit libs/core capability packages
3. reuse if possible
4. create a clearly named capability package if necessary
5. do not create common/utils/helpers/misc at core root
6. add tests
```

---

# 111. Code Review Checklist

Every backend PR should verify:

- correct package/layer
- SRP preserved
- OCP extension path considered
- real duplication avoided
- business logic isolated
- public contracts versioned
- cross-service calls use public clients
- no internal cross-service import
- no cross-service DB access
- infrastructure hidden behind interfaces where appropriate
- no vendor branching in business logic
- interfaces remain small
- dependencies are constructor-injected
- controllers contain no business rules
- validators do not query repositories
- repositories are persistence-oriented
- mappers are pure
- errors are semantic/normalized
- context/cancellation is propagated
- transaction boundary is correct
- concurrency is bounded/cancellable
- secrets are protected
- logs/metrics/traces are considered
- OpenAPI/protobuf/events are updated
- generated files are fresh and untouched manually
- migrations exist and are tested
- tests cover executable behavior
- compatibility is preserved
- future extraction boundary remains clean

---

# 112. Definition of Done — Backend Unit

A backend change is complete when:

```text
[ ] Responsibility/layer is clear
[ ] SRP is preserved
[ ] OCP extension point is correct where variation exists
[ ] No unnecessary duplication exists
[ ] No forbidden cross-service import exists
[ ] No cross-service DB access exists
[ ] No provider/vendor conditionals leak into business logic
[ ] Public contracts are versioned appropriately
[ ] OpenAPI/protobuf/events match implementation
[ ] Generated code was regenerated, not hand-edited
[ ] Controller has no business logic
[ ] Validator has no repository/business authorization logic
[ ] Mapper is pure
[ ] Repository contains persistence only
[ ] Service owns business orchestration
[ ] Transaction boundary is correct
[ ] Errors are semantic and normalized
[ ] Context/deadlines propagate
[ ] Concurrency has explicit lifecycle
[ ] Security and secret handling are correct
[ ] Structured observability is present
[ ] Migration exists for schema change
[ ] Backend unit tests exist under backend/testing/unit-test/ with mirrored ownership
[ ] Sanity covers real infrastructure/public contract/provider behavior where required
[ ] Regression covers user-critical behavior where appropriate
[ ] gofmt passes
[ ] go vet passes
[ ] golangci-lint passes
[ ] tests pass
[ ] coverage was actually measured
[ ] architecture checks pass
[ ] compatibility checks pass
```

---

# 113. Phase 0 Backend Reference Implementation

Phase 0 should include one small reference module proving the rules.

It should demonstrate:

1. service root `New(...)`
2. launcher-owned construction
3. entity
4. business service
5. repository interface
6. PostgreSQL/mock repository implementation
7. HTTP endpoint + OpenAPI
8. mapper + validator
9. semantic errors
10. public client with in-process implementation
11. one capability/provider interface
12. provider registry/factory
13. mock/reference provider
14. provider conformance validation under the canonical unit-test/sanity model
15. structured logging + OTel
16. unit-test + sanity + regression proof using the canonical three-suite model
17. Docker Compose PostgreSQL
18. CI architecture checks

This reference module becomes the template for future backend contributors and AI agents.

---

# 114. Suggested Engineering Documentation Layout

```text
docs/
└── engineering/
    ├── README.md
    ├── FRONTEND_CODE_GUIDELINES.md
    ├── BACKEND_CODE_GUIDELINES.md
    ├── API_CONTRACT_GUIDELINES.md
    ├── PROVIDER_SDK_GUIDELINES.md
    ├── DATABASE_GUIDELINES.md
    ├── TESTING_GUIDELINES.md
    ├── DEPLOYMENT_GUIDELINES.md
    ├── DOCKER_GUIDELINES.md
    ├── KUBERNETES_GUIDELINES.md
    ├── SHELL_SCRIPT_GUIDELINES.md
    ├── SECURITY_GUIDELINES.md
    ├── OBSERVABILITY_GUIDELINES.md
    ├── GIT_AND_PR_GUIDELINES.md
    ├── RELEASE_GUIDELINES.md
    └── DOCUMENTATION_GUIDELINES.md
```

The root `CONTRIBUTING.md` should link to these standards.

---

# 115. Final Mandatory Backend Rules

1. **One monorepo and modular monolith first.**
2. **Strong service boundaries even in one process.**
3. **Cross-service dependencies use versioned public clients.**
4. **Never import another service's internals.**
5. **Never query another service's database.**
6. **Business logic belongs in services/entities, not transport/persistence layers.**
7. **Repositories are persistence-only.**
8. **Controllers are transport-only.**
9. **Mappers are pure.**
10. **Validators are structural.**
11. **Infrastructure is behind focused interfaces where variation/testability requires it.**
12. **Provider selection is strategy/registry driven, not vendor branching.**
13. **Business code asks for capabilities, not brand names.**
14. **Provider registry is not a business-layer service locator.**
15. **Shared technical code belongs in explicit `libs/core/<capability>` packages.**
16. **Domain logic never belongs in `libs/core`.**
17. **No giant common/utils/misc packages.**
18. **Constructor dependency injection is the default.**
19. **Interfaces are small and intentional.**
20. **Public APIs/events/clients are versioned; internals are not duplicated by API version.**
21. **Generated code is never edited manually.**
22. **Errors are semantic and normalized.**
23. **Context cancellation/timeouts propagate end-to-end.**
24. **Retries are bounded and only used when safe.**
25. **Concurrency is bounded and cancellable.**
26. **Transactions are owned by the business use case.**
27. **Schema ownership is explicit.**
28. **Every schema change has a migration.**
29. **Observability is built in from the beginning.**
30. **Secrets never appear in logs or normal responses.**
31. **Pure/business logic is meaningfully unit-tested.**
32. **Persistence semantics that require a real database are validated in sanity against real PostgreSQL.**
33. **Public/provider contract and conformance responsibilities are validated inside unit-test or sanity, never as extra top-level suites.**
34. **Architecture rules are enforced in CI.**
35. **If code does not clearly fit a layer, re-evaluate responsibility before creating a folder.**

---

# 116. Guiding Principle

The target backend shape is:

```text
Stable Public Contract
        ↓
Transport Adapter
        ↓
Business Use Case
        ↓
Domain Model
        ↓
Small Port / Repository Interface
        ↑
Infrastructure / Provider Adapter
```

For replaceable infrastructure:

```text
Logical Capability
       ↓
Injected Resolver / Interface
       ↓
Provider Binding
       ↓
Registered Strategy
       ↓
Concrete Provider Adapter
```

The result should be:

> **A clean Go backend that is cheap to run as one process, easy to clone and deploy, strongly modular, provider-neutral, testable, secure, observable, and capable of becoming distributed later without rewriting its business logic.**

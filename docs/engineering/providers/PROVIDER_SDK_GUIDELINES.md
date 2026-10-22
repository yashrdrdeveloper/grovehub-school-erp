# Provider SDK Guidelines

> **Status:** Mandatory Phase 0 engineering standard  
> **Location:** `docs/engineering/providers/PROVIDER_SDK_GUIDELINES.md`  
> **Applies to:** provider contracts, capability interfaces, provider metadata, registry, resolver, bindings, capability detection, configuration, lifecycle, health, middleware/decorators, conformance tests, reference providers, external provider packages, and marketplace installation  
> **Primary implementation:** Go Provider SDK first; additional Java/Python/TypeScript SDKs may follow when real provider authors need them  
> **Core principles:** capability-first design, strategy pattern, explicit bindings, provider neutrality, small interfaces, conformance, safe lifecycle, no service locator, and replaceable infrastructure

---

## 1. Objective

Perihellion must integrate with infrastructure without coupling product/business modules to vendors.

The rule is:

> **Platform capability = what Perihellion requires. Provider implementation = what the organization already owns or chooses.**

---

## 2. Target Architecture

```text
Product Module / Service
        ↓
Platform Capability SDK
        ↓
Capability Resolver
        ↓
Provider Binding
        ↓
Provider Registry
        ↓
Registered Strategy
        ↓
Concrete Provider
```

Business code asks for a capability.

It never branches on provider brand.

---

## 3. Mandatory Rules

1. Business code depends on logical capabilities, not vendor SDKs.
2. Provider implementations live behind capability contracts.
3. Provider registry is infrastructure/composition, not a business service locator.
4. Resolve provider from explicit context/binding.
5. Do not probe all providers on every operation.
6. Provider type registration and provider instance binding are separate concepts.
7. Provider instance configuration is validated before use.
8. Provider secrets are references/securely resolved values, not casually persisted/logged.
9. Interfaces are small and capability-specific.
10. A provider may implement multiple capability interfaces.
11. Capability support is explicit.
12. Vendor-specific optional behavior uses capability detection, not type assertions scattered through product code.
13. Binding resolution follows one deterministic hierarchy.
14. Provider errors map to normalized platform errors.
15. Provider calls propagate context/deadline/cancellation.
16. Retry is centralized and capability-safe.
17. Health and readiness are explicit.
18. Providers have lifecycle ownership.
19. Provider implementations pass conformance tests.
20. Product modules do not import provider implementations.
21. Provider implementations do not import product business modules.
22. Provider SDK/contracts are versioned.
23. Default/reference providers may remain in the main monorepo initially.
24. Optional ecosystem providers may be extracted later only when independent release/install/ownership is useful.
25. Out-of-process providers must use a stable RPC contract.
26. Marketplace installation never bypasses trust/signature/policy checks.
27. Provider configuration changes are auditable.
28. Provider selection is environment/project/org/global binding-driven.
29. Missing capability is explicit, not silently emulated.
30. No vendor `switch` in business code.

---

## 4. Capability Families

Initial families include:

```text
SOURCE
BUILD
TEST
ARTIFACT
RUNTIME
RESOURCE
PROVISIONING
ROUTING
SECRET
OBSERVABILITY
COST
POLICY
EXECUTION
ORCHESTRATION
MODEL
CACHE
EVENTS
OBJECT_STORAGE
```

Families may grow through governed contract evolution.

---

## 5. Capability vs Provider

Example:

```text
Capability:
CI_RUN

Provider implementations:
Woodpecker
Jenkins
GitHub Actions
GitLab CI
Buildkite
```

Product code asks:

```go
ci, err := resolver.CI(ctx, scope)
```

not:

```go
if provider == "jenkins" { ... }
```

---

## 6. Provider Metadata

Every provider type exposes metadata.

Concept:

```go
type ProviderMetadata struct {
    Type        string
    Name        string
    Version     string
    Capabilities []Capability
}
```

Actual contract should use canonical SDK types.

---

## 7. Provider Base Contract

A minimal provider contract may include:

```go
type Provider interface {
    Metadata() ProviderMetadata
    Health(ctx context.Context) error
    Close(ctx context.Context) error
}
```

Only include lifecycle methods that all registered providers can meaningfully support.

---

## 8. Small Capability Interfaces

Example cache capability:

```go
type Cache interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}
```

Do not create one giant `InfrastructureProvider` interface.

---

## 9. Capability IDs

Capability IDs are stable machine constants.

Example:

```text
CACHE_KV
CACHE_STREAMS
CI_RUN
CI_LOGS
OBJECT_PUT
OBJECT_GET
SECRET_RESOLVE
OBSERVABILITY_LOGS
```

Do not make provider names capability IDs.

---

## 10. Capability Detection

Optional behavior:

```go
if provider.Supports(CapabilityCacheStreams) {
    ...
}
```

Better still, resolve the narrower capability interface directly when that is the operation needed.

Do not use:

```go
if _, ok := provider.(*RedisProvider); ok { ... }
```

in business code.

---

## 11. Provider Type Registration

Registry stores provider factories/types.

Concept:

```text
provider type: valkey
factory: creates configured instances
declared capabilities
configuration schema
SDK compatibility
```

Registration does not select the provider for a tenant/project/environment.

---

## 12. Provider Instance

A provider instance represents configured infrastructure.

Example:

```text
id
provider type
display name
scope
config reference
secret references
enabled state
health
metadata
```

Instance IDs are opaque.

---

## 13. Provider Binding

Binding maps:

```text
scope + capability
        ↓
provider instance
```

Example:

```text
environment production
capability CI_RUN
provider instance company-jenkins-prod
```

---

## 14. Binding Hierarchy

Default precedence:

```text
Environment override
        ↓
Project override
        ↓
Organization / Business Unit override
        ↓
Global default
```

The nearest explicit applicable binding wins.

Do not add hidden fallback orders.

---

## 15. Resolve Once per Context

Within one request/execution/context, resolve the binding once and reuse it safely.

Avoid repeated DB/config lookups for every low-level call.

Cache resolution with correct invalidation/versioning.

---

## 16. Scope

A provider resolution scope may include:

```text
organization
project
application
environment
execution
region/site when applicable
```

The exact capability contract defines required scope.

---

## 17. No Global Tenant-Specific Singleton

A global registry may hold provider **types/factories**.

It must not expose one mutable tenant-specific provider singleton for all requests.

Instance lifecycle/caching must be keyed by explicit binding/config identity.

---

## 18. Dependency Injection

Business services receive a focused resolver/capability dependency.

Example:

```go
type Service struct {
    cache CacheResolver
}
```

Avoid global `platform.Get("cache")` service-locator APIs inside business code.

---

## 19. SDK Package Shape

Recommended:

```text
sdk/provider/
├── capability/
├── metadata/
├── registry/
├── resolver/
├── binding/
├── config/
├── errors/
├── lifecycle/
├── conformance/
└── version/
```

The concrete repository layout should remain consistent with backend standards.

---

## 20. Provider Package Shape

Example:

```text
providers/valkey/
├── provider.go
├── config.go
├── factory.go
├── health.go
├── cache.go
├── types/
├── tests/
└── README.md
```

A provider can split further when it implements several capabilities.

---

## 21. Configuration Schema

Every provider type declares a typed/configurable schema.

Separate:

```text
non-secret configuration
secret references
runtime-derived values
```

Do not accept arbitrary unvalidated `map[string]any` deep into provider code.

---

## 22. JSON/JSONB Metadata

Flexible provider metadata may use JSON/JSONB at persistence boundaries.

SDK/core logic should map it into typed provider configuration before operation.

---

## 23. Secrets

Provider config stores:

```text
secret reference
credential profile reference
workload identity binding
```

where possible.

Raw secret values should be short-lived in memory and redacted from logs/errors.

---

## 24. Provider Factory

A factory creates/initializes a provider instance from validated config and dependencies.

Factory responsibilities:

```text
validate config
resolve secret references through approved secret capability
construct client
perform safe initialization
return provider instance
```

Business modules never call provider constructors directly.

---

## 25. Health

Provider health should distinguish where useful:

```text
configured
reachable
authenticated
ready for capability
degraded
```

Avoid expensive destructive health probes.

---

## 26. Readiness

A provider may be configured but not ready.

Resolution should surface meaningful `Unavailable`/configuration errors rather than nil pointer behavior.

---

## 27. Lifecycle

Long-lived provider clients/pools have explicit ownership.

Define:

```text
create
reuse
refresh/rotate
close
evict
```

Avoid unbounded connection/client caches.

---

## 28. Configuration Revision

Provider instance configuration should have a revision/version.

Executions can record which provider binding/config revision was used for evidence/reproducibility.

---

## 29. Binding Changes

Binding changes are audited.

For long-running executions, decide whether execution remains pinned to resolved provider revision or re-resolves at safe boundaries.

Default durable execution behavior should prefer reproducibility/pinning unless the use case explicitly requests dynamic resolution.

---

## 30. Normalized Errors

Provider SDK uses normalized semantic errors:

```text
ErrUnauthorized
ErrForbidden
ErrUnavailable
ErrTimeout
ErrNotFound
ErrConflict
ErrQuotaExceeded
ErrInvalidConfiguration
ErrUnsupportedCapability
```

Provider-specific error details stay internal/diagnostic.

---

## 31. Error Metadata

Diagnostic provider errors may preserve safe structured fields:

```text
provider instance ID
operation
retryable
upstream status category
```

Never include secret values.

---

## 32. Timeout

Provider calls respect `context.Context`.

Capability-level defaults may exist, but caller deadlines win when earlier.

---

## 33. Retry

Retry belongs in a shared decorator/resilience layer.

Decision uses:

```text
operation idempotency
normalized error
attempt count
deadline
provider guidance
```

Do not make every provider implement ad hoc retry differently.

---

## 34. Circuit Breaker

Use only for remote dependencies where it solves a measured/reasonable failure mode.

Circuit state should be scoped by provider instance/capability, not globally across unrelated customers.

---

## 35. Rate Limit

Providers may expose normalized quota/rate-limit errors/metadata.

Perihellion should not depend on vendor-specific header names outside the provider adapter.

---

## 36. Middleware / Decorators

Provider call path may include:

```text
Capability
   ↓
Tracing
   ↓
Metrics
   ↓
Audit where required
   ↓
Timeout
   ↓
Retry
   ↓
Circuit Breaker where justified
   ↓
Provider
```

Keep ordering documented and consistent.

---

## 37. Telemetry

Every provider operation should expose:

```text
provider type
provider instance ID (bounded cardinality policy)
capability
operation
duration
success/failure category
retry count where useful
```

Do not label metrics with raw resource/user IDs.

---

## 38. Audit

Administrative actions:

```text
create provider
update config
change secret reference
change binding
disable/delete provider
```

are auditable.

Ordinary high-volume reads may not require audit events unless security policy requires it.

---

## 39. Conformance Harness

A capability defines reusable tests.

Example:

```go
func RunCacheConformance(
    t *testing.T,
    factory ProviderFactory,
)
```

Each provider runs the same semantic checks.

---

## 40. Conformance Scope

Test:

```text
happy path
not found
timeout/cancellation where practical
invalid config
health
idempotent semantics
capability-specific edge cases
```

---

## 41. Provider-Specific Tests

Conformance proves common semantics.

Provider-specific tests cover:

```text
vendor configuration
mapping
special optional capabilities
SDK edge behavior
```

---

## 42. Mock Provider

Phase 0 includes a mock/memory provider to prove extension without external infrastructure.

It is deterministic and implements the same capability contract.

Do not add test-only methods to the production capability interface.

---

## 43. Reference Providers

Reference/open defaults may include:

```text
Valkey
Woodpecker CI
NATS JetStream
MinIO/S3-compatible
OpenBao
Prometheus
Loki
Tempo
OpenTofu
Flux
```

Only enable/install the capabilities a deployment needs.

---

## 44. Provider Semantics Are Not Automatically Interchangeable

Two technologies may both be called "cache" but differ in:

```text
TTL
transactions
streams
consistency
atomic operations
limits
```

A provider is compatible only if it satisfies the semantic capability contract.

---

## 45. Capability Granularity

Do not make capabilities so broad that every provider must implement unrelated features.

Prefer:

```text
CacheKV
CacheStreams
```

over one giant cache interface when semantics differ.

---

## 46. No Leaky Vendor Types

Capability methods use SDK/platform types.

Do not return:

```text
*redis.Client
*jenkins.Build
*kubernetes.Pod
```

from generic provider interfaces.

Map vendor data to normalized contracts.

---

## 47. Escape Hatches

Provider-specific advanced operations may exist only through an explicitly provider-specific extension API/package.

Core product modules must not depend on it for generic behavior.

---

## 48. External / Out-of-Process Providers

When providers become external packages/processes:

```text
Core
 ↓
Provider RPC contract
 ↓
Provider process/plugin service
```

Use stable gRPC/Protobuf where appropriate.

Do not load arbitrary untrusted code into the core process.

---

## 49. Provider Protocol Handshake

Out-of-process providers should negotiate:

```text
SDK/protocol version
provider metadata
capabilities
health
configuration schema/version
```

Reject incompatible providers clearly.

---

## 50. Provider Packaging

Initially:

```text
main repo/providers/
```

Later optional providers may be distributed as:

```text
OCI artifact/package
signed marketplace package
separate repo when independently maintained
```

Source repository count and install experience are separate concerns.

---

## 51. Marketplace Trust

Provider package installation should eventually verify:

```text
publisher identity
signature
checksum/digest
SDK compatibility
declared permissions/capabilities
configuration schema
```

Air-gapped installation must support offline trust metadata.

---

## 52. Provider Upgrade

Upgrade must define:

```text
config compatibility
protocol/SDK compatibility
rollback
health verification
running execution behavior
```

Do not replace a provider under active durable work without defined semantics.

---

## 53. Provider Removal

Prevent deletion when bindings/running executions still depend on the provider unless a controlled force workflow exists.

Return dependency information.

---

## 54. Binding Validation

Before activating binding, verify:

```text
instance exists/enabled
capability supported
scope allowed
configuration healthy enough
policy permits binding
```

---

## 55. Provider Policy

Governance may restrict:

```text
which provider types are allowed
which scopes can bind them
regions/endpoints
secret sources
capabilities
```

Provider SDK exposes metadata needed for policy but does not implement business governance itself.

---

## 56. Provider Configuration UI

Frontend renders provider configuration from normalized typed schema/metadata where possible.

UI depends on capability/provider metadata contract, not backend implementation folders.

---

## 57. Provider API

Administrative APIs expose:

```text
provider types
instances
capabilities
health
bindings
config schema
```

Business APIs should resolve bindings automatically from context.

---

## 58. Provider Context in Evidence

Record enough to reproduce/understand an operation:

```text
capability
provider type
provider instance
binding/config revision
```

Redact secrets.

---

## 59. Provider SDK Versioning

Use semantic/major compatibility rules.

Breaking capability interface changes require a new SDK/protocol major or supported migration layer.

---

## 60. Backward Compatibility

Core should support a documented provider SDK compatibility window.

Do not silently break installed ecosystem providers.

---

## 61. Provider PR Checklist

```text
[ ] Capability contract reused/extended deliberately
[ ] No vendor branch in business module
[ ] Config typed/validated
[ ] Secrets are references/redacted
[ ] Factory/lifecycle correct
[ ] Context/deadlines propagate
[ ] Errors normalized
[ ] Health implemented
[ ] Telemetry present
[ ] Conformance suite passes
[ ] Provider-specific tests pass
[ ] Binding/capability metadata correct
[ ] No import of product business internals
[ ] Docs/runbook updated
```

---

## 62. Definition of Done

```text
[ ] Provider satisfies declared capabilities
[ ] Conformance tests pass
[ ] Invalid config fails clearly
[ ] Secret handling reviewed
[ ] Binding resolution works at correct scope
[ ] Health/readiness observable
[ ] Lifecycle bounded/closed
[ ] Errors normalized
[ ] No business-code vendor dependency introduced
[ ] Compatibility documented
```

---

## 63. Phase 0 Deliverables

Phase 0 should include:

```text
provider SDK package
capability IDs
metadata contract
registry
resolver
binding model
binding precedence
normalized provider errors
config validation
mock/memory provider
one conformance suite
one reference capability/provider proof
telemetry decorator proof
```

Do not implement every real provider in Phase 0.

---

## 64. Agent Rules

Before provider work inspect:

```text
docs/AGENT.md
this file
capability interface
registry/resolver
binding model
existing providers
conformance harness
security/observability guidelines
```

If an existing capability fits, extend/register; do not create a parallel vendor-specific path.

---

## 65. Final Mandatory Rules

1. Business code asks for capability, not vendor.
2. Resolver + binding selects provider.
3. Registry holds provider types/factories, not business dependencies.
4. Bindings are deterministic and scoped.
5. Provider interfaces are small.
6. Capabilities express semantics, not brand.
7. Provider configs are typed and secure.
8. Secrets are references/redacted.
9. Errors are normalized.
10. Context/deadlines propagate.
11. Retry/telemetry are standardized decorators.
12. Providers pass conformance tests.
13. Default providers can live in monorepo.
14. Optional providers may later ship independently.
15. No service locator in business code.
16. Provider selection is auditable and reproducible.

---

## 66. Guiding Principle

```text
Business Intent
      ↓
Logical Capability
      ↓
Contextual Binding
      ↓
Registered Strategy
      ↓
Replaceable Provider
```

> **Infrastructure may change. Product semantics must not.**

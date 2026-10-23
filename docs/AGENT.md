# AGENT.md

> **Status:** Mandatory repository context and execution contract for AI coding agents  
> **Location:** `docs/AGENT.md`  
> **Applies to:** AI coding agents, code generators, refactoring agents, automated reviewers, and any tool that changes or evaluates Perihellion source  
> **Purpose:** Define exactly how repository context is acquired, which engineering standards apply, which boundaries must be inspected, and how work is validated before completion

---

## 0. GroveHub Current Mandatory Overrides

These GroveHub-specific rules are newer than older generic examples in this document and the long-form engineering standards. **When there is a conflict, these rules and the dedicated GroveHub standards below win.**

Mandatory current standards:

```text
docs/engineering/backend/GROVEHUB_BACKEND_STRUCTURE.md
docs/engineering/testing/BACKEND_UNIT_TEST_STRUCTURE.md
docs/engineering/observability/TRACE_CONTEXT_GUIDELINES.md
docs/engineering/observability/GROVEHUB_OBSERVABILITY.md
docs/engineering/scripts/GROVEHUB_LOCAL_DEVELOPMENT.md
```

### Backend structure

The root `backend/launcher/` is the platform edge/composition boundary.

Backend services do **not** have service-local `launcher/` folders. Current service structure uses the service root plus explicit `controller/`, `routes/`, `constants/`, `dto/`, `api/`, and other responsibility folders only when needed.

```text
backend/services/<service>/
├── <service>.go
├── controller/
├── routes/
├── constants/
├── dto/
├── api/
└── tests/
```

Do not recreate `backend/services/<service>/launcher/` from an older blueprint.

### Backend unit tests

The old centralized path `backend/testing/unit-test/` is superseded and must not be recreated.

Unit tests live with their owner:

```text
backend/launcher/tests/unit-test/...
backend/services/<service>/tests/unit-test/...
backend/libs/core/tests/unit-test/...
backend/providers/<provider>/tests/unit-test/...
```

Inside `tests/unit-test/`, mirror the production responsibility.

Every backend unit-test package folder must contain `suite_test.go` and use Ginkgo v2 + Gomega:

```go
import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)
```

`suite_test.go` owns suite registration plus shared constants, fixtures, builders, fakes/stubs, setup/cleanup, and reusable helper functions for that test package. Other `_test.go` files should stay focused on behavior specs and use the suite support.

Repository-wide `testing/sanity/` and `testing/regression/` remain repository-wide and are not moved into one backend owner.

### Trace context and API-to-API propagation

OpenTelemetry trace context is mandatory across the GroveHub backend ecosystem.

The launcher/edge starts or continues the inbound trace. Every nested service/business/dependency operation must receive the existing `context.Context` and create child spans rather than a disconnected trace.

GroveHub HTTP responses expose the active trace identifier in:

```text
X-Trace-ID
```

Outbound HTTP calls also forward `X-Trace-ID` for support/debug correlation. W3C `traceparent` remains the authoritative OpenTelemetry propagation header because parent/child relationships require both trace and span context. `X-Trace-ID` is never added to the JSON response envelope.

Request/execution logs inside a trace include:

```text
traceId
spanId
traceFlags
```

On failures, spans/logs capture safe request/response metadata such as method, route/path, status, request size, response size, duration, service/component, traceId, and spanId. Generic middleware must not dump raw request/response bodies, Authorization headers, cookies, passwords, tokens, credentials, or arbitrary sensitive query data.

Before adding an HTTP client, service call, worker, messaging boundary, provider call, or middleware, verify that trace context crosses that boundary.

### GroveHub observability identity, metrics, and Grafana

The canonical logical telemetry services are defined centrally. Current identities are:

```text
service.namespace = grovehub

service.name = LAUNCHER (EDGE)
grovehub.service.key = launcher-edge

service.name = IDENTITY
grovehub.service.key = identity
```

Do not emit ad-hoc variants such as `edge`, `launcher`, `platform`, `identity-service`, or differently cased service names. Future logical services must be registered centrally before service-owned telemetry is emitted.

Metrics, traces, and logs must agree on logical service identity. The same trace ID continues across `LAUNCHER (EDGE) -> IDENTITY -> real dependency` boundaries while each child span has its own span ID.

HTTP latency uses OpenTelemetry histogram metrics so Grafana computes P50/P95/P99. Do not calculate request percentiles inside application code. Metric dimensions must be bounded; use route templates and never use trace IDs, request IDs, user/student IDs, emails, or raw resource URLs as metric labels.

The initial GroveHub observability surface includes:

```text
RPS / request count
5xx error rate
active requests
P50 / P95 / P99 HTTP latency
service-operation latency/errors
service availability
Go runtime metrics
WARN/ERROR logs
trace waterfall/correlation
Grafana overview dashboard + filters
bootstrap alert rules
```

When canonical telemetry names, resources, labels, or dashboard query assumptions change, update the Grafana/dashboard contract and `GROVEHUB_OBSERVABILITY.md` in the same task.

### Local developer hot reload

GroveHub development commands have fixed semantics:

```text
make backend-start
    -> one-time backend binary
    -> no hot reload

make backend-dev
    -> Air in foreground
    -> hot reload

make dev-start
    -> Air-backed backend in background + Swagger
    -> hot reload

make dev-start-with-logs
    -> Air-backed backend in background + Swagger + local OTel collector
    -> hot reload
```

Do not route `dev-start` or `dev-start-with-logs` through a one-time backend binary.

Air watches backend `.go`, `go.mod`, and `go.sum` changes. After a developer runs `git pull`, watched file changes must trigger dependency refresh, rebuild, and backend restart automatically. Air does not perform `git pull` itself.

---

## 1. Purpose

This is the **first repository document an AI agent reads**.

It is not a replacement for detailed standards.

It tells an agent:

```text
what to read
what not to read
how to identify ownership
how to find the public boundary
how to search before creating
how to understand consumers
how to choose tests
how to validate changes
how to report actual results
```

An agent must not implement from the task prompt alone when repository context is available.

---

## 2. Canonical Engineering Standards

```text
docs/
├── AGENT.md
│
└── engineering/
    ├── README.md
    ├── backend/
    │   ├── BACKEND_CODE_GUIDELINES.md
    │   └── GROVEHUB_BACKEND_STRUCTURE.md
    ├── frontend/
    │   └── FRONTEND_CODE_GUIDELINES.md
    ├── testing/
    │   ├── TESTING_GUIDELINES.md
    │   ├── BACKEND_UNIT_TEST_STRUCTURE.md
    │   └── REGRESSION_TESTING.md
    ├── scripts/
    │   ├── SCRIPT_AND_MAKEFILE_GUIDELINES.md
    │   └── GROVEHUB_LOCAL_DEVELOPMENT.md
    ├── api/
    │   └── API_CONTRACT_GUIDELINES.md
    ├── database/
    │   └── DATABASE_GUIDELINES.md
    ├── providers/
    │   └── PROVIDER_SDK_GUIDELINES.md
    ├── security/
    │   └── SECURITY_GUIDELINES.md
    ├── observability/
    │   ├── OBSERVABILITY_GUIDELINES.md
    │   ├── TRACE_CONTEXT_GUIDELINES.md
    │   └── GROVEHUB_OBSERVABILITY.md
    ├── ci/
    │   └── CI_CD_GUIDELINES.md
    ├── deployment/
    │   ├── DEPLOYMENT_GUIDELINES.md
    │   ├── DOCKER_GUIDELINES.md
    │   └── KUBERNETES_GUIDELINES.md
    ├── release/
    │   └── RELEASE_GUIDELINES.md
    ├── git/
    │   └── GIT_AND_PR_GUIDELINES.md
    └── documentation/
        └── DOCUMENTATION_GUIDELINES.md
```

`docs/engineering/README.md` is the human-readable standards index.

---

## 3. Do Not Read Everything Blindly

Use progressive context.

```text
Task
 ↓
Owning area
 ↓
Primary engineering standard
 ↓
Owning package/service
 ↓
Public contract
 ↓
Target implementation
 ↓
Tests/stories/mocks/runbook
 ↓
Consumers/dependencies only when relevant
```

The objective is the **smallest complete context required for correctness**, not the smallest number of files.

---

## 4. First Classification

Before opening implementation files, classify the task:

```text
Backend
Frontend
API / contract
Database / migration
Provider / capability
Testing
Playwright regression
Scripts / Makefile
Security
Observability
CI
Docker
Kubernetes
Deployment
Release
Git / repository governance
Documentation
Cross-cutting / full stack
Unknown
```

If more than one applies, load each relevant standard.

---

## 5. Primary Context Matrix

### Backend

Always:

```text
docs/AGENT.md
docs/engineering/backend/BACKEND_CODE_GUIDELINES.md
docs/engineering/backend/GROVEHUB_BACKEND_STRUCTURE.md
docs/engineering/testing/TESTING_GUIDELINES.md
docs/engineering/testing/BACKEND_UNIT_TEST_STRUCTURE.md
```

Add:

```text
API change       -> API_CONTRACT_GUIDELINES
DB change        -> DATABASE_GUIDELINES
Provider change  -> PROVIDER_SDK_GUIDELINES
Security change  -> SECURITY_GUIDELINES
Telemetry change -> OBSERVABILITY_GUIDELINES + TRACE_CONTEXT_GUIDELINES + GROVEHUB_OBSERVABILITY
```

---

### Frontend

Always:

```text
docs/AGENT.md
docs/engineering/frontend/FRONTEND_CODE_GUIDELINES.md
docs/engineering/testing/TESTING_GUIDELINES.md
```

Add:

```text
API contract     -> API_CONTRACT_GUIDELINES
critical journey -> REGRESSION_TESTING
security/auth    -> SECURITY_GUIDELINES
telemetry        -> OBSERVABILITY_GUIDELINES
```

---

### API / Contract

Read:

```text
docs/engineering/api/API_CONTRACT_GUIDELINES.md
docs/engineering/backend/BACKEND_CODE_GUIDELINES.md
docs/engineering/backend/GROVEHUB_BACKEND_STRUCTURE.md
docs/engineering/testing/TESTING_GUIDELINES.md
docs/engineering/testing/BACKEND_UNIT_TEST_STRUCTURE.md
```

Inspect:

```text
OpenAPI/proto source
DTOs
client/vN
controllers/mappers
contract tests
all known consumers
```

If frontend consumes it, inspect relevant frontend service/reader.

---

### Database

Read:

```text
docs/engineering/database/DATABASE_GUIDELINES.md
docs/engineering/backend/BACKEND_CODE_GUIDELINES.md
docs/engineering/backend/GROVEHUB_BACKEND_STRUCTURE.md
docs/engineering/testing/TESTING_GUIDELINES.md
docs/engineering/testing/BACKEND_UNIT_TEST_STRUCTURE.md
```

Inspect:

```text
owning service
entity
repository
PostgreSQL adapter
migrations
queries/indexes
real DB tests
```

If shipping migration, also read release/deployment guidelines.

---

### Provider

Read:

```text
docs/engineering/providers/PROVIDER_SDK_GUIDELINES.md
docs/engineering/backend/BACKEND_CODE_GUIDELINES.md
docs/engineering/backend/GROVEHUB_BACKEND_STRUCTURE.md
docs/engineering/testing/TESTING_GUIDELINES.md
docs/engineering/testing/BACKEND_UNIT_TEST_STRUCTURE.md
docs/engineering/security/SECURITY_GUIDELINES.md
docs/engineering/observability/OBSERVABILITY_GUIDELINES.md
docs/engineering/observability/TRACE_CONTEXT_GUIDELINES.md
docs/engineering/observability/GROVEHUB_OBSERVABILITY.md
```

Inspect:

```text
capability interface
registry
resolver
binding model
existing strategies
config schema
conformance harness
```

---

### General Testing

Read:

```text
docs/engineering/testing/TESTING_GUIDELINES.md
docs/engineering/testing/BACKEND_UNIT_TEST_STRUCTURE.md when backend Go is involved
```

Then the implementation-area guideline.

---

### Playwright Regression

Read:

```text
docs/engineering/testing/REGRESSION_TESTING.md
docs/engineering/testing/TESTING_GUIDELINES.md
```

Also load frontend/API/backend standards when the test touches those boundaries.

---

### Scripts / Makefile

Read:

```text
docs/engineering/scripts/SCRIPT_AND_MAKEFILE_GUIDELINES.md
docs/engineering/scripts/GROVEHUB_LOCAL_DEVELOPMENT.md when local start/stop/dev behavior is involved
docs/engineering/observability/GROVEHUB_OBSERVABILITY.md when local telemetry/export behavior is involved
```

Inspect:

```text
Make target
affected script
scripts/lib
tool version source
scripts/docker if container runtime involved
CI usage
```

Add Docker/CI/testing/backend/frontend guideline according to what the script orchestrates.

---

### Security

Read:

```text
docs/engineering/security/SECURITY_GUIDELINES.md
```

Then the owning implementation guideline.

Security-sensitive triggers include:

```text
authentication
authorization
tenant isolation
secret handling
token/session
file upload
webhook
remote command/terminal
provider credentials
marketplace package execution
AI tool permissions
cryptography
```

---

### Observability

Read:

```text
docs/engineering/observability/OBSERVABILITY_GUIDELINES.md
docs/engineering/observability/TRACE_CONTEXT_GUIDELINES.md
docs/engineering/observability/GROVEHUB_OBSERVABILITY.md
```

Inspect shared telemetry/logger primitives, service registry, collector configuration, dashboards, and alert query assumptions before creating new signal conventions.

---

### CI

Read:

```text
docs/engineering/ci/CI_CD_GUIDELINES.md
docs/engineering/scripts/SCRIPT_AND_MAKEFILE_GUIDELINES.md
docs/engineering/git/GIT_AND_PR_GUIDELINES.md
```

Add testing/security/release depending on job purpose.

---

### Docker

Read:

```text
docs/engineering/deployment/DOCKER_GUIDELINES.md
docs/engineering/deployment/DEPLOYMENT_GUIDELINES.md
docs/engineering/security/SECURITY_GUIDELINES.md
```

For development Compose, also read script guideline.

---

### Kubernetes

Read:

```text
docs/engineering/deployment/KUBERNETES_GUIDELINES.md
docs/engineering/deployment/DEPLOYMENT_GUIDELINES.md
docs/engineering/deployment/DOCKER_GUIDELINES.md
docs/engineering/security/SECURITY_GUIDELINES.md
```

Add release/observability as relevant.

---

### Deployment

Read:

```text
docs/engineering/deployment/DEPLOYMENT_GUIDELINES.md
```

Then Docker/Kubernetes and state/security/observability guidelines as applicable.

---

### Release

Read:

```text
docs/engineering/release/RELEASE_GUIDELINES.md
docs/engineering/git/GIT_AND_PR_GUIDELINES.md
docs/engineering/ci/CI_CD_GUIDELINES.md
docs/engineering/deployment/DEPLOYMENT_GUIDELINES.md
docs/engineering/testing/TESTING_GUIDELINES.md
docs/engineering/testing/REGRESSION_TESTING.md
docs/engineering/security/SECURITY_GUIDELINES.md
```

---

### Git / PR Governance

Read:

```text
docs/engineering/git/GIT_AND_PR_GUIDELINES.md
```

Inspect `.github/`, CI config, contribution docs, and release rules if tags/branches are affected.

---

### Documentation

Read:

```text
docs/engineering/documentation/DOCUMENTATION_GUIDELINES.md
```

Then inspect the actual code/contract that the documentation describes.

Never silently "correct" a source-backed contract from generic knowledge.

---

## 6. Context Loading Levels

### Level 0 — Repository Instructions

```text
docs/AGENT.md
```

### Level 1 — Applicable Standard

One or more files under:

```text
docs/engineering/
```

### Level 2 — Owning Unit

Examples:

```text
backend/services/<service>/
backend/launcher/
backend/libs/core/
frontend/packages/domain/<domain>/
frontend/packages/shell/
providers/<provider>/
tests/regression/domains/<domain>/
deployments/
scripts/
```

Read nearest `README.md` / `RUNBOOK.md`.

### Level 3 — Public Boundary

Frontend:

```text
index.ts
```

Backend:

```text
api/openapi/vN
api/proto/vN
dto/vN
client/vN
```

Provider:

```text
capability interface
metadata/config schema
```

Deployment:

```text
values schema
release manifest
public Make target
```

### Level 4 — Target Implementation

Read target file and adjacent:

```text
types
constants
helpers
services
interfaces
mappers
repositories
tests
stories
mocks
```

### Level 5 — Consumers / Dependencies

Expand only when:

```text
target consumes it
it consumes target
it defines target contract
change crosses the boundary
compatibility is affected
```

---

## 7. Stop Condition for Context Gathering

Implementation can begin only when the agent can answer:

```text
Who owns this behavior?
Which guideline governs it?
Where does the code belong?
What is its public boundary?
What existing unit should be reused?
Who consumes it?
What must remain compatible?
Which tests prove current behavior?
Which tests must prove the change?
What validation commands must run?
```

If a materially important answer is unknown, gather more context.

---

## 8. Search Before Create

Before creating:

```text
helper
util
component
hook
HOC
service
reader
repository
port
mapper
validator
provider
strategy
fixture
builder
page object
script helper
CI job
deployment template
```

search by:

```text
responsibility
business noun
verb/use case
type/interface
endpoint
route
error
constant
similar behavior
```

Do not search only for the filename you already intend to invent.

---

## 9. Existing Code Is Context, Not Automatic Approval

Legacy code can violate current standards.

Do not assume:

```text
existing = correct architecture
```

When legacy code conflicts with a canonical guideline:

```text
do not duplicate the legacy pattern
do not refactor unrelated code automatically
preserve compatibility
isolate unavoidable interaction
implement new code according to current standard where possible
```

---

## 10. Scope Discipline

Implement the requested outcome.

Do not opportunistically:

```text
rewrite unrelated modules
rename unrelated packages
extract microservices
introduce a framework
migrate all legacy code
add future providers
restructure the repository
```

unless the task explicitly requires it.

---

## 11. Frontend Context Rules

Canonical tiers:

```text
packages/base
packages/business
packages/components
packages/domain
packages/shell
frontend/app
```

Dependency intent:

```text
app
 ↓
shell
 ├────────→ components → base
 └────────→ domain
               ├──────→ business → base
               ├──────→ components → base
               └──────→ base
```

Forbidden:

```text
domain -> shell
business -> shell
components -> shell
cross-domain deep imports
```

Shell imports domains through root public APIs only.

---

## 12. Frontend Public Boundary

Cross-domain:

```ts
import { something } from '@domain/catalog';
```

Never:

```ts
import { something } from '@domain/catalog/internal/path';
```

If a needed unit is private, decide whether the owning domain should deliberately expose it.

Do not bypass the boundary.

---

## 13. Frontend Shared Placement

Before promotion:

```text
pure generic technical       -> base
shared business concept      -> business
shared visual/platform unit  -> components
domain behavior              -> domain
global composition           -> shell
```

Do not move code upward merely because it might be reusable someday.

---

## 14. Frontend Services

API calls live at:

```text
packages/domain/<domain>/services
packages/business/services
```

never inside component folders or shell.

Inspect existing response readers and authoritative request/response types before adding new ones.

---

## 15. Frontend Shell

Shell owns:

```text
global providers
Redux store composition
route aggregation
navigation aggregation
domain registrations
global layout
startup skeleton
global fallbacks
```

Shell owns **no product business logic**.

---

## 16. Frontend Mocking

Runtime public API:

```text
@domain/<domain>
```

Test/development mock API:

```text
@domain/<domain>/mock
```

Production runtime code never imports mock packages.

---

## 17. Frontend Stories / Runbooks

Before changing a public shared visual component, inspect:

```text
stories
tests
public props/types
consumers
```

Update Storybook when visual/public states change.

Inspect/update package/domain/shell `RUNBOOK.md` when public usage changes.

---

## 18. Backend Context Rules

Dependency direction:

```text
Controller
   ↓
Business Service / Inbound Port
   ↓
Domain Entity
   ↓
Repository / Outbound Port
   ↑
Infrastructure Adapter
```

Do not bypass layers because code is in one process.

Use `GROVEHUB_BACKEND_STRUCTURE.md` for current service folder placement.

---

## 19. Backend Cross-Service

```text
Service A
   ↓
Service B client/vN
   ↓
in-process OR remote transport
   ↓
Service B
```

Never:

```text
A imports B/src
A queries B database
```

If the public client lacks an operation, extend B's public contract first.

Trace context must continue through the public client using `context.Context`; remote HTTP transport must use the core traced transport so `traceparent` and `X-Trace-ID` are propagated.

---

## 20. Backend Placement

```text
OpenAPI              -> api/openapi/vN
Protobuf             -> api/proto/vN
DTO                   -> dto/vN
Public client         -> client/vN
Business use case     -> src/services
Domain entity         -> src/entity
Persistence contract  -> src/repository
Outbound port         -> src/ports
Persistence adapter   -> src/infrastructure/persistence
Controller            -> controller/ or service-owned controller path
Routes                -> routes/
Mapper                -> src/mapper
Structural validator  -> src/validator
Messaging             -> src/messaging
Construction          -> service root New(...)
Unit tests            -> owner/tests/unit-test/<mirrored path>
```

The root platform composition remains in `backend/launcher/`. Service-local `launcher/` folders are not used in current GroveHub services.

Do not invent a random folder when responsibility is unclear.

---

## 21. API Context Rule

Before a public contract change inspect:

```text
machine-readable contract
DTO
client/vN
mapper
controller
contract tests
known consumers
```

Classify:

```text
additive compatible
deprecated
breaking
```

Do not create a new major version for internal refactor.

---

## 22. Database Context Rule

Before persistence changes inspect:

```text
owning service/schema
entity
repository interface
adapter
migration history
constraints/indexes
transaction boundary
real PostgreSQL tests
```

Never infer schema solely from frontend shape.

---

## 23. Provider Context Rule

Product code asks for capability.

Provider selection:

```text
scope
 ↓
resolver
 ↓
binding precedence
 ↓
registry/factory
 ↓
provider instance
```

Provider Registry is infrastructure, not a global business service locator.

---

## 24. Security Context Rule

For security-sensitive behavior inspect both:

```text
trust source
authorization boundary
```

Do not make frontend checks authoritative.

Do not solve test failures by weakening authorization or secret handling.

---

## 25. Observability Context Rule

Before adding logs/metrics/traces inspect shared telemetry abstractions and naming plus both `TRACE_CONTEXT_GUIDELINES.md` and `GROVEHUB_OBSERVABILITY.md`.

Do not invent:

```text
new request-id spelling
new trace-id spelling
new metric unit
new provider telemetry format
new service.name variant
new high-cardinality metric label
```

without extending the shared standard.

Every request-scoped backend boundary must preserve `context.Context` and the active trace.

Use the canonical telemetry service registry. Current service resources are `LAUNCHER (EDGE)` and `IDENTITY`; machine keys are `launcher-edge` and `identity`.

For latency percentiles, emit histogram measurements and let Grafana/Prometheus calculate P50/P95/P99. Do not store precomputed percentile values as per-request metrics.

---

## 26. Testing Context Rule

Ask:

```text
What is the lowest effective layer?
Who owns the production behavior?
Does real infrastructure semantics matter?
Is there already a fixture/builder?
Is this a public-contract test?
Does this need browser regression?
```

Backend unit-test routing is owner-local as defined in `BACKEND_UNIT_TEST_STRUCTURE.md`; do not use `backend/testing/unit-test`.

Do not duplicate every scenario in E2E.

---

## 27. Playwright Context Rule

Inspect:

```text
failing/target spec
domain page/flow objects
fixtures
builders
API setup helpers
cleanup
trace/failure artifact
actual UI state
```

Never fix flakiness with arbitrary `waitForTimeout`.

---

## 28. Scripts / Make Context Rule

Trace:

```text
Make target
  ↓
small script
  ↓
scripts/lib helper
  ↓
toolchain / Compose / application command
```

Makefile remains thin.

For GroveHub local development startup, also follow `docs/engineering/scripts/GROVEHUB_LOCAL_DEVELOPMENT.md`.

No duplicated bootstrap logic.

No broad `pkill` / `killall`.

No silent `sudo`.

---

## 29. Docker Context Rule

Inspect:

```text
Dockerfile
.dockerignore
runtime user
build source/tool versions
Compose profile if relevant
health
security scan
release image usage
```

Do not add secrets or environment-specific config to image.

---

## 30. Kubernetes Context Rule

Inspect:

```text
Helm values/schema
Deployment/Job
RBAC/ServiceAccount
probes
resources
secrets
migration flow
network exposure
upgrade impact
```

Do not add cluster-admin because a provider needs one new action.

---

## 31. CI Context Rule

CI must use existing Make/scripts when practical.

Before adding a job inspect:

```text
local command
existing stage
required checks
secret scope
fork behavior
cache
artifacts
release impact
```

No second CI-only architecture.

---

## 32. Release Context Rule

Never change publishing/version behavior without inspecting:

```text
Git/tag policy
CI release flow
artifact manifest
deployment/upgrade tests
API/DB/provider compatibility
security scans
```

Never claim support without executed evidence.

---

## 33. Documentation Context Rule

When asked to document current behavior, derive it from:

```text
code
machine-readable contract
canonical guideline
actual command/config
```

Do not silently fill missing source facts from general knowledge.

Mark future/proposed behavior clearly.

---

## 34. Public Boundary Change Rule

Compatibility-sensitive boundaries include:

```text
frontend root index.ts exports
shared component props
domain registrations
backend OpenAPI/protobuf
DTO/client/vN
events
provider capability interfaces
provider SDK/protocol
database schema used by rolling versions
CLI commands
Make targets
configuration keys
Helm values
release artifact names
```

Before modification:

```text
search consumers
classify compatibility
prefer additive change
add migration/deprecation when required
update tests/docs
```

---

## 35. Generated Code

Never hand-edit generated output.

Workflow:

```text
edit source contract/schema/interface
run pinned generator
review generated diff
run freshness/compatibility checks
```

---

## 36. Bug-Fix Workflow

```text
1. read AGENT + applicable standard
2. understand observed failure
3. locate owning unit
4. inspect existing tests
5. trace execution path
6. identify root cause
7. add protection at lowest effective layer
8. implement minimal fix
9. run targeted tests
10. run relevant broader gates
11. report actual results
```

---

## 37. Feature Workflow

```text
1. identify owner
2. inspect public contract
3. inspect neighboring implementation
4. identify reuse
5. inspect data/security/telemetry impact
6. implement in dependency direction
7. add tests/stories/regression as needed
8. update docs/runbooks/contracts
9. validate
```

---

## 38. Refactor Workflow

```text
identify behavior that must not change
inspect public boundaries/consumers
preserve compatibility
run existing tests before/after where practical
do not mix new behavior unless requested
```

---

## 39. Shared-Abstraction Promotion

Before moving code into:

```text
frontend/base
frontend/business
frontend/components
backend/libs/core
provider SDK
shared script lib
```

answer:

```text
Is it truly shared?
Is the contract stable?
Does promotion reduce conceptual duplication?
Does it introduce upward/circular dependency?
Can existing public path remain compatible?
```

---

## 40. New Dependency Rule

Before adding a third-party dependency:

```text
search existing stack
justify capability
check maintenance/security/license
check bundle/image impact
prefer standard library/existing dependency when sufficient
```

Do not add a package for a trivial helper already supported by existing tools.

---

## 41. New Architecture Rule

An agent must not introduce:

```text
new framework
new database
new queue
new state manager
new provider mechanism
new test framework
new CI platform
new deployment system
```

without explicit task/approved architecture decision.

---

## 42. Context Budget

Do not load the whole monorepo for a one-file fix.

Expand when:

```text
public compatibility
dependency direction
consumer behavior
security/data boundary
test behavior
```

requires it.

---

## 43. Validation — Backend

As applicable:

```text
gofmt
go test
go vet
golangci-lint
race tests
real PostgreSQL integration
OpenAPI validation
Buf lint/breaking
generated freshness
architecture/forbidden imports
```

---

## 44. Validation — Frontend

As applicable:

```text
format
ESLint
TypeScript
Jest
React Testing Library
Storybook build/test
package-boundary rules
circular dependency checks
mock-production import checks
```

---

## 45. Validation — Regression

As applicable:

```text
target Playwright test
domain suite
sanity/critical-risk coverage as defined by current testing standards
test independence
trace artifacts
no arbitrary sleeps
no test.only
```

---

## 46. Validation — Scripts

As applicable:

```text
make help
make doctor
affected Make target
ShellCheck
shfmt when adopted
docker compose config
start/readiness/stop
idempotent rerun
non-interactive CI behavior
```

For hot-reload changes also validate that a watched backend file change triggers an Air rebuild and restart without rerunning Make.

For observability changes, validate collector configuration, signal export assumptions, dashboard query syntax, and no-data behavior before claiming Grafana verification.

---

## 47. Validation — API / DB / Provider

As applicable:

```text
OpenAPI/Buf compatibility
client contract tests
migration clean+upgrade test
real PostgreSQL integration
provider conformance
provider health/error behavior
```

---

## 48. Validation — Deployment / Release

As applicable:

```text
Docker build/scan/smoke
Helm lint/template/schema
Kubernetes schema/security validation
fresh install
supported upgrade
rollback/forward recovery
air-gap verification
checksums/SBOM/signature
```

---

## 49. Reporting Rule

At completion report only what actually happened:

```text
files changed
behavior changed
tests/checks actually run
actual results
checks not run and why
compatibility/migration impact
known limitations
```

Never claim:

```text
tests passed
coverage achieved
browser verified
migration verified
air-gap verified
Grafana data verified
```

without executing the relevant validation.

---

## 50. Forbidden Agent Behaviors

An agent must not:

- implement from the prompt without inspecting existing code
- invent a parallel architecture
- bypass frontend public domain boundaries
- import frontend domain internals from shell/another domain
- import another backend service's `src`
- query another service's database
- recreate service-local `launcher/` folders in GroveHub services
- recreate `backend/testing/unit-test/`
- add a backend unit-test package without `suite_test.go`
- add backend Go unit tests outside the Ginkgo v2 + Gomega standard
- drop trace context across a request-scoped service/API/dependency boundary
- use `X-Trace-ID` as a substitute for W3C `traceparent` propagation
- put trace IDs inside API JSON response envelopes
- invent an ad-hoc GroveHub telemetry `service.name` instead of using the canonical registry
- calculate P95/P99 inside individual application requests instead of recording histograms
- use trace IDs, request IDs, user/student IDs, emails, raw URLs, raw errors, or SQL as metric labels
- log raw secrets, authorization headers, cookies, credentials, or arbitrary request/response bodies from generic tracing middleware
- route `make dev-start` or `make dev-start-with-logs` through a one-time backend binary instead of Air
- add frontend API services inside components or shell
- create generic dumping-ground `common` / `misc`
- hand-edit generated code
- put business logic in scripts/CI/deployment templates
- add arbitrary Playwright sleeps
- hide flakiness with retry
- seed E2E by direct application DB mutation
- leak provider vendor types into generic capability contracts
- use provider registry as business service locator
- weaken security to make implementation easier
- log secrets
- add unbounded metric labels
- silently run privileged system commands
- use broad process killing
- publish mutable release artifacts
- claim validation that was not run
- refactor unrelated code without scope
- build future architecture without current requirement

---

## 51. Agent Pre-Edit Checklist

```text
[ ] Task owner identified
[ ] Applicable guideline(s) read
[ ] GroveHub current override docs checked when backend/testing/telemetry/local-dev is involved
[ ] Target implementation inspected
[ ] Nearest public boundary inspected
[ ] Existing equivalent searched
[ ] Authoritative types/contracts identified
[ ] Tests/stories/mocks/runbook inspected
[ ] Consumers inspected if public behavior changes
[ ] Trace-context boundary checked for request-scoped backend changes
[ ] Canonical telemetry service identity checked when observability is involved
[ ] Metric cardinality/route-template safety checked when metrics are involved
[ ] Dashboard/alert query impact checked when telemetry names or labels change
[ ] Local dev command semantics preserved when start/stop scripts change
[ ] Security/data/provider impact understood
[ ] Minimal requested scope understood
[ ] Validation commands identified
```

---

## 52. Definition of Done

```text
[ ] Requested behavior implemented
[ ] Correct layer/package owns it
[ ] Public boundaries respected
[ ] No duplicate abstraction introduced
[ ] Compatibility considered
[ ] Required tests added/updated in the correct owner-local location
[ ] Every backend unit-test package has suite_test.go and uses Ginkgo/Gomega
[ ] Trace context preserved across relevant backend boundaries
[ ] Canonical GroveHub telemetry service identity preserved
[ ] Metrics use bounded labels and percentile-compatible histograms where relevant
[ ] Grafana/dashboard/alert artifacts updated when telemetry contracts changed
[ ] Dev-start commands remain Air-backed when local development startup is changed
[ ] Stories/mocks/runbooks/docs updated where relevant
[ ] Security/observability included where relevant
[ ] Relevant commands actually run
[ ] Actual results reported
[ ] No unrelated architecture introduced
```

---

## 53. Final Context Rule

> **Repository context is a dependency graph, not a pile of files.**

Use:

```text
Task
  ↓
Owner
  ↓
Standard
  ↓
Public Contract
  ↓
Implementation
  ↓
Consumers / Dependencies
  ↓
Tests / Stories / Mocks / Runbooks
  ↓
Minimal Correct Change
```

---

## 54. Final Principle

> **Understand ownership before placement. Understand contracts before modification. Search before creating. Preserve boundaries, canonical telemetry identity, and trace context before optimizing. Validate before reporting completion.**

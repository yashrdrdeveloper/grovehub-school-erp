# Perihellion Engineering Standards

> **Status:** Canonical engineering documentation index  
> **Location:** `docs/engineering/README.md`  
> **Audience:** Engineers, reviewers, AI coding agents, maintainers, provider authors, release engineers, and operators  
> **Purpose:** Route every engineering task to the correct detailed standard without duplicating those standards here

---

## 1. Principle

Perihellion engineering standards are intentionally split by responsibility.

Do not read every file for every task.

Use:

```text
Task
  ↓
Owning area
  ↓
Primary guideline
  ↓
Cross-cutting guidelines only when boundary is affected
```

`docs/AGENT.md` defines the full context-loading workflow for AI agents.

For GroveHub backend work, the newer GroveHub-specific structure/testing/tracing standards listed below take precedence over conflicting generic examples in older standards.

---

## 2. Canonical Structure

```text
docs/
├── AGENT.md
│
└── engineering/
    ├── README.md
    │
    ├── backend/
    │   ├── BACKEND_CODE_GUIDELINES.md
    │   └── GROVEHUB_BACKEND_STRUCTURE.md
    │
    ├── frontend/
    │   └── FRONTEND_CODE_GUIDELINES.md
    │
    ├── testing/
    │   ├── TESTING_GUIDELINES.md
    │   ├── BACKEND_UNIT_TEST_STRUCTURE.md
    │   └── REGRESSION_TESTING.md
    │
    ├── scripts/
    │   └── SCRIPT_AND_MAKEFILE_GUIDELINES.md
    │
    ├── api/
    │   └── API_CONTRACT_GUIDELINES.md
    │
    ├── database/
    │   └── DATABASE_GUIDELINES.md
    │
    ├── providers/
    │   └── PROVIDER_SDK_GUIDELINES.md
    │
    ├── security/
    │   └── SECURITY_GUIDELINES.md
    │
    ├── observability/
    │   ├── OBSERVABILITY_GUIDELINES.md
    │   └── TRACE_CONTEXT_GUIDELINES.md
    │
    ├── ci/
    │   └── CI_CD_GUIDELINES.md
    │
    ├── deployment/
    │   ├── DEPLOYMENT_GUIDELINES.md
    │   ├── DOCKER_GUIDELINES.md
    │   └── KUBERNETES_GUIDELINES.md
    │
    ├── release/
    │   └── RELEASE_GUIDELINES.md
    │
    ├── git/
    │   └── GIT_AND_PR_GUIDELINES.md
    │
    └── documentation/
        └── DOCUMENTATION_GUIDELINES.md
```

---

## 3. What Each File Owns

| Guideline | Owns |
| --- | --- |
| `BACKEND_CODE_GUIDELINES.md` | Generic Go service architecture, layering, clients, ports, repositories, errors, messaging, provider consumption |
| `GROVEHUB_BACKEND_STRUCTURE.md` | **Current GroveHub backend placement**: root launcher/edge, service root + `controller/` + `routes/`, no service-local launcher, owner-local tests |
| `FRONTEND_CODE_GUIDELINES.md` | React/TypeScript package architecture, `base/business/components/domain/shell`, public boundaries, services/readers, hooks, Storybook, mocks |
| `TESTING_GUIDELINES.md` | Broad test strategy across unit/sanity/regression and supporting validation responsibilities |
| `BACKEND_UNIT_TEST_STRUCTURE.md` | **Current GroveHub Go unit-test placement**: owner-local `tests/unit-test`, mirrored paths, mandatory `suite_test.go`, Ginkgo v2 + Gomega |
| `REGRESSION_TESTING.md` | Playwright browser regression/E2E, page objects, fixtures, isolation, traces, release regression |
| `SCRIPT_AND_MAKEFILE_GUIDELINES.md` | `Makefile`, bootstrap, toolchain downloads, Docker/Colima detection, scripts, local start/stop/test commands |
| `API_CONTRACT_GUIDELINES.md` | REST/OpenAPI, gRPC/Protobuf, DTOs, public clients, errors, pagination, idempotency, events, compatibility |
| `DATABASE_GUIDELINES.md` | PostgreSQL ownership, schema, migrations, constraints, transactions, indexes, JSONB, backup/restore |
| `PROVIDER_SDK_GUIDELINES.md` | Capabilities, registry, resolver, bindings, provider lifecycle, configuration, conformance, marketplace providers |
| `SECURITY_GUIDELINES.md` | Authn/authz, secrets, input/output safety, tenant isolation, supply chain, audit, secure development |
| `OBSERVABILITY_GUIDELINES.md` | General OTel, logs, metrics, traces, correlation, health/readiness, dashboards, redaction |
| `TRACE_CONTEXT_GUIDELINES.md` | **Current GroveHub distributed tracing contract**: `context.Context`, `X-Trace-ID`, W3C `traceparent`, edge/service/outbound spans, log correlation, safe failure diagnostics |
| `CI_CD_GUIDELINES.md` | Woodpecker/reference CI, pipeline gates, caching, artifacts, secrets, promotion, release evidence |
| `DEPLOYMENT_GUIDELINES.md` | Deployment profiles, config/state, migrations, upgrades, backup/restore, Compose/Kubernetes selection |
| `DOCKER_GUIDELINES.md` | Dockerfiles, images, Compose behavior, non-root runtime, scans, multi-arch |
| `KUBERNETES_GUIDELINES.md` | Helm, RBAC, securityContext, resources, probes, migrations, HA, networking, GitOps |
| `RELEASE_GUIDELINES.md` | SemVer, tags, artifacts, checksums, SBOM, signatures, upgrade gates, air-gap bundle |
| `GIT_AND_PR_GUIDELINES.md` | Branches, commits, PRs, CODEOWNERS, reviews, merge policy, open-source contribution |
| `DOCUMENTATION_GUIDELINES.md` | README/runbook/ADR/API documentation, ownership, source-of-truth and documentation style |

---

## 4. Task-to-Guideline Routing

### GroveHub backend feature

```text
GROVEHUB_BACKEND_STRUCTURE
+ BACKEND_CODE_GUIDELINES for generic layering/domain rules
+ BACKEND_UNIT_TEST_STRUCTURE
+ TESTING_GUIDELINES
+ API_CONTRACT_GUIDELINES if API/event/client changes
+ DATABASE_GUIDELINES if persistence changes
+ SECURITY_GUIDELINES if authorization/secret/trust boundary changes
+ OBSERVABILITY_GUIDELINES + TRACE_CONTEXT_GUIDELINES for request/dependency boundaries
```

### Frontend feature

```text
FRONTEND_CODE_GUIDELINES
+ API_CONTRACT_GUIDELINES if backend contract consumption changes
+ SECURITY_GUIDELINES for auth/permission/sensitive rendering
+ TESTING_GUIDELINES
+ REGRESSION_TESTING for critical user journey
```

### Provider

```text
PROVIDER_SDK_GUIDELINES
+ GROVEHUB_BACKEND_STRUCTURE
+ BACKEND_CODE_GUIDELINES
+ SECURITY_GUIDELINES
+ OBSERVABILITY_GUIDELINES
+ TRACE_CONTEXT_GUIDELINES when provider calls are request/execution scoped
+ BACKEND_UNIT_TEST_STRUCTURE
+ TESTING_GUIDELINES
```

### Database

```text
DATABASE_GUIDELINES
+ GROVEHUB_BACKEND_STRUCTURE
+ BACKEND_CODE_GUIDELINES
+ BACKEND_UNIT_TEST_STRUCTURE
+ TESTING_GUIDELINES
+ RELEASE_GUIDELINES when migration ships
```

### Browser regression

```text
REGRESSION_TESTING
+ TESTING_GUIDELINES
+ FRONTEND_CODE_GUIDELINES when UI helpers/domain behavior are touched
+ API_CONTRACT_GUIDELINES when API setup/assertion semantics are touched
```

### Local tooling

```text
SCRIPT_AND_MAKEFILE_GUIDELINES
+ DOCKER_GUIDELINES when container/Compose changes
+ CI_CD_GUIDELINES when same command is used by CI
```

### CI

```text
CI_CD_GUIDELINES
+ SCRIPT_AND_MAKEFILE_GUIDELINES
+ GIT_AND_PR_GUIDELINES
+ TESTING_GUIDELINES
+ SECURITY_GUIDELINES
+ RELEASE_GUIDELINES if publishing
```

### Deployment

```text
DEPLOYMENT_GUIDELINES
+ DOCKER_GUIDELINES or KUBERNETES_GUIDELINES
+ DATABASE_GUIDELINES for migrations/state
+ SECURITY_GUIDELINES
+ OBSERVABILITY_GUIDELINES
+ TRACE_CONTEXT_GUIDELINES when telemetry transport changes
+ RELEASE_GUIDELINES
```

### Release

```text
RELEASE_GUIDELINES
+ GIT_AND_PR_GUIDELINES
+ CI_CD_GUIDELINES
+ DEPLOYMENT_GUIDELINES
+ TESTING_GUIDELINES
+ REGRESSION_TESTING
+ SECURITY_GUIDELINES
```

---

## 5. Style Shared by All Standards

Engineering standards intentionally use the same pattern:

```text
Status / Applies To / Core Principles
        ↓
Objective
        ↓
Architecture / Ownership
        ↓
Mandatory Rules
        ↓
Detailed Rules + Good/Bad Examples
        ↓
Dependency / Compatibility Rules
        ↓
Testing / Enforcement
        ↓
PR Checklist
        ↓
Definition of Done
        ↓
Phase 0 Deliverables
        ↓
Agent Rules
        ↓
Final Mandatory Rules
        ↓
Guiding Principle
```

This consistency is intentional.

Humans and AI agents should be able to predict where a rule is located.

---

## 6. General Engineering Principles

Across all standards:

```text
SRP
OCP
DRY
explicit ownership
stable public contracts
composition over coupling
configuration/strategy/registry over giant conditionals
capability over provider brand
small focused interfaces
strong typing
deterministic tests
automation over reviewer memory
```

---

## 7. Public Boundary Principle

The repository repeatedly uses:

```text
Public Contract
      ↓
Private Implementation
```

Examples:

```text
Frontend domain       -> root index.ts
Backend service       -> client/vN + API/DTO contract
Provider              -> capability interface
Database              -> repository/port
Deployment            -> versioned artifact/config
```

Consumers must not bypass these boundaries because internal files are visible in the monorepo.

---

## 8. Provider Principle

The platform rule is:

```text
Business Code
     ↓
Capability
     ↓
Resolver
     ↓
Binding
     ↓
Provider Strategy
```

Vendor branching does not belong in business modules.

---

## 9. Testing Principle

```text
lowest effective layer
        +
owner-local backend unit tests
        +
real infrastructure when semantics matter
        +
small high-value regression suite
```

Do not use browser regression as the only correctness layer.

Backend unit-test ownership follows `BACKEND_UNIT_TEST_STRUCTURE.md`.

---

## 10. Trace Context Principle

```text
edge/server span
      ↓ same traceId
service/business spans
      ↓ same traceId
outbound dependency/client span
```

`X-Trace-ID` is the support/debug header. W3C `traceparent` is the authoritative distributed parent/child propagation contract.

Do not drop `context.Context` across request-scoped backend boundaries.

---

## 11. Automation Principle

```text
make <intent>
     ↓
small script
     ↓
shared helpers
     ↓
pinned tools / runtime
```

Users remember commands, not setup procedures.

---

## 12. Documentation Principle

One authoritative source owns each rule.

This index routes; it does not duplicate detailed standards.

For current GroveHub backend placement/testing/tracing, the GroveHub-specific standards are authoritative over conflicting historical generic examples.

---

## 13. Updating Standards

A standard change should be reviewed as architecture.

When changing one standard, inspect cross-references in:

```text
docs/AGENT.md
docs/engineering/README.md
related standards
```

Do not update one dependency rule while leaving the context router stale.

---

## 14. Final Rule

> **If ownership is unclear, do not invent a new folder or abstraction. Determine the responsibility and public boundary first.**

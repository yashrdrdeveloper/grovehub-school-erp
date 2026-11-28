# Shared Library Ownership and Capability Design

## 1. Library rule

Shared libraries are **capability-first**. A library exists only when a reusable contract has a distinct owner and proven stable semantics. Duplication within one service is not a reason to publish a global helper.

```text
backend/libs/
├── core/
│   ├── <capability>/
│   │   ├── <capability>.go       # stable public facade
│   │   └── internal/
│   │       ├── src/             # contracts, constants, helpers, impl, validators, mappers
│   │       └── resources/tests/ # private unit tests, mocks, fixtures and helpers
│   ├── annotation/
│   ├── database/
│   └── docs/engineering/
└── <technical-library>/
    ├── <capability>/
    └── docs/engineering/
```

For core capabilities, use the standard internal/src categories operations, helpers, constants, contracts and config when needed. Do not create empty categories. Expose a single root facade. Create only folders with real responsibilities. Services import the public capability facade, **never** a library's private helpers, contracts, mappers or validators unless explicitly part of its API. Libraries must follow Go `internal` import restrictions.

## 2. Core capabilities

Core owns domain-neutral mechanics: typed ID generation, normalization, generic error representation, request/HTTP pipeline, header context, structured logging/tracing, typed search, event envelopes, service pipeline, transaction UnitOfWork, JSON diff, annotation parser/generator and test infrastructure.

Core must not know vendor names, business entity names, application permissions, transport-specific broker storage or event delivery policy. Domain validation remains in its owning service.

## 3. Technical infrastructure contracts

Infrastructure-owned shared contracts, such as broker-delivery DTOs, event-store records, publish options and transport configuration, belong in a separate technical library rather than `libs/backend/core/events`. The actual vendor SDK/transport adapter lives in Infrastructure's `externalService/<capability>/<provider>`.

## 4. Avoiding false reuse

- Shared, cohesive, domain-agnostic conversions and normalization belong in `libs/backend/core/utils` and are exposed exclusively through `utils.go`. Name normalization functions `Normalize<Kind>` (for example `NormalizeEmail`, `NormalizePhone`, `NormalizeUUIDv7`). Do not use separate public core/helpers or core/normalize packages or move capability-local helpers into utils.
- No service-specific event/error/DTO in core.
- No generic repository replacing aggregate-specific SQL and lifecycle rules.
- No external app importing another service's private/internal files.
- No hardcoded production broker URL, passwords or tokens in library defaults.
- No global mutable registry unless concurrency, lifetime and cleanup are rigorously defined.
- No duplicate DTO copied solely to bypass a private package import; publish a deliberate stable alias or boundary contract.

## 5. Documentation for libraries

Each owned library has `docs/engineering/README.md`: public entry points, import rules, SRP, initialization, error/timeout semantics, lifecycle, security, usage and tests. Complex capabilities can add topic-specific engineering guides. **Product execution guides belong to consuming services**, not shared libraries.

Breaking library changes must update consumers, generated code, tests, and owner-local docs together.

## 6. School-specific gateway contract

The school module owns `backend/libs/schoolgateway` (verified session principal, AccessVerifier and the school cookie names). `libs/backend/core/gateway` owns only neutral route tiers and workload verification. The legacy school Admin claim is not a universal permission and must never bypass RBAC. Core must not import any school-owned package, including in tests or generated overlays. [Core boundary audit](../../product/internal/CORE_BOUNDARY_AUDIT.md) records the migration inventory. The Core module has moved and passed standalone build/test checks. Published releases and second-application integration remain pending.

## 7. Independent Core workspace

Core has its own go.mod and go.sum under `libs/backend/core`, with module path `github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core`. Root go.work includes the school module and Core. The school uses a checkout-relative replace for reproducible monorepo builds with GOWORK=off; Core itself has no replace back into school. Use the capability import prefix, not the former school module path. A released tag/private module distribution must be designed and validated separately. Application-specific technical adapters remain under that application's libs. The standalone test procedure is documented in Core's owner runbook.

## Phase 1 capability testing and visibility

Runtime capabilities expose a single public facade. Private implementation packages live under `<capability>/internal/src/`. Capability-specific tests belong under `<capability>/internal/resources/tests/unit-test/` and may import that capability's own internal source. Public behavior is tested under `backend/libs/core/tests/unit-test/`. Never import private packages from other capabilities or production code from test resources. Avoid empty folders and redundant tests. Annotation CLI subcommands are dispatched through annotation/cmd. The annotation capability has no compatibility import paths. Other capabilities undergoing migration must not introduce new legacy imports.

## Operation-first internal organization

A capability with multiple functional areas uses a named folder and matching facade per area, e.g. `internal/src/transaction/transaction.go` and `internal/src/postgres/postgres.go`. Individual use cases are grouped under `operations/<operation>/<operation>.go`. Where an interface is required, the interface lives in the operation entry file and its concrete implementation lives in `operations/<operation>/impl/<operation>.go`. A constructor returns the operation interface and its implementation asserts conformance. Do not create a category-wide `impl/` folder for unrelated operations. Pure functions do not need artificial interfaces or impl folders. Constants, helpers, validators and mappers belong to their owning operation unless genuinely shared across operations. Category facades delegate and do not contain operation logic. Every internal source file must have a matching relative-path unit-test file under `internal/resources/tests/unit-test/`; use separate suites for public API behavior and real infrastructure sanity tests. Do not create empty directories.

Build-time executables under `core/tools` retain their command entry points rather than exposing a runtime facade. Go CLI logic belongs in the tool's `internal/src/operations/` with mirrored tests. Standalone shell tools exported into release bundles must remain self-contained and have command-level regression tests.

### Mandatory Core coverage

Every executable statement in each Core capability must be covered. Internal implementation and external public facade tests are measured independently; neither suite may compensate for missing statements in the other. `backend/libs/core/tools/coverage/check.sh` and the Core 100 percent coverage workflow enforce the exact uncovered block count, not rounded percentages. Declaration-only files have no measurable statements. CLI entry points and shared test support are separately measured; standalone shell commands use command-level regression tests. Do not add exclusions or meaningless tests just to increase coverage.

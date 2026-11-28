# Backend Unit Test Ownership and Structure

**Scope:** Go unit-test placement, fixture ownership and test-framework policy. Generic backend architecture: [backend.md](../backend/backend.md); detailed quality levels: [testing.md](../backend/testing.md). This document is generic and names no specific service.

## 1. Core rule

Place tests with the same **capability/resource/operation owner** as the production code they verify; do not create a global `backend/testing/unit-test/`. Use `internal/tests/unit-test/` for resource-first modules, while respecting established tests in a module that has not yet migrated. Never move unrelated production code or another module's tests merely for symmetry.

```text
backend/
├── launcher/tests/unit-test/
├── services/<module>/
│   └── internal/
│       └── tests/
│           ├── unit-test/
│           │   ├── controller/<resource>/
│           │   ├── validators/<resource>/
│           │   ├── mappers/<resource>/
│           │   ├── repository/<resource>/
│           │   ├── services/<resource>/operations/<verb>/
│           │   └── events/producers/<aggregate>/
│           ├── mocks/
│           └── fixtures/
└── libs/<library>/tests/unit-test/
```

When a test must use unexported symbols, a deliberate same-package `*_test.go` can be placed next to production code instead; this is a Go visibility choice, not an excuse to duplicate the test tree.

## 2. Layer ownership

| Production owner | Focused test responsibility | Not a valid substitute for |
|---|---|---|
| Entity | Schema/serialization compatibility where behavior exists | Trivial field-name assertions |
| Validator | Type/normalization/rule rejection and PATCH presence | Authorization through IAM |
| Mapper | Representation, field mapping, null/secret behavior | Persistence correctness |
| Repository | Constraint error translation and typed query construction | Real PostgreSQL transaction/locking tests |
| Service operation | Lifecycle, idempotency, business rules, orchestration | HTTP binding and provider SDK tests |
| Controller | HTTP status, binding, envelope, principal handoff | Domain invariants |
| Event producer | Version, payload, once-only eligibility | Broker durability or consumer dedup |
| Event consumer | Version validation, idempotency, handler errors | Runtime subscription/replay without broker |
| External adapter | Retry/timeout, provider error mapping, DTO normalization | Live sandbox conformance |
| Bootstrap/worker | Dependency checks, cancellation and shutdown order | Successful business workflow |

## 3. Frameworks and annotations

Use the repository's existing **Ginkgo v2 + Gomega** suite conventions for already established suites. `@Testing` can generate a supported suite registration; `@Mock` may generate typed mocks only when a processor supports its target. Plain Go `testing` and table-driven tests are acceptable for focused pure helpers and repository tests; do not create an artificial suite wrapper solely to meet an outdated layout rule.

When using Ginkgo, maintain one proper suite entry point for the test package; share immutable fixtures/mocks through coherent `fixtures_test.go`/`helpers_test.go` or a small `suite_test.go`. Avoid a huge monolithic suite file.

## 4. Test boundaries and assurance

- Unit tests use deterministic synthetic data, controlled clock/IDs, fake typed interfaces.
- HTTP tests must invoke actual auth/tenant middleware to claim access-control coverage.
- Database constraints, ACID, locks and migration upgrades require PostgreSQL integration tests.
- Broker delivery/replay and transactional outbox guarantees require broker+DB tests.
- Provider OAuth/webhook callbacks require signature/state/replay verification and sandbox tests.
- Cover success, invalid input, not-found, unauthorized, conflict, rollback, retry, concurrent, cancellation, and no-op paths.
- Never put production secrets or user personal data into fixtures/log output.

## 5. Commands and generated code

Use live Makefile targets. Annotation-backed packages may require `-overlay=.generated/overlay.json` when running `go test` from `backend/`; do not assume generated methods exist in maintained source. Test the exact commit under review; report separate build, lint, unit, integration, and CI results.

## 6. Review checklist

Test locations mirror owners, no empty generated directories, API/event compatibility is covered, secrets remain redacted, and owner engineering/product documentation describes meaningful new behavior. Source annotations are not assumed to enforce runtime policy.

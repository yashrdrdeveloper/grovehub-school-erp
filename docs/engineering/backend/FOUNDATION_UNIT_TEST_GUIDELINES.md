# Foundation unit testing — implementation and quality gates

> **Status:** Test framework implemented; Foundation test cases are being
> filled in. Do not claim 100% coverage without the strict gate passing.

## Scope and mirror rule

Every production file with executable behavior under
`backend/services/foundation/internal/app/<relative>.go` must have a maintained
test file under
`backend/services/foundation/internal/tests/unit-test/<relative>_test.go`.
Keep the relative folders and original file stem; tests are separate from app
packages and import production code (never vice versa). Interfaces/constants
with no executable behavior are validated through consumers or compilation.

Example:

```text
internal/app/services/person/operations/create/helpers/person.go
internal/tests/unit-test/services/person/operations/create/helpers/person_test.go
internal/tests/helpers/person/fixtures.go
internal/tests/constants/person/person.go
internal/tests/mocks/person/interfaces.go
```

Mocks, fixtures and constants must be domain owned. Shared generic HTTP/SQL
mock adapters, deterministic clocks and controller helpers live under
`backend/libs/core/tests/`; never import Foundation into core helpers.

## Ginkgo, Gomega and annotations

Use Ginkgo v2 + Gomega for test assertions and GoMock for typed interface
expectations. Comments are recognized by the
`libs/core/annotation/generation/testing` processor:

- `// @Testing suite=PersonCreate` on an unexported or exported no-argument
  test spec registration function, e.g. `func Specs(){Describe(...)}`; the
  generator emits a **virtual** `zz_suite_test.go` containing one `TestX`,
  `RegisterFailHandler`, `RunSpecs` and an enclosing Ginkgo `Describe`.
- `// @Mock` on test-owned interfaces declared under
  `internal/tests/mocks/`; the generator invokes `go.uber.org/mock/mockgen`
  and overlays `zz_mocks.gen.go`. Keep interfaces narrow and rely on GoMock
  expectations, not hand-written stubs that return unconditional success.
- `// @BeforeAll` on one no-argument fixture function per mirrored test
  package; register a Ginkgo `BeforeSuite` callback.
- `// @BeforeEach` on any number of no-argument fixture functions in that
  test package; register them inside the generated `Describe`.

The annotation generator rejects duplicate suites with different names,
multiple BeforeAll hooks, malformed annotations and unexpected placements.
For clarity, suites can spread multiple `@Testing` spec functions across
`*_test.go` files with the same `suite=` value.

## Database, routes and messaging isolation

- GORM repositories use `libs/core/tests/mocks/database.NewPostgres()`,
  an SQLMock-backed GORM dialector; no network or real PostgreSQL instance.
- HTTP controllers and clients use `httptest` and generated service mocks;
  no listening sockets.
- Service operations use mocked repository, mapping, validation, transaction
  and event-producing interfaces; assert calls, ordering, and failures.
- Event producers use a mocked outbox `Store`; do not start NATS.
- Mock context data, UUIDs, clock values and faults explicitly. Check each
  SQLMock expectation and close the fake DB.

## Commands

```bash
make backend-build              # Produces production and test overlays
make foundation-unit            # Runs all mirrored tests and reports actual coverage
make foundation-unit-coverage   # Fails unless 100% statement coverage is achieved
make foundation-test-inventory  # Lists missing mirrored source functions
```

`foundation-unit.sh` executes all
`internal/app/...` packages **and** mirrored suites so packages with no test
coverage are counted as zero, rather than omitted. Go's cover instrumenter
cannot read virtual generated production Go files; the script materializes
the overlay temporarily and removes every file via a cleanup trap. Generated
source remains ignored and non-maintained.

Coverage uses Go's instrumented statement blocks, counted from the
`.generated/foundation-unit-coverage.out` profile. A 100% statement score is
necessary but insufficient for true branch and behavioral coverage; branch,
error and nil paths require independently authored cases. Do not generate
fake "coverage" tests or exclude difficult packages from the denominator.

## Enforcement

The standard backend CI runs the mirrored tests and prints the real coverage
and method inventory. The strict command rejects every uncovered statement;
it must become a mandatory CI gate before considering the full testing
migration complete. Existing suite tests in core/Identity/Integration are
preserved. New suites never need handwritten `suite_test.go` files.

Outstanding until strict gate passes: test every service strategy/flow,
repository persistence and search, controller error path, generated client,
outbox and lifecycle, bootstrap and route registration, and remaining
helpers/validators/DTOs. Never report 100% just because suite execution passes.

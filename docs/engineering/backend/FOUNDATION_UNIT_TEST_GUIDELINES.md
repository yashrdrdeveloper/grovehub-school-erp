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
make backend-build                                  # Regenerates all annotation sources and compiles
make unit-test package=foundation                      # Report Foundation coverage, without strict gate
make unit-coverage package=foundation                  # Requires exactly 100% statement coverage
make unit-coverage packages='["foundation","identity"]' # Independently enforces 100% for each module
make unit-coverage                                     # All services, core and launcher
make test-inventory package=foundation                 # Lists missing mirrored source functions
```

`scripts/testing/unit-coverage.sh` delegates to one shared Go runner. It runs all
selected production packages **and** mirrored suites so packages with no test
coverage are counted as zero, rather than omitted. Go's cover instrumenter
cannot reliably resolve virtual generated Go files when coverage is
instrumented (particularly on macOS). The runner materializes ALL generated
Go packages referenced by the overlay, including transitive dependencies in
`libs/core/` and sibling services. It then runs `go list` and `go test`
WITHOUT `-overlay`, records the physical files in a cleanup manifest,
generates reports, and removes the temporary files on exit. On a subsequent
run, interrupted materialization is detected and safely cleaned.
The annotation generator clears stale ignored module-local `.generated/` output
before every build and Air recompilation so deleted or renamed packages cannot
remain in the overlay. Generated source remains ignored and non-maintained.

Coverage uses Go's instrumented statement blocks, counted from the
`.generated/unit-coverage/<package>.out` profile, with an HTML report under
`.runtime/reports/unit-coverage/<package>.html`. A 100% statement score is
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

## Shared runner and selectors

All backend services use the same unit-test runner: `make unit-test`
for an informational coverage report and `make unit-coverage` for the strict
100% statement gate. `package=foundation` selects one service;
`packages='["foundation","identity"]'` selects multiple services. The
unquoted `packages=[foundation,identity]` form is also accepted. With no
selector or an empty list, all services plus `core` and `launcher` are
selected. Coverage and HTML files are separate per package, and no service is
silently dropped from the denominator. Missing unit-test suites count as 0%.

On every `make backend-build`, `make backend-test`, `make backend-dev`
rebuild, or `make unit-coverage`, annotation generation starts with a fresh
ignored output tree. Air watches Go and resource assets including SQL, YAML,
JSON and template files; editing such assets also rebuilds the binary with
fresh embedded resources. The code generators do not write maintained Go
source files; generation paths stay isolated inside `.generated/`.

## macOS generated-source coverage troubleshooting

The Go build overlay remains the default for normal `backend-build` and Air
compilation. For **unit coverage only**, Go cover must read generated source
from its original package directory. Running `go test -overlay` with coverage
can produce errors such as `open services/.../.generated/.../zz_suite_test.go:
no such file or directory` on macOS. The shared runner avoids that combination;
it never leaves generated test or mock files in maintained directories after a
completed invocation. CI exercises the same runner on both Linux and macOS.

`make unit-coverage package=foundation` will still report a failure until
all Foundation statements are genuinely covered. Do not confuse that expected
100% threshold failure with a source-generation or compiler failure.

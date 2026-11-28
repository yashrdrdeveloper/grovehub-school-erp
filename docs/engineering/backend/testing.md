# Backend Testing, CI and Release Acceptance

## 1. Test ownership

```text
services/<module>/internal/resources/tests/
├── unit-test/
│   ├── controller/<resource>/
│   ├── validators/<resource>/
│   ├── services/<resource>/operations/<verb>/
│   ├── repository/<resource>/
│   └── events/producers/<aggregate>/
├── mocks/                      # module-owned fakes
└── fixtures/                   # deterministic valid/invalid data

libs/<library>/tests/unit-test/ # generic capability tests when relevant
launcher/tests/unit-test/       # composition, middleware, entrypoints
```

Placement mirrors code ownership. Same-package tests are allowed for Go `internal` visibility. Do not create empty test trees or revive a global `backend/testing/unit-test` bucket.

## 2. Testing levels

| Level | Purpose | Dependencies |
|---|---|---|
| Unit | One validator, mapper, helper, strategy or operation | Fakes/fixtures |
| Contract | Clients, OpenAPI, provider DTO, event schema | Schema + typed contracts |
| Repository integration | SQL, constraints, indexes, migrations | Real PostgreSQL |
| Messaging integration | Outbox, broker retries, deduplication | Real broker + DB |
| HTTP integration | Middleware, route and response envelope | Local server/handler |
| E2E | User-visible workflow | Isolated test environment |
| Performance | Search, hydration, dispatch, throughput | Representative loads |

Mocks do **not** prove ACID, broker durability, or IAM enforcement.

## 3. Test annotations

`@Testing` (suite generation), `@BeforeAll`, `@BeforeEach`, and `@Mock` may be used only with their implemented generators. A cataloged proposed directive is not automatically functional. Prefer Ginkgo/Gomega for established project suites; plain Go table tests are acceptable when clearer.

## 4. Minimum test matrix

- Happy path, missing resource, invalid input, malformed/unknown JSON.
- Absent vs explicit-null PATCH, Unicode normalization and query field allowlist.
- Principal missing/expired, forbidden action, wrong tenant/owner.
- No-op and idempotent replay; same key with conflicting fingerprint.
- Concurrent version conflicts, duplicate contact/identifier, transaction rollback.
- Event schema/version, duplicate/out-of-order delivery, outbox redelivery crash window.
- Provider timeout, 429/retry, webhook replay/signature/credential failure.
- Cancellation and worker shutdown; secrets absent in logging.
- Clean-install migration and upgrade from last released version.

Tests must assert both results and side effects, not just a nil error.

## 5. Generation and build gates

Use live repository scripts and Makefile, including the Go overlay for generated annotation code. Typical entry points:

```bash
make backend-build
make backend-lint
make unit-test
make backend-test
```

Check Makefile definitions before running; avoid concurrent generator jobs rewriting one overlay. A `go test` that runs without an overlay may not compile generated client/mapper/search calls.

## 6. Safe delivery

Confirm API/event compatibility, secret handling, database migration history, build/test status for the **exact commit**, background job shutdown, and updated engineering/product docs. Report each of annotation checks, compilation, unit tests, integration tests, database tests and CI as **passed**, **failed**, **pending** or **not run**; don't claim success based on an older commit's green job.

## 7. Independent Core verification

`make backend-test` runs both the school module and Core module test trees; generation continues to emit owner-local files. `make unit-test package=core` and `make unit-coverage package=core` use the separate Core module root. The latter retains the strict 100% requirement. Core-only validation from the monorepo root:

```bash
cd libs/backend/core
export GOWORK=off
go run ./annotation/cmd testing "$PWD" "$PWD"
go run ./annotation/cmd overlay "$PWD" "$PWD"
go test -vet=off -overlay=.generated/overlay.json ./...
```

A Core-only empty overlay is valid; school generation still rejects absent generated application output. CI also builds the school with GOWORK=off using the committed monorepo-relative replacement. Live migration, database and broker integration tests remain separate gates.

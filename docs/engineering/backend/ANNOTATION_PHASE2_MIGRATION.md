# GrooveHub Annotation Framework — Phase 2 Migration Ledger

> **Status:** Partial migration in progress; **not** a complete backend refactor.
> **Branch:** `phase-2.1-identity-account`.
> **Phase 1:** [complete annotation catalog](ANNOTATION_PHASE1_REFERENCE.md).
> **Architecture:** [master plan](ANNOTATION_FRAMEWORK_MASTER_PLAN.md).

## 1. Non-negotiable boundaries

- Do not replace business workflows with annotations. `services/impl` remains responsible for application rules, cross-service orchestration, ID creation, lifecycle, authorization, and appropriate event publishing.
- Do not generate repository SQL from unverified method names. Existing GORM queries, transactions, optimistic concurrency, search projections, constraint handling and tenant scope remain explicit.
- An annotation in a comment is **not an enforced runtime policy** until a generator or wrapper explicitly implements it and tests prove equivalent behavior.
- Generated sources physically live only in `services/<service>/.generated/<package>/`; `go build -overlay=...` maps them into existing packages.
- Preserve the core libraries (`database`, `events`, `http`, `errors`, `logging`, `telemetry`, `normalize`) rather than introducing a competing runtime container.

## 2. Implemented migration slices

| Area | Code moved or replaced | Annotation / API | Current limitation |
| --- | --- | --- | --- |
| Foundation HTTP routes | Replaced five manual mux entries with generated registrations | `@Routes`, `@POST`, `@GET`, `@PATCH` on `controller.PersonController` | Needs full clean-checkout build and runtime smoke test |
| Identity HTTP routes | Replaced manual registration for `POST /auth/register` | `@Routes`, `@POST` | Needs auth registration HTTP regression test |
| Integration HTTP routes | Replaced 13 manual connection, webhook and provider route entries | `@Routes`, `@GET`, `@POST`, `@DELETE` on two controller interfaces | Needs OAuth callback, webhook, and ownership regression tests |
| HTTP build tooling | New `generation/routes` tool, generalized overlay maps service-local generated package files | Go AST, `-overlay` | Gopls and standard direct Go invocations need explicit configuration |
| Foundation, Integration and Infrastructure clients | 15 in-process forwarding methods now generated; constructors and injected service contracts remain maintained Go | `@Client` and `@Proxy` | Must verify all delegations compile and preserve context, signatures, auth boundaries |
| Integration config | Gmail and Outlook env parsing moved to typed tagged loader; preserved callback checks and error semantics | `@Configuration`, `env`, `required`, `default`, `secret` | `GmailPush` and encryption-key loaders remain explicit |
| Foundation DTO validation | Create/update name trimming, notblank and rune-length limits moved to DTO tags and shared validator | `@Validator`, `normalize`, `validate` tags | Contact uniqueness/primary/search validation stays custom |
| Test tooling | Added overlay-aware `make backend-test`; sanity tests prepare overlay | `go test -vet=off -overlay=...` | Go vet cannot automatically read non-existent virtual paths; see below |

The source retains service methods, service interfaces, client constructors, GORM repositories, domain event decorators, provider strategies, authentication decisions, and middleware behavior.

### Endpoint parity

The `generation/routes/main_test.go` suite asserts that all **19** canonical route patterns are present:

- Foundation: 5 routes
- Identity: 1 route
- Integration: 13 routes

The generator checks duplicate method/path pairs, interface/method signatures, and incomplete route declarations before emitting source. Generation is idempotent.

## 3. Current implementation structure

```text
backend/
  libs/core/annotation/
    annotation.go                     # Phase 1 metadata registry
    annotations_*.go                  # all declared annotation categories
    configuration/configuration.go    # typed env/default/required loader
    validation/validation.go          # DTO string/*string normalization and validation
    generation/
      expand/                         # existing Foundation @Entity preprocessor
      gormtable/                      # existing GORM table generator
      routes/                         # NEW Go AST -> traced HTTP route generation
      overlay/                        # expanded to all service-local generated packages
  services/
    foundation/
      controller/person.go            # annotated route contract
      dto/request/person.go           # declarative validation tags
      validators/impl/person.go       # application error mapping and contact checks
      routes/http.go                  # composition; delegates to generated routes
      .generated/routes/zz_routes.gen.go
    identity/
      controller/auth.go              # annotated route contract
      routes/http.go                  # generated-route composition
      .generated/routes/zz_routes.gen.go
    integration/
      controller/{connections,provider}.go  # annotated route contracts
      config/config.go                # typed env contracts
      config/impl/config.go           # startup loader and provider-specific validation
      routes/http.go                  # composition; delegates to generated routes
      .generated/routes/zz_routes.gen.go
  .generated/overlay.json             # generated Go -overlay manifest

scripts/backend/generate-entities.sh   # entity -> routes -> overlay
scripts/backend/build.sh               # overlay-aware build
scripts/backend/air-build.sh           # overlay-aware hot reload
scripts/backend/test.sh                # overlay-aware go test
scripts/testing/sanity.sh              # overlay-aware sanity test
```

Generated files are **not committed**; directories above show runtime outputs after the generation step.

`generation/clients` adds `.generated/client/impl/zz_clients.gen.go` within Foundation, Integration and Infrastructure. The Go overlay scanner now supports nested generated package paths. `client/impl` retains only the constructor and injected service dependencies; generated forwarders are defined by `@Proxy` method bindings on each public client interface.

## 4. Development and verification commands

```bash
git pull origin phase-2.1-identity-account
make backend-build
make backend-test
make sanity
```

When invoking Go directly:

```bash
bash scripts/backend/generate-entities.sh
cd backend
go build -overlay=.generated/overlay.json ./cmd/platform
go test -vet=off -overlay=.generated/overlay.json ./...
```

**Important Go tooling limitation:** A virtual generated `.go` file whose logical path does not physically exist can cause the implicit `go vet` step in `go test` to fail trying to open that logical path. The project test wrapper uses `-vet=off` for generated packages. This does not replace dedicated static analysis or make every IDE overlay-aware.

### Verification status

- Source updates committed on the integration branch.
- Pure generator test code and route parity expectations added; startup config and DTO validation unit tests added.
- The full Go 1.27.x backend build, Ginkgo/Gomega suite, and HTTP/integration smoke tests **have not been executed** in this environment.
- No claim of full runtime compatibility or 232 executable annotations is made.

## 5. Remaining backend migration work

| Responsibility | Next safe migration | Required checks |
| --- | --- | --- |
| Entity/GORM models | Consider `@Entity` for Identity and Integration only after allowing safe accessor exclusions | No accidental credential/PII setters; schema/table parity |
| Mapper/DTO | Typed MapStruct-style mapper with explicit nested/converter mappings | No secret export, exhaustive destination fields, error DTO semantics |
| Validators | Extend field tag rules for email/phone/UUID, cross-field rules, PATCH presence | Preserve existing resource errors and normalization |
| Clients | Pilot applied to Foundation, Integration and Infrastructure; extend only after compile tests | No authorization bypass, no signature drift |
| Repositories | Opt-in CRUD, transaction and query contracts where safe | Preserve SQL/GORM, tenant boundaries, search, conflict mapping |
| Events | Typed domain event descriptors and handlers | Keep `Envelope[T]`, event topics, idempotent/no-op semantics |
| Messaging | Producer/consumer registration descriptors | NATS/Watermill ack/retry/trace/backpressure semantics |
| External services | Register provider Strategy/Factory/Adapter metadata | Gmail/Outlook OAuth, webhook and token safety |
| HTTP middleware | Verified wrappers for authentication, authorization, limits, tracing | Stable ordering, no trust bypass |
| Logging, metrics, security | Typed bounded instrumentation and redaction | No secrets/PII or unbounded labels |
| Launcher/server/jobs | Typed lifecycle config, dependency checks and graceful shutdown | Explicit shutdown and development-only safeguards |
| Testing/utilities | Reusable fixture/mock generation and cache contracts | Determinism, safe concurrency, stable compatibility |
| GoF patterns | Start with existing Adapter/Strategy/Decorator/Proxy roles | Patterns stay explicit when generation adds no value |

These remain work items, **not implemented features**.

## 6. Acceptance gate before widespread migration

1. Run a clean checkout build using the pinned Go toolchain and generator versions.
2. Run the full test suite and compare route patterns with the former manual definitions.
3. Verify Foundation create/update/search/deactivation and all event publishing behavior.
4. Verify Identity registration, credential handling and cross-module Foundation Person calls.
5. Verify Integration Gmail/Outlook connection lifecycle, OAuth callbacks, webhook authenticity and provider actions.
6. Verify NATS consumer delivery, durable subscriptions, tracing context and failure/retry behavior.
7. Verify safe handling of tenant boundaries, audit data, PII/credentials and the local-only development principal.
8. Measure generation overhead, Air reload latency and config/validation performance before further migrations.

**Do not proceed with automatic replacement of service, repository, messaging or security logic until these gates pass.**

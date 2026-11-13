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
| Foundation HTTP routes | Moved Person controller contracts and implementation to `controller/protected/`; five endpoints now use `/foundation/protected/persons` | `// @Routes /foundation` on `routes.RegisterHTTP`; `// @Routes /persons` on protected PersonController | Route tiers are URL grouping only, without authentication enforcement; update clients and enforce gateway policy before production |
| Identity HTTP routes | Replaced manual registration for `POST /auth/register` | `@Routes`, `@POST` | CI compilation verified; auth registration HTTP regression pending |
| Integration HTTP routes | Replaced 13 manual connection, webhook and provider route entries | `@Routes`, `@GET`, `@POST`, `@DELETE` on two controller interfaces | Needs OAuth callback, webhook, and ownership regression tests |
| HTTP build tooling | New `generation/routes` tool, generalized overlay maps service-local generated package files | Go AST, `-overlay` | Gopls and standard direct Go invocations need explicit configuration |
| Foundation, Integration and Infrastructure clients | Generated 15 typed in-process forwarding methods; maintained constructors and service contracts retained | `@Client`, `@Proxy` | CI compilation verified; runtime context and auth tests pending |
| Integration config | Gmail and Outlook env parsing moved to typed tagged loader; preserved callback checks and error semantics | `@Configuration`, `env`, `required`, `default`, `secret` | `GmailPush` and encryption-key loaders remain explicit |
| Shared PATCH DTO generation | Replaced Foundation Person-only PATCH helper and custom JSON serializer with `annotation/dto.Patch[T]` and `@PatchRequestDto` | `generation/patchdto` emits JSON serialization using service-local Go overlay | Omitted/null/value semantics preserved; validation of clearable fields stays in the service |
| Foundation and Identity request validation | Shared typed HTTP executor now binds path/query/body inputs, runs canonicalization + DTO/custom validation, calls typed service adapters and writes uniform responses; Foundation Person and Identity Register controllers migrated | `rest.Run(w,r,Config[P,Q,B,R])`, `rest.Input`, `rest.Result`, `annotationvalidation.Rules[T]` | No authentication added; HTTP callback adapters retain existing service contracts; input binding supports common primitives, optional pointers and repeated query slices |
| Test tooling | Added overlay-aware `make backend-test`; sanity tests prepare overlay | `go test -vet=off -overlay=...` | Go vet cannot automatically read non-existent virtual paths; see below |
| Annotation safety check | Build-time AST scanner rejects unimplemented directives such as `@Authorize` before generating code | `generation/check` | Does not replace actual IAM middleware or authorization decisions |
| Strict DTO projection | Identity registration response and Foundation Person error metadata generated; contact grouping remains custom | `@Mapper`, `@Mapping` | Explicit fields fail generation when destination gains unmapped fields |
| GORM table-only entities | Removed manual `TableName()` methods for seven Identity/Integration models; no new credential setters | `@Entity ... mode=table` | CI compilation verified; DB schema unchanged |
| Passive entity refactor | Removed generated entity builders and entity-owned New/Get/Set/helper methods in Foundation and Identity; replaced usage with service-owned constructors, copies and explicit updates | `@Entity` now maps tables only | Go 1.27 CI compilation and existing unit suites passed; live DB regression remains pending |
| Annotated repository search | Concrete `repository/impl/PersonRepository.Search` delegates to a generated private `searchGenerated` helper using `core/search.Execute`; filters/sorts/cursor definitions remain in `repository/helpers/search/person.go` | `@Searchable config=search.PersonConfig alias=p` above the implementation method | The public `@Repository` interface has no search implementation annotation; optional `@Filterable`, `@Sortable`, `@Cursor` now belong on the concrete Search method |
| Repository responsibility cleanup | Moved shared `UnitOfWork` to `libs/core/database/transaction`, sentinel errors to `repository/constants/errors/`, non-query error helpers to `repository/helpers/errors/`, named PostgreSQL constraint constants to `repository/constants/{person,contact}.go`, and Person search configuration plus its private contact filter helper to `repository/helpers/search/person.go` | Thin interfaces, resource-owned implementation details, JavaDoc-style comments above annotations | Full backend compilation and Foundation unit test gates; live PostgreSQL regression pending |
| Person/contact separation | Removed non-column `Person.Contacts`; repository and service pass child `[]Contact` explicitly, while mapper emits unchanged response DTO; renamed Go idempotency fields with old database column tags | GORM `@Entity` and existing `PersonRepository` | Go build/CI verification; live DB regression remains pending |
| Event identity declarations | Foundation Person changed, Identity Account created, Integration Mailbox changed topics/types/versions generated | `@Event` | Event publication timing, after-commit, error policies unchanged |
| CI regression gate | GitHub Actions Go 1.27.3 binary compilation, package resolution, core/launcher/Identity/Integration unit tests, and generated-output checks | `.github/workflows/backend-annotations.yml` | Live database/NATS integration regression remains separate |

The source retains service methods, service interfaces, client constructors, GORM repositories, domain event decorators, provider strategies, authentication decisions, and middleware behavior.

**Generator change:** `generation/entities` no longer implicitly generates constructors/builders/getters/setters from `@Entity` and fails on combinations of `@Entity` with `@Builder`/`@Getter`/`@Setter`. Explicit opt-in builder/accessor directives remain available only for non-entity types. `generation/gormtable` remains the sole `@Entity` codegen for all services. The deprecated Gombok staging/expand path remains inactive.

### Endpoint parity

The route generator regression suite covers the Foundation five-method protected Person contract, service-prefix changes, example public and restricted controllers, legacy Identity URLs, duplicate routes, missing category/prefix errors, and deterministic output. The build still generates **19 endpoints total**: Foundation 5 (now under `/foundation/protected/persons`), Identity 1 (legacy `/auth/register`), and Integration 13 (legacy `/integrations`). The generator checks duplicate method/path pairs and controller signatures; it does **not** implement access control.

## 3. Current implementation structure

```text
backend/
  libs/core/annotation/
    annotation.go                     # Phase 1 metadata registry
    contract/contract.go             # shared acyclic annotation contracts
    inventory.go                     # aggregates category-owned packages
    entity/annotations.go            # entity declarations (builder also here)
    event/annotations.go             # event contracts
    validation/annotations.go        # validation directives
    utilities/annotations.go         # utility directives
    patterns/{creational,structural,behavioral,patterns}.go
    mapping/ persistence/ composition/ http/ middleware/ security/
    external/ messaging/ observability/ server/ testing/  # other catalogs
    configuration/configuration.go    # typed env/default/required loader
    validation/validation.go          # DTO string/*string normalization and validation
    generation/
      entities/                       # Gombok-compatible AST builder/accessor generator
      gormtable/                      # @Entity TableName generator (full or table-only)
      patchdto/                       # @PatchRequestDto JSON serialization generator
      mappers/                        # exact-field DTO projection generator
      search/                         # typed repository search and allowlist checker
      events/                         # versioned topic/type constant generator
      routes/                         # Go AST -> traced HTTP route generation
      clients/                        # typed @Client/@Proxy delegation
      overlay/                        # expanded to all service-local generated packages
  services/
    foundation/
      controller/protected/person.go  # protected Person HTTP contract
      dto/request/person.go           # Person create/PATCH DTO, using @PatchRequestDto
      dto/response/person.go          # Person response DTO
      # shared PATCH field: backend/libs/core/annotation/dto/patch.go
      validators/impl/person.go       # application error mapping and contact checks
      routes/http.go                  # @Routes /foundation service prefix and generated registrations
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

scripts/backend/generate-annotations.sh # entity -> tables -> routes -> clients -> search -> mappers -> events -> overlay
scripts/backend/generate-entities.sh    # compatibility wrapper
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
bash scripts/backend/generate-annotations.sh
cd backend
go build -overlay=.generated/overlay.json ./cmd/platform
go test -vet=off -overlay=.generated/overlay.json ./...
```

**Important Go tooling limitation:** A virtual generated `.go` file whose logical path does not physically exist can cause the implicit `go vet` step in `go test` to fail trying to open that logical path. The project test wrapper uses `-vet=off` for generated packages. This does not replace dedicated static analysis or make every IDE overlay-aware.

### Verification status

- Source updates committed on the integration branch.
- Annotation-framework test cases were removed as requested; other existing backend unit tests remain in the repository.
- The **complete backend binary has compiled successfully** on GitHub Actions using Go **1.27.3** and clean-checkout generation. CI checks the annotation packages resolve and runs existing core/launcher/Foundation/Identity/Integration unit suites. The **live database/NATS/HTTP integration smoke tests have not yet been executed**.
- No claim of full runtime compatibility or 232 executable annotations is made.

## 5. Remaining backend migration work

| Responsibility | Next safe migration | Required checks |
| --- | --- | --- |
| Entity/GORM models | Identity and Integration tables now use `@Entity schema.table mode=table`; remaining work is parity and migration tooling verification | No accidental credential/PII setters; GORM schema/table parity |
| Mapper/DTO | Strict `@Mapper` generator active for Identity registration and Foundation error metadata; complex mappings remain handwritten | No secret export, exhaustive destination fields, error DTO semantics |
| Validators | Extend field tag rules for email/phone/UUID, cross-field rules, PATCH presence | Preserve existing resource errors and normalization |
| Clients | Pilot applied to Foundation, Integration and Infrastructure; extend only after compile tests | No authorization bypass, no signature drift |
| Repositories | Opt-in CRUD, transaction and query contracts where safe | Preserve SQL/GORM, tenant boundaries, search, conflict mapping |
| Events | `@Event` generates Foundation, Identity and Integration topic/type/version identifiers; event handlers and publishers remain explicit | Keep `Envelope[T]`, event topics, idempotent/no-op semantics |
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

**Confirmed boundary:** `@Repository` validates interface entity ownership but does not generate constructors, transactions, or persistence. `@Searchable` belongs to the concrete `Search` implementation and generates only a private `searchGenerated` helper using a resource-owned Config (which provides field allowlists and stable cursor definition); optional `@Filterable`, `@Sortable`, and `@Cursor` check parity with that configuration. `@Service`, `@Consumer`, `@Authorize`, `@Transactional`, `@Retry`, `@Audit`, and most design-pattern annotations remain metadata-only, not effective runtime policy. No endpoint must rely on an inert annotation for protection.

## Annotation folder refactor

The 16 category definitions were moved from root-level `annotations_*.go` files into separate category packages under `libs/core/annotation/`, including `event/`, `utilities/`, and `patterns/`. Shared types live in `contract/`, and the original root package remains a compatibility facade. All annotation-framework `*_test.go` files and the former `libs/core/tests/unit-test/annotation/` suite were removed as requested. No Foundation/Identity/Infrastructure/Integration business logic, generators, or runtime imports were changed by this folder-only migration. CI still compiles the backend and runs unrelated existing unit suites.

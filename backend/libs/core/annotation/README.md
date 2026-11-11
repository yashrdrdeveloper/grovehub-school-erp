# GrooveHub Core Annotation Framework

> **Status:** Phase 1 catalog completed; Phase 2 generators active in selected backend layers. Other annotations remain metadata only.
>
> [Complete master plan](../../../../docs/engineering/backend/ANNOTATION_FRAMEWORK_MASTER_PLAN.md) · [Phase 1 catalog](../../../../docs/engineering/backend/ANNOTATION_PHASE1_REFERENCE.md) · [Phase 2 migration ledger](../../../../docs/engineering/backend/ANNOTATION_PHASE2_MIGRATION.md)

Go doesn't natively support Java-style annotations. GrooveHub accepts source comments such as `// @Entity ...`, parses them at build time and generates ordinary typed Go where a processor exists. No annotation alters a function automatically at runtime.

## Inventory and code structure

The Phase 1 registry defines the catalog of supported annotations, including the new executable `@PatchRequestDto`, and all 23 GoF design patterns. each category's `annotations.go` owns its declarations; `inventory.go` and `annotation.go` expose `All`, `Lookup`, `Uses`, `Patterns`, `ParseComment` and `ValidateDirective`.

```text
backend/libs/core/annotation/
  annotation.go              # stable public registry and API
  inventory.go               # imports category-owned definitions
  contract/contract.go       # shared types; avoids circular dependencies
  entity/                    # entity.go + annotations.go
  validation/                # validation.go + annotations.go
  dto/                       # shared Patch[T] and @PatchRequestDto catalog
  configuration/             # configuration.go + annotations.go
  mapping/annotations.go
  persistence/annotations.go
  composition/annotations.go
  http/annotations.go
  middleware/annotations.go
  security/annotations.go
  external/annotations.go
  event/annotations.go
  messaging/annotations.go
  observability/annotations.go
  server/annotations.go
  testing/annotations.go     # definitions, not unit tests
  utilities/annotations.go
  patterns/
    creational.go
    structural.go
    behavioral.go
    patterns.go              # all 23 GoF decisions
  generation/
    entities/ gormtable/ routes/ clients/ mappers/ events/
    search/ patchdto/ check/ parser/ overlay/ expand/
```

All code generated for a service is placed under **that service's ignored `.generated/` directory**, not beside its maintained Go files.

## Active annotation examples

### Entities and GORM tables

```go
// @Entity foundation.person_contacts
type Contact struct {
    Value string
}
```

This generates **only** GORM `TableName()`. No `NewContactBuilder`, `NewPersonBuilder`, entity constructors, implicit getters or setters are produced. Foundation and Identity entity construction, optional strings, mutable contact copies and lifecycle updates now belong to their service implementations/helpers.

For sensitive records:

```go
// @Entity identity.password_credentials mode=table
type PasswordCredential struct {
    PasswordHash string
}
```

**Only `TableName()` is generated** for every `@Entity` model. The `mode=table` suffix remains a backward-compatible alias for the default.

### HTTP routes

```go
// @Routes /auth
type AuthController interface {
    // @POST /register
    Register(http.ResponseWriter, *http.Request)
}
```

`generation/routes` validates method/path contracts and emits `net/http` handlers with the existing tracing middleware. All 19 Foundation/Identity/Integration endpoints were preserved.

### Typed clients and mappers

```go
// @Client messagingClient
type MessagingClient interface {
    // @Proxy service.Health
    Health(ctx context.Context) error
}
```

This generates a signature-preserving delegation to an explicitly injected implementation. No service locator or network call is invented.

```go
// @Mapper AuthMapper
type AuthMapper interface {
    // @Mapping Email=account.GetEmail()
    // @Mapping Status=string(account.GetStatus())
    // ... every other destination field must also be declared
    ToRegisterResponse(account entity.Account) response.Register
}
```

The mapper generator requires **complete destination-field coverage** and fails rather than silently dropping fields. Complex Foundation contact grouping stays handwritten.

### Annotated repository search

```go
/**
 * PersonRepository owns foundation.persons and the derived SearchDetails.
 * ContactRepository owns contact rows; PersonService coordinates mutations.
 */
// @Repository entity=Person
type PersonRepository interface {
    /**
     * Search uses trusted field mappings from search.PersonConfig.
     */
    Search(ctx context.Context, query coresearch.Query[string]) (coresearch.QueryResult[entity.Person], error)
}
```

Place `@Searchable` on the concrete `Search` method in `repository/impl/person.go`, not the interface:

```go
// @Searchable config=search.PersonConfig alias=p
func (repository *PersonRepository) Search(
    ctx context.Context, query coresearch.Query[string],
) (coresearch.QueryResult[entity.Person], error) {
    return repository.searchGenerated(ctx, query)
}
```

`generation/search` generates only the private `searchGenerated` helper in the existing service-local source overlay. This retains an explicit, reviewable `Search` implementation in `impl/` while avoiding handwritten GORM query plumbing. The build-time `generation/search` processor resolves external repository search configuration (including nested `repository/helpers/search/`), verifies that the cursor expression matches a stable sort key, optionally checks `@Filterable`/`@Sortable`/`@Cursor` declarations when provided, and generates the private `searchGenerated` helper into the service's ignored `.generated/repository/impl/` folder. Execution delegates to the shared `core/search.Execute`. The resource validator still controls operators and normalized values; repository-owned `coresearch.Config` keeps all SQL expressions explicit and reviewed. `@Repository` is ownership metadata validation only, **not** DI or transaction generation.

### Reusable PATCH request DTOs

Annotate a DTO to generate its `MarshalJSON` method; no resource-local
`helpers/patch.go` or handwritten `MarshalJSON` is needed:

```go
import annotationdto "grovehub/libs/core/annotation/dto"

// @PatchRequestDto
type UpdatePerson struct {
    FirstName annotationdto.Patch[string]   `json:"firstName"`
    MiddleName annotationdto.Patch[string]   `json:"middleName"`
    Contacts  annotationdto.Patch[Contacts] `json:"contacts"`
}
```

The shared `Patch[T]` distinguishes omitted, explicit `null`, and supplied values
using `Set` and `Value`. The build-time `generation/patchdto` tool validates
the annotated DTO shape and generates JSON serialization into the service's
ignored `.generated/dto/request/zz_patchdto.gen.go` file for flat request packages (or a corresponding nested path for other services), mounted by the Go overlay.
Business validators still decide which fields may be cleared; annotations
do not create authorization or validation rules. No generated code is checked in.

### Configuration, validation and events

Foundation uses a plain, non-generic `Person` name DTO embedded in `CreatePerson`; `UpdatePerson` uses `@PatchRequestDto` and shared presence-aware `annotation/dto.Patch[T]` fields for PATCH. Validation is enforced by `normalize` / `validate` struct tags through the existing validator, not redundant field-level annotation comments. `headercontext` captures validated HTTP metadata once at the edge rather than storing it in JSON DTOs.


Integration config uses `// @Configuration` and `env`, `required`, `default`, `secret` struct tags with the typed startup loader. Foundation and Identity request DTOs use `// @Validator` and explicit `normalize` / `validate` tags for the supported rules. Cross-field business checks remain in their existing validators.

```go
// @Event topic=identity_account_created type=identity.account.created.v1 source=identity
type Created struct {
    AccountID string
}
```

This generates `CreatedTopic`, `CreatedType`, `CreatedVersion` and `Source`, **not** the event publisher, consumer or after-commit transaction policy. The original NATS/Watermill event envelopes and publication behavior are unchanged.

## Build and verification

```bash
make backend-build
make backend-test
bash scripts/backend/generate-annotations.sh
```

The unified generator first checks annotation validity, then runs entities → GORM table metadata → routes → clients → repository searches → PATCH DTO JSON → mappers → events → Go overlay. The `generation/check` tool blocks unsupported or metadata-only directives from being used as if they enforce security or transaction policy. CI at `.github/workflows/backend-annotations.yml` compiles the backend with Go 1.27.3, verifies annotation packages resolve, and runs existing non-annotation backend unit suites. Annotation-framework test cases were removed as requested; this does not disable other backend tests. Direct Go commands need the `-overlay` option; the test wrapper uses `-vet=off` for virtual files and is not a substitute for separate static analysis.

## Security and migration limits

`@Service`, `@Consumer`, `@Producer`, `@Authorize`, `@Transactional`, `@Audit`, `@Retry` and most other cataloged directives **do not enforce behavior yet**. `@Repository` only validates ownership metadata; it does not enforce transaction boundaries. Do not use them as security, durability, or transaction guarantees. Keep authorization, domain decisions, GORM predicates, credential/OAuth policy, message acknowledgments and event publication explicit until tested processors are separately implemented.

The complete remaining work and acceptance criteria are in the Phase 2 migration ledger.

## Ownership and dependency direction

Each responsibility directory owns an exported `UseCases []contract.UseCase` catalog. The shared `contract/` package has no dependency on the registry, processors or services. The root `annotation` package imports the domain catalogs, exposes the unchanged `All`, `Lookup`, `Uses`, `Patterns`, `ParseComment` APIs and re-exports their contract types. The pattern catalog lives only in `patterns/`, split by the three GoF families. Build-time generators and their output directories remain unchanged.

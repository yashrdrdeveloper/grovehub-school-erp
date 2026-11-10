# GrooveHub Core Annotation Framework

> **Status:** Phase 1 catalog completed; Phase 2 generators active in selected backend layers. Other annotations remain metadata only.
>
> [Complete master plan](../../../../docs/engineering/backend/ANNOTATION_FRAMEWORK_MASTER_PLAN.md) · [Phase 1 catalog](../../../../docs/engineering/backend/ANNOTATION_PHASE1_REFERENCE.md) · [Phase 2 migration ledger](../../../../docs/engineering/backend/ANNOTATION_PHASE2_MIGRATION.md)

Go doesn't natively support Java-style annotations. GrooveHub accepts source comments such as `// @Entity ...`, parses them at build time and generates ordinary typed Go where a processor exists. No annotation alters a function automatically at runtime.

## Inventory and code structure

The Phase 1 registry defines **226 use cases, 232 distinct annotation names and all 23 GoF design patterns**. `annotations_*.go` organizes metadata by responsibility; `inventory.go` and `annotation.go` expose `All`, `Lookup`, `Uses`, `Patterns`, `ParseComment` and `ValidateDirective`.

```text
backend/libs/core/annotation/
  annotation.go                    # metadata registry and parser
  annotations_*.go                # 16 responsibility catalogs
  configuration/                  # typed startup env loader
  validation/                     # safe request field rules and normalization
  entity/                         # existing generic builder facade
  generation/
    entities/                     # Gombok-compatible builders, getters and setters
    gormtable/                    # GORM TableName() generation
    routes/                       # traced net/http registration
    clients/                      # typed in-process client forwarding
    mappers/                      # exact-field DTO projections
    events/                       # stable domain event topic/type/version constants
    overlay/                      # virtual Go package source injection
    parser/                       # read-only catalog AST scan
    expand/                       # retained legacy tool; not in active build path
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

This generates `NewContactBuilder()`, supported getters and setters, and GORM `TableName()`. Existing custom Person value getters/setters are preserved; the generator suppresses conflicting accessor categories.

For sensitive records:

```go
// @Entity identity.password_credentials mode=table
type PasswordCredential struct {
    PasswordHash string
}
```

**Only `TableName()` is generated**; no builder or mutator is emitted. Identity and Integration use this mode for seven persistence models.

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

### Configuration, validation and events

Integration config uses `// @Configuration` and `env`, `required`, `default`, `secret` struct tags with the typed startup loader. Foundation and Identity request DTOs use `// @Validator` and explicit `normalize` / `validate` tags for the supported rules. Cross-field business checks remain in their existing validators.

```go
// @Event topic=identity_account_created type=identity.account.created.v1 source=identity
type Created struct {
    AccountID string
}
```

This generates `CreatedTopic`, `CreatedType`, `CreatedVersion` and `Source`, **not** the event publisher, consumer or after-commit transaction policy. The original NATS/Watermill event envelopes and publication behavior are unchanged.

## Build and tests

```bash
make backend-build
make backend-test
bash scripts/backend/generate-annotations.sh
```

The unified generator runs entities → GORM table metadata → routes → clients → mappers → events → Go overlay. CI at `.github/workflows/backend-annotations.yml` compiles the backend with Go 1.27.3 and runs annotation tests. Direct Go commands need the `-overlay` option; the test wrapper uses `-vet=off` for virtual files and is not a substitute for separate static analysis.

## Security and migration limits

`@Service`, `@Repository`, `@Consumer`, `@Producer`, `@Authorize`, `@Transactional`, `@Audit`, `@Retry` and most other cataloged directives **do not enforce behavior yet**. Do not use them as security, durability, or transaction guarantees. Keep authorization, domain decisions, GORM predicates, credential/OAuth policy, message acknowledgments and event publication explicit until tested processors are separately implemented.

The complete remaining work and acceptance criteria are in the Phase 2 migration ledger.

# Annotation Contract and Generation Guide

**Generic only.** Module-specific annotation decisions belong to the module's `docs/engineering/`. Read [backend.md](backend.md) first.

## 1. What an annotation means

Go has no runtime Java-style annotations. The syntax `// @Name ...` is a source comment recognized by a **specific** GrooveHub scanner/checker/generator. Its presence alone does not create runtime functionality. Every directive has four dimensions:

| Dimension | Meaning |
|---|---|
| **Target** | Type, interface, method, function or field supported by the processor |
| **Mode** | `GEN`: generated typed Go; `CHECK`: compile-time validation; `REG`: registry metadata; `TAG`: validated struct tags; `WRAP`: deliberate wrapper |
| **Status** | Existing, Pilot, Optional, Proposed, Discouraged or historical catalog |
| **Runtime obligation** | Actual constructor, auth check, transaction, consumer lifecycle, or adapter still required |

A directive documented as Proposed must **not** be added to code expecting enforcement. `Pilot` means there is a limited checked/generated implementation, not a broad runtime framework. Inspect the active processor, tests, and generated overlay before relying on it.

## 2. Source layout

```text
libs/core/annotation/
├── annotation.go                    # public lookup / parsing facade
├── inventory.go                     # category definitions aggregation
├── contract/                        # annotation metadata and targets
├── entity/                           # entity / data metadata
├── dto/                              # presence-aware Patch[T]
├── mapping/                          # typed mapper directives
├── validation/                       # validation/tag rules
├── persistence/                      # repo/search declarations
├── http/                             # routes / methods / client request lines
├── external/                         # typed clients and external providers
├── configuration/                    # env/config contracts
├── event/                            # domain event metadata
├── messaging/                        # future producer/consumer registration
├── composition/                      # service assembly metadata
├── strategy/ and chain/              # checked operation annotations
├── testing/                          # test registration / mocks
├── security/, middleware/, server/   # enforcement catalog
├── observability/, utilities/        # observability / helpers catalog
├── patterns/                         # GoF pattern documentation/registry
└── generation/
    ├── parser/ and check/
    ├── entities/ and gormtable/
    ├── routes/ and clients/
    ├── mappers/ and events/
    ├── search/ and patchdto/
    └── overlay/ and expand/
```

Category-specific `annotations.go` describes a directive. The processor under `generation/` determines executable semantics; these are **not interchangeable**.

## 3. Safe/common usage matrix

| Annotation | Target, purpose | Actual or expected boundary |
|---|---|---|
| `@Entity <schema>.<table>` | Persisted struct; GORM `TableName` | Generates table name **only**; no mutations, constructors or validations |
| `@Builder`, `@Getter`, `@Setter` | Optional typed construction/access | Only on supported targets; do not auto-generate mutable sensitive fields |
| `@Data` | Passive data marker | Checked metadata, not magic value validation |
| `@PatchRequestDto` | Presence-aware PATCH DTO | Generates supported JSON marshal/unmarshal semantics; validate absent/null/value |
| `@Mapper`, `@Mapping` | Typed representation converter | Generate mapping on supported interface; explicit sensitive/null behavior |
| `@Repository entity=<Type>` | Repository ownership | Checked metadata; **does not** create a DB or transaction |
| `@Searchable ...` | Repository search implementation | Generated only against explicit allowlisted config; no inferred SQL |
| `@Filterable`, `@Sortable`, `@Cursor` | Search allowlist parity | Checked fields; policy stays in resource validator |
| `@Service` | Composed service | Checked ownership; explicit bootstrap constructors |
| `@Strategy` | Typed operation strategy | Checked; actual behavior implemented and wired explicitly |
| `@ChainOfResponsibility` | Ordered operation chain | Checked; invoke real typed runner yourself |
| `@Transactional` | Write operation intent | Checked; **explicit** UnitOfWork still required |
| `@Client` | Public typed client interface | Processor-generated proxy or request adapter |
| `@Proxy service.Operation` | Delegation in typed client | Calls supplied service implementation; no implicit network |
| `@RequestLine METHOD /path` | In-process typed HTTP client | Generated against registered methods and REST contracts |
| `@Routes` + `@GET/POST/PUT/PATCH/DELETE` | HTTP handler routing | Generated registration, not permission checks |
| `@Event topic=... type=... source=...` | Domain event contract | Generated/checked metadata, not delivery guarantees |
| `@ExternalService` | Third-party capability interface | Declares external boundary; bootstrap registers concrete adapter |
| `@Provider id=NAME category=kind` | Concrete provider type | **Checked** uppercase identifier, lowercase category, unique IDs; no auto-init |
| `@Configuration`, `@Env`, `@Secret` | Typed runtime configuration | Use only implemented features; secret redaction and validation still required |
| `@Validator`, `@Valid`, `@NonNull`, `@Size` | Input validation/tag metadata | Match real validation processor and field tags, not comments alone |
| `@Testing`, `@Mock` | Approved test-suites / interface mocks | Generate in test-only ownership, not production code |

Consult the live catalog to distinguish a registered *name* from the exact parser/generator behavior. Other catalog annotations (audit, cron, tracing, retry, auth scopes, middleware, caching, etc.) **must not be treated as implemented** until a tested processor enforces them.

## 4. Examples with placeholder types

```go
// @Entity app.records
type Record struct {
    ID string
}

// @Repository entity=Record
type RecordRepository interface {
    Create(ctx context.Context, record Record) error
}

// @Service
type RecordService struct {
    // Explicitly constructed operation dependencies.
}

// @Provider id=VENDOR_A category=storage
type vendorStrategy struct {
    // Concrete external service client supplied by bootstrap.
}
```

These snippets illustrate *annotation placement*; they are not complete compilable modules. Use real schema names and only annotations supported by the currently installed processor.

For HTTP:

```go
// @Routes /resources
type ResourceController interface {
    // @GET /{id}
    Get(writer http.ResponseWriter, request *http.Request)
}
```

An access tier inside a route does **not** create identity verification; authorization must be wired independently.

## 5. Generators, overlay and ownership

1. Maintained Go definitions are the **source of truth**.
2. The annotation checker scans maintained source, rejects unknown directives, malformed targets and metadata-only directives that have no processor.
3. A generator verifies method signatures and emits typed Go beneath the owning service/library `.generated/` directory.
4. The Go build/test invokes the approved overlay JSON to see generated files without committing them beside maintained source.
5. Never edit an emitted file by hand. Fix the source annotation or generator and regenerate.
6. The generated output must preserve API signature, public behavior, error mapping, source ownership and import rules.
7. Generated code must be reproducible across clean local checkout and CI; do not rely on hidden local runtime environment.
8. Do not add generated code to the module's public boundary unless the contract explicitly requires it.

## 6. Provider annotation rules

The concrete external adapter owns `// @Provider id=... category=...`. Provider IDs are static, uppercase `A-Z`/`0-9`/`_`, unique across checked services. Category is lowercase and stable. The checker validates metadata; explicit bootstrap still constructs clients, validates credentials, chooses adapters and registers capabilities. Runtime registration must report duplicate/conflicting providers; a source check alone does not establish a running catalog.

## 7. Annotation design checklist

Before adding a directive, answer:

- Does an actual processor recognize this **exact target and argument syntax**?
- Is it generating useful **typed code**, or would plain Go be clearer?
- Who owns its generated file? Will a private overlay compile without consumer-side internal imports?
- Is the operation's auth, transaction, retry and lifecycle visible without reading generated implementation?
- Can malformed metadata or collision fail the build predictably?
- Are tests present for valid syntax, invalid syntax, target mismatch, generator conflicts and runtime behavior?
- Are the master annotation guide and the owning library's engineering docs updated?

**Never use an annotation as a substitute for a business invariant or a security control.**

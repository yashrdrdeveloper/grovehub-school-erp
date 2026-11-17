# Backend Annotation and Code Generation Standard

Scope: Go source annotations and deterministic generated code. An annotation is a **build-time declaration**, not a replacement for architectural boundaries or runtime validation.

## Design constraints

1. Prefer idiomatic, handwritten Go if generation does not remove meaningful repetition.
2. Annotations declare metadata; generators produce code without hidden network or database side effects.
3. Generated code must be deterministic, reproducible and versioned by a documented toolchain.
4. Do not generate business decisions, authorization rules, orchestration or provider-specific policy.
5. Do not turn annotations into an implicit service container.
6. Keep visible source minimal, with generated output in module-approved `.generated/` or other explicitly documented locations where build/test tooling can still compile it.
7. Go ignores dot-prefixed directories in ordinary package discovery; if output lives there, generation must also arrange a supported compile/import workflow. Never promise the output will compile automatically.

## Annotation families (proposals require implementation verification)

| Annotation | Intended responsibility | Boundary |
|---|---|---|
| `@Entity` | Map a declared entity to schema/table metadata | No business validation |
| `@Table` | Describe table/column mapping where needed | No migration execution at runtime |
| `@Repository` | Generate or register repetitive persistence glue | No hidden transactions |
| `@Getter`, `@Setter`, `@Data` | Explicit data-access boilerplate only if justified | Do not break invariants or expose private state |
| `@Builder` | Construct values with optional fields and declared validation | Not a default for every struct |
| `@Mappers` | Deterministic conversion between known shapes | No lossy implicit conversion |
| `@Validators`, `@Valid`, `@NonNull` | Generate structural checks | Application semantics stay explicit |
| `@PatchRequestDto` | Track presence/null/value for patchable fields | Never conflate omission and explicit null |
| `@Searchable`, `@Filterable`, `@Sortable`, `@Cursor` | Declare trusted searchable field metadata | No user-provided raw SQL |
| `@Routes` | Generate transport registration metadata when adopted | No authorization policy inference |
| `@Client` | Repetitive typed client plumbing if required | No secret management or retry policy |
| `@Service` | Not recommended for generation in this architecture | Application services remain explicit |

These names form a **design vocabulary**, not a claim that all processors currently exist. Keep the implementation status in the generator's module engineering docs and do not annotate source with unsupported features.

## Generation process

```text
Handwritten source + annotation declarations
             |
    Parse / validate declarations
             |
    Generate deterministic Go
             |
  gofmt -> compile -> contract tests
```

A generator must reject invalid declarations with file/line diagnostics, detect collisions, preserve explicit field and null semantics, and refuse to overwrite handwritten files. Generated output should contain a provenance header. CI must detect drift between source and generated artifacts.

## Architectural placement

- Generic annotation parser/generator packages belong in a shared tooling/library boundary.
- Each module owns its entity and repository declarations.
- Module engineering docs describe the annotations actually enabled, generation command, compile visibility and regeneration procedure.
- Product docs describe resulting API behavior, never annotations.

## Non-goals

Generating complete services, implicit dependency injection, magic transaction management, silent entity mutation, undocumented SQL generation, cross-provider business adapters or treating generation as a substitute for tests.

## Checklist

Before adding an annotation: verify generator support, name/ownership, determinism, Go visibility, error diagnostics, concurrency safety, migration compatibility, opt-out path and documented regeneration.

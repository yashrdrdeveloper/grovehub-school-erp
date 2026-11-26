# Core Annotation Engine — Implementation Notes

The generic semantics are defined in [annotation.md](../../../../../docs/engineering/backend/annotation.md). This file is for the **Core library owner** to describe the actual engine boundary.

## Source packages

`annotation/annotation.go`, `inventory.go`, `contract/` and category-owned `annotations.go` declare and expose metadata. `annotation/generation/parser` scans source comments and targets; `generation/check` rejects unknown or unsupported use cases; specialized generators handle entities, routes, typed clients, mappers, events, search and PATCH DTOs; `generation/overlay` exposes emitted files to the Go toolchain.

## Build contract

```text
maintained Go + // @Directive
  -> parser (name/arguments/target)
  -> catalog resolution and checker (allowed status and target)
  -> specialized generator (typed code where implemented)
  -> ignored <owner>/.generated/
  -> overlay.json
  -> go build / go test using overlay
```

The checker is a safety gate, not a runtime security provider. A `CHECK` directive must never be described as if it creates auth, transactions, retry policies, scheduling or dependency injection.

## Concrete guardrails

- `@Entity`: table name only for supported structs; not automatic setters/builders.
- `@Repository`: ownership metadata; GORM and transactions are explicit.
- `@Searchable`: explicit typed config; SQL is not inferred from untrusted fields.
- `@Routes`: handler registration, not authentication/authorization.
- `@Client/@Proxy`: typed delegation; does not automatically run public HTTP middleware.
- `@Provider`: static ID/category and duplicate validation, not runtime adapter registration.
- `@Transactional`: intent metadata; tests must verify explicit UnitOfWork.
- `@Testing/@Mock`: tests generated under approved test ownership only.

Any proposed annotation must obtain a tested parser/generator and be moved to an executable status **before** callers rely on its effects. Generated code must remain reproducible, overlay-backed and out of maintained source.

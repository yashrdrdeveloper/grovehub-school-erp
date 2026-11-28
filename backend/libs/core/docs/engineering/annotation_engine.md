# Core Annotation Engine — Implementation Notes

The generic semantics are defined in [annotation.md](../../../../../docs/engineering/backend/annotation.md). This file is for the **Core library owner** to describe the actual engine boundary.

## Source packages

`annotation/annotation.go` exposes metadata from `internal/src/registry/`, `internal/src/contracts/` and `internal/src/catalog/`. `internal/src/parser/` scans source comments and targets; `cmd check` rejects unknown or unsupported use cases; specialized generators handle entities, routes, typed clients, mappers, events, search and PATCH DTOs; `cmd overlay` exposes emitted files to the Go toolchain.

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

## Explicit roots and module identity

Application generator entrypoints take the application root. Testing and overlay accept `<application-directory> <core-directory>`; Core-only mode uses the same directory for both. The application must have its own go.mod. The shell runner at `scripts/backend/generate-annotations.sh` accepts the same two roots, invokes the canonical Core import path and preserves owner-local generated outputs. Coverage accepts explicit Core, generation script and reports paths. The empty-overlay exception is restricted to standalone Core mode.

## Phase 1 capability encapsulation

The physical annotation owner is `backend/libs/core/annotation/`. Runtime clients import the root facade; the single command under `cmd/main.go` dispatches through `internal/src/cli/` to `internal/src/generators/`. Catalog metadata lives under `internal/src/catalog/`, typed metadata contracts under `internal/src/contracts/`, and runtime operations under their responsibility-specific packages in `internal/src/`. Private unit suites belong in `internal/resources/tests/unit-test/`; public API suites remain under Core's `tests/unit-test/`. The testing generator, method inventory and coverage runner discover the capability-local tree. Legacy annotation compatibility packages have been removed; the root facade is the sole runtime import.

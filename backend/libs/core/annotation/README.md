# GrooveHub Annotation Core — Phase 1

> **Status:** Metadata and declaration contracts only. Phase 2 has not begun.
>
> **Source of truth:** [Backend annotation master plan](../../../../docs/engineering/backend/ANNOTATION_FRAMEWORK_MASTER_PLAN.md) and [Phase 1 annotation reference](../../../../docs/engineering/backend/ANNOTATION_PHASE1_REFERENCE.md).

This directory owns `@...` annotation definitions. It **does not** replace GORM, Watermill, Go's HTTP server, any application service, or any business logic.

## Inventory

- **226 use cases** across 16 backend responsibility categories.
- **232 unique directive names**, including 21 GoF-only annotations.
- **All 23 GoF design patterns** classified with automation decisions; `@Builder` and `@Adapter` already appear in the primary inventory.
- `annotations_*.go` files own a category, `inventory.go` assembles it, and `patterns.go` records pattern decisions.
- `annotation.go` exposes typed definitions, safe catalog lookup, source comment parsing, and basic syntax validation.
- `generation/parser/parser.go` scans ordinary Go source using the Go AST; it never changes or executes source.

## Phase 1 API

```go
definition, ok := annotation.Lookup("Mapper")
entries := annotation.All()
contexts := annotation.Uses()
patterns := annotation.Patterns()

directive, found, err := annotation.ParseComment("// @Entity foundation.person_contacts")
occurrences, err := parser.ScanFile("path/to/contact.go")
```

`ParseComment` rejects unknown annotation names and validates the existing `@Entity schema.table` syntax. Other directive arguments are currently opaque strings: they **must not be treated as validated** until Phase 2 provides the processor and schema. `ScanFile` reports structural target and source position for each occurrence.

## One recognized syntax, no implicit execution

```go
// @Entity foundation.person_contacts   // existing generator: Foundation only
// @Mapper                               // metadata only
// @NonNull identity.InvalidRegistration // metadata only; does NOT validate nil
// @Repository                           // metadata only
// @Service                              // metadata only; do NOT apply in services in Phase 1
// @Producer                             // metadata only; does NOT publish
```

**Phase 1 does not apply any new directives to any service.** For example, the comment `// @Authorize` does not enforce an access-control policy; neither does `// @Transactional` begin a transaction. This distinction is a security requirement, not a missing optional convenience.

The existing `@Entity` and upstream Gombok annotations retain their prior generator behavior; they were already present before this phase. Their existing `generation/expand`, `generation/gormtable`, `generation/overlay` tools are not rewritten here.

## Contracts and migration

Phase 2 will create actual typed processors, parameter schemas and explicit runtime wrappers, then migrate each backend layer **after behavior tests**, preserving all explicit service rules. Do not substitute metadata annotations for working validation, authorization, durable message delivery, secret management, or transaction handling.

Code-generation outputs stay under the owning module's ignored `.generated/` directory and compile using the established Go overlay strategy. The annotation registry is **not** a runtime reflection container or a service locator.

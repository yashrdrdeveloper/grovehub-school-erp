# Core Annotation Framework

**Owner:** `backend/libs/core/annotation/` • **Purpose:** compile-time metadata, validation and typed Go generation.

## Public API and ownership

Consumers use `annotation.go`, importing `github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation`. Public operations include catalog lookup, directive parsing and validation, source scanning, typed PATCH fields, entity builders, startup configuration and DTO validation.

```text
annotation/
├── annotation.go
├── generation/                      # stable command-line entry points only
└── internal/
    ├── src/
    │   ├── contracts/
    │   ├── catalog/
    │   ├── impl/
    │   ├── helpers/
    │   ├── parser/
    │   └── generators/
    └── resources/tests/unit-test/   # private implementation tests
```

Core-level public API tests live in `backend/libs/core/tests/unit-test/`. Capability-local tests mirror private source responsibilities and may import their own internal packages. Production code must never import test resources. Both test trees are discovered by the annotation testing and coverage tools. The generators retain their original CLI paths for scripts and CI.

Legacy category and runtime import paths are compatibility bridges during the consumer migration; do not use them in new code. The PATCH generator accepts both legacy and facade imports while migration is in progress.

## Maintained source of truth

The metadata catalog is declared under `internal/src/catalog/` and composed by the registry under `internal/src/impl/registry/`. Actual generation logic lives under `internal/src/generators/`. Go comments do not enforce runtime transactions, authorization, retries or network registration.

Generated code remains ignored beneath the owning service or library and is injected via the approved Go overlay. Modify source declarations and generators, not generated artifacts.

A checked `@Provider id=... category=...` validates static metadata; registration remains explicit. A checked `@Transactional` requires a real UnitOfWork. `@Entity` generates TableName, not business constructors or setters.

See [Core annotation engine](../docs/engineering/annotation_engine.md) and [Core capabilities index](../docs/engineering/README.md).

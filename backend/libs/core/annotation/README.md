# Core Annotation Framework

**Owner:** `libs/backend/core/annotation/` • **Purpose:** compile-time metadata, validation and typed Go generation.

- **Generic usage and annotation matrix:** [docs/engineering/backend/annotation.md](../../../../docs/engineering/backend/annotation.md)
- **Backend architecture:** [docs/engineering/backend/backend.md](../../../../docs/engineering/backend/backend.md)
- **Library implementation details:** [Core annotation engine](../docs/engineering/annotation_engine.md)
- **Core capabilities index:** [Core engineering](../docs/engineering/README.md)

## Maintained source of truth

The catalog is declared by category-owned `annotations.go` files and accessed through `annotation.go`/`inventory.go`; the processor code under `generation/` determines whether a directive is **implemented**, **checked** or still a proposal. Go comments do not enforce runtime transactions, authorization, retries or network registration.

Code generation writes ignored `.generated/` output beneath the **owning service or library** and injects it via the approved Go overlay. Modify the maintained Go declarations and generators, not generated artifacts. Tests must verify the exact parser target, metadata constraints, generated signature and runtime behavior.

A checked `@Provider id=... category=...` validates IDs/categories and duplicate declarations; adapter construction and registration remain explicit in bootstrap. A checked `@Transactional` requires a real UnitOfWork. The `@Entity` generator provides TableName only; no automatic business constructors/setters.

For adding or changing a directive, start with the [generic annotation contract](../../../../docs/engineering/backend/annotation.md), modify the category registry and executable processor together, add negative tests, regenerate the overlay, and run the full backend build.

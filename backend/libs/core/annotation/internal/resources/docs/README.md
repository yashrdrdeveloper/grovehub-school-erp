# Core Annotation Framework

**Owner:** `backend/libs/core/annotation/` • **Purpose:** compile-time metadata, validation and typed Go generation.

## Public API and ownership

Consumers use `annotation.go`, importing `github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation`. Public operations include catalog lookup, directive parsing and validation, source scanning, typed PATCH fields, entity builders, startup configuration and DTO validation.

```text
annotation/
├── annotation.go
├── cmd/main.go                      # single executable command
└── internal/
    ├── src/
    │   ├── contracts/
    │   ├── catalog/
    │   ├── registry/
    │   ├── cli/
    │   ├── configuration/, dto/, entity/, validation/
    │   ├── helpers/
    │   ├── parser/
    │   └── generators/
    └── resources/tests/unit-test/   # private implementation tests
```

Core-level public API tests live in `backend/libs/core/tests/unit-test/`. Capability-local tests mirror private source responsibilities and may import their own internal packages. Production code must never import test resources. Both test trees are discovered by the annotation testing and coverage tools. All fourteen processors are dispatched by one command registry. Run `go run ./annotation/cmd --help` from the Core module, or `go run ./annotation/cmd routes <backend-directory>` for a specific processor. Coverage and inventory are subcommands too.

The root facade is the only runtime import. Legacy category and runtime subpackages have been removed. The PATCH generator requires the facade import.

## Maintained source of truth

The metadata catalog is declared under `internal/src/catalog/` and composed by the registry under `internal/src/registry/`. Actual generation logic lives under `internal/src/generators/`. Go comments do not enforce runtime transactions, authorization, retries or network registration.

Generated code remains ignored beneath the owning service or library and is injected via the approved Go overlay. Modify source declarations and generators, not generated artifacts.

A checked `@Provider id=... category=...` validates static metadata; registration remains explicit. A checked `@Transactional` requires a real UnitOfWork. `@Entity` generates TableName, not business constructors or setters.

See [Core annotation engine](../../../../internal/resources/docs/engineering/annotation_engine.md) and [Core capabilities index](../../../../internal/resources/docs/engineering/README.md).

## Unit-test source mirroring

Every maintained Go file under `internal/src/` must have a corresponding `*_test.go` file at exactly the same relative path under `internal/resources/tests/unit-test/`. Each test package has one Ginkgo suite, and test-only fixtures belong in `internal/resources/tests/helpers/` or a test-local fixture file. Tests of private functions exercise their observable behavior through the package API; do not export implementation details solely for tests. Metadata-only files are checked for their declared values and registry integration. Mirrored file presence is a structural requirement, not proof of statement coverage.

CI enforces the mapping and reports statement coverage across all internal source packages. From the Core module:

```bash
go test -count=1 ./annotation/internal/resources/tests/unit-test/...
go test -covermode=atomic -coverpkg=./annotation/internal/src/... -coverprofile=/tmp/annotation-coverage.out ./annotation/internal/resources/tests/unit-test/... ./tests/unit-test/annotation/...
go tool cover -func=/tmp/annotation-coverage.out
```

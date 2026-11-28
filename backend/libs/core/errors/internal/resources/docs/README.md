# Core Errors

Consumers import only `grovehub/libs/core/errors` (package
`coreerrors`). Error codes, public messages and domain policy are owned by
consumers, not this library.

```text
errors/
├── errors.go
└── internal/
    ├── src/
    │   ├── contracts/
    │   │   ├── code.go
    │   │   ├── definition.go
    │   │   └── error.go
    │   └── operations/
    │       ├── definition/
    │       │   ├── definition.go
    │       │   └── validators/status.go
    │       ├── create/create.go
    │       └── wrap/wrap.go
    └── resources/tests/unit-test/  # exact mirror of src/
```

The facade exposes `Code`, `Definition`, `Error`, `NewDefinition`, `NewError`,
`NewErrorWithMetadata`, `Wrap`, and `WrapWithMetadata`. The definition operation
validates optional HTTP status declarations (400–599, default 500). Creation and
wrapping share one private contract constructor. No interface or impl folder is
needed for these deterministic operations. The previous public contracts and
helpers packages are removed.

`Error()` returns only the declared message; `Unwrap()` preserves the cause
for standard Go errors.Is/errors.As. Metadata is retained as supplied; consumers
must handle sensitive fields and transport serialization. Definitions are
value types; the zero value retains the existing HTTP 500 fallback.

## Testing

Every internal source file has an exact relative-path Ginkgo/Gomega test file.
Public behavior remains tested in Core's `tests/unit-test/errors/`. No mocks
or infrastructure are required.

```bash
cd backend/libs/core
go test ./errors/... ./tests/unit-test/errors/...
go test -coverpkg=./errors/... -coverprofile=/tmp/errors-coverage.out ./errors/internal/resources/tests/unit-test/... ./tests/unit-test/errors/...
go tool cover -func=/tmp/errors-coverage.out
```

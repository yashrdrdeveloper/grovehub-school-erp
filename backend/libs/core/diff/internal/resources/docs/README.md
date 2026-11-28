# Core Diff

The public consumer import remains `grovehub/libs/core/diff`.
This capability compares JSON-compatible values. Business identity, event publication,
field redaction and audit policy stay with consumers.

```text
diff/
├── diff.go
└── internal/
    ├── src/
    │   ├── contracts/change.go
    │   └── operations/compare/
    │       ├── compare.go
    │       └── helpers/
    │           ├── json.go
    │           ├── walk.go
    │           ├── object.go
    │           └── array.go
    └── resources/tests/unit-test/  # exact relative-path mirror of src/
```

The `Compare[T]` facade delegates to a single pure operation, which uses JSON
normalization and recursive object/array comparison. No interface or impl folder
is required because there are no interchangeable implementations or external
dependencies. `Change` is a public alias of the private data contract.

## Behavior

- JSON field names and array indices form paths; root scalar changes use `$`.
- Object keys are sorted; array positions are visited in increasing order.
- Numbers retain `json.Number` rather than being converted to floating point.
- Added and removed values use null; missing and explicit null object fields
  retain their existing equivalence. Empty objects/arrays compared with null
  produce no leaf changes.
- Different JSON types are reported as replacements at their containing path.
- Marshal failures identify the old/new side and preserve the original error.
- Comparison neither mutates inputs nor opens connections.

Public tests live in Core's `tests/unit-test/diff/`. Private Ginkgo/Gomega suites
mirror every internal source file; fixtures use only in-memory values.

```bash
cd backend/libs/core
go test ./diff/... ./tests/unit-test/diff/...
go test -coverpkg=./diff/... -coverprofile=/tmp/diff-coverage.out ./diff/internal/resources/tests/unit-test/... ./tests/unit-test/diff/...
go tool cover -func=/tmp/diff-coverage.out
```

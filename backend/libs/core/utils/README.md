# Core Utils

Consumers import only `github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/utils`.
The former public core/helpers package is removed.

```text
utils/
├── utils.go
├── README.md
└── internal/
    ├── src/operations/strings/strings.go
    └── resources/tests/unit-test/operations/strings/
        ├── strings_test.go
        └── suite_test.go
```

OptionalString maps exactly the empty string to nil and preserves every other
value. NullableString maps nil to the empty string and dereferences any other
pointer unchanged. Neither operation trims, normalizes or performs I/O.

Only cohesive, generic conversions belong here; capability-local helpers
remain private to their owning capability. No config, constants, contracts,
helpers or impl directories are necessary for these two pure operations.

```bash
cd backend/libs/core
GOWORK=off go test ./utils/... ./tests/unit-test/utils/...
GOWORK=off go test -covermode=atomic -coverpkg=./utils/... -coverprofile=/tmp/utils-coverage.out ./utils/internal/resources/tests/unit-test/... ./tests/unit-test/utils/...
go tool cover -func=/tmp/utils-coverage.out
```

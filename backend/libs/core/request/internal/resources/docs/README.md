# Core Request

Consumers import only `github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/request`. The root facade
exposes Request, Contract, Sections, Empty, Parts and Inspect.

```text
request/
├── request.go
├── README.md
└── internal/
    ├── src/
    │   ├── contracts/request.go
    │   └── operations/inspect/inspect.go
    └── resources/tests/unit-test/
        ├── contracts/
        └── operations/inspect/
```

Concrete endpoint DTOs embed Request and declare only Params, Query and
RequestBody as needed. Inspect validates the request shape and discovers
those fields. Pointer inputs preserve writable reflection references for
HTTP binding. Empty and Parts remain available for the older typed REST
pipeline. The package does not decode HTTP input or apply domain validation.

Public type aliases preserve the source-level API; reflection reports the
internal contracts package as the declaring package. Existing regression
tests are retained. No additional configuration or helper packages are needed.

```bash
cd backend/libs/core
GOWORK=off go test ./request/... ./tests/unit-test/request/...
GOWORK=off go test -covermode=atomic -coverpkg=./request/... -coverprofile=/tmp/request-coverage.out ./request/internal/resources/tests/unit-test/... ./tests/unit-test/request/...
go tool cover -func=/tmp/request-coverage.out
```

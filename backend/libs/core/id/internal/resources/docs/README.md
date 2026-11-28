# Core ID

Consumers import only `grovehub/libs/core/id`. The root id.go facade
exposes NewUUIDv7; generation stays under internal/src/operations/uuid.
The former public id/helpers package is removed.

```text
id/
├── id.go
├── README.md
└── internal/
    ├── src/operations/uuid/uuid.go
    └── resources/tests/unit-test/operations/uuid/
        ├── uuid_test.go
        └── suite_test.go
```

The operation preserves google/uuid v7 generation, canonical string formatting
and original error propagation. No extra interface, configuration, constants
or helper folders are required. Tests exercise successful generation and
entropy failure. The public regression suite remains under tests/unit-test/id.

```bash
cd backend/libs/core
GOWORK=off go test ./id/... ./tests/unit-test/id/...
GOWORK=off go test -covermode=atomic -coverpkg=./id/... -coverprofile=/tmp/id-coverage.out ./id/internal/resources/tests/unit-test/... ./tests/unit-test/id/...
go tool cover -func=/tmp/id-coverage.out
```

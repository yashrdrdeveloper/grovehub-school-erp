# Core Header Context

Consumers import only `grovehub/libs/core/headercontext`.
The public facade preserves Values, ErrInvalidHeader, WithValues, FromContext,
IdempotencyKey, CorrelationID, RequestID, WithoutIdempotencyKey and Middleware.

```text
headercontext/
├── headercontext.go
├── README.md
└── internal/
    ├── src/
    │   ├── constants/metadata.go
    │   ├── contracts/
    │   │   ├── errors.go
    │   │   └── values.go
    │   ├── helpers/
    │   │   ├── header.go
    │   │   └── validation.go
    │   └── operations/
    │       ├── metadata/metadata.go
    │       ├── middleware/middleware.go
    │       └── read/read.go
    └── resources/tests/unit-test/  # exact mirror of src/
```

Only idempotency, correlation and request ID headers are captured. Missing
headers are optional; duplicates and blank or invalid values are denied
with HTTP 400 and a generic response. Metadata IDs retain the existing
128-byte ASCII allowlist; idempotency validation remains owned by core
idempotency. Trusted in-process metadata may propagate event correlation,
and downstream commands can clear the inherited idempotency key without
changing the parent context. Arbitrary inbound trace headers are not trusted.

No config folder is necessary because this capability has no runtime configuration.

```bash
cd backend/libs/core
GOWORK=off go test ./headercontext/... ./tests/unit-test/headercontext/...
GOWORK=off go test -covermode=atomic -coverpkg=./headercontext/... -coverprofile=/tmp/headercontext-coverage.out ./headercontext/internal/resources/tests/unit-test/... ./tests/unit-test/headercontext/...
go tool cover -func=/tmp/headercontext-coverage.out
```

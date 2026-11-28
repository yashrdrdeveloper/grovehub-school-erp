# Core Idempotency

Consumers import only `github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/idempotency`. The public
facade exposes MaxKeyLength, ErrInvalidKey, ErrFingerprint, ValidateKey and
Fingerprint. The former public contracts, helpers and validators packages
are removed.

```text
idempotency/
├── idempotency.go
├── README.md
└── internal/
    ├── src/
    │   ├── constants/key.go
    │   ├── contracts/errors.go
    │   └── operations/
    │       ├── key/key.go
    │       └── fingerprint/fingerprint.go
    └── resources/tests/unit-test/  # exact mirrors
```

Key validation preserves the existing 128-byte limit, whitespace trimming
and ASCII allowlist. Fingerprints preserve the lowercase SHA-256 digest of
the original JSON encoding, including the existing error wrapping semantics.
No storage, replay policy, business logic, runtime configuration or artificial
interfaces are introduced. Persistence and request replay remain consumer
responsibilities.

```bash
cd backend/libs/core
GOWORK=off go test ./idempotency/... ./tests/unit-test/idempotency/...
GOWORK=off go test -covermode=atomic -coverpkg=./idempotency/... -coverprofile=/tmp/idempotency-coverage.out ./idempotency/internal/resources/tests/unit-test/... ./tests/unit-test/idempotency/...
go tool cover -func=/tmp/idempotency-coverage.out
```

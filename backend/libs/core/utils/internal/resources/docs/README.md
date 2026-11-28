# Core Utils

Consumers import only `github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/utils`.
The former core/helpers and core/normalize packages are removed.

```text
utils/
├── utils.go
├── README.md
└── internal/
    ├── src/
    │   ├── constants/normalization.go
    │   ├── contracts/normalization.go
    │   └── operations/
    │       ├── strings/strings.go
    │       └── normalize/
    │           ├── normalize.go
    │           └── time.go
    └── resources/tests/unit-test/  # exact source mirrors
```

OptionalString and NullableString preserve their existing conversion behavior.
Normalization is exposed through NormalizeEmail, NormalizePhone,
NormalizePhoneCharacters, NormalizeUUIDv7 and NormalizeRFC3339UTC. Existing
validation rules, canonical representations and error identities are preserved.
Only cohesive, domain-independent utilities belong here; capability-local
helpers remain private.

```bash
cd backend/libs/core
GOWORK=off go test ./utils/... ./tests/unit-test/utils/...
GOWORK=off go test -covermode=atomic -coverpkg=./utils/... -coverprofile=/tmp/utils-coverage.out ./utils/internal/resources/tests/unit-test/... ./tests/unit-test/utils/...
go tool cover -func=/tmp/utils-coverage.out
```

# Core Logging

Import only `github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/logging`. The single public facade
exposes Levels, NewLevels, DefaultLevels, ParseLevels and NewLevelFilterHandler.

```text
logging/
├── logging.go
├── README.md
└── internal/
    ├── src/
    │   ├── config/levels.go
    │   ├── constants/levels.go
    │   ├── contracts/levels.go
    │   ├── helpers/levels.go
    │   └── operations/filter/filter.go
    └── resources/tests/unit-test/  # exact source-file mirrors
```

Logging uses explicit level membership, not a minimum threshold. The default
set is INFO, WARN and ERROR. Configuration accepts case-insensitive comma-separated
standard levels, trims whitespace, ignores empty entries and deduplicates;
unknown names and delimiter-only configuration return errors. Derived handlers
retain the filter and delegate attributes, groups and writes to the underlying
slog handler. NewLevels also exposes the former contracts constructor at the
root so callers can build custom level sets without importing internals.

No service-specific logging rules, environment reads or artificial interfaces
are introduced.

```bash
cd backend/libs/core
GOWORK=off go test ./logging/... ./tests/unit-test/logging/...
GOWORK=off go test -covermode=atomic -coverpkg=./logging/... -coverprofile=/tmp/logging-coverage.out ./logging/internal/resources/tests/unit-test/... ./tests/unit-test/logging/...
go tool cover -func=/tmp/logging-coverage.out
```

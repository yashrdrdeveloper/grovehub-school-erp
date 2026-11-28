# Core Search

Import only `grovehub/libs/core/search`. All query, cursor, sort, pagination,
validation and filter contracts are exposed by search.go.

```text
search/
├── search.go
├── README.md
└── internal/
    ├── src/
    │   ├── constants/pagination.go
    │   ├── contracts/{search,validation,errors}.go
    │   ├── helpers/query.go
    │   └── operations/
    │       ├── cursor/cursor.go
    │       ├── page/pagination.go
    │       ├── query/{query,multi_sort}.go
    │       └── validation/{cursor,multi_sort,operators,query,request}.go
    └── resources/tests/unit-test/  # mirrored source files
```

The public API retains all existing root signatures, operator tokens, JSON
field names, default/max rows, cursor encoding, and offset/keyset behavior.
Resource-owned configuration continues to allowlist SQL expressions and
operators; user values are bound as parameters. Relation set filters bind
fixed conditions before list values, matching SQL placeholder order.
Service-specific filter normalization and search policies remain in services.

The old public contracts, helpers, mappers and validators directories are
removed. Regression tests previously adjacent to production source now run
from tests/unit-test/search. Internal mirrored tests cover SQL building,
cursor encoding, validation, pagination, and query execution using sqlmock.

```bash
cd backend/libs/core
GOWORK=off go test ./search/... ./tests/unit-test/search/...
GOWORK=off go test -covermode=atomic -coverpkg=./search/... -coverprofile=/tmp/search-coverage.out ./search/internal/resources/tests/unit-test/... ./tests/unit-test/search/...
go tool cover -func=/tmp/search-coverage.out
```

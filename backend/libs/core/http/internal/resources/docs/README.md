# Core HTTP

Consumers import only `grovehub/libs/core/http` (package corehttp).
The root http.go facade exposes response contracts, correlation, endpoint
dispatch, header names, logging levels, readiness, request decoding, response
writing and tracing.

```text
http/
├── http.go
└── internal/
    ├── src/
    │   ├── constants/
    │   ├── contracts/
    │   ├── helpers/
    │   └── operations/
    │       ├── correlation/
    │       ├── endpoint/
    │       ├── readiness/
    │       ├── request/
    │       ├── response/
    │       └── tracing/
    └── resources/tests/unit-test/  # exact source-file mirrors
```

All former public subpackages have been removed. Inbound metadata header names
are owned by headercontext and re-exported through HTTP constants, preventing
the endpoint/headercontext import cycle. Runtime behavior, response JSON,
request types, trace propagation and security constraints remain unchanged.
There are no artificial interfaces or empty config folders.

Readiness calls use ReadinessHandler and ReadinessHTTP in the root facade.
The old correlation.Header constant maps to RequestID.

```bash
cd backend/libs/core
GOWORK=off go test ./http/... ./tests/unit-test/http/...
GOWORK=off go test -covermode=atomic -coverpkg=./http/internal/src/...,./http -coverprofile=/tmp/http-coverage.out ./http/internal/resources/tests/unit-test/... ./tests/unit-test/http/...
go tool cover -func=/tmp/http-coverage.out
```

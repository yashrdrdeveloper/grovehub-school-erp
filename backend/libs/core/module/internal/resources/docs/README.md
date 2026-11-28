# Core Module

Import only `github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/module`.
The public module.go facade exposes the shared Module lifecycle interface.

```text
module/
├── module.go
├── README.md
└── internal/
    ├── src/contracts/module.go
    └── resources/tests/unit-test/contracts/
        ├── module_test.go
        └── suite_test.go
```

The interface intentionally defines only RegisterHTTP(*http.ServeMux) and
Close(). Service-specific modules own bootstrapping, typed clients, resource
management and implementations. This package does not execute lifecycle
operations and does not need operations, helpers, constants, config or impl
directories.

The facade aliases the internal contract, preserving consumer method sets and
structural interface compatibility. Contract tests verify the exact method set;
public regression tests exercise an independently defined service module.
There are no executable production statements to measure for statement coverage.

```bash
cd backend/libs/core
GOWORK=off go test ./module/... ./tests/unit-test/module/...
```

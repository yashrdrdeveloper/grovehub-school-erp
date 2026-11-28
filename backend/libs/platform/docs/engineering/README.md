# Platform shared library

Module: `grovehub/libs/platform`.

The Platform-specific shared library owns security contracts and session context shared by the Platform launcher and services. It does not own business services, HTTP routing, persistence, or workload orchestration. Session is its first capability.

Import only `grovehub/libs/platform/session`. The capability's constants, contracts, and operations live under `session/internal/src`; other packages cannot import them directly. Principals must originate from an authenticated adapter, and tenant checks must originate from trusted registry policy. Authentication, cookies, and generation checks retain their existing semantics.

The module is standalone, uses only the Go standard library, and does not depend on Core or application projects. Its local module path is deliberately not a GitHub URL. Platform's project module uses a checkout-relative `replace`; no workspace is required. Remote module distribution is not configured.

From this module:

```sh
GOWORK=off go build ./...
GOWORK=off go test ./...
bash tools/coverage/check.sh
```

Tests are capability-owned and run without annotation generation or an application database. Launcher retains an annotated consumer-contract suite. The independent coverage gate rejects every uncovered production statement; Launcher coverage measures only Launcher-owned production packages.

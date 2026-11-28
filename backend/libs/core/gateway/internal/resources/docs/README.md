# Core Gateway

Consumers import only `github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway`.
The public facade exposes existing route tiers, workload claims, policy and
verifier contracts, Ed25519 credentials, restricted-route middleware and
signed HTTP transport. Concrete public types are aliases of internal types,
so existing struct literals, interfaces and method calls remain valid.

```text
gateway/
├── gateway.go
├── README.md
└── internal/
    ├── src/
    │   ├── config/config.go
    │   ├── constants/credential.go
    │   ├── contracts/contracts.go
    │   ├── helpers/validation.go
    │   └── operations/
    │       ├── credential/credential.go
    │       ├── restricted/restricted.go
    │       ├── transport/transport.go
    │       └── workload/workload.go
    └── resources/tests/unit-test/  # exact source mirrors
```

The five source categories are config (configuration construction), constants
(fixed protocol values), contracts (public types and interfaces), helpers
(pure reusable validation), and operations (runtime behavior).

The existing ghw1 Ed25519 wire format, expiry, scope, delegation, route,
origin, header filtering, redirect and fail-closed checks are preserved.
Credential verification establishes authenticity; scoped workload verification
authorizes a registered route and attaches verified claims to context.
RestrictedOnly is a route boundary, not a substitute for authorization.
Transport signs only policy-derived claims. Key configuration must be immutable
while serving traffic. No service-specific policies or private keys belong here.

```bash
cd backend/libs/core
GOWORK=off go test ./gateway/... ./tests/unit-test/gateway/...
GOWORK=off go test -covermode=atomic -coverpkg=./gateway/... -coverprofile=/tmp/gateway-coverage.out ./gateway/internal/resources/tests/unit-test/... ./tests/unit-test/gateway/...
go tool cover -func=/tmp/gateway-coverage.out
```

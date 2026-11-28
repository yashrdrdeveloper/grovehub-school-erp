# Core Telemetry

Import only `github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry`. The root facade exposes logical service
identity, runtime initialization and shutdown, provider selection, metric
recording, trace context, structured logging and configuration validation.

```text
telemetry/
├── telemetry.go
├── README.md
└── internal/
    ├── src/
    │   ├── constants/config.go
    │   ├── contracts/service.go
    │   ├── config/{environment,validation}.go
    │   ├── helpers/service.go
    │   └── operations/
    │       ├── runtime/{runtime,metrics}.go
    │       └── tracing/trace.go
    └── resources/tests/unit-test/  # mirrored source files
```

Runtime lifecycle and metric instruments share a private runtime package
because they own the same provider state. Tracing uses that runtime's provider
selection. Environment configuration and pure service-name resolution are
separate. Existing OTLP endpoints, metric names, resource attributes, service
identities, sampling defaults and shutdown behavior are preserved. Former
public contracts, helpers, mappers and validators are removed. Their supported
configuration APIs are now available directly from the facade.

```bash
cd backend/libs/core
GOWORK=off go test ./telemetry/... ./tests/unit-test/telemetry/...
```

# GroveHub Current Backend Structure

> **Status:** Mandatory GroveHub architecture standard  
> **Applies to:** the current Go modular monolith and all new backend services/modules  
> **Precedence:** This file is authoritative for current GroveHub service folder structure when an older generic backend blueprint conflicts with it.

---

## 1. Current Runtime

GroveHub currently runs as one Go platform process with a root launcher acting as the HTTP edge/composition boundary.

```text
cmd/platform
   ↓
launcher (edge/composition)
   ↓
services
```

The root launcher is allowed and required.

Service-local `launcher/` folders are not part of the current GroveHub service structure.

---

## 2. Service Structure

A service owns its construction through its root package and exposes its HTTP routes/controllers directly through explicit folders.

Current canonical shape:

```text
backend/services/<service>/
├── <service>.go
├── api/
├── config/                  # when required
├── constants/               # when required
│   └── errorcode/
├── controller/
├── routes/
├── dto/
│   └── vN/
│       ├── request/
│       └── response/
├── client/                  # when a public cross-service client exists
├── src/                     # domain/business/infrastructure as complexity requires
└── tests/
    └── unit-test/
```

Use only folders with a current responsibility.

Do not add:

```text
services/<service>/launcher/
```

unless the architecture decision is explicitly changed and this standard is updated first.

---

## 3. Identity Example

Current Identity shape:

```text
backend/services/identity/
├── identity.go
├── api/
│   └── openapi/v1/openapi.yaml
├── constants/
│   └── errorcode/
├── controller/
├── routes/
├── dto/
│   └── v1/
│       └── response/
└── tests/
    └── unit-test/
```

`identity.go` owns module construction and route registration exposure.

`controller/` owns transport controllers.

`routes/` owns HTTP route registration.

---

## 4. Root Launcher

The root launcher owns platform-wide concerns such as:

```text
HTTP server lifecycle
edge routing/fallback
cross-cutting middleware
OpenTelemetry server instrumentation
service module registration
platform bootstrap
```

It must not absorb service business logic.

The launcher is the current edge/API-gateway boundary and starts the top-level HTTP trace span.

---

## 5. Shared Technical Libraries

Cross-service technical primitives belong under:

```text
backend/libs/core/<capability>/
```

Examples:

```text
libs/core/errors
libs/core/dto
libs/core/http
libs/core/telemetry
```

Shared business/domain concepts do not belong in `libs/core` merely because more than one service might eventually use them.

Owner-local core tests live under:

```text
backend/libs/core/tests/unit-test/<mirrored-path>/
```

---

## 6. Error Codes

Service-specific semantic error definitions live with the owning service, for example:

```text
backend/services/identity/constants/errorcode/health_code.go
```

A semantic error definition owns stable:

```text
Code
Message
```

HTTP status remains a transport/controller decision.

Optional metadata is supplied at the failing request boundary when needed.

---

## 7. Trace Context

The root launcher starts/continues the edge trace.

Service route boundaries add child spans while preserving the same trace ID.

All nested calls must pass `context.Context`.

See:

```text
docs/engineering/observability/TRACE_CONTEXT_GUIDELINES.md
```

---

## 8. Unit Tests

Backend unit tests live inside their owner, not under `backend/testing`.

See:

```text
docs/engineering/testing/BACKEND_UNIT_TEST_STRUCTURE.md
```

---

## 9. Agent Rule

Before creating a backend folder, determine whether it belongs to:

```text
root platform edge/launcher
service-owned API/business code
shared technical core
owner-local tests
```

Do not copy an older generic service-local `launcher/` blueprint into GroveHub.

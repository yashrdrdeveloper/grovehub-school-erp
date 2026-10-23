# GroveHub Trace Context Guidelines

> **Status:** Mandatory GroveHub engineering standard  
> **Applies to:** launcher/edge HTTP traffic, backend services, in-process service calls, outbound HTTP/gRPC, logs, traces, error diagnostics, and OpenTelemetry exporters  
> **Precedence:** This file is authoritative for GroveHub trace propagation when an older generic guideline conflicts with it.

---

## 1. Objective

Every request entering GroveHub must be traceable through the complete backend execution path.

The expected waterfall is conceptually:

```text
client
  ↓
GroveHub edge / launcher span
  ↓
service span
  ↓
service/business/dependency span
  ↓
outbound API / database / provider span
```

A nested call must continue the existing trace instead of starting an unrelated trace.

---

## 2. Trace Identity

OpenTelemetry owns trace and span identifiers.

Use:

```text
traceId -> identifies the complete distributed request/execution
spanId  -> identifies one operation inside that trace
```

Never generate a replacement trace ID in the middle of an existing request.

Pass `context.Context` through every backend call that participates in request execution.

---

## 3. HTTP Headers

GroveHub exposes the active trace identifier through:

```text
X-Trace-ID
```

Rules:

1. The edge sets `X-Trace-ID` on backend HTTP responses when a valid trace exists.
2. GroveHub outbound HTTP calls forward `X-Trace-ID` for support/debug correlation.
3. `X-Trace-ID` is not part of the JSON response envelope.
4. Header matching is case-insensitive; Go may render the canonical casing as `X-Trace-Id`.
5. The authoritative distributed propagation format is W3C Trace Context using `traceparent`.
6. Do not attempt to reconstruct a parent span from `X-Trace-ID` alone because a trace ID does not contain the parent span identity.
7. An inbound `X-Trace-ID` without a valid W3C parent is not enough to continue a caller's span hierarchy; GroveHub creates an authoritative OpenTelemetry trace and returns its trace ID.

Therefore every traced outbound HTTP request carries both concepts:

```text
X-Trace-ID  -> human/support correlation
traceparent -> authoritative OpenTelemetry parent/child propagation
```

---

## 4. Edge Ownership

The root GroveHub launcher is the current HTTP edge/API-gateway boundary.

It owns the inbound server span and must:

```text
extract incoming W3C trace context when present
create a new trace when no valid parent exists
attach the trace context to request.Context()
return X-Trace-ID
record request timing and status
mark failed operations on the span
emit a trace-correlated request completion log
```

Edge span names use the matched route pattern rather than arbitrary path values so trace naming remains bounded as dynamic routes are introduced.

Example:

```text
edge GET /identity/health
```

---

## 5. Service Spans

Every service boundary reached during a request creates a child span from the current context.

Example:

```text
edge GET /identity/health
  └── identity GET /identity/health
```

In the modular monolith these are logical service/component spans inside the same process. Do not fake a separate deployment identity merely to create a service-map node.

When a service is later independently deployed, the same W3C context continues over the remote transport and naturally produces a cross-service trace.

---

## 6. In-Process Calls

Cross-service or internal calls must accept and forward `context.Context`.

Correct:

```go
func (service *Service) Execute(ctx context.Context, input Input) error
```

Incorrect:

```go
func (service *Service) Execute(input Input) error
```

when the operation belongs to an existing request/execution trace.

A nested operation starts a child span using the current context.

---

## 7. Outbound HTTP

Do not construct uninstrumented shared HTTP clients for GroveHub service/provider traffic.

Use the core traced HTTP transport/client so it:

```text
creates or continues an outbound client span
injects W3C traceparent
forwards X-Trace-ID from the active client span
preserves cancellation/deadlines from request.Context()
```

The traced transport must also create and forward a valid trace when a legitimate background/root outbound operation has no parent trace yet.

Do not manually copy trace IDs into query parameters or JSON payloads.

---

## 8. Logs

Request/execution logs emitted inside a valid trace must include:

```text
traceId
spanId
traceFlags
```

The local OpenTelemetry Collector maps `traceId` and `spanId` into the OpenTelemetry LogRecord trace context so Grafana can navigate between a log entry and its trace. `traceFlags` remains useful structured metadata and does not need to be parsed by every collector implementation to preserve log-to-trace correlation.

Do not invent alternate spellings such as:

```text
trace_id
traceIDValue
requestTrace
```

unless a provider protocol requires a mapped external field.

---

## 9. Failure Diagnostics

When an HTTP operation fails, GroveHub records safe diagnostic metadata on the log/span.

Default safe fields include:

```text
method
route/path
HTTP status
request content length
response byte count
duration
traceId
spanId
service/component/operation
semantic GroveHub error code when available
safe semantic error message when available
```

The span is marked as an error and receives an error/failure event where appropriate.

HTTP request handlers that emit a GroveHub semantic error must use the context-aware common error writer so the error definition is attached to the active service/edge span without parsing the response body.

Example:

```go
_ = httpresponse.WriteErrorContext(
    request.Context(),
    writer,
    http.StatusNotFound,
    errorcode.RouteNotFound,
)
```

For metadata-bearing errors use the corresponding context-aware metadata writer.

### Never capture by default

Do not automatically record raw:

```text
request bodies
response bodies
Authorization headers
Cookie / Set-Cookie
passwords
access tokens
refresh tokens
session identifiers
provider credentials
uploaded files
arbitrary query strings containing user data
```

A specific endpoint may add explicitly reviewed, redacted diagnostic metadata when necessary. Generic tracing middleware must remain payload-safe.

---

## 10. Error Logging

Avoid duplicate error logs at every layer.

The layer that owns actionable failure context should log the failure. Parent spans still show the failure through span status/events.

The edge request-completion log provides the final transport outcome and trace IDs. Semantic errors remain visible on the active failing span through `grovehub.error.code` and related safe attributes/events.

Current HTTP request-log severity policy:

```text
2xx / 3xx -> INFO
other 4xx -> WARN
404       -> ERROR
5xx       -> ERROR
```

---

## 11. Sampling

The current GroveHub foundation records all traces while usage is small.

Sampling must remain centrally configurable as traffic grows. Business code must never decide sampling independently.

A sampled-out trace must not break request processing.

---

## 12. Telemetry Failure

Observability export failure must not crash GroveHub.

If the collector/backend is unavailable:

```text
application request processing continues
trace context can still exist locally
export failure is reported safely
business behavior is not changed
```

Controlled process shutdown must flush the OpenTelemetry provider before process exit where lifecycle ownership permits it.

---

## 13. Grafana Expectations

A traced request should allow an operator to move from a failed log to a trace and inspect a waterfall such as:

```text
edge GET /identity/login
  ├── identity validate request
  ├── identity load user
  │     └── database query
  └── notification/email API
```

Each child operation must use the same `traceId` and its own `spanId`.

A failed semantic operation should expose the stable GroveHub error code on the failing span in addition to the transport status.

---

## 14. Tests

Unit tests must cover shared trace primitives and boundaries, including:

```text
child span keeps parent traceId
child span receives a different spanId
X-Trace-ID response header
outbound X-Trace-ID forwarding
outbound traceparent injection
root/background outbound call creates both headers
traceId/spanId log correlation
failure span status/event
semantic error code attached to active span
```

Tests follow the owner-local backend test structure defined in:

```text
docs/engineering/testing/BACKEND_UNIT_TEST_STRUCTURE.md
```

---

## 15. Mandatory Agent Rule

Before changing any HTTP client, service boundary, middleware, logging abstraction, worker, messaging consumer/producer, or provider call, check whether trace context must cross that boundary.

Never introduce an API-to-API path that drops `context.Context`, `traceparent`, or the GroveHub trace correlation model.

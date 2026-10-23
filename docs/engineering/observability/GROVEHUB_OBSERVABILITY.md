# GroveHub Observability

> **Status:** Mandatory GroveHub observability standard  
> **Applies to:** GroveHub backend logs, traces, metrics, service identity, dashboards, alerts, sampling, and telemetry cost controls  
> **Precedence:** This standard is authoritative for GroveHub-specific observability when older generic examples conflict with it.

---

## 1. Objective

A GroveHub request must be diagnosable from the platform edge through every real service/dependency boundary using correlated logs, traces, and metrics.

The operator flow is:

```text
Dashboard / alert
        ↓
slow or failing route/service
        ↓
trace waterfall
        ↓
failed/slow span
        ↓
correlated logs by traceId/spanId
```

---

## 2. Canonical Logical Services

Current logical services are centrally defined in `backend/libs/core/telemetry/service.go`.

```text
key             OpenTelemetry service.name
launcher-edge   LAUNCHER (EDGE)
identity        IDENTITY
```

All GroveHub logical services use:

```text
service.namespace = grovehub
```

`service.name` is a resource-level identity. Do not create ad-hoc variants such as `edge`, `launcher`, `identity-service`, or `platform` in telemetry.

Future services must be added to the canonical service registry before emitting service-owned telemetry.

---

## 3. Resource Attributes

Every trace and metric resource should contain, where available:

```text
service.name
service.namespace
grovehub.service.key
service.version
service.instance.id
deployment.environment.name
```

Logs are normalized by the local collector to the same resource identity.

---

## 4. Trace Model

One user/API operation keeps one trace ID across logical services.

```text
LAUNCHER (EDGE) server span
        ↓ same traceId
IDENTITY internal/service span
        ↓ same traceId
future DB / HTTP / provider client spans
```

Each span gets a different span ID and the correct parent span ID.

The platform edge starts or continues the inbound trace. Internal service spans must inherit `context.Context`; they must not create disconnected traces.

---

## 5. Propagation

HTTP uses:

```text
traceparent   authoritative W3C parent/child propagation
X-Trace-ID    GroveHub support/debug correlation header
```

`X-Trace-ID` is returned in backend HTTP response headers and forwarded on outbound HTTP requests.

Do not add the trace ID to the API JSON response envelope.

---

## 6. Span Boundaries

Create spans around meaningful boundaries:

```text
incoming HTTP request
service/business operation
outbound HTTP request
database operation when persistence exists
provider call when provider integration exists
message publish/consume when messaging exists
worker/job operation when workers exist
```

Do not create spans for trivial getters, mappers, or tiny helpers.

---

## 7. Errors in Traces

Failed operations should set span error status and attach safe diagnostic attributes/events.

Useful attributes include:

```text
grovehub.error.code
grovehub.operation
http.response.status_code
dependency
provider type
```

Do not capture secrets or arbitrary raw request/response bodies.

---

## 8. Logs

Application logs are structured and correlated with the active service/trace.

Common request fields include:

```text
service
serviceKey
traceId
spanId
traceFlags
request.method
request.route
request.path
response.status
durationMs
errorCode when relevant
```

Local application levels are controlled through:

```env
GROVEHUB_LOG_LEVELS=ERROR,WARN
```

Supported exact levels:

```text
DEBUG
INFO
WARN
ERROR
```

Air/build tooling output is local developer output and must not be exported to Grafana as GroveHub application logs.

---

## 9. Sensitive Data

Generic telemetry must never record:

```text
Authorization headers
cookies
passwords
access/refresh tokens
API keys
client secrets
private keys
database credentials
raw request bodies
raw response bodies
arbitrary sensitive query values
```

Prefer opaque identifiers and explicit redaction.

---

## 10. HTTP Metrics

The edge uses OpenTelemetry HTTP server metrics. The primary latency histogram is:

```text
http.server.request.duration
```

It is recorded in seconds with explicit histogram boundaries suitable for web/API latency.

The HTTP metric set is used for:

```text
request rate / RPS
request count
active requests
status/error rate
P50 latency
P95 latency
P99 latency
```

Percentiles are calculated by the metrics backend/Grafana from histogram buckets. Application code must not calculate P95/P99 request-by-request.

---

## 11. Internal Service Metrics

Meaningful service operations use:

```text
grovehub.service.operation.duration
grovehub.service.operation.requests
grovehub.service.operation.errors
grovehub.service.up
```

The operation-duration metric is a histogram in seconds.

Low-cardinality dimensions may include:

```text
grovehub.operation
grovehub.result
http.response.status_code
```

---

## 12. Runtime Metrics

The launcher process emits OpenTelemetry Go runtime metrics, including supported Go memory, goroutine, GC, and scheduling/runtime signals.

Runtime metrics belong to `LAUNCHER (EDGE)` because the current modular monolith runs in one Go process.

Do not pretend each in-process logical service has independent process CPU or heap usage.

---

## 13. Metric Cardinality

Never use these as metric labels:

```text
traceId
spanId
requestId
userId
studentId
email
raw URL with identifiers
raw error message
SQL statement
```

Use normalized route templates such as:

```text
/students/{studentId}
```

rather than raw resource paths.

The metric SDK applies a bounded cardinality limit as an additional guardrail; this does not replace correct label design.

---

## 14. Outbound HTTP

Every request-scoped outbound HTTP client must use the shared traced HTTP transport/client.

When the caller is known, use the service-aware helper so client metrics/spans are owned by the calling logical service.

Outbound HTTP continues:

```text
traceparent
X-Trace-ID
```

---

## 15. Future Dependencies

Add DB, Redis, provider, messaging, or worker telemetry only when the dependency actually exists in the product.

When introduced, instrument:

```text
latency
error/result
pool/saturation where relevant
meaningful client span
trace propagation where applicable
```

Do not create fake dependency metrics or service graph nodes before the dependency exists.

---

## 16. Sampling

Local development defaults to full trace sampling:

```env
GROVEHUB_TRACE_SAMPLE_RATIO=1.0
```

Allowed range is `0.0` through `1.0`.

Production may lower sampling after observing real traffic/cost. Parent sampling decisions are preserved across child services.

Errors/slow requests may receive higher-priority sampling later when a deliberate tail-sampling policy is introduced.

---

## 17. Metrics Export Interval

Local metric export is configurable:

```env
GROVEHUB_METRIC_EXPORT_INTERVAL=15s
```

Do not configure intervals below the repository minimum without an explicit observability decision because shorter intervals increase telemetry volume.

---

## 18. Grafana Cloud Datasources

Current GroveHub Cloud stack uses:

```text
Prometheus/Mimir -> metrics
Loki             -> logs
Tempo            -> traces
```

Dashboards must use datasource UIDs configured in Grafana rather than embedding credentials.

---

## 19. Dashboard Filters

The GroveHub overview dashboard should expose filters for:

```text
environment
service
route
HTTP method
```

Filters must support one, many, and All values using regex-safe variable expansion.

---

## 20. Initial Dashboard Signals

The initial overview should include:

```text
RPS
5xx error rate
P50 latency
P95 latency
P99 latency
active requests
latency trend
traffic by route
errors by route/status
service-operation P95
service-operation errors
service availability
Go runtime health
recent WARN/ERROR logs
trace navigation/correlation
```

---

## 21. Initial Alerts

Initial alert definitions are:

```text
high 5xx rate
high P95 HTTP latency
service telemetry unavailable
```

Thresholds are bootstrap values only. Review them after baseline production traffic exists.

Before telemetry has been deployed and observed, alert rules may be created paused to avoid false notifications.

---

## 22. Cost Controls

Keep telemetry within practical limits using:

```text
WARN/ERROR production log policy when appropriate
trace sampling
bounded metric cardinality
no payload logging
normal route templates
reasonable metric export intervals
retention/sampling policy in Grafana Cloud
```

Metrics are the primary detection signal; traces explain representative requests; logs provide detailed correlated events.

---

## 23. Local Development

Use:

```bash
make dev-start-with-logs
```

This starts:

```text
Air-backed GroveHub backend
Swagger
local OpenTelemetry Collector
```

Changes to `.env` telemetry settings require a restart because Air does not watch the root `.env` file.

---

## 24. Testing

Observability unit tests follow the owner-local Ginkgo v2 + Gomega standard.

Test the behavior that is architectural or safety-critical:

```text
canonical service identity
same trace ID across child service spans
new span ID per child span
X-Trace-ID and traceparent propagation
HTTP error severity mapping
log/trace correlation
metric/service ownership where practical
redaction/absence of sensitive fields
```

Do not write brittle tests for every human-readable log message.

---

## 25. Mandatory Agent Rules

When adding or modifying backend behavior:

1. Use the canonical service registry for telemetry identity.
2. Preserve `context.Context` across request-scoped boundaries.
3. Use `traceparent` for distributed parent/child relationships.
4. Keep `X-Trace-ID` as the GroveHub support/debug header.
5. Emit route templates, not raw identifiers, as metric dimensions.
6. Do not add unbounded metric labels.
7. Do not log secrets or arbitrary payloads.
8. Reuse common telemetry and HTTP tracing helpers.
9. Add telemetry for real dependencies only when they exist.
10. Update dashboard/query expectations when canonical metric names change.
11. Add/update owner-local Ginkgo/Gomega tests for critical telemetry behavior.
12. Report only telemetry validation that was actually performed.

---

## 26. Definition of Done

An observability change is complete when:

```text
[ ] logical service identity is correct
[ ] trace propagation remains connected
[ ] logs are trace/service correlated
[ ] metrics use bounded semantic dimensions
[ ] P95/P99-compatible histogram exists for latency behavior
[ ] sensitive data is not introduced
[ ] dashboards/alerts are updated when relevant
[ ] owner-local tests are updated
[ ] docs/AGENT routing remains current
[ ] actual validation is reported
```

---

## 27. Guiding Principle

> **Metrics tell us where the system is unhealthy, traces show where time/failure occurred, and logs explain the correlated event. All three must agree on service identity and request context.**

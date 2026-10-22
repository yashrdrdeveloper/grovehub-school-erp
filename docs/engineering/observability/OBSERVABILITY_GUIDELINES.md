# Observability Guidelines

> **Status:** Mandatory engineering standard  
> **Location:** `docs/engineering/observability/OBSERVABILITY_GUIDELINES.md`  
> **Applies to:** backend logs, metrics, traces, frontend telemetry, request/correlation IDs, provider telemetry, Mission Runtime, health/readiness, alerting, dashboards, sampling, redaction, and OpenTelemetry integration  
> **Reference stack:** OpenTelemetry + Prometheus + Loki + Tempo + Grafana  
> **Core principles:** structured signals, end-to-end correlation, bounded cardinality, safe redaction, semantic consistency, actionable diagnostics, and provider-neutral instrumentation

---

## 1. Objective

A production failure should be diagnosable without adding temporary logs and redeploying.

Observability must answer:

```text
What happened?
Where did it happen?
Which request/execution caused it?
Which organization/project/environment was involved?
Which provider/dependency was involved?
How long did it take?
Why did it fail?
What should the operator do next?
```

---

## 2. Signal Model

Perihellion uses three operational telemetry signals:

```text
Logs
Metrics
Traces
```

and two product-governance signals:

```text
Audit
Evidence
```

These are related but not interchangeable.

---

## 3. Mandatory Rules

1. OpenTelemetry is the default instrumentation model.
2. Production logs are structured.
3. `requestId`, `correlationId`, and trace identifiers have distinct semantics.
4. Context propagates across HTTP, gRPC, messaging, provider calls, workers, and durable executions where applicable.
5. Secrets are never logged.
6. Metric labels must have bounded cardinality.
7. Raw user IDs, request IDs, correlation IDs, resource IDs, URLs, SQL, and error messages are not metric labels.
8. Logs contain enough safe context to diagnose.
9. Avoid duplicate error logging at every layer.
10. Health and readiness are separate.
11. Provider operations use normalized telemetry.
12. Database pools and query behavior are observable.
13. Mission Runtime emits structured execution telemetry.
14. Frontend telemetry goes through one abstraction.
15. Vendor analytics/observability SDKs are not scattered through product code.
16. Sampling is configurable.
17. Telemetry export failure must not crash the product.
18. Telemetry queues/buffers are bounded.
19. Audit events are not debug logs.
20. Evidence is not dependent on observability retention.
21. Metric names/units are stable.
22. PII and secret redaction is centralized.
23. Telemetry configuration is environment/profile driven.
24. Correlation survives asynchronous execution boundaries.
25. Alerting should focus on user impact and SLO symptoms.
26. Observability provider replacement must not require business-code changes.
27. Health/telemetry endpoints do not expose sensitive configuration.
28. Large payloads are not logged by default.
29. Instrumentation should be tested where context propagation/redaction is critical.
30. Every service/component exposes version/build metadata.

---

## 4. Request ID

A request ID identifies one inbound request/attempt.

```text
browser request
API request
gRPC invocation
webhook request
```

Generate one if missing or untrusted.

Do not reuse one request ID across retries as though they are the same attempt.

---

## 5. Correlation ID

A correlation ID links related work.

Example:

```text
User presses "Deploy"
        ↓
HTTP request
        ↓
Mission execution
        ↓
provider calls
        ↓
events/workers
```

All may share one correlation ID while each request/span gets its own identifiers.

---

## 6. Trace ID

Trace IDs come from the tracing system.

Do not use a trace ID as durable business identity.

Traces may be sampled/expired; execution IDs and correlation IDs remain product concepts.

---

## 7. Standard Context Fields

Common safe structured fields:

```text
service
component
version
environment
requestId
correlationId
traceId
spanId
organizationId
projectId
applicationId
environmentId
executionId
capability
providerType
providerInstanceId
operation
```

Only include fields relevant to the signal.

---

## 8. Cardinality

Metrics must never use arbitrary/unbounded labels such as:

```text
user email
full URL
requestId
correlationId
executionId
resource name
raw error
SQL statement
```

Use those in logs/traces instead.

---

## 9. Log Format

Production logs are JSON/structured.

Example:

```json
{
  "level": "INFO",
  "message": "execution step completed",
  "service": "execution",
  "executionId": "opaque-id",
  "stepType": "DEPLOY",
  "durationMs": 842,
  "correlationId": "opaque-id"
}
```

---

## 10. Log Levels

Use consistently:

```text
DEBUG
INFO
WARN
ERROR
```

`ERROR` means an operation/system behavior failed, not merely that the user requested an invalid resource.

---

## 11. Duplicate Logging

Avoid:

```text
repository logs error
service logs same error
controller logs same error
middleware logs same error
```

Prefer one ownership-level error log plus trace/span error metadata.

---

## 12. Error Fields

Use stable fields:

```text
errorCode
errorCategory
retryable
operation
dependency
providerType
```

Human error text is not a metric dimension.

---

## 13. Redaction

Central redaction must cover:

```text
Authorization
Cookie
password
access token
refresh token
API key
private key
client secret
provider credential
database credential
signed URL secrets
```

Do not depend on every caller remembering redaction manually.

---

## 14. PII

Prefer opaque IDs over names/emails.

Raw PII in telemetry requires explicit need, policy, and retention consideration.

---

## 15. Frontend Logging

Frontend code uses one logging abstraction.

Avoid committed production:

```ts
console.log(...)
```

The abstraction should support:

```text
level
domain
request/correlation context
safe metadata
redaction
environment behavior
```

---

## 16. Metric Naming

Use stable semantic names.

Example:

```text
perihellion_http_server_requests_total
perihellion_http_server_duration_seconds
perihellion_provider_operations_total
perihellion_execution_steps_total
```

Follow OTel semantic conventions when available.

---

## 17. Units

Use standard units consistently:

```text
seconds
bytes
count
ratio
```

Do not publish milliseconds and seconds under the same semantic metric.

---

## 18. Counters

Good uses:

```text
requests
errors
executions
retries
provider operations
events consumed
```

---

## 19. Histograms

Good uses:

```text
request duration
provider operation duration
DB query duration
queue wait
payload size
execution step duration
```

---

## 20. Gauges

Good uses:

```text
in-flight requests
worker concurrency
queue depth
DB pool connections
provider client cache size
```

---

## 21. HTTP Metrics

At minimum:

```text
method
route template
status class/code where appropriate
duration
request count
```

Use route template:

```text
/applications/{applicationId}
```

not raw resource path.

---

## 22. gRPC Metrics

Track:

```text
service
method
status
duration
```

Message size only where useful.

---

## 23. Database Metrics

Track:

```text
pool open/in-use/idle
wait count/duration
query duration by operation name
transaction rollback
serialization/deadlock errors
availability
```

---

## 24. Provider Metrics

Normalize:

```text
capability
operation
provider type
result category
duration
retry count
```

Provider instance IDs require cardinality review before becoming metric labels.

---

## 25. Mission Runtime Metrics

Track:

```text
executions started
completed
failed
cancelled
execution duration
step duration
queue wait
retry count
approval wait
rollback count
worker saturation
```

---

## 26. Event/Queue Metrics

Track:

```text
publish rate
consume rate
lag/backlog
redelivery
dead-letter count
processing duration
```

---

## 27. Cache Metrics

Normalize:

```text
hit
miss
error
duration
eviction where available
```

Vendor-specific metrics can coexist in provider dashboards.

---

## 28. Frontend Performance

Capture user-perceived behavior where useful:

```text
route load
Workbench panel load
initial shell bootstrap
mission start latency
critical API latency
web vitals
```

Do not collect invasive analytics by default.

---

## 29. Tracing

Create spans around meaningful boundaries:

```text
HTTP/gRPC request
business use case
DB operation
provider call
message publish/consume
execution step
```

Do not create spans around trivial getters/helpers.

---

## 30. Span Naming

Good:

```text
catalog.application.create
provider.cache.get
execution.step.run
```

Avoid IDs in span names.

---

## 31. Span Attributes

Use bounded semantic attributes.

High-cardinality identifiers may exist in traces where needed, subject to privacy/retention policy.

---

## 32. Trace Propagation

Use W3C Trace Context for HTTP where supported.

Use standard metadata propagation for gRPC/messaging.

Validate/bound incoming baggage.

---

## 33. Messaging Propagation

Carry trace/correlation metadata in event envelopes.

A consumer creates the appropriate linked/child span based on messaging semantics.

---

## 34. Durable Execution

Mission Runtime must not rely on one continuous trace.

Persist:

```text
executionId
correlationId
step IDs
attempt numbers
```

Each resumed worker may create a new trace linked to durable identity.

---

## 35. Sampling

Development may use 100%.

Production sampling is configurable.

Higher priority sampling may apply to:

```text
errors
security-relevant flows
critical missions
```

within cost/privacy bounds.

---

## 36. Audit vs Logs

Audit asks:

> Who performed a governed action on what target, and what was the result?

Operational logs ask:

> What was the system doing and why did it fail?

Do not use ordinary logs as the only audit trail.

---

## 37. Evidence vs Telemetry

Evidence is a product/governance object.

It can reference:

```text
trace
logs
artifact digest
test result
approval
provider revision
```

but remains durable according to governance requirements even when raw telemetry expires.

---

## 38. Liveness

Liveness answers:

```text
Is this process alive enough to avoid restart?
```

Do not fail liveness because an optional remote provider is unavailable.

---

## 39. Readiness

Readiness answers:

```text
Can this instance safely receive work now?
```

It may depend on:

```text
DB
required config
mandatory internal initialization
```

---

## 40. Optional Provider Failure

An optional provider outage should degrade only its capability where possible.

Do not make whole-platform readiness fail for unused optional tooling.

---

## 41. Alerts

Prefer:

```text
SLO burn
user-visible error increase
mission failure rate
queue backlog
DB exhaustion
provider capability outage
resource saturation
```

Avoid alerting on every transient exception.

---

## 42. SLO Definition

An SLO must define:

```text
service/user experience
population
SLI
target
window
exclusions
```

Do not publish an uptime percentage without semantics.

---

## 43. Dashboard Layers

Recommended:

```text
Product / User Experience
Service
Mission Runtime
Capability / Provider
Infrastructure
```

Operators should not have to start at node CPU to understand a failed deployment.

---

## 44. Provider-Neutral Dashboards

Core dashboards show normalized capabilities.

Provider-specific dashboards may add deeper vendor metrics.

---

## 45. Workbench Correlation

Workbench should expose authorized diagnostic identifiers such as:

```text
executionId
requestId
correlationId
trace link
```

where useful.

---

## 46. Log Viewer

The product log viewer must preserve:

```text
timestamp
level
source
structured fields
correlation
```

and apply redaction.

---

## 47. Retention

Retention differs for:

```text
debug logs
operational logs
metrics
traces
audit
evidence
```

Retention is deployment/governance policy.

---

## 48. Cost Controls

Observability cost is controlled through:

```text
sampling
retention
cardinality
log level
payload suppression
aggregation
```

Do not dump entire objects on every request.

---

## 49. Request/Response Bodies

Do not log full bodies by default.

If diagnostic body capture is enabled, it must be:

```text
explicit
size-bounded
redacted
access-controlled
short-lived
```

---

## 50. SQL Logging

Log normalized repository/query operation names.

Do not log sensitive bind values.

---

## 51. Provider Logging

Include safe provider/capability identifiers.

Never log credentials or full secret-bearing configuration.

---

## 52. Local Reference Stack

Reference development stack:

```text
OpenTelemetry Collector
Prometheus
Loki
Tempo
Grafana
```

Lite mode may omit optional components.

---

## 53. External Observability Products

External tools integrate through:

```text
OTel exporters
observability providers
approved adapters
```

Core code remains provider-neutral.

---

## 54. Export Failure

Telemetry export failures are bounded and non-fatal.

Avoid blocking user requests indefinitely on telemetry backend availability.

---

## 55. Buffers

Queues/buffers have maximum bounds.

Prefer controlled telemetry loss/degradation over product OOM.

---

## 56. Testing

Test critical telemetry behavior:

```text
request/correlation propagation
redaction
trace context propagation
critical metric emission
health/readiness semantics
```

Do not write brittle tests for every log sentence.

---

## 57. Naming Registry

Centralize standard field names.

Avoid repository-wide variants:

```text
request_id
requestId
reqID
```

unless transport mapping requires it.

---

## 58. Build Metadata

Every process exports safe metadata:

```text
component
version
commit SHA
build ID
```

---

## 59. Deployment Metadata

Useful resource attributes:

```text
environment
region/site
instance/pod
deployment profile
```

---

## 60. Observability PR Checklist

```text
[ ] Structured logs used
[ ] Secrets/PII redacted
[ ] Request/correlation/trace context propagated
[ ] Metric labels bounded
[ ] Metric unit/name consistent
[ ] Spans cover meaningful boundary
[ ] Provider telemetry normalized
[ ] Health/readiness semantics correct
[ ] Duplicate error logging avoided
[ ] Failure output actionable
[ ] Cost/volume considered
```

---

## 61. Definition of Done

```text
[ ] Important operation is observable
[ ] Failure can be correlated across layers
[ ] Critical metric/log/trace exists
[ ] No sensitive telemetry leak
[ ] Cardinality is bounded
[ ] Exporter outage does not break product
[ ] Relevant telemetry tests pass
```

---

## 62. Phase 0 Deliverables

```text
OpenTelemetry initialization
structured logger
request/correlation middleware
trace propagation
HTTP/gRPC baseline metrics
DB pool metrics
provider telemetry decorator
Mission Runtime telemetry skeleton
health/readiness endpoints
local OTel export path
redaction helper/tests
Grafana reference dashboard skeleton
```

---

## 63. Agent Rules

Before telemetry work inspect:

```text
docs/AGENT.md
this file
shared telemetry/logger implementation
standard field registry
security redaction rules
existing spans/metrics
```

Do not invent one-off naming.

---

## 64. Final Mandatory Rules

1. OTel-first.
2. Structured logs.
3. Request/correlation/trace are distinct.
4. Context propagates.
5. No secrets in telemetry.
6. Metric cardinality is bounded.
7. Names and units are stable.
8. Core telemetry is provider-neutral.
9. Mission execution has structured telemetry.
10. Liveness/readiness are explicit.
11. Telemetry export failure is non-fatal and bounded.
12. Audit and evidence are not ordinary logs.

---

## 65. Guiding Principle

```text
Structured Signals
      +
Shared Context
      +
Bounded Cardinality
      +
Safe Redaction
      +
Actionable Views
      =
Operable Platform
```

> **When Perihellion fails, an engineer should be able to move from user symptom to execution, service, provider, and root cause without guessing.**

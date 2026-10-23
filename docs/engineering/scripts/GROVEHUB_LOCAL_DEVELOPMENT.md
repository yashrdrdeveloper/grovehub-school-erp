# GroveHub Local Development

> **Status:** Mandatory GroveHub local-development standard  
> **Applies to:** backend developer startup, Air hot reload, Swagger startup, local observability startup, process ownership, Go module refresh, and local backend telemetry configuration  
> **Precedence:** This file is authoritative for GroveHub local development command semantics when an older generic example conflicts with it.

---

## 1. Command Semantics

The stable commands have different purposes:

```text
make backend-start
    -> build and run one backend binary
    -> no hot reload

make backend-dev
    -> run Air in the foreground
    -> hot reload enabled

make dev-start
    -> run backend through Air in the background
    -> start Swagger
    -> hot reload enabled

make dev-start-with-logs
    -> run backend through Air in the background
    -> start Swagger
    -> start the local OpenTelemetry collector
    -> export GroveHub logs, traces, and metrics to Grafana Cloud
    -> hot reload enabled
```

Do not route `dev-start` or `dev-start-with-logs` through a one-time backend binary.

---

## 2. Git Pull and File Changes

Air watches the local backend source tree.

After a developer runs:

```bash
git pull
```

changes to watched backend files must trigger a rebuild and process restart automatically.

Current watched extensions include:

```text
.go
.mod
.sum
```

Air does not perform `git pull` itself. It reacts to filesystem changes produced by the developer, editor, generator, or Git operation.

---

## 3. Dependency Refresh

The backend must synchronize the Go module graph before compiling when dependencies may have changed.

Canonical module synchronization is:

```bash
go mod tidy
```

This is intentionally different from `go mod download`.

`go mod download` may populate the module cache without producing every checksum entry required by a package build. A fresh clone or a branch that added new imports can therefore still fail with:

```text
missing go.sum entry for module providing package ...
```

Canonical Air behavior is:

```text
.go / go.mod / go.sum changes
        ↓
Air detects change
        ↓
scripts/backend/air-build.sh
        ↓
go mod tidy
        ↓
go build
        ↓
old backend process stopped
        ↓
new backend process started
```

The same shared module synchronization helper is used by ordinary backend build/setup paths.

Do not use `go get` as automatic development bootstrap because it can intentionally change selected dependency versions.

`backend/go.sum` is repository state and should be committed whenever dependency resolution changes it. Do not gitignore `go.sum`.

---

## 4. Process Ownership

The repository owns one backend PID file:

```text
.runtime/pids/backend.pid
```

That PID may represent either:

```text
one-time GroveHub backend binary
or
Air supervising the development backend
```

Stop/status logic must recognize both forms and must never use broad `pkill` or `killall` commands.

`make dev-stop` must safely stop the Air supervisor and its managed backend process.

---

## 5. Logs

Both one-time and Air-backed development startup write backend output to:

```text
.runtime/logs/backend.log
```

To watch rebuilds locally:

```bash
tail -f .runtime/logs/backend.log
```

A successful Air reload should show the module-sync, build, and restart sequence in this log.

Shell messages such as `[INFO] Starting GroveHub backend` and Air watcher/build lines are developer tooling output and are not controlled by the application log-level setting. The local OpenTelemetry filelog receiver drops non-JSON developer-tooling lines before exporting to Grafana, so the Grafana application-log stream contains only structured GroveHub records.

The collector must translate the structured GroveHub `level` field into the OpenTelemetry log record severity (`SeverityNumber`/`SeverityText`) explicitly. Grafana/Loki level detection must not be left to inference from the JSON body. For example, a GroveHub request record with `level=ERROR` for HTTP 404 must arrive in Grafana as OpenTelemetry `ERROR` severity.

The collector also maps structured GroveHub `service` and `serviceKey` fields into resource attributes so logs use the same logical service identity as traces and metrics.

---

## 6. Application Log Levels

Local backend application log levels are configured through the root `.env` file using:

```env
GROVEHUB_LOG_LEVELS=ERROR,WARN
```

The value is a comma-separated exact-match list. Supported values are:

```text
DEBUG
INFO
WARN
ERROR
```

Rules:

1. Matching is case-insensitive.
2. Whitespace around values is ignored.
3. Duplicate values are harmless.
4. Unknown values are configuration errors and the backend must fail startup clearly.
5. If `GROVEHUB_LOG_LEVELS` is absent or empty, the compatibility default is `INFO,WARN,ERROR`.
6. This is exact selection, not minimum-threshold semantics. `ERROR,WARN` emits only WARN and ERROR. `DEBUG,ERROR` emits only DEBUG and ERROR.
7. Filtering happens at the GroveHub application logger before logs are written to `.runtime/logs/backend.log`; filtered-out application records therefore are not shipped to Grafana by the local collector.

---

## 7. Trace and Metric Controls

Local tracing and metric export are configured through the root `.env` file:

```env
GROVEHUB_TRACE_SAMPLE_RATIO=1.0
GROVEHUB_METRIC_EXPORT_INTERVAL=15s
OTEL_SERVICE_VERSION=dev
OTEL_DEPLOYMENT_ENVIRONMENT=local
```

Rules:

1. `GROVEHUB_TRACE_SAMPLE_RATIO` is `0.0` through `1.0`; local development defaults to `1.0`.
2. `GROVEHUB_METRIC_EXPORT_INTERVAL` must remain at least 5 seconds; the default is 15 seconds.
3. `OTEL_SERVICE_VERSION` is attached to telemetry resources.
4. `OTEL_SERVICE_INSTANCE_ID` is optional; when absent the backend creates a local `<hostname>-<pid>` instance ID.
5. The application owns logical service identities such as `LAUNCHER (EDGE)` and `IDENTITY`; `OTEL_SERVICE_NAME` in local collector configuration is only a fallback for records that lack a service identity.
6. Changing root `.env` telemetry values does not trigger Air. Restart the backend/dev stack after changing them.

The backend startup scripts may load the root `.env`, but only application runtime variables explicitly approved for the backend are exported to the ordinary backend process. Grafana Cloud credentials are used by the collector and must not be exported merely to configure application logging/metrics.

---

## 8. Local Observability Flow

```text
GroveHub backend
  ├── structured application logs
  ├── OpenTelemetry traces
  └── OpenTelemetry metrics
          ↓
local OpenTelemetry Collector
          ↓
Grafana Cloud
  ├── Loki
  ├── Tempo
  └── Prometheus/Mimir
```

Use `docs/engineering/observability/GROVEHUB_OBSERVABILITY.md` for canonical service identities, metrics, dashboard filters, sampling, cardinality, and alert rules.

---

## 9. Mandatory Agent Rules

When changing local developer startup:

1. Preserve the command semantics in this document.
2. Keep the root Makefile thin.
3. Put orchestration in `scripts/`.
4. Reuse shared process/toolchain/module/environment helpers.
5. Keep `backend-start` production-like and non-reloading.
6. Keep `dev-start` and `dev-start-with-logs` Air-backed.
7. Ensure pulled `.go`, `go.mod`, and `go.sum` changes trigger a rebuild.
8. Use `go mod tidy` to synchronize the module graph before builds that must self-heal after dependency changes.
9. Treat `backend/go.sum` as committed repository state.
10. Preserve `GROVEHUB_LOG_LEVELS` exact-list semantics and supported values.
11. Do not silently fall back when an explicitly configured log level is invalid.
12. Keep telemetry runtime variables explicit; do not export unrelated `.env` secrets into the backend process.
13. Keep local collector filelog export limited to structured GroveHub application records; do not ship Air/watch/build chatter as application logs.
14. Map GroveHub structured `level` to OpenTelemetry severity explicitly before export; do not rely on Grafana/Loki inferred level detection for application severity.
15. Keep log resource service identity aligned with trace/metric service identity.
16. Preserve traces + metrics export in `dev-start-with-logs`, not only log shipping.
17. Never add broad process termination.
18. Report the actual start/stop/hot-reload/telemetry validation performed.

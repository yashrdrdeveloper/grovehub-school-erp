# GroveHub Local Development

> **Status:** Mandatory GroveHub local-development standard  
> **Applies to:** backend developer startup, Air hot reload, Swagger startup, local observability startup, process ownership, Go module refresh, and local backend log configuration  
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
8. Changing the root `.env` does not itself trigger Air because Air watches the backend source/module tree. Restart the backend/dev stack after changing `GROVEHUB_LOG_LEVELS`.

The backend startup scripts may load the root `.env`, but only application runtime variables explicitly approved for the backend should be exported to the backend process. Grafana Cloud credentials must not be exported to the ordinary backend merely to configure logging.

---

## 7. Mandatory Agent Rules

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
12. Do not export unrelated `.env` secrets into the backend process.
13. Keep local collector filelog export limited to structured GroveHub application records; do not ship Air/watch/build chatter as application logs.
14. Never add broad process termination.
15. Report the actual start/stop/hot-reload validation performed.

# GroveHub Local Development

> **Status:** Mandatory GroveHub local-development standard  
> **Applies to:** backend developer startup, Air hot reload, Swagger startup, local observability startup, process ownership, and Go module refresh during development  
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

The Air build command must resolve the modules declared in `go.mod` before compiling so a pulled dependency change does not leave the running development process stale.

Canonical behavior:

```text
go.mod / go.sum changes
        ↓
Air detects change
        ↓
go mod download
        ↓
go build
        ↓
old backend process stopped
        ↓
new backend process started
```

Do not use `go get` as an automatic development bootstrap because it may change dependency versions.

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

A successful Air reload should show the build/restart sequence in this log.

---

## 6. Mandatory Agent Rules

When changing local developer startup:

1. Preserve the command semantics in this document.
2. Keep the root Makefile thin.
3. Put orchestration in `scripts/`.
4. Reuse shared process/toolchain helpers.
5. Keep `backend-start` production-like and non-reloading.
6. Keep `dev-start` and `dev-start-with-logs` Air-backed.
7. Ensure pulled `.go`, `go.mod`, and `go.sum` changes trigger a rebuild.
8. Never add broad process termination.
9. Report the actual start/stop/hot-reload validation performed.

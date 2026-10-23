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

---

## 6. Mandatory Agent Rules

When changing local developer startup:

1. Preserve the command semantics in this document.
2. Keep the root Makefile thin.
3. Put orchestration in `scripts/`.
4. Reuse shared process/toolchain/module helpers.
5. Keep `backend-start` production-like and non-reloading.
6. Keep `dev-start` and `dev-start-with-logs` Air-backed.
7. Ensure pulled `.go`, `go.mod`, and `go.sum` changes trigger a rebuild.
8. Use `go mod tidy` to synchronize the module graph before builds that must self-heal after dependency changes.
9. Treat `backend/go.sum` as committed repository state.
10. Never add broad process termination.
11. Report the actual start/stop/hot-reload validation performed.

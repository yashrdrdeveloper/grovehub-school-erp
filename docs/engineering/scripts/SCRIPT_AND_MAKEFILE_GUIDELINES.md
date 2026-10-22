# Scripts & Makefile Engineering Guidelines

> **Status:** Mandatory Phase 0 engineering standard  
> **Location:** `docs/engineering/scripts/SCRIPT_AND_MAKEFILE_GUIDELINES.md`  
> **Applies to:** root `Makefile`, `scripts/`, local development bootstrap, toolchain download, Docker Compose orchestration, frontend/backend start-stop commands, test runners, and CI-facing command wrappers  
> **Supported developer hosts:** macOS, Linux, and Windows through WSL2  
> **Primary shell:** Bash  
> **Core principle:** **one predictable command for the user; small, explicit, reusable scripts underneath**

---

# 1. Objective

A contributor should be able to clone Perihellion and run common workflows without manually installing or remembering project-specific versions of Node, pnpm, Go, Playwright, database tools, or development services.

The intended user experience is:

```bash
git clone <perihellion-repository>
cd <perihellion-repository>

make frontend-start
```

or:

```bash
make backend-start
```

or:

```bash
make dev-start
```

The automation layer owns:

```text
host detection
toolchain version resolution
repo-local tool download
package dependency installation
container-runtime detection
Docker Compose startup
health checks
PID/process lifecycle
test orchestration
consistent environment wiring
diagnostics
```

The developer should not need to know the internal sequence.

---

# 2. Core Philosophy

Use:

```text
Makefile
   ↓
small command script
   ↓
shared script library
   ↓
tool/runtime/application command
```

Do not place complex orchestration directly in the Makefile.

Do not duplicate shell logic across frontend/backend/testing scripts.

The Makefile is a **stable human interface**.

The scripts are the **implementation**.

---

# 3. Canonical Repository Structure

```text
perihellion/
├── Makefile
│
├── scripts/
│   ├── bootstrap.sh
│   ├── doctor.sh
│   ├── testing.sh
│   ├── clean.sh
│   │
│   ├── lib/
│   │   ├── common.sh
│   │   ├── os.sh
│   │   ├── paths.sh
│   │   ├── process.sh
│   │   ├── toolchain.sh
│   │   ├── download.sh
│   │   ├── container.sh
│   │   └── compose.sh
│   │
│   ├── toolchain/
│   │   ├── versions.env
│   │   ├── checksums/
│   │   └── README.md
│   │
│   ├── assets/
│   │   ├── README.md
│   │   └── ...
│   │
│   ├── docker/
│   │   ├── compose.yml
│   │   ├── compose.dev.yml
│   │   ├── compose.test.yml
│   │   ├── start.sh
│   │   ├── stop.sh
│   │   ├── reset.sh
│   │   └── status.sh
│   │
│   ├── frontend/
│   │   ├── setup.sh
│   │   ├── start.sh
│   │   ├── stop.sh
│   │   ├── dev.sh
│   │   ├── build.sh
│   │   └── clean.sh
│   │
│   ├── backend/
│   │   ├── setup.sh
│   │   ├── start.sh
│   │   ├── stop.sh
│   │   ├── dev.sh
│   │   ├── build.sh
│   │   ├── migrate.sh
│   │   └── clean.sh
│   │
│   └── testing/
│       ├── frontend/
│       │   ├── unit.sh
│       │   ├── component.sh
│       │   ├── lint.sh
│       │   └── typecheck.sh
│       │
│       ├── backend/
│       │   ├── unit.sh
│       │   ├── integration.sh
│       │   ├── contract.sh
│       │   ├── lint.sh
│       │   └── vet.sh
│       │
│       └── regression/
│           ├── smoke.sh
│           ├── critical.sh
│           ├── all.sh
│           └── domain.sh
│
├── frontend/
├── backend/
├── tests/
│   └── regression/
│
├── .tools/          # generated; gitignored
├── .runtime/        # generated; gitignored
└── .state/          # generated; gitignored
```

Do not create every optional file immediately. Add a script when a real command needs it.

---

# 4. Generated Runtime Directories

These are repository-local generated directories:

```text
.tools/
.runtime/
.state/
```

They must be gitignored.

## `.tools/`

Repo-local installed project toolchains.

Example:

```text
.tools/
├── node/
├── pnpm/
├── go/
└── bin/
```

## `.runtime/`

Process and runtime artifacts.

Example:

```text
.runtime/
├── pids/
│   ├── frontend.pid
│   └── backend.pid
├── logs/
│   ├── frontend.log
│   └── backend.log
└── tmp/
```

## `.state/`

Small bootstrap state.

Example:

```text
.state/
├── frontend-lock.hash
├── backend-toolchain.hash
└── bootstrap.version
```

Never use these as product persistence.

---

# 5. Shared Download Cache

Downloaded archives should not be repeatedly fetched.

Use a shared cache.

Recommended environment variable:

```text
PERIHELLION_CACHE_HOME
```

Default:

```text
macOS/Linux/WSL:
$HOME/.cache/perihellion
```

Possible contents:

```text
$PERIHELLION_CACHE_HOME/
├── node/
├── go/
├── pnpm/
├── playwright/
└── downloads/
```

A cached artifact must still pass checksum validation.

---

# 6. `scripts/assets/`

`scripts/assets/` contains **small, checked-in, immutable support assets** used by scripts.

Examples:

```text
config templates
small seed templates
development certificates when non-secret and explicitly intended
example env fragments
test bootstrap templates
```

Do not store downloaded binaries here.

Do not store secrets here.

Downloaded artifacts belong in the shared cache.

---

# 7. Toolchain Version Sources

Versions must have one authoritative source.

Recommended:

```text
Node:
.node-version

pnpm:
frontend/package.json -> packageManager

Go:
backend/go.mod -> toolchain/go directive
or .go-version when explicitly standardized

Playwright:
frontend/test package lock / regression package lock

Container images:
scripts/docker/compose*.yml
```

If a central shell-readable manifest is needed:

```text
scripts/toolchain/versions.env
```

may contain derived/bootstrap-only values.

Do not maintain the same version manually in five scripts.

---

# 8. Toolchain Bootstrap Rule

Project language toolchains should be bootstrapped automatically into `.tools/`.

Example:

```text
make frontend-start
        ↓
ensure Node
        ↓
ensure pnpm
        ↓
ensure frontend dependencies
        ↓
start frontend
```

The user should not need:

```bash
nvm install ...
npm install -g pnpm@...
```

manually.

---

# 9. Node Bootstrap

`scripts/lib/toolchain.sh` should resolve the required Node version and install the exact OS/architecture build into:

```text
.tools/node/
```

Supported architecture resolution should include as required:

```text
darwin-arm64
darwin-x64
linux-arm64
linux-x64
```

Windows development through WSL uses the Linux build.

Node downloads must:

```text
1. resolve exact version
2. determine OS/architecture
3. check shared cache
4. download from approved source if absent
5. verify checksum
6. extract atomically
7. expose through .tools/bin or PATH wrapper
```

Do not mutate a developer's global Node installation.

---

# 10. npm

npm normally ships with Node.

Do not separately install a random global npm version unless the project explicitly requires it.

If npm must be pinned independently, define its source of truth and install it inside the repo-local toolchain.

---

# 11. pnpm Bootstrap

pnpm must use the version declared by the project.

Preferred source:

```json
{
  "packageManager": "pnpm@<version>"
}
```

The bootstrap script may install the exact pnpm version using the bootstrapped Node/npm or another verified official distribution.

Install into repo-owned tooling, not system-global state.

A developer's globally installed pnpm must not change project behavior.

---

# 12. Frontend Dependency Installation

`scripts/frontend/setup.sh` owns frontend setup.

Conceptual flow:

```text
ensure Node
ensure pnpm
calculate lockfile/package manifest state
install when required
validate installation
```

Use:

```bash
pnpm install --frozen-lockfile
```

in CI and reproducible bootstrap paths.

Do not silently rewrite the lockfile during startup.

---

# 13. Go Bootstrap

The backend bootstrap should resolve the required Go version and place it under:

```text
.tools/go/
```

Conceptual flow:

```text
resolve version
detect OS/arch
check cache
download
verify checksum
extract
add .tools/go/bin to command environment
```

Do not require developers to change their global Go installation.

---

# 14. Tool Download Security

Every direct binary/archive download must have:

```text
approved HTTPS source
exact version
checksum verification
atomic extraction
failure on mismatch
```

Never:

```bash
curl <url> | bash
```

for project-language toolchains in normal bootstrap.

If an official installer must be used for a system-level prerequisite, make the action explicit and user-visible.

---

# 15. Offline / Mirror Support

Bootstrap must allow an enterprise/local mirror later.

Reserve:

```text
PERIHELLION_ASSET_MIRROR
```

When configured:

```text
approved mirror
    ↓
cache
    ↓
.tools
```

The same checksums still apply.

This supports air-gapped evolution without redesigning the bootstrap architecture.

---

# 16. System-Level vs Repo-Level Dependencies

There are two classes.

## Repo-level

Automatically managed:

```text
Node
npm bundled with Node
pnpm
Go
project JavaScript dependencies
Playwright package dependencies
generated local binaries where practical
```

## Host/system-level

Examples:

```text
container runtime
virtualization
WSL2
kernel capabilities
```

These may require elevated privileges or OS configuration.

Scripts must detect them and provide or offer a safe bootstrap path.

Do not silently make privileged system changes.

---

# 17. Official Host Support

Official developer environments:

```text
macOS
Linux
Windows via WSL2
```

Windows-native PowerShell support may be added later, but Bash/Make behavior is canonical.

This avoids maintaining two divergent orchestration systems.

---

# 18. macOS Container Runtime

On macOS:

```text
1. use an already working Docker daemon if present
2. otherwise use/start Colima when installed
3. if Colima is missing and an approved package manager is available,
   bootstrap may offer to install it
4. verify Docker Compose after startup
```

Colima is an implementation of the local Docker-compatible runtime, not a product dependency.

Scripts must not hardcode application behavior based on `colima`.

---

# 19. Linux Container Runtime

On Linux:

```text
1. detect Docker CLI
2. detect Docker daemon availability
3. detect Docker Compose plugin
4. use current user permissions
5. provide clear diagnostics if daemon/socket access fails
```

Do not automatically prefix every Docker command with `sudo`.

A developer environment should be configured so normal project commands do not require root.

---

# 20. Windows Support

Canonical Windows development uses:

```text
Windows
  ↓
WSL2
  ↓
Bash + Make
  ↓
Docker Desktop WSL integration
or supported Docker Engine inside WSL
```

Scripts detect WSL where relevant.

Do not create a separate Windows architecture unless native Windows support becomes a deliberate product requirement.

---

# 21. Container Runtime Abstraction

All scripts requiring Compose should use one shared helper:

```text
scripts/lib/container.sh
scripts/lib/compose.sh
```

Individual scripts should not repeatedly implement:

```bash
docker info
docker compose version
colima status
colima start
```

The shared helper owns this detection.

---

# 22. Docker Compose Location

Development Compose files belong under:

```text
scripts/docker/
```

because they are part of the zero-setup developer orchestration.

Production deployment manifests remain separate under:

```text
deployments/
```

Do not treat `scripts/docker/compose.yml` as the production deployment definition.

---

# 23. Compose Files

Recommended:

```text
scripts/docker/
├── compose.yml
├── compose.dev.yml
└── compose.test.yml
```

## `compose.yml`

Base local dependencies.

Example:

```text
PostgreSQL
optional local object storage
optional cache/events as capabilities are enabled
```

## `compose.dev.yml`

Development overrides.

## `compose.test.yml`

Isolated test infrastructure.

---

# 24. Compose Profiles

Use profiles when dependencies are optional.

Example:

```text
lite
team
observability
testing
```

Do not start every possible Perihellion dependency for every command.

A frontend-only developer should not need the full backend infrastructure unless the workflow actually requires it.

---

# 25. Compose Image Rules

Never use floating production-sensitive tags such as:

```text
latest
```

Pin versions.

Prefer exact versions and, where security/reproducibility needs justify it, digests.

---

# 26. Compose Health Checks

Infrastructure containers should expose meaningful health checks.

A script must not report:

```text
PostgreSQL ready
```

merely because the container process exists.

Wait for health/readiness.

---

# 27. Compose Project Name

Use a stable project name to prevent accidental collisions.

Example:

```text
perihellion
```

For parallel test runs, derive an isolated project name:

```text
perihellion-test-<run-id>
```

---

# 28. Compose Volumes

Development volumes must be intentionally named.

Test suites should prefer disposable/isolated volumes.

`make docker-reset` may delete local development data only after an explicit confirmation or force flag.

Never delete arbitrary Docker volumes.

---

# 29. Script Size Rule

Scripts must stay small.

Target:

```text
command wrapper:   ~20–80 lines
shared helper:     ~20–150 lines
```

When logic grows:

```text
extract a named helper
```

Do not create a 700-line `setup.sh`.

---

# 30. Shell Script Baseline

Every Bash script starts with:

```bash
#!/usr/bin/env bash
set -Eeuo pipefail
```

Use:

```bash
IFS=$'\n\t'
```

only when the script's behavior specifically benefits from it.

Do not add ceremony without understanding it.

---

# 31. Root Resolution

Every script must resolve repository root reliably.

Do not assume the caller's current working directory.

Shared helper example:

```bash
SCRIPT_DIR="$(
  cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1
  pwd
)"

REPO_ROOT="$(
  cd -- "${SCRIPT_DIR}/../.." >/dev/null 2>&1
  pwd
)"
```

Prefer centralizing this in `scripts/lib/paths.sh`.

---

# 32. Quoting

Always quote variable expansions unless word splitting is deliberately required.

Good:

```bash
rm -rf -- "${TARGET_DIR}"
```

Bad:

```bash
rm -rf $TARGET_DIR
```

---

# 33. Temporary Files

Use:

```bash
mktemp
```

and cleanup traps.

Do not create predictable temporary files in `/tmp`.

---

# 34. Cleanup Traps

For temporary resources:

```bash
cleanup() {
  rm -f -- "${TMP_FILE:-}"
}

trap cleanup EXIT
```

Long-lived processes should have deliberate lifecycle management instead of broad cleanup traps.

---

# 35. Error Helpers

Use shared helpers.

Example:

```bash
log_info "Starting frontend"
log_warn "Docker daemon is not available"
die "Required Node version could not be installed"
```

Do not manually reproduce color/logging code in every script.

---

# 36. Non-Interactive Behavior

Scripts must support CI.

Do not block on a prompt when:

```text
CI=true
PERIHELLION_NON_INTERACTIVE=1
```

If a privileged/system-level operation would require confirmation, fail with actionable output in non-interactive mode.

---

# 37. Idempotency

Setup/start/stop scripts should be idempotent where possible.

Examples:

```text
setup twice → same valid environment
start when already running → report running, do not duplicate
stop when stopped → success
docker-start when healthy → success
```

---

# 38. No Broad Process Killing

Forbidden:

```bash
pkill node
killall node
pkill go
```

Use process ownership.

Store PIDs under:

```text
.runtime/pids/
```

Verify the process before sending a signal.

---

# 39. Process Start Model

Provide two modes.

## Background

```bash
make frontend-start
make backend-start
make dev-start
```

These return after readiness and store:

```text
PID
log path
port
```

## Foreground development

```bash
make frontend-dev
make backend-dev
```

These attach to the process for interactive logs/HMR.

---

# 40. Frontend Start

`make frontend-start` calls:

```text
scripts/frontend/start.sh
```

Conceptual flow:

```text
resolve repository root
ensure frontend toolchain
ensure dependencies
load approved local environment
check target port
start frontend in background
write PID/log
wait for frontend readiness
print URL
```

It should not require the developer to manually install Node or pnpm.

---

# 41. Frontend Stop

`make frontend-stop` calls:

```text
scripts/frontend/stop.sh
```

It:

```text
reads owned PID
verifies process
requests graceful termination
waits
escalates only if required
removes stale PID
```

It must not kill unrelated Node processes.

---

# 42. Backend Start

`make backend-start` calls:

```text
scripts/backend/start.sh
```

Conceptual flow:

```text
ensure Go toolchain
ensure required development infrastructure
apply/verify migrations according to dev policy
check backend port
start backend
write PID/log
wait for /ready
print URL/status
```

The backend script does not duplicate Docker Compose logic; it calls shared Compose helpers.

---

# 43. Backend Stop

`make backend-stop`:

```text
graceful application shutdown
owned PID only
idempotent when already stopped
```

Local infrastructure is not necessarily destroyed when backend stops.

Use separate:

```bash
make infra-stop
```

for dependencies.

---

# 44. Development Start

`make dev-start` should provide the common full local experience.

Recommended flow:

```text
doctor
  ↓
toolchain bootstrap
  ↓
infra-start
  ↓
backend-start
  ↓
frontend-start
  ↓
health summary
```

Output example:

```text
Perihellion development environment

Frontend   UP   http://localhost:<frontend-port>
Backend    UP   http://localhost:<backend-port>
Postgres   UP
```

Do not print success before readiness checks pass.

---

# 45. Development Stop

Recommended:

```text
frontend-stop
backend-stop
infra-stop
```

Stop in dependency-safe reverse order.

---

# 46. Testing Structure

Required paths:

```text
scripts/testing/frontend/
scripts/testing/backend/
scripts/testing/regression/
scripts/testing.sh
```

`scripts/testing.sh` is the aggregate test entry point.

---

# 47. Frontend Test Scripts

Recommended:

```text
scripts/testing/frontend/
├── unit.sh
├── component.sh
├── lint.sh
└── typecheck.sh
```

Each script should run one clear responsibility.

Example:

```bash
make test-frontend
```

may orchestrate:

```text
lint
typecheck
unit/component tests
```

Do not mix regression Playwright into frontend unit scripts.

---

# 48. Backend Test Scripts

Recommended:

```text
scripts/testing/backend/
├── unit.sh
├── integration.sh
├── contract.sh
├── lint.sh
└── vet.sh
```

Integration scripts can start isolated Compose test dependencies through shared helpers.

They must not silently use a developer's production/shared database.

---

# 49. Regression Test Scripts

Recommended:

```text
scripts/testing/regression/
├── smoke.sh
├── critical.sh
├── all.sh
└── domain.sh
```

These call the Playwright harness described in:

```text
docs/engineering/testing/REGRESSION_TESTING.md
```

Do not duplicate Playwright configuration inside shell scripts.

---

# 50. Aggregate Testing

`scripts/testing.sh` should stay small.

Conceptually:

```bash
#!/usr/bin/env bash
set -Eeuo pipefail

run frontend tests
run backend tests
run required regression gate
```

The exact suite can be controlled by Make targets/environment.

Do not turn `testing.sh` into the testing framework.

---

# 51. Makefile Responsibility

The root `Makefile` is the stable command catalog.

It should:

```text
declare targets
declare target dependencies
pass variables
call scripts
provide help
```

It should not contain:

```text
large inline shell blocks
OS detection
download logic
complex Docker logic
business logic
large test orchestration
```

---

# 52. Makefile Style

Use portable make syntax where practical.

Define targets as `.PHONY` when they do not represent files.

Example:

```make
.PHONY: frontend-start frontend-stop

frontend-start:
	@./scripts/frontend/start.sh

frontend-stop:
	@./scripts/frontend/stop.sh
```

Keep recipes thin.

---

# 53. Required Make Targets

Recommended Phase 0 public interface:

```text
make help

make bootstrap
make doctor

make frontend-setup
make frontend-start
make frontend-stop
make frontend-dev
make frontend-build

make backend-setup
make backend-start
make backend-stop
make backend-dev
make backend-build
make backend-migrate

make infra-start
make infra-stop
make infra-status
make infra-reset

make dev-start
make dev-stop
make dev-status

make test
make test-frontend
make test-backend
make test-backend-unit
make test-backend-integration
make test-backend-contract
make test-regression
make test-regression-smoke
make test-regression-critical

make clean
```

Names are stable user-facing contracts.

Do not rename them casually.

---

# 54. Example Makefile

```make
SHELL := /usr/bin/env bash

.DEFAULT_GOAL := help

.PHONY: help
help:
	@./scripts/help.sh

.PHONY: bootstrap doctor
bootstrap:
	@./scripts/bootstrap.sh

doctor:
	@./scripts/doctor.sh

.PHONY: frontend-setup frontend-start frontend-stop frontend-dev frontend-build
frontend-setup:
	@./scripts/frontend/setup.sh

frontend-start:
	@./scripts/frontend/start.sh

frontend-stop:
	@./scripts/frontend/stop.sh

frontend-dev:
	@./scripts/frontend/dev.sh

frontend-build:
	@./scripts/frontend/build.sh

.PHONY: backend-setup backend-start backend-stop backend-dev backend-build backend-migrate
backend-setup:
	@./scripts/backend/setup.sh

backend-start:
	@./scripts/backend/start.sh

backend-stop:
	@./scripts/backend/stop.sh

backend-dev:
	@./scripts/backend/dev.sh

backend-build:
	@./scripts/backend/build.sh

backend-migrate:
	@./scripts/backend/migrate.sh

.PHONY: infra-start infra-stop infra-status
infra-start:
	@./scripts/docker/start.sh

infra-stop:
	@./scripts/docker/stop.sh

infra-status:
	@./scripts/docker/status.sh

.PHONY: dev-start dev-stop
dev-start: infra-start backend-start frontend-start

dev-stop:
	@./scripts/frontend/stop.sh
	@./scripts/backend/stop.sh
	@./scripts/docker/stop.sh

.PHONY: test test-frontend test-backend test-regression
test:
	@./scripts/testing.sh

test-frontend:
	@./scripts/testing/frontend/unit.sh

test-backend:
	@./scripts/testing/backend/unit.sh
	@./scripts/testing/backend/integration.sh
	@./scripts/testing/backend/contract.sh

test-regression:
	@./scripts/testing/regression/all.sh
```

The production Makefile may use target dependencies and shared variables more extensively, but complex behavior stays in scripts.

---

# 55. Make Variables

Allow controlled overrides.

Examples:

```bash
make frontend-start PORT=5173
make test-regression DOMAIN=catalog
make infra-start PROFILE=lite
```

The Makefile should pass them to scripts as environment variables.

Do not build large conditional trees in Make.

---

# 56. Environment Files

Provide:

```text
.env.example
```

Never commit real secrets.

Generated local values may live under:

```text
.runtime/
```

or a gitignored:

```text
.env.local
```

Scripts should define clear precedence.

Recommended:

```text
explicit environment
    ↓
.env.local
    ↓
.env.example/defaults
```

Do not silently read random shell-specific configuration files.

---

# 57. Environment Variable Prefix

Automation-specific variables should use:

```text
PERIHELLION_
```

Examples:

```text
PERIHELLION_CACHE_HOME
PERIHELLION_ASSET_MIRROR
PERIHELLION_NON_INTERACTIVE
PERIHELLION_AUTO_INSTALL_SYSTEM_TOOLS
PERIHELLION_COMPOSE_PROFILE
```

Application runtime variables may follow backend/frontend configuration standards.

---

# 58. `make doctor`

`make doctor` is mandatory.

It should report:

```text
OS
architecture
WSL status
repository root
Node required/installed
pnpm required/installed
Go required/installed
Docker CLI
Docker daemon
Docker Compose
Colima state on macOS when relevant
required ports
frontend dependencies
backend dependencies
Playwright readiness
```

Use status:

```text
OK
WARN
ERROR
```

Do not mutate the system in `doctor`.

---

# 59. `make bootstrap`

`make bootstrap` prepares the repository.

Conceptual flow:

```text
detect host
create generated dirs
ensure repo toolchains
install project dependencies
verify container runtime
prepare local config
run basic doctor
```

It should be safe to rerun.

---

# 60. Automatic System Tool Installation

Repo-local tools may be automatically installed.

System-level tools that require package managers, admin privileges, virtualization, or OS changes must be explicit.

Recommended behavior:

```text
interactive:
offer to install/fix supported prerequisite

non-interactive:
fail with exact actionable reason
```

Optional:

```text
PERIHELLION_AUTO_INSTALL_SYSTEM_TOOLS=1
```

may enable approved unattended host installation paths.

Never silently run `sudo`.

---

# 61. Logs

Background development processes log to:

```text
.runtime/logs/
```

Provide:

```bash
make frontend-logs
make backend-logs
```

if useful.

Do not discard startup output that is needed for diagnostics.

---

# 62. Readiness

Every start script must distinguish:

```text
process started
```

from:

```text
service ready
```

A start target succeeds only when readiness succeeds or the command clearly documents asynchronous startup.

For Perihellion, prefer waiting for readiness.

---

# 63. Ports

Do not hardcode ports in multiple scripts.

Define them in one configuration source.

Before starting:

```text
check whether port is already owned
identify likely Perihellion process when possible
fail clearly on conflict
```

Do not kill the conflicting process automatically.

---

# 64. Migrations

Backend startup behavior regarding migrations must be explicit.

For local development, automatic migration may be allowed by policy.

For production, migration execution is a deployment concern.

Do not embed irreversible production migration behavior into generic `backend-start.sh`.

---

# 65. Cleanup

`make clean` may remove:

```text
build output
runtime PIDs/logs
temporary generated state
repo-local tool output when explicitly requested
```

Separate destructive cleanup:

```text
make clean-tools
make infra-reset
```

Do not make normal `clean` destroy local databases unexpectedly.

---

# 66. ShellCheck

Bash scripts must be checked with ShellCheck.

CI should run something equivalent to:

```bash
shellcheck scripts/**/*.sh
```

Account for shell globbing portability through the CI implementation.

Intentional suppressions require a nearby explanation.

---

# 67. Formatting

Use `shfmt` for Bash formatting when adopted.

Formatting is automated.

Do not spend review time debating indentation.

---

# 68. Compose Validation

CI should validate Compose files.

Example:

```bash
docker compose \
  -f scripts/docker/compose.yml \
  -f scripts/docker/compose.test.yml \
  config
```

Do not wait until runtime to discover malformed Compose YAML.

---

# 69. Makefile Validation

At minimum:

```text
make help
make doctor
make -n <selected-targets>
```

should remain valid.

Where practical, run quick bootstrap/start/stop smoke tests in CI.

---

# 70. Script Tests

Critical shared script helpers should be tested when logic is non-trivial.

Examples:

```text
OS/architecture detection
version parsing
checksum verification
PID validation
Compose project naming
```

Do not create unit tests for trivial one-line wrappers.

---

# 71. Download Retries

Network downloads may use bounded retry with backoff.

Never infinite retry.

Errors must show:

```text
artifact
version
source
cache path
reason
```

Do not leak credentials embedded in private mirror URLs.

---

# 72. Proxy Support

Respect standard proxy environment variables:

```text
HTTP_PROXY
HTTPS_PROXY
NO_PROXY
```

Do not invent a parallel proxy configuration unless needed.

---

# 73. CI Uses the Same Entry Points

CI should call the same Make/script entry points developers use where practical.

Good:

```bash
make test-backend
make test-frontend
make test-regression-smoke
```

Avoid maintaining a second unrelated set of CI-only commands that drift from local behavior.

---

# 74. CI Does Not Need Developer Background Mode

Scripts may support:

```text
developer background process mode
CI foreground/service mode
```

but should share the same setup/bootstrap logic.

Do not fork the architecture into `local` and `CI` copies.

---

# 75. Platform-Specific Logic

OS-specific differences belong in:

```text
scripts/lib/os.sh
scripts/lib/container.sh
```

not scattered across application scripts.

Example:

```text
macOS → Colima handling
Linux → Docker daemon handling
WSL2 → Docker Desktop integration checks
```

---

# 76. Architecture Dependency

Scripts may orchestrate code.

Code must not depend on scripts for product semantics.

Forbidden:

```text
business rule only implemented in shell
API contract generated from ad-hoc runtime script state
application behavior dependent on Make target
```

Scripts are tooling/orchestration, not product-domain logic.

---

# 77. No Duplicate Package Manager Logic

Only one shared frontend setup path should decide:

```text
Node version
pnpm version
install command
lockfile policy
```

Frontend testing/build/start scripts reuse it.

They must not each install Node separately.

---

# 78. No Duplicate Backend Toolchain Logic

Only one backend toolchain path should decide:

```text
Go version
Go path
module download setup
```

Build/test/start scripts reuse it.

---

# 79. No Duplicate Compose Logic

Only:

```text
scripts/lib/compose.sh
scripts/docker/*
```

should own Compose invocation details.

Other scripts request:

```text
ensure dependency X
start profile Y
```

through shared functions.

---

# 80. Command Output

User-facing automation output should be concise and structured.

Good:

```text
[frontend] Node 24.x ready
[frontend] pnpm ready
[frontend] dependencies unchanged
[frontend] starting...
[frontend] ready: http://localhost:5173
```

Do not print pages of debug output by default.

Allow:

```text
PERIHELLION_VERBOSE=1
```

for detailed diagnostics.

---

# 81. Failure Output

A failed command must say:

```text
what failed
why it likely failed
where logs are
what the user can run next
```

Example:

```text
ERROR: Backend did not become ready within 30s.
Log: .runtime/logs/backend.log
Run: make doctor
```

---

# 82. `make help`

`make help` should list only stable public commands with short descriptions.

Example:

```text
Development
  make dev-start              Start local Perihellion
  make dev-stop               Stop local Perihellion

Frontend
  make frontend-start         Start frontend in background
  make frontend-dev           Start frontend interactively
  make frontend-stop          Stop frontend

Backend
  make backend-start          Start backend in background
  make backend-dev            Start backend interactively
  make backend-stop           Stop backend

Testing
  make test                   Run standard test gate
  make test-frontend          Run frontend tests
  make test-backend           Run backend tests
  make test-regression        Run Playwright regression

Infrastructure
  make infra-start            Start local dependencies
  make infra-stop             Stop local dependencies

Diagnostics
  make doctor                 Validate local development environment
```

Internal Make targets may be undocumented and prefixed consistently if needed.

---

# 83. Target Naming

Use:

```text
<noun>-<verb>
```

where it remains natural.

Examples:

```text
frontend-start
frontend-stop
backend-build
infra-reset
test-regression
```

Keep names predictable.

---

# 84. Backward Compatibility of Make Targets

Once documented publicly, Make targets are part of the developer interface.

Avoid unnecessary renames.

If a target is renamed, provide a temporary alias and deprecation message where practical.

---

# 85. Open-Source Quick Start

The release/repository goal is:

```text
clone
 ↓
make dev-start
 ↓
browser opens/URL printed
```

No manual:

```text
install exact Node
install exact pnpm
install exact Go
create PostgreSQL
remember migration command
remember frontend command
remember backend command
```

The automation owns those details.

---

# 86. Minimal Host Prerequisite Reality

A completely dependency-free host is not realistic because virtualization/container execution is host-level.

The project therefore distinguishes:

```text
project dependencies → automatically managed
host capabilities    → detected and safely bootstrapped/offered where possible
```

Do not claim zero prerequisites when the operating system itself must provide virtualization/runtime capability.

---

# 87. Phase 0 Deliverables

Phase 0 automation is complete when:

```text
[ ] root Makefile exists
[ ] make help works
[ ] make doctor works
[ ] make bootstrap is idempotent
[ ] frontend toolchain is repo-local/pinned
[ ] backend toolchain is repo-local/pinned
[ ] pnpm version is pinned
[ ] shared download cache works
[ ] checksums are verified
[ ] scripts/docker/compose.yml exists
[ ] local PostgreSQL starts through shared Compose logic
[ ] macOS Docker/Colima detection works
[ ] Linux Docker detection works
[ ] WSL2 Docker detection works
[ ] frontend-start/start-stop lifecycle works
[ ] backend-start/start-stop lifecycle works
[ ] dev-start/dev-stop works
[ ] scripts/testing/frontend exists
[ ] scripts/testing/backend exists
[ ] scripts/testing/regression exists
[ ] scripts/testing.sh exists
[ ] make test-frontend works
[ ] make test-backend works
[ ] make test-regression-smoke works
[ ] scripts are ShellCheck-clean
[ ] Compose config validation passes
[ ] clean clone bootstrap is documented and tested
```

---

# 88. Architecture Enforcement

CI should eventually verify:

```text
Makefile remains thin
scripts are ShellCheck-clean
Compose config is valid
tool versions are pinned
download checksums exist
no secrets are committed
no curl-pipe-shell project tool bootstrap
no broad pkill/killall
no floating Compose latest tags
testing entry points work
quick-start smoke works
```

---

# 89. Agent Rules for Script Changes

Before modifying automation, an AI agent must read:

```text
docs/AGENT.md
docs/engineering/scripts/SCRIPT_AND_MAKEFILE_GUIDELINES.md
```

Then inspect:

```text
Makefile
affected script
shared scripts/lib helpers
toolchain source of truth
Compose files if relevant
commands that consume the script
CI usage if relevant
```

If changing frontend bootstrap, also read the frontend guideline.

If changing backend bootstrap, also read the backend guideline.

If changing test orchestration, also read the regression/testing guideline as applicable.

---

# 90. Script Change Checklist

```text
[ ] User-facing target remains simple
[ ] Makefile recipe is thin
[ ] Script responsibility is singular
[ ] Shared helper reused
[ ] OS-specific behavior is centralized
[ ] Tool version comes from canonical source
[ ] Download is version-pinned
[ ] Download checksum is verified
[ ] System-level changes are not silently privileged
[ ] Script is idempotent where practical
[ ] Stop logic owns only its process
[ ] No arbitrary process killing
[ ] Docker/Compose logic is centralized
[ ] Errors are actionable
[ ] CI/non-interactive behavior works
[ ] ShellCheck passes
[ ] Compose validates when affected
[ ] README/RUNBOOK/help updated when public command changes
```

---

# 91. Definition of Done

A script/Make change is complete only when:

```text
[ ] Fresh clone path still works
[ ] Existing documented Make targets remain compatible
[ ] Relevant target was actually executed
[ ] Start reports readiness, not just process creation
[ ] Stop is safe/idempotent
[ ] Required toolchain is resolved automatically
[ ] No global developer toolchain mutation is required
[ ] Logs/diagnostics are available
[ ] Relevant tests passed
[ ] Actual results are reported
```

---

# 92. Final Principle

> **The user should remember commands, not setup procedures.**

The intended architecture is:

```text
make <intent>
     ↓
small script
     ↓
shared bootstrap/runtime helpers
     ↓
pinned tools + Compose + application
```

The public developer experience should remain:

```bash
make frontend-start
make frontend-stop

make backend-start
make backend-stop

make test
```

while all environment, version, download, runtime, and orchestration complexity stays behind those stable commands.

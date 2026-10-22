# Docker Guidelines

> **Status:** Mandatory engineering standard  
> **Location:** `docs/engineering/deployment/DOCKER_GUIDELINES.md`  
> **Applies to:** Dockerfiles, build contexts, multi-stage builds, image metadata, runtime users, Compose, health checks, volumes, networks, local development, CI image builds, release images, and container scanning  
> **Core principles:** minimal images, reproducible builds, non-root runtime, immutable artifacts, small contexts, explicit health/state, multi-architecture support, and no secrets in images

---

## 1. Objective

Docker must provide reproducible packaging, not a second hidden build system.

---

## 2. Mandatory Rules

1. Multi-stage builds are preferred.
2. Final images contain only runtime necessities.
3. Containers run as non-root where practical.
4. No secrets are copied into images.
5. Build context is minimized with `.dockerignore`.
6. Base images are pinned to supported versions; release builds should record immutable digest.
7. `latest` is not the deployment contract.
8. Every image includes OCI metadata/version/commit.
9. Runtime config is injected at runtime.
10. Build-time arguments do not carry secrets.
11. Health checks reflect application semantics where useful.
12. One container runs one primary application responsibility.
13. Logs go to stdout/stderr.
14. Application process receives signals directly.
15. Use exec-form `ENTRYPOINT`/`CMD`.
16. Avoid package-manager caches in final layer.
17. Build output is reproducible from source + lockfiles.
18. Images are vulnerability scanned.
19. Multi-arch builds are validated.
20. Compose service names/networks/volumes are stable and explicit.
21. Development bind mounts are not assumed in production.
22. Test Compose state is isolated.
23. Privileged containers are forbidden unless a documented feature absolutely requires it.
24. Docker socket mounting is forbidden by default.
25. Container images do not contain SSH server/debug daemons.
26. Writable filesystem paths are explicit.
27. Resource constraints are supported by deployment profile.
28. Container user/group IDs are documented where host-volume interaction matters.
29. Startup scripts are small and do not hide application logic.
30. Dockerfiles are reviewed like code.

---

## 3. Dockerfile Location

Recommended:

```text
backend/Dockerfile
frontend/Dockerfile
```

or a clearly centralized build location if the monorepo build system benefits.

Avoid dozens of nearly identical Dockerfiles.

---

## 4. Build Context

Use repository-root context only when necessary.

`.dockerignore` excludes:

```text
.git
node_modules
coverage
playwright-report
.runtime
.tools
.state
local secrets
IDE files
large docs/artifacts not needed in image
```

---

## 5. Backend Multi-Stage Build

Concept:

```dockerfile
FROM golang:<pinned> AS build
WORKDIR /src
# copy dependency metadata first
# download modules
# copy source
# build static/minimal binary

FROM <minimal-runtime>
COPY --from=build /out/perihellion /usr/local/bin/perihellion
USER <non-root>
ENTRYPOINT ["/usr/local/bin/perihellion"]
```

Do not ship compiler/toolchain in runtime image unless required.

---

## 6. Frontend Build

Static frontend may be built in a Node stage and served via the selected production web strategy.

Node/pnpm versions match repository source of truth.

Do not run Vite dev server as production frontend server.

---

## 7. Distroless / Minimal Runtime

Use minimal runtime bases where operational/debug requirements are satisfied.

Do not choose minimalism that makes TLS/cert/timezone behavior broken.

---

## 8. Base Image Governance

Approved base images have:

```text
known source
supported version
security updates
multi-arch support
license compatibility
```

---

## 9. Pinning

Development may pin major/minor for convenience only if CI/release resolves exact digest.

Release manifest records digest.

---

## 10. OCI Labels

Include:

```text
org.opencontainers.image.title
org.opencontainers.image.version
org.opencontainers.image.revision
org.opencontainers.image.source
org.opencontainers.image.licenses
```

where applicable.

---

## 11. Non-Root

Create/use explicit runtime user.

Do not rely on root because writing directories were not designed.

---

## 12. Filesystem

Prefer read-only root filesystem at runtime.

Explicit writable locations:

```text
/tmp
specific cache dir
```

Durable product data belongs outside application container filesystem.

---

## 13. Signals

Use exec form:

```dockerfile
ENTRYPOINT ["..."]
```

Avoid shell wrappers that swallow SIGTERM.

If wrapper needed, it must `exec` the application.

---

## 14. Init

Add an init process only if child reaping/signal requirements justify it.

Go services usually should not need a heavyweight supervisor.

---

## 15. Healthcheck

Image-level health check may exist, but orchestrator probes remain authoritative.

Health command must be lightweight.

---

## 16. Environment

Container accepts documented config via environment/files/secret references.

Do not bake environment-specific config into image.

---

## 17. Secrets in Build

Never:

```dockerfile
ARG NPM_TOKEN
RUN echo "$NPM_TOKEN" ...
```

without BuildKit secret mechanisms and explicit security review.

Prefer no secret-dependent public build when possible.

---

## 18. BuildKit Secrets

If private dependencies require credentials, use ephemeral secret mounts.

Ensure resulting layers/cache do not contain the secret.

---

## 19. Dependency Caching

Order Dockerfile layers to cache dependencies without sacrificing correctness.

Use locked dependency installs.

---

## 20. Go Build

Use reproducible flags/metadata.

Do not download runtime dependencies during application startup that belong in release build.

---

## 21. Frontend Build

Use:

```text
frozen pnpm lockfile
pinned Node
pinned pnpm
```

Build artifacts are immutable.

---

## 22. SBOM

Release image build emits SBOM.

Do not rely only on registry scanner.

---

## 23. Vulnerability Scan

Scan:

```text
base OS packages
application dependencies
binary/image
```

Critical findings require triage.

---

## 24. Image Size

Track size but do not optimize at the expense of security/correctness.

Unexpected large increases should be reviewed.

---

## 25. Compose Ownership

Developer Compose:

```text
scripts/docker/
```

Production/reference Compose may live under:

```text
deployments/compose/
```

Do not mix purposes invisibly.

---

## 26. Compose Base

Use a base Compose file plus explicit profile/override files.

Example:

```text
compose.yml
compose.dev.yml
compose.test.yml
```

---

## 27. Compose Profiles

Examples:

```text
lite
team
testing
observability
```

Do not start every optional dependency by default.

---

## 28. Compose Project Name

Use stable project naming.

Testing uses unique project name per run/worker when isolation requires it.

---

## 29. Compose Volumes

Name durable volumes intentionally.

Test volumes are disposable.

Never `docker volume prune` as a project cleanup command.

---

## 30. Compose Networks

Use one or more explicit project networks.

Do not expose internal ports to host unnecessarily.

---

## 31. Dependency Readiness

`depends_on` process startup is not sufficient.

Use health checks/readiness scripts before application/test begins.

---

## 32. PostgreSQL Compose

Reference DB has:

```text
pinned image
health check
named volume for dev
isolated volume for tests
explicit credentials through dev-safe config
```

Production credentials are not committed.

---

## 33. Valkey/NATS/MinIO/OpenBao

Only start when profile/capability requires them.

Each has health checks and pinned versions.

---

## 34. Observability Compose

Observability stack can be separate profile to keep Lite cheap.

---

## 35. Development Mounts

Bind mount source only in development.

Do not mount repository source into release containers.

---

## 36. Hot Reload

Frontend/backend hot reload is developer tooling.

Release images do not contain watchers/dev dependencies unless explicitly required.

---

## 37. Container Privileges

Forbidden by default:

```text
--privileged
host PID
host network
Docker socket
broad device mounts
```

Any exception requires security/design review.

---

## 38. Capabilities

Drop Linux capabilities by default.

Add only those required.

---

## 39. Seccomp/AppArmor

Use runtime defaults or tighter profiles where deployment supports them.

Do not disable security profiles casually.

---

## 40. Resource Limits

Compose reference production profiles may define CPU/memory guidance.

Application must behave gracefully under limits.

---

## 41. OOM

Container OOM is observable.

Do not configure silent unlimited memory in production guidance.

---

## 42. Restart Policy

Use suitable restart policy for deployment profile.

Do not use restart loops to hide deterministic startup/config errors.

---

## 43. Image Registry

Reference registry may be GHCR or configured release registry.

Enterprise/air-gap can mirror/import images.

---

## 44. Image Tags

Publish:

```text
v1.2.3
1.2
1
```

only if release policy supports aliases.

Immutable digest is authoritative.

---

## 45. Multi-Architecture

Build/test:

```text
linux/amd64
linux/arm64
```

where supported.

Do not publish architecture variant that has not passed smoke validation.

---

## 46. Container Test

Release CI should run image smoke:

```text
start image
wait readiness
call version/health
run critical minimal request
shutdown
```

---

## 47. Dockerfile Lint

Use a Dockerfile linter/security scanner where practical.

Rules should be CI-enforced.

---

## 48. Image Rebuild

Security rebuild with same application version should have a clear revision/build identifier.

Do not silently replace immutable release digest.

---

## 49. Compose Upgrade

For production/reference Compose:

```text
pull/import target images
preflight
backup
migration
recreate app containers
verify
```

Do not recreate DB volume.

---

## 50. Compose Reset

Destructive reset must be explicit.

Example:

```text
make infra-reset FORCE=1
```

No surprise data deletion.

---

## 51. Docker PR Checklist

```text
[ ] Multi-stage build
[ ] Final image minimal
[ ] Non-root
[ ] No secrets
[ ] Base/version pinned
[ ] OCI metadata
[ ] Exec-form process
[ ] Read-only/writable paths considered
[ ] .dockerignore correct
[ ] Compose profiles appropriate
[ ] Health checks meaningful
[ ] No broad privileges/socket mounts
[ ] Scan passes/triaged
[ ] Multi-arch impact considered
```

---

## 52. Definition of Done

```text
[ ] Image builds from clean checkout
[ ] Container starts
[ ] Health/readiness passes
[ ] Graceful SIGTERM works
[ ] Image scan passes/triaged
[ ] No secret in history/layers
[ ] Compose validates
[ ] Test profile isolated
[ ] Release digest recorded
```

---

## 53. Phase 0 Deliverables

```text
server Dockerfile
web Dockerfile
.dockerignore
non-root runtime
OCI metadata
developer Compose
test Compose
PostgreSQL health
optional profiles
image smoke test
container scan
amd64/arm64 build proof where possible
```

---

## 54. Agent Rules

Before Docker changes read:

```text
docs/AGENT.md
this file
deployment guideline
scripts guideline for developer Compose
security guideline
release guideline for published images
```

Do not put production orchestration logic into Dockerfile.

---

## 55. Final Mandatory Rules

1. Minimal multi-stage images.
2. Non-root runtime.
3. No secrets.
4. Immutable/pinned release artifacts.
5. Explicit config/state.
6. Correct signal handling.
7. Health/readiness.
8. No privileged/socket shortcuts.
9. Developer/test Compose isolated by purpose.
10. Scan and smoke every release image.

---

## 56. Guiding Principle

```text
Reproducible Build
      +
Minimal Runtime
      +
Non-Root
      +
Explicit State
      +
Immutable Digest
      =
Safe Container Artifact
```

> **A Perihellion container should be boring: deterministic to build, safe to run, easy to inspect, and disposable without losing business state.**

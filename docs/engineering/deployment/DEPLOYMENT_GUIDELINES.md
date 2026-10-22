# Deployment Guidelines

> **Status:** Mandatory engineering standard  
> **Location:** `docs/engineering/deployment/DEPLOYMENT_GUIDELINES.md`  
> **Applies to:** self-hosted installation, deployment profiles, runtime configuration, secret references, database migrations, Docker Compose, Kubernetes/Helm, health/readiness, scaling, backups, upgrades, rollback, air-gap packaging, and operational topology  
> **Supported profiles:** Lite, Team, Enterprise, Air-Gapped Enterprise  
> **Core principles:** smallest viable footprint, immutable releases, explicit state ownership, safe upgrades, provider-driven infrastructure, reproducibility, health-based orchestration, and no hidden production mutation

---

## 1. Objective

Perihellion should be easy to run on one machine and capable of scaling to enterprise environments without changing product semantics.

---

## 2. Deployment Philosophy

```text
Perihellion application stack
        ↓
stable product contracts
        ↓
provider-driven infrastructure
```

We are opinionated about application behavior.

We are replaceable where customers already own infrastructure.

---

## 3. Profiles

```text
Lite
Team
Enterprise
Air-Gapped Enterprise
```

Profiles enable capabilities; they are not separate products/codebases.

---

## 4. Lite

Minimum practical product:

```text
Perihellion Server
Perihellion Web
PostgreSQL
```

Optional dependencies are not mandatory unless a capability uses them.

---

## 5. Team

Reference stack may add:

```text
Valkey
NATS JetStream
MinIO/S3-compatible
Woodpecker CI
OpenBao
basic observability
```

These are reference providers, not mandatory vendor lock-in.

---

## 6. Enterprise

Typical:

```text
Kubernetes + Helm
external HA PostgreSQL
customer cache/events/object storage
enterprise OIDC
secret provider
full observability
GitOps
HA/DR
governance/policy
```

---

## 7. Air-Gapped Enterprise

Requires:

```text
offline OCI images
offline Helm/Compose bundle
checksums
signatures/trust metadata
SBOM
offline provider packages
offline docs
upgrade bundle
```

Installation must not silently reach public internet.

---

## 8. Mandatory Rules

1. Production deploys immutable, versioned artifacts.
2. `latest` is not a production deployment contract.
3. Runtime config is externalized.
4. Secrets are referenced/externalized.
5. Stateful systems have explicit ownership and backup strategy.
6. Database migration is a controlled release/deployment step.
7. Readiness gates traffic.
8. Liveness does not depend on every remote optional dependency.
9. Graceful shutdown is supported.
10. Production resource requests/limits are defined.
11. Containers run non-root where practical.
12. Images contain no secrets.
13. Development Compose is not automatically production Compose.
14. One product release ships compatible component versions.
15. Upgrade compatibility is documented and tested.
16. Rollback limitations are documented.
17. Backups are tested through restore.
18. External infrastructure integrates through provider contracts.
19. Production changes are reproducible from versioned config/manifests.
20. Environment-specific secret files are never committed.
21. Migration execution does not race across replicas.
22. Horizontal scaling cannot depend on process-local durable state.
23. Durable files use object storage/volumes explicitly.
24. Air-gap artifacts are verifiable.
25. Health is observable.
26. Small deployments do not require enterprise HA components.
27. Upgrade/preflight jobs are idempotent where practical.
28. Destructive operational actions require explicit intent.
29. Unsupported deployment profiles are not advertised.
30. Every supported profile has a validation path.

---

## 9. Core Deployables

Initial deployable units:

```text
server
web
```

Distributed tools:

```text
CLI
Local Development Agent
```

Workers can remain in the backend process until independent scale/lifecycle is justified.

---

## 10. Artifact Identity

Every artifact exposes:

```text
product version
commit SHA
build ID
architecture
digest/checksum
```

---

## 11. Configuration

Document one precedence model.

Example:

```text
explicit runtime secret/config
        ↓
environment variables
        ↓
config file
        ↓
safe defaults
```

Do not keep conflicting defaults across templates.

---

## 12. Secrets

Use approved secret stores/mechanisms.

Prefer:

```text
secret reference
workload identity
short-lived credential
```

over copying raw secret values.

---

## 13. PostgreSQL

PostgreSQL is required for core transactional state.

Production must define:

```text
backup
restore
HA as profile requires
connection budget
storage
monitoring
```

---

## 14. Migrations

Production flow:

```text
preflight
  ↓
compatible migration
  ↓
application rollout
  ↓
verification
```

Schema contraction occurs only after old application versions no longer depend on old shape.

---

## 15. Migration Ownership

One controlled actor/job migrates.

Do not let every replica race migration.

---

## 16. Startup

Recommended sequence:

```text
load configuration
initialize telemetry
connect critical dependencies
verify schema compatibility
initialize registries
start listeners
pass readiness
```

---

## 17. Graceful Shutdown

Sequence:

```text
mark unready
stop accepting new work
drain requests
stop/hand off workers
flush bounded telemetry
close provider clients/pools
exit
```

---

## 18. Stateless Application Tier

API/server replicas should be stateless except for explicit external state.

Do not store durable workflow state only in memory.

---

## 19. Durable Execution

Mission Runtime checkpoints/state are persisted so browser/process restarts do not lose missions.

---

## 20. Object Storage

Large artifacts/evidence/exports belong in S3-compatible storage when enabled.

Local filesystem mode, if supported for Lite, must document backup/HA limitations.

---

## 21. Cache

Cache is an optimization/capability, not source of truth unless a domain explicitly defines durable cache semantics.

Valkey is the reference provider.

---

## 22. Events

NATS JetStream is the Team/reference event provider.

Other implementations must satisfy the event capability contract.

---

## 23. Observability

Reference:

```text
OpenTelemetry
Prometheus
Loki
Tempo
Grafana
```

Customers may bind alternatives.

---

## 24. Health Endpoints

Expose:

```text
/live
/ready
```

or repository-standard equivalents.

---

## 25. Version Endpoint

Expose safe build/version metadata for operators/support.

Do not expose configuration or secret values.

---

## 26. Resources

Production defines:

```text
CPU
memory
ephemeral storage
worker concurrency
DB connection budget
```

---

## 27. Database Connection Budget

Replica scaling must respect global DB connection capacity.

Do not multiply large default pools across many replicas.

---

## 28. Autoscaling

Scale on meaningful signals:

```text
CPU
memory
request concurrency
queue backlog
worker utilization
```

Avoid scaling purely on a noisy single metric.

---

## 29. Rolling Upgrades

Old and new application versions may overlap.

Contracts/schema/config must support the documented overlap window.

---

## 30. Rollback

Rollback is safe only while:

```text
schema
data
config
provider protocol
```

remain compatible.

Use forward recovery for irreversible data transformations.

---

## 31. Upgrade Preflight

Verify:

```text
current version
target version
supported path
DB backup state
schema compatibility
required provider SDK/protocol
resources/disk
air-gap package completeness
```

---

## 32. Backup

Backup scope may include:

```text
PostgreSQL
object storage
critical configuration
provider binding metadata
encryption metadata according to design
```

---

## 33. Restore

Periodically restore to isolation and verify actual product behavior.

A successful archive command alone does not prove recoverability.

---

## 34. Disaster Recovery

Enterprise profiles define:

```text
RPO
RTO
restore/failover runbook
off-site copy
regional/site strategy
```

---

## 35. TLS

Use TLS for user-facing and untrusted-network communication.

Certificate lifecycle belongs to ingress/platform integration.

---

## 36. Ingress/Gateway

Perihellion integrates with existing ingress/gateway.

It does not require replacing an organization's gateway.

---

## 37. Network Exposure

Expose only necessary ports.

Database/cache/event/provider administrative ports remain private.

---

## 38. Network Policy

Enterprise deployments should constrain:

```text
ingress
east/west
egress
```

according to actual dependencies.

---

## 39. Air-Gap Registry

Support import into private OCI registry or offline runtime.

Every required image is listed by immutable digest.

---

## 40. Architecture Targets

Support:

```text
amd64
arm64
```

where all selected dependencies support them.

Publish architecture support clearly.

---

## 41. Persistent Volumes

Only explicit stateful components use persistent volumes.

Application logs should go to stdout/telemetry, not durable local files.

---

## 42. Read-Only Filesystem

Application containers should support read-only root filesystem where practical.

Writable temp paths are explicit.

---

## 43. Init Containers / Init Jobs

Use for deterministic setup only.

Do not hide domain business logic in deployment scripts.

---

## 44. Capability Enablement

Deployment enables/binds capabilities.

No compile-time forks for each customer infrastructure choice.

---

## 45. Optional Dependency Failure

Unavailable optional provider should degrade relevant capability, not necessarily the entire platform.

---

## 46. Installation UX

Target product installation:

```text
release bundle
   ↓
one documented install workflow
   ↓
bootstrap/configure
   ↓
browser
```

Developer `make dev-start` remains a development path, not the enterprise production installer.

---

## 47. Docker Compose

Compose supports:

```text
development
Lite
small Team
reference deployments
```

Detailed rules are in `DOCKER_GUIDELINES.md`.

---

## 48. Kubernetes

Kubernetes + Helm is the enterprise/reference orchestrated deployment.

Detailed rules are in `KUBERNETES_GUIDELINES.md`.

---

## 49. Helm Dependency Strategy

Allow:

```text
bundled/reference dependency
or external provider endpoint
```

Do not force users to deploy duplicate infrastructure they already operate.

---

## 50. GitOps

Flux is the reference GitOps integration where desired.

GitOps is not mandatory for all profiles.

---

## 51. Upgrade Compatibility Matrix

Document:

```text
from version
to version
DB requirements
configuration changes
provider SDK/protocol requirements
Kubernetes/runtime requirements
breaking notes
```

---

## 52. Configuration Deprecation

Deprecated keys have:

```text
replacement
warning
support window
removal version
```

Do not silently reinterpret old config.

---

## 53. Pre-Deploy Validation

Validate:

```text
rendered manifests
images/digests
secret references
DB connectivity
migration compatibility
resource sanity
provider endpoint prerequisites
```

---

## 54. Post-Deploy Validation

Verify:

```text
readiness
version
critical API smoke
login/basic navigation
migration status
provider health summary
```

---

## 55. Progressive Delivery

Can be used when infrastructure supports it.

Core product must not depend on one specific progressive-delivery tool.

---

## 56. Drift

Desired deployment config is version-controlled.

Manual mutations are discouraged/detected.

---

## 57. Immutable Infrastructure

Replace instances/images instead of manually modifying live containers/VMs.

---

## 58. Production Shell Access

Interactive shell is exceptional.

Prefer audited Workbench/operational capabilities.

---

## 59. Deployment PR Checklist

```text
[ ] Deployment profile identified
[ ] Images pinned
[ ] Config external
[ ] Secrets external
[ ] Migration safe
[ ] Probes correct
[ ] Graceful shutdown tested
[ ] Resource budget considered
[ ] Durable state explicit
[ ] Backup/restore impact considered
[ ] Upgrade/rollback documented
[ ] Air-gap impact considered
[ ] Observability present
```

---

## 60. Definition of Done

```text
[ ] Fresh install works
[ ] Supported upgrade path works
[ ] Readiness/liveness verified
[ ] Migration verified
[ ] Critical smoke passes
[ ] Recovery path understood
[ ] Artifacts immutable/verifiable
[ ] No hidden internet dependency for air-gap
```

---

## 61. Phase 0 Deliverables

```text
Lite/reference Compose deployment
PostgreSQL
server/web images
health/readiness
configuration baseline
secret-reference baseline
migration execution path
version metadata
clean clone -> development start proof
Helm skeleton
deployment smoke test
```

---

## 62. Agent Rules

Before deployment changes read:

```text
docs/AGENT.md
this file
Docker or Kubernetes guideline
database guideline when schema changes
security guideline
observability guideline
release guideline
```

Do not conflate development scripts with production deployment.

---

## 63. Final Mandatory Rules

1. Immutable versioned artifacts.
2. Smallest viable profile first.
3. Optional capabilities are optional.
4. Config/secrets externalized.
5. Migrations are controlled.
6. Health gates rollout.
7. Graceful shutdown is mandatory.
8. State/backup ownership explicit.
9. Rolling compatibility is considered.
10. Air-gap has no hidden internet dependency.
11. External infrastructure remains provider-driven.
12. Upgrade/recovery is tested, not assumed.

---

## 64. Guiding Principle

```text
One Product Release
      +
Immutable Artifacts
      +
Explicit State
      +
Safe Upgrade
      +
Provider-Driven Infrastructure
      =
Reliable Self-Hosting
```

> **Perihellion should be inexpensive to run small, predictable to operate, and capable of using infrastructure the customer already owns.**

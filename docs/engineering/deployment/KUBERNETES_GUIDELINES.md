# Kubernetes Guidelines

> **Status:** Mandatory enterprise deployment standard  
> **Location:** `docs/engineering/deployment/KUBERNETES_GUIDELINES.md`  
> **Applies to:** Kubernetes manifests, Helm charts, namespaces, service accounts, RBAC, security contexts, probes, resources, autoscaling, PDBs, topology, network policies, Jobs, migrations, secrets, storage, ingress, GitOps, upgrades, and air-gap clusters  
> **Primary packaging:** Helm  
> **Reference GitOps:** Flux when GitOps is desired  
> **Core principles:** least privilege, declarative desired state, immutable images, safe rollout, explicit resources, workload identity, health-based routing, state externalization, and portable cluster integration

---

## 1. Objective

Kubernetes is an enterprise deployment target, not a dependency embedded into Perihellion business logic.

---

## 2. Mandatory Rules

1. Helm is the reference Kubernetes package.
2. Images are immutable/pinned.
3. Pods run non-root where practical.
4. Use dedicated ServiceAccounts.
5. No cluster-admin for normal application workloads.
6. Security context is explicit.
7. Readiness/liveness/startup probes are meaningful.
8. Resource requests are defined.
9. Resource limits are defined where appropriate.
10. Horizontal scaling respects DB/worker budgets.
11. Pod-local disk is not durable business state.
12. Secrets come from Kubernetes/external secret provider, never ConfigMap.
13. Migrations run in a controlled Job/hook/workflow.
14. App replicas do not race migrations.
15. Graceful termination honors `terminationGracePeriodSeconds`.
16. PodDisruptionBudget is used for HA profiles where meaningful.
17. Topology spreading/anti-affinity is used for HA.
18. NetworkPolicy is supported/recommended for enterprise.
19. Ingress is configurable and does not force one controller.
20. TLS is supported through standard ingress/gateway integration.
21. Persistent volumes are used only by components that own persistent state.
22. Stateful infrastructure can be external.
23. Helm values are typed/documented and avoid secret defaults.
24. Helm templates do not contain business logic.
25. CRDs are avoided unless a capability truly requires them.
26. Kubernetes API is not called directly by generic business code.
27. Cluster integration happens through runtime/provisioning providers.
28. RBAC follows least privilege.
29. GitOps is optional.
30. Upgrade/rollback behavior is tested.

---

## 3. Namespace

Default installation should support a dedicated namespace:

```text
perihellion
```

Allow configurable namespace.

Do not assume cluster-wide ownership.

---

## 4. Release Name

Helm release name is configurable.

Resource names derive predictably without exceeding Kubernetes limits.

---

## 5. Helm Structure

Recommended:

```text
deployments/helm/perihellion/
├── Chart.yaml
├── values.yaml
├── values.schema.json
├── templates/
├── crds/                # only if truly needed
├── README.md
└── tests/
```

---

## 6. Values Schema

Use `values.schema.json` for validation.

Fail early on invalid values.

---

## 7. Values Philosophy

Values configure deployment, not application business semantics.

Avoid huge nested vendor-specific values when provider configuration belongs in Perihellion/provider contracts.

---

## 8. Images

Values specify:

```text
repository
tag
digest
pullPolicy
pullSecrets
```

Release deployment should prefer digest-verifiable references.

---

## 9. Service Account

Create dedicated ServiceAccount by default.

Do not use default ServiceAccount for privileged operations.

---

## 10. RBAC

Grant only resources/verbs actually required.

If Perihellion provider needs broad Kubernetes access, scope it by namespace/credential/binding and document why.

---

## 11. Workload Identity

Prefer cloud/workload identity federation over static cloud credentials where available.

---

## 12. Pod Security Context

Configure:

```text
runAsNonRoot
runAsUser/group where appropriate
fsGroup only if needed
seccompProfile
```

---

## 13. Container Security Context

Prefer:

```text
allowPrivilegeEscalation: false
readOnlyRootFilesystem: true
capabilities:
  drop: ["ALL"]
```

Add exceptions only with documented requirement.

---

## 14. Privileged

Privileged pods are forbidden by default.

Local Development Agent may run outside cluster where host access is required instead of granting core control plane privilege.

---

## 15. Probes

Use:

```text
startupProbe
readinessProbe
livenessProbe
```

when startup can be slow enough to require distinction.

---

## 16. Readiness

Readiness failure removes pod from service traffic.

Do not use a static TCP check if app-level readiness exists.

---

## 17. Liveness

Liveness detects wedged process.

Do not couple to optional external provider availability.

---

## 18. Graceful Shutdown

On SIGTERM:

```text
mark unready
drain
finish/handoff work
close
```

Configure enough termination grace time.

---

## 19. PreStop

Use only if it adds necessary drain behavior.

Do not add arbitrary sleeps as rollout synchronization.

---

## 20. Requests

Every production workload defines CPU/memory requests.

Scheduler needs realistic values.

---

## 21. Limits

Memory limits are strongly recommended.

CPU limits should be chosen with workload characteristics in mind.

Do not copy arbitrary defaults across all components.

---

## 22. HPA

HorizontalPodAutoscaler may use:

```text
CPU
memory
custom metrics
queue backlog
```

Scale decisions must respect external resource budgets.

---

## 23. Worker Scaling

Mission workers may scale separately in future.

If still in modular monolith, ensure scaling does not duplicate singleton tasks/migrations.

---

## 24. PDB

HA profile uses PodDisruptionBudget so voluntary disruption does not remove all replicas.

Lite single-replica deployments should not create impossible PDB constraints.

---

## 25. Replicas

Values profile controls replicas.

Do not claim HA with one replica.

---

## 26. Topology Spread

Enterprise HA spreads across nodes/zones where infrastructure supports it.

---

## 27. Affinity

Use affinity/anti-affinity intentionally.

Avoid hard scheduling rules that make small clusters unschedulable.

---

## 28. Service

Use ClusterIP internally by default.

External access through ingress/gateway/load balancer integration.

---

## 29. Ingress

Support generic Kubernetes Ingress/Gateway configuration.

Do not hardcode NGINX/Traefik-specific behavior into core chart unless under optional annotations/values.

---

## 30. TLS

Allow:

```text
existing TLS secret
cert-manager integration
external ingress TLS
```

Do not generate insecure self-signed production certs automatically.

---

## 31. NetworkPolicy

Enterprise chart should support policies for:

```text
web/server ingress
DB egress
provider endpoints
observability
DNS
```

Default can be opt-in until dependency graph is complete, but target is least egress.

---

## 32. ConfigMap

Use for non-secret application configuration.

Large dynamic business configuration belongs in application config service/database, not giant ConfigMaps.

---

## 33. Secrets

Kubernetes Secret is a delivery mechanism, not automatically a secure secret manager.

External Secrets/OpenBao/cloud providers are preferred for enterprise rotation.

---

## 34. Secret Mount vs Env

Choose based on library/runtime.

Mounted files can support rotation better than process env for some secrets.

Application design should support secret provider abstractions.

---

## 35. PostgreSQL

Chart may provide a reference PostgreSQL for Lite/Team if desired.

Enterprise should support external PostgreSQL cleanly.

Do not force in-cluster DB.

---

## 36. Stateful Dependencies

Valkey/NATS/MinIO/OpenBao may be:

```text
bundled reference
external
disabled
```

according to profile/capability.

---

## 37. PersistentVolume

Use StorageClass-configurable PVC.

Document size/growth/backup.

---

## 38. Database Migration Job

Migration Job:

```text
uses target release migration image/binary
uses migration credentials
runs once
reports status
```

Application Deployment waits on release workflow, not ad hoc race.

---

## 39. Helm Hooks

Use hooks carefully.

Failed hooks can make recovery confusing.

An explicit release Job controlled by CI/GitOps may be clearer.

---

## 40. Init Containers

Do not run DB migrations in every pod init container.

Use init only for pod-local prerequisites.

---

## 41. Config Compatibility

Chart validates required values.

Application validates semantic configuration at startup.

---

## 42. Observability

Support annotations/ServiceMonitor/OpenTelemetry collector integration without requiring one monitoring operator.

---

## 43. Prometheus Operator

ServiceMonitor/PodMonitor support can be optional.

Core metrics remain standard HTTP/OpenMetrics compatible.

---

## 44. Logging

Containers log to stdout/stderr.

Cluster logging provider collects them.

Do not require hostPath log volumes.

---

## 45. Tracing

OTLP endpoint is configurable.

TLS/auth supported.

---

## 46. Image Pull

Support private registry credentials.

Air-gap deployment uses internal registry/digest.

---

## 47. Registry Mirrors

No application assumption about public registry hostname after install bundle is rewritten/configured.

---

## 48. PriorityClass

Use only if cluster operators choose it.

Do not ship aggressive global priority by default.

---

## 49. LimitRange/ResourceQuota

Perihellion chart should coexist with namespace policies.

Do not assume unlimited cluster permissions/resources.

---

## 50. Node Selectors/Tolerations

Expose generic values.

Do not hardcode cloud vendor labels.

---

## 51. Service Mesh

Perihellion should work with/without service mesh.

Do not require mesh-specific sidecars in core chart.

---

## 52. GitOps

Flux can reconcile HelmRelease/Kustomization.

The chart remains usable without Flux.

---

## 53. Argo CD

Users can deploy with Argo CD if standard Helm behavior is preserved.

No need for product-specific dependency.

---

## 54. Rollout

Default Kubernetes Deployment rolling strategy.

Use readiness to gate.

Advanced canary/blue-green tools are optional.

---

## 55. Revision History

Set reasonable `revisionHistoryLimit`.

Do not rely on deployment revision as full product rollback plan.

---

## 56. Horizontal Scaling and Jobs

Leader/singleton tasks require explicit coordination.

Do not assume one pod when replicas > 1.

---

## 57. Leader Election

Introduce only for workloads that need singleton ownership.

Prefer durable queue/DB claiming when it naturally fits.

---

## 58. CronJobs

Use for deployment/operational jobs only when schedule belongs to cluster ops.

Product scheduling belongs to Mission Runtime/application where appropriate.

---

## 59. Kubernetes Provider

When Perihellion manages customer workloads on Kubernetes, that functionality belongs behind a provider capability.

The control-plane deployment chart does not become the runtime provider implementation.

---

## 60. Cluster Credentials

Kubernetes provider credentials are scoped and secret-managed.

Never assume the control-plane ServiceAccount can administer arbitrary customer clusters.

---

## 61. Multi-Cluster

Model as provider instances/bindings.

Do not hardcode one cluster into business modules.

---

## 62. Air-Gap

Chart dependencies/images are vendored/listed.

No remote chart dependency resolution during offline install unless bundle contains it.

---

## 63. Helm Tests

Add Helm tests/smoke where useful.

Render tests validate:

```text
Lite
Team
external DB
custom ingress
air-gap registry overrides
```

---

## 64. Static Validation

CI should run:

```text
helm lint
helm template
values schema validation
kubeconform/kubeval equivalent
security policy scan
```

---

## 65. Kubernetes Security Scan

Validate:

```text
non-root
capabilities
privilege
host mounts
service account
secrets
resource limits
```

---

## 66. Upgrade Test

Release gate should install previous supported chart/version, create state, then upgrade to target and run smoke/regression.

---

## 67. Rollback Test

Where rollback is supported by schema compatibility, test Helm/application rollback.

---

## 68. Chart Version

Helm chart version and appVersion follow release policy.

Chart changes not tied to app code may still require chart versioning.

---

## 69. CRDs

Avoid CRDs for core product configuration unless Kubernetes-native extension semantics clearly justify them.

A database/API is often better for product state.

---

## 70. Operator Pattern

Do not build a Kubernetes Operator in Phase 0.

Add one only when reconciliation of Kubernetes-native resources is a real product requirement.

---

## 71. Kubernetes PR Checklist

```text
[ ] Images pinned
[ ] ServiceAccount/RBAC least privilege
[ ] Non-root/securityContext
[ ] Probes meaningful
[ ] Resources defined
[ ] Graceful termination
[ ] HA/PDB/topology profile appropriate
[ ] Secrets external
[ ] Migration not racing
[ ] External dependencies configurable
[ ] Ingress vendor-neutral
[ ] Network exposure minimal
[ ] Helm lint/render/schema pass
[ ] Upgrade impact reviewed
```

---

## 72. Definition of Done

```text
[ ] Helm install works on clean cluster
[ ] Upgrade path works
[ ] Pods become ready
[ ] SIGTERM drains correctly
[ ] Security checks pass
[ ] Lite profile can schedule on small cluster
[ ] Enterprise profile supports external dependencies
[ ] Air-gap rendering does not reference public internet unexpectedly
```

---

## 73. Phase 0 Deliverables

```text
Helm chart skeleton
server/web Deployment
Service
ConfigMap/Secret reference model
ServiceAccount
securityContext
probes
resources
external PostgreSQL values
reference PostgreSQL option if chosen
migration Job pattern
helm lint/template CI
```

---

## 74. Agent Rules

Before Kubernetes changes read:

```text
docs/AGENT.md
this file
deployment guideline
Docker guideline
security guideline
observability guideline
release guideline
```

Do not add cluster-wide permissions for convenience.

---

## 75. Final Mandatory Rules

1. Helm is reference packaging.
2. Non-root/least privilege.
3. Dedicated ServiceAccount.
4. Explicit probes/resources.
5. Controlled migrations.
6. External state.
7. Provider-driven customer cluster integrations.
8. Ingress/observability remain vendor-neutral.
9. HA behavior is profile-aware.
10. Air-gap is render/install verifiable.
11. Kubernetes config is declarative.
12. Upgrade is tested.

---

## 76. Guiding Principle

```text
Declarative State
      +
Least Privilege
      +
Health-Based Rollout
      +
External Durable State
      +
Portable Integrations
      =
Enterprise Kubernetes Deployment
```

> **Kubernetes should operate Perihellion; it should not leak into Perihellion's business architecture.**

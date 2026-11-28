# Multi-project refactor

Baseline: phase-2.1-identity-account, bbd54e2a6f8217119bb72d7eccdb5dc1e2504e13. Static inventory only; build and test baseline pending.

## Target
- frontend/libs/{base,styles,components}: independently packaged and tested.
- frontend/projects/{school,internal-tool}: independent manifests, lockfiles, builds.
- backend/libs/{core,contracts}: versioned Go modules; contracts contain public DTOs, client interfaces and events, not implementations.
- backend/projects/{platform,school,internal-tool}: independent Go modules, launchers, migrations and services.
- backend/deployment: shared deployment orchestration, no root language tooling.
- Root retains repository metadata and .github/workflows only alongside frontend/backend.
- Proposed canonical namespace github.com/groovehub/groovehub and @groovehub npm scope; verify ownership before changing imports or releasing.

## Existing ownership
Move Foundation, Identity, Infrastructure, Integration and Audit together into Platform. Preserve existing public routes, migrations, generated code conventions, operator commands and in-process clients. Communication is future work. School owns student, academic, attendance, CRM, finance, task and HR domains; Accounting deferred. Internal Tool owns separate operator identity, approvals, API templates, execution and environment configuration. Do not scaffold unused services merely to fill directories.

Existing backend/libs/schoolgateway is application-specific: retain school session policy in Platform and expose only necessary verified contracts. Extract generic messaging from backend/libs/infrastructure into Core after dependency audit. Preserve Audit-owned producer facade and transactional semantics.

## Migration gates
1. Capture build/test baseline, route inventory, migration history, imports and aliases. Do not claim baseline green without execution.
2. Move source and configuration with atomic path updates. Preserve existing behavior and old namespace until target ownership verified. Move root scripts/docs to their owners; preserve cross-project orchestration under backend/deployment. Remove stray root artifacts only after inspection.
3. Decouple annotation generation from writable Core source. Ship generated Core source required by consumers; version generator CLI; application-local overlays. Validate isolated builds with GOWORK=off and no local replace.
4. Compile frontend libraries to JS/declarations/CSS, declare peer dependencies, preserve public exports, test packed packages in isolated consumers. Publish Base and Styles before Components.
5. Verify organization and registry permissions, transfer/rename only with authorization, release Go submodule tags and frontend packages. Pin versions, remove workspace and local path dependencies.
6. Stabilize Platform and School frontend; preserve auth, RBAC, invitation, audit, integration and foundation regression coverage.
7. Introduce typed HTTP transport alongside existing in-process clients; restricted routes require workload identity, verified actor context, tenant authorization, deadlines and idempotency. Preserve session revocation and CSRF. Internal sessions and cookies must be distinct.
8. Each project owns database credentials and local transactional outbox; no cross-project SQL or remote outbox writes. Share NATS with scoped subjects, durable consumers, idempotency, retries and replay. Maintain existing event compatibility.
9. Add School and Internal Tool launchers only after shared contracts stable. Independent builds and deployment; one VPS, Platform 8080, School 8081, Internal 8082; dev frontends 4173/4174; private backend listeners. Resolve existing Swagger 8081 collision. PostgreSQL one instance with separate logical databases; optional audit DSN.
10. CI isolated source builds, package compatibility, security regressions, migrations, load tests, backup restore and rollback. Update all engineering guidelines.

## Acceptance
No root package.json, node_modules, pnpm lock/workspace, go.work or go.work.sum. No application imports another project implementation. No registry secrets in artifacts. Preserve existing APIs, migration history and event delivery. Build and test every changed project before merge. No production release until namespace, credentials and isolated builds verified.

## Status
Branch initialized; documentation only. No source migration, package release, namespace transfer, or executed build/test yet.

## Verified CI baseline (2026-11-27)
- Frontend run 45955769607: install, shared dependency inventory, guideline audit, Base/Components/application typecheck, Jest and ESLint passed. Architecture check failed; build, prerender and browser regression skipped. Formatting remediation artifact generated. Exact failure logs unavailable through connector; inspect run before modifying architecture rules.
- Backend run 45955765539: PostgreSQL integration/RBAC job passed; annotation and macOS coverage jobs pending at inspection.
- Do not disable architecture checks to make migration appear green. No physical migration yet.

## Mandatory API boundary (supersedes earlier direct Platform ingress proposal)
All Platform business APIs are restricted, accessible only by authenticated and authorized workloads. No browser-facing public or protected Platform routes. Health endpoints are private operational exceptions, not business APIs. School and Internal Tool expose their own protected use-case wrappers; application-owned Identity edge exposes public login/recovery/invitation endpoints and protected session/account operations as appropriate. Platform Identity remains a restricted backend capability. School/Identity wrappers enforce session, CSRF, origin, tenant, resource and capability policies before calling Platform; Platform independently verifies workload scope and delegated actor context. No generic unrestricted proxy or browser access to Platform. Provider callbacks terminate at a dedicated verified ingress adapter, not a public Platform business route. Nginx must not forward external traffic directly to Platform. Preserve existing browser URLs through application-owned wrappers where possible. Migration must inventory every current public/protected Platform route, define its wrapper and restricted counterpart, and switch routing atomically with regression tests. Tests must prove browser cookies cannot authenticate restricted routes, forged actor headers fail, cross-tenant access fails, and wrappers preserve response contracts. No direct Platform exposure during transition in production; current branch has not yet implemented these changes.

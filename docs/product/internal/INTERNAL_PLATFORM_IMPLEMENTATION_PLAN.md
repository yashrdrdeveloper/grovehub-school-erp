# GrooveHub Internal Platform — Implementation Plan

Status: PROPOSED — not implemented. Owner: platform engineering. Audience: engineers and coding agents.
Baseline: phase-2.1-identity-account @ 9fa83d6ca850069d9243e67a1a4d4355adce483b, inspected 26 November 2026.
Companion: [School Identity RBAC implementation plan](../users/RBAC_IMPLEMENTATION_PLAN.md).
This plan coordinates cross-cutting work. It is not a replacement for service-owned OpenAPI, migrations, engineering architecture or product execution guides.

## 1. Objective and non-negotiable decisions

Build a separate internal engineering/operations control plane for GrooveHub employees. The first release covers internal Identity/RBAC, API registry, typed templates, execution, saved requests, universal approvals and audit. Later it must support deployments, releases, monitoring, incidents, Kubernetes/orchestration, jobs and database operations without replacing the approval or identity foundation.

- Extract reusable frontend Base and Components out of school frontend; extract neutral Styles alongside them to preserve design tokens. Extract Go Core only after proving it has no business dependencies.
- Use **internal/backend** and **internal/frontend** (not backend-internal/frontend-internal). School backend/ and frontend/ retain their paths.
- Independent internal PostgreSQL URL, database, credentials, migrations, sessions, users, permissions, approval data, runtime and deployment. No school database reads/writes by internal application.
- Templates are independent of environment and school. Execution selects environment and, only for SCHOOL scope, school. GLOBAL execution has no tenant header.
- Registration approval and execution approval are independent. Every execution receives its own approval request/decision even when auto-approved.
- If requester is an eligible approver and the applicable policy explicitly allows it, auto-approve the request. Sensitive/privilege-escalating operations require independent approvers. No universal self-approval.
- Use authenticated workload SYSTEM actor; -1 is an audit/display convention, never a UUID or an authorization credential.
- Keep existing school blue design, internal orange theme; light/dark/system modes, existing mobile list and full-screen Menu behavior. Do not recreate UI primitives from scratch.

## 2. Read before implementing

Mandatory source of truth:
- [Agent contract](../../AGENT.md), [engineering index](../../engineering/README.md).
- [Backend master](../../engineering/backend/backend.md), [libraries](../../engineering/backend/libraries.md), [annotations](../../engineering/backend/annotation.md), [API](../../engineering/backend/api.md), [patterns](../../engineering/backend/patterns.md), [persistence](../../engineering/backend/persistence.md), [events](../../engineering/backend/events.md), [gateway](../../engineering/backend/gateway.md), [providers](../../engineering/backend/providers.md), [searchDetails](../../engineering/backend/searchDetails.md), [hybrid pagination](../../engineering/backend/HYBRID_SEARCH_PAGINATION.md), [testing](../../engineering/backend/testing.md).
- [Frontend guidelines](../../engineering/frontend/FRONTEND_CODE_GUIDELINES.md), [universal table](../../engineering/frontend/UNIVERSAL_TABLE.md), [compact table](../../engineering/frontend/COMPACT_TABLE_LAYOUT.md), [workspace theme](../../engineering/frontend/WORKSPACE_THEME.md), [auth UI](../../engineering/frontend/AUTH_PAGE_DESIGN_SYSTEM.md).
- [API contracts](../../engineering/api/API_CONTRACT_GUIDELINES.md), [response envelope](../../engineering/api/API_RESPONSE_ENVELOPE_GUIDELINES.md), [database](../../engineering/database/DATABASE_GUIDELINES.md), [security](../../engineering/security/SECURITY_GUIDELINES.md), [test standards](../../engineering/testing/TESTING_GUIDELINES.md), [backend test layout](../../engineering/testing/BACKEND_UNIT_TEST_STRUCTURE.md), [regression tests](../../engineering/testing/REGRESSION_TESTING.md).
- [Scripts](../../engineering/scripts/SCRIPT_AND_MAKEFILE_GUIDELINES.md), [local development](../../engineering/scripts/GROVEHUB_LOCAL_DEVELOPMENT.md), [CI/CD](../../engineering/ci/CI_CD_GUIDELINES.md), [deployment](../../engineering/deployment/DEPLOYMENT_GUIDELINES.md), [observability](../../engineering/observability/GROVEHUB_OBSERVABILITY.md), [documentation](../../engineering/documentation/DOCUMENTATION_GUIDELINES.md).
- School Identity [architecture](../../../backend/services/identity/docs/engineering/architecture.md), [API workflows](../../../backend/services/identity/docs/product/api_workflows.md), [temporary authorization warning](../users/ACCOUNT_AUTHORIZATION.md).

Precedence: deployed contracts and applied migrations > canonical standards > owner-local docs > proposed design. Reopen relevant full standards and actual owners at each phase. Some guidelines contain legacy paths/status; update path references alongside real moves, never relax gates to make a move pass.

## 3. Current-state findings (verified sample, not complete dependency audit)

| Existing item | Finding / implication |
|---|---|
| frontend/packages/{base,components,styles,business,domain,shell,assets} | Existing app imports use @Base, @Components, @Styles etc. Vite and tsconfig resolve app-local paths. Need package boundaries, not a blind move. |
| Components/SchoolMonogram.tsx | Imports @Assets/brand RDR artwork. Move into school owner before library extraction. Audit other school-specific exports and tests. |
| useWorkspaceTheme | Writes groovehub.theme and owns document theme; inject application preference namespace and policy, keep one active theme owner. |
| frontend/package.json | pnpm 10.14.0, React 19, Vite 7, TanStack Query/Table, Radix, Jest and numerous architecture checks. Preserve pinned dependency versions during move. |
| frontend/app/styles.css | Legacy global styles include shell/branding and responsive selectors. Extract actual generic consumers into owner CSS Modules; do not publish all global CSS. |
| backend/go.mod | One Go module grovehub; core imports grovehub/libs/core/...; no standalone Core go.mod. |
| backend/libs/core | 121 non-tests-directory Go files at baseline; mostly technical capabilities by package inventory, but not individually audited. |
| core/gateway/gateway.go | Principal has AccountID, PersonGID, TenantID, SessionID, Admin and hardcoded school cookie names. Requires contract review/config injection before extraction. |
| annotation generators + scripts/backend/generate-annotations.sh | Assume backend/services and backend/libs/core locations, including .generated overlay. Generalize target root, module path and library generation, then test both apps. |
| school gateway external operation | Protected uses cookie verification; restricted uses workload verification. No blanket workload access to protected routes; introduce explicit authorized bridge. |
| school Identity | Existing six role codes and primary/secondary assignment; temporary authenticated-tenant-wide account access. Production release blocked until RBAC enforcement. |

No runtime build or full source audit was performed for this planning document. Phase A starts by establishing an exact-commit baseline and dependency inventory.

## 4. Target repository layout

~~~text
grovehub-school-erp/
├── backend/                        # School Go module; existing
│   ├── launcher/
│   ├── libs/                       # Only school-owned technical contracts, if needed
│   └── services/
├── frontend/                       # School React application; existing
│   ├── app/
│   └── packages/{assets,business,domain,shell}/
├── internal/
│   ├── backend/                    # Independent Go module
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── launcher/
│   │   └── services/
│   └── frontend/                   # Independent React application
│       ├── package.json
│       ├── app/
│       └── packages/{assets,business,domain,shell}/
├── libs/
│   ├── backend/
│   │   └── core/                   # Standalone Go module, conditional on audit
│   └── frontend/
│       ├── base/
│       ├── styles/
│       ├── components/
│       └── tooling/                # Shared check configs, only if truly reusable
├── scripts/                       # App-specific entrypoints, shared mechanics
├── docs/
├── go.work                        # Local workspace: backend and Core
├── pnpm-workspace.yaml
├── pnpm-lock.yaml
└── Makefile
~~~

Current shared libraries: @groovehub/base, @groovehub/styles, @groovehub/components; Core module `grovehub/libs/core` (private local-source stage; versioned publishing pending). Existing school module name grovehub stays until a separately planned migration. Internal module grovehub/internal-platform, distinct from school and Core. Go module names must not collide; verify Go internal visibility and standalone builds. No empty placeholder directories. Generic workspace primitives may later move from shell to a separate library after proving no domain dependencies; do not move full shell now.

## 5. Phase A — extract shared frontend libraries first

### A1. Inventory and classify
- Enumerate every source, test, barrel export, CSS class, asset, translation, route, storage key and import under base/components/styles. Build import graph with TS compiler or existing boundary checker; trace indirect imports and side effects.
- Base owns pure formatting/i18n, generic utilities, readers, events, types; no app/API/auth/branding or DOM side effects unless explicitly separate browser capability.
- Components owns Button/Input/Select/AsyncSelect, DatePicker/DateRangePicker, Sheet/Dialog/Popover, Toast, Avatar/StatusPill, Universal Table, mobile list and toolbar, chart/editor primitives as needed, router adapter only if app-neutral.
- Styles owns neutral semantic tokens, base defaults, overlay scrollbars, primitives and configurable theme presets. School RDR colors/branding and internal orange accent selected by each app. Extract existing CSS consumers before removing aliases.
- Move SchoolMonogram and any other school assets/copy into school domain/assets. Inject route targets, translations, storage namespace, feature labels, app catalog and auth behavior from each application. Components never import @Business/@Domain/@Shell/@Assets of either product.
- Preserve school public landing pages, login/recovery and existing table behavior.

### A2. Workspace migration
- Create root pnpm-workspace.yaml with frontend, internal/frontend and libs/frontend/* (only existing packages during transition). Single lockfile, pinned package manager. Each package owns package.json, tsconfig, explicit exports, tests, README/RUNBOOK and build output.
- Establish @groovehub/base as dependency of Components; Styles exports CSS explicitly; React/ReactDOM peer dependencies, deduplicated runtime and compatible Radix/TanStack versions. No cycles. App-only query clients/providers remain app-owned.
- Move files with history; migrate imports to public package entrypoints or intentional subpath exports. During transition app aliases may resolve workspace source but published outputs MUST NOT contain unresolved @Base or app aliases. Update Vite, TS, Jest, CSS/Tailwind source scanning, prerender, lint, checks and CI path filters together.
- Preserve literal __() translation contract: each app initializes formatter before importing eagerly translated constants; isolated test bootstrap uses real formatter. Library must explicitly document boot requirement and types.
- Remove obsolete source copies and aliases only after consumer search. No compatibility forks or dual implementations. Baseline source architecture gates must remain at least as strict and scan library plus both apps.
- Release ESM, types, CSS exports and peer dependency ranges. Pack tarballs and smoke-test in isolated consumer without workspace aliases. Use workspace:* locally; CI can publish private GitHub Packages after verifying scope/registry ownership and auth. Version changesets or equivalent semver/changelog, tags and rollback. No public publication of proprietary code by default.

### A3. Design system acceptance
- Current school blue/zinc surfaces remain; internal orange/zinc theme with semantic statuses and adequate contrast. Theme keys and pinned app preferences namespaced; independent auth cookies.
- Desktop Universal Table: full-width fill layout, stable footer total records left and rows/page + page count + arrows right, filters/sort/search/settings/refresh, sticky and pinned columns, overlay scrollbars without layout shift.
- Mobile: compact title + icon-only create, hidden description, full search + settings/sliders opens single bottom sheet. Filter/sort/settings/batch size open dedicated sheets, refresh closes sheet and reloads. Custom mobile list rows and cursor infinite loading; no numbered footer, final count, stray zero, chevrons or row menus unless resource opts in. Account detail sheet holds actions.
- Preserve 16/14/12px username/email/status design tokens and 32px avatar; email single-line ellipsis only when needed. No badge border for plain status. Full-page Menu, bottom Home/Apps/Menu navigation, pin/unpin catalog, safe-area insets, no white strip. Desktop sidebar remains collapsible.
- Reusable auth layouts, form alignment, 1px focus, no desktop panel on phones, mobile keyboard/short viewport handling. Test light/dark/system, 320/375/393/720 and desktop, real Safari/device where available.
- Future component catalog/Storybook optional; do not block MVP on new tooling. Verify existing audit, architecture, typecheck, Jest, lint, build, prerender and browser suite before/after extraction.

## 6. Phase B — Core extraction (conditional, not blind move)

### B1. Audit each package and classify
Candidates: annotation mechanics, database/postgres/transaction, diff, errors, event envelope (not business events), headercontext, HTTP request/response/logging/tracing, IDs, idempotency mechanics, module lifecycle, normalization, pipeline/rest/service, request, search, telemetry and generic test infrastructure.
Inspect source and tests for school/service/vendor imports, permission or persona enums, hardcoded schema/route/cookie/header values, identity business claims, provider-specific behavior, environment defaults and implicit filesystem assumptions. Generate a table of every package, import graph, classification and disposition. Acceptance: no imports from school backend/services/launcher or internal app into extracted Core; no tenant-specific authorization or broker policy. Do not move business code just because it is in Core.

Known remediation:
- Move cookie names into each Identity/gateway application's configuration or neutral injectable cookie contract. Preserve existing school cookie names/wire behavior. Internal uses separate names and host/audience.
- Review Principal fields: generic verified actor/subject, optional tenant scope and session metadata, typed workload context; school legacy Admin and PersonGID semantics stay adapter/owner-owned. Introduce backward-compatible adapter during migration, remove Admin privilege after RBAC. Do not create a fake account to satisfy PrincipalFrom.
- Generators must discover module import path from target go.mod and accept explicit root/owner directories. Audit hardcoded grovehub, libs/core and service lists in every generator; parameterize output owner and import prefixes. Keep generated overlays owner-local, deterministic, never copy into maintained sources. Separate school/internal/core generation and test concurrency; do not write into downloaded read-only module cache.
- Domain-specific event definitions, outbox delivery, providers, configuration secrets and service authorization remain in respective modules. Keep current transaction context compatibility and error codes.

### B2. Move and validate
Once neutral, move Core to libs/core with go.mod/go.sum, documentation and generic tests. Use root go.work and explicit local workspace build commands; school/internal require the same tagged module version. For deployment without workspace use canonical private module path or reproducible vendored/relative build context as selected in ADR. Verify GOWORK=off release path; do not rely on local replace paths in published dependencies. Select canonical module identity before migration, replace all imports atomically, ensure context keys/types aren't duplicated by importing both old/new packages. Move generator commands/config and scripts in same change, validate generated routes, mappers, entities, search, clients and tests for both target roots. Retain existing school public API, migration and event versions. Update Docker/CI/build/Air/Make paths and owner docs. No old core facade containing duplicate types. If audit finds unresolved product coupling, keep affected package in school ownership until decoupled; block extraction, not school RBAC security work.

## 7. Phase C — independent internal application foundation

Suggested local defaults:

| Setting | School | Internal |
|---|---|---|
| Backend | 8080 | 8081 |
| Frontend | 4173 | 4174 |
| PostgreSQL | existing config | INTERNAL_DATABASE_URL (required, no school fallback) |
| Cookie / audience | existing school identity | distinct internal cookie namespace and session realm |
| Runtime | existing .runtime paths | .runtime/internal/{logs,pids} |
| NATS | existing | separate internal subjects, stream, durable and credentials |
| Host | existing | separate host; protected network/Zero Trust in real environments |

Verify actual existing port allocations first (Swagger etc). Separate DB server is supported but not required; different database/credentials mandatory. Reject accidental identical database target. Independent DB migrations/history, backups/restore, health, telemetry, secret configuration, build images and deployment. Use managed secrets; no URLs containing passwords in logs.
Proposed Make targets: internal-backend-start/stop/status/build/test/migrate, internal-frontend-start/stop/status/build/test, internal-dev-start/stop/status/start-with-logs, internal-identity-bootstrap-admin. Implement with existing scripts conventions, independent PID ownership and shutdown. School commands retain behavior. Do not expose internal tool via unauthenticated Cloudflare Quick Tunnel by default; use local/private access or authenticated named tunnel. No cross-app proxy or cookie forwarding.

### C1. Backend services and ownership
| Service | First release |
|---|---|
| identity | Internal people/accounts, credentials, invitations, sessions, MFA, internal roles/permissions/personas/features |
| registry | Modules, template definitions/versions, variable declarations, registration lifecycle |
| environment | L1/L2 configured target environments, credential references, environment entitlements |
| organisation | Authoritative school discovery from target API and optional cached references, school entitlements |
| approval | Generic policies/steps/requests/decisions, inbox and requester views |
| execution | Saved requests, typed resolution, dispatch, results, replay control |
| audit | User-facing immutable operational records, access logs |
| infrastructure | Internal outbox, messaging, worker adapters and external integration when needed |

Independent internal identity must not require a school Foundation person record. Model an internal employee/person under internal ownership (decide whether a distinct person module is justified). Reuse mechanical libraries, not school tables, roles, auth state or repository implementations. Use same annotated service layout, REST pipeline and established invitation/setup/reset flows; preserve non-enumeration, token expiry, secure HttpOnly opaque sessions, Argon2id, refresh/revocation, CSRF and Origin. Mandatory MFA before privileged access, recent reauthentication, recovery codes hashed, explicit initial admin bootstrap. Separate internal personas (DEVELOPER, TECH_LEAD, PLATFORM_ADMIN, SUPPORT, OPERATIONS, QA, DEVOPS, SECURITY, AUDITOR, MANAGEMENT) proposed; role/persona separation same as school plan. No employee email collision grants cross-app access.

### C2. Resource-level backend structure
For every actual resource use existing backend.md blueprint: entity/<resource>.go, dto/<resource>/{request,response,params,operations,search_fields}.go as required, repository/<resource>/{interface,impl,helpers/search}, validators, mappers, services/<resource>/operations/<verb>/{flows,helpers,strategies} only when appropriate, controller/<resource>/<tier>/{interface,impl}, routes/http.go, client/<resource> facade, bootstrap/impl typed wiring, migrations/OpenAPI, mirrored internal/tests/unit-test. Use @Entity/@Repository/@Searchable/@Routes/@RequestLine etc only where implemented by actual generators. One-line signatures, GORM and UnitOfWork, service authorization, no DB in validators, no cross-service private imports. Shared response envelope and error mapping. No generic repository or god service.
All search endpoints use exact Phase 2.1 hybrid contract: searchText, allowed filters/sort, pageInfo.pageNumber/rowSize/cursor, resource SQL allowlist, stable unique tie-breaker, authorization/tenant predicates before count and rows. Desktop numbered; mobile cursor. Use the existing operator types, not invented variants.

## 8. Phase D — Universal approval/request engine

One approval service per runtime/data boundary. School RBAC has its own instance/implementation and school-only requests; internal approvals never imply school employee permissions. Share only neutral mechanics after demonstrated reuse; do not create a cross-database approval engine.

Entities (internal approval schema):
- policies: id, name, requestType, selector conditions (module/template/environment/school/scope/risk as applicable), priority, mode, status, version, timestamps and actor.
- policy_steps: policy/version, sequence, eligible role/permission/user rules, required distinct count, allowSelfApproval, expiry.
- requests: requestId, type, owning service/resourceId/version, requester, reason, normalized immutable payload hash/snapshot reference, resolved policy/version, optional execution scope, status, expiry, version, timestamps.
- request_steps: resolved stage policy snapshot and state, required count, sequence.
- decisions: immutable decision id, request/step, actor, MANUAL/AUTO_ELIGIBLE/AUTO_POLICY, decision, comment, timestamp; unique eligible actor per step. System policy auto decision is not attributed as human review.
- Optional effect status on owning resource: approval approved is not same as action successfully executed. Durable outbox links request transitions to owner effects.

Modes: MANUAL, AUTO_IF_ELIGIBLE, AUTOMATIC (explicit policy). Missing/ambiguous policy fails closed. Auto eligibility checks active requester, exact step scope, permissions, policy and risk; only satisfies permitted stages. No auto-approval of own privileges, policy changes or critical operations. Reject/cancel/expire closes pending stages; revocation, stale resource or policy change triggers revalidation/reapproval before effect. Optimistic locking + unique decisions + idempotency; bulk approval per-item independently authorized and transactional, returns individual result rather than silently skipping or reporting global success. Require rejection reason and sensitive operation approval evidence.

Public UI **Create request** selects request type, then mounts the owner-defined form. Server dispatches to typed owning service (e.g. registration or execution); frontend never sends arbitrary resource type + payload to privileged restricted approval endpoint. Owner creates intent and approval atomically in its transaction where same DB, otherwise durable outbox/reconciliation; only owner applies approved effects. Approval handler registry is typed and allowlisted, no arbitrary callback URL or function name. Generic approval API:
- POST /approval/restricted/requests (trusted owner only if needed; prefer typed in-process boundary)
- POST /approval/protected/requests/search (own requests; explicit audit permission for wider scope), GET /requests/{id}, PATCH /requests/{id}/cancel
- POST /approval/protected/approvals/search (eligible inbox/history), GET /approvals/{id}, POST /approvals/bulk, GET /approvals/summary
- POST /approval/protected/policies, POST /policies/search, GET/PATCH /policies/{id}, PATCH /policies/{id}/activate or /deactivate; policy changes themselves separately approved.
All shortened paths above inherit /approval/protected. Request types MVP: API_TEMPLATE_REGISTRATION, API_TEMPLATE_UPDATE, API_TEMPLATE_DEACTIVATION, API_EXECUTION, ROLE_ACCESS, PERMISSION_ACCESS, ENVIRONMENT_ACCESS, SCHOOL_ACCESS, FEATURE_FLAG_CHANGE, APPROVAL_POLICY_CHANGE. Future types in §13. Define enum in owner contracts, not frontend literals alone.

## 9. Phase E — Registry and template variables

Module catalog: moduleId, unique code, name, description, active status, version, timestamps. Register via protected authorized API; modules identify school service groups (IDENTITY, FOUNDATION, INTEGRATION, INFRASTRUCTURE; later others). No school or environment in module/template registration.
Templates split stable identity (templateId, moduleId, name/description, status, activeVersion, audit) and immutable version (version, method, relative endpoint, scope GLOBAL/SCHOOL, operation READ/WRITE, headers, query/path/body defaults, schema, variables, required execution permission, risk, timeout, redaction rules, idempotency policy, status, creator/time). Store variable declarations as owned structured JSON with schema/version or normalized rows if query needs prove it. DB uniqueness on template/version, active pointer references approved version. Register creates DRAFT or PENDING_APPROVAL; submit creates registration approval, active only after approval. Old active version remains usable until replacement activates atomically. No edits to approved versions; deactivation is separately authorized, invalidates queued executions. Pending executions on retired versions require new approval/version.

API prefix /registry/protected:
- POST /modules, GET /modules, POST /modules/search, GET/PATCH /modules/{id}, GET /modules/{id}/api-templates (active, authorized catalog; bounded).
- POST /api-templates, POST /api-templates/search, GET/PATCH /api-templates/{id}, GET /api-templates/{id}/versions, GET /api-templates/{id}/versions/{version}, POST /api-templates/{id}/submit, PATCH /api-templates/{id}/deactivate.
Use Idempotency-Key on commands and expected versions on mutation. Endpoint allowlist derived from target registered contracts, explicit method and path; relative only, no raw host/scheme, redirects or arbitrary path expansion.

### Template definition example
~~~json
{
  "moduleId": "11111111-1111-4111-8111-111111111111",
  "name": "Search school accounts",
  "method": "POST",
  "endpoint": "/identity/protected/accounts/search",
  "scope": "SCHOOL",
  "operationType": "READ",
  "riskLevel": "LOW",
  "headers": {"Accept-Language": "${locale}"},
  "queryParams": {},
  "requestBody": {
    "searchText": "${searchText}",
    "filters": [],
    "sort": [],
    "pageInfo": {"pageNumber": "${pageNumber}", "rowSize": "${rowSize}"}
  },
  "variables": [
    {"key":"locale","label":"Language","type":"STRING","required":true,"defaultValue":"en"},
    {"key":"searchText","label":"Search","type":"STRING","required":false,"defaultValue":""},
    {"key":"pageNumber","label":"Page","type":"INTEGER","required":true,"defaultValue":1,"validation":{"min":1}},
    {"key":"rowSize","label":"Rows","type":"INTEGER","required":true,"defaultValue":25,"validation":{"min":1,"max":100}}
  ]
}
~~~
Path example: /identity/protected/accounts/${accountId}; query value example: status=${status}. Only values/path segments, not method, host, header names or JSON keys. Define every referenced key once, reject unknown/unused/conflicting definitions. Types STRING, TEXT, INTEGER, DECIMAL, BOOLEAN, ENUM, MULTI_ENUM, UUID, DATE, DATE_RANGE, DATETIME, ARRAY, OBJECT, SECRET_REFERENCE, LOOKUP (later types only when supported). Required/default/enum/range/length/pattern, input validation, sensitive flag, UI label and help text. Server enforces schema; client preview is not trusted.
Structured JSON traversal: exact placeholder substitutes actual typed JSON value, mixed literal string interpolation only scalar encoded string; reject incompatible types/missing required values. Path segments escape independently and reject traversal/encoded separators; query via URL encoder; header reject CRLF and reserved auth/tenant/system/forwarding names. No eval, script engine, shell interpolation, recursive expansion or arbitrary template expressions. Limits for nesting, lengths, regex cost and variable count. Secrets use references only and are resolved server-side at dispatch; bind secret version/allowed scope to approval and never log actual values. LOOKUP providers registered, read-only, scoped and independently authorized. Reuse neutral substitution mechanics in future typed deployment/job templates after tests.

## 10. Phase F — environment, schools, saved requests and execution

Environment: environmentId, code (L1/L2 etc), name, baseURL/approved target binding, credentialRef, status, version. Internal config owner manages trusted environment URLs; no user-editable execution URL. L1/L2 are deployment targets, not approval stages. Actual URLs and production mapping unspecified. School directory fetched from trusted target global list API; cache only necessary references keyed by environment and school ID/tenant ID. Never trust a user-supplied tenant header. Environment/school access explicit, optionally time-bound; verify membership/current status at submission and dispatch.
Endpoints: /environment/protected/environments (create/list/search/detail/update), /environment/protected/environments/{id}/schools (authorized discovery via organisation client), access grant/revoke routes owned by relevant resource, approved when required. Separate privileged configuration from ordinary catalog reads.

Saved requests: owner, visibility PRIVATE/SHARED, name, approved templateId/version, environmentId, nullable schoolId, variable values (no plaintext secrets), version/timestamps. Share with explicit audiences; shared visibility never grants execution. CRUD/search/duplicate under /execution/protected/saved-requests. Version mismatch must revalidate and show diff. Save is not execute.

Execution API POST /execution/protected/executions:
~~~json
{
  "templateId":"22222222-2222-4222-8222-222222222222",
  "templateVersion":1,
  "environmentId":"33333333-3333-4333-8333-333333333333",
  "schoolId":"44444444-4444-4444-8444-444444444444",
  "variables":{"locale":"en","searchText":"teacher","pageNumber":1,"rowSize":25},
  "reason":"Investigating support ticket TEST-001"
}
~~~
SchoolId required for SCHOOL, forbidden for GLOBAL. Frontend does not submit arbitrary method/URL/headers/body overrides; values are resolved using approved definition. Always create execution ID and independently resolve approval policy. Payload hash binds method/path/query/safe headers/body, template version, environment target/version, scope/school, actor, reason, credential reference/version and policy snapshot. After approval revalidate all entitlements and template state, atomically claim execution using durable worker, dispatch once with timeout and correlation. Record QUEUED/RUNNING/SUCCEEDED/FAILED/REJECTED/CANCELLED/EXPIRED/UNKNOWN_OUTCOME; approval pending/approved separate. Unknown outcome on lost response for non-idempotent operation; no blind retry. Idempotency key propagated only where target supports semantics. Rerun creates new execution/approval; cancellation after dispatch is best-effort, not undo. Bound concurrency, per-target rate limits, max response size; sensitive response fields redacted before persistence/display/export. Retention and deletion policy configured, not arbitrary forever; audit metadata retained per policy.
Execution entity: executionId, immutable template/version, environment/school, requester, approvalId, payload fingerprint, encrypted snapshot or safe reference, status/version, timestamps, correlation/trace and failure summary. execution_results: status code, duration, safe headers/body, truncated flag, result classification, timestamps. Access to result and export separately authorized.
Endpoints: POST /executions/search, GET /executions/{id}, GET /executions/{id}/request, /response, /approvals, /timeline, PATCH /executions/{id}/cancel, POST /executions/{id}/rerun (all /execution/protected).

### Trusted system invocation
Internal backend never connects to school DB. Use environment-specific short-lived workload credential, exact audience/operation/tenant/execution binding, TLS, replay controls. Existing school gateway protected routes currently require session cookies; implement explicit restricted workload dispatch/adapter mapping only approved target operations to typed service clients with resource authorization. Alternatively extend protected gateway for explicitly registered workload-capable operations after security review; never disable browser origin/CSRF or accept arbitrary X-User-ID. Reserved transport headers injected by verified adapter only; X-Tenant-ID for school scope, X-Actor-Type SYSTEM, X-User-ID -1 display convention, X-Initiated-By internal UUID, X-Execution-ID, X-Correlation-ID. School verifies workload scope and derives actor independently from signed claims; no fake UUID session. Strip spoofed headers. No browser cookies forwarded. Block unapproved hosts, IP/DNS rebinding, unsafe redirects, credentials in URL; validate network allowlist at connect time. Log both SYSTEM and real initiator/approvers. Fail closed if target route contract, permission or identity is unknown.

## 11. Phase G — Internal frontend: full page inventory

Structure mirrors frontend/app and packages/{business,domain,shell,assets}. Shared imports only from extracted libraries; keep auth/API readers, app route registry, permissions, feature navigation and Orange theme app-owned. Start with school templates/containers and reuse actual components, not a new parallel design system. Pages use dedicated create/edit/detail routes, sheet for small contextual actions, table list for catalog/execution history. Universal request form registry maps request type to domain-owned renderer, validator, submit service; no untyped giant switch in component. Use typed renderer/helper with explicit branch returns and tested pure logic.

| Area | Pages and smaller views (MVP unless marked future) |
|---|---|
| Auth | /login, /forgot-password, /set-password, /reset-password, /auth/mfa, /auth/mfa/setup, recovery codes, expired proof, session expired, 403, 404, 503 |
| Home | Overview, pending approvals, recent executions, health summary, notifications (only real data), app catalog/search/pin/unpin |
| Identity | Users list/create/invite/detail/edit/status actions/resend, roles list/create/copy/edit/detail, permissions catalog/details, role-permission matrix, personas, feature flags, access requests, profile/security/session revoke/MFA |
| Registry | Modules list/create/detail/edit; module API list; templates list/create/edit/duplicate/detail; versions list/detail/diff; registration request history |
| API Explorer | Full-width table of active templates filtered by module/scope/method, Create opens dedicated execution form, details and saved requests; no persistent split explorer |
| Execution form | Module -> API -> environment -> conditional school -> dynamic variables -> resolved request preview -> reason -> submit; only registered input controls, no reserved headers editing |
| Saved requests | List, private/shared filter, create from template, detail/edit/duplicate/delete/share; execute creates new intent |
| Execution history | List, detail tabs Overview/Request/Response/Headers/Approval/Timeline, rerun, response export when allowed; failure/redaction states |
| Approval dashboard | Page title + Create request above Approvals and Requests tabs; pending/approved/rejected/auto/all views, bulk action, review detail with exact payload/diff/steps, requester cancellation |
| Approval policies | List/create/edit/duplicate/detail, stage config, simulation, version history; delegation future |
| Environments | List/create/detail/edit, target health, scoped access; secrets metadata only |
| Schools | Environment-scoped list/detail, tenant references/access and discovery failure state |
| Audit | List, detail, actor/target/correlation filters, export only with permission |
| Settings | Profile, appearance, sessions, security, app pin preferences |

Every list has loading, error/retry, empty/no-match, permission denied, stale/disabled, selected/filter state, desktop/mobile variants. Use current Universal Table search/filter/sort/date range/selection and mobile behavior. Header Create compact icon on mobile, one-line row summary, no excessive cards. On mobile form sheets full width, menu full page; bottom nav safe-area-aware. Orange only via semantic theme tokens; no raw hardcoded color in domain. Avoid arbitrary HTML, owner CSS Modules, no mutable let, no ternary rendering, literal __() with interpolation, typed hooks/renderers, pure functions unmocked, tests describe named function. Follow exact current guidelines and remaining debt obligations. Render only authorized and enabled modules; frontend guard is UX, backend authorization mandatory.

## 12. Phase H — auditing, reliability and operations

Audit append-only domain records with actor kind, initiator/approver, resource, operation, environment/school scope, correlation, safe before/after references, time, outcome. Separate audit from debug logs and outbox retention. Outbox event envelope existing generic version fields; internal namespace e.g. internal.registry.template.submitted.v1, internal.approval.request.approved.v1, internal.execution.completed.v1 (final subjects match events.md). Events contain minimum data, no secrets/raw sensitive payload; consumers deduplicate, handle out-of-order, reconcile missing outcomes. Internal messaging failure must not turn unapproved intent into executable command. Trace across services; redaction before logs, metrics avoid high-cardinality PII. Backups and restore drills, operational alerts, worker recovery, explicit migration rollout.
Test concurrency: double submit, duplicate approval, concurrent decisions, stale payload/role/policy, replay token, cross-school/environment, disabled feature, worker crash before/after dispatch, unknown outcome, refresh/revocation, redaction, SSRF, malformed placeholders, pagination duplicate/terminal cursor, mobile sheet/focus/no overflow. Use real PostgreSQL and NATS integration, not mock claims of ACID. Contract tests with school target sandbox; no production data. Preserve release gates.

## 13. Future control plane (do not implement prematurely)

| Module | Future pages/typed operations |
|---|---|
| deployment/release | services/release catalog, deploy/promote/rollback requests, pipeline details/logs/history, environment configuration |
| monitoring | environment/service health, API latency/errors/throughput, DB/NATS metrics, logs/traces, alert rules and active alerts |
| orchestration | cluster/namespace/node/workload/pod/service/statefulset/volume/config metadata, scaling/restart/rollout approval, pod logs/events |
| incident | incidents list/detail/timeline/assignment/resolution |
| scheduler | job templates, schedules, run/trigger/retry history, logs |
| database | instance health, migration registry, approved migrate/restore, backup/restore history |

Reuse identity, entitlements, environment registry, typed variable mechanics, approvals and audit; separate resource-specific owners/workers/providers. GitHub Actions, Grafana, Kubernetes may be providers through existing Strategy/Provider conventions (not required MVP integrations). Future request types DEPLOYMENT, DEPLOYMENT_ROLLBACK, SERVICE_RESTART, WORKLOAD_SCALING, DATABASE_MIGRATION, JOB_EXECUTION, SECRET_ROTATION. No arbitrary shell/kubectl/SQL console or raw secret display. Future navigation hidden until implemented and feature-enabled. Avoid duplicating approval UI per module.

## 14. Delivery plan and acceptance gates

| Work package | Output | Gate / dependency |
|---|---|---|
| A0 baseline | Exact source dependency inventory, source/test/CI baseline, ADR paths/package names | Do not treat prior green CI as current |
| A frontend libs | Shared Base/Styles/Components, school uses packages, private release readiness | School build, all tests, visual regressions, no domain imports |
| B Core | Audit disposition, neutral Core module and generator migration | School generation/build/test/migration compatibility; standalone module tests |
| R school RBAC | Companion plan, migrate legacy roles and enforce permissions | Can run alongside A/B once coordinated; mandatory before production |
| C internal foundation | Independent DB/runtime/identity/MFA, frontend shell | Separate keys/DB/ports; internal login and cross-app denial |
| D approval | Typed engine + dashboard, policy/step tests | Manual/auto eligibility and bulk concurrency, no execution bypass |
| E registry | Modules/templates/versions/variables | Register, approve, activate, reject, version diff; no school in template |
| F execution | Environments/school lookup, saved requests, workload bridge/results | School RBAC + gateway adapter; end-to-end sandbox, idempotency |
| G frontend | All MVP pages and responsive states | Shared components only; full light/dark/mobile test |
| H hardening/release | Audit, observability, security review, recovery, CI/release | No production launch with temporary permissive policy |
| Future | Operations modules | Separate approved plans |

Suggested phase labels 2.2 shared libs, 2.3 school RBAC, 2.4 internal foundation, 2.5 approvals/registry, 2.6 execution/UX, 2.7 release hardening. Labels are proposals, not branch names. PRs small, separate mechanical moves from behavior; coordinate shared dependency migration with RBAC branch to avoid conflicts. Expand/backfill/validate/switch/contract migrations; never edit applied history. Rollback via compatible app versions and backups, not destructive DB down. Document new service-owned architecture/OpenAPI/runbooks as each service is implemented. Update docs/AGENT.md and engineering index path references only when moves land.

## 15. Definition of done and open decisions

- [ ] Full Core and frontend dependency audit, boundary CI and standalone package release tests.
- [ ] Existing school functionality and contracts preserved except intentional RBAC access restriction.
- [ ] Separate internal database, credentials, identity, session realm, runtime, migrations, NATS and deployment.
- [ ] Dynamic variables server-validated; template registration and every execution separately approved; GLOBAL/SCHOOL isolation proven.
- [ ] Universal request UI and approvals/requests tabs; bulk, auto-approval, detail and audit flows tested.
- [ ] Browser light/dark and real mobile tests; existing architecture checks and backend generated overlay tests pass; coverage target 100% where project requires, report actual gaps.
- [ ] Exact commit build/test/security/CI evidence and rollback runbook; no public internal Quick Tunnel.

Open before implementation: private publishing permissions and release tag/version policy, actual L1/L2 targets and credential mechanism, school discovery API owner, internal MFA enrollment/bootstrap and recovery procedure, approval policy owners/risk matrix, audit/response retention and RPO/RTO. Confirm with maintainers; proposed defaults are not deployed facts.

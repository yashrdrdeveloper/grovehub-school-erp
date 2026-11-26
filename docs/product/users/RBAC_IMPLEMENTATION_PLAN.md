# School Identity — RBAC Implementation Plan

Status: PROPOSED. Owner: School Identity engineering. Audience: backend/frontend/security engineers.
Baseline: phase-2.1-identity-account @ 9fa83d6ca850069d9243e67a1a4d4355adce483b (26 November 2026).
Related [internal platform plan](../internal/INTERNAL_PLATFORM_IMPLEMENTATION_PLAN.md), [current interim authorization](ACCOUNT_AUTHORIZATION.md), [Identity architecture](../../../backend/services/identity/docs/engineering/architecture.md), [backend standard](../../engineering/backend/backend.md), [frontend standard](../../engineering/frontend/FRONTEND_CODE_GUIDELINES.md), [gateway](../../engineering/backend/gateway.md), [hybrid pagination](../../engineering/backend/HYBRID_SEARCH_PAGINATION.md), [persistence](../../engineering/backend/persistence.md), [security](../../engineering/security/SECURITY_GUIDELINES.md).

## 1. Existing contract and release blocker

Verified current sources: backend/services/identity/internal/app/entity/account.go; constants/role/catalog.go; dto/account/{roles,request,response,search_fields}.go; validators/account/helpers/roles.go; repository/account/helpers/search/config.go; controller/account/protected/account.go; migration 202611250001_account_role_assignment.sql. Existing Account UUID, PersonGID, tenant, nullable primary_role_id TEXT and secondary_role_ids JSONB array, status/version/search_details/timestamps. Role catalog GET /identity/protected/accounts/roles returns six ID/name pairs. Current assignment accepts only hardcoded codes, validates primary and distinct secondary role set, orders secondary roles deterministically. Existing account create/update accept roleAssignment; response primaryRoleId/secondaryRoleIds. Search has primaryRoleId filter and uses AccountConfig with stable account ID cursor.
The current temporary policy allows any authenticated tenant member to manage accounts in the tenant. Do not release to actual school users until all protected operations enforce permission and tenant ownership. Older Identity architecture prose about admin allowlist may not reflect the newer interim policy. Inspect implementation and tests rather than assume legacy Admin flag enforces access.

## 2. Model and vocabulary

Persona is fixed school job category, not a permission or assignable role. Role is a school/tenant-owned dynamic named bundle of permissions, linked to one or more personas (proposed multiple allowed). Permission is code-defined enum representing one backend capability, registered globally via restricted API; registration alone grants nobody access. Feature flag definition is global, tenant enablement is separate. Account has one primary and zero/more secondary roles; primary is for UX/default context, effective grants union across active assignments. No implicit permission from persona, role name or primary priority. Roles have explicit permission mappings, no runtime inheritance from similar role. Optional similarRoleId is copy-on-create within tenant, not a live parent. User-selected explicit permission list wins; if omitted on copy, snapshot source grants and authorize every copied grant. Reject ambiguity/duplicates; preview source version.

### Initial hardcoded personas
Preserve existing six codes as first catalog: SCHOOL_ADMINISTRATION, TEACHER, COORDINATOR, ACCOUNTANT, LIBRARIAN, RECEPTIONIST. Proposed future additions PRINCIPAL, VICE_PRINCIPAL, ADMISSIONS, HR, TRANSPORT, IT_SUPPORT, COUNSELLOR, SUPPORT_STAFF, STUDENT, PARENT_GUARDIAN only when actual modules/actors are introduced. Do not seed future personas as active privileges. Expose id/name/description from constants/persona, frontend uses server catalog, no persona DB CRUD. Existing role IDs are legacy role codes; migration creates real tenant-specific role IDs with same display names.

## 3. Persistence design (proposed)

All tables in school identity schema except school approval engine tables in approval schema or Identity-owned approval schema per ownership ADR. Use UUIDs for new aggregate IDs; tenant foreign-key boundaries via composite (tenant_id, role_id/account_id) constraints. Existing tenant data/type must be inspected before DDL; no assumptions that tenant_id is UUID.

| Table | Fields / constraints |
|---|---|
| permissions | permission_id UUID, code enum unique, module, name, description, optional feature_id, status, risk/grant policy metadata, version, created/updated; code immutable |
| features | feature_id UUID, code enum unique, module, name/description, defaultEnabled (safe false for new protected capability), status/version/timestamps |
| tenant_features | tenant_id + feature_id unique, enabled, version, actor/timestamps; explicit enablement, disabled/inactive denies |
| roles | role_id UUID, tenant_id non-null, name, normalized_name unique per tenant (choose inactive reuse policy), description, status, version, search_details, created_by/updated_by/timestamps |
| role_personas | role_id + tenant_id + persona_code unique; validate enum, at least one on role creation |
| role_permissions | tenant_id + role_id + permission_id unique, granted_by/time; permission global, role tenant-scoped |
| account_roles | tenant_id + account_id + role_id unique, assignment_type PRIMARY/SECONDARY, grant metadata, optional expiry only if requirement; partial unique primary per account, no secondary without primary enforced in operation and DB where possible |
| authorization_revisions | tenant/account revision or equivalent transactional version source for cache/session invalidation; design after measuring |
| approval policies/steps/requests/request_steps/decisions | school-local typed requests and tenant-scoped approvers, immutable decisions; follow shared workflow semantics in companion plan |

Normalize account_roles to enforce tenant-safe FKs and uniqueness. Existing text/jsonb columns remain during migration compatibility; don't maintain two independent sources permanently. No direct account_permissions grants in MVP: all permissions come through roles, including access requests, unless explicitly approved design later. Permission/feature code definitions owned by Identity or domain-owned public contract contribution via checked registration manifests; use typed Permission enum and module/feature mapping, never accept arbitrary strings as new runtime capability. For featureless permission, check active registration and role; for associated feature, require enabled tenant flag. Register only known code and correct metadata, immutable binding changes separately reviewed. Feature toggle is not a grant. Show disabled permissions in administrative matrix as disabled with reason (subject to management read permission); hide unusable actions in ordinary UI.

## 4. API contracts (proposed paths)

Existing prefix /identity/{public,protected,restricted}. Preserve existing account endpoints and envelope; update OpenAPI, DTOs, clients, tests and frontend together. All reads authenticated and scoped except public auth unchanged. Every command has explicit capability, expected version when modifying, idempotency where retriable, structured errors. Restricted means verified workload, not browser admin.

| Method | Route | Required capability/behavior |
|---|---|---|
| GET | /identity/protected/personas | role catalog view; fixed enum |
| POST | /identity/restricted/permissions | register known enum definitions, batch supported, idempotent |
| POST | /identity/restricted/features | register known feature enum definitions, batch supported |
| PATCH | /identity/restricted/permissions/{permissionId} | controlled metadata/lifecycle with version, no arbitrary enum |
| PATCH | /identity/restricted/features/{featureId} | controlled definition lifecycle |
| POST | /identity/protected/permissions/search | permission catalog view; feature/module/status filters |
| GET | /identity/protected/permissions/{permissionId} | catalog detail |
| POST | /identity/protected/features/search | feature list with effective tenant state |
| GET | /identity/protected/features/{featureId} | detail |
| PATCH | /identity/protected/features/{featureId}/state | tenant enable/disable with approval when policy requires; school or internal workload policy explicit |
| POST | /identity/protected/roles | create tenant role, personas, permissions, optional similarRoleId |
| POST | /identity/protected/roles/search | role view |
| GET | /identity/protected/roles/{roleId} | role detail |
| PATCH | /identity/protected/roles/{roleId} | name/persona/description changes, expected version |
| PATCH | /identity/protected/roles/{roleId}/deactivate | reject active assignments or approved reassignment plan |
| GET | /identity/protected/roles/{roleId}/permissions | whole role permission set, no per-permission API calls |
| PUT | /identity/protected/roles/{roleId}/permissions | atomic replacement of complete set, expected version |
| GET | /identity/protected/role-permissions | whole authorized tenant role-to-permission matrix in one response |
| GET | /identity/protected/accounts/roles | existing compatibility assignable role catalog now dynamic |
| GET | /identity/protected/accounts/{accountId}/roles | account assignment detail |
| PUT | /identity/protected/accounts/{accountId}/roles | replace primary/secondary assignment with approval policy |
| GET | /identity/protected/accounts/{accountId}/permissions | effective permissions for authorized management |
| GET | /identity/protected/auth/permissions | current user's effective capabilities/features/revision |

PUT new method is proposed; generator supports it, follow current API standard or use PATCH if replacement contract mandates it. For tenant feature toggles distinguish internal restricted entitlement administration from school-managed features; school cannot enable purchased/platform-restricted entitlement simply by possessing local ROLE_ADMIN. Default denies when entitlement off. Batch registry input bounded, deterministic, reject unknown enums, same key/different payload conflict; DB uniqueness and one transaction. No production registration via arbitrary console unless workload has exact registry capability and approval.

Role create request example (proposed):
~~~json
{
  "name": "Senior Teacher",
  "description": "Teaching and attendance responsibilities",
  "personas": ["TEACHER"],
  "permissions": ["ACCOUNT_SELF_VIEW", "ATTENDANCE_VIEW"],
  "similarRoleId": null
}
~~~
ATTENDANCE_VIEW is future illustrative permission, not registered until its actual backend operation exists. Optional similarRoleId uses same-tenant approved source and no hidden inheritance; source changes after creation do not affect copy. For role mapping replacement use permissionCodes + version, no incremental partial update ambiguity. Return whole role detail including assigned permission codes, persona codes and version. Matrix example payload:
~~~json
{
  "roles": [{"roleId":"11111111-1111-4111-8111-111111111111","name":"Teacher","version":2}],
  "permissions": [{"code":"ACCOUNT_SELF_VIEW","module":"IDENTITY","enabled":true}],
  "mappings": [{"roleId":"11111111-1111-4111-8111-111111111111","permissionCodes":["ACCOUNT_SELF_VIEW"]}]
}
~~~
Use batch reads and compact mapping arrays, no N+1; bound tenant role/permission catalog sizes or reject oversized request with explicit alternative. Do not silently truncate matrix. Exclude other tenants even for shared global permission list. Ordinary users cannot enumerate other roles unless authorized.

## 5. Permission inventory and enforcement

Start from actual existing account controller methods: Create, Search, Roles, Get, Update, Suspend, Activate, Deactivate, Resend. Define service-owned enum names (proposed) IDENTITY_ACCOUNT_CREATE, IDENTITY_ACCOUNT_SEARCH, IDENTITY_ACCOUNT_VIEW, IDENTITY_ACCOUNT_UPDATE, IDENTITY_ACCOUNT_SUSPEND, IDENTITY_ACCOUNT_ACTIVATE, IDENTITY_ACCOUNT_DEACTIVATE, IDENTITY_ACCOUNT_RESEND_INVITATION, IDENTITY_ACCOUNT_SELF_VIEW, IDENTITY_ROLE_VIEW/CREATE/UPDATE/DEACTIVATE, IDENTITY_ROLE_PERMISSION_VIEW/MANAGE, IDENTITY_ACCOUNT_ROLE_VIEW/ASSIGN, IDENTITY_PERMISSION_VIEW/REGISTER, IDENTITY_FEATURE_VIEW/MANAGE/REGISTER, APPROVAL_REQUEST_CREATE/VIEW, APPROVAL_VIEW/APPROVE/REJECT/POLICY_MANAGE. Example payload abbreviations above must be reconciled to final enum; prefer explicit prefix (ACCOUNT_SELF_VIEW -> IDENTITY_ACCOUNT_SELF_VIEW). Do not register permissions for nonimplemented routes except explicit disabled declarations. Audit ALL protected and restricted routes across school modules, not just Identity; route-to-capability map uses registered method+pattern, service-level invariant and tenant predicate. Public authentication/refresh/setup/reset semantics remain unchanged, limited grant context cannot access normal protected operations. Self/profile operations explicitly own account only; never use generic management permission for own profile.
Decision order: verified active principal/session -> authorized tenant scope -> registered operation permission -> active role assignments and active permission -> effective feature/entitlement -> resource ownership and business policy -> approval when required. Unknown missing context/permission/policy/feature fail closed. No trusted principal from request body or header. Workload identity separate, explicit allowed operation/tenant and no universal -1 bypass. Internal in-process typed callers also enforce equivalent capability. Deny by default for unmapped protected operations (explicitly enumerate safe self-service routes).
Privilege escalation: role creators/managers must possess delegable grants or a distinct constrained delegation capability, cannot assign privileges they cannot grant; sensitive grants independently approved. Enforce checks on copy, bulk, account create/update role assignment and matrix replace. Last active authorized administrator protection; lock tenant authority changes so concurrent requests cannot remove both last admins. Revalidate on approval apply. Role assignment, permission or feature changes immediately affect existing sessions (authoritative evaluation per request for first release; optional revision-keyed cache later, no stale grant window). Token/session refresh must not restore revoked access. Return 401 unauthenticated, 403 forbidden, 404 for hidden cross-tenant resource as established contract, 409 stale/duplicate conflicting changes, validation errors in existing envelope. Never leak other tenants via counts/search.

## 6. Universal school requests and approval management

School approval service is separate from internal platform (own DB/tenant and actors). Generic typed workflow, resource owners apply effects; no arbitrary approval payload from frontend. Request types proposed ROLE_ACCESS, PERMISSION_ACCESS, ROLE_ASSIGNMENT_CHANGE, ROLE_PERMISSION_CHANGE, FEATURE_FLAG_CHANGE and APPROVAL_POLICY_CHANGE. Permission access request resolves to a role-based assignment or approved role change; do not silently add direct grants. Policy selector can depend on tenant, operation, role, permission risk; manual/auto-if-eligible/automatic explicit. If requester is authorized approver for exact policy stage and self-approval allowed, auto decision recorded; own privileged access or policy changes require independent approver. Reject unknown policy; expiry, version checks, eligibility at decision and effect, immutable decisions and actor/time. Approved does not mean effect applied; track applied/failed/stale separately. No cross-tenant approval. Request creation by typed owner or universal protected request dispatcher with registered typed validators.

| Method | Route (school /approval/protected) |
|---|---|
| POST | /requests (typed request type and payload, reason, idempotency) |
| POST | /requests/search (own requests, scoped audit access separately) |
| GET | /requests/{requestId} |
| PATCH | /requests/{requestId}/cancel |
| POST | /approvals/search (eligible inbox and authorized history) |
| GET | /approvals/{requestId} |
| POST | /approvals/bulk |
| GET | /approvals/summary |
| POST | /policies and /policies/search |
| GET/PATCH | /policies/{policyId} |

Bulk example:
~~~json
{
  "decisions":[
    {"requestId":"22222222-2222-4222-8222-222222222222","version":1,"decision":"APPROVE","comment":"Reviewed"}
  ]
}
~~~
Return per-item result and updated status, not blanket success. Bulk max size configurable; verify each item tenant, policy stage, active reviewer and version. Unique decision constraints and transaction, no duplicate side effect. Rejection reason mandatory. Owner records intent and outbox atomically; approval outcome idempotently applied only to exact payload/version. Initial policies bootstrapped through controlled deployment; no empty permissive default. School-local approval implementation can share neutral mechanics after duplication proves useful, not a dependency on internal service database.

## 7. Search, persistence and backend implementation details

Use current resource-first structure under backend/services/identity/internal/app:
- constants/{persona,permission,feature,error,events}; existing role catalog becomes legacy migration adapter then remove when no consumers.
- entity/{role,role_persona,permission,feature,tenant_feature,role_permission,account_role,...}.go annotated @Entity.
- dto/{persona,permission,feature,role,rolepermission,accountrole}/ relevant request,response,params,operations,search_fields. Validators syntax/type/enum; services perform active/ownership/delegation/approval policy; mappers expose safe response.
- repository/<resource> interface, impl GORM, resource search config, tenant predicate, joined or batch hydration. Add indexes normalized name, tenant/role, tenant/account, active assignment lookup, feature lookup and stable search/order. Do not put authorization in generic Core.
- services/<resource>/operations/{create,search,get,update,replace,deactivate,...} with explicit transaction and multistep flows only where needed. Typed client façade, generated annotated controllers/routes. Gateway injected authorizer and service enforcement, no duplicate permission lists in controller implementations.
- New migrations and OpenAPI. Event producers for role, permission mapping, account assignment and feature changes with versioned envelope, old/new safe data or references, outbox same transaction, audit record actor/approval; events idempotent. No passwords/secrets.
- Search exact Phase 2.1 contract: pageInfo.pageNumber,rowSize,cursor, rowSize default25 max100, legacy rows alias; searchText, existing filter operator enum, stable multi-sort + ID tie-breaker. Role name/persona/status/createdAt/updatedAt; permission code/module/feature/status; approval type/status/requester/date; account existing fields and primaryRoleId, update SQL mapping with migration. Tenant/eligibility applied before totalRecords. Validate search allowlist, no interpolated SQL. Mobile load more cursor dedup/reset. Use actual request/response shape from shared code, not invented copies.

## 8. Migration and rollout (critical)

1. Inventory live tenant values, legacy codes, null primary and secondary arrays; take backup and rehearse restore. Record existing schema/history; never alter 202611250001 or destructive initial migration.
2. Add new tables, indexes, constraints, disabled feature/permission registry definitions and revision mechanism. Create six seeded tenant roles per actual tenant with explicit legacy-code-to-new-UUID mapping; personas preserve old role names. Do not auto-grant admin permission based only on SCHOOL_ADMINISTRATION label or all authenticated accounts.
3. Backfill account_roles deterministically using mapping, quarantine unknown/invalid roles/null tenant and report; verify counts, duplicates and one-primary invariant. Approved explicit per-tenant initial admin assignment via controlled bootstrap and known operator identity, protect against removing last administrator. Other roles start least privilege, approve mapping with school owner.
4. Compatibility: existing Account DTO uses strings for role IDs; during transition translate allowed legacy IDs at boundary for deployed clients if necessary, return existing shape; migrate frontend role catalog to dynamic IDs and assignment endpoints together. Temporary dual read/write only one transaction/source-of-truth, parity checks and cutover deadline. No separate uncoordinated mutation paths.
5. Deploy authorization in test/staging, dry-run decision logging without PII, validate role access and deny paths. Switch all operations to fail-closed RBAC and remove temporary permissive policy/legacy Admin authority. Do not deploy dry-run to real school users as security enforcement. Revoke/invalidate stale sessions and verify audit.
6. After compatible release and parity, remove legacy hardcoded assignability and unused columns via separate forward migration only when no clients/readers need them. Rollback to secured compatible build; never restore tenant-wide permissive access to production. DB restore plan tested.

## 9. School frontend scope

Reuse extracted Base/Components/Styles; school blue theme. Existing User Setup desktop/mobile layouts, search, create and details retained; primary/secondary role select uses dynamic role catalog. Role assignment changes that need approval show request status, not optimistic privilege. Show effective permissions in authorized detail, disabled-feature explanation.
Pages: Identity > User Setup, Roles table/create/edit/detail/duplicate (similar role select + persona multiselect + permission selector grouped module/feature), Permissions read-only catalog/search/details, Role Permission Matrix (roles vs permissions, bulk draft changes and atomic save with version; mobile role selector + compact permission list rather than horizontal matrix), Personas catalog, Feature flags state list/detail/confirmation, Access Requests and Approval Dashboard. Approval dashboard title and Create request above Approvals / Requests tabs; typed request form changes per request type. Pending, approved, rejected, auto-approved history, bulk selection, detail exact request and effect status, cancellation. Only show actions allowed by fetched current permission snapshot and flags; API denies regardless of UI. Update auth/session cache on changes or forbidden responses; no token storage in localStorage. Navigation catalog filters permissions/features, pinned unauthorized app removed/hidden safely. Same universal table, bottom sheets and mobile infinite list behavior, full-screen menu. Follow existing frontend owner structure, translation, types, hooks, renderer, CSS Modules and tests. No raw role-name comparison as authorization.

## 10. Test and delivery checklist

- [ ] Inventory every gateway method pattern and protected service/in-process operation; tests prove default deny, missing mapping, tenant A/B search/count, self vs management, public and restricted separation.
- [ ] Enum registration idempotency, unknown code, immutable code, feature off/on/retired and tenant scope, no grant by registration.
- [ ] Dynamic role create/copy, duplicate names, multiple personas, invalid source tenant/version, copy permission escalation, primary/secondary duplicate rejection, union, disabled role/permission.
- [ ] Whole role-permission matrix and replace concurrency, no N+1, atomic rollback/outbox, stale version.
- [ ] Account create/update cannot bypass role assignment policy; last-admin concurrent protection, revoked access in live session, workload -1 cannot spoof.
- [ ] Approval manual/auto eligibility, independent high-risk review, self escalation denied, bulk partial results, duplicate decisions, expiry/cancel/stale payload, effect retry and audit.
- [ ] Clean and upgrade DB migration, backfill parity and restore; real PostgreSQL constraint/concurrency tests, real outbox/broker integration.
- [ ] Annotated generation/overlay, go build/lint, Ginkgo/Gomega tests mirrored, test target 100% coverage where required (report actual coverage and gaps). Use make unit-coverage package=identity if confirmed in Makefile.
- [ ] Contract/OpenAPI, existing account/activation/forgot/reset/login/session regression, frontend architecture/audit/typecheck/Jest/lint/build, desktop and real mobile light/dark/browser.
- [ ] Remove temporary authorization warning only when actual tested enforcement replaces it; update owning Identity docs, API workflows and release notes.

Order: permission/persona/feature registry -> role and mappings -> account migration and assignment -> gateway/service enforcement and safe bootstrap -> school approval workflows (required for sensitive changes before release) -> frontend -> end-to-end hardening. Coordinate with shared library work, but do not delay security fix behind internal UI. Internal platform executes school restricted registrations only through the scoped workload bridge after this plan's enforcement exists.

Open decisions: actual tenant ID type/directory, final permission naming and grant policy, approved initial administrator mapping, feature entitlement owner and defaults, school approval service location, exact API compatibility window, cache strategy and approval risk policy. No implementation or test execution is claimed by this plan.

## Implementation checkpoint (26 November 2026)

The additive RBAC schema, code-owned vocabularies, pure evaluator and tenant-scoped snapshot repository are implemented. Snapshot loading uses one SQL statement and rejects missing or inconsistent scope; seven SQL-mock regression scenarios pass. The school-specific principal and cookies are owned by `backend/libs/schoolgateway`. The evaluator ignores the legacy Admin claim. No live migration, registry seed, role backfill, protected route switch or approval workflow has been applied. Next implement controlled catalog registration/bootstrap and an authorization service that loads the verified actor's snapshot; then integrate operation-by-operation after staged migration and test verification. See [implementation progress](../internal/IMPLEMENTATION_PROGRESS.md).

## Registry and authorization checkpoint

The registry service accepts only known feature and permission codes, maximum 100 per batch, canonical metadata and deterministic order. It registers new definitions inactive, preserves previous activation/metadata, and never modifies tenant feature or role grant state. Registration retries by code are safe; API-level idempotency keys and workload authorization remain unimplemented. An unknown code, duplicate, missing feature, or conflicting persisted feature association fails closed. Internal registration is not exposed through an endpoint.

The authorization service provides Require and Effective, loads a fresh scoped snapshot using the verified school principal, checks active account and registered capability, and ignores legacy Admin. Both services are constructed by Identity bootstrap but not attached to live routes. The PostgreSQL 17 CI integration gate applies the actual Identity migrations only to an explicitly provisioned fresh disposable database, and verifies rollback, tenant boundaries, feature/permission enforcement and revocation. Next: role ownership/read-write model, safe delegation and approval policy, controlled bootstrap/backfill, then endpoint-by-endpoint enforcement and frontend integration. The temporary tenant-wide policy remains a release blocker.

## Role persistence checkpoint

Added the internal role entity and tenant-scoped repository for atomic, inactive draft creation with fixed personas and single-query detail retrieval. Drafts cannot accept permission grants, activation, or assignments; duplicate normalized names are rejected within a tenant. Added PostgreSQL integration assertions for scope, duplicate names, persona validation and absence of grants. No endpoint or live authorization policy changes. CI verification is pending. Next: verified-principal service, delegation and permission mapping, versioned edits, search, outbox and approval enforcement before route exposure.

## Role service checkpoint

Added internal CreateDraft and Get operations, separate DTOs, validation, response mapping and bootstrap wiring with the existing live authorization service. Tenant and creator come from verified principal only; create requires IDENTITY_ROLE_CREATE and detail requires IDENTITY_ROLE_VIEW. Name and description limits are provisionally 120 and 1000 Unicode code points. Inactive drafts only; no activation, grants, approval or endpoint exposure. Added service and validator Ginkgo tests. CI verification pending. Next: versioned metadata edits and scoped search; then permission delegation and approval/outbox before exposure.

## Versioned role metadata checkpoint

Internal UpdateMetadata requires IDENTITY_ROLE_UPDATE and a positive expected version. It replaces name, description and complete persona set atomically with a tenant/role/version predicate and version increment. Status, grants and assignments are untouched. Missing or foreign IDs are not found; stale versions conflict. Added service/validator tests and PostgreSQL assertions for tenant boundaries, persona replacement and two concurrent edits with a single winner. Full replacement is internal only, not a final PATCH contract. CI pending. Search, outbox, approval and grant management remain outstanding.

## Role search checkpoint

Added internal RoleView-authorized search using the shared core search framework with tenant-scoped totals, name/description search text, role ID/name/status/personas/timestamp filters, stable multi-sort and numbered or legacy cursor pagination. Persona filtering uses tenant-correlated EXISTS; list hydration returns persona sets in the same select, not N+1 queries. Added validator/service tests and real PostgreSQL pagination and tenant isolation assertions. No HTTP route exposed. CI pending. Next: permission delegation, versioned grant replacement, approval and outbox.

## Delegation preview checkpoint

Added tenant-scoped whole-role permission reads and a read-only versioned replacement preview. Effective authority is loaded once per preview; the caller needs mapping management and must hold every current and proposed permission (including removals). Unknown/duplicate codes and reserved registry/feature-management permissions are rejected, even if held. Explicit school delegation allowlist defaults to deny for new capabilities. Any actual change requires approval under the provisional conservative policy; preview is not a grant or approval token. No write endpoint or persistence mutation added. Approval effect must revalidate authority and expected version within a serialized tenant authority transaction, enforce independent review and last-admin safety, and persist outbox/audit before exposure. Added unit and PostgreSQL mapping read assertions; CI pending.

## Transactional draft permission persistence checkpoint

Added a separate, unwired DraftPermissionWriter and transaction-scoped tenant advisory lock. Replaces mappings only for an inactive role with zero account assignments, locks the role row, requires the expected version, validates canonical registered active permission bindings, and atomically writes mapping and version. Unchanged sets do not bump version. Reserved, duplicate, missing or mismatched permissions fail closed; rollback restores both mappings and version. PostgreSQL assertions cover stale version, no-op, unregistered grant, outer rollback, tenant scope and active/assigned roles. CI pending. This is an internal persistence primitive, not an approval or authorization boundary: no service, API or approval effect invokes it. The advisory protocol must be shared by future authority writers; last-administrator protection, approved effect revalidation, outbox/audit and activation remain mandatory before any live mutation. No production migration or grant changes performed.

## Authority protection checkpoint

Added an internal last-administrator guard using the existing transaction-scoped tenant advisory lock, with pre/post checks and rollback on zero remaining effective administrators. Provisional recovery policy requires both IDENTITY_ROLE_PERMISSION_MANAGE and IDENTITY_ACCOUNT_ROLE_ASSIGN on the same active account through active tenant roles, registered active permissions and the enabled RoleManagement feature. Union across that account's roles is valid; union across different accounts is not. SQL count is tenant-scoped and does not trust personas or legacy admin flags. Added unit tests and PostgreSQL assertions for concurrent account suspension, split grants, disabled entitlements and rollback. No live mutation path invokes the guard yet. All authority writers, including account lifecycle, must adopt the protocol; global capability lifecycle needs coordinated cross-tenant locking. Bootstrap, role lifecycle, assignment validation, approved execution and transactional audit/outbox remain outstanding. No production changes. CI pending.

## Role lifecycle checkpoint

Added isolated versioned ACTIVE/INACTIVE transition persistence and an internal lifecycle operation guarded by the tenant authority invariant. Both transitions reject assigned roles, stale versions, foreign IDs and repeated state requests. Activation requires at least one persona and grant, canonical active registered school-delegable permissions and enabled associated tenant features. These are provisional conservative rules; deactivation with assignments needs an approved reassignment workflow, and bootstrap is separate. Added unit and PostgreSQL lifecycle assertions. The lifecycle operation is not wired into the service graph or HTTP: approval-bound execution, live requester/approver delegation checks, transactional audit/outbox, account assignment eligibility and shared global capability locking remain mandatory before exposure. CI pending; no production mutations.

## Account role assignment read and eligibility checkpoint

Added normalized account-role reads and batched role candidate lookup, scoped by verified tenant in the internal service. Preview requires IDENTITY_ACCOUNT_ROLE_ASSIGN, checks a positive account version, one primary role, no duplicates, eligible target account status (ACTIVE or PENDING_SETUP), active nonempty proposed roles and effective delegation authority over every current and proposed role grant. Returns account and role versions for later approval revalidation. Role reads require IDENTITY_ACCOUNT_ROLE_VIEW. Provisional limit of 20 secondary roles; all changes require approval. No normalized or legacy account writes; no service graph or route wiring yet. Added pure validator/preview tests and PostgreSQL read assertions. Atomic assignment replacement, legacy backfill/cutover, explicit approved-effect revalidation, last-admin guard and outbox remain required. CI pending.

## Atomic account assignment checkpoint

Added an isolated normalized account-role writer and guarded internal operation. It acquires the tenant authority lock, locks the account and candidate roles, requires exact account and old/new role version stamps, validates eligible active targets and atomically replaces primary/secondary rows while incrementing the account version. No-op preserves the version. The last-administrator guard rolls back unsafe changes. Added validator and PostgreSQL assertions for primary/secondary swap, stale account/role stamps, no-op and last-admin rollback. The writer changes normalized assignments only and intentionally does not update legacy role columns or search projections: migration/backfill and cutover compatibility must be implemented before wiring it to any route. Approved execution must revalidate requester/approver effective authority, feature eligibility, payload and versions under the same lock and enqueue audit/outbox atomically. No live service or route invokes the writer. CI pending.

## Legacy migration planning checkpoint

Added a read-only, repeatable-read tenant inventory and pure deterministic migration planner. Maps six fixed persona codes to tenant-specific UUIDv5 IDs in a dedicated namespace; preserves assignment type, account version and reports accounts without roles. Rejects malformed legacy data, unknown/duplicate codes, any existing normalized assignment, or role ID/name collisions. No writes, permission grants, role activation, legacy cutover or admin bootstrap. UUIDv5 is an explicit exception for deterministic migration IDs: existing role-search cursor validation expects UUIDv7 and must be adapted before using migrated IDs. Existing operator bootstrap still relies on the legacy admin allowlist and is unchanged. Before applying any migration, add operator confirmation, same-transaction inventory revalidation, explicit grant manifest and administrator bootstrap, parity checks, durable audit, maintenance lock across legacy writers, and rollback/cutover procedure. Unit and PostgreSQL inventory assertions added; CI pending.

## Offline transactional backfill checkpoint

Added an unwired, maintenance-only backfill primitive. Acquires write-blocking table locks before rebuilding the reviewed plan; rejects stale plans or prior normalized assignments, creates only INACTIVE roles and normalized primary/secondary rows, checks exact assignment parity, unchanged legacy account projections/version, expected row counts, persona mappings and absence of grants before commit. Failure rolls back the entire transaction. Reapplication is intentionally rejected, not silently merged. UUIDv5 role IDs now accepted by role cursor validation alongside UUIDv7. Legacy columns remain unchanged and no effective permissions are created. Tests cover stale plan, successful inert backfill, repeated apply conflict and no effective admin. CI pending. The operator entrypoint, grant manifest, bootstrap/activation of assigned migration roles (special controlled path), audit record and runtime cutover are still outstanding. Run only with application writers stopped and a backup; table locks are database-wide for these tables, not tenant-local. No production backfill executed.

## Controlled bootstrap checkpoint

Fixed Phase 11B CI compile failure: inventory hydration now uses PrimaryCode. Added explicit reviewed grant manifest and offline bootstrap primitive, which revalidates original migration plan and inert backfill parity under maintenance table locks, validates registered active school-delegable grants and already-enabled tenant entitlements, activates only explicitly listed roles, and requires an active nominated account and exactly one effective administrator at commit. Recovery grants can only be placed on roles exclusively assigned to the nominee, including checks for other inactive holders; shared legacy administrator roles need a separately reviewed dedicated role workflow. No implicit registry/entitlement changes. All writes roll back on failure. Added unit and PostgreSQL assertions for missing entitlement, explicit bootstrap and repeat refusal. No CLI, durable audit/outbox, gateway cutover or production bootstrap yet. Both backfill and bootstrap require offline maintenance; CI pending.

## Operator migration CLI checkpoint

Added a dedicated root backend command and script with read-only plan output and explicit reviewed-file backfill/bootstrap commands. Writes require exact tenant/database, operator UUID, SHA-256 of reviewed bytes and offline-maintenance attestation. Added an additive receipt table and same-transaction receipt write, recording operation, tenant, operator, input digest and result. No production execution or runtime wiring. The runbook is RBAC_MIGRATION_RUNBOOK.md. Operator identity is supplied by the trusted database operator, not independently authenticated by CLI; writers must actually be stopped, and the expected database name is not a substitute for checking server identity. CI pending. Remaining: legacy compatibility, explicit permission enforcement readiness/cutover, approval and live audit/outbox.

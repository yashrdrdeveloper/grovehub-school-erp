# Internal platform — execution progress

Status: IN PROGRESS. See [implementation plan](INTERNAL_PLATFORM_IMPLEMENTATION_PLAN.md) and [school RBAC plan](../users/RBAC_IMPLEMENTATION_PLAN.md).

## A0 — Baseline and dependency inventory

Source baseline: `b24c0e627538ad32795087d0cbba783d0f7f4916`. The latest source-changing frontend CI [run 45800103347](https://github.com/yashrdrdeveloper/grovehub-school-erp/actions/runs/45800103347) failed at check:architecture before typecheck, Jest, build, browser and lint. The first failing gate reports multiple describe blocks in Calendar navigation, DateRange calendar, and table scrollbar geometry helper tests. The latest earlier backend annotation run 45797793768 passed at a different commit; not evidence for the current revision.

Tracked file inventory from the complete Git tree:

| Candidate | Files | Status |
|---|---:|---|
| frontend/packages/base | 81 | Requires complete dependency audit |
| frontend/packages/components | 470 | RDR-specific SchoolMonogram imports @Assets/brand |
| frontend/packages/styles | 10 | Theme/global style ownership audit required |
| backend/libs/core | 169 | Application-specific gateway cookie/principal fields and generator root assumptions require review |

Vite, TS, Jest, package scripts and CI use frontend-relative paths. Core annotation generation assumes backend/services and backend/libs/core. Do not move them before the boundary audit and passing baseline. Direct local checkout unavailable in this execution environment (GitHub DNS resolution failed); GitHub source access works through the connector.

## A0.1 — First baseline fixes

- Split the three failing helper suites into one describe per named export, retaining all assertions and actual pure helper imports. No test gate or allowlist relaxation.
- Remove Vite tunnel host ternaries, preserving the same exact-host security behavior.
- New CI validation required; these changes are not evidence that later checks pass.

## Next gates

1. Recheck CI at the source-changing remediation commit and repair remaining baseline failures without weakening checks.
2. Finish import/asset/CSS/storage and Go package audits; record every package disposition. Extract school-specific UI from reusable Components.
3. Move Base, Styles and Components, migrate aliases, build and test consumers, then Core only when neutral.
4. Follow the school RBAC plan and subsequent internal foundation, approval, registry, execution and UI phases.

No library extraction, schema migration, service implementation or publishing has occurred yet. This log must be updated with exact commit and actual test results at each step.

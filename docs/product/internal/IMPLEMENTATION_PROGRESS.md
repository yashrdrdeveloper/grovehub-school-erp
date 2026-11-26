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

## A0.2 — Tooltip automation contract

CI [45804844951](https://github.com/yashrdrdeveloper/grovehub-school-erp/actions/runs/45804844951) passed boundary and unit structure checks, then failed Tooltip's automation contract. Added typed dataTestId propagation through the trigger Box and a regression test. Remaining gates and tests are not yet verified.

## A0.3 — Calendar and Table conditional expressions

CI [45807088192](https://github.com/yashrdrdeveloper/grovehub-school-erp/actions/runs/45807088192) verified boundaries, units, UI automation contracts, styles and layout/assets, then rejected Calendar and Table ternaries. Replaced the reported branches without changing rendering semantics. Remaining checks require CI verification.

## A1 — Base source migration

Inspected all 80 TypeScript source/test files in the 81-file Base package: imports are confined to relative Base modules and lodash-es, clsx and tailwind-merge. School-name strings appeared only in two generic fixture tests, not runtime logic. Moved the 81 tracked files to libs/frontend/base (no source duplication), with a private workspace package manifest and independent TypeScript check. Added a root pnpm workspace/lockfile, preserved existing dependency resolutions, rewired school Vite/TypeScript/Jest/format/lint and CI, and adapted owner checks to retain coverage of extracted files. Existing @Base aliases are retained as source compatibility during the workspace transition. Private registry publishing and production ESM/declaration packaging remain future work. CI and runtime verification are pending; the prior baseline still had later unverified checks.

CI [45807842172](https://github.com/yashrdrdeveloper/grovehub-school-erp/actions/runs/45807842172) confirmed frozen workspace installation and retained source counts (1,222 boundary-checked files, 42 operation units). It reached the same pre-existing missing User Setup renderer test gate as before the Base move. Added tests for mobile row, status, username/email and both timestamps. Full CI remains pending.

## A2 — Styles extraction

Moved the ten original styles files to libs/frontend/styles, isolated school-only public-site palette and RDR legacy aliases in frontend/app/theme, and updated the workspace and application stylesheet. CI [45808534498](https://github.com/yashrdrdeveloper/grovehub-school-erp/actions/runs/45808534498) passed dependency installation, boundary/unit/UI/style checks (113 light/dark tokens), and stopped at Prettier with 141 files requiring formatting. Added read-only dependency inventory for Components and Core and independent CI typecheck/test/lint diagnostics on failure; no gate was relaxed. Components and Core have not been relocated. School RBAC and internal services remain unimplemented.

## R1 — School RBAC foundation (not enforced)

Verified the existing account tenant column is TEXT and account ID UUID. Added the six fixed personas as the single source for the legacy role catalog, explicit permission and feature vocabularies, a fail-closed pure evaluator and unit tests. Added an additive normalized RBAC schema migration with tenant-safe composite foreign keys, unique primary assignments and inactive/disabled defaults. No grants are seeded, existing role columns are untouched, and the evaluator is not yet wired into the current gateway or service operations. Repository hydration, API/approval workflows, migration backfill, controlled administrator bootstrap and full authorization enforcement are still required. The temporary permissive account access release blocker remains. Migration and tests have not been run here.

### Verification and follow-up

- Frontend run [45808791918](https://github.com/yashrdrdeveloper/grovehub-school-erp/actions/runs/45808791918): frozen install, shared Base typecheck, ESLint and all 191 Jest suites (356 tests) passed. Overall run failed formatting (142 files) and school typechecking; build and browser checks did not run.
- Corrected all reported unchecked-index diagnostics in DateRange and DragList tests, filter validation and selection, and User Setup date parsing. Reverification is required on the new commit.
- Full extraction inventory: 473 component sources, 165 Core Go sources, 132 static findings. Components retains a school-brand dependency and theme preference key; Core extraction remains conditional on its domain-neutral dependency audit.
- RBAC schema and pure authorization evaluator are foundations only: no migration has been applied, no controller or service has been switched to the evaluator, and no permission grants are seeded. Existing temporary authenticated-tenant access remains a production release blocker.

## A0.4 — Reverified baseline and exact formatter remediation

Frontend [45809337646](https://github.com/yashrdrdeveloper/grovehub-school-erp/actions/runs/45809337646) at `4fe0cdae38373f849d6344cdc2562b04a14d182c` confirmed frozen install, shared Base and school TypeScript checks, ESLint, and all 191 Jest suites (356 tests) passed. The architecture gate still failed at formatting; build and browser verification were skipped. Added a CI failure diagnostic to generate and upload the exact pinned Prettier diff, restoring the worktree before subsequent validation. The existing formatting and architecture gates remain unchanged; the patch must be reviewed and applied separately. No automatic commit or publishing is enabled.

## R1.1 — Permission feature binding integrity

The pure evaluator now rejects any mismatch between the code-owned feature binding and registered state, including unexpected feature associations on featureless permissions. Added regression assertions for valid featureless grants and mismatched bindings. No grants, migrations, or route enforcement were enabled. Backend tests need verification at the new revision.

## A0.5 — Full formatter remediation and verification

Applied the four Base blobs from the first CI artifact, matching every generated blob hash. A later run revealed remaining school frontend formatting debt. On isolated branch `maintenance/frontend-format-073c903`, [workflow 45810539101](https://github.com/yashrdrdeveloper/grovehub-school-erp/actions/runs/45810539101) ran frozen pnpm installation, the pinned formatter, format:check and git diff --check successfully. Verified all 144 resulting changed files were tracked frontend/library files and transferred the formatter tree to the implementation branch at `09e7eb784c1b4a89516c5ccd02cd2c6afcf94b0f`, excluding the temporary workflow. No architecture, test or lint gate has been relaxed. Full normal frontend CI on the new tree remains required; passing formatting alone is not a release signal.

## A3.1 — Theme namespace preparation

Commit `073c9030ca02cb43da3c26aee64b806220c3d61a` allows Components theme storage readers, writers and hook to receive an explicit stable application storage key, with two isolation tests. The legacy school default is retained temporarily to preserve current preferences. Still required: move the school default into the school owner, inject the key consistently in all consumers, audit table preference keys and remove the school artwork dependency before moving Components. Internal application and extracted Core remain unimplemented.

## R1.2 — Backend verification

[Backend run 45810062545](https://github.com/yashrdrdeveloper/grovehub-school-erp/actions/runs/45810062545) at `2cac9c6a5b25fb4722e3969690e0c5d2c2345aba` passed Linux annotation build, available service/core unit suites, coverage-report generation, and macOS generated-source verification. This does not prove complete coverage, migration execution, production authorization or end-to-end RBAC. Repository hydration and controlled bootstrap are prerequisites for enforcement; interim permissive account authorization remains a release blocker.

## A3.2 — School branding and preference boundaries

[59c5e18](https://github.com/yashrdrdeveloper/grovehub-school-erp/commit/59c5e185f81b16a17a66f942b543e9568e16fa68) relocates the RDR monogram into `frontend/packages/domain/brand`. Updated public site, auth and topbar consumers; Components can no longer import app-owned Assets. The UI contract checker now verifies the school-owned selector through Box. [42a8153](https://github.com/yashrdrdeveloper/grovehub-school-erp/commit/42a8153f5f235b4aea2d1d534197498ad5c6eb5c) relocates the original three-path artwork assertion to the domain test, removing the shared renderer's dependency on school branding. [edbe55c](https://github.com/yashrdrdeveloper/grovehub-school-erp/commit/edbe55c6fd172c5dbe8000a871eaef80239618b8) makes application theme storage keys mandatory and moves the preserved school storage key into `@Domain/preferences`; application bootstrap, login and shell supply it explicitly. Table preferences already accept a caller-owned storage key and have no hardcoded product namespace in their read/write paths.

## A0.6 — Remaining architecture debt

[Run 45810667801](https://github.com/yashrdrdeveloper/grovehub-school-erp/actions/runs/45810667801) passed the formatting and earlier checks, then surfaced 20 audit violations (15 LARGE_UNIT, 5 NATIVE_COLLECTION). [ef8c226](https://github.com/yashrdrdeveloper/grovehub-school-erp/commit/ef8c226d35fb8c28568da926a0b25a428657a8c0) replaces the five native collection usages with Base utilities; [5f56eab](https://github.com/yashrdrdeveloper/grovehub-school-erp/commit/5f56eab0acbd2b105ba66e4a54a42be67ccbbadd) adds the mobile pagination deduplication regression test. [Run 45811360922](https://github.com/yashrdrdeveloper/grovehub-school-erp/actions/runs/45811360922) passed typechecks, ESLint and 193 Jest suites / 359 tests after brand test relocation, but failed one formatting check. Applied the exact generated formatter blob at [88870ca](https://github.com/yashrdrdeveloper/grovehub-school-erp/commit/88870cada3607b3dc8d7db4f22e6c01da8a00374). Full combined CI [45811529676](https://github.com/yashrdrdeveloper/grovehub-school-erp/actions/runs/45811529676) must still be checked. No allowlist or size threshold changes. Large units remain to be decomposed, followed by actual shared Components workspace move, Core boundary work, RBAC enforcement and internal app implementation.

## A0.7 — Responsibility-based oversized-unit remediation (first batch)

Split nine of the fifteen flagged units without changing the 155-line threshold or audit allowlist. Calendar delegates day/month/year selection to CalendarGrid; DateRange delegates preset scrolling and selection to DateRangePresets; Filters separates mobile condition list/editor and desktop actions; SelectionPanel delegates options rendering; Sort delegates footer actions. Table preference/action contracts are separated while preserving existing type exports; UserSetup delegates account details, and UserForm delegates lifecycle actions. Original CSS classes, selectors, callback ownership, API behavior and interaction semantics are retained. Commits: `8de7572`, `caf11b7`, `5b9c76e`.

[CI 45818478975](https://github.com/yashrdrdeveloper/grovehub-school-erp/actions/runs/45818478975) passed Base and app typechecks, ESLint and 193 Jest suites / 359 tests, but reported six formatting deviations. Applied the exact CI formatter artifact (all six resulting blob hashes matched) in `2b41f521acaa35090f373c72a7aa0c058bca3874`. [Combined CI 45818660603](https://github.com/yashrdrdeveloper/grovehub-school-erp/actions/runs/45818660603) is pending at the time of this record. Remaining expected oversized units: Table, MobileTableToolbar, useTable, DesktopTable, MobileTable and useUserSetup. No production build or browser smoke result yet. Components physical move, Core separation, RBAC enforcement and the internal platform remain outstanding.

## A0.8 — Remaining six oversized units

[CI 45818660603](https://github.com/yashrdrdeveloper/grovehub-school-erp/actions/runs/45818660603) confirmed exactly six remaining LARGE_UNIT findings, with 193 Jest suites / 359 tests passing. Refactored Table into desktop toolbar and column option calculation, MobileTableToolbar into mobile options, useTable into a dedicated column preference hook, DesktopTable into head and body components, MobileTable into row presentation, and useUserSetup into query and pagination units. Commits: `9de043b`, `4fd3f79`, `261891c`. All six original units and new units are below or at the existing 155-line threshold after formatting. [CI 45819258666](https://github.com/yashrdrdeveloper/grovehub-school-erp/actions/runs/45819258666) passed Base and school typechecks, ESLint and all 193 suites / 359 tests, with six formatter deviations. Applied the CI-generated patch with all six blob hashes verified in `81e65ddabc026e5d72e77be155d1a920628210e5`. [Combined CI 45819428622](https://github.com/yashrdrdeveloper/grovehub-school-erp/actions/runs/45819428622) remains the required gate for the complete architecture check, build, browser smoke and SEO verification. No audit exemptions or thresholds changed. Local Git checkout remains unavailable due to DNS, retested this turn. Physical shared Components extraction remains the next phase after baseline validation.

## A0.9 — Frontend baseline fully green

All fifteen oversized units were decomposed and all five native-collection findings removed; the guideline regression gate passed without threshold or allowlist changes. The final CSS Module check rejected dynamic HSL in Avatar. Moved the same hue-dependent gradient and text formulas to the shared Styles utility on the avatar element and added two regressions covering the hue, compact mode, initials and image failure (`1afc425`); applied the exact CI formatter blob (`89fce33`). Made UI, layout and ownership scripts resolve the physical Components owner and added a relocation contract test (`fb4a265`). [Baseline CI 45826161838](https://github.com/yashrdrdeveloper/grovehub-school-erp/actions/runs/45826161838) completed successfully: frozen install, all architecture gates, Base and app typechecks, 194 Jest suites / 361 tests, production Vite/SEO build, Chromium mobile/accessibility/theme smoke, SEO artifacts and ESLint. The audit still reports 26 files / 33 historical findings; these remain migration debt, not exemptions for new work.

## A3.3 — Components workspace extraction

[Isolated migration run 45826276734](https://github.com/yashrdrdeveloper/grovehub-school-erp/actions/runs/45826276734) successfully relocated 515 original Components files, the global i18n type declaration and two new package configuration files (518 files total). All original Components file paths were matched; only the runbook and primitive CSS changed during the relocation. No old physical Components folder remains. The primitives now use shared semantic surface and muted tokens. Private source-workspace `@groovehub/components` declares Base/Styles, Radix/TanStack/CVA/Lucide dependencies and React/ReactDOM/Router/React Query peers. The existing dependency package resolutions and snapshots in pnpm-lock.yaml are byte-identical. Frozen install, standalone Components typecheck, school/Base typechecks, all architecture gates, 194 suites / 361 tests, ESLint and production build passed on candidate `71f13b34e944e81cfc63ef9896296b40d8248b83`. Transferred the verified candidate onto `phase-2.1-identity-account` in `67bee65b577bbbbb387718c5870e2b59b128440a`, removed the one-time workflow from that tree and added standalone Components typechecking to normal CI. Normal post-move browser/SEO smoke still requires verification on the implementation branch. Existing `@Components`, `@Base` and `@Styles` aliases are retained during the private-source workspace stage; publishing compiled ESM/declarations and fully independent consumer packaging are not complete. Temporary migration branch retained for traceability; no migration workflow is present on the implementation branch. Next: conditional Go Core package-by-package audit and decoupling, and school RBAC repository hydration/enforcement without delaying the latter behind UI work.

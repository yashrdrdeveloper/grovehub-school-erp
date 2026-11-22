# Frontend refactor migration — Revision 5

> **Status:** Planned code changes. This document inventories architectural debt and defines acceptance checks. It does not claim that the source code has already been migrated to Revision 5.
>
> **Authority:** [Frontend Code Guidelines](./FRONTEND_CODE_GUIDELINES.md). This is a migration checklist, not a second specification.

## 1. Baseline

Existing code on `phase-2.1-identity-account` has Tailwind v4, shadcn-compatible shared primitives, root domain APIs, a shared `packages/styles` design-token package, basic architecture checks and security-sensitive frontend regression tests. Revision 5 now requires a **new** flat helper/constant architecture and more rigorous selector/import/no-ternary enforcement. Passing current CI does **not** prove Revision 5 conformance.

## 2. Repository inventory to inspect in code migration

| Area | Known current debt | Target |
|---|---|---|
| `domain/auth/pages/login/config/` | `index.ts`, `login.config.ts`, `loginViewport.constant.ts` | Remove generic config; split `constants/login.fields.ts`, `login.copy.ts`, `login.routes.ts`, `login.testIds.ts`, validation patterns, viewport and view metadata |
| `domain/auth/pages/login/helpers/` | Multiple nested one-helper folders with their own `index.ts` | Flat named `login.<useCase>.ts` helpers, sibling `types/` and `__tests__/`, no barrels |
| `domain/auth/pages/login/hooks/useLoginForm/` | State separated from peer actions/submission hooks | Compose focused nested subhooks and return one form interface. No duplicate `isSubmitting` |
| `domain/auth/pages/login/hooks/useLoginViewport/` | Extra helper folders with per-helper barrel files | Flat hook helpers, media query and browser event constants; only retain hook wrapper if useful |
| `domain/auth/pages/login/components/LoginForm/` | Render-time variant branching and inline markup; test beside implementation | Pure typed view-model helpers and constants, `@Components` primitives, test under component `__tests__/`, stable page test IDs |
| `domain/auth/pages/hooks/useAuthViewportFit/` | Local test not in `__tests__/` and mixed viewport concerns | Scoped helpers/constants/event listener cleanup and colocated hook tests |
| `domain/auth/pages/recovery/hooks/useRecovery/` | One hook coordinates link verification, grant status, editing, saving, retry and browser listeners | Distinct form/link/grant/submission/retry hooks with narrow contracts and one coordinator; keep server-authoritative failures |
| `domain/auth/pages/recovery/helpers/` | Nested helper units with `index.ts` | Flat helper files, tests/types adjacent in sibling folders |
| `domain/auth/pages/forgotPassword/` | Nested helpers and constants, submission/render-time choices | Apply same flat and nested-hook standard |
| `packages/components/` | Generic `Form`, `Text`, button/card/input variants exist; other components and domains still have styling duplication | Curated icon API, `useBoolean`, deterministic selectors passed to native DOM, all controls using shared design tokens |
| `packages/styles/` | Theme tokens, typography, spacing, breakpoints, shadows and Tailwind theme are now central; app still has large legacy screen CSS | Keep root light/dark tokens centralized and progressively move screen CSS to owning components/domains; browser visual regression required |
| `packages/base/` | Only some generic utility leaf exports exist | Supply/test exact `_trim`, `_find`, `_size`, `_keys`, `_map`, `_isUndefined`, `_isFunction`, etc. as needed; browser event names |
| `domain/dashboard/` and `domain/publicSite/` | Styling/JSX still contains inline choices and native controls | Common selectors, shared primitives and view helpers without layout regression |
| `scripts/check-*.mjs` | Existing checks enforce only some Revision 4 rules | Add exact static enforcement for Revision 5 as each module migrates |

The table is a known-first audit, not a complete proof that all other files are clean. Search the entire `frontend/` tree for unused helpers/constants, stale imports and violations before deleting anything.

## 3. Implementation batches and gates

1. **Infrastructure:** base leaf Lodash exports/tests; `@Components/icons`; shared `useBoolean` and test IDs; browser event constant map; static checks for direct vendor imports.
2. **Login:** rewrite constants; flatten page and hook-local helpers; replace wrapper normalizers with `_trim`; consolidate state; nested subhooks; view-model selection; use shared controls; tests.
3. **Auth remainder:** shared viewport fit, recovery, forgot-password and guards; preserve user messaging, retry limits and security behavior.
4. **Other domains:** dashboard, public school site, shell, business/remaining domains; no dead code, retained routes and responsive behavior.
5. **Style ownership:** keep `packages/styles` as the source for theme/type/spacing tokens; move remaining old app CSS to component/domain owners without changing auth/public page behavior.
6. **Full enforcement:** no ternary/switch dispatch, typed IDs for page-owned DOM, no direct Lucide/Lodash outside owners, no deprecated helper layout; add browser UI regression coverage.

After every batch, run frozen install, architecture checks, TypeScript, Jest, lint and build/prerender. Run real viewport tests before signing off styling or selector changes.

## 4. No misleading completion claims

**Docs first:** this documentation change updates standards and migration planning only. It neither moves helper files nor installs new utilities/hooks. Each later implementation commit must reference the rules it satisfies and record its actual CI results. A non-executed browser test must never be reported as passing.

**Done** means architecture, feature behavior, security states, accessibility, deterministic selectors and tooling all pass; then remove this area's debt from the table.

# Frontend Code Guidelines

> **Status:** Mandatory, authoritative frontend engineering standard — **Revision 7 (November 2026)**.
>
> **Scope:** Every source file under `frontend/`, including `app`, `packages/base`, `packages/business`, `packages/components`, `packages/styles`, `packages/assets`, `packages/domain`, `packages/shell`, unit tests, integration tests, styles, generated code boundaries, and scripts.
>
> **Stack:** React 19, TypeScript strict, pnpm, Vite, Tailwind CSS 4, shadcn-compatible Radix primitives, Lucide icons through the components package, Jest, and the existing project architecture checks.
>
> **Precedence:** Revision 7 sections 23–35 are the newest, most specific acceptance rules. They override older examples where they conflict. This document replaces previous frontend guideline revisions and examples. Existing code that differs is **migration debt**, not an exception to this standard. Do not claim implementation is complete until checks and behavior regressions pass.

## 1. Non-negotiable principles

SRP, DRY, OCP, explicit ownership, composition, type safety, stable accessibility, automation selectors, domain isolation and predictable import boundaries are mandatory. A file has one discoverable purpose. A hook does not become a hiding place for a giant component. A helper must not restate an existing library operation without adding real semantics. Do not introduce a directory, abstraction, wrapper, factory, or duplicated state merely for the appearance of architecture.

**Required working sequence:** identify the owning unit; reuse an existing shared primitive or utility; define types and constants; implement focused pure helpers or hooks; compose a presentational component; add the correct level of tests; delete unused code; run all applicable gates.

**Prefer explicit behavior over clever syntax.** No ternary operator in application TypeScript/TSX, including JSX props and nested conditional rendering. No multi-branch `switch` for behavior selection. For two cases, use readable `if / else` in a named helper. For three or more variants of the **same behavior**, use an exhaustive typed strategy registry keyed by a named value/enum.

## 2. Required package ownership and dependency direction

~~~text
frontend/
├── app/                       # bootstrap only; mounts the shell
├── packages/
│   ├── base/                  # pure, framework-neutral technical utilities/constants
│   ├── business/              # reusable business contracts, API transport and operations
│   ├── components/            # reusable accessible UI, icons, hooks and component recipes
│   ├── styles/                # global design tokens, themes, typography and Tailwind
│   ├── assets/                # SVGs, images and illustrative assets by domain
│   ├── domain/
│   │   ├── auth/              # authentication pages and session-aware UI
│   │   ├── dashboard/         # school home page
│   │   ├── publicSite/        # school website, enquiry UX and public pages
│   │   └── <domain>/
│   └── shell/                 # routing, providers and application/workspace composition
├── components.json            # shadcn configuration
├── scripts/
└── package.json
~~~

The dependency graph is `app → shell → domain / components / styles / base` and `domain → business / components / styles / base`; `business → base`; `components → assets / styles / base`; `assets → nothing higher`; `styles → nothing higher`; `base → nothing higher`. Business and components are independent of domain and shell. A domain never imports shell or another domain's private internals. Shell imports domain **root APIs only**, such as `@Domain/auth`. Type-only imports obey the same boundaries. API calls must not originate from the shell or reusable presentation components.

Use **the exact case-sensitive aliases actually declared in tsconfig/Vite**: `@Base`, `@Business`, `@Components`, `@Domain`, `@Styles`, `@Assets`, `@Shell`. No obsolete `@Auth` or `@LandingPage` alias. Within a unit use short relative imports; outside the unit use its intentional public boundary.

`frontend/app` owns no business screens. `packages/base` must not import React, React Router, the DOM, business entities, or UI. Base browser-event **name constants** are pure; adding actual event listeners belongs in a browser-aware component/domain hook.

## 3. Canonical directory layout: one page, one composition root

~~~text
packages/domain/auth/
├── pages/
│   └── login/
│       ├── Login.tsx                     # presentation only
│       ├── index.ts                      # page public entry
│       ├── containers/
│       │   └── loginContainer/
│       │       ├── LoginContainer.tsx
│       │       ├── types/
│       │       └── index.ts
│       ├── components/
│       │   └── LoginForm/
│       │       ├── LoginForm.tsx
│       │       ├── types/loginForm.props.ts
│       │       ├── helpers/               # only if LoginForm owns pure view logic
│       │       ├── constants/             # only if LoginForm owns fixed values
│       │       ├── __tests__/LoginForm.test.tsx
│       │       └── index.ts
│       ├── templates/
│       │   ├── desktopLogin/
│       │   └── mobileLogin/
│       ├── hooks/
│       │   ├── useLoginForm/
│       │   │   ├── useLoginForm.ts
│       │   │   ├── types/useLoginForm.ts
│       │   │   ├── hooks/
│       │   │   │   ├── useFormActions/
│       │   │   │   └── useFormSubmission/
│       │   │   ├── helpers/               # flat use-case helper files
│       │   │   │   ├── login.validateCredentials.ts
│       │   │   │   ├── types/login.validateCredentials.ts
│       │   │   │   └── __tests__/login.validateCredentials.test.ts
│       │   │   ├── constants/
│       │   │   └── index.ts
│       │   └── useLoginViewport/
│       │       ├── useLoginViewport.ts    # only if a reusable hook is warranted
│       │       ├── helpers/login.viewport.ts
│       │       ├── constants/
│       │       ├── types/
│       │       └── index.ts
│       ├── helpers/                       # page-wide pure operations only
│       │   ├── login.getFormViewModel.ts
│       │   ├── types/login.getFormViewModel.ts
│       │   └── __tests__/login.getFormViewModel.test.ts
│       ├── constants/
│       │   ├── login.fields.ts
│       │   ├── login.copy.ts
│       │   ├── login.routes.ts
│       │   ├── login.testIds.ts
│       │   ├── login.passwordVisibility.ts
│       │   └── types/                     # only for genuinely shared contracts
│       └── types/
│           ├── login.props.ts
│           └── login.state.ts
└── index.ts
~~~

For a **migrated page root**, only `<Page>.tsx` and `index.ts` are permitted as runtime code. Containers, templates, hooks, HOCs and visual subunits belong under their owning folders. Preserve one active form in the DOM across mobile/desktop variants; do not mount duplicate hidden credential forms.

For `packages/shell`, retain `Shell.tsx` and `index.ts` at root; container/provider/registry/workspace view units belong in named subfolders.

**Only create a child `helpers/`, `constants/`, `hooks/`, `hocs/`, `components/` or `types/` folder when its owner actually needs it.** Hook-owned helpers stay next to that hook; component-owned visual helpers stay next to that component; page-wide helpers live at page level. Never duplicate the same operation at page level and hook level.

## 4. Public entrypoints: precisely where `index.ts` is allowed

A unit's **public boundary** (domain, page, component, hook, HOC, container, service) can have an intentional `index.ts` when callers need that API. These exports must be explicit; do not use catch-all re-export barrels.

**`helpers/` and `constants/` have a stricter rule: no `index.ts`, no one-folder-per-helper.** Their implementations live directly in their owning `helpers/` or `constants/` folder. Import the exact file, for example `../helpers/login.getFormViewModel` or `../constants/login.copy`. A `types/` folder may contain multiple adjacent named contracts. Tests belong in the one flat sibling `__tests__/` folder.

~~~text
helpers/
├── login.getFormViewModel.ts
├── login.getPasswordViewModel.ts
├── login.getFailureMessage.ts
├── types/
│   ├── login.getFormViewModel.ts
│   ├── login.getPasswordViewModel.ts
│   └── login.getFailureMessage.ts
└── __tests__/
    ├── login.getFormViewModel.test.ts
    ├── login.getPasswordViewModel.test.ts
    └── login.getFailureMessage.test.ts

constants/
├── login.fields.ts
├── login.copy.ts
├── login.routes.ts
├── login.testIds.ts
├── login.passwordVisibility.ts
└── types/
~~~

Forbidden in a migrated helpers/constants directory: `helpers/index.ts`, `constants/index.ts`, `helpers/getSomething/getSomething.ts`, `constants/someConstant/index.ts`, `helpers.ts`, `common.ts`, `utils.ts`, `misc.ts`.

A self-contained business operation or reusable component is different: its own public `index.ts` can remain. Do not delete **domain/page/component/hook** public entrypoints while removing helper/constant barrels.

## 5. File naming, types, placement and unused-code removal

Use filenames that identify owner **and** use case: `login.getPasswordViewModel.ts`, `login.validateCredentials.ts`, `login.fieldIds.ts`, `recovery.classifyFailure.ts`, `person.request.mapper.ts`. Avoid anonymous `getData.ts` or `helpers.ts`. Exported functions have named callable contracts in the nearest `types/<useCase>.ts`. Reuse the nearest authoritative page/entity types via indexed access, `Pick` and `Omit` instead of duplicating state shapes.

~~~ts
// hooks/useLoginForm/helpers/types/login.validateCredentials.ts
import type { LoginCredentials } from '../../../../types/login.state';
export type ValidateCredentials = (input: LoginCredentials) => LoginFieldErrors;

// hooks/useLoginForm/helpers/login.validateCredentials.ts
import type { ValidateCredentials } from './types/login.validateCredentials';
export const validateCredentials: ValidateCredentials = (input) => {
  // focused pure validation
};
~~~

Do not use a UI props type as an API transport contract. Shared backend entities live under `business/<domain>/types/baseEntities/`; per-operation service requests/responses stay with the business service. Backend readers should be typed, pure and free from page decisions.

**Delete unused code** only after checking all importers: code search, public barrel exports, test references, route registrations, generated imports, dynamic loading and builds. Remove superseded imports, files, obsolete tests and empty directories in the same migration commit; update documentation and architecture checks for the new layout. Do not remove an indirectly referenced route or public API merely because a text search returns zero import lines.

## 6. Constants, identifiers, copy and configuration

**Literal/static values belong in `constants/`, not a generic `config/` folder.** There is no page-level `config/index.ts`. Do not combine unrelated concerns into one huge `login.config.ts` or `Object.freeze({ ...everything })`. Separate:

- `constants/login.fields.ts` — typed field identity and metadata.
- `constants/login.copy.ts` — user-visible messages/labels.
- `constants/login.routes.ts` — stable route names and paths.
- `constants/login.passwordVisibility.ts` — visibility states and their icon/label/type metadata.
- `constants/login.submitStatus.ts` — ready/submitting labels and icon behaviors.
- `constants/login.testIds.ts` — deterministic automation selectors.
- `constants/login.validation.ts` — only literal regexes, limits, rule identifiers; any predicate/function belongs in a helper.
- `constants/login.viewport.ts` — viewport media query.
- `packages/base/events/constants/browserEvents.ts` — shared browser event names.

Use string-valued `as const` maps, or a real TypeScript `enum` when stable interoperability requires it. Do not repeat raw `'email'`, `'password'`, `'change'`, `'hashchange'` or status strings across the code. Example:

~~~ts
// constants/login.fields.ts
export const LOGIN_FIELD = {
  EMAIL: 'email',
  PASSWORD: 'password',
} as const;

export type LoginField = (typeof LOGIN_FIELD)[keyof typeof LOGIN_FIELD];

export const LOGIN_FIELD_ORDER: readonly LoginField[] = [
  LOGIN_FIELD.EMAIL,
  LOGIN_FIELD.PASSWORD,
];

export const LOGIN_FIELDS = {
  [LOGIN_FIELD.EMAIL]: { label: 'Email address', placeholder: 'Email', autoComplete: 'username' },
  [LOGIN_FIELD.PASSWORD]: { label: 'Password', placeholder: 'Password', autoComplete: 'current-password' },
} as const;
~~~

A `config/` directory is **not** the default and must not be a synonym for constants. Allow it only when a named use case genuinely composes runtime-configurable providers, registries or policy from smaller constants/strategies. An executable calculation belongs in a helper; a static map belongs in constants. Do not add a `config/index.ts` as an import shortcut.

For password visibility and submit labels, use constants for **fixed configurations**, and a pure `login.getPasswordViewModel.ts` / `login.getSubmitViewModel.ts` to select the configuration. Do not compute `VisibilityIcon`, `visibilityLabel`, `passwordType` or `submitLabel` via ternaries inside JSX or the component render function.

## 7. Lodash and `@Base`: the only generic transformation gateway

All reusable Lodash operations **must** be exposed by `packages/base` and imported by exact operation from `@Base`, not from `lodash-es`, `lodash`, or from a broad base barrel that eagerly imports unrelated modules. The imported name carries a leading underscore: `_trim`, `_find`, `_keys`, `_size`, `_map`, `_isUndefined`, `_isFunction`, `_isEmpty`, `_get`, `_includes`, `_isEqual`.

~~~ts
// packages/base/utilities/trim/trim.ts
import trim from 'lodash-es/trim';
export { trim as _trim };

// usage in domain/business code
import { _trim } from '@Base/utilities/trim';
import { _keys } from '@Base/utilities/keys';
import { _size } from '@Base/utilities/size';

const normalizedEmail = _trim(email);
const errorCount = _size(_keys(fieldErrors));
~~~

**Do not create `normalizeLoginEmail.ts` that only returns `email.trim()`, or a custom `isUndefined` wrapper mirroring Lodash.** Call the existing named `@Base` operation directly. A dedicated helper is justified only when it adds real validation, context, return mapping or policy. A function named `_isValid` is **not assumed to exist in Lodash**; use a real domain validation helper or define an intentionally named pure base predicate if its semantics are genuinely generic and tested.

Prefer `@Base` collection/object/string operations consistently for transformation and inspection; avoid scattered native `Object.keys`, `Array.prototype.find/map`, `String.prototype.trim` and duplicated comparison logic when the equivalent base utility is available. This does **not** ban language constructs or platform APIs: use React hooks, promises, `addEventListener`, `URL`, `FormData`, `Date` and native functionality where Lodash has no meaningful equivalent. Never import or expose a huge Lodash namespace: individual small leaf modules preserve tree-shaking and Jest compatibility.

Before using a requested `@Base` operation, implement/export/test it in base; **do not assume a documented name already exists**. `packages/base` may depend on Lodash but may not depend on React or a domain. Business operations must never import Lodash directly.

## 8. Hooks, ownership of state and nested hook composition

**Exactly one owner for each state domain.** For login, `useLoginForm` owns `email`, `password`, field errors, request error, visibility and `isSubmitting`. A separate submission hook may coordinate a request but must not create a second `isSubmitting` state. Pass the typed form model or its minimal selected contract to subordinate hooks.

Reusable shared booleans belong to `@Components/hooks/useBoolean`. Required API:

~~~ts
const [isVisible, show, hide] = useBoolean(false);
show();   // sets true
hide();   // sets false
~~~

The real hook must return stable, memoized action functions. A hook that needs additional toggle semantics may expose a separate named `toggle` action under an intentional contract. Do not replace business flows with many unrelated global booleans; the parent state owner still coordinates the use case.

A complex parent hook may contain **nested hook units**:

~~~text
hooks/
└── useLoginForm/
    ├── useLoginForm.ts
    ├── types/useLoginForm.ts
    ├── hooks/
    │   ├── useFormActions/
    │   │   ├── useFormActions.ts
    │   │   ├── types/useFormActions.ts
    │   │   └── index.ts
    │   └── useFormSubmission/
    │       ├── useFormSubmission.ts
    │       ├── helpers/
    │       │   ├── login.validateCredentials.ts
    │       │   ├── login.getErrorMessage.ts
    │       │   ├── types/
    │       │   └── __tests__/
    │       ├── types/
    │       └── index.ts
    ├── constants/
    └── index.ts
~~~

`useLoginForm` composes these focused subhooks and returns **one cohesive form interface** (state, actions and onSubmit) to the container. Do not have `useLogin` assemble a second copy of the form interface, or scatter duplicate setters across page and component units. Avoid `useCallback`/`useMemo` when they add no stable identity or measurable benefit; favor named handlers and small dependency sets. Hook calls must be unconditional and obey React rules.

The recovery flow **must** be split into focused units: `useRecoveryForm` (single state owner), `useRecoveryLink` (URL proof/exchange lifecycle), `useRecoveryGrant` (server-authoritative grant status), `useRecoverySubmission` (validation and password operation), `useRecoveryRetry` (retry orchestration), and a thin composition hook/container. Move API error classification, proof parsing and validation into the specific hook's flat `helpers/`. Never let a single hook own proof parsing, browser listeners, status verification, form state, submission and retry in one 200-line body. Preserve the backend's two-attempt policy: the client must **not** infer attempt count from transient errors.

The same structure applies to `useAuthViewportFit` and other auth page hooks: pure calculations/event constants inside the hook's own `helpers/` and `constants/`; DOM subscriptions in the hook; tests next to pure helpers. Do not create a generic hook that merely forwards a single hook call and has only one caller. Extract an abstraction when multiple consumers or independent semantics justify it.

## 9. Viewport subscriptions, constants and event listeners

Viewport query strings are **constants**. Snapshot calculations, environment guards and subscription mechanics are **purely focused helpers** (each under one flat `helpers/` directory). A React hook invokes `useSyncExternalStore` directly on those functions.

~~~ts
// hooks/useLoginViewport/constants/login.viewport.ts
export const PHONE_VIEWPORT = '(max-width: 720px)';

// hooks/useLoginViewport/helpers/login.getMobileSnapshot.ts
import { _isFunction } from '@Base/utilities/isFunction';
import { _isUndefined } from '@Base/utilities/isUndefined';

// hooks/useLoginViewport/useLoginViewport.ts
export const useLoginViewport: UseLoginViewport = () => {
  return useSyncExternalStore(subscribeToViewport, getMobileSnapshot, getServerSnapshot);
};
~~~

The above shows responsibility placement, **not** a complete implementation: `_isUndefined` is imported separately from `@Base/utilities/isUndefined` and `_isFunction` from `@Base/utilities/isFunction`. Helpers handle server rendering and unavailable `matchMedia` without throwing. A direct, clear expression is better than chains of named helpers that simply wrap one line. Do not create a second `useLoginViewport` wrapper if the owner can call `useSyncExternalStore` directly and no reuse is needed.

Define browser event names in pure `@Base` constants and reuse everywhere:

~~~ts
export const BROWSER_EVENT = {
  CHANGE: 'change',
  HASH_CHANGE: 'hashchange',
  KEY_DOWN: 'keydown',
  RESIZE: 'resize',
} as const;
~~~

A hook must register/unregister the **same** callback, handle SSR and listener availability, ensure cleanup, and not leak listeners across rerenders. Define proper no-op unsubscribe callbacks in helpers; do not put browser event names or media queries as inline magic strings in JSX.

## 10. Component composition and view-model selection

A shared input, button, form field, dialog, checkbox, select, loader, icon or error-feedback primitive belongs under `@Components`. A page/component owns only feature-specific layout and labels; it must not reconstruct the native input behavior, password-toggle accessibility, submit-loading state or error-row conventions already provided by shared primitives.

**All Lucide imports are centralized in `packages/components/icons/`**. Feature code imports icons through `@Components/icons` (or through a shared icon component), **never directly from `lucide-react`**. The components package can import Lucide internally and publish a curated, typed API. Do not copy SVG paths into domains.

A login form is a presentational component in `pages/login/components/LoginForm/` and is rendered by its page/desktop/mobile templates. The reusable Button, InputControl, form field, visibility action and accessibility behavior should come from `@Components`; use the existing auth CSS classes as transitional styling props rather than reimplementing primitives. Screen-specific branches and string/icon selection live in typed helpers, constants, or strategy registries.

~~~ts
// constants/login.passwordVisibility.ts
export const PASSWORD_VISIBILITY = {
  HIDDEN: 'HIDDEN',
  VISIBLE: 'VISIBLE',
} as const;

export const PASSWORD_VIEW = {
  [PASSWORD_VISIBILITY.HIDDEN]: {
    inputType: 'password',
    toggleLabel: 'Show password',
    icon: 'eye',
  },
  [PASSWORD_VISIBILITY.VISIBLE]: {
    inputType: 'text',
    toggleLabel: 'Hide password',
    icon: 'eye-off',
  },
} as const;

// helpers/login.getPasswordViewModel.ts
export const getPasswordViewModel: GetPasswordViewModel = (isVisible) => {
  if (isVisible) {
    return PASSWORD_VIEW[PASSWORD_VISIBILITY.VISIBLE];
  }
  return PASSWORD_VIEW[PASSWORD_VISIBILITY.HIDDEN];
};
~~~

`LOGIN_COPY`, field definitions, visible/hidden UI, submission text, selected icon and layout variants must not be computed by ternaries in components. A view helper returns the required typed fields. Static variants may be looked up via a typed `Record`; behavior uses a strategy `Map`. Always assign a stable DOM test ID to resulting elements.

## 11. Branching, registries and strategy/factory execution

**Forbidden:** ternary syntax `condition ? yes : no` anywhere in application code; nested conditional props; branching chains of `if / else if` for an extendable set of modes; `switch` for variant dispatch. Guard clauses and a simple `if / else` for a true binary case are allowed.

For three or more variants, use a `Map<Mode, Strategy>`, keyed by a typed enum/constant-derived value. Put the registry next to its owning use case (usually under `constants/<useCase>.strategies.ts` when it is a static map, with functions defined in helpers/services). Do not create a universal strategy registry for unrelated behavior.

~~~ts
export const RECOVERY_MODE = {
  INVITATION: 'INVITATION',
  RESET: 'RESET',
  FORGOT: 'FORGOT',
} as const;
export type RecoveryMode = (typeof RECOVERY_MODE)[keyof typeof RECOVERY_MODE];

type RecoveryStrategy = (proof: string) => Promise<void>;

const strategies = new Map<RecoveryMode, RecoveryStrategy>([
  [RECOVERY_MODE.INVITATION, validateInvitation],
  [RECOVERY_MODE.RESET, validateReset],
  [RECOVERY_MODE.FORGOT, validateForgot],
]);

export const runRecoveryStrategy: RunRecoveryStrategy = (mode, proof) => {
  const strategy = strategies.get(mode);
  if (!strategy) {
    throw new UnsupportedRecoveryModeError(mode);
  }
  return strategy(proof);
};
~~~

A typed `Record<Mode, Configuration>` is preferred for static data/configuration. A typed `Map<Mode, Function>` is preferred for runtime function dispatch. Validate untrusted selector values before dispatch and fail explicitly on unsupported values. Keep business-side dispatch in `@Business` if it is shared business behavior.

In rendering, use focused components, explicit view-model flags, named handlers and preselected component strategies. Do not put case selection or transformations inside JSX; do not replace a ternary with unreadable short-circuit chains or nested immediately invoked functions.

## 12. Typed test selectors: every owned DOM element

**Mandatory prop:** `dataTestId?: string` on every reusable DOM-producing base/component contract (including button, input, form field, icon, dialog, table, text and container). Each shared component forwards it to the **final native DOM element** as `data-testid`. Composite components pass it down to the leaf control; composite surfaces use deterministic suffixed selectors for internal label, hint, error, visibility action, spinner and wrapper. Do **not** output the camelCase prop directly onto DOM.

**Mandatory page contract:** define a centralized page test ID catalog, e.g. `constants/login.testIds.ts`. Every **page-owned rendered DOM element** must carry a stable, meaningful `data-testid`: page root, sections, divs, labels, headings, form, individual fields, actions, feedback/error/status and interactive elements. For shared nested primitives, provide `dataTestId` rather than putting a second selector on a React wrapper. The requirement includes noninteractive owned container elements; decorative SVG paths created by an icon library inherit the parent icon's test ID and do not each receive separate IDs.

~~~ts
export const LOGIN_TEST_IDS = {
  PAGE: 'login-page',
  FORM: 'login-form',
  EMAIL_FIELD: 'login-email',
  PASSWORD_FIELD: 'login-password',
  PASSWORD_VISIBILITY: 'login-password-visibility',
  SUBMIT: 'login-submit',
  FEEDBACK: 'login-feedback',
} as const;

// In a domain component:
<Button dataTestId={LOGIN_TEST_IDS.SUBMIT} type="submit">Sign in</Button>
<InputControl dataTestId={LOGIN_TEST_IDS.EMAIL_FIELD} />

// Native page-owned markup:
<form data-testid={LOGIN_TEST_IDS.FORM}>...</form>
~~~

**Unique and deterministic:** every test ID should be unique among currently rendered elements; repeated list items derive IDs from a stable non-sensitive business key (`... + '-' + record.gid`), never array index, text label or a random/`useId()` value. Names must not expose passwords, email addresses, tokens, phone numbers or other PII. A component rendered more than once must receive distinct caller-provided IDs; don't use a hardcoded default that creates duplicates. Use an owner prefix and role suffix (e.g. `login-submit-spinner`). For a test needing to choose among repeated widgets, it must be able to locate the specific element using page/widget ancestry and a meaningful identifier.

**Accessibility is independent:** semantic role, label, keyboard activation and focus order remain required. Test IDs are for stable automation, not a substitute for accessible names or proper `id`/`htmlFor` relationships. Check every rendered page DOM node (excluding third-party implementation internals) in the migration audit; document justified exceptions.

## 13. Test ownership and regression policy

Place tests beside the **owning unit**, never in an unrelated global bucket.

~~~text
pages/login/components/LoginForm/
├── LoginForm.tsx
├── types/
├── __tests__/
│   └── LoginForm.test.tsx
└── index.ts

pages/login/hooks/useLoginForm/helpers/
├── login.validateCredentials.ts
├── types/
└── __tests__/
    └── login.validateCredentials.test.ts

pages/login/__tests__/
└── login.flow.test.tsx
~~~

- **Pure functions** (helpers, validations, mappers, selectors, view-model selectors, base utilities) require direct Jest tests for success, edge and error cases.
- **Presentational components** receive focused tests for props, callbacks, role, error placement, conditional states, stable test IDs and keyboard focus; component tests live in that component's own `__tests__/`.
- **Hooks** receive tests when they own meaningful state transitions, subscriptions or async orchestration. Do not write tests merely to prove a trivial hook exists.
- **Containers/pages** receive integration-level tests for important user flows; do not mock away the behavior being validated. Preserve security flows, cookie-state transitions, recovery retries, and mobile breakpoint behavior.
- **Shared primitives** have focused tests proving `dataTestId` reaches the native leaf. Avoid snapshots as the only contract.
- **Browser tests** verify all required viewports, zoom, mobile keyboards, OS theme and dark mode, reduced motion, unusually short screens and focus management.

No fabricated API responses/activities in production. Mock/test fixtures live in a separate domain `mock/` entrypoint and are never imported by production code.

## 14. Forms, networking, auth and error states

Frontend validation must run before a network call. A form component only renders controlled values, field error messages and named actions. The owning state hook manages its state, a nested submission hook manages asynchronous sequencing, and `@Business` owns HTTP transport, backend request/response types and idempotency/security logic.

For identity, use HttpOnly cookie sessions. Never store credentials, password reset tokens, bearer proofs or secrets in localStorage/sessionStorage. Clear one-time URL fragments after exchange and use only a non-secret grant marker where required. A transient server/network failure is **not** evidence that an activation attempt was consumed. Server-issued expired/invalid and `IDENTITY_SETUP_ATTEMPTS_EXCEEDED` statuses are authoritative.

Use a typed error classifier helper inside the operation that handles the error. Avoid exposing internal gateway error details to end users. Keep form errors in reserved readable slots and use accessible status/alert messages without layout-jumping.

Async operations require cancellation/stale-request handling, duplicate submit protection, and cleanup of subscriptions/intervals. Define retries explicitly; do not retry non-idempotent mutations blindly.

## 15. Shared component design and styling

`@Components` owns reusable shadcn-compatible/Radix components, Tailwind utility composition, CVA variants, icons, focus rings, form accessibility and shared UI hooks such as `useBoolean`. Domains pass semantic props and feature-level test IDs. Do not duplicate input/loader/dialog/visibility/validation rendering in multiple pages.

Tailwind 4 is installed through Vite, with the **single authoritative stylesheet** `packages/styles/index.css`, semantic foundations in `packages/styles/tokens/`, light and dark theme values in `packages/styles/themes/`, and `components.json` pointing to that same CSS file. RDR brand uses navy `#1F2A44`, neutral/light and charcoal/dark palettes. **Tailwind preflight is intentionally disabled during migration** to preserve legacy public-site and auth layouts; styling utilities and tokens still work. Do not remove old CSS until screenshot and browser regression checks pass. Always make labels, controls and submit buttons reachable on mobile; avoid horizontal page scrolling; short viewports may need a vertical-scroll fallback instead of clipped buttons.

Use system theme on mobile according to existing product UX; do not overwrite the user's desktop stored appearance preference. Do not replace the approved RDR monogram or historical school photography with artificial graphics. Respect source consent and disclosure validity.

Use accessible HTML semantics: a native anchor for navigation, a native button for actions, proper labels and required state, visible focus, screen reader announcements, and no click-only non-focusable `div` interactions.

## 16. Errors, performance, state, permissions and security

Distinguish server state (TanStack Query/cache) from local UI state (component hooks). Do not mirror the same data or boolean into multiple stores. Each state domain has one owner. Shared application state is deliberate; avoid Redux for every local field. Keep permission decisions server-authoritative and have UI capabilities reflect trusted grants rather than role-name heuristics.

No synchronous heavy transformations in JSX. Use pure tested selectors/mappers and memoize only expensive work or identity-sensitive derived props. Keep lists keyed by stable identity, prevent unnecessary rerenders and avoid dynamic circular imports. Avoid full Lodash barrels. No secrets in logs; do not log raw form input/authorization headers. Normalize API errors at business boundaries and public messages at use-case boundaries.

Date/time logic must account for the intended locale/timezone and server representation; centralize reusable date helpers. Feature flags and provider differences should use typed registries; avoid vendor-name checks in presentation components. Observability belongs in shared infrastructure hooks/adapters, not repeated `console.log` statements.

## 17. ESLint, architecture gates and automated enforcement

Migration must add or update real static checks; documentation alone does not guarantee conformance. Enforce:

1. Exact package aliases and allowed import graph; no deep cross-domain imports or circular dependencies.
2. No `lucide-react` import outside `packages/components`.
3. No direct `lodash`/`lodash-es` import outside `packages/base`; prefer one exact `@Base/utilities/<operation>` import.
4. No ternary operators in application code; restrict `switch` and multi-variant `if/else if` to explicit justified exceptions.
5. No `index.ts` below `helpers/` or `constants/`, no per-helper subdirectories, no giant `*.config.ts` containing mixed concerns.
6. One state owner: review duplicate mutable state in hooks, with direct tests for meaningful transitions.
7. Every shared component contract exposes `dataTestId` and forwards it to the native leaf; every migrated page-owned DOM element has a deterministic selector.
8. No unused code, obsolete exports/imports, disabled tests, unsafe `any`, circular dependency or auth credential persistence.

Some of these checks do **not yet exist** in the repo's current architecture scripts. Introduce them as migration code work; do not imply the existing CI already enforces all rules. Rule changes must be committed with the implementation they guard, and existing violations must be explicitly enumerated/managed rather than silently exempted.

Run from `frontend/`:

~~~bash
pnpm install --frozen-lockfile
pnpm check:architecture
pnpm typecheck
pnpm exec jest --runInBand
pnpm lint
pnpm build
~~~

CI verifies architecture, typecheck, Jest, build and static prerender routes. Add targeted responsive/visual tests separately; a passing compiler does not prove layout correctness.

## 18. Migration order and definition of done

The status/inventory and phased plan are integrated in **section 22 of this document**. There is no separate migration document. Do not relocate the same helper twice: ownership review happens before moving files.

**Required migration order**:

1. Add base leaf Lodash utilities and shared browser event constants; centralize icon exports and implement/test `useBoolean`; add architecture rules.
2. Refactor `domain/auth/pages/login` as the reference unit: split constants, flatten helpers, move validation and error classification to owning nested hooks, centralize state in `useLoginForm`, remove redundant hook wrappers, replace render-time ternaries with pure view models, add selectors and tests.
3. Refactor shared auth viewport fit and other cross-page hooks, then recovery into focused nested hooks, then forgot-password; preserve all auth behavior and UI.
4. Migrate dashboard, school website, shell and remaining domains to the same units/selector/icon rules; remove dead files and obsolete import facades.
5. Adopt shared shadcn/Tailwind components without breaking visual layout; extend end-to-end browser/automation coverage; remove superseded CSS only when safe.

**A migrated unit is done only when** its owner and public contract are clear, helpers/constants are flat, state is not duplicated, pure transformations have tests, components expose deterministic `dataTestId` on real DOM, icons and Lodash use package boundaries, no ternary dispatch remains, unused files are removed, docs are accurate, and all applicable checks pass. Preserve public URLs and security invariants.

## 19. Engineering review questions

- Is this code in the closest correct package and parent hook/component?
- Does a helper add value beyond calling an existing `@Base` utility?
- Are constants/enums reused instead of string literals?
- Is one hook the authoritative owner for this form's state?
- Can a new behavior be registered via a typed `Map` instead of another branch?
- Does a visual component delegate to existing `@Components` primitives/icons?
- Does each owned DOM element have a stable `data-testid`, unique among repeated items and free of private data?
- Are source owners obvious from filenames, and are tests located in the owner's `__tests__/`?
- Were old implementations, unused exports and obsolete barrels actually removed?
- Are accessibility, mobile layouts, errors, auth contracts and build/prerender unaffected?

**Guiding principle:** make ordinary code explicit, reusable logic truly reusable, and every screen predictable to develop, test and automate.

## 20. Mandatory shared styles package and reusable component styling

`frontend/packages/styles` is the single source for **generic design decisions**, never feature-specific markup or API logic:

```text
packages/styles/
├── index.css                 # imported once by app/main.tsx; Tailwind v4 theme entry
├── index.ts                  # optional pure TypeScript tokens such as breakpoints
├── tokens/
│   ├── foundations.css      # fonts, type scale, spacing, radii, control sizes, shadows, motion
│   └── breakpoints.ts       # shared viewport queries for browser-aware consumers
├── themes/
│   ├── light.css            # semantic light palette
│   └── dark.css             # semantic dark palette
├── base.css                 # HTML/browser defaults only
├── utilities.css            # small generic classes only
└── RUNBOOK.md
```

**Ownership is non-negotiable.** The `--gh-*` CSS custom properties are canonical; `--rdr-*`, `--canvas`, `--surface`, `--accent` and `--ui-*` compatibility aliases are declared only in `@Styles` until legacy styles are migrated. **No** feature, component or app stylesheet should redefine these root/theme tokens. Light/dark theme changes must be made in the two corresponding theme files, keeping their semantic keys in sync. App startup and the existing theme hook choose `data-theme`; the styles package merely interprets it.

**Styling choice:** use **Tailwind CSS v4 + semantic CSS variables + CVA** for reusable component variants. **New component layout styles target scoped `.module.scss`** using Dart Sass as the compiler. Sass must not become a second source for theme tokens: all color/typography/spacing values remain in `@Styles`. The migration must install Sass and verify Vite, TypeScript, Jest and browser behavior before replacing existing CSS. Shared `@Styles` should not depend on React, components, domains, business, shell or the DOM API. Components may use pure `@Styles` token values/queries; business and base may not.

**Three-layer rule:**

1. **`@Styles`** defines what `primary`, `muted`, `danger`, typography sizes, spacing, height, border, shadow and light/dark appearance *mean*. It does not implement buttons or fields.
2. **`@Components`** defines *how* `Button`, `InputControl`, `Input`, `Text`, `Typography`, `Form`, `FormField`, `Card`, `Container` and `Dialog` render, with accessible semantics, typed props, native `dataTestId`, CVA variants and component-specific CSS. All colors, font sizes, focus rings, dimensions and spacing come from the style tokens.
3. **`@Domain/<name>`** composes these components and supplies content, `variant`, `controlSize`, `padding`, `tone`, event handlers and stable test IDs. Domain CSS is allowed for truly feature-specific layout/artwork, not a second button/input/typography system.

**Example — no per-screen CSS decisions:**

```tsx
import { Form, Input, Button, Card, Typography, Text } from '@Components';

<Card variant="outline" padding="lg" dataTestId="admissions-card">
  <Typography as="h2" variant="heading" size="2xl" dataTestId="admissions-heading">Enquiry</Typography>
  <Text tone="muted" size="sm" dataTestId="admissions-hint">We will contact you.</Text>
  <Form layout="stack" spacing="md" dataTestId="admissions-form" onSubmit={onSubmit}>
    <Input label="Email" controlSize="md" dataTestId="admissions-email" />
    <Button variant="primary" size="lg" type="submit" dataTestId="admissions-submit">Send</Button>
  </Form>
</Card>
```

Shared component APIs should expose sensible defaults. Prefer semantic `variant`, `tone`, `size`/`controlSize`, `spacing`, `layout`, `weight` or `padding` props to arbitrary style overrides. Consumers should not have to choose font pixels, border values or hover colors. **Do not create dozens of boolean variant props**; use a concise typed union backed by a CVA recipe or small CSS class variants. Do not put domain/business logic in the component just to avoid a prop.

**Safe legacy migration:** Existing `app/styles.css` still contains old feature and shell CSS; migrate it gradually to its owning domain/component. The shared button/field/card classes are now owned by `packages/components/styles/primitives.css`, and global theme declarations are owned only by `packages/styles`. Tailwind Preflight stays disabled until responsive and visual regression checks prove it can be safely introduced.

**Required verification:** `pnpm check:architecture` must enforce theme parity, no root-token redefinition, correct shadcn entrypoint and legal package imports; `pnpm typecheck`, Jest and build must pass. Visual verification across desktop, mobile, short viewports, light/dark and high zoom is still mandatory before declaring a CSS migration complete.


## 21. Assets, artwork, page components and persistent shell boundaries

### 21.1 `packages/assets` owns the visual source, not screen behavior

All **new** reusable SVG files, SVG-only TSX vectors, illustration source assets and reusable image assets must live under `frontend/packages/assets/<domain>/`, such as `packages/assets/auth/`. Prefer `.svg` for static standalone art; use a typed SVG-only `.tsx` component only when `currentColor`, theme CSS variables, SVG property forwarding or inline styling must be inherited. Do not copy SVG markup into `pages/`, shell layouts, JSX helper files, or domain feature pages. Do not create duplicate `.svg` and `.tsx` versions of identical art merely to satisfy naming.

A shared renderer `@Components/components/SVGRenderer` receives either `src` for SVG files or `component` for inline SVG-only assets. It owns native `dataTestId` forwarding, the `img`/SVG leaf and accessible roles. Decorative art **must** be `decorative`, hidden from assistive tech, and has no synthetic alt text; meaningful art **must** specify `alt`. Keep status copy, surrounding layout and behavior in `@Domain/auth/components`, not in `@Assets`. The `SVGRenderer` component tests must exercise both rendering modes.

**Approved auth example:** `@Assets/auth/AcademicPathwaysSvg` defines the vector, `@Components/components/SVGRenderer` renders the native SVG, and `@Domain/auth/components/authAmbientGraphic` owns the caption and page-level markup. A composite `AuthStatusArt` used across recovery and forgot-password belongs in the shared *auth-domain component layer*, not `pages/` and not `@Assets` because it contains status-specific presentation.

**No dead prototypes:** the retired `AuthFrame`, `AuthMobileArtwork`, `loginv2/SoffitGradient` and standalone old `LoginV2` must not be restored or imported. Keep the single current RDR auth experience, its approved monogram, styles and accessible form. Delete all unused SVGs, import facades, component wrappers, styles and tests only after proving no live importer exists and updating the owning tests.

### 21.2 Strict ownership by directory

`pages/` is for actual routed pages and page-local `components`, `hooks`, `containers`, `templates`, `helpers`, `constants` and `__tests__`. Any UI component shared by multiple pages in the **same domain** belongs in `packages/domain/<domain>/components/`. A domain-wide viewport hook belongs in `packages/domain/<domain>/hooks/`; a session security guard and proof lifecycle belong in `packages/domain/auth/core/guards` and `core/session`.

Do **not** leave `AuthAmbientGraphic.tsx`, `LoginLayout.tsx`, `RecoveryFlowBoundary.tsx`, random SVG elements, or compatibility stubs at `auth/pages/` root. Only true routing page files go under their own `pages/<page>/`. Existing tests at old flat locations must be relocated to the owning unit's `__tests__` during migration.

`@Assets` has no business, hooks, styles, or renderer logic; `@Components` can import `@Assets`; `@Domain` can import `@Components` and `@Assets`. `@Styles`, `@Base` and `@Business` must not depend on `@Assets`. Vite, TypeScript and Jest must resolve every new alias consistently. Never expose sensitive or user-created content as public static files.

### 21.3 Public routes have **no workspace layout**

The school website is a public site. Its routes render `@Domain/publicSite` directly in `AppRoutes`: **no shared `PublicLayout` or `Shell` wrapper around school pages, and no session bootstrap requirement for public pages**. `LandingPage.tsx` is presentational; `containers/landingContainer` and `hooks/useLandingPage` own route navigation, enquiry sheet state, scroll restoration and SEO updates. The site header/footer may be composed inside that public page, without creating a shell-level layout.

**Content and configuration:** school identity, school photographs and alt text, disclosures, navigation labels, enquiry options, panel mappings, SEO metadata and any fixed copy live in **named files under `constants/`**; never in generic `config/content.ts` or one enormous `*.config.ts`. `registry/schoolPages.registry.tsx` owns the React component-to-public-route table. `PUBLIC_SITE_PATHS` is derived from that registry and exported through `@Domain/publicSite`; the shell must not maintain a second copy of the route list. `config/` is reserved for genuine runtime composition, not static content or dictionaries. Keep all existing public SEO paths, disclosure caveats and enquiry behavior.

### 21.4 Protected routes share **one persistent** layout

Protected routes must follow this exact structural ownership:

```text
AppRoutes
└── RequireSession            # session proof and redirect gate
    └── ProtectedLayout       # mounts Shell once for authenticated workspace
        └── Shell              # sidebar/topbar/app catalog/mobile navigation
            └── Outlet         # ONLY this area changes with child route
                ├── /dashboard
                └── future real domain routes
```

The `ProtectedLayout` route is **pathless**, and no domain screen creates its own `Shell`. Within protected navigation the shell's local state, sidebar, pinned apps, theme and topbar survive without remounting. Do **not** put `key={pathname}` on `ProtectedLayout` or `Shell`, render distinct shells per child route, wrap protected child routes in another `BrowserRouter`, or rebuild the provider tree on route changes. React may rerender components when their state/context changes; the invariant is **no unmount/remount** of the shared shell when switching protected children. Add tests proving the same DOM shell remains mounted.

`@Shell` composes only real route implementations exported by domain root APIs. Remove dead route entries and unimplemented People, Accounts, Integrations, Settings or Security pages; **do not** show clickable phantom navigation or catalog entries for them. The catalog may add a module only after that module has a routed, authorized implementation. Account/theme preferences inside the shell are controls, not automatically a separate Settings route. Retain the Dashboard domain at `/dashboard`; it is the implemented child page, **not** the persistent layout itself.

### 21.5 Safe pruning and migration tests

Before deleting source, verify actual live imports, dynamic imports, route registry, CSS consumers, public API exports and tests. Remove dead source + imports + unused tests + CSS as one change; never silently redirect an old production authentication URL to a non-equivalent screen. Keep `/loginv2` redirecting to the canonical sign-in URL for bookmarks, `/setup` for invitation links, all thirteen public SEO paths and `/dashboard`. The unused standalone `LoginV2` prototype is distinct from the preserved `/loginv2` compatibility *route*.

The backend owns recovery retry limits; routing and artwork changes must not turn invalid links into sign-in states or leak one-time proof tokens to persistent storage. Test accessible artwork, original CSS hooks, mobile theme, login navigation, enquiry modal, public SEO routes, and protected layout identity; CI does not replace manual visual testing.

## 22. Unified migration inventory and enforcement status

This section replaces the retired `FRONTEND_REFACTOR_MIGRATION.md` and remains inside the **single authoritative** frontend engineering document. Status is based on implemented source changes, **not** a blanket claim that Revision 6 is fully enforced.

| Area | Current state | Remaining strict requirement |
|---|---|---|
| `packages/assets/auth` and `@Components/SVGRenderer` | Active academic-pathways SVG moved; native renderer and its tests present; unused auth graphics/prototypes deleted | Move any remaining screen-owned new vectors to assets as touched |
| `domain/auth/pages` | Shared login layout moved to `domain/auth/components/loginLayout`; viewport hook moved to domain hooks; recovery guard/session moved into core; redundant page wrappers deleted | Enforce owner-local `__tests__`, flatten remaining page helper directories |
| `domain/auth/pages/login` | Functional form and dedicated hooks exist | Split remaining `config/` into flat `constants/`, flatten helpers, complete `@Base` Lodash migration and view-model/strategy selection |
| `domain/auth/pages/recovery` | Server-authoritative activation recovery preserved | Split its large useRecovery hook into link/grant/form/submission/retry units; flatten helper units |
| `domain/auth/pages/forgotPassword` | Functional form and countdown preserved | Flatten remaining helpers/constants, resolve branch-only view logic, add owner-local tests |
| `domain/publicSite` | School content/SEO/panel data in constants, routed views in domain registry, page root presentational, navigation/enquiry owned by hook | Finish finer-grained section component splitting and remove duplicate CSS decisions only after visual verification |
| `domain/dashboard` | Data preview split into constants, typed date helper and tests flattened; implemented dashboard route kept | Integrate real business APIs only when they exist, avoiding fabricated counts |
| `shell` | Pathless `ProtectedLayout` persists the shell, tested across child route changes; only implemented applications in catalog; no public layout | Add routed modules only with working domain/page contracts; no speculative links |
| `packages/styles` | Central light/dark tokens, typography, spacing, breakpoints and Tailwind mappings | Progressively move legacy `app/styles.css` screen rules to owning components/domains |
| `packages/components` | Shared Form, Text, SVGRenderer, Input, Typography, Button and other primitives with typed selectors | Complete shared icon exports, `useBoolean`, native leaf selector coverage and consistent variants |
| `packages/base` | Several typed leaf utilities already exist | Add/test requested Lodash `_trim`, `_keys`, `_size`, `_find`, `_isUndefined`, `_isFunction` as needed, never import Lodash outside Base |
| `scripts/check-*.mjs` | Existing boundaries, UI contracts, styles parity, test/build checks | Add full no-ternary, strict flat helper/constant directory and icon import enforcement as code is migrated |

**Mandatory batch sequence:** shared Base/icons/useBoolean improvements → login as the reference unit → recovery + forgot-password + auth viewport rules → remaining domain sections/components → cross-repo static enforcement and browser visual tests. After each batch run frozen dependency install, architecture, TypeScript, Jest, lint and production build/prerender. Do not claim `pnpm lint` is already part of CI unless that step is actually present in its workflow. Route behavior, accessible selectors, real browser responsive checks, backend security states and intentional dead-file removal are all part of the acceptance criteria.

**Definition of done:** source paths meet the rules, no obsolete exports/barrels remain, components get stable leaf test IDs, helpers and static copy have a single owner, only real routes appear in the app launcher, the shell survives protected navigation, tests are colocated and pass, and visual checks verify light/dark/mobile. This is a live migration ledger: completed rows stay factual and uncompleted enforcement remains visible.

## 23. Revision 7 strict classification, folders and ownership

**Before touching any file, read it line by line with its direct and indirect importers, tests, routes, CSS, accessible DOM and live behavior.** A refactor must reassign responsibility, not merely move the same 200-line component into a subfolder. Unused-code deletion requires proof there is no importer, dynamic route, indirect export, generated dependency, stylesheet consumer or required test.

| Code responsibility | Required owner | Explicitly forbidden |
|---|---|---|
| Design tokens, system/light/dark palette, fonts, spacing, radii, motion, focus and breakpoints | `@Styles` | Per-feature hex colors and token duplication |
| SVG asset paths / raw vector markup | `@Assets/<feature>` | Auth pages/components embedding vector paths |
| Accessible HTML, Text, Typography, Heading, Button, Input, Form, Label, Row, Box, Icon, Link | `@Components` | Native HTML controls and generic `div/span/p/h3` inside migrated feature presentation |
| Reused auth visual layout/status | `domain/auth/components/<unit>` | Loose `auth/pages/*.tsx` components |
| Page-only visual | `pages/<page>/components/<unit>` | Route registry or business layer |
| JSX-producing selector, grouped label, status badge, variant content, repeated option | `renderer/<owner>.<useCase>.tsx` | Inline JSX ternary or nested callbacks |
| Non-JSX pure computation: class selection, filters, validation, copy fallback | `helpers/<owner>.<useCase>.ts` | A `renderer` for non-JSX values |
| Enums, static option arrays, metadata, copy, routes, test IDs, timeouts | `constants/<owner>.<concern>.ts` | Inline magic strings, huge mixed config |
| State, stable actions, subscriptions and side effects | **One** owning `hooks/use<Feature>/` | Peer hooks duplicating state or passing many setters around |
| Auth session context | `auth/core/providers/authProvider/` | Flat multipurpose provider |
| Recovery proof/grant lifecycle | `auth/core/session` and owned helpers/constants | Generic page helper, persistent bearer token |
| App routing/persistent workspace | `@Shell` plus shared route constants | Separate Shell per page; duplicate route strings |
| Feature-scoped style | `styles/<unit>.module.scss` using `@Styles` vars | Global BEM strings, domain-owned theme colors |
| Tests | Owner's `__tests__/` | Scattered page-root test files |

**Flat file rules:** `helpers/`, `constants/`, and **`renderer/`** have no `index.ts`, no folder per single function and no generic `utils.ts`. Each may contain named flat files and optional sibling `types/` and `__tests__/`. Actual hook/component/provider units may have an intentional public `index.ts`, and can nest their own helpers, constants, renderers, hooks, styles and tests when their ownership is clear.

Examples: `helpers/recovery.getMobileDescription.ts`, `renderer/credentialField.label.tsx`, `renderer/loginLayout.appearanceOption.tsx`, `constants/loginLayout.variant.ts`, `constants/forgotPassword.testIds.ts`. Never conflate **config** with constants: a static map is a constant; an executable selector is a helper; JSX selection is a renderer; a runtime provider factory is a factory helper. No extra directory if only an unconditional expression exists.

**Zero ternaries** is a strict syntax rule for application TypeScript/TSX (all `ConditionalExpression` nodes). No nested JSX ternaries, inline ternary class/label/icon, ternary inside helper, or substitute nested `&&` / immediately invoked functions to hide complex branching. Binary choice uses a named `if/else` helper or renderer. Three or more executable strategies use a typed `Map<Enum, Function>`; static mappings use `Record<Enum, Descriptor>`. Never claim this is enforced today until the new AST lint check is actually installed.

## 24. Renderers: JSX-only pure functions with tests

The `renderer/` directory contains **pure typed functions returning JSX**, distinct from full reusable components. A renderer accepts values/callbacks and returns `ReactNode` / `ReactElement`. It may call another renderer or shared component. It **must not** call hooks, APIs, storage, router navigation, subscribe to events or mutate state. It does not own general layouts shared across pages. If a renderer is large/reusable and requires lifecycle, it should become an actual component with its own folder.

~~~text
auth/components/credentialField/
├── CredentialField.tsx
├── constants/credentialField.testIds.ts
├── renderer/
│   ├── credentialField.label.tsx
│   ├── types/credentialField.label.ts
│   └── __tests__/credentialField.label.test.tsx
├── styles/credentialField.module.scss
├── types/credentialField.props.ts
├── __tests__/CredentialField.test.tsx
└── index.ts
~~~

The grouped label currently uses a forbidden ternary. The correct target pattern is:

~~~tsx
// renderer/credentialField.label.tsx
import { FieldLabel, Row } from '@Components';
import styles from '../styles/credentialField.module.scss';
import type { RenderCredentialFieldLabel } from './types/credentialField.label';

export const renderCredentialFieldLabel: RenderCredentialFieldLabel = (props) => {
  const label = (
    <FieldLabel htmlFor={props.id} dataTestId={props.ids.LABEL}>
      {props.label}
    </FieldLabel>
  );

  if (props.grouped) {
    return (
      <Row className={styles.passwordLabel} dataTestId={props.ids.LABEL_ROW}>
        {label}
        {props.action}
      </Row>
    );
  }
  return label;
};
~~~

`FieldLabel` and `Row` are **target shared primitives, not currently guaranteed to exist**. First implement them under `@Components` with correctly forwarded native test IDs and label semantics, then adopt them. Test both grouped/non-grouped render and `htmlFor` relationship. No one-off raw `<label>`/`<div>` in an auth renderer. The caller simply renders `{renderCredentialFieldLabel(model)}`. Do not place the helper inside JSX or call a React hook from it.

Use `renderer/recovery.content.tsx` for choosing form vs state, `renderer/forgotPassword.notice.tsx` for status UI, `renderer/authStatusArt.badge.tsx` for status decoration, `renderer/appLauncher.app.tsx` for repeated catalog row JSX, and `renderer/appRoutes.route.tsx` for mapping route descriptors. Use a helper instead when the output is **not** JSX. Tests belong to the renderer in `renderer/__tests__/` or its nearest owning unit.

## 25. Typed status enums, exhaustive maps, and zero inline icon/color literals

No raw string status comparisons (`kind === 'email'`, `props.view === 'ready'`, `variant='login'`) in migrated feature logic or feature props. Define explicit typed enums or stable `as const` runtime maps in nearest `constants/` file. Prefer string enums for server-aligned or cross-unit state identifiers. Every known state must have an explicit visual description; the compiler must detect missing mappings.

~~~ts
// auth/components/authStatusArt/constants/authStatusArt.kind.ts
export enum AuthStatusKind {
  EMAIL = 'email',
  SUCCESS = 'success',
  EXPIRED = 'expired',
  MISSING = 'missing',
  CHECKING = 'checking',
  UNAVAILABLE = 'unavailable',
  ATTEMPTS_EXHAUSTED = 'attempts-exhausted',
}

export enum AuthBadgeKind {
  CONFIRMED = 'confirmed',
  INVALID = 'invalid',
  NONE = 'none',
}

// constants/authStatusArt.badge.ts
export const AUTH_STATUS_BADGE: Readonly<Record<AuthStatusKind, AuthBadgeKind>> = {
  [AuthStatusKind.EMAIL]: AuthBadgeKind.CONFIRMED,
  [AuthStatusKind.SUCCESS]: AuthBadgeKind.CONFIRMED,
  [AuthStatusKind.EXPIRED]: AuthBadgeKind.INVALID,
  [AuthStatusKind.MISSING]: AuthBadgeKind.INVALID,
  [AuthStatusKind.CHECKING]: AuthBadgeKind.NONE,
  [AuthStatusKind.UNAVAILABLE]: AuthBadgeKind.NONE,
  [AuthStatusKind.ATTEMPTS_EXHAUSTED]: AuthBadgeKind.NONE,
};
~~~

~~~tsx
// renderer/authStatusArt.badge.tsx
import { Icon, Text } from '@Components';
import { AUTH_STATUS_BADGE, AuthBadgeKind } from '../constants/authStatusArt.badge';
import styles from '../styles/authStatusArt.module.scss';
import type { RenderAuthStatusBadge } from './types/authStatusArt.badge';

export const renderAuthStatusBadge: RenderAuthStatusBadge = (kind) => {
  const badge = AUTH_STATUS_BADGE[kind];
  if (badge === AuthBadgeKind.CONFIRMED) {
    return <Icon name="check-circle-2" decorative className={styles.confirmedBadge}/>;
  }
  if (badge === AuthBadgeKind.INVALID) {
    return <Text as="span" className={styles.invalidBadge} aria-hidden="true">×</Text>;
  }
  return null;
};
~~~

**Before this example compiles**, add the required shared `Icon` abstraction, typed icon identifiers and import the enum from the correct local file. Renderer functions are pure. Every dynamic icon, color and case is tested. Do not import `lucide-react` in domain/shell. Do not store a React component or JSX expression under a generic helper path when it is a renderer.

**CSS colors:** `fill="#039C90"` and `color="#fff"` are forbidden. Use `styles.confirmedBadge` with `color: var(--gh-color-on-success)` and `fill: var(--gh-color-success)`. Add any new semantic token (e.g. `--gh-color-on-success`) in **both** `@Styles/themes/light.css` and `dark.css` before use. Background, border, title, ring, illustration and disabled colors all follow the same rule; no new theme palette in domain SCSS.

## 26. LoginLayout exact decomposition and theme-option renderer

~~~text
auth/components/loginLayout/
├── LoginLayout.tsx
├── constants/
│   ├── loginLayout.variant.ts
│   ├── loginLayout.copy.ts
│   ├── loginLayout.appearance.ts
│   └── loginLayout.testIds.ts
├── helpers/
│   ├── loginLayout.className.ts
│   ├── loginLayout.mobileDescription.ts
│   ├── types/
│   └── __tests__/
├── renderer/
│   ├── loginLayout.appearanceOption.tsx
│   ├── loginLayout.visual.tsx
│   ├── types/
│   └── __tests__/
├── styles/loginLayout.module.scss
├── types/loginLayout.props.ts
├── __tests__/LoginLayout.test.tsx
└── index.ts
~~~

~~~ts
// constants/loginLayout.variant.ts
export enum AuthLayoutVariant {
  LOGIN = 'login',
  FORGOT = 'forgot',
  RECOVERY = 'recovery',
}

// constants/loginLayout.copy.ts
export const LOGIN_LAYOUT_COPY = {
  DEFAULT_EYEBROW: 'YOUR SCHOOL WORKSPACE',
  RECOVERY_EYEBROW: 'SECURE YOUR ACCOUNT',
  SCHOOL_NAME: 'RDR School',
} as const;

// constants/loginLayout.appearance.ts (shared theme values/icons are imported)
import { Monitor, Sun, Moon } from '@Components/icons';
import { THEME_PREFERENCE } from '@Components/hooks/theme/constants/themePreference';

export const APPEARANCE_OPTIONS = [
  { value: THEME_PREFERENCE.SYSTEM, label: 'System', Icon: Monitor },
  { value: THEME_PREFERENCE.LIGHT, label: 'Light', Icon: Sun },
  { value: THEME_PREFERENCE.DARK, label: 'Dark', Icon: Moon },
] as const;
~~~

`@Components/icons` and `THEME_PREFERENCE` are **planned APIs**. If current theme ownership only has a type union, add one runtime constant **in the shared theme unit**, not a second status map here. Default destructuring uses named constants, not `variant = 'login', eyebrow = '...'`. A separate `getMobileDescription(mobile, desktop)` helper uses explicit `if` and preserves an intentionally empty mobile string:

~~~ts
export const getMobileDescription = (
  mobileDescription: string | undefined,
  description: string,
): string => {
  if (mobileDescription !== undefined) return mobileDescription;
  return description;
};
~~~

Appearance rendering: the options live in constants; `renderer/loginLayout.appearanceOption.tsx` receives typed option, selected preference and a **stable bound action**. It uses `@Components/Button`, an icon from `@Components/icons`, `dataTestId` from `constants/loginLayout.testIds.ts`, `cn` from `@Base/utilities/cn`, and CSS-module classes `styles.appearanceOption`/`styles.appearanceSelected`. It derives `selected` as a boolean and calls an already-named action, not an inline JSX handler. The layout maps with `_map(APPEARANCE_OPTIONS, renderAppearanceOption)` from `@Base/utilities/map` (implement this leaf utility first). No `.map`, raw `<button>`, direct `'is-active'` class, ternary, or inline icon selection.

`helpers/loginLayout.className.ts` constructs the layout classes using `cn`, a typed variant-to-module-class map, and an `if` for `minimalMobile`. Avoid `${minimalMobile ? '...' : ''}` in JSX. `LoginLayout.tsx` composes named `Box`/`Section`/`Heading`/`Text`/`AppLink`/`Icon` primitives instead of manually creating arbitrary `div/h1/h2/p/span/button`. Keep one actual form in the DOM across desktop/mobile. Test theme appearance buttons, dark/light/system behavior, media queries, compact states, and short-viewport reachability.



## 27. No native generic markup in feature presentation: shared primitive contract

In migrated domain and shell presentation/renderer components, **do not author generic native HTML directly**: `<div>`, `<span>`, `<p>`, `<h1>`–`<h6>`, `<small>`, `<strong>`, `<label>`, `<button>`, `<input>`, `<form>`, `<a>`. The design system handles appearance, semantics, native event behavior, accessibility, loading, disabled states and responsive spacing. A feature provides content, typed variant, state and named actions—not its own low-level control implementation.

| Semantic requirement | Shared API to use/build | Required behavior |
|---|---|---|
| Layout/container | `Box` / `Row` / `Stack` / `Container` | Typed `as` where needed, layout/gap/size, forwarded `dataTestId` |
| Section / main / header / footer / navigation | `Section`, `Main`, `Header`, `Footer`, `Navigation`, or a typed semantic `Box as="..."` | Preserve landmark semantics and accessible name |
| Heading | `Heading` or `Typography as="h2"` | Correct heading hierarchy, size and token-based color |
| Body, inline, bold or small text | `Text` / `Typography` with correct `as` | Own typography and native text semantics |
| Input, label, password | `Input` / `InputControl` / `FieldLabel` / `FormField` / `PasswordField` | Unique id, htmlFor, accessible name, errors, autofill |
| Actions | `Button` / `IconButton` | Native `type`, disabled, loading, focus, tests |
| Form | `Form` | Native submit, aria, spacing and test ID |
| In-app link | `AppLink` wrapping React Router Link | Typed route `to`, no navigation on a button |
| Alerts/status/loading | `Alert` / `StatusMessage` / `Spinner` | `role`/`aria-live` and reserved geometry |
| Icons/illustrations | `Icon` / `SVGRenderer` | Decorative/informative accessibility and theme-driven style |
| Modal, panel, table, card | Shared `Dialog` / `Card` / `DataTable` | Radix/native accessibility and named variants |

Some target primitives (`FieldLabel`, `Row`, `AppLink`, `Icon`, semantic landmarks, etc.) **are not presently in the repository**. First implement/export/test the missing shared primitive before migrating feature JSX. Never write an invalid import just to satisfy these guidelines. If a primitive cannot preserve the required HTML semantics, fix that primitive instead of substituting `div onClick` or a misleading `role`. **Exceptions:** native HTML is required *inside* shared primitive implementations, React Router/third-party adapters, and app bootstrap. A route renderer may create `<Route>` / `<Navigate>`; SVG asset vector files naturally own SVG path nodes. This rule targets **feature markup duplication**, not correct HTML implementation inside component boundaries.

**Correct migration for `ForgotPasswordForm`:** its form, input, notice, alert, back link, resend action, error slot and success state are assembled from shared UI. Its status icon/body and class selection belong to `renderer/forgotPassword.notice.tsx`, `renderer/forgotPassword.sentState.tsx`, `helpers/forgotPassword.getNoticeViewModel.ts`. Its labels, messages, email input metadata and IDs belong in named constants. Use `Button loading` for standard spinner instead of emitting `LoaderCircle` in every page; use semantic `AppLink to={APP_ROUTE.LOGIN}` for navigation and a proper submit button for resend. Preserve cooldown, button reachability, the privacy-preserving "if eligible" response and live-region announcements.

**Correct migration for `Recovery`:** `Recovery.tsx` composes `LoginLayout`, a scoped `styles.recovery` wrapper and the result of `renderer/recovery.content.tsx`, which selects **exactly one** of `RecoveryPasswordForm` and `RecoveryState`. A typed `recovery.getViewModel.ts` helper resolves form visibility, mobile title, description and status from enums; do not do `props.invitation ? ...`, `formVisible ? ...` or inline `view === 'ready' || view === 'saving'` in the page. Pass `AuthLayoutVariant.RECOVERY`, not raw `variant="recovery"`, and use named copy constants, not `eyebrow="SECURE YOUR ACCOUNT"`.

### View renderer acceptance tests

For every conditional renderer assert: initial state; each variant; no extra/hidden duplicate form; correct icon/accessible labels; no invalid native controls; correct test IDs; keyboard and Enter/Space support; reserved error geometry; light/dark token classes. Pure renderer tests may use React Testing Library without mocking the business layer because the renderer must have no business imports.

## 28. CSS ownership: migrate the supplied 1,095-line authentication stylesheet to SCSS modules

**Chosen stack:** Tailwind CSS v4 + CVA for shared component variants and **Dart Sass `*.module.scss`** for substantial component-specific styling. Sass preprocesses and scopes; **`@Styles` remains the only source of design decisions** (light/dark tokens, brand palette, font families/scales, spacing, motion, focus, borders and shadow). The repo currently has no `sass` compiler dependency, so **install and test it during the code refactor**—this document update itself does not make Sass work.

Target ownership:

~~~text
packages/styles/
├── index.css                         # one Tailwind entry
├── tokens/foundations.css
├── tokens/breakpoints.ts
├── themes/light.css
├── themes/dark.css
└── scss/                             # shared mixins only; no duplicated palette
    ├── _media.scss
    ├── _a11y.scss
    └── _layout.scss

domain/auth/components/loginLayout/styles/loginLayout.module.scss
domain/auth/components/credentialField/styles/credentialField.module.scss
domain/auth/components/authStatusArt/styles/authStatusArt.module.scss
domain/auth/pages/forgotPassword/components/forgotPasswordForm/styles/forgotPasswordForm.module.scss
domain/auth/pages/recovery/styles/recovery.module.scss
~~~

The present `auth/components/loginLayout/styles/optionOneAuth.css` is **not one logical component's styling**. The uploaded stylesheet includes layout, branding, SVG art, input shell, password toggle, form errors, forgot/password reset status, recovered-state illustrations, both theme palettes, autofill, responsive layout, very short windows and reduced motion. **Split by owner**, do not rename the same 1,095-line file to `loginLayout.module.scss` and leave all responsibilities mixed.

Required ownership during extraction:

- **Generic shared primitives:** typography, Button variants, input focus/disabled, error alert, control sizing, loader, layout primitives → `@Components` recipes using `@Styles` values.
- **LoginLayout:** desktop/mobile visual geometry, toolbar, brand, content centering, visual copy and measured fit → `loginLayout.module.scss`.
- **CredentialField:** label group, input shell, password action, autofill, reserved error row → `credentialField.module.scss`.
- **AuthStatusArt:** halo/symbol/badge/sparks, accessible decorative placement → `authStatusArt.module.scss`.
- **ForgotPasswordForm:** notice, accepted-request state, actions and responsive form blocks → `forgotPasswordForm.module.scss`.
- **Recovery:** proof state, requirement/attempt status, form presentation and responsive actions → `recovery.module.scss`.
- **Pure design values:** gradients, error/success/brand colors, placeholders, font sizes/weights, spacing, shadows, radii and focus tokens → `@Styles/themes/*` / `tokens/foundations.css`.

CSS Modules are imported with e.g. `import styles from './styles/loginLayout.module.scss'`, then `className={styles.root}`. Combine variants with `cn` imported from `@Base/utilities/cn`, e.g. `cn(styles.root, styles.compact)`. **Any branching CSS class selection is in a named helper**, not JSX string interpolation, nested ternary or raw class string `'is-active'`. Use `className={styles.recovery}` rather than `className="rdr-recovery"`. Domain-specific CSS modules are **owned next to their components**, while only generic design utilities/mixins/tokens belong in `packages/styles`. Do not centralize one-use feature selectors in the global styles package just to make an import shorter.

~~~scss
/* authStatusArt.module.scss; both tokens must exist in @Styles light/dark themes. */
.root {
  color: var(--gh-color-text);
}
.confirmedBadge {
  fill: var(--gh-color-success);
  color: var(--gh-color-on-success);
}
.invalidBadge {
  color: var(--gh-color-primary-foreground);
  background: var(--gh-color-danger);
  border-radius: var(--gh-radius-full);
}
~~~

**Forbidden in migrated feature SCSS:** hex/rgba/hsla color literals; duplicated light/dark color maps; global `.rdr-*` selectors without scoped module ownership; `!important` without documented interoperability reason; global `:root` overrides; magic gradient colors; inline `style={{ fill: ... }}`. New semantic roles require adding corresponding token keys to both `light.css` and `dark.css` and passing theme-parity checks. Numeric geometry can be local **only when it is layout-specific** (for example, artwork aspect ratio or measured safe-area offset), never a second generic type/spacing scale.

**Migration trap:** converting `.rdr-login-v2__input-shell input` to a module class without changing the corresponding JSX will silently break styling. Move rules and their matching DOM classes **atomically**, and verify DOM screenshots before deleting old selectors. For controlled third-party internals, use a scoped `:global(...)` child under a local root, never a blanket global selector. Keep autofill, focus rings, password visibility, short viewport scroll fallback, high zoom, RDR monogram and system-mode mobile visuals working.

**Required visual scenarios:** width 320/375/390/720 px, tablets, desktop, short-height browser window, 200% zoom, dark/light/system, reduced motion, iOS/Chrome autofill, focused/errors/expired/success/loading states, mobile keyboard open. Normal-height mobile screens should fit without unwanted scrolling, but unusually short/zoomed views **must scroll rather than clip actions**. Passing TypeScript/Jest/Vite alone is not visual verification.

## 29. Routes, PUBLIC/PROTECTED variants and the recovery-flow effect hook

All navigation paths are named once. The user-selected canonical location is `frontend/app/core/routes/` with a narrowly scoped **pure-data leaf** allowed as a dependency of domain/shell (no React, business code, API, component import or state). Configure the `@AppRoutes` alias explicitly in TypeScript, Vite, Jest and boundaries **before adoption**. Other `frontend/app` imports remain disallowed from packages. If the project architecture cannot support this leaf without a dependency cycle, move canonical constants into `@Base/constants/routes` and keep `app/core/routes` as integration-only; **never duplicate the dictionaries**.

~~~ts
// app/core/routes/constants/appRoutes.paths.ts
export const APP_ROUTE = {
  HOME: '/',
  LOGIN: '/login',
  LOGIN_LEGACY: '/loginv2',
  FORGOT_PASSWORD: '/forgot-password',
  RESET_PASSWORD: '/reset-password',
  SET_PASSWORD: '/set-password',
  SETUP_LEGACY: '/setup',
  DASHBOARD: '/dashboard',
} as const;

// app/core/routes/constants/appRoutes.access.ts
export enum RouteAccess {
  PUBLIC = 'PUBLIC',
  PROTECTED = 'PROTECTED',
}
export enum PublicRoutePolicy {
  OPEN = 'OPEN',
  GUEST_ONLY = 'GUEST_ONLY',
}
~~~

**IMPORTANT:** PUBLIC does **not** always mean guest-only. Authenticated users must still view the school website (`PublicRoutePolicy.OPEN`). Login and certain auth pages may be `GUEST_ONLY`. Protected pages require an initialized session. Model this policy deliberately, not by checking pathname substrings.

Create `PublicRoute` and `ProtectedRoute` reusable route wrappers/guards under their owning auth/shell route unit with typed `RouteAccess` and optional `PublicRoutePolicy`. Only guards trigger authenticated session initialization; normal public marketing and token exchange must not eagerly call `/me`. Route redirect and loading UI must use shared components and named route constants. Preserve `/loginv2` redirect, `/setup` invitation links and all public SEO paths. A protected layout remains a **single pathless mounted parent**:

~~~text
BrowserRouter
└── AuthProvider
    └── AppRoutes
        ├── public school pages (no workspace layout)
        ├── PublicRoute (guest-only login / other allowed public auth pages)
        └── ProtectedRoute
            └── ProtectedLayout -> Shell -> Outlet
                ├── /dashboard
                └── future implemented protected domain pages
~~~

**Do not** put a separate `Shell` inside every `Route` element: that unmounts sidebar/theme/pins on navigation. Route JSX registration is allowed to use React Router primitives; grouping and rendering route descriptors belongs in `renderer/appRoutes.route.tsx`, with `_map` from `@Base/utilities/map`. The route metadata itself is under `registry/` or constants; a renderer must not decide auth by looking for `'/protected'` in a string. Keep `PUBLIC_SITE_PATHS` exported from `@Domain/publicSite`; avoid a second hand-maintained thirteen-path copy.

~~~tsx
// Example: typed route descriptor renderer only, not a second shell instance
export const renderRoute = (route: AppRouteDescriptor) => (
  <Route key={route.path} path={route.path} element={route.element}/>
);

// AppRoutes maps descriptors via the base leaf
const publicRouteElements = _map(publicRoutes, renderRoute);
const workspaceRouteElements = _map(workspaceRoutes, renderRoute);

// ONE protected parent wraps the entire workspace group; no per-page Shell
<Route element={<ProtectedRoute/>}>
  <Route element={<ProtectedLayout/>}>{workspaceRouteElements}</Route>
</Route>
~~~

The outer `AppRoutes` becomes a thin container of named route arrays and renderers; no native `PUBLIC_SITE_PATHS.map(...)` or inline mapping arrow that emits a large JSX tree.

`RecoveryFlowBoundary` is a **side-effect hook**, not a visual component. Replace the `RecoveryFlowBoundary.tsx` null component with `useRecoveryFlowBoundary` under `auth/core/hooks/useRecoveryFlowBoundary/`, with its own types, helper for purpose selection and tests. Call it unconditionally from an already-rendered component **inside BrowserRouter and AuthProvider**, such as the app-routes controller. On leaving/switching recovery paths it clears pending proof state; do not clear a valid flow merely because a layout rerenders. Never reintroduce a no-op JSX component whose sole purpose is executing effects. No token gets put in persistent browser storage.


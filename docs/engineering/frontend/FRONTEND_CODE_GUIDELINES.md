# Frontend Code Guidelines

> **Status:** Mandatory, authoritative frontend engineering standard — **Revision 6 (November 2026)**.
>
> **Scope:** Every source file under `frontend/`, including `app`, `packages/base`, `packages/business`, `packages/components`, `packages/styles`, `packages/assets`, `packages/domain`, `packages/shell`, unit tests, integration tests, styles, generated code boundaries, and scripts.
>
> **Stack:** React 19, TypeScript strict, pnpm, Vite, Tailwind CSS 4, shadcn-compatible Radix primitives, Lucide icons through the components package, Jest, and the existing project architecture checks.
>
> **Precedence:** This document replaces previous frontend guideline revisions and examples. Existing code that differs is **migration debt**, not an exception to this standard. Do not claim implementation is complete until checks and behavior regressions pass.

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

The status/inventory and phased plan are integrated in **§22 of this document**. There is no separate migration document. Do not relocate the same helper twice: ownership review happens before moving files.

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

**Styling choice:** use **Tailwind CSS v4 + semantic CSS variables + CVA** for reusable component variants. Use scoped `.module.css` for unusually complex component-only selectors. Do not introduce `.module.scss`/Sass without a concrete demonstrated need: Sass adds a compiler and risks a second source for tokens. Shared `@Styles` should not depend on React, components, domains, business, shell or the DOM API. Components may use pure `@Styles` token values/queries; business and base may not.

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

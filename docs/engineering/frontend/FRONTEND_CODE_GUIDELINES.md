# Frontend Code Guidelines

> **Status:** Mandatory, authoritative frontend engineering standard — **Revision 8 (November 2026)**.
>
> **Scope:** Every source file under `frontend/`, including `app`, `packages/base`, `packages/business`, `packages/components`, `packages/styles`, `packages/assets`, `packages/domain`, `packages/shell`, unit tests, integration tests, styles, generated code boundaries, and scripts.
>
> **Stack:** React 19, TypeScript strict, pnpm, Vite, Tailwind CSS 4, shadcn-compatible Radix primitives, Lucide icons through the components package, Jest, and the existing project architecture checks.
>
 > **Precedence:** Revision 8 rules on route ownership, tests, and flat helper/constant units override older examples. Generic examples illustrate reusable patterns; actual GrooveHub migration paths are recorded separately in the audit inventory. This document replaces previous frontend guideline revisions and examples. Existing code that differs is **migration debt**, not an exception to this standard. Do not claim implementation is complete until checks and behavior regressions pass.

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
packages/domain/catalog/
├── pages/
│   └── item/
│       ├── Item.tsx                     # presentation only
│       ├── index.ts                      # page public entry
│       ├── containers/
│       │   └── itemContainer/
│       │       ├── ItemContainer.tsx
│       │       ├── types/
│       │       └── index.ts
│       ├── components/
│       │   └── ItemForm/
│       │       ├── ItemForm.tsx
│       │       ├── types/itemForm.props.ts
│       │       ├── helpers/               # only if ItemForm owns pure view logic
│       │       ├── constants/             # only if ItemForm owns fixed values
│       │       ├── __tests__/ItemForm.test.tsx
│       │       └── index.ts
│       ├── templates/
│       │   ├── desktopItem/
│       │   └── mobileItem/
│       ├── hooks/
│       │   ├── useItemForm/
│       │   │   ├── useItemForm.ts
│       │   │   ├── types/useItemForm.ts
│       │   │   ├── hooks/
│       │   │   │   ├── useFormActions/
│       │   │   │   └── useFormSubmission/
│       │   │   ├── helpers/               # flat use-case helper files
│       │   │   │   ├── item.validateCredentials.ts
│       │   │   │   ├── types/item.validateCredentials.ts
│       │   │   │   └── __tests__/item.validateCredentials.test.ts
│       │   │   ├── constants/
│       │   │   └── index.ts
│       │   └── useItemViewport/
│       │       ├── useItemViewport.ts    # only if a reusable hook is warranted
│       │       ├── helpers/item.viewport.ts
│       │       ├── constants/
│       │       ├── types/
│       │       └── index.ts
│       ├── helpers/                       # page-wide pure operations only
│       │   ├── item.getFormViewModel.ts
│       │   ├── types/item.getFormViewModel.ts
│       │   └── __tests__/item.getFormViewModel.test.ts
│       ├── constants/
│       │   ├── item.fields.ts
│       │   ├── item.copy.ts
│       │   ├── item.routes.ts
│       │   ├── item.testIds.ts
│       │   ├── item.descriptionVisibility.ts
│       │   └── types/                     # only for genuinely shared contracts
│       └── types/
│           ├── item.props.ts
│           └── item.state.ts
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
├── item.getFormViewModel.ts
├── item.getDescriptionViewModel.ts
├── item.getFailureMessage.ts
├── types/
│   ├── item.getFormViewModel.ts
│   ├── item.getDescriptionViewModel.ts
│   └── item.getFailureMessage.ts
└── __tests__/
    ├── item.getFormViewModel.test.ts
    ├── item.getDescriptionViewModel.test.ts
    └── item.getFailureMessage.test.ts

constants/
├── item.fields.ts
├── item.copy.ts
├── item.routes.ts
├── item.testIds.ts
├── item.descriptionVisibility.ts
└── types/
~~~

Forbidden in a migrated helpers/constants directory: `helpers/index.ts`, `constants/index.ts`, `helpers/getSomething/getSomething.ts`, `constants/someConstant/index.ts`, `helpers.ts`, `common.ts`, `utils.ts`, `misc.ts`.

A self-contained business operation or reusable component is different: its own public `index.ts` can remain. Do not delete **domain/page/component/hook** public entrypoints while removing helper/constant barrels.

## 5. File naming, types, placement and unused-code removal

Use filenames that identify owner **and** use case: `login.getPasswordViewModel.ts`, `login.validateCredentials.ts`, `login.fieldIds.ts`, `recovery.classifyFailure.ts`, `person.request.mapper.ts`. Avoid anonymous `getData.ts` or `helpers.ts`. Exported functions have named callable contracts in the nearest `types/<useCase>.ts`. Reuse the nearest authoritative page/entity types via indexed access, `Pick` and `Omit` instead of duplicating state shapes.

~~~ts
// hooks/useItemForm/helpers/types/item.validateCredentials.ts
import type { ItemCredentials } from '../../../../types/item.state';
export type ValidateCredentials = (input: ItemCredentials) => ItemFieldErrors;

// hooks/useItemForm/helpers/item.validateCredentials.ts
import type { ValidateCredentials } from './types/item.validateCredentials';
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
// constants/item.fields.ts
export const ITEM_FIELD = {
  TITLE: 'title',
  DESCRIPTION: 'description',
} as const;

export type ItemField = (typeof ITEM_FIELD)[keyof typeof ITEM_FIELD];

export const ITEM_FIELD_ORDER: readonly ItemField[] = [
  ITEM_FIELD.TITLE,
  ITEM_FIELD.DESCRIPTION,
];

export const ITEM_FIELDS = {
  [ITEM_FIELD.TITLE]: { label: 'Title address', placeholder: 'Title', autoComplete: 'username' },
  [ITEM_FIELD.DESCRIPTION]: { label: 'Description', placeholder: 'Description', autoComplete: 'current-description' },
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
└── useItemForm/
    ├── useItemForm.ts
    ├── types/useItemForm.ts
    ├── hooks/
    │   ├── useFormActions/
    │   │   ├── useFormActions.ts
    │   │   ├── types/useFormActions.ts
    │   │   └── index.ts
    │   └── useFormSubmission/
    │       ├── useFormSubmission.ts
    │       ├── helpers/
    │       │   ├── item.validateCredentials.ts
    │       │   ├── item.getErrorMessage.ts
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
// hooks/useItemViewport/constants/item.viewport.ts
export const PHONE_VIEWPORT = '(max-width: 720px)';

// hooks/useItemViewport/helpers/item.getMobileSnapshot.ts
import { _isFunction } from '@Base/utilities/isFunction';
import { _isUndefined } from '@Base/utilities/isUndefined';

// hooks/useItemViewport/useItemViewport.ts
export const useItemViewport: UseItemViewport = () => {
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
// constants/item.descriptionVisibility.ts
export const DESCRIPTION_VISIBILITY = {
  HIDDEN: 'HIDDEN',
  VISIBLE: 'VISIBLE',
} as const;

export const DESCRIPTION_VIEW = {
  [DESCRIPTION_VISIBILITY.HIDDEN]: {
    inputType: 'description',
    toggleLabel: 'Show description',
    icon: 'eye',
  },
  [DESCRIPTION_VISIBILITY.VISIBLE]: {
    inputType: 'text',
    toggleLabel: 'Hide description',
    icon: 'eye-off',
  },
} as const;

// helpers/item.getDescriptionViewModel.ts
export const getDescriptionViewModel: GetDescriptionViewModel = (isVisible) => {
  if (isVisible) {
    return DESCRIPTION_VIEW[DESCRIPTION_VISIBILITY.VISIBLE];
  }
  return DESCRIPTION_VIEW[DESCRIPTION_VISIBILITY.HIDDEN];
};
~~~

`LOGIN_COPY`, field definitions, visible/hidden UI, submission text, selected icon and layout variants must not be computed by ternaries in components. A view helper returns the required typed fields. Static variants may be looked up via a typed `Record`; behavior uses a strategy `Map`. Always assign a stable DOM test ID to resulting elements.

## 11. Branching, registries and strategy/factory execution

**Forbidden:** ternary syntax `condition ? yes : no` anywhere in application code; nested conditional props; branching chains of `if / else if` for an extendable set of modes; `switch` for variant dispatch. Guard clauses and a simple `if / else` for a true binary case are allowed.

For three or more variants, use a `Map<Mode, Strategy>`, keyed by a typed enum/constant-derived value. Put the registry next to its owning use case (usually under `constants/<useCase>.strategies.ts` when it is a static map, with functions defined in helpers/services). Do not create a universal strategy registry for unrelated behavior.

~~~ts
export const PROCESS_MODE = {
  CREATE: 'CREATE',
  UPDATE: 'UPDATE',
  ARCHIVE: 'ARCHIVE',
} as const;
export type ProcessMode = (typeof PROCESS_MODE)[keyof typeof PROCESS_MODE];

type ProcessStrategy = (proof: string) => Promise<void>;

const strategies = new Map<ProcessMode, ProcessStrategy>([
  [PROCESS_MODE.CREATE, validateInvitation],
  [PROCESS_MODE.UPDATE, validateReset],
  [PROCESS_MODE.ARCHIVE, validateForgot],
]);

export const runProcessStrategy: RunProcessStrategy = (mode, proof) => {
  const strategy = strategies.get(mode);
  if (!strategy) {
    throw new UnsupportedProcessModeError(mode);
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
export const ITEM_TEST_IDS = {
  PAGE: 'item-page',
  FORM: 'item-form',
  TITLE_FIELD: 'item-title',
  DESCRIPTION_FIELD: 'item-description',
  DESCRIPTION_VISIBILITY: 'item-description-visibility',
  SUBMIT: 'item-submit',
  FEEDBACK: 'item-feedback',
} as const;

// In a domain component:
<Button dataTestId={ITEM_TEST_IDS.SUBMIT} type="submit">Sign in</Button>
<InputControl dataTestId={ITEM_TEST_IDS.TITLE_FIELD} />

// Native page-owned markup:
<form data-testid={ITEM_TEST_IDS.FORM}>...</form>
~~~

**Unique and deterministic:** every test ID should be unique among currently rendered elements; repeated list items derive IDs from a stable non-sensitive business key (`... + '-' + record.gid`), never array index, text label or a random/`useId()` value. Names must not expose passwords, email addresses, tokens, phone numbers or other PII. A component rendered more than once must receive distinct caller-provided IDs; don't use a hardcoded default that creates duplicates. Use an owner prefix and role suffix (e.g. `login-submit-spinner`). For a test needing to choose among repeated widgets, it must be able to locate the specific element using page/widget ancestry and a meaningful identifier.

**Accessibility is independent:** semantic role, label, keyboard activation and focus order remain required. Test IDs are for stable automation, not a substitute for accessible names or proper `id`/`htmlFor` relationships. Check every rendered page DOM node (excluding third-party implementation internals) in the migration audit; document justified exceptions.

## 13. Test ownership and regression policy

Place tests beside the **owning unit**, never in an unrelated global bucket.

~~~text
pages/item/components/ItemForm/
├── ItemForm.tsx
├── types/
├── __tests__/
│   └── ItemForm.test.tsx
└── index.ts

pages/item/hooks/useItemForm/helpers/
├── item.validateCredentials.ts
├── types/
└── __tests__/
    └── item.validateCredentials.test.ts

pages/item/__tests__/
└── item.flow.test.tsx
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

<Card variant="outline" padding="lg" dataTestId="inventory-card">
  <Typography as="h2" variant="heading" size="2xl" dataTestId="inventory-heading">Enquiry</Typography>
  <Text tone="muted" size="sm" dataTestId="inventory-hint">Complete the item details.</Text>
  <Form layout="stack" spacing="md" dataTestId="inventory-form" onSubmit={onSubmit}>
    <Input label="Title" controlSize="md" dataTestId="inventory-title" />
    <Button variant="primary" size="lg" type="submit" dataTestId="inventory-submit">Send</Button>
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
auth/components/filterField/
├── FilterField.tsx
├── constants/filterField.testIds.ts
├── renderer/
│   ├── filterField.label.tsx
│   ├── types/filterField.label.ts
│   └── __tests__/filterField.label.test.tsx
├── styles/filterField.module.scss
├── types/filterField.props.ts
├── __tests__/FilterField.test.tsx
└── index.ts
~~~

The grouped label currently uses a forbidden ternary. The correct target pattern is:

~~~tsx
// renderer/filterField.label.tsx
import { FieldLabel, Row } from '@Components';
import styles from '../styles/filterField.module.scss';
import type { RenderFilterFieldLabel } from './types/filterField.label';

export const renderFilterFieldLabel: RenderFilterFieldLabel = (props) => {
  const label = (
    <FieldLabel htmlFor={props.id} dataTestId={props.ids.LABEL}>
      {props.label}
    </FieldLabel>
  );

  if (props.grouped) {
    return (
      <Row className={styles.groupedLabel} dataTestId={props.ids.LABEL_ROW}>
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
// auth/components/taskStatusVisual/constants/taskStatusVisual.kind.ts
export enum TaskStatusKind {
  CREATED = 'created',
  SUCCESS = 'success',
  STALE = 'stale',
  INVALID = 'invalid',
  CHECKING = 'checking',
  UNAVAILABLE = 'unavailable',
  REJECTED = 'rejected',
}
~~~

~~~ts
// constants/taskStatusVisual.badge.ts
import { TaskStatusKind } from './taskStatusVisual.kind';

export enum TaskBadgeKind {
  CONFIRMED = 'confirmed',
  INVALID = 'invalid',
  NONE = 'none',
}

export const TASK_STATUS_BADGE: Readonly<Record<TaskStatusKind, TaskBadgeKind>> = {
  [TaskStatusKind.CREATED]: TaskBadgeKind.CONFIRMED,
  [TaskStatusKind.SUCCESS]: TaskBadgeKind.CONFIRMED,
  [TaskStatusKind.STALE]: TaskBadgeKind.INVALID,
  [TaskStatusKind.INVALID]: TaskBadgeKind.INVALID,
  [TaskStatusKind.CHECKING]: TaskBadgeKind.NONE,
  [TaskStatusKind.UNAVAILABLE]: TaskBadgeKind.NONE,
  [TaskStatusKind.REJECTED]: TaskBadgeKind.NONE,
};
~~~

~~~tsx
// renderer/taskStatusVisual.badge.tsx
import { Icon, Text } from '@Components';
import { TASK_STATUS_BADGE, TaskBadgeKind } from '../constants/taskStatusVisual.badge';
import styles from '../styles/taskStatusVisual.module.scss';
import type { RenderTaskStatusBadge } from './types/taskStatusVisual.badge';

export const renderTaskStatusBadge: RenderTaskStatusBadge = (kind) => {
  const badge = TASK_STATUS_BADGE[kind];
  if (badge === TaskBadgeKind.CONFIRMED) {
    return <Icon name="check-circle-2" decorative className={styles.confirmedBadge}/>;
  }
  if (badge === TaskBadgeKind.INVALID) {
    return <Text as="span" className={styles.invalidBadge} aria-hidden="true">×</Text>;
  }
  return null;
};
~~~

**Before this example compiles**, add the required shared `Icon` abstraction, typed icon identifiers, and the renderer's callable type under `renderer/types/`; the snippets now show separate, correctly imported enum and mapping files. Renderer functions are pure. Every dynamic icon, color and case is tested. Do not import `lucide-react` in domain/shell. Do not store a React component or JSX expression under a generic helper path when it is a renderer.

**CSS colors:** `fill="#039C90"` and `color="#fff"` are forbidden. Use `styles.confirmedBadge` with `color: var(--gh-color-on-success)` and `fill: var(--gh-color-success)`. Add any new semantic token (e.g. `--gh-color-on-success`) in **both** `@Styles/themes/light.css` and `dark.css` before use. Background, border, title, ring, illustration and disabled colors all follow the same rule; no new theme palette in domain SCSS.

## 26. LoginLayout exact decomposition and theme-option renderer

~~~text
auth/components/itemLayout/
├── ItemLayout.tsx
├── constants/
│   ├── itemLayout.variant.ts
│   ├── itemLayout.copy.ts
│   ├── itemLayout.appearance.ts
│   └── itemLayout.testIds.ts
├── helpers/
│   ├── itemLayout.className.ts
│   ├── itemLayout.mobileDescription.ts
│   ├── types/
│   └── __tests__/
├── renderer/
│   ├── itemLayout.appearanceOption.tsx
│   ├── itemLayout.visual.tsx
│   ├── types/
│   └── __tests__/
├── styles/itemLayout.module.scss
├── types/itemLayout.props.ts
├── __tests__/ItemLayout.test.tsx
└── index.ts
~~~

~~~ts
// constants/itemLayout.variant.ts
export enum ItemLayoutVariant {
  LIST = 'list',
  DETAIL = 'detail',
  EDIT = 'edit',
}

// constants/itemLayout.copy.ts
export const ITEM_LAYOUT_COPY = {
  DEFAULT_EYEBROW: 'YOUR WORKSPACE',
  EDIT_EYEBROW: 'EDIT ITEM',
  SCHOOL_NAME: 'Workspace',
} as const;

// constants/itemLayout.appearance.ts (shared theme values/icons are imported)
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

domain/auth/components/itemLayout/styles/itemLayout.module.scss
domain/auth/components/filterField/styles/filterField.module.scss
domain/auth/components/taskStatusVisual/styles/taskStatusVisual.module.scss
domain/auth/pages/forgotPassword/components/itemEditor/styles/itemEditor.module.scss
domain/auth/pages/itemDetails/styles/itemDetails.module.scss
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
/* taskStatusVisual.module.scss; both tokens must exist in @Styles light/dark themes. */
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
  ENTRY: '/entry',
  LEGACY_ENTRY: '/old-entry',
  HELP: '/help',
  ITEM_EDIT: '/items/edit',
  ITEM_CREATE: '/items/new',
  LEGACY_ITEM: '/old-items',
  WORKSPACE: '/workspace',
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
                ├── /workspace
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



## 30. AuthProvider and recovery-flow session: real units, one state owner, zero secret leakage

### 30.1 Provider directory

A provider is a functional **unit**, not a single do-everything file:

~~~text
domain/resource/core/
├── providers/
│   └── resourceProvider/
│       ├── ResourceProvider.tsx
│       ├── hooks/
│       │   └── useResourceSession/
│       │       ├── useResourceSession.ts
│       │       ├── helpers/
│       │       ├── constants/
│       │       ├── types/
│       │       └── __tests__/
│       ├── helpers/
│       │   └── resourceProvider.requireContext.ts
│       ├── constants/resourceProvider.messages.ts
│       ├── types/
│       │   ├── resourceProvider.props.ts
│       │   └── resourceProvider.context.ts
│       ├── __tests__/ResourceProvider.test.tsx
│       └── index.ts
├── hooks/useOperationFlowBoundary/
├── session/
│   └── operationFlowSession/
│       ├── operationFlowSession.ts
│       ├── constants/operationFlowSession.keys.ts
│       ├── constants/operationFlowSession.lifetime.ts
│       ├── helpers/operationFlowSession.parseMarker.ts
│       ├── helpers/operationFlowSession.isActive.ts
│       ├── types/operationFlowSession.ts
│       └── __tests__/
└── guards/
    ├── publicRoute/
    └── protectedRoute/
~~~

`AuthProvider.tsx` should create/provide context, not mix nested `try/catch`, HTTP refresh fallback, pending request deduplication, sign-in and sign-out in the same 80-line component. The named `useAuthSession` hook owns the **single** account/initialized/loading state and cohesive actions, with testable request helpers/classifiers where appropriate. The context contract is typed in `types/`; the exported `useAuth` consumer hook throws if used outside its provider, tested at the provider boundary.

**No duplicated session state.** A public route must not bootstrap session just because the provider mounted; a protected route requests session initialization. No access tokens in JS storage; session remains HttpOnly-cookie based. On logout remove cached account state even when the network request fails; on login and refresh handle stale responses and duplicate requests. Preserve CSRF and session continuity rules established by the backend rather than trying to infer them from redirects.

### 30.2 Recovery proof adapter

The current `auth/core/session/recoveryFlowSession.ts` is **not dead simply because it has no JSX**. It holds security-sensitive lifecycle behavior. Classify it as a browser session adapter, with pure helpers for parsing/validation and constants for storage keys, limits and purpose enum. The only permitted persisted continuation is a short-lived **nonsecret grant marker** when required to resume an HttpOnly-cookie grant. Raw password recovery/invitation bearer proof stays in memory, is removed from URL fragments after exchange and is cleared on leaving/switching the flow; it is never written to localStorage/sessionStorage or printed in logs.

Do **not** split a stateful browser adapter into unrelated globally mutable "helpers". An adapter entry may legitimately group coordinated `read`, `write`, `clear` operations; split out pure suboperations when independently testable. Test invalid/malformed markers, expiry, storage exceptions, reload, wrong purpose, URL navigation, transient server errors and server-owned exhausted-attempt decisions. Do not count attempts on the frontend.

### 30.3 Full `useRecovery` split (currently not done)

~~~text
pages/operation/hooks/useOperation/
├── useOperation.ts                 # returns one cohesive typed interface
├── hooks/
│   ├── useOperationForm/           # one form-state owner
│   ├── useOperationLink/           # hash proof handling + cleanup + sequence
│   ├── useOperationGrant/          # grant status from server
│   ├── useOperationSubmission/     # validation + mutation + duplicate protection
│   └── useOperationRetry/          # retry orchestration
├── helpers/
│   ├── operation.classifyFailure.ts
│   ├── operation.getViewModel.ts
│   ├── operation.validatePassword.ts
│   ├── types/
│   └── __tests__/
├── constants/
│   ├── operation.state.ts
│   ├── operation.copy.ts
│   ├── operation.testIds.ts
│   └── operation.routes.ts         # use common paths instead of duplicating values
├── types/useOperation.ts
└── __tests__/useOperation.test.ts
~~~

**One owner:** only one hook owns each mutable form field, fieldErrors and request-in-flight flag; subordinate hooks accept typed state and actions rather than creating second mirrors of it. If a child only transforms inputs or reads a map, it is a *helper*, not a hook. Hooks are permitted to be nested under a parent hook's `hooks/` to group focused effects. They must be called unconditionally. Do not add five hooks merely to shrink one file while leaving them mutually coupled through many setters.

Migrate the existing direct `invitation ? API_A : API_B` or `purpose === 'invitation'` branches to a named enum plus a typed strategy registry. The validation and exchange sequence should use a `Map<RecoveryPurpose, RecoveryOperation>` **when selecting between registered API strategies**, with explicit missing-strategy handling. View state labels/copy/icons use data maps, not function maps. Strictly preserve: verify link before displaying form; clear proof fragment promptly; grant verification after refresh; two server-enforced unsuccessful password attempts; expired/missing/unavailable/attempts-exhausted distinct states; transient failures not counted by guesswork; and success clearing credentials and old sessions.

## 31. `useShell`: one local state/actions interface and persistent layout

The shell's current `useShellState` and `useShellActions` are two peer hooks with extra setter plumbing, later combined by `useShell`. **Target:** one `hooks/useShell/useShell.ts` owns state and named actions and returns one typed model. Nested hook units are allowed for independently meaningful storage sync or theme integration, **not** for duplicating setters or making redundant open/close wrappers.

~~~text
workspace/
├── Workspace.tsx                           # presentation with <Outlet/>
├── containers/workspaceContainer/
│   ├── WorkspaceContainer.tsx
│   ├── constants/workspaceContainer.queryOptions.ts
│   ├── helpers/workspaceContainer.createQueryClient.ts
│   └── types/
├── hooks/useWorkspace/
│   ├── useWorkspace.ts
│   ├── hooks/
│   │   ├── usePinnedAppsPersistence/
│   │   └── useSidebarPreference/
│   ├── helpers/workspace.togglePinnedPath.ts
│   ├── types/useWorkspace.ts
│   └── __tests__/useWorkspace.test.ts
├── components/
│   ├── appRoutes/{renderer,types,__tests__}/
│   ├── catalogLauncher/
│   │   ├── CatalogLauncher.tsx
│   │   ├── hooks/useCatalogLauncher/
│   │   ├── helpers/workspace.filterCatalog.ts
│   │   ├── renderer/catalogLauncher.app.tsx
│   │   ├── constants/catalogLauncher.testIds.ts
│   │   ├── styles/catalogLauncher.module.scss
│   │   ├── types/
│   │   └── __tests__/
│   ├── protectedLayout/
│   ├── sidebar/
│   ├── topbar/
│   └── ...
└── registry/
~~~

If `useBoolean` or `useBooleanState` is available from `@Components/hooks`, use its stable `[value, setTrue, setFalse]` actions in the shell model instead of duplicating `openSidebar = useCallback(() => setSidebarOpen(true))` and `closeSidebar = ...`. **Current availability must be checked:** do not add imports for imaginary exports. A named `openCatalog` action is still warranted because it opens the catalog **and** closes the sidebar; similarly `openAccountPreferences` may legitimately coordinate several UI states.

~~~ts
// hooks/useWorkspace/helpers/workspace.togglePinnedPath.ts
import { _filter } from '@Base/utilities/filter';
import { _includes } from '@Base/utilities/includes';

export const togglePinnedPath = (pinned: readonly string[], path: string): string[] => {
  if (_includes(pinned, path)) {
    return _filter(pinned, (entry) => entry !== path);
  }
  return [...pinned, path];
};
~~~

This helper is pure and tested. `_filter`/`_includes` must first exist as narrow `@Base` exports. **Never perform `writePinnedApps` / `writeSidebarCollapsed` inside a `setState(previous => ...)` updater**: React StrictMode may invoke the updater more than once. Keep updaters pure and persist through one controlled effect or a separately tested external-store abstraction. Use lazy initialization from storage with browser exceptions handled; no duplicated `pinned` copies. Keep the shell mounted across protected child navigation and ensure only its Outlet page changes.

### 31.1 QueryClient options are constants, creation is a factory helper

The current module-level construction mixes immutable options and runtime instance creation. Store defaults in the owning container constants, and construct once through a helper. Do **not** create a QueryClient every render.

~~~ts
// workspace/containers/workspaceContainer/constants/workspaceContainer.queryOptions.ts
export const SHELL_QUERY_OPTIONS = {
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
} as const;

// workspace/containers/workspaceContainer/helpers/workspaceContainer.createQueryClient.ts
import { QueryClient } from '@tanstack/react-query';
import { SHELL_QUERY_OPTIONS } from '../constants/workspaceContainer.queryOptions';

export const createWorkspaceQueryClient = (): QueryClient => new QueryClient(SHELL_QUERY_OPTIONS);
~~~

A constructor is **not static data** and does not belong in `constants/` or `renderer/`. For the current SPA, a stable app-instance is sufficient; future SSR must scope the QueryClient per request/runtime. Query options, retry rules and timeouts must be named constants, not bare numbers in component markup.

### 31.2 AppLauncher must be decomposed, not relabeled

Current `shell/layout/AppLauncher.tsx` still owns dialog, search state, effect, filtering, pins, UI cards, native input/buttons, direct icons, literal CSS class branches and render JSX; this is a **high-priority violation**.

- `components/appLauncher/hooks/useAppLauncher`: single owner of query/search/category state, reset-on-open effect and typed actions.
- `helpers/shell.filterCatalog.ts`: pure case normalization, `_trim`, `_filter`, `_includes` and `_size` against implemented-only app descriptors.
- `renderer/appLauncher.app.tsx`: one catalog-item view returning accessible shared Card/AppLink/Button/Icon/Text primitives; no anonymous event closure or business logic.
- `constants/appLauncher.copy.ts` / `appLauncher.testIds.ts` / `appLauncher.category.ts`: display strings, filters and stable IDs.
- `styles/appLauncher.module.scss`: scoped layout consuming `@Styles`, not raw literal class names.
- `AppLauncher.tsx`: typed, presentational composition from these units only.

Remove stale façade `shell/layout/AppLauncher.tsx` and old exports after caller migration; verify pinned-navigation persistence and dialog focus trapping. **Only implemented apps** belong in launcher and navigation. People, Accounts, Integrations and Settings are not functional domain routes and may not appear as clickable features. All map/filter/find/length/key transformations use tested `@Base` leaf utilities where applicable.

## 32. Lodash / Base, zero ternaries, meaningful class merging and readable formatting

Use `@Base/utilities/<operation>` as the only application gateway to Lodash functions. **Required leaf exports when consumed** include `_map`, `_filter`, `_find`, `_trim`, `_includes`, `_keys`, `_size`, `_isUndefined`, `_isFunction` and `_isEmpty`. Each leaf has an explicit tested type and does not eagerly import the whole Lodash barrel (important for Jest CJS compatibility). Only `packages/base` can import `lodash-es` directly. Do not make a `normalizeLoginEmail` wrapper that only delegates to `_trim`. `_isValid` is not a standard Lodash function; a true business validation predicate lives under domain/business helpers.

~~~tsx
import { _map } from '@Base/utilities/map';
import { renderAppRoute } from './renderer/appRoutes.route';

// No direct .map in AppRoutes
const landing = _map(PUBLIC_SITE_PATHS, renderPublicLandingRoute);
const publicPages = _map(publicRoutes, renderAppRoute);
const workspacePages = _map(workspaceRoutes, renderAppRoute);
~~~

Avoid raw `Object.keys`, `Object.entries`, `.find`, `.filter`, `.map`, `.includes` and `.trim` where the equivalent authorized Base leaf exists. **Don't** apply Lodash to platform-only APIs (React hooks, URL, streams, Intl, DOM event registration). A helper must add real policy, not one extra wrapper layer.

**No ternary operator in any owned application TS/TSX**, including hooks, renderers, classes, styles, routing, icon selection and helper returns. Enforce using AST `ConditionalExpression` lint, not a naïve regex that rejects legitimate optional TS property syntax (`type?: string`). No inline `mobileDescription ?? description` in JSX: an explicit named helper owns defaulting and intentional empty-string handling. An inline `&&` is acceptable only for an already-computed, single-purpose visibility flag; when it selects a different JSX tree, use a renderer. No deeply nested conditional expressions hidden inside closures.

Use `cn` imported from `@Base/utilities/cn` to combine **CSS module** classes: `cn(styles.root, styles.selected)`. The class decision itself may be a `helper` that takes typed variant enum. No raw `'is-active'`, magic global BEM strings, manual string concatenation or arbitrary colors directly in migrated UI. Do not use `className={styles['computed' + value]}` to bypass a typed variant map.

**Readability:** normal two-space TypeScript indentation; descriptive exported names; sorted/import groups (React/router → `@Base`/`@Components`/`@Business` → domain-local); sensible multiline JSX; no squeezed one-line logic, repeated props, deep anonymous callbacks, inline async handlers or giant JSX return bodies. Reusable handler logic belongs in hooks. A renderer may receive a prebound stable action and forward it as `onClick={onSelect}`; do not create closures for state mutations in every list row.

## 33. Test IDs, route IDs, icon names and UI variants: no literal leakage

Every page or reusable feature defines one typed `constants/<owner>.testIds.ts`. Test ID naming follows:

~~~text
<domain>-<screen>-<component>-<element>-<useCase>

auth-item-form-input-title
auth-item-form-input-description
auth-item-form-button-submit
auth-item-layout-button-theme-system
auth-forgot-description-form-button-resend
auth-forgot-description-form-alert-feedback
auth-process-form-input-confirmation
shell-app-launcher-dialog-input-search
shell-app-launcher-card-button-pin-dashboard
dashboard-home-card-status-tasks
~~~

~~~ts
// pages/item/constants/item.testIds.ts
export const ITEM_TEST_IDS = {
  TITLE: 'auth-item-form-input-title',
  DESCRIPTION: 'auth-item-form-input-description',
  SUBMIT: 'auth-item-form-button-submit',
  AUTH_ERROR: 'auth-item-form-alert-feedback',
} as const;
~~~

No `dataTestId="forgot-password-submit"` or `data-testid="catalog-search"` written directly in a migrated feature. Instead reference a named constant. For dynamic entries use a **central factory**, backed by a stable non-sensitive path/GID and an explicit prefix; no email, person name, password, token, index or random `useId` embedded in the selector. Test IDs are unique among simultaneously rendered elements; a hidden second form does not excuse duplicate IDs. Shared component `dataTestId` forwards to the native leaf `data-testid`. Composite children use documented stable suffixes, and semantic `id`/`htmlFor` must still work. No CSS based on automation IDs.

**Typed variant convention:** domain-specific variants are enums (`AuthLayoutVariant.RECOVERY`, `AuthStatusKind.EMAIL`, `RouteAccess.PROTECTED`). Generic Button variants/size and Typography variants may use the shared component's documented string-union interface but should be supplied through named static descriptors when repeated. Do not write direct route literals `to="/login"`; use `APP_ROUTE.LOGIN` from the single canonical route contract. Do not write direct Lucide JSX in a domain; use `@Components/icons`/`Icon`. No magic icon size/color in feature JSX; expose semantic icon size/intent from the shared icon contract and CSS tokens.

### Required tests by owner

- `renderer/__tests__`: each JSX branch, unsupported status, icon and aria semantics, stable test IDs.
- `helpers/__tests__`: pure class/URL/validation/selection logic, edge cases and missing/null inputs.
- `components/<unit>/__tests__`: shared primitives forward `dataTestId` to the native leaf, focus order, label relations, submit and loading behavior, correct native role.
- `hooks/use<Feature>/__tests__`: state transitions, stale async results, StrictMode cleanup, persistence without duplicate writes, timers, viewport matching and recovery sequence.
- `auth/core/providers/authProvider/__tests__`: protected-only initialization, login/refresh/reload/logout and duplicate in-flight reload protection.
- `shell/components/protectedLayout/__tests__`: same mounted Shell element across child protected navigation.
- `publicSite`: all current SEO routes, enquiry dialog, forms and honest mailto behavior; no fabricated backend calls.
- **Browser visual/a11y**: light/dark, 320–720px/mobile/tablet/desktop, short viewport, autofocus/autofill, high zoom and keyboard input.



## 34. Full file-by-file audit, dead domain deletion and legacy CSS migration ledger

**Do not infer that prior CI success implies Revision 7 compliance.** Read all frontend source paths and their call sites. Before implementation record an inventory with: path, owner, live importers, feature behavior, DOM/test IDs, status/route strings, JSX ternary and other conditional branches, native elements, JSX-returning function location, direct Lucide/Lodash imports, CSS classes/colors, provider state/side effects, tests, and exact move/remove plan. This audit belongs in the implementation work, not a claim that this docs-only commit scanned every line. Affected migrations require tests first or an explicit regression plan.

The following **existing, unused UI domains** are removal candidates explicitly requested for deletion:

~~~text
frontend/packages/domain/accounts/
frontend/packages/domain/people/
frontend/packages/domain/integrations/
frontend/packages/domain/settings/
~~~

They are still present at the time of the documentation audit. Before deletion, verify all static/dynamic imports, domain exports, tests, links, app catalog, persisted pinned paths, search indexes, generated dependencies and route registries. Then delete each directory and unused corresponding presentation references **in the same code migration**, without removing real reusable `@Business/accounts` or other backend/business clients needed by authentication. Do not reintroduce phantom People/Accounts/Integrations/Settings pages or "coming soon" clickable cards. If a domain has a real live route or critical use, disclose that blocker and retain it until safe to remove.

**Observed concrete debt (selected verified files, not a complete repository-wide scan):**

| Current file/area | Exact issue | Required change |
|---|---|---|
| `auth/components/authStatusArt/AuthStatusArt.tsx` | Raw string states, direct Lucide, JSX ternary, hard-coded badge fill and foreground | Enums/metadata, scoped status renderer, shared Icon and named theme colors |
| `auth/components/credentialField/CredentialField.tsx` | GroupLabel JSX ternary, native `div/label`, BEM strings and inline error layout | `renderer/credentialField.label.tsx`, shared label/row/error controls, CSS Module |
| `auth/components/loginLayout/LoginLayout.tsx` | Default text literals, appearance array, native button + native .map, raw routes, class ternary, inline mobile-description fallback and direct Lucide | Constants, typed theme options, `renderer/loginLayout.appearanceOption.tsx`, helpers, primitives and module SCSS |
| `auth/components/loginLayout/styles/optionOneAuth.css` | Large mixed auth screens, hard-coded colors, theme/style values, responsive and status-art selectors in one global file | `loginLayout.module.scss` and separate credential/status/forgot/recovery modules; generic tokens to `@Styles` |
| `auth/pages/forgotPassword/components/forgotPasswordForm/ForgotPasswordForm.tsx` | Direct icons, notice JSX ternary, conditional class string, native controls and status text/URLs | Typed view model/renderer, named ID/route/copy constants and shared primitives |
| `auth/pages/recovery/Recovery.tsx` | Ternary title/description/form selection, raw `variant="recovery"` and `className="rdr-recovery"` | Typed enum, `recovery.getViewModel`, `renderer/recovery.content.tsx`, CSS Module |
| `auth/pages/recovery/hooks/useRecovery/useRecovery.ts` | Combined proof parsing, grant, form state, submission, API verification and retry | Typed nested hooks with **one state owner**, pure local helpers and security regression tests |
| `auth/core/guards/RecoveryFlowBoundary.tsx` | Side-effect component returns null | Convert to `useRecoveryFlowBoundary` hook and invoke under active Router |
| `auth/core/guards/RequireSession.tsx` | PublicOnly/RequireSession mixed, raw loader UI, ternary and literal redirect | Typed PublicRoute/ProtectedRoute/RouteAccess, shared loader and central paths |
| `auth/core/providers/AuthProvider.tsx` | Flat context/provider and multi-responsibility session actions, network request retries | Provider unit + `useAuthSession`, typed helpers, tests |
| `auth/core/session/recoveryFlowSession.ts` | Proof/grant adapter holds purpose/expiry/storage/constants in one file | Preserve security, extract named constants, pure helpers, types/tests under owned folder |
| `shell/components/appRoutes/AppRoutes.tsx` | Native `.map`, route JSX mapping and route-policy composition | `@Base` `_map`, typed route registries and `renderer/appRoutes.route.tsx` |
| `shell/containers/shellContainer/ShellContainer.tsx` | Inline `QueryClient` options/runtime construction | Typed options constant, createQueryClient helper, stable per-runtime instance |
| `shell/hooks/useShellState`, `useShellActions`, `useShell` | Peer state and actions and duplicate boolean handlers | One public `useShell` model, nested persistence hooks, shared boolean actions and pure pin selector |
| `shell/layout/AppLauncher.tsx` | Large mixed search/dialog/catalog UI, native HTML, map/filter loops, raw class/route/icon | Real `components/appLauncher` unit, state hook, renderer, helpers, test IDs and module SCSS |
| `packages/base/utilities` | `cn`/`isEqual` are present; many requested Lodash leaves are not | Add tested `_map`/`_filter`/`_find`/`_trim`/`_includes`/`_keys`/`_size`/`_isUndefined`/`_isFunction`, then migrate callers |
| `packages/components` | Shared UI primitives exist, but generic semantic Box/Row/FieldLabel/Link/Icon/Status adapters are incomplete | Implement missing primitives before replacing feature-native tags |
| `packages/domain/{accounts,people,integrations,settings}` | Code folders remain although they are not implemented protected pages | Remove after safe importer/route audit |
| `frontend/scripts/check-*.mjs` | Earlier boundary/UI/style checks do not enforce all Revision 7 syntax and folder rules | Introduce AST, import, route, token, primitive and selector checks |

The supplied auth CSS includes repeated later media queries and short-window override rules. **Do not mechanically strip them** until every overlapping selector, specificity conflict, accessibility state and mobile/desktop breakpoint behavior has been reviewed. For every migrated selector, identify its DOM consumer and its replacement CSS Module class. No blanket global replacement. Keep dark/light theme tokens centralized, mobile theme system behavior and responsive fit intact.

### Audit commands and interpretation

These are discovery commands, **not** standalone correctness checks. Run at the frontend root and review each finding, including library/vendor/test files separately where necessary:

~~~bash
rg -n 'from .lucide-react.|from .lodash-es.|\\.map\\(|\\.filter\\(|\\.find\\(|\\.includes\\(|\\.trim\\(|Object\\.keys\\(' packages app
rg -n 'className="|\\bto="/|fill="#|color="#|style=\\{\\{|<div|<button|<input|<p\\b|<h[1-6]\\b' packages
rg -n 'AuthProvider|RecoveryFlowBoundary|RequireSession|PublicOnly|AppLauncher|optionOneAuth' packages
rg -n 'data-testid|dataTestId|const .*\\?|\\?.*:' packages/domain packages/shell
rg -n 'accounts|people|integrations|settings' packages/shell packages/domain
~~~

Text searches can produce false positives for type syntax or tests; **true no-ternary enforcement uses TypeScript/ESLint AST**. Use TypeScript/lint to discover unused imports, but cross-check dynamic imports, route registration and public exports before deletion.

## 35. Enforceability, migration order and definition of done

### 35.1 New/extended automated gates to implement with source migration

1. **AST no-ternary rule** across app/packages source (TypeScript `ConditionalExpression`), with explicit distinction from optional property syntax.
2. **Import boundary rules**: only `@Base` imports Lodash; only `@Components` imports Lucide and Radix leaf implementations; `@Assets` and `@Styles` depend on no domain/shell code; the narrowly scoped `app/core/routes` contract is pure data.
3. **Directory rules**: no `renderer/index.ts` / `helpers/index.ts` / `constants/index.ts`, no one-folder-per-helper, correctly named `*.tsx` renderers, no miscellaneous `pages/*.tsx` components.
4. **Semantic UI boundary**: no raw generic/native HTML in migrated feature presentation/renderer, except documented shared primitive/Router/asset owners; no component-specific global CSS in migrated modules.
5. **Style ownership**: enforce `*.module.scss`/SCSS compiler support when enabled, no feature hard-coded colors, parity of light/dark semantic tokens, canonical size/focus/contrast roles and CSS Module imports.
6. **Enum/config and IDs**: no raw domain status/path/variant literals outside definitions, test ID catalogue for each screen, unique native test IDs and deterministic repeated controls.
7. **Hooks/side effects**: state owner and lifecycle review, no storage writes inside React setState functional updater, no null-only side-effect components, no API calls in renderers.
8. **Routing**: single path dictionary, `PUBLIC/PROTECTED` policy and guest-only distinction, one protected pathless shell and no unimplemented catalog entries.
9. **Unused code/dependency removal**: no unreachable unused UI domains, dead helpers, obsolete barrels or duplicate styles after importer checks.
10. **Tests**: pure helpers/renderers; form controls and accessibility; token/link/route enums; auth security; persistent shell; CSS/browser regressions. Document which checks are automated and which require manual review.

The existing CI **does not yet enforce all ten**. Add the rules incrementally when their target code is migrated; do not add broad exemptions or claim "fully strict" with only this documentation update. If introducing a new rule would fail thousands of known legacy lines, either migrate the affected batch in the same commit or report the exact remaining violations with a tracked follow-up. Do not suppress lint errors globally.

### 35.2 Mandatory implementation sequence

1. Audit and classify frontend files, importers, routes, old test locations and stylesheet selectors; verify dead domains and their deletion impact.
2. Implement shared `@Base` Lodash leaf utilities, `@Components` icons/semantic primitives/boolean state hooks; establish centralized app route constants and Sass tooling; add focused tests.
3. Reorganize AuthProvider/core session/route guards and replace the recovery boundary component with a hook; test security and public/protected behavior.
4. Refactor AuthStatusArt/CredentialField/LoginLayout (enums, constants, renderers, SCSS modules, IDs). Preserve SVG assets and exact merged RDR artwork.
5. Refactor login, forgot-password and recovery screens with shared UI, renderers and pure helpers; **complete useRecovery split** before marking auth done.
6. Consolidate `useShell` state/actions into one owner; decompose AppLauncher, centralize route mapping and QueryClient options; delete verified unused UI domains; preserve the mounted protected shell.
7. Migrate remaining dashboard/publicSite layout and presentation to the same rules without changing verified school content, CBSE disclosure warnings, SEO paths or enquiry behavior.
8. Enable/reinforce all checks, remove obsolete exports and CSS only after tests and browser visual verification. Update the single document's ledger with actual completion evidence.

### 35.3 Required validation commands

~~~bash
cd frontend
pnpm install --frozen-lockfile
pnpm check:architecture
pnpm typecheck
pnpm exec jest --runInBand
pnpm lint
pnpm build
~~~

Also verify the deployed/prerendered public route HTML and run a **real browser** visual/accessibility matrix: 320/375/390/720px mobile/tablet, desktop and short windows, both themes and system theme, reduced motion, keyboard focus/tab/Enter, password reveal, login/password reset link validity, cooldown, form errors, autofill, zoom and screen reader labels. Do not claim that an unexecuted browser test passed merely because CI is green.

### 35.4 Definition of done, per file and per feature

A file is conformant only when its owner is obvious, it uses typed contracts, names static data in constants, puts JSX-returning variation in `renderer/`, puts pure non-JSX computations in `helpers/`, never contains ternary expressions, never reimplements generic HTML or imports vendor icons/Lodash in a domain, uses shared CSS tokens and scoped module styling, centralizes test IDs/routes, has tests where behavior warrants and has no obsolete imports. A feature is conformant only after the associated state/action hooks, providers, adapters, tests, accessibility and viewport behavior also conform.

**The document is now the single authoritative rulebook. The implementation still has known migration debt.** A docs commit is complete when the document has been updated and verified in the repository; it does **not** delete domains, move hooks, convert CSS to Sass, implement missing shared components or fix the shown JSX by itself. Those are separate source changes that must be executed and verified before they can be reported as done.

## 36. Unit test structure: exact exported name, no mock of pure code

One **top-level** `describe` per unit test file. The `describe` name must be the exact exported symbol being tested, matching case, not a prose capability:

~~~ts
// helpers/__tests__/item.getErrorMessage.test.ts
import { getItemErrorMessage } from '../item.getErrorMessage';

describe('getItemErrorMessage', () => {
  it('maps invalid credentials to safe copy', () => {
    expect(getItemErrorMessage(new Error('Invalid item')))
      .toBe('The request was rejected.');
  });
  it('hides unknown backend failures', () => {
    expect(getItemErrorMessage(new Error('internal secret')))
      .toBe('The service is currently unavailable.');
  });
});
~~~

**Never** `describe('login error mapping')` or `describe('login credential validation')`. Use `describe('validateLoginCredentials')`, `describe('LoginForm')`, `describe('useLogin')`, `describe('renderCredentialFieldLabel')`, etc. No nested `describe` for methods or error cases; use multiple named `it` cases. Different exports get separate test files.

**Never mock pure helpers, validators, selectors, renderer functions or base utilities**: execute their real input/output semantics. Even when testing an upstream caller, prefer real pure dependencies; do not replace validation with a Jest spy to make the test pass. `jest.fn` is allowed as an interaction callback spy in a component test, not as a replacement for a pure helper. Mock only genuine external IO boundaries where necessary: network/API, browser storage failures, clock or framework navigation; prefer stateful integration tests when feasible. Pure tests must not require API or auth-context mocks.

Tests belong under the owning `helpers/__tests__/`, `renderer/__tests__/`, `hooks/useX/__tests__/` or `components/X/__tests__/`. Never put tests under a nested single-function helper folder. Add targeted tests for valid, invalid, empty, malformed, case-preserving email, unchanged password, generic backend failure, and failure confidentiality. Preserve exact user-facing copy.

## 37. Page-level constants and helpers must be flat; delete `config/`

**Required login unit, reference for all future page refactors:**

~~~text
pages/item/
├── Item.tsx
├── constants/
│   ├── item.fields.ts
│   ├── item.fieldConfig.ts
│   ├── item.fieldOrder.ts
│   ├── item.validation.ts
│   ├── item.copy.ts
│   ├── item.errorRules.ts
│   ├── item.testIds.ts
│   └── types/item.fieldConfig.ts
├── helpers/
│   ├── item.getErrorMessage.ts
│   ├── item.validateCredentials.ts
│   ├── item.isValidTitle.ts
│   ├── types/
│   │   ├── item.getErrorMessage.ts
│   │   └── item.validateCredentials.ts
│   └── __tests__/
│       ├── item.getErrorMessage.test.ts
│       ├── item.validateCredentials.test.ts
│       └── item.isValidTitle.test.ts
├── renderer/                     # only JSX-producing pure selectors
├── hooks/
├── components/
├── templates/
├── types/
└── __tests__/
~~~

**No `pages/login/config/`, `helpers/index.ts`, `constants/index.ts`, `renderer/index.ts`, `helpers/getLoginErrorMessage/` or any one-folder-per-helper nesting.** Constants have no executable operations, helpers contain no JSX, and renderers return JSX without hook/IO side effects. Nested `types/` and `__tests__/` are permitted only for ownership. Name source files `<owner>.<useCase>.ts[x]`. Update every direct importer and delete obsolete wrappers/barrels after validation.

### 37.1 Typed fields and frozen static descriptors

Use exactly **one authoritative stable field identifier contract**. Values remain lowercase to match existing state/error keys; enum *members* are uppercase to make identity explicit.

~~~ts
// constants/item.fields.ts
export enum ItemField {
  TITLE = 'title',
  DESCRIPTION = 'description',
}
~~~

~~~ts
// constants/types/item.fieldConfig.ts
export interface ItemFieldConfig {
  readonly label: string;
  readonly placeholder: string;
  readonly autoComplete: 'name' | 'current-description';
}

// constants/item.fieldConfig.ts
import { ItemField } from './item.fields';
import type { ItemFieldConfig } from './types/item.fieldConfig';

export const ITEM_FIELD_CONFIG: Readonly<Record<ItemField, Readonly<ItemFieldConfig>>> = Object.freeze({
  [ItemField.TITLE]: Object.freeze({
    label: 'Title address', placeholder: 'Title', autoComplete: 'name',
  }),
  [ItemField.DESCRIPTION]: Object.freeze({
    label: 'Description', placeholder: 'Description', autoComplete: 'current-description',
  }),
});
~~~

~~~ts
// constants/item.fieldOrder.ts
import { ItemField } from './item.fields';
export const ITEM_FIELD_ORDER = Object.freeze([
  ItemField.TITLE,
  ItemField.DESCRIPTION,
] as const satisfies readonly ItemField[]);
~~~

Components and validators read `LOGIN_FIELD_CONFIG[LoginField.EMAIL]` instead of inventing a separate email/password config. Do not store a duplicate object `{ email: ..., password: ... }` with unrelated string keys. The enum can type `LoginFieldErrors = Partial<Record<LoginField, string>>`, but avoid new circular runtime imports from page type modules.

### 37.2 Validation policies are data; predicates belong to helpers

No `isValid: value => ...` functions inside `constants/login.validation.ts`. Represent rule ordering as typed **data**, and implement operations under `helpers/`:

~~~ts
// constants/item.validation.ts
import { ItemField } from './item.fields';
export enum ItemValidationKind {
  REQUIRED = 'REQUIRED',
  TITLE_FORMAT = 'TITLE_FORMAT',
}
export const TITLE_PATTERN = /^.{3,}$/;
export const ITEM_VALIDATION_RULES: Readonly<Record<ItemField, readonly ItemValidationKind[]>> = Object.freeze({
  [ItemField.TITLE]: Object.freeze([ItemValidationKind.REQUIRED, ItemValidationKind.TITLE_FORMAT]),
  [ItemField.DESCRIPTION]: Object.freeze([ItemValidationKind.REQUIRED]),
});
export const ITEM_REQUIRED_MESSAGES: Readonly<Record<ItemField, string>> = Object.freeze({
  [ItemField.TITLE]: 'Please enter your title address.',
  [ItemField.DESCRIPTION]: 'Please enter your description.',
});
export const ITEM_TITLE_FORMAT_MESSAGE = 'Enter a valid title address.';
~~~

~~~ts
// helpers/item.isValidTitle.ts
import { TITLE_PATTERN } from '../constants/item.validation';
export const isValidItemTitle = (value: string): boolean => TITLE_PATTERN.test(value);
~~~

`login.validateCredentials.ts` takes typed credentials, uses real `_trim` from `@Base/utilities/trim` to normalize **email only**, iterates `LOGIN_FIELD_ORDER` and evaluates rules through pure helper(s), returns typed field errors, and never logs/trims the password. For multiple behavioral rule strategies use a typed helper strategy registry; static policy metadata stays in constants. Tests execute real validators, not mocks.

Other exact responsibilities: `constants/login.copy.ts` for screen copy, `constants/login.errorRules.ts` for safe matcher metadata, `helpers/login.getErrorMessage.ts` for mapping unknown failures to nonleaking copy, `constants/login.testIds.ts` for stable automation IDs. The old `normalizeLoginEmail` pure wrapper is redundant with `_trim`: remove the wrapper and its obsolete test when migrating, and test the actual `_trim` leaf under `@Base`. Do not keep a compatibility `helpers/index.ts`.

## 38. Central route contracts live under `@Business/auth/routes` and `@Business/app/routes`

**Supersedes the earlier tentative `app/core/routes` placement in section 29.** Frontend navigation paths are immutable framework-neutral contracts:

~~~text
packages/business/
├── resource/routes/
│   ├── constants/resource.routes.ts       # sign in, forgot, reset, invitation, legacy auth redirects
│   └── index.ts                       # intentional package public entrypoint
└── app/routes/
    ├── constants/workspace.routes.ts        # home, dashboard and implemented workspace destinations
    ├── constants/workspace.routeAccess.ts   # PUBLIC / PROTECTED and OPEN / GUEST_ONLY
    └── index.ts
~~~

**Auth route vs API endpoint distinction:** `@Business/auth/routes` stores browser-facing `/login` and `/forgot-password`, **not** `/identity/public/auth/*` transport URLs. `@Business/app/routes` stores `/`, `/dashboard` and route access policy. No duplicating auth path strings in both modules or hand-maintained `LOGIN_ROUTES = { success: '/dashboard', forgotPassword: '/forgot-password' }` inside login pages. Reference `APP_ROUTE.DASHBOARD` on login success and `AUTH_ROUTE.FORGOT_PASSWORD` in login template links. Keep publicSite SEO route registry as the owner of school paths (no duplicate shell copy).

For route policy, use named enum `RouteAccess.PUBLIC` / `RouteAccess.PROTECTED` plus an explicit `PublicRoutePolicy.OPEN` / `PublicRoutePolicy.GUEST_ONLY`: public school pages remain accessible while signed in; login can use guest-only behavior. Shell keeps a **single persistent ProtectedLayout and Outlet**, without remounting per protected route. Keep compatibility `/loginv2` redirect and `/setup` invitation links; no new phantom domain URLs. Do not import React, router, hooks, fetch or UI components into these pure route contracts. Frontend `app` bootstrap may use the same source of truth but may not redefine it.

## 39. Revision 8 completion criteria

A conformant login unit has **no `config` directory, no helper/constant/renderer barrels, no nested helper unit folders and no normalize-wrapper**. It has enum-backed field metadata, immutable order, data-only validation rules, pure typed helper predicates, exact-name test `describe`, real pure utilities, named copy/test IDs, and no page-local route dictionary. All imports resolve to specific leaf files. Login/forgot/recovery visible behavior and server error confidentiality remain intact.

**Implementation sequence:** documentation **first**; then real Base Lodash leaves and pure auth/app route contracts; then flat login constants/helpers and their tests; then hook and component/template imports; remove old directories and aliases; run architecture, TypeScript, Jest and production build and document results. Strict zero-ternary, renderer and CSS-module migration from sections 23–35 remain required; a successful helper refactor alone is not a claim the whole frontend is fully compliant.

Only declare code migration done after tests and meaningful checks have run. For documentation-only commits, verify the actual document and report that code/CI were not yet changed.

## 40. Comprehensive tracked-file inventory and gap review (Revision 8)

**GitHub snapshot:** `0eebd64a32760f520a94da31a75d7326b0bde6c1`, `473` tracked files under `frontend/`. Every tracked file has exactly one entry below, including source, tests, styles, configs, docs and the generated lockfile. **Source reviewed** means inspected in this audit; **Structural finding** means directory evidence; **Content review required** means no assertion of source compliance. Some other source-reviewed findings were established before the snapshot and should be reconfirmed after any code change. This is a rolling snapshot, so recent commits may have resolved some previously reported issues. No static scan substitutes for browser visual, security, import reachability or CI testing.

| Tracked file (relative to `frontend/`) | Gap or verification requirement | Evidence |
|---|---|---|
| `.env.example` | Verify configuration, import aliases, scripts, build/Jest/ESLint and boundary enforcement | Config review required |
| `README.md` | Check documentation examples/references against current source | Documentation review |
| `app/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `app/main.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `app/styles.css` | Large legacy global CSS; split by feature/primitive while retaining responsive behavior | Source reviewed |
| `components.json` | Verify configuration, import aliases, scripts, build/Jest/ESLint and boundary enforcement | Config review required |
| `eslint.config.mjs` | Verify configuration, import aliases, scripts, build/Jest/ESLint and boundary enforcement | Config review required |
| `index.html` | Verify configuration, import aliases, scripts, build/Jest/ESLint and boundary enforcement | Config review required |
| `jest.config.cjs` | Verify configuration, import aliases, scripts, build/Jest/ESLint and boundary enforcement | Config review required |
| `package.json` | Verify configuration, import aliases, scripts, build/Jest/ESLint and boundary enforcement | Config review required |
| `packages/assets/README.md` | Check documentation examples/references against current source | Documentation review |
| `packages/assets/auth/AcademicPathwaysSvg.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/assets/auth/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/assets/brand/RdrSchoolMonogramSvg.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/assets/brand/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/assets/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/base/RUNBOOK.md` | Check documentation examples/references against current source | Documentation review |
| `packages/base/constants/index.ts` | Forbidden barrel; direct leaf imports | Structural finding |
| `packages/base/events/createEventEmitter/createEventEmitter.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/base/events/createEventEmitter/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/base/events/createEventEmitter/types/createEventEmitter.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/base/events/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/base/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/base/readers/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/base/readers/readProperty/__tests__/readProperty.test.ts` | Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/base/readers/readProperty/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/base/readers/readProperty/readProperty.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/base/readers/readProperty/types/readProperty.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/base/utilities/cn/__tests__/cn.test.ts` | Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/base/utilities/cn/cn.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/base/utilities/cn/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/base/utilities/cn/types/cn.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/base/utilities/find/__tests__/find.test.ts` | **Source confirmed:** correct `describe('_find')`; real pure Lodash implementation, no mock | Source reviewed |
| `packages/base/utilities/find/find.ts` | **Source confirmed:** typed pure Lodash `_find` bridge; verify compile/runtime contract and tests | Source reviewed |
| `packages/base/utilities/find/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/base/utilities/find/types/find.ts` | Review prop/enum contract, data ownership and consumer typechecks | Content review required |
| `packages/base/utilities/index.ts` | **Confirmed:** only cn/isEqual exported; Base _map/_filter/_trim etc not implemented | Source reviewed |
| `packages/base/utilities/isEqual/__tests__/isEqual.test.ts` | **Source confirmed:** pure test mocks removed; tests execute real equality | Source reviewed |
| `packages/base/utilities/isEqual/index.ts` | Check pure Lodash test/Jest compatibility | Structural finding |
| `packages/base/utilities/isEqual/isEqual.ts` | Check pure Lodash test/Jest compatibility | Structural finding |
| `packages/base/utilities/isEqual/types/isEqual.ts` | Check pure Lodash test/Jest compatibility | Structural finding |
| `packages/base/utilities/keys/__tests__/keys.test.ts` | Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/base/utilities/keys/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/base/utilities/keys/keys.ts` | **Source confirmed:** typed pure Lodash `_keys` bridge | Source reviewed |
| `packages/base/utilities/keys/types/keys.ts` | Review prop/enum contract, data ownership and consumer typechecks | Content review required |
| `packages/base/utilities/map/__tests__/map.test.ts` | **Source confirmed:** correct `describe('_map')`; real pure Lodash implementation, no mock | Source reviewed |
| `packages/base/utilities/map/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/base/utilities/map/map.ts` | **Source confirmed:** typed `_map` bridge; verify JSX callback typing | Source reviewed |
| `packages/base/utilities/map/types/map.ts` | Review prop/enum contract, data ownership and consumer typechecks | Content review required |
| `packages/base/utilities/trim/__tests__/trim.test.ts` | Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/base/utilities/trim/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/base/utilities/trim/trim.ts` | **Source confirmed:** typed `_trim` bridge; must preserve case | Source reviewed |
| `packages/base/utilities/trim/types/trim.ts` | Review prop/enum contract, data ownership and consumer typechecks | Content review required |
| `packages/business/RUNBOOK.md` | Check documentation examples/references against current source | Documentation review |
| `packages/business/accounts/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/business/accounts/services/accounts.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/business/accounts/types.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/business/api/cookies.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/business/api/http.test.ts` | Move test inside owner __tests__; Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/business/api/http.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/business/api/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/business/api/types.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/business/app/routes/__tests__/app.routes.test.ts` | Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/business/app/routes/constants/app.routeAccess.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/business/app/routes/constants/app.routes.ts` | **Source confirmed:** central app/browser paths implemented; check all consumers remove raw URLs | Source reviewed |
| `packages/business/app/routes/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/business/auth/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/business/auth/routes/__tests__/auth.routes.test.ts` | Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/business/auth/routes/constants/auth.routes.ts` | **Source confirmed:** central auth/browser paths implemented; backend API endpoint constants remain separate | Source reviewed |
| `packages/business/auth/routes/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/business/auth/services/authentication.ts` | **Confirmed:** backend API path literals; separate transport endpoint constants from frontend page routes | Source reviewed |
| `packages/business/auth/services/login/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/business/auth/services/login/login.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/business/auth/services/login/types/login.request.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/business/auth/services/login/types/login.response.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/business/auth/services/login/types/login.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/business/auth/types.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/business/auth/types/baseEntities/account.entity.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/business/auth/types/baseEntities/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/business/auth/types/baseEntities/session.entity.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/business/foundation/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/business/foundation/services/person.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/business/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/business/integrations/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/business/integrations/services/integrations.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/components/RUNBOOK.md` | Check documentation examples/references against current source | Documentation review |
| `packages/components/components/Button/Button.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/components/components/Button/config/buttonVariants.config.ts` | Config-folder ownership debt; static descriptors should be constants | Structural finding |
| `packages/components/components/Button/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/components/components/Button/types/button.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/components/components/Card/Card.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/components/components/Card/__tests__/Card.test.tsx` | Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/components/components/Card/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/components/components/Card/types/card.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/components/components/Container/Container.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/components/components/Container/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/components/components/Container/types/container.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/components/components/DataTable/DataTable.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/components/components/DataTable/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/components/components/DataTable/types/dataTable.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/components/components/Dialog/Dialog.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/components/components/Dialog/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/components/components/Dialog/types/dialog.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/components/components/EmptyState/EmptyState.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/components/components/EmptyState/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/components/components/EmptyState/types/emptyState.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/components/components/Form/Form.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/components/components/Form/__tests__/Form.test.tsx` | Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/components/components/Form/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/components/components/Form/types/form.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/components/components/FormField/FormField.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/components/components/FormField/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/components/components/FormField/types/formField.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/components/components/Input/Input.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/components/components/Input/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/components/components/Input/types/input.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/components/components/InputControl/InputControl.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/components/components/InputControl/__tests__/InputControl.test.tsx` | Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/components/components/InputControl/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/components/components/InputControl/types/inputControl.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/components/components/PageHeading/PageHeading.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/components/components/PageHeading/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/components/components/PageHeading/types/pageHeading.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/components/components/PasswordField/PasswordField.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/components/components/PasswordField/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/components/components/PasswordField/types/passwordField.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/components/components/SVGRenderer/SVGRenderer.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/components/components/SVGRenderer/__tests__/SVGRenderer.test.tsx` | Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/components/components/SVGRenderer/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/components/components/SVGRenderer/types/svgRenderer.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/components/components/SchoolMonogram.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/components/components/StatusPill/StatusPill.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/components/components/StatusPill/config/statusPillTones.config.ts` | Config-folder ownership debt; static descriptors should be constants | Structural finding |
| `packages/components/components/StatusPill/helpers/getStatusPillTone/__tests__/getStatusPillTone.test.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts}; Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/components/components/StatusPill/helpers/getStatusPillTone/getStatusPillTone.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/components/components/StatusPill/helpers/getStatusPillTone/index.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/components/components/StatusPill/helpers/getStatusPillTone/types/getStatusPillTone.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/components/components/StatusPill/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/components/components/StatusPill/types/statusPill.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/components/components/Text/Text.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/components/components/Text/__tests__/Text.test.tsx` | Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/components/components/Text/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/components/components/Text/types/text.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/components/components/TextField/TextField.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/components/components/TextField/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/components/components/TextField/types/textField.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/components/components/Typography/Typography.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/components/components/Typography/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/components/components/Typography/types/typography.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/components/hooks/query/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/components/hooks/theme.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/components/hooks/useWorkspaceTheme/helpers/resolveWorkspaceTheme/__tests__/resolveWorkspaceTheme.test.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts}; Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/components/hooks/useWorkspaceTheme/helpers/resolveWorkspaceTheme/index.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/components/hooks/useWorkspaceTheme/helpers/resolveWorkspaceTheme/resolveWorkspaceTheme.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/components/hooks/useWorkspaceTheme/helpers/resolveWorkspaceTheme/types/resolveWorkspaceTheme.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/components/hooks/useWorkspaceTheme/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/components/hooks/useWorkspaceTheme/types/useWorkspaceTheme.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/components/hooks/useWorkspaceTheme/useWorkspaceTheme.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/components/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/components/styles/primitives.css` | Review token colors, theme parity, media breakpoints and global selector ownership | Content review required |
| `packages/domain/accounts/index.ts` | Remove unused UI domain after route/importer check | Structural finding |
| `packages/domain/accounts/pages/AccountDetailsPage.tsx` | Remove unused UI domain after route/importer check | Structural finding |
| `packages/domain/accounts/pages/AccountsPage.tsx` | Remove unused UI domain after route/importer check | Structural finding |
| `packages/domain/auth/RUNBOOK.md` | Check documentation examples/references against current source | Documentation review |
| `packages/domain/auth/components/authAmbientGraphic/AuthAmbientGraphic.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/auth/components/authAmbientGraphic/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/components/authStatusArt/AuthStatusArt.tsx` | **Confirmed:** raw status literals, ternary/icon selection, inline hex colors, direct Lucide | Source reviewed |
| `packages/domain/auth/components/authStatusArt/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/components/credentialField/CredentialField.tsx` | **Confirmed:** label JSX ternary, native label/wrapper and BEM classes | Source reviewed |
| `packages/domain/auth/components/credentialField/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/components/credentialField/types/credentialField.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/components/loginLayout/LoginLayout.tsx` | **Confirmed:** inline copy/theme options, native buttons and .map, conditional class logic, direct Lucide | Source reviewed |
| `packages/domain/auth/components/loginLayout/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/components/loginLayout/styles/optionOneAuth.css` | **Confirmed:** large mixed-scope auth stylesheet; split ownership into SCSS modules | Source reviewed |
| `packages/domain/auth/core/guards/RecoveryFlowBoundary.tsx` | **Confirmed:** hook returning null used as component; convert to useRecoveryFlowBoundary | Source reviewed |
| `packages/domain/auth/core/guards/RequireSession.test.tsx` | Move test inside owner __tests__; Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/domain/auth/core/guards/RequireSession.tsx` | **Confirmed:** mixed guards and native loading/ternary; PublicRoute/ProtectedRoute | Source reviewed |
| `packages/domain/auth/core/guards/__tests__/RecoveryFlowBoundary.test.tsx` | Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/domain/auth/core/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/core/providers/AuthProvider.tsx` | **Confirmed:** provider mixes session state/actions and context; unit ownership | Source reviewed |
| `packages/domain/auth/core/session/recoveryFlowSession.ts` | **Confirmed:** security-sensitive adapter needs owned types/keys/validation helpers, no proof leakage | Source reviewed |
| `packages/domain/auth/hooks/useAuthViewportFit/__tests__/useAuthViewportFit.test.tsx` | Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/domain/auth/hooks/useAuthViewportFit/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/hooks/useAuthViewportFit/useAuthViewportFit.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/domain/auth/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/forgotPassword/ForgotPassword.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/auth/pages/forgotPassword/__tests__/ForgotPasswordPage.test.tsx` | Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/domain/auth/pages/forgotPassword/components/forgotPasswordForm/ForgotPasswordForm.tsx` | **Confirmed:** native controls, literal routes and selectors, ternary notice icon and classes | Source reviewed |
| `packages/domain/auth/pages/forgotPassword/components/forgotPasswordForm/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/forgotPassword/components/forgotPasswordForm/types/forgotPasswordForm.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/forgotPassword/constants/forgotPasswordCooldown/forgotPasswordCooldown.constant.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/domain/auth/pages/forgotPassword/constants/forgotPasswordCooldown/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/forgotPassword/containers/forgotPasswordContainer/ForgotPasswordContainer.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/auth/pages/forgotPassword/containers/forgotPasswordContainer/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/forgotPassword/containers/forgotPasswordContainer/types/forgotPasswordContainer.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/forgotPassword/helpers/getForgotPasswordViewModel/__tests__/getForgotPasswordViewModel.test.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts}; Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/domain/auth/pages/forgotPassword/helpers/getForgotPasswordViewModel/getForgotPasswordViewModel.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/auth/pages/forgotPassword/helpers/getForgotPasswordViewModel/index.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/auth/pages/forgotPassword/helpers/getForgotPasswordViewModel/types/getForgotPasswordViewModel.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/auth/pages/forgotPassword/helpers/validateForgotPasswordEmail/__tests__/validateForgotPasswordEmail.test.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts}; Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/domain/auth/pages/forgotPassword/helpers/validateForgotPasswordEmail/index.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/auth/pages/forgotPassword/helpers/validateForgotPasswordEmail/types/validateForgotPasswordEmail.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/auth/pages/forgotPassword/helpers/validateForgotPasswordEmail/validateForgotPasswordEmail.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/auth/pages/forgotPassword/hooks/useForgotPassword/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/forgotPassword/hooks/useForgotPassword/types/useForgotPassword.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/forgotPassword/hooks/useForgotPassword/useForgotPassword.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/domain/auth/pages/forgotPassword/hooks/useForgotPasswordCountdown/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/forgotPassword/hooks/useForgotPasswordCountdown/types/useForgotPasswordCountdown.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/forgotPassword/hooks/useForgotPasswordCountdown/useForgotPasswordCountdown.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/domain/auth/pages/forgotPassword/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/forgotPassword/types/forgotPassword.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/forgotPassword/types/forgotPassword.state.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/forgotPassword/types/forgotPassword.viewModel.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/forgotPassword/types/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/login/Login.tsx` | **Source confirmed:** mobile/desktop template ternary; move selection to typed renderer, preserve one form | Source reviewed |
| `packages/domain/auth/pages/login/__tests__/LoginPage.test.tsx` | Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/domain/auth/pages/login/components/LoginForm/LoginForm.test.tsx` | **Source confirmed:** exact `describe('LoginForm')` fixed; still adjacent to component instead of `__tests__`, test mixes mobile template behavior | Source reviewed |
| `packages/domain/auth/pages/login/components/LoginForm/LoginForm.tsx` | **Source confirmed, urgent:** uses undefined `LOGIN_FIELDS` after switching to `LOGIN_FIELD_CONFIG`; multiple ternaries, raw native controls, direct Lucide, legacy BEM CSS, untyped icon variants | Source reviewed |
| `packages/domain/auth/pages/login/components/LoginForm/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/login/components/LoginForm/types/loginForm.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/login/constants/login.copy.ts` | Inspect state/logic ownership, utilities, importers and typed identifiers | Content review required |
| `packages/domain/auth/pages/login/constants/login.errorRules.ts` | Inspect state/logic ownership, utilities, importers and typed identifiers | Content review required |
| `packages/domain/auth/pages/login/constants/login.fieldConfig.ts` | **Source confirmed:** frozen typed `Record<LoginField, Config>` now implemented | Source reviewed |
| `packages/domain/auth/pages/login/constants/login.fieldOrder.ts` | **Source confirmed:** immutable enum-backed field order now implemented | Source reviewed |
| `packages/domain/auth/pages/login/constants/login.fields.ts` | **Source confirmed:** enum-backed email/password identity now implemented; verify all consumers use enum | Source reviewed |
| `packages/domain/auth/pages/login/constants/login.testIds.ts` | Inspect state/logic ownership, utilities, importers and typed identifiers | Content review required |
| `packages/domain/auth/pages/login/constants/login.validation.ts` | **Source confirmed:** enum-backed rule metadata now pure; email regex corrected in newer commit; add valid/invalid Unicode/whitespace test coverage | Source reviewed |
| `packages/domain/auth/pages/login/constants/types/login.fieldConfig.ts` | Review prop/enum contract, data ownership and consumer typechecks | Content review required |
| `packages/domain/auth/pages/login/containers/loginContainer/LoginContainer.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/auth/pages/login/containers/loginContainer/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/login/containers/loginContainer/types/loginContainer.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/login/helpers/__tests__/login.getErrorMessage.test.ts` | **Source confirmed:** `describe('getLoginErrorMessage')`, unmocked real helper test now correct | Source reviewed |
| `packages/domain/auth/pages/login/helpers/__tests__/login.getValidationError.test.ts` | **Source confirmed:** `describe('getLoginValidationError')`, no pure-function mock | Source reviewed |
| `packages/domain/auth/pages/login/helpers/__tests__/login.isValidEmail.test.ts` | **Source confirmed:** `describe('isValidLoginEmail')`, no mock; limited positive domain cases | Source reviewed |
| `packages/domain/auth/pages/login/helpers/__tests__/login.validateCredentials.test.ts` | **Source confirmed:** `describe('validateLoginCredentials')`, no pure-function mock; preserve password bytes | Source reviewed |
| `packages/domain/auth/pages/login/helpers/login.getErrorMessage.ts` | **Source confirmed:** flat, safe message mapper using `_find`, no pure mock | Source reviewed |
| `packages/domain/auth/pages/login/helpers/login.getValidationError.ts` | **Source confirmed:** validation branching resides in helper and maps enum rule kinds, not in constants | Source reviewed |
| `packages/domain/auth/pages/login/helpers/login.isValidEmail.ts` | **Source confirmed:** pure predicate uses centralized email regex; add edge-case tests | Source reviewed |
| `packages/domain/auth/pages/login/helpers/login.validateCredentials.ts` | **Source confirmed:** flat helper calls Base `_trim` for email only; password preserved | Source reviewed |
| `packages/domain/auth/pages/login/helpers/types/login.getErrorMessage.ts` | Review prop/enum contract, data ownership and consumer typechecks | Content review required |
| `packages/domain/auth/pages/login/helpers/types/login.getValidationError.ts` | Review prop/enum contract, data ownership and consumer typechecks | Content review required |
| `packages/domain/auth/pages/login/helpers/types/login.validateCredentials.ts` | Review prop/enum contract, data ownership and consumer typechecks | Content review required |
| `packages/domain/auth/pages/login/hooks/useLogin/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/login/hooks/useLogin/types/useLogin.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/login/hooks/useLogin/useLogin.test.tsx` | Move test inside owner __tests__; Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/domain/auth/pages/login/hooks/useLogin/useLogin.tsx` | **Confirmed:** composes peer hooks; unify ownership under useLoginForm | Source reviewed |
| `packages/domain/auth/pages/login/hooks/useLoginActions/index.ts` | Move peer hook under owning useLoginForm hook or combine into it | Structural finding |
| `packages/domain/auth/pages/login/hooks/useLoginActions/types/useLoginActions.ts` | Move peer hook under owning useLoginForm hook or combine into it | Structural finding |
| `packages/domain/auth/pages/login/hooks/useLoginActions/useLoginActions.ts` | **Source confirmed:** field enum and if-based error clearing implemented; still peer state/action hook, should be nested under useLoginForm | Source reviewed |
| `packages/domain/auth/pages/login/hooks/useLoginForm/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/login/hooks/useLoginForm/types/useLoginForm.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/login/hooks/useLoginForm/useLoginForm.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/domain/auth/pages/login/hooks/useLoginSubmission/index.ts` | Move peer hook under owning useLoginForm hook or combine into it | Structural finding |
| `packages/domain/auth/pages/login/hooks/useLoginSubmission/types/useLoginSubmission.ts` | Move peer hook under owning useLoginForm hook or combine into it | Structural finding |
| `packages/domain/auth/pages/login/hooks/useLoginSubmission/useLoginSubmission.ts` | **Source confirmed:** uses `APP_ROUTE.DASHBOARD`, `_keys`, flat helpers; peer submission hook still outside one owner | Source reviewed |
| `packages/domain/auth/pages/login/hooks/useLoginViewport/helpers/getMobileSnapshot/getMobileSnapshot.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/auth/pages/login/hooks/useLoginViewport/helpers/getMobileSnapshot/index.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/auth/pages/login/hooks/useLoginViewport/helpers/getMobileSnapshot/types/getMobileSnapshot.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/auth/pages/login/hooks/useLoginViewport/helpers/subscribeToViewport/index.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/auth/pages/login/hooks/useLoginViewport/helpers/subscribeToViewport/subscribeToViewport.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/auth/pages/login/hooks/useLoginViewport/helpers/subscribeToViewport/types/subscribeToViewport.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/auth/pages/login/hooks/useLoginViewport/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/login/hooks/useLoginViewport/types/useLoginViewport.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/login/hooks/useLoginViewport/useLoginViewport.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/domain/auth/pages/login/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/login/templates/desktopLogin/DesktopLogin.tsx` | **Source confirmed:** central `AUTH_ROUTE.FORGOT_PASSWORD` fixed; still native Link and raw global class | Source reviewed |
| `packages/domain/auth/pages/login/templates/desktopLogin/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/login/templates/desktopLogin/types/desktopLogin.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/login/templates/mobileLogin/MobileLogin.tsx` | **Source confirmed:** central `AUTH_ROUTE.FORGOT_PASSWORD` fixed; still native Link and raw global class | Source reviewed |
| `packages/domain/auth/pages/login/templates/mobileLogin/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/login/templates/mobileLogin/types/mobileLogin.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/login/types/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/login/types/login.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/login/types/login.state.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/recovery/Recovery.tsx` | **Confirmed:** raw status strings, title and JSX ternaries, inline class | Source reviewed |
| `packages/domain/auth/pages/recovery/__tests__/RecoveryPage.test.tsx` | Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/domain/auth/pages/recovery/components/recoveryArtwork/RecoveryArtwork.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/auth/pages/recovery/components/recoveryArtwork/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/recovery/components/recoveryArtwork/types/recoveryArtwork.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/recovery/components/recoveryPasswordForm/RecoveryPasswordForm.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/auth/pages/recovery/components/recoveryPasswordForm/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/recovery/components/recoveryPasswordForm/types/recoveryPasswordForm.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/recovery/components/recoveryState/RecoveryState.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/auth/pages/recovery/components/recoveryState/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/recovery/components/recoveryState/types/recoveryState.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/recovery/containers/recoveryContainer/RecoveryContainer.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/auth/pages/recovery/containers/recoveryContainer/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/recovery/containers/recoveryContainer/types/recoveryContainer.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/recovery/helpers/classifyRecoveryFailure/__tests__/classifyRecoveryFailure.test.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts}; Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/domain/auth/pages/recovery/helpers/classifyRecoveryFailure/classifyRecoveryFailure.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/auth/pages/recovery/helpers/classifyRecoveryFailure/index.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/auth/pages/recovery/helpers/classifyRecoveryFailure/types/classifyRecoveryFailure.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/auth/pages/recovery/helpers/getRecoveryArtworkCopy/__tests__/getRecoveryArtworkCopy.test.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts}; Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/domain/auth/pages/recovery/helpers/getRecoveryArtworkCopy/getRecoveryArtworkCopy.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/auth/pages/recovery/helpers/getRecoveryArtworkCopy/index.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/auth/pages/recovery/helpers/getRecoveryArtworkCopy/types/getRecoveryArtworkCopy.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/auth/pages/recovery/helpers/validateRecoveryPassword/__tests__/validateRecoveryPassword.test.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts}; Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/domain/auth/pages/recovery/helpers/validateRecoveryPassword/index.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/auth/pages/recovery/helpers/validateRecoveryPassword/types/validateRecoveryPassword.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/auth/pages/recovery/helpers/validateRecoveryPassword/validateRecoveryPassword.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/auth/pages/recovery/hooks/useRecovery/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/recovery/hooks/useRecovery/types/useRecovery.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/recovery/hooks/useRecovery/useRecovery.ts` | **Confirmed:** long combined form/proof/grant/submission/retry hook; split with one state owner | Source reviewed |
| `packages/domain/auth/pages/recovery/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/recovery/types/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/auth/pages/recovery/types/recovery.artwork.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/auth/pages/recovery/types/recovery.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/dashboard/RUNBOOK.md` | Check documentation examples/references against current source | Documentation review |
| `packages/domain/dashboard/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/dashboard/pages/dashboard/Dashboard.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/dashboard/pages/dashboard/components/dashboardAside/DashboardAside.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/dashboard/pages/dashboard/components/dashboardAside/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/dashboard/pages/dashboard/components/dashboardAside/types/dashboardAside.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/dashboard/pages/dashboard/components/dashboardHeader/DashboardHeader.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/dashboard/pages/dashboard/components/dashboardHeader/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/dashboard/pages/dashboard/components/dashboardHeader/types/dashboardHeader.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/dashboard/pages/dashboard/components/dashboardMain/DashboardMain.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/dashboard/pages/dashboard/components/dashboardMain/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/dashboard/pages/dashboard/components/dashboardMain/types/dashboardMain.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/dashboard/pages/dashboard/components/dashboardNotice/DashboardNotice.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/dashboard/pages/dashboard/components/dashboardNotice/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/dashboard/pages/dashboard/components/dashboardNotice/types/dashboardNotice.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/dashboard/pages/dashboard/components/dashboardOverview/DashboardOverview.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/dashboard/pages/dashboard/components/dashboardOverview/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/dashboard/pages/dashboard/components/dashboardOverview/types/dashboardOverview.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/dashboard/pages/dashboard/constants/dashboard.metrics.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/domain/dashboard/pages/dashboard/constants/dashboard.workspaces.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/domain/dashboard/pages/dashboard/containers/dashboardContainer/DashboardContainer.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/dashboard/pages/dashboard/containers/dashboardContainer/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/dashboard/pages/dashboard/containers/dashboardContainer/types/dashboardContainer.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/dashboard/pages/dashboard/helpers/__tests__/dashboard.getDateParts.test.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts}; Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/domain/dashboard/pages/dashboard/helpers/dashboard.getDateParts.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/domain/dashboard/pages/dashboard/helpers/types/dashboard.getDateParts.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/dashboard/pages/dashboard/hooks/useDashboardCalendar/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/dashboard/pages/dashboard/hooks/useDashboardCalendar/types/useDashboardCalendar.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/dashboard/pages/dashboard/hooks/useDashboardCalendar/useDashboardCalendar.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/domain/dashboard/pages/dashboard/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/dashboard/pages/dashboard/types/dashboard.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/dashboard/pages/dashboard/types/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/integrations/index.ts` | Remove unused UI domain after route/importer check | Structural finding |
| `packages/domain/integrations/pages/ConnectionDetailsPage.tsx` | Remove unused UI domain after route/importer check | Structural finding |
| `packages/domain/integrations/pages/IntegrationsPage.tsx` | Remove unused UI domain after route/importer check | Structural finding |
| `packages/domain/people/index.ts` | Remove unused UI domain after route/importer check | Structural finding |
| `packages/domain/people/pages/PeoplePage.tsx` | Remove unused UI domain after route/importer check | Structural finding |
| `packages/domain/publicSite/RUNBOOK.md` | Check documentation examples/references against current source | Documentation review |
| `packages/domain/publicSite/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/publicSite/pages/landing/LandingPage.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/publicSite/pages/landing/__tests__/LandingPage.test.tsx` | Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/domain/publicSite/pages/landing/components/aboutContent/aboutContent.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/publicSite/pages/landing/components/aboutContent/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/publicSite/pages/landing/components/campusContent/campusContent.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/publicSite/pages/landing/components/campusContent/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/publicSite/pages/landing/components/enquiryForm/constants/enquiry.initial.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/domain/publicSite/pages/landing/components/enquiryForm/containers/enquiryFormContainer/EnquiryFormContainer.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/publicSite/pages/landing/components/enquiryForm/containers/enquiryFormContainer/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/publicSite/pages/landing/components/enquiryForm/enquiryForm.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/publicSite/pages/landing/components/enquiryForm/helpers/getEnquiryMailto/__tests__/getEnquiryMailto.test.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts}; Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/domain/publicSite/pages/landing/components/enquiryForm/helpers/getEnquiryMailto/getEnquiryMailto.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/publicSite/pages/landing/components/enquiryForm/helpers/getEnquiryMailto/index.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/publicSite/pages/landing/components/enquiryForm/helpers/getEnquiryMailto/types/getEnquiryMailto.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/publicSite/pages/landing/components/enquiryForm/helpers/validateEnquiryContact/__tests__/validateEnquiryContact.test.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts}; Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/domain/publicSite/pages/landing/components/enquiryForm/helpers/validateEnquiryContact/index.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/publicSite/pages/landing/components/enquiryForm/helpers/validateEnquiryContact/types/validateEnquiryContact.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/publicSite/pages/landing/components/enquiryForm/helpers/validateEnquiryContact/validateEnquiryContact.ts` | Nested helper layout; flatten into helpers/{types,__tests__,<owner>.<useCase>.ts} | Structural finding |
| `packages/domain/publicSite/pages/landing/components/enquiryForm/hooks/useEnquiryForm/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/publicSite/pages/landing/components/enquiryForm/hooks/useEnquiryForm/types/useEnquiryForm.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/publicSite/pages/landing/components/enquiryForm/hooks/useEnquiryForm/useEnquiryForm.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/domain/publicSite/pages/landing/components/enquiryForm/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/publicSite/pages/landing/components/enquiryForm/types/enquiryForm.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/publicSite/pages/landing/components/homeContent/homeContent.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/publicSite/pages/landing/components/homeContent/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/publicSite/pages/landing/components/publicSheet/PublicSheet.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/publicSite/pages/landing/components/publicSheet/components/infoPanels/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/publicSite/pages/landing/components/publicSheet/components/infoPanels/infoPanels.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/publicSite/pages/landing/components/publicSheet/components/panelPhoto/PanelPhoto.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/publicSite/pages/landing/components/publicSheet/components/panelPhoto/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/publicSite/pages/landing/components/publicSheet/components/storyPanel/StoryPanel.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/publicSite/pages/landing/components/publicSheet/components/storyPanel/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/publicSite/pages/landing/components/publicSheet/constants/school.sheetCopy.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/domain/publicSite/pages/landing/components/publicSheet/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/publicSite/pages/landing/components/publicSheet/types/schoolPanel.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/publicSite/pages/landing/components/resourceContent/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/publicSite/pages/landing/components/resourceContent/resourceContent.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/publicSite/pages/landing/components/siteActions/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/publicSite/pages/landing/components/siteActions/siteActions.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/publicSite/pages/landing/components/siteChrome/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/publicSite/pages/landing/components/siteChrome/siteChrome.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/publicSite/pages/landing/components/siteMedia/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/publicSite/pages/landing/components/siteMedia/siteMedia.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/publicSite/pages/landing/constants/school.content.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/domain/publicSite/pages/landing/constants/school.navigation.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/domain/publicSite/pages/landing/constants/school.panels.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/domain/publicSite/pages/landing/constants/school.seo.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/domain/publicSite/pages/landing/containers/landingContainer/LandingContainer.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/publicSite/pages/landing/containers/landingContainer/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/publicSite/pages/landing/hooks/useLandingPage/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/publicSite/pages/landing/hooks/useLandingPage/types/useLandingPage.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/domain/publicSite/pages/landing/hooks/useLandingPage/useLandingPage.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/domain/publicSite/pages/landing/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/domain/publicSite/pages/landing/registry/schoolPages.registry.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/domain/publicSite/pages/landing/styles/landingPage.css` | Large global feature CSS; tokenize and split scoped ownership | Source reviewed |
| `packages/domain/publicSite/pages/landing/styles/optionOne.css` | Legacy non-module CSS; map selectors to owning .module.scss with @Styles tokens | Structural finding |
| `packages/domain/publicSite/pages/landing/styles/optionOnePages.css` | Legacy non-module CSS; map selectors to owning .module.scss with @Styles tokens | Structural finding |
| `packages/domain/settings/index.ts` | Remove unused UI domain after route/importer check | Structural finding |
| `packages/domain/settings/pages/SecurityPage.tsx` | Remove unused UI domain after route/importer check | Structural finding |
| `packages/shell/RUNBOOK.md` | Check documentation examples/references against current source | Documentation review |
| `packages/shell/Shell.tsx` | **Source confirmed:** native layout/button; sidebarCollapsed ternary builds BEM class; use Box/Button and named class helper | Source reviewed |
| `packages/shell/components/appRoutes/AppRoutes.tsx` | **Confirmed:** direct .map route JSX, mixed route policy; central route constants and renderer | Source reviewed |
| `packages/shell/components/appRoutes/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/shell/components/appRoutes/types/appRoutes.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/shell/components/mobileNav/MobileNav.tsx` | **Source confirmed:** direct Lucide, literal `/dashboard`, conditional class ternary and native nav/button | Source reviewed |
| `packages/shell/components/mobileNav/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/shell/components/mobileNav/types/mobileNav.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/shell/components/pinnedNavigation/PinnedApp.tsx` | **Source confirmed:** direct Lucide, NavLink class ternary, raw button/div/span; semantic components | Source reviewed |
| `packages/shell/components/pinnedNavigation/PinnedNavigation.tsx` | **Source confirmed:** native flatMap/find/map, entries/empty ternary, native nav/p; move to Base and renderer | Source reviewed |
| `packages/shell/components/pinnedNavigation/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/shell/components/pinnedNavigation/types/pinnedNavigation.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/shell/components/protectedLayout/ProtectedLayout.tsx` | Review native JSX, ternaries, icons, test IDs, route literals and scoped styles | Content review required |
| `packages/shell/components/protectedLayout/__tests__/ProtectedLayout.test.tsx` | Check exact exported-name describe; no nested describe/pure mocks | Structural finding |
| `packages/shell/components/protectedLayout/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/shell/components/protectedLayout/types/protectedLayout.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/shell/components/sidebar/Sidebar.tsx` | **Source confirmed:** direct Lucide, native buttons/aside, raw routes and class/label ternaries | Source reviewed |
| `packages/shell/components/sidebar/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/shell/components/sidebar/types/sidebar.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/shell/components/themeChoices/ThemeChoices.tsx` | **Source confirmed:** inline theme option literals and selection ternary, native button, direct icons; constants + renderer | Source reviewed |
| `packages/shell/components/themeChoices/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/shell/components/themeChoices/types/themeChoices.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/shell/components/topbar/Topbar.tsx` | **Source confirmed:** direct Lucide, native header/button/span markup and global BEM | Source reviewed |
| `packages/shell/components/topbar/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/shell/components/topbar/types/topbar.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/shell/containers/shellContainer/ShellContainer.tsx` | **Confirmed:** QueryClient construction/options mixed; separate constant and factory | Source reviewed |
| `packages/shell/containers/shellContainer/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/shell/containers/shellContainer/types/shellContainer.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/shell/hooks/useShell/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/shell/hooks/useShell/types/useShell.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/shell/hooks/useShell/useShell.ts` | **Confirmed:** peer useShellState + useShellActions, consolidate | Source reviewed |
| `packages/shell/hooks/useShellActions/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/shell/hooks/useShellActions/types/useShellActions.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/shell/hooks/useShellActions/useShellActions.ts` | **Confirmed:** duplicated boolean setters; native includes/filter; writes persistence in state updater | Source reviewed |
| `packages/shell/hooks/useShellState/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/shell/hooks/useShellState/types/useShellState.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/shell/hooks/useShellState/useShellState.ts` | **Source confirmed:** independent peer state owner; fold into authoritative useShell contract | Source reviewed |
| `packages/shell/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/shell/layout/AppLauncher.tsx` | **Confirmed:** monolithic catalog, native button/input, inline filter and raw class/route/icon | Source reviewed |
| `packages/shell/registry/appCatalog.registry.ts` | **Source confirmed:** direct Lucide import, literal `/dashboard` duplicated and hard-coded pin descriptors; use shared route/icon constants | Source reviewed |
| `packages/shell/registry/domains.registry.tsx` | **Confirmed:** inline route strings, raw purpose and legacy redirect; business path contracts | Source reviewed |
| `packages/shell/registry/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/shell/storage/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/shell/storage/types/workspacePreferences.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/shell/storage/workspacePreferences.ts` | **Source confirmed:** inline storage keys, native filter/map and Set policy; constant keys and helper, preserve SSR/privacy handling | Source reviewed |
| `packages/shell/types/shell.props.ts` | Type contract: verify enum keys, props, boundary ownership and consumer compatibility | Structure only |
| `packages/styles/RUNBOOK.md` | Check documentation examples/references against current source | Documentation review |
| `packages/styles/base.css` | Review token colors, theme parity, media breakpoints and global selector ownership | Content review required |
| `packages/styles/index.css` | Review token colors, theme parity, media breakpoints and global selector ownership | Content review required |
| `packages/styles/index.ts` | Public API/entrypoint: verify exports used; no path-level violation detected | Structure only |
| `packages/styles/themes/dark.css` | Review token colors, theme parity, media breakpoints and global selector ownership | Content review required |
| `packages/styles/themes/light.css` | Review token colors, theme parity, media breakpoints and global selector ownership | Content review required |
| `packages/styles/tokens/breakpoints.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |
| `packages/styles/tokens/foundations.css` | Review token colors, theme parity, media breakpoints and global selector ownership | Content review required |
| `packages/styles/utilities.css` | Review token colors, theme parity, media breakpoints and global selector ownership | Content review required |
| `pnpm-lock.yaml` | Generated dependency snapshot; validate with frozen pnpm install, not manual code refactoring | Tracked metadata |
| `scripts/check-boundaries.mjs` | Verify configuration, import aliases, scripts, build/Jest/ESLint and boundary enforcement | Config review required |
| `scripts/check-layout-assets.mjs` | **Source confirmed:** resolves route source by regex; catalog `/dashboard` check breaks with route constants instead of literal paths (latest failing CI) | Source reviewed |
| `scripts/check-styles.mjs` | Verify configuration, import aliases, scripts, build/Jest/ESLint and boundary enforcement | Config review required |
| `scripts/check-ui-contracts.mjs` | Verify configuration, import aliases, scripts, build/Jest/ESLint and boundary enforcement | Config review required |
| `scripts/check-units.mjs` | **Source confirmed:** enforces nested helper index.ts/contracts for Base and some page helpers, contradicting Rev 8 flat-helper rule; rewrite checker by ownership | Source reviewed |
| `scripts/download-school-images.mjs` | Verify configuration, import aliases, scripts, build/Jest/ESLint and boundary enforcement | Config review required |
| `scripts/prerender-public-pages.mjs` | Verify configuration, import aliases, scripts, build/Jest/ESLint and boundary enforcement | Config review required |
| `scripts/test/polyfills.cjs` | Verify configuration, import aliases, scripts, build/Jest/ESLint and boundary enforcement | Config review required |
| `scripts/test/styleMock.cjs` | Verify configuration, import aliases, scripts, build/Jest/ESLint and boundary enforcement | Config review required |
| `scripts/test/transformEsm.cjs` | Verify configuration, import aliases, scripts, build/Jest/ESLint and boundary enforcement | Config review required |
| `tsconfig.json` | Verify configuration, import aliases, scripts, build/Jest/ESLint and boundary enforcement | Config review required |
| `vite.config.ts` | Review raw strings, native loops, API/IO boundaries, exact helpers and tests | Content review required |

## 41. Cross-cutting confirmed blockers and review priorities

**Current priority evidence:**

1. `pages/login/components/LoginForm/LoginForm.tsx` still refers to `LOGIN_FIELDS` despite importing `LOGIN_FIELD_CONFIG`; this is an unresolved symbol and should be found by TypeScript once the architecture gate passes. The field config already exists; fix both field descriptors without altering autofill.
2. New Base `_find`/`_keys`/`_map`/`_trim` types were missing in commit `200d35de`, but the subsequent `b8340847` commit added contracts. Do **not** keep reporting that older failure as unresolved without checking the latest head.
3. `scripts/check-layout-assets.mjs` currently compares string catalog paths against route registration that now uses `APP_ROUTE.DASHBOARD`. Its static regex does not correctly resolve that identifier, producing a false `catalog path /dashboard has no implemented protected route` architecture failure. Use a shared typed descriptor/route manifest rather than source-code regex parsing.
4. `scripts/check-units.mjs` still requires nested helper folders, their `index.ts` and callable `types`; Revision 8 explicitly requires flat page helpers and no helper barrels. The checker must distinguish Base leaf utility units from flat domain page helpers instead of enforcing opposite conventions.
5. Login helpers and exact suite names are now mostly migrated; login presentation, auth provider, recovery hooks, shell JSX, app launcher and legacy CSS remain substantial nonconformant areas. Do not mark a package complete when only constants or tests were reorganized.
6. Public site/dashboard CSS has large legacy/global files; style migration requires CSS selector-to-DOM ownership mapping and real dark/light/mobile browser checks. Deleting CSS by regex or mass rename is unsafe.
7. `domain/accounts`, `domain/people`, `domain/integrations`, `domain/settings` remain unimplemented UI-domain deletion candidates. Check all static/dynamic dependencies and exports before deleting, while preserving any real business/auth services.

**Audit method limitation:** Every path is inventoried, but only rows explicitly marked **Source reviewed** were read/verified at the code-content level in this audit. The other entries are obligations, not a completed manual line-by-line review. Full AST no-ternary, import ownership, pure-test mocking, dead-symbol detection, JSX primitive enforcement and SCSS module validations must be implemented as automated repository checks before claiming exhaustive compliance.

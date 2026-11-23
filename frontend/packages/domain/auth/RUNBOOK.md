# Identity UI domain

Owns login, forgot password, invitation activation, password reset and session-aware route guards.
Canonical path: `frontend/packages/domain/auth/`. The shell imports only `@Domain/auth`; business transport remains in `@Business/auth` and `@Business/accounts`.

## Page structure

- `pages/login/Login.tsx` is presentational, while `pages/login/containers/loginContainer/` binds state and screen size.
- `hooks/useLogin/` composes `useLoginForm`, `useLoginActions` and `useLoginSubmission`; it does not own one giant stateful use case.
- `pages/login/helpers/` contains independent pure validation, email normalization and safe error mapping, each with a colocated test and callable type.
- The desktop and mobile login templates render **one** form in the DOM; forgot-password placement follows the established desktop/mobile layout.
- `pages/forgotPassword/` separates the countdown hook, request orchestration, pure email validation, copy selection and form presentation.
- `pages/recovery/` shares invitation and reset presentation, pure password validation and recovery-state artwork copy.
- The `LoginV2Page`, `ForgotPasswordPage` and `RecoveryPage` names remain compatibility exports without duplicate implementations.

## Security and attempts

- Sessions rely on HttpOnly cookies; no auth tokens or passwords are stored in localStorage or sessionStorage.
- One-time invitation/reset bearer proofs are held **in memory** during exchange, then removed from URL fragments. Only a non-secret marker of a successfully exchanged grant cookie may persist in sessionStorage.
- The server alone determines when activation attempts have been exhausted, using `IDENTITY_SETUP_ATTEMPTS_EXCEEDED`. A transient client/network error must never produce a fabricated local attempt count.
- Local validation rejects malformed forms before contacting the backend. Route guard session bootstrap happens only when a protected route is opened.

## Checks

From `frontend/`: `pnpm check:architecture`, `pnpm typecheck`, `pnpm test -- --runInBand`, `pnpm build`. Browser verification remains required for keyboard focus, narrow-height mobile screens and reset/activation token behavior.

## Layout and art ownership (Revision 6)

The active `LoginLayout` now lives under `components/loginLayout/`, with its scoped stylesheet. `AuthAmbientGraphic` and `AuthStatusArt` are shared domain components, and the academic-pathways vector is supplied by `@Assets/auth` through the reusable `@Components/components/SVGRenderer`. Deprecated `AuthFrame`, mobile-school SVG mockups, the standalone `LoginV2` prototype and the Soffit shader were deleted as unused.

`hooks/useAuthViewportFit/` is domain-wide; `core/guards/RecoveryFlowBoundary` and `core/session/recoveryFlowSession` own the route/session proof boundary. The `/loginv2` URL remains a redirect to `/login`, not a separate screen. The remaining login/recovery helper flattening and recovery hook decomposition are tracked in the unified engineering guideline §22.

## Global copy, status enums and presentation

All visible authentication copy is resolved while rendering through bare global `__()`, without `import { __ }` in pages, templates, components or renderers. Stable English source messages and fallbacks remain in the nearest `constants/<owner>.copy.ts`. Dynamic names, counts and addresses use named `{{value}}` parameters. Backend auth error codes, proof tokens, cookies, route paths, input types and `dataTestId` values remain technical and must not be translated. The shell initializes the resolver **before** public login, forgot password, set-password and recovery routes mount.

Use the component's exported typed values for closed properties (`LoginLayoutVariant.LOGIN`, `BUTTON_SIZE.LARGE`, `BOX_TAG.H3`, etc.). `LoginLayout` delegates school branding and trust/footer presentation to dedicated `components/visual` and `components/footer` owners to stay below the 155-line SRP architecture limit. Renderer functions are pure and have owner-local tests. Preserve a11y labels, one-form responsive behavior, token-visit safety, 30-second resend cooldown, and identical visible English copy.

## Literal-source auth copy and pure renderer choices

`LOGIN_COPY.desktopTitle = __('Welcome back.')` and `<LoginLayout title={LOGIN_COPY.desktopTitle} />` are required. The pattern `__(LOGIN_COPY.desktopTitle)` is forbidden. Server/auth error identifiers remain technical values. `pages/login/renderer/login.formTemplate.ts` returns DesktopLogin or MobileLogin via explicit branch returns, preserving a single mounted form without `let FormTemplate`. Forgot-password countdown/email messages use `__('Resend in {{seconds}}s', { seconds })` and a pure selector. Recovery password, state, and artwork labels follow the same rule. See master guideline Revision 14 and CI `check:no-let`.

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

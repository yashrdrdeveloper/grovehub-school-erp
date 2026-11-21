# Auth feature
Owns sign-in, forgot password, invitation setup and session runtime. `@Business/auth` owns transport, `packages/auth/pages` owns UI, and the shell only imports the public `@Domain/auth` facade. Session secrets are HttpOnly cookies; do not store them in localStorage, sessionStorage or Redux.

## Login unit — container / presentation

The login unit lives in `pages/login/` and is exported through its `index.ts`. `Login.tsx` is presentational; `LoginContainer.tsx` is the route entry point.
The legacy `LoginV2Page` name is maintained as a compatibility export.

- `Login.tsx` is a pure page composition that selects one active desktop or mobile template from supplied props. `LoginContainer.tsx` invokes the state and breakpoint hooks and binds them to `Login`.
- `hooks/useLogin/useLogin.tsx` owns the controlled fields, client validation, password visibility, pending/submission guard, existing `useAuth().signIn` call and successful redirect.
- `hooks/useLoginViewport/` subscribes to viewport changes; it renders **one** form rather than duplicate hidden forms (autofill and IDs remain unique).
- `helpers/login.credentials.validation.ts` and `helpers/login.authentication.error.ts` are pure and depend on `config/login.config.ts` for field rules, public messages and paths.
- `components/LoginForm/` is stateless and renders the shared accessible fields, errors, and submit button.
- `templates/desktopLogin/` places Forgot password beside the password label; `templates/mobileLogin/` places it below the inputs.
- `types/login.props.ts` and `types/login.state.ts` are the local contracts. Backend transport remains in existing auth business/provider modules.

**Boundary:** Templates and `LoginForm` cannot import auth providers, call services, or navigate. Do not add password or session storage in this unit.

**Run checks:** From `frontend/`, use `pnpm typecheck`, `pnpm test -- --runInBand`, `pnpm lint`, and `pnpm check:architecture`. Also verify the 720px breakpoint and dark/light mobile appearances in a real browser.

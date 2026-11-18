# GrooveHub frontend

React 19 + TypeScript + Vite + pnpm. The application uses the repository's mandatory `docs/engineering/frontend/FRONTEND_CODE_GUIDELINES.md`.

## Commands

From the monorepo root:

```bash
make frontend-setup
make frontend-start
make frontend-stop
make frontend-status
make frontend-build
make dev-start
```

The frontend runs at **http://127.0.0.1:4173**. Vite proxies `/identity`, `/foundation`, `/integrations`, and `/health` to Launcher on `127.0.0.1:8080`, so the browser uses same-origin cookie authentication.

## Modules

- `packages/base`: pure technical utilities. Use `@Base/utilities/isEqual` for the Lodash shim, never import Lodash directly in domain components.
- `packages/business`: typed API clients and business contracts shared by domains.
- `packages/components`: reusable presentational components, no API calls.
- `packages/auth/pages`: login, forgot/reset password, invitation setup.
- `packages/domain`: public domain facade and domain-owned screens.
- `packages/shell`: provider and route composition, responsive workspace frame, no business logic.

## Security and actual backend capabilities

The backend sets opaque HttpOnly access/refresh cookies and a readable CSRF cookie. The frontend uses `credentials: 'include'`, sends `X-CSRF-Token` for mutations after a cookie is established, and never stores access/refresh tokens in browser storage.

The frontend wires all 22 current Identity operations through typed clients. Accounts and sessions have live screens; Foundation Person and Integration endpoints have typed API clients. Student lifecycle, finance, tasks and calendar are not yet implemented backend modules; their dashboard fields deliberately show unavailable placeholders rather than fake metrics.

Invitations/recovery use one-time proof fragments: the UI removes `#proof=...` from browser history before exchanging it for a short-lived, HttpOnly grant cookie.

Production rollout requires a configured trusted origin, secure cookies over HTTPS, reverse-proxy settings, enabled Identity accounts, and a working Integration outbound email connection. The `frontend/packages/business/api` layer is the single place to adapt future API schema changes.

## Design

Login: tenant inferred from `*.prod.groovehub.com`, no organization-code field, compact GrooveHub branding and a single credential form. Workspace: understated green/off-white, small chrome, navigation sidebar on desktop and bottom bar on mobile, user information inside Quick Settings rather than the global header.

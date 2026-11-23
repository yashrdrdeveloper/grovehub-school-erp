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
- `packages/domain/auth/pages`: login, forgot/reset password and invitation setup.
- `packages/domain`: public domain facade and domain-owned screens.
- `packages/shell`: provider and route composition, responsive workspace frame, no business logic.

## Security and actual backend capabilities

The backend sets opaque HttpOnly access/refresh cookies and a readable CSRF cookie. The frontend uses `credentials: 'include'`, sends `X-CSRF-Token` for mutations after a cookie is established, and never stores access/refresh tokens in browser storage.

The frontend wires all 22 current Identity operations through typed clients. Accounts and sessions have live screens; Foundation Person and Integration endpoints have typed API clients. Student lifecycle, finance, tasks and calendar are not yet implemented backend modules; their dashboard fields deliberately show unavailable placeholders rather than fake metrics.

Invitations/recovery use one-time proof fragments: the UI removes `#proof=...` from browser history before exchanging it for a short-lived, HttpOnly grant cookie.

Production rollout requires a configured trusted origin, secure cookies over HTTPS, reverse-proxy settings, enabled Identity accounts, and a working Integration outbound email connection. The `frontend/packages/business/api` layer is the single place to adapt future API schema changes.

## Design

Login and recovery use shared RDR School branding, one credential form and system appearance on phones. Workspace uses RDR navy tokens, a collapsible desktop sidebar, and compact phone navigation.

Dependencies are pinned in `pnpm-lock.yaml`; setup and CI use `pnpm install --frozen-lockfile`.

## Global message API and component contracts

The app uses a single browser-global `__()` for user-visible copy. `ShellContainer` calls its owned `useInitI18n` hook and **blocks rendering of both public and protected routes** until `window.__` is installed. Do not import `__` from `@Base/i18n` inside features. TypeScript exposes the function through `app/types/i18n.global.d.ts`; Jest installs the real pure formatter from `scripts/test/setupI18n.ts`.

```tsx
import { Box, BOX_TAG } from '@Components/components/Box';

export const Greeting = ({ name }: { name: string }) => (
  <Box as={BOX_TAG.P}>{__('Hello {{name}}', { name })}</Box>
);
```

`@Base/i18n` is currently an **English-only interpolation utility**; using `__()` is not equivalent to shipping multilingual dictionaries. Never format API identifiers, route paths, CSS classes, test IDs or HTML/ARIA protocol tokens. Design-system closed prop values use component-owned exports (`BUTTON_VARIANT`, `BUTTON_SIZE`, `BOX_TAG`, `TEXT_TAG`, etc.), not arbitrary strings.

From `frontend/`, verify changes with `pnpm check:architecture`, `pnpm typecheck`, `pnpm exec jest --runInBand` and `pnpm build`. `pnpm check:architecture` includes a strict global i18n boundary checker and a measured zero-new-raw-JSX-copy baseline.

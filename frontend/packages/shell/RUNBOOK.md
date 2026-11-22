# Shell — application composition

The shell owns providers, authenticated workspace chrome, route composition, view templates, and the browser-only workspace preferences. It does **not** own Identity transport endpoints, backend APIs or route path definitions. Use business route descriptors from `@Business/auth/routes` and `@Business/app/routes`; import domain pages from their root APIs only.

## Canonical owner layout

- `Shell.tsx`, `index.ts`, `types/`: tiny public presentation root.
- `containers/shellContainer/`: BrowserRouter, query client and authentication provider composition.
- `containers/workspaceShellContainer/`: one persistent `useShell` state owner.
- `templates/workspaceShell/`: pure full workspace composition; nested `Outlet` persists across child navigation.
- `components/topbar/`, `sidebar/`, `mobileNav/`, `appLauncher/`, `themeChoices/`, `pinnedNavigation/`: one accessible visual responsibility per component.
- `hooks/useShell/`: state/actions; `helpers/shell.preferences.ts` does best-effort localStorage; `constants/shell.storageKeys.ts` owns keys.
- `constants/shell.catalog.ts`: only implemented protected applications, shared by the sidebar and launcher.
- `renderer/shell.routes.tsx`: maps business route IDs to view elements without duplicating URLs.
- `helpers/` and `constants/` are flat. No `registry/`, `storage/`, `layout/` catch-all or helper barrel.

Navigation links and routing primitives must use `@Components/router`; generic DOM is rendered through `@Components` leaves. `Dialog` subcomponents are exposed from `@Components/components/Dialog`, never from `Dialog/parts`. Shared route error boundaries prevent one failed page from taking down the entire app.

## RDR School workspace behavior

The existing demo label is **RDR School**; future tenant-aware labels must come from domain data. A persistent sidebar supports collapse and pin/unpin actions; phones show Home, Apps and Menu with accessible dismissal. Catalog search lists **implemented apps only**—currently Dashboard. Pinned paths and collapse state remain local to the browser, not synchronized server-side.

Theme preference uses `ThemePreference.SYSTEM`, `LIGHT`, `DARK` from `@Components/hooks/useWorkspaceTheme`. Mobile follows the operating-system appearance without overwriting the desktop preference. Frontend bootstrap resolves the theme before render.

The Dashboard presents unavailable/empty states without invented school analytics. Public school pages are **not** wrapped in workspace chrome. Public authentication routes use `PublicOnly`; protected pages use `RequireSession → ProtectedLayout → WorkspaceShellContainer → Shell → WorkspaceShell → Outlet`.

## Checks and outstanding migration

Run `pnpm check:architecture`, `pnpm typecheck`, `pnpm exec jest --runInBand`, and `pnpm build`. Architectural checks also block misplaced business transport contracts, native shell markup, direct router imports, and new audit findings. The remaining legacy selector/CSS ownership debt is still tracked; a successful build alone does not verify 320–720 px behavior, keyboard focus, dark mode contrast or browser autofill.

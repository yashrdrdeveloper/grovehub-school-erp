# Shell
The shell is a composition-only package. It constructs React Router, TanStack Query, auth provider and the workspace chrome. Feature screens are imported from `@Domain/<name>` public APIs. No shell file imports `@Business`, calls `fetch`, or performs a use case.


## RDR School experience (November 2026)

The workspace displays **RDR School**; this is the requested school demo identity, not a claim that dynamic multi-tenant school branding exists. Later tenant-aware school metadata should replace the static label in `components/protectedLayout/ProtectedLayout.tsx` without changing the shared shell.

### Visual system

Use screenshot-inspired neutral colors: clean white and pale-gray surfaces for light mode; near-black indigo-charcoal surfaces with hairline dividers for dark mode. Use bright pink, blue, and green sparingly as information accents, not large page backgrounds. Avoid oversized colorful hero gradients, fake charts, or fabricated metrics. Shared global design tokens are in `packages/styles`; component recipes belong in `packages/components/styles`. Follow components-first composition, accessible focus states, and clear contrast.

### Appearance

`@Components/hooks/useWorkspaceTheme/` controls `system | light | dark` and persists explicit preference using `groovehub.theme` in browser local storage. The default is `system`; OS preference changes are observed. `app/main.tsx` resolves the preferred theme before rendering including public authentication pages. Theme controls are available in account preferences.

### Navigation and application catalog

The shared shell owns the narrow desktop sidebar, mobile navigation (Home, Apps, Menu) and searchable, categorized app catalog. Active routes come from `registry/domains.registry.tsx`, and available shortcuts come from `registry/navigation.registry.ts`. The catalog shows **implemented routes only**. Placeholder People, Accounts, Integrations, Security, Settings and future modules have been removed until working routes and business logic exist. Prioritize lightweight touch targets and a mobile bottom sheet for the catalog.

### Operational dashboard

The Dashboard domain shows four compact metrics and six at-a-glance cards: Pending for You, Admissions & Dues, Today & Upcoming, Reminders, Complaints & Requests, Audit & Compliance. Cards are responsive, compact, and ready for future domain APIs. Show explicit unavailable states for missing data; never display arbitrary pretend counts or fabricate school activity. As services are connected, their domain hooks/queries should replace static empty states, not fetch directly from the shell.

### Mobile usability

At smaller widths, sidebar becomes a dismissible drawer, app catalog a bottom sheet, KPI cards become two columns and operational panels become a single column on phones. Keep top-level tasks discoverable without horizontal page scrolling; individual data tables can scroll within their own containers.


## Workspace navigation and app catalog

- The top-left four-square button opens a **full-workspace catalog**. On desktop it fills the area to the right of the persistent sidebar; on mobile it uses the entire viewport. Escape or the close button dismisses it.
- **Pinned apps** are shown as the first catalog section and populate the sidebar in saved order. Use the pin/unpin control on any available application tile or hover over a pinned sidebar item to unpin it.
- Unimplemented applications are **not listed**. Introduce an application only after it has an implemented, authorized route.
- The desktop sidebar collapses to an icon rail. Its expanded/collapsed preference and pinned routes are stored in local browser storage (`groovehub.workspace.sidebar.collapsed.v1` and `groovehub.workspace.pinned.v1`).
- The mobile sidebar remains a drawer even if the desktop collapse preference is on. The mobile Apps entry opens the full-screen catalog.
- Navigation and catalog application metadata live in `registry/appCatalog.registry.ts`. The catalog UI is in `layout/AppLauncher.tsx`. Shell reads the same catalog, so app labels/icons/links must remain consistent. Pinning is currently per-browser, not server-synchronized.
- Appearance remains **System** by default, with user-selected Light and Dark modes. The visual palette follows the November 18 clean-analytics screenshots: minimal white light mode, near-black dark mode, hairline borders, and restrained vivid status/chart accents. Light and dark share one layout; the visual theme follows System by default.


### Reference-driven theme (November 18)

Dashboard remains a layout preview, never invented analytics. KPI cards form a compact hairline-divided row on desktop; sections stay readable on mobile. The mobile shell uses a compact floating bottom navigation; collapse, pin/unpin and catalog overlay behavior remain unchanged. The reference screenshots inform styling only and do not authorize adding operational screens, chart metrics or routes.

## Architecture after refactor

Only `Shell.tsx` and `index.ts` live at the package root. `containers/shellContainer/` creates providers; `components/appRoutes/` composes domain-root routes; `components/protectedLayout/`, `sidebar/`, `topbar/`, `mobileNav/`, `themeChoices/` and `pinnedNavigation/` own view markup. `hooks/useShellState/`, `hooks/useShellActions/` and `storage/` isolate state and persistence. Theme resolution is in `@Components/hooks/useWorkspaceTheme/` so mobile system appearance does not overwrite a desktop preference.

## Persistent protected layout, no public layout

`AppRoutes` renders public school routes directly and public auth pages under `PublicOnly`. The protected pathless route nests `RequireSession → ProtectedLayout → Shell → Outlet`. The layout and its sidebar, toolbar, launcher and theme state are mounted once across protected child-route navigation; only the outlet's page changes. A regression test under `components/protectedLayout/__tests__` checks that the same shell DOM node survives child route changes.

The Dashboard is an outlet child at `/dashboard`, **not** the workspace wrapper. Do not create a `Shell` per domain or key it by pathname. The catalog and pinned navigation contain **only implemented apps** (currently Dashboard); no route to People/Integrations/Settings is exposed until built.

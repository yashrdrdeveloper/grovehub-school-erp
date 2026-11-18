# Shell
The shell is a composition-only package. It constructs React Router, TanStack Query, auth provider and the workspace chrome. Feature screens are imported from `@Domain/<name>` public APIs. No shell file imports `@Business`, calls `fetch`, or performs a use case.


## RDR School experience (November 2026)

The workspace displays **RDR School**; this is the requested school demo identity, not a claim that dynamic multi-tenant school branding exists. Later tenant-aware school metadata should replace the static label in `ShellContainer.tsx` without changing the shared shell.

### Visual system

Use a neutral SaaS canvas (`#F8FAFC`), ink (`#0B1220`), operational blue (`#2563EB`), soft borders (`#E2E8F0`), and semantic success/warning/danger tokens. Do not reintroduce the previous olive/green theme as the brand default. Shared visual tokens are in `app/styles.css`. Follow components-first composition, accessible focus states, and clear contrast.

### Appearance

`theme.ts` controls `system | light | dark` and persists explicit preference using `groovehub.theme` in browser local storage. The default is `system`; OS preference changes are observed. `app/main.tsx` resolves the preferred theme before rendering including public authentication pages. Theme controls are available in account preferences.

### Navigation and application catalog

The shared shell owns the narrow desktop sidebar, mobile navigation (Home, People, Apps, Accounts, Menu) and searchable, categorized app catalog. Active routes come from `registry/domains.registry.tsx`, and available shortcuts come from `registry/navigation.registry.ts`. The catalog may describe planned apps but must not create dead navigation to unimplemented routes or claim unimplemented operations work. Prioritize lightweight touch targets and a mobile bottom sheet for the catalog.

### Operational dashboard

The Dashboard domain shows four compact metrics and six at-a-glance cards: Pending for You, Admissions & Dues, Today & Upcoming, Reminders, Complaints & Requests, Audit & Compliance. Cards are responsive, compact, and ready for future domain APIs. Show explicit unavailable states for missing data; never display arbitrary pretend counts or fabricate school activity. As services are connected, their domain hooks/queries should replace static empty states, not fetch directly from the shell.

### Mobile usability

At smaller widths, sidebar becomes a dismissible drawer, app catalog a bottom sheet, KPI cards become two columns and operational panels become a single column on phones. Keep top-level tasks discoverable without horizontal page scrolling; individual data tables can scroll within their own containers.


## Workspace navigation and app catalog

- The top-left four-square button opens a **full-workspace catalog**. On desktop it fills the area to the right of the persistent sidebar; on mobile it uses the entire viewport. Escape or the close button dismisses it.
- **Pinned apps** are shown as the first catalog section and populate the sidebar in saved order. Use the pin/unpin control on any available application tile or hover over a pinned sidebar item to unpin it.
- Planned applications are visible for discovery but cannot be opened or pinned until their routes are implemented. Never simulate functional navigation to unimplemented modules.
- The desktop sidebar collapses to an icon rail. Its expanded/collapsed preference and pinned routes are stored in local browser storage (`groovehub.workspace.sidebar.collapsed.v1` and `groovehub.workspace.pinned.v1`).
- The mobile sidebar remains a drawer even if the desktop collapse preference is on. The mobile Apps entry opens the full-screen catalog.
- Navigation and catalog application metadata live in `layout/AppLauncher.tsx`. Shell reads the same catalog, so app labels/icons/links must remain consistent. Pinning is currently per-browser, not server-synchronized.
- Appearance remains **System** by default, with user-selected Light and Dark modes. The visual palette uses warm porcelain, charcoal and a restrained iris accent; the dark version uses deep aubergine surfaces.

# Dashboard domain

`@Domain/dashboard` exports only the implemented dashboard page. Its route is `/dashboard`, inside the persistent shell-level `ProtectedLayout`—the dashboard page **is not** itself the layout wrapper.

`pages/dashboard/Dashboard.tsx` is stateless presentation. `containers/dashboardContainer/` owns composition, `hooks/useDashboardCalendar/` owns calendar state, `helpers/dashboard.getDateParts.ts` is the pure typed helper with `helpers/__tests__`, and `constants/dashboard.metrics.ts` and `constants/dashboard.workspaces.ts` own static preview metadata. There is no `config/` folder or per-helper barrel.

**No fabricated live activity**: cards explicitly show unavailable/preview states. User-specific counts are not invented. Add a domain only after its route and business services exist, and never trigger business fetching from the shared protected shell.

Conformance and remaining work are maintained in `docs/engineering/frontend/FRONTEND_CODE_GUIDELINES.md` §22. Run architecture, TypeScript, Jest and production build checks before merge.

## Display strings and static preview states

Dashboard visible headings, hints, empty-state messages and accessibility labels call bare global `__()` during render. Stable copy should live in the owning `constants/<owner>.copy.ts`; any varying school data is passed as named `{{value}}` parameters, not concatenated. Static preview metadata remains owned by the dashboard and is **never** treated as live student or fee data. Closed UI props use component-exported constants such as `BOX_TAG` and `TEXT_TONE`. Do not import `@Base/i18n` into domain UI; the shell registers the browser-global resolver before mounting the dashboard.

## Display content at its static owner

`constants/dashboard.metrics.ts` and `constants/dashboard.workspaces.ts` translate human-readable names/descriptions when declared using literal `__('message')`, while technical tone classes and icon tokens stay raw. JSX renders these values directly and does not call `__(metric.name)`. Conditional view state belongs to owner-local tested helpers; no `let` declarations.

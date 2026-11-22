# Dashboard domain

`@Domain/dashboard` exports only the implemented dashboard page. Its route is `/dashboard`, inside the persistent shell-level `ProtectedLayout`—the dashboard page **is not** itself the layout wrapper.

`pages/dashboard/Dashboard.tsx` is stateless presentation. `containers/dashboardContainer/` owns composition, `hooks/useDashboardCalendar/` owns calendar state, `helpers/dashboard.getDateParts.ts` is the pure typed helper with `helpers/__tests__`, and `constants/dashboard.metrics.ts` and `constants/dashboard.workspaces.ts` own static preview metadata. There is no `config/` folder or per-helper barrel.

**No fabricated live activity**: cards explicitly show unavailable/preview states. User-specific counts are not invented. Add a domain only after its route and business services exist, and never trigger business fetching from the shared protected shell.

Conformance and remaining work are maintained in `docs/engineering/frontend/FRONTEND_CODE_GUIDELINES.md` §22. Run architecture, TypeScript, Jest and production build checks before merge.

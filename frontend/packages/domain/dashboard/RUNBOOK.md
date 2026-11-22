# Dashboard domain

The route exposed by `@Domain/dashboard` is composed in `pages/dashboard/`. The page root contains only `Dashboard.tsx` and `index.ts`; `containers/dashboardContainer/` gathers presentation inputs. `hooks/useDashboardCalendar/` owns the Today/Upcoming selection. `helpers/getDashboardDateParts/` is pure, typed and has adjacent Jest tests. Visual sections live in focused components, and preview metadata lives under `config/`.

**No live activity has been fabricated.** Empty metrics and cards explicitly indicate preview states. Any future domain service data must be introduced through domain hooks and business-owned API contracts, without fetching in presentational components or the shell.

Run `pnpm check:architecture`, `pnpm typecheck`, `pnpm test -- --runInBand`, and `pnpm build` from `frontend/`.

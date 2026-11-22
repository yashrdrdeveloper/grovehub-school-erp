# GroveHub styles package

`@Styles` owns **generic design decisions**: typography families and sizes, semantic light/dark colors, spacing, control heights, radii, shadows, focus, motion, breakpoints and Tailwind utilities.

Import `@Styles/index.css` **once** in `frontend/app/main.tsx`. Use `@Styles/tokens/breakpoints` only in browser-aware hooks to share viewport values. CSS files are deliberately separate from React so styles can be used by CSS Modules, Tailwind and plain CSS without importing components.

The canonical root theme overrides were removed from the legacy app stylesheet. The architecture check now verifies light/dark token parity and single style ownership.

The central theme is the document `data-theme="light" | "dark"`. The existing theme resolver respects the OS on mobile and a desktop preference elsewhere. No design token changes that behavior. The `--gh-*` variables are authoritative; `--canvas`, `--surface`, `--accent`, and `--rdr-*` are temporary aliases for legacy CSS during migration.

`@Components` owns actual `Button`, `Input`, `FormField`, `Typography`, `Card`, `Container` and other reusable controls, their prop variants, accessibility, native leaf `dataTestId` and component-specific styles. `@Styles` never imports React, business data, domain pages or component code. Domains should use `@Components` rather than copy dimensions and colors.

**Preferred styling:** Tailwind v4 utility classes plus shared CSS variables and CVA in components. CSS Modules (`.module.css`) are acceptable for complex component-specific selectors. Do not introduce Sass/SCSS without a genuine need; we currently need neither a Sass compiler nor a second source for design tokens.

**Preservation:** Tailwind Preflight is intentionally not imported yet. Existing `frontend/app/styles.css` contains legacy shell/feature rules; relocate those gradually to their owning components/domains. It must not reintroduce global theme tokens. Do not globally restyle the RDR auth/public site during infrastructure migration without visual regression evidence.

**Verification:** `pnpm check:architecture`, `pnpm typecheck`, `pnpm exec jest --runInBand`, `pnpm build`. Responsive and dark/light browser testing is required before deleting legacy screen styles.

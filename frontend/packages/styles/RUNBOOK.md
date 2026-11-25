# GroveHub styles package

`@Styles` owns **generic design decisions**: typography families and sizes, semantic light/dark colors, spacing, control heights, radii, shadows, focus, motion, breakpoints and Tailwind utilities.

Import `@Styles/index.css` **once** in `frontend/app/main.tsx`. Use `@Styles/tokens/breakpoints` only in browser-aware hooks to share viewport values. CSS files are deliberately separate from React so styles can be used by CSS Modules, Tailwind and plain CSS without importing components.

The canonical root theme overrides were removed from the legacy app stylesheet. The architecture check now verifies light/dark token parity and single style ownership.

The central theme is the document `data-theme="light" | "dark"`. The existing theme resolver respects the OS on mobile and a desktop preference elsewhere. No design token changes that behavior. The `--gh-*` variables are authoritative; `--canvas`, `--surface`, `--accent`, and `--rdr-*` are temporary aliases for legacy CSS during migration.

`@Components` owns actual `Button`, `Input`, `FormField`, `Typography`, `Card`, `Container` and other reusable controls, their prop variants, accessibility, native leaf `dataTestId` and component-specific styles. `@Styles` never imports React, business data, domain pages or component code. Domains should use `@Components` rather than copy dimensions and colors.

**Preferred styling:** Tailwind v4 utility classes plus shared CSS variables and CVA in components. CSS Modules (`.module.css`) are acceptable for complex component-specific selectors. Do not introduce Sass/SCSS without a genuine need; we currently need neither a Sass compiler nor a second source for design tokens.

**Preservation:** Tailwind Preflight is intentionally not imported yet. Existing `frontend/app/styles.css` contains legacy shell/feature rules; relocate those gradually to their owning components/domains. It must not reintroduce global theme tokens. Do not globally restyle the RDR auth/public site during infrastructure migration without visual regression evidence.

**Verification:** `pnpm check:architecture`, `pnpm typecheck`, `pnpm exec jest --runInBand`, `pnpm build`. Responsive and dark/light browser testing is required before deleting legacy screen styles.

## Text and tokens

Styling owns colors, typography, spacing, focus and responsive tokens, **not translated copy**. Do not generate essential user-facing messages with CSS `content:`; text belongs in accessible JSX and is resolved through the shell-provided global `__()`. Translatable text must be able to expand without clipping and respect 320/375/390/720 px responsive layouts. A locale-aware future must handle font coverage, directionality and text expansion without creating per-domain global palettes.

## No mutable display state

Copy expansion and responsive text remain presentation concerns. Never implement alternate CSS/class/token selection by `let` reassignment inside JSX; use a pure owner-specific renderer/helper returning an immutable variant. Translation calls always receive literal source messages. Styling technical class strings are not translated.

## Native scrollbars

`@Styles/scrollbars.css` owns the application-wide scrollbar appearance. The root ShellContainer installs one passive capture-phase scroll observer for document and nested scroll containers, including portal content. Each scroll owner reveals its thumb independently; after 1000 ms without a scroll event, the active marker is removed. No React render or wheel-event interception is involved.

Use native overflow on scroll owners, and `scrollbar-gutter: stable` for vertical or two-axis scroll containers where the appearance of overflow must not change content width. Do not toggle scrollbar widths, overflow or padding to show and hide a scrollbar. Do not define component-local scrollbar colors, scrollbar widths or hidden scrollbar pseudo-elements. The global desktop style uses fixed 8px scrollbar geometry and rounded thumbs in Chromium/WebKit; native thin scrollbars in other engines retain their platform geometry and reveal on hover/focus/scroll. Touch and forced-colors environments retain their native appearance and behavior. Classic scrollbars reserve space; OS overlay scrollbars need not. A newly appearing horizontal overflow can still consume height on classic-scrollbar platforms.

Verify wheel/trackpad, thumb dragging, keyboard scrolling, nested dialog scroll, table horizontal/vertical scrolling, dark/light themes, and content width before/after thumb reveal in browsers. No third-party scroll library is needed.

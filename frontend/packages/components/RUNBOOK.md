# Shared UI components

**Ownership:** `packages/components` owns all reusable React UI and component-specific styling. It imports design tokens from **`@Styles`** and pure utilities from **`@Base`**; it never imports `@Business`, `@Domain` or `@Shell`.

**Theme/typography/spacing:** Colors, light/dark mode, font families, text scales, focus, shadows, radii, control heights and spacing are owned by `frontend/packages/styles`. Shared component CSS in `styles/primitives.css` consumes `--gh-*` variables and should not invent a second palette.

**Public components:**
- `Button`: `variant=` `primary | secondary | outline | ghost | danger`, `size=sm | md | lg`, loading, asChild, dataTestId
- `InputControl`: lowest native input, `controlSize=sm | md | lg`, dataTestId; use `unstyled` only to preserve a legacy feature layout
- `Input`/`TextField`: label/hint/error/required with `controlSize`, dataTestId and a shared InputControl leaf
- `PasswordField`: shared InputControl plus visibility action and accessible labels
- `Typography`/`Text`: semantic element, named variant/tone, `size` and `weight` mapped to design tokens
- `Form`: native form wrapper, `layout=stack | inline`, `spacing=sm | md | lg` and dataTestId; no form submission business logic
- `FormField`, `Card`, `Container`, `Dialog`, `DataTable` and feedback primitives: accessible structure/variants, with feature-owned stable test IDs

`dataTestId` must reach the last native DOM element as `data-testid`; typed component APIs own that propagation. Do not auto-generate one fixed ID for all instances. Components are **presentational**. Domain code supplies controlled values, handlers, loading/error states and test IDs; business APIs remain outside this package.

**Component structure:** Compound components use one named subfolder under `components/` per primitive. `Dialog/components/root`, `portal`, `overlay`, `content`, `title`, `description` and `close` own their respective `*.tsx`, props contracts and entries; consumers import them from `@Components/components/Dialog`. The old `Dialog/parts.ts` has been deleted. `FormField` owns pure tested label/test-ID and ARIA-description assembly helpers. Shared inputs and password fields reuse them and `Box`/`Button` leaves instead of inline regex, native button composition or duplicated code. `ErrorBoundary` uses a genuinely scoped `ErrorBoundary.module.css` and recovers on a route reset.\n\n**Styling standard:** Tailwind v4 utilities + CSS custom properties + CVA are preferred. Use owner-scoped `.module.css` when component selectors need isolation; existing global CSS remains documented migration debt and must not be re-labeled as scoped. No Sass dependency is needed. Do not copy visual CSS into each page: promote reusable behavior to `@Components`, and retain layout-only CSS locally to the domain.

**Verification:** `pnpm check:architecture`, `pnpm typecheck`, `pnpm exec jest --runInBand`, `pnpm build`. Responsive/manual/visual checks are still required when removing existing CSS.

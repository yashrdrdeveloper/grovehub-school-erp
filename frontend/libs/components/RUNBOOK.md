# Shared UI components

**Ownership:** `libs/frontend/components` owns all reusable React UI and component-specific styling. It imports design tokens from **`@Styles`** and pure utilities from **`@Base`**; it never imports `@Business`, `@Domain` , `@Shell` or `@Assets`.

**Theme/typography/spacing:** Colors, light/dark mode, font families, text scales, focus, shadows, radii, control heights and spacing are owned by `libs/frontend/styles`. Shared component CSS in `styles/primitives.css` consumes `--gh-*` variables and should not invent a second palette.

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

## Closed component props and displayed messages

Every component must export its own closed-set runtime options and derive public prop types from the options: `BUTTON_VARIANT` / `BUTTON_SIZE`, `BOX_TAG`, `TYPOGRAPHY_TAG`, `TEXT_TAG`, `TEXT_TONE`, `TEXT_SIZE`, `CARD_VARIANT`, `FORM_LAYOUT`, and other owner-specific values. Re-export the **same underlying token** for a forwarded prop, not a duplicate string map; for example `TEXT_TONE` aliases Typography's actual tone choices. No `variant="primary"`, `size="lg"`, `as="h3"` or other raw closed-option strings in feature JSX.

Presentational features call bare `__()` for user-facing labels, headings, accessible names, placeholders, errors and tooltips. They must not import `__` or mutate `window.__`. The global resolver is provided by `ShellContainer` before routes mount, while Jest installs the real pure resolver in test setup. Component contracts still forward actual DOM attributes and `dataTestId` precisely. See the unified frontend guideline Revision 13 for the literal-versus-technical-string rules.

**Enforced enums:** `pnpm check:component-enums` rejects raw closed prop literals for shared components. `Input`, `PasswordField` and `TextField` publicly alias the same underlying InputControl sizing values as `INPUT_SIZE`, `PASSWORD_FIELD_SIZE` and `TEXT_FIELD_SIZE`. The aliases are public conveniences, not new size scales; always re-use the owner's actual runtime value and derive types from it.

## Pure display selection and immutable bindings

Presentation copy in owned `constants/*.copy.ts` is declared as `__('literal')`; components later use `COPY.key` without re-wrapping. Per-instance dynamic labels call the global `__('literal')` inside branch-return helpers/renderer functions. Components must not import the formatter. Mutable `let` declarations, JSX component reassignment and ternaries are forbidden. `PasswordField/renderer/passwordField.visibility.ts` demonstrates a typed pure icon/label/input-type selector.

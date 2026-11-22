# RDR authentication UI — design and acceptance standard

> **Engineering ownership:** Follow [Frontend Code Guidelines, Revision 5](./FRONTEND_CODE_GUIDELINES.md) for file layout, pure-helper ownership, flat constants, nested hooks, central icons, stable `dataTestId`/`data-testid` selectors and no-ternary rendering. This document governs auth **visual and interaction acceptance**, not engineering folder structure.

## Scope and source

Applies to sign-in, forgot-password, set-password, reset-password, invitation, link verification, success, unavailable, expired, and exhausted-attempt states.

This design standard adapts the [UI/UX Pro Max skill](https://github.com/nextlevelbuilder/ui-ux-pro-max-skill/blob/main/.claude/skills/ui-ux-pro-max/SKILL.md) (accessibility, interaction/touch, layout, performance, forms, stack guidelines) to **existing RDR School branding**, not the skill's generic example palettes. This skill is design guidance, not an authentication package.

### Invariants — do not change without product approval

- **Original** `SchoolMonogram` for both desktop and mobile; never redraw it.
- No username or public self-registration; preserve all auth, refresh, invitation, reset, and two-attempt flows.
- Charcoal mobile dark mode (`#171E2A`), ivory/white light mode, RDR navy (`#1F2A44`) primary, and restrained gold text accent.
- One shared `LoginLayout`; identical labels, success/error actions, and security messaging across devices. The responsive layout may rearrange control positions, not change behavior.
- Security footer: lock icon + “A secure workspace for RDR School” on one line when space permits, followed by “© RDR School” and “Haluhera · Rewari”.
- Keep headings and action text untruncated, including “Forgot your password?”.
- Never store tokens or passwords in local storage. Retain existing auth handling.

## Design decisions

| Layer | Light | Dark |
|---|---|---|
| Form surface | `#FFFFFF` | `#171E2A` |
| Main text | `#1F2A44` | `#F3F6FF` |
| Secondary text | `#5F6B7D` | `#BDC8D7` |
| Input | `#F8FAFC` | `#222E40` |
| Border | `#BAC6D6` | `#596A82` |
| CTA | `#1F2A44` | `#202B49` |
| Focus outline | `#506E9B` | `#ACC7ED` |
| Link | `#465F8A` | `#A8BDE2` |

These values are **auth-specific semantic tokens** in `optionOneAuth.css` and should not create a second competing global theme.

Choose the skill's **Accessible & Ethical** + **Minimalism & Swiss Style** characteristics for hierarchy and restraint; treat **Trust & Authority** as a content priority. Avoid generic dashboard/glassmorphism cards, neon gradients, and ornamental layouts that interfere with the form. The desktop SVG is decorative and hidden on mobile.

## Layout and state rules

1. Limit the form to 440px on desktop and 380px on phones; do not stretch controls across the right panel.
2. Use the available viewport (`100dvh`), but **never hide controls** to force an artificial no-scroll layout. Standard phones should fit compactly; keyboard, zoom, large text, or landscape may scroll **inside** the form panel.
3. No horizontal scrolling at 320px, 375px, 393px, 820px, 1024px, or wider. Screens must respect safe-area insets.
4. All icon buttons, recovery links, and primary actions should have a usable minimum ~44px target. Icons must have an accessible name when actionable.
5. Mobile text inputs are 16px to prevent automatic Safari zoom. Keep associated labels in the DOM even when visually hidden by the approved mobile placeholder treatment.
6. Reserve a two-line region below credential fields, including when there is no error; a validation message must never move the next field or CTA.
7. Reserve a fixed region for form-level and API errors. Long text must remain reachable and announced, without shifting primary controls.
8. Keep “Forgot your password?” and “Forgot password?” unwrapped at standard phone sizes. At very small text zoom widths, preserving legibility and access is higher priority than forced one-line text.
9. Success / expired / missing-link / checking states use the **same** content and action order on desktop and mobile; only spacing, text size and decorative art vary.
10. Use accessible `aria-invalid`, `aria-describedby`, and polite announcements for inline errors. Never rely on red alone to communicate failure.
11. Honor reduced-motion preference and never animate an element's dimensions to reveal errors.
12. The original school artwork may be decorative; it must not replace instructions or use fake stock people.

## Responsive verification matrix

| Device | Size / condition | Expected |
|---|---|---|
| Small phone | 320×568 | No horizontal overflow, CTA remains reachable; scroll only if necessary |
| iPhone 14 Pro | 393×852 | Sign-in, empty-field errors, forgot, forgot-success, reset and invitation actions visible without routine scrolling |
| Android phone | 360×800 | Same action order and system-derived dark/light palette |
| iPad Air portrait | 820×1180 | Compact split layout; academic art, form and footer never clip |
| Tablet landscape | 1180×820 | Form max 440px; no clipped success/error messaging |
| Desktop | 1440×900 | Split visual hierarchy, centered readable form |
| Accessibility | 200% text zoom / software keyboard / landscape | Scrolling permitted so focus and controls stay reachable |

### Manual acceptance scenarios

- Submit blank login; both error messages appear **without shifting email/password inputs**.
- Correct only the email; remaining password error stays associated with its input.
- Invalid credentials; feedback is visible and screen-reader announced, while the Sign in button stays in position.
- Toggle password; the focusable visibility control remains in the same location.
- Forgot password: empty email, invalid email, success, countdown, resend, and another address.
- Reset and invite: token check loading, ready, save error, two-attempt exhaustion, expired link and success.
- Toggle desktop theme, then resize to mobile; mobile uses the system preference.
- Test keyboard tab order, focus styles, contrast, and reduced-motion in both themes.
- Confirm no asset requests are required for the desktop illustration; keep the RDR logo component unchanged.

GitHub Actions typechecking and unit tests are **not** a substitute for these browser viewport checks.

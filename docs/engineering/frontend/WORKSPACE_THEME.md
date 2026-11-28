# Workspace theme

The school workspace adopts the compact shadcn-inspired internal platform layout with a blue primary, not orange. Light mode uses zinc canvas #FAFAFA, white surfaces #FFFFFF, heading #F4F4F5 and borders #E4E4E7. Dark mode uses canvas #09090B, surface #101012, heading #18181B and borders #29292D. Blue primary #2563D9 and focus #3E83F6 are shared by both themes; dark-mode text accents use #93C5FD.

Semantic theme files own colors; component CSS Modules own layout. Tables have compact typography, neutral hover, sticky headers, a full-height body and a stable footer for fillHeight pages. Preserve existing toolbar order, search, sort, filtering, selection, pagination, pinning, resizing, overlay scrollbars, mobile cards and backend contracts. The app catalog uses compact cards with separate pin actions. The sidebar and topbar keep existing navigation and theme behavior.

Existing grouped header and status color families, public-site and authentication illustration tokens remain intact. This change does not add an internal frontend or change authentication or API behavior. Browser, accessibility and architecture gates remain required before release.

## Mobile workspace chrome

At widths up to 720px, hide the workspace topbar completely. Home, Apps and Menu in the bottom navigation are the only persistent workspace navigation. Account preferences, appearance and sign-out remain available through Menu. The shared shell reserves safe-area-aware bottom spacing so content and table pagination can scroll above the bottom navigation. Desktop chrome remains unchanged.

The document must not reserve a scrollbar gutter on phones (including desktop device emulation). Hide the native document scrollbar without disabling vertical scrolling; preserve component scroll areas and forced-colors accessibility. This rule applies to all frontend routes, not only authentication. Validate light and dark modes at 320, 375, 393 and 720px, including sidebar and catalog overlays and the last row/footer of each mobile table.

On phones, Menu opens the existing shared sidebar as a full-viewport navigation page (100% width, 100dvh), not a narrow drawer. The underlying workspace and bottom navigation are hidden while the menu is open. The menu has its own safe-area-aware scrolling, a close action, full-size app labels and visible unpin actions even when the desktop sidebar preference is collapsed. Account and appearance controls remain accessible in the menu. Desktop layout and collapse persistence are unchanged.

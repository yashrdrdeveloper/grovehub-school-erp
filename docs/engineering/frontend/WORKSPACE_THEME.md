# Workspace theme

The school workspace adopts the compact shadcn-inspired internal platform layout with a blue primary, not orange. Light mode uses zinc canvas #FAFAFA, white surfaces #FFFFFF, heading #F4F4F5 and borders #E4E4E7. Dark mode uses canvas #09090B, surface #101012, heading #18181B and borders #29292D. Blue primary #2563D9 and focus #3E83F6 are shared by both themes; dark-mode text accents use #93C5FD.

Semantic theme files own colors; component CSS Modules own layout. Tables have compact typography, neutral hover, sticky headers, a full-height body and a stable footer for fillHeight pages. Preserve existing toolbar order, search, sort, filtering, selection, pagination, pinning, resizing, overlay scrollbars, mobile cards and backend contracts. The app catalog uses compact cards with separate pin actions. The sidebar and topbar keep existing navigation and theme behavior.

Existing grouped header and status color families, public-site and authentication illustration tokens remain intact. This change does not add an internal frontend or change authentication or API behavior. Browser, accessibility and architecture gates remain required before release.

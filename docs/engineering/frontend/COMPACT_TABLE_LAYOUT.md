# Compact desktop table layout

Table supports opt-in fillHeight. The workspace supplies its available height; title, controls and pagination stay outside the scrolling data body. The desktop header uses pointer-capture resize dividers and arrow-key resizing with column min/max constraints. A narrow fixed Actions column remains at the right. Theme-owned resize colors distinguish normal and focused dividers.

ColumnManager uses List, DragList and Checkbox. DragList is a controlled reusable list; a drag handle reorders items, and arrow keys provide keyboard reordering. The lock toggles left pinning; existing right-pinned columns can be unlocked. The default lock limit is three, excluding the fixed Actions column. Table owns and enforces ordering, visibility, pinning and width preferences. Mobile touch drag is not implemented; desktop HTML drag and keyboard are supported.

Filters use inline Field / Operator / Value rows on desktop, and Manage / Add / Edit on mobile. Apply commits only valid drafts. The current server supports AND conditions, so OR and active filter chips are not rendered.

User Setup renders actual account creation and update timestamps via @Base/localization/getFormattedDate, using Intl and the browser timezone by default. Sorting uses the original backend field values; missing or invalid timestamps render a dash.

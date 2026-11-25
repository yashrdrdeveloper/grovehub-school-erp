# Universal Table — Phase 2.1 implementation status

The initial shared Table, immutable Column builder, Filters, Sort, QuickFilters, Pagination, local and async selections, and responsive Sheet have been implemented. The design extends the existing Login theme tokens with a common control radius and heights. Desktop dialogs and mobile sheets share state.

User Setup is registered at /users/setup. The verified account API supports search, invitation-based creation, username update, lifecycle transitions, and invitation resend. Account search accepts one sort; email is filterable but not sortable. Verified statuses are PENDING_SETUP, ACTIVE, SUSPENDED and DEACTIVATED. Primary role (single Select) and secondary roles (MultiSelect) are provided as a reusable UserRoleFields component; actual role assignment must wait for server-owned role lookup and persistence APIs. No status or password is editable in the user form.

Filter UX: manage -> add/edit -> field selector -> operator -> value -> save draft -> apply. No applied chips are rendered on the table page.

Outstanding: full anchored popovers, calendar engine and timezone-aware date/time, option-label resolution and virtualization, client-side and advanced table engines, column resizing/pinning/reordering/persistence, bulk selection, full mobile card renderer, backend role and personal detail integration, and production validation. The existing architecture gates remain required; this is incremental source work, not release certification.

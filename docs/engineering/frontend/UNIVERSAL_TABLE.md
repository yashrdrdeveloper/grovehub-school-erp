# Universal Table — Phase 2.1

Shared Table supports DEFAULT, VIRTUALIZED, EXPANDABLE, GROUPED and TREE row structures, server and client data modes, typed action dispatch, configurable header/subheader/footer, column builder, selection, bulk-action dispatch, pinning, visibility, ordering, step-based column width changes and density. Preferences contain only display settings and must be scoped to a signed-in account. Virtual rendering uses fixed-height rows on desktop; mobile uses bounded incremental card rendering. Expandable rows deliberately use their natural height.

Filters and Sort use anchored desktop overlays and responsive bottom sheets. Filter changes are drafts until Apply; no active filter chips appear on the table page. Async selection uses cursor pagination, cancellation, cached labels and optional selected-value resolution. Native date input and custom Calendar coexist; range controls use a single collapsed trigger; DateTime sends UTC ISO values.

User Setup uses cursor pagination, server-supported search fields and one sort. Roles are returned from a server-owned starter catalog. Primary role is required for the new form and secondary roles are optional and distinct. Existing legacy account requests remain compatible. Role assignment is tenant-admin authorized, persisted, version checked, and does not alter the configured gateway administrator policy. Full RBAC permission evaluation and custom role management are separate work.

Profile reads and name updates go through Identity's tenant-scoped account boundary to Foundation. Account search batches Foundation profile reads for display. Personal name edits and account role edits have separate saves because they belong to separate services. Email and phone changes require a verified contact workflow; they are not silently overwritten.

Source implementation requires Prettier, architecture checks, typecheck, unit tests, annotation generation, migrations and browser verification before release. No CI or security gate should be disabled.

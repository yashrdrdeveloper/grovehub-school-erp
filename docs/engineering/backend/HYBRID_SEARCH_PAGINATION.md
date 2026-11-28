# Shared hybrid search pagination

The shared search request accepts `pageInfo.pageNumber` (1-based), `pageInfo.rowSize` (1–100, default 25) and an optional opaque `pageInfo.cursor`. The previous `rows` request property remains a compatibility alias; conflicting sizes are rejected.

- With a cursor: use the existing ordered keyset predicate and never apply an offset. Numbered cursors must match the requested target page and row size.
- Without a cursor: apply `OFFSET (pageNumber - 1) * rowSize` after the same search predicates and stable ordering. This enables arbitrary jumps without sequential API calls.
- Numbered responses contain totalRecords, totalPages, pageNumber, rowSize, rows, hasPrevious, hasNext and nextCursor when available. The count excludes both cursor and offset constraints.
- Identity and Foundation accept the shared numbered contract. Existing unnumbered Foundation clients retain their cursor-only behavior.
- Frontend User Setup caches cursor anchors by page number, resets them on query or page size changes, and uses a direct offset fallback for unknown pages.
- Multi-sort cursor order and the unique-ID tie-breaker remain intact. Counts and page results are separate queries and may differ during concurrent writes; numbered offset jumps are more costly at deep pages. This is not a snapshot guarantee.

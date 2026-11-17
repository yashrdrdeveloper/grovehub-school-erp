# Backend Search Contract and Query Design

Scope: generic search behavior for backend services, independent of resource or product.

## Ownership and responsibility

- Transport DTO represents user-supplied search terms, filters, ordering, page size and cursor.
- Validation verifies field/operator allowlists, value types, sizes, combinations, authorization and limits.
- Search configuration maps public field names to trusted server-side SQL expressions, relations, sort keys and cursor keys.
- Repository implementation builds parameterized queries from validated configuration.
- Application service enforces tenant and actor scope **before** applying caller filters and cannot be bypassed by a query.
- Response DTO exposes records, a `hasNext` indicator and an opaque next cursor if supported.

## Query stages

```text
Decode -> Normalize -> Validate -> Authorize/Scope -> Compile -> Execute -> Map -> Encode cursor
```

Never accept raw SQL column names, table aliases, operators or expressions from clients. Distinguish full-text search from exact matching. Search configurations are implementation-owned and reviewed with indexes.

## Field and operator model

Every searchable field declares: public name, value type, permitted operators, indexed/expensive status, nullable semantics and relation/cardinality. Filter operators might include `eq`, `ne`, `in`, `range`, `contains`, `isNull`; enable only those safe for that specific field. A missing filter is not equivalent to matching SQL NULL. Escape LIKE wildcards when literals are intended.

Relations should compile as trusted `EXISTS` subqueries when this avoids join multiplicity. Fixed predicates and parent-to-child links are declared by code, not copied from user input.

## Sorting and pagination

- Use stable deterministic ordering: requested sortable field(s) plus a unique immutable tie-breaker.
- Allowlist directions and null handling. Cursor must capture every sort key required to reproduce the ordering.
- Prefer keyset pagination for large or changing datasets; avoid huge OFFSET scans.
- Cap page size and query complexity; fetch `limit + 1` to determine `hasNext` if suitable.
- Cursor is opaque, versioned and integrity-protected where client tampering could affect scope.
- A cursor must not override tenant scope, permissions or validated filter context.
- Define how deleted/updated rows affect traversal; do not promise snapshot consistency without a snapshot strategy.

## Performance and safety

Use parameterized values, scoped SELECT projections and appropriate DB indexes. Review query plans, count queries, relation filters, text-search indexes and worst-case selectivity. Set DB/request timeouts and avoid N+1 fetching. Require an explicit policy for unbounded or cross-tenant searching.

## Example (generic only)

```json
{
  "searchTerms": ["sample"],
  "filters": [{"field": "status", "operator": "eq", "values": ["ACTIVE"]}],
  "sort": {"field": "createdAt", "direction": "desc"},
  "limit": 25,
  "cursor": null
}
```

The example illustrates a **contract shape only**, not a guaranteed endpoint or actual supported field. Actual DTOs, operators and response property names must be taken from implementation and API contracts. Service engineering docs record their configured fields, permissions, indexes and deviation decisions; product docs give real endpoint requests and expected results.

## Acceptance checklist

A search feature is ready only when it covers invalid operators, injection attempts, unauthorized fields, empty inputs, null values, stable ordering, cursor tampering, duplicate relation rows, page boundaries, concurrent writes and execution-plan review.

# Search, SearchDetails and Cursor Pagination

**Generic reference** • Search remains a shared technical capability; every resource owns its own field allowlist, projection inputs and hydration policy.

## 1. Vocabulary and separation

- **Source of truth:** normalized aggregate tables. Never substitute a search projection for authoritative uniqueness or domain state.
- **Search projection / `search_details`:** denormalized, rebuildable text used for free-text matches. It contains only approved fields and is not authoritative.
- **Search request:** untrusted public filters, text, sorting and opaque cursor.
- **Validated query:** typed, normalized and allowlisted query accepted by a repository.
- **Search config:** developer-maintained trusted database expressions, joins, cursor rules and output mapping.
- **Hydration:** repository or service uses stable identifiers to load complete required relation data; avoids N+1.
- **Cursor:** stable, opaque keyset marker for deterministic forward pagination. Not a raw SQL offset.

## 2. Canonical placement

```text
libs/core/search/
├── search.go                  # only public capability facade
├── contracts/                 # private generic request/query contracts
├── validators/                # generic operator/value/cursor rules
├── helpers/                   # trusted keyset/parsing mechanics
└── mappers/                   # generic query conversions if needed

services/<module>/internal/app/
├── dto/<resource>/request.go   # resource's public search DTO, if distinct
├── validators/<resource>/
│   ├── <resource>.go           # field/operator limits and normalization
│   └── impl/<resource>.go
├── repository/<resource>/
│   ├── <resource>.go           # Search(ctx, search.Query[ID])
│   ├── impl/<resource>.go
│   └── helpers/search/
│       ├── config.go          # trusted field expressions
│       └── cursor.go          # resource-specific keyset extraction if needed
└── services/<resource>/
    └── operations/search/     # Prepare -> Query -> Hydrate when multistep
```

Do **not** import a shared capability's private `contracts/`, `helpers/`, `mappers/`, or `validators/` from services; import only `libs/core/search`.

## 3. Pipeline and owners

```text
public search request
  -> resource validator: allowlisted fields/operators, normalized values,
                         safe row cap, sort parity, cursor token format
  -> typed search.Query[ID]
  -> repository: trusted SQL expressions and parameterized conditions
  -> keyset execution and result IDs
  -> relation hydration in batches
  -> response {items, pageInfo: {rows, hasNext, nextCursor?}}
```

The validator rejects unknown filters, unsupported operators, too many values and user-supplied raw columns. The repository rejects untrusted SQL even if a caller bypasses the HTTP validator. `search_details` writes belong to the aggregate owner, **not** to the read-only Search operation.

## 4. Projection invariants

1. Rebuild from the complete aggregate already available in the write-side operation; don't add an expensive second SQL aggregate query for data already in memory.
2. Normalize/canonicalize values with the same shared normalizer used by persistence.
3. Include only explicitly permitted searchable fields. Avoid secrets, credentials and unnecessarily sensitive data.
4. Preserve deterministic ordering/encoding to prevent non-idempotent projection churn.
5. Update projection in **the same UnitOfWork** as the aggregate and its relations. A transient mismatch should not be committed.
6. Document rebuild/backfill strategy if fields, tokenizer, format or search algorithm change.
7. Authoritative uniqueness and conflict detection always query canonical columns/constraints, not projection text.
8. Restrict projection visibility and indexing policy where stored data is sensitive.

## 5. Allowed fields, operators and sorts

Document one resource-owned table of `field -> value type -> operators -> normalization -> SQL mapping`. A contact-like relation may require a fixed relation predicate, trusted alias and typed values. SQL identifiers/expression fragments are code-maintained constants; user input only supplies parameter values.

Operator behavior is shared where semantics are actually generic. Resource validators can further restrict operators. Unknown sort fields or sort directions must fail rather than fall back to a risky default.

## 6. Cursor and stable order

Keyset pagination uses a stable ordered tuple, for example `(sort_value, unique_id)`. The cursor includes enough trusted, validated data to reproduce the ordering and ties; it must be opaque to clients, typed on decode, and rejected if malformed or incompatible with filters/sort. When a unique ID is used as a tie-breaker it must have matching direction and type. The generic search engine cannot assume all module IDs are UUIDs.

- Use `LIMIT rows + 1` (or equivalent) to determine `hasNext`.
- Return `nextCursor` only when another page exists.
- Avoid offset pagination, which shifts under concurrent writes.
- Sort by deterministic, indexed expressions when possible.
- Document consistency under concurrent updates; keyset avoids offset duplication but cannot guarantee snapshot isolation across requests.

## 7. Search annotations

`@Repository` identifies resource ownership; `@Searchable` may generate a typed repository search helper when it references a developer-maintained config and supported method signature. Optional `@Filterable`, `@Sortable`, `@Cursor` provide parity checks. They do not infer safe SQL, authorize a filter, or validate a caller.

Example only:

```go
// @Repository entity=Record
type RecordRepository interface {
    // @Searchable impl=Repository entity=Record config=search.RecordConfig alias=r
    Search(ctx context.Context, query search.Query[string]) (search.Result[Record], error)
}
```

Check the live generator for exact accepted signatures/config placement; this is a schematic resource-neutral illustration.

## 8. Indexing, cost and maintenance

- Design `btree` indexes for equality/range/sort, `gin` or trigram indexes when text predicates justify them. Validate plans with representative volumes.
- Avoid unbounded `ILIKE` scans on growth paths and indexing raw PII without risk review.
- Count query + hydration query count and bytes. N+1 is a correctness/performance defect.
- Rebuild safely with backfill, dual-read/check or cutover on projection format changes.
- Treat search text and cursor tokens as untrusted and potentially sensitive when logging.

## 9. Tests

Test empty query, one relation, multiple relation filters, unknown field, invalid operator, escaping/SQL injection strings, normalization parity, sorting ties, forward cursor pages, deleted/updated rows, 0/1/max page sizes, projection rebuild, concurrent updates, and parameterized SQL. Include a PostgreSQL-backed query-plan test for large/high-cardinality tables where warranted.


## Field identifiers, enum rules and helper decomposition (mandatory)

- The resource DTO boundary MUST declare canonical filter/sort field-name constants (or a reviewed generator must derive them). Validator maps, repository search configuration, cursor extraction and documented request schemas MUST share those constants.
- Go JSON struct tags are literal values; a type-safe constant does not automatically update a tag. Add reflection-based parity assertions in **tests** (not the runtime request path) and update OpenAPI in the same commit.
- Domain status values MUST be checked using the owning entity's enum constants; never duplicate status string values inside search validators.
- Keep `validators/<resource>/impl/<resource>.go` as the thin entrypoint; substantive field specs, normalization functions, cursor/sort policies belong in `validators/<resource>/helpers/search.go` or other focused purpose files.
- Only trusted repository search config maps identifiers to SQL expressions. All parameters remain bound, and a filter unsupported by either validator or repository must fail closed.

# Foundation Person service operations

The Person service composes all public methods from operation-owned packages.
The same HTTP controller, client, DTOs, service interface and error codes remain
in place. The PersonService implementation contains only dependency composition.

## Folder ownership

Under backend/services/foundation/services/person/operations/:

- create/create.go: Create chain (Prepare -> Idempotency -> Contacts -> Persist), with flows, private helpers and strategy interface and implementation.
- getbygid/getbygid.go: GetByGID entry point and a single @Strategy interface and implementation. No chain or flow is required.
- update/update.go: Update chain (Load -> Apply -> Persist), with a typed flow, strategies and private helpers/patch.go.
- search/search.go: Search chain (Prepare -> Query -> Hydrate), with a typed flow, strategies and private helpers/cursor.go.
- deactivate/deactivate.go: Deactivate chain (Load -> Prepare -> Persist), with a typed flow and strategies.

services/person/helpers/ contains helpers shared between operations: contact
construction/comparison, contact association, search projection, ReadPerson
(GetByGID/Update/Deactivate), and MapUpdateError (Update/Deactivate).
Nullable/optional string converters remain in libs/core/helpers.

## Annotations and runtime guarantees

- @Service marks the composed PersonService in services/person/impl.
- @Strategy marks typed strategy interfaces and concrete implementations in every operation.
- @ChainOfResponsibility marks Create, Update, Search, and Deactivate. Each
  method explicitly invokes the typed core servicepipeline.Run.
- @Transactional marks Create and Update write handlers; those explicitly
  use the existing transaction coordinator.

Go comment annotations are checked metadata, not implicit execution or DI.
Each operation is explicitly constructed by NewPersonService.

## Behavior preserved

- The same service interface methods: Create, GetByGID, Update, Search, Deactivate.
- Controller validation and existing Update service-level normalization and
  validation (important for direct service callers and PATCH semantics).
- Create idempotency-key and contact uniqueness recovery after rollback.
- PATCH absent/null/value semantics, contact ownership and verified flags,
  optimistic versions, conflict recovery, atomic Person and Contact writes.
- Search normalization, contact hydration, cursor and error mapping.
- Deactivate returns an already-deactivated Person without rewriting and
  re-reads after optimistic conflicts.

Single-step GetByGID uses a strategy directly rather than an artificial chain.
Operation-private helpers are not promoted to core and actually shared logic
remains under person/helpers.

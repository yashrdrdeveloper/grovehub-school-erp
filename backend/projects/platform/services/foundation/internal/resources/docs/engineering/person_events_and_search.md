# Foundation — Person Events and Search

This file contains **module-specific** details. Generic [events](../../../../../../../../../docs/engineering/backend/events.md) and [search](../../../../../../../../../docs/engineering/backend/searchDetails.md) rules remain shared.

## Changed event

- **Type:** `foundation.person.changed.v1`
- **Topic:** `foundation_person_changed`
- **Payload:** `{operation, old, new}` inside a generic event envelope.
- **Operations:** `CREATED`, `UPDATED`, `DEACTIVATED`.
- **Semantics:** Creation has `old = null`; updates and deactivations carry full old and new Person response snapshots. No computed diff is included.
- No event on a no-op update or idempotent replay. A deactivation emits only `DEACTIVATED`.
- There is no Person `DELETED` event today.

The event producer is under `internal/app/events/producers/person/`, with Changed pipeline and strategy owned by the Person aggregate. `EventStore.Enqueue` writes to `infrastructure.messaging` inside the Person transaction. Infrastructure asynchronously publishes the pending row; consumers must deduplicate by `eventId` and order by aggregate version if required. Full contact snapshots include PII, so consumers must be authorized and raw payloads must never be logged.

## Search projection

Person stores a rebuildable `search_details` text projection. The canonical Person and Contact tables are authoritative; the projection is for free-text matching only. Create and Update rebuild it using the aggregate in memory and persist it in the same transaction.

## Allowed search fields

Person scalar fields include `gid`, `firstName`, `middleName`, `lastName`, `status`, `createdAt` and `updatedAt`. Contact relation filters include `contacts.email`, `contacts.phone`, `contacts.primaryEmail`, `contacts.primaryPhone`, `contacts.verifiedEmail` and `contacts.verifiedPhone`.

The resource validator owns the exact operator/value allowlist. `repository/person/helpers/search/` supplies trusted table expressions and joins to the generic search capability. The repository accepts only a validated typed search query; the service hydrates related contacts in batch before returning a response. Opaque cursor pagination is mandatory; never accept direct SQL expressions or client-defined offsets.

## Verification

Test event snapshots and once-only emission, rollback without queued event, contact PII logging policy, projection consistency after contact replacement/clear, relation filters, sorting ties, invalid cursor, and repeated forward pagination. Use the live OpenAPI contract as the source for JSON shapes and response envelopes.

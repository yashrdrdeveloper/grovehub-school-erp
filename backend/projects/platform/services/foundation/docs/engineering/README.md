# Foundation — Engineering

This directory records **implemented module-specific architecture**. Generic rules belong to [Backend Engineering Standard](../../../../../docs/engineering/backend/backend.md), not here.

## Module responsibilities

Foundation owns the permanent Person GID and Person–Contact aggregate; it does not own login credentials, session tokens, or external OAuth provider connections. Data is persisted in PostgreSQL through GORM repositories, exposed over generated HTTP routes and a typed internal client. Its durable Person events are written through Infrastructure's shared transactional outbox.

## Read in order

- [architecture.md](architecture.md) — bootstrap, folders, data ownership, access boundary and migration.
- [person_events_and_search.md](person_events_and_search.md) — event schema, transaction, search projection and cursor specifics.
- [../product/README.md](../product/README.md) — business surface and execution scenarios.
- Canonical OpenAPI: `internal/resources/api/openapi/v1/openapi.yaml`.

When a folder, API, table, event or behavior changes, update the specific document in the same commit. The engineering docs record **what exists**; [generic API rules](../../../../../docs/engineering/backend/api.md), [persistence](../../../../../docs/engineering/backend/persistence.md) and [search](../../../../../docs/engineering/backend/searchDetails.md) remain authoritative across modules.

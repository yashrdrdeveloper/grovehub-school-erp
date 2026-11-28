# Infrastructure — Engineering

Infrastructure owns provider-neutral broker messaging, optional PostgreSQL transactional outbox, asynchronous delivery, retries, and background worker lifecycle.

Generic [Backend](../../../../../docs/engineering/backend/backend.md), [Events](../../../../../docs/engineering/backend/events.md), [Persistence](../../../../../docs/engineering/backend/persistence.md) and [Shared Libraries](../../../../../docs/engineering/backend/libraries.md) guidelines apply.

- [messaging_contracts.md](messaging_contracts.md) — public broker contracts, configuration and ownership.
- [architecture.md](architecture.md) — exact runtime, adapter, repository, worker and config ownership.
- [delivery_runbook.md](delivery_runbook.md) — retry, leases, cleanup, operational diagnostics.
- [../product/README.md](../product/README.md) — consumer-facing messaging workflows and expectations.

This is not an HTTP business module. There are intentionally no Infrastructure messaging HTTP controllers/routes.

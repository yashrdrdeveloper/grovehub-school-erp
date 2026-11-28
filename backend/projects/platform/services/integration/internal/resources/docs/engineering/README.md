# Integration — Engineering

Integration owns provider connection lifecycle, OAuth, provider capability discovery, connection tests, generic actions, and webhook registration. It does **not** own persistent message ledgers, mailbox synchronization cursors or user-facing communication workflows.

Read the generic [Backend Standard](../../../../../../../../../docs/engineering/backend/backend.md), [Provider Guide](../../../../../../../../../docs/engineering/backend/providers.md), [Events](../../../../../../../../../docs/engineering/backend/events.md), and [API](../../../../../../../../../docs/engineering/backend/api.md). This directory contains only **module-specific implementations**:

- [architecture.md](architecture.md) — exact resource, controller and operation wiring.
- [providers_and_events.md](providers_and_events.md) — Gmail/Outlook specialization, credentials, notification semantics and worker lifecycle.
- [../product/README.md](../product/README.md) — supported use cases and API execution.

OpenAPI: `internal/resources/api/openapi/v1/openapi.yaml`. The public `integration.New(...)` lifecycle remains stable; the module's `internal/app/config` owns provider configuration. Do not expose that private package through an import.

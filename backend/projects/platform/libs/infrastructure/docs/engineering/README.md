# Infrastructure Shared Library — Engineering

`libs/infrastructure` owns **technical contracts** that Infrastructure publishes to backend modules. It is distinct from domain-neutral `libs/core` and from the live broker/DB adapters owned by the Infrastructure service.

- Generic [library rules](../../../../../docs/engineering/backend/libraries.md).
- [messaging_contracts.md](messaging_contracts.md) — exact shared messaging structures, configuration and ownership.
- Consuming module interface: `services/infrastructure/client/messaging/messaging.go`.
- Owner's runtime: `services/infrastructure/docs/engineering/architecture.md`.

This library intentionally has no `docs/product/`; the runtime owner documents consumer-visible workflows. Never import a vendor SDK from this contract package.

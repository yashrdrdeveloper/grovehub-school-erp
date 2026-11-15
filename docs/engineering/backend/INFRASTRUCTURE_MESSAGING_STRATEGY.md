# GroveHub Infrastructure Messaging Strategy

> **Status:** Mandatory GroveHub backend infrastructure messaging standard  
> **Applies to:** `backend/services/infrastructure` messaging capabilities and provider implementations

## 1. Consumer Boundary

Other GroveHub services consume only the public Infrastructure client:

```text
consumer service
    ↓
infrastructure/client.MessagingClient
    ↓
infrastructure/client/impl
    ↓
infrastructure/services.MessagingService
```

Consumers publish/consume logical topics through the same public `MessagingClient` (`Publish`, `PublishWithOptions` and `Consume`). Event listeners implement domain handlers; an additional HTTP client or separate consumer transport client is not required. They never select NATS, Kafka, RabbitMQ, or another provider directly. Domain event meaning/contracts and core JSON diffing are outside Infrastructure; see `EVENT_GUIDELINES.md`.

## 2. Strategy Resolution

Messaging provider selection uses a generic strategy registry:

```text
MessagingService
    ↓
Registry[MessagingStrategy]
    ↓
configured provider key
    ↓
MessagingStrategy
```

Do not grow provider `if/else` or `switch` trees inside business/service methods. New providers register another strategy during module composition.

## 3. Provider-Specific Placement

Provider-specific SDK/API mechanics live below:

```text
backend/services/infrastructure/externalService/<provider>/
```

For NATS:

```text
externalService/nats/adapters/
```

Generic Infrastructure service/client contracts must not expose NATS or Watermill types.

## 4. Current NATS Boundary

The first implementation uses Watermill with NATS JetStream:

```text
MessagingClient
    ↓
MessagingService
    ↓
NATS MessagingStrategy
    ↓
externalService/nats/adapters
    ↓
Watermill
    ↓
NATS JetStream
```

JetStream owns durable broker message/consumer state. The separate `infrastructure.messaging` PostgreSQL outbox is used only for messages explicitly requiring atomic persistence before publishing; it is not a replacement for JetStream's own storage.

## 5. Per-message Persistence

`Publish(ctx, topic, payload, metadata)` remains a direct NATS publish for existing callers.

`PublishWithOptions(ctx, topic, payload, metadata, messagingclient.PublishOptions{Persist: true})` inserts a `PENDING` row in `infrastructure.messaging`. The existing dispatcher publishes it to NATS asynchronously and retries failures; sent rows are pruned after 24 hours. An absent or false `Persist` sends directly to NATS without a PostgreSQL outbox row.

The caller may pass a stable `metadata["event-id"]`; otherwise Infrastructure generates a UUID and includes it in the stored metadata. Consumers must deduplicate at-least-once deliveries by the event ID. The original `Publish` method is unchanged for compatibility.

**Atomicity:** Call `PublishWithOptions(..., Persist: true)` *inside the same `coretransaction.Run` / `UnitOfWork.Run` context* as the business database write. The outbox and business tables must be in the same PostgreSQL database. Calling outside that transaction still queues a durable message but cannot atomically commit it with an earlier business write. `Persist: false` should be used after commit if it describes a committed state change.

Foundation Person already uses the existing `EventStore.Enqueue` in its UnitOfWork; do not additionally call `PublishWithOptions` for the same event or it would be queued twice.

## 6. Trace Context

Publish propagates the active W3C trace context through message metadata and includes `X-Trace-ID` for GroveHub support correlation. Consumers extract W3C context before invoking handlers so downstream work continues the same distributed trace.

## 7. Hard Rules

```text
YES consumers depend only on infrastructure/client
YES provider selection is registry/strategy driven
YES provider-specific code stays in externalService/<provider>
YES context.Context crosses publish/consume boundaries
YES durable broker delivery state belongs to JetStream
YES critical business events opt into the transactional outbox inside the business transaction

NO NATS/Watermill types in public client/service contracts
NO provider switch trees in service methods
NO business service importing externalService/nats
NO outbox row for a direct (Persist: false) publish
NO assumption that an outbox insert outside UnitOfWork is atomic with a business write
```

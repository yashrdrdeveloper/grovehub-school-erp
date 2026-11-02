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

Consumers publish/consume logical topics. They never select NATS, Kafka, RabbitMQ, or another provider directly.

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

JetStream owns durable broker message/consumer state. A GroveHub database table is not required merely to register events or persist broker delivery state.

## 5. Trace Context

Publish propagates the active W3C trace context through message metadata and includes `X-Trace-ID` for GroveHub support correlation. Consumers extract W3C context before invoking handlers so downstream work continues the same distributed trace.

## 6. Hard Rules

```text
YES consumers depend only on infrastructure/client
YES provider selection is registry/strategy driven
YES provider-specific code stays in externalService/<provider>
YES context.Context crosses publish/consume boundaries
YES durable messaging state belongs to JetStream

NO NATS/Watermill types in public client/service contracts
NO provider switch trees in service methods
NO business service importing externalService/nats
NO database table solely for event registration/delivery state
```

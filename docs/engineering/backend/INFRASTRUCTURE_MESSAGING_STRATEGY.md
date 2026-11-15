# GroveHub Infrastructure Messaging Strategy

> **Applies to:** `backend/services/infrastructure`.

Infrastructure is the owner of transport, provider selection, outbox delivery,
and retry mechanics. Business modules own event contracts and snapshots.

## Client, service and provider boundary

```text
Foundation / Identity / Integration
        ↓
infrastructure/client/messaging.MessagingClient
        ↓  generated @Proxy
internal/app/client/messaging/impl
        ↓
internal/app/services/messaging/impl
        ├─ operations/publish
        │  └─ strategies/impl + helpers  ─┐
        ├─ operations/consume              ├─ provider-neutral Strategy
        └─ operations/health               │       ↓
                                           └─ externalService/messaging/nats/adapters
                                                      ↓
                                              Watermill / NATS JetStream
```

`bootstrap/impl/messaging.go` registers the concrete NATS adapter directly
against the provider-neutral `services/messaging/strategies.Strategy`
interface. There must be **no** `services/messaging/strategies/nats.go`:
the NATS-specific implementation belongs entirely under
`externalService/messaging/nats/`.

The generic registry does not branch on vendor names inside business
or messaging operations.

## Config and bootstrap

Configuration lives at `internal/app/config/messaging/messaging.go` with a
stable public alias at `client/messaging/config.go`. The Foundation-style
bootstrap builder owns `bootstrap.go`, `messaging.go`, `repository.go`,
`client.go` and `events.go`. Infrastructure exposes `New` and `Runtime`;
the dispatcher is started and stopped by the runtime.

## Per-event persistence

```go
// Direct to NATS, no PostgreSQL outbox row.
messaging.Publish(ctx, topic, payload, metadata)

// Explicitly non-persistent, also direct to NATS.
messaging.PublishWithOptions(ctx, topic, payload, metadata,
    messagingclient.PublishOptions{Persist:false})

// Enqueue into infrastructure.messaging; publish asynchronously after commit.
messaging.PublishWithOptions(ctx, topic, payload, metadata,
    messagingclient.PublishOptions{Persist:true})
```

Use `Persist:true` *within the business `UnitOfWork`* if the event and
entity mutation need to commit atomically. An outbox insert in an independent
transaction is not atomic with a previously committed business write.

An explicit `metadata["event-id"]` is preserved. Otherwise Infrastructure
generates a UUID. Consumers must deduplicate by the event ID because a
successful publish followed by a failed `MarkSent` can cause delivery again.

Foundation already calls `EventStore.Enqueue` from its Person event pipeline.
Do not double-publish by calling `PublishWithOptions` for the same change.

## Retention and trace propagation

JetStream owns broker storage, consumer acknowledgements and replay.
The optional PostgreSQL outbox is only delivery coordination. Successfully
published `SENT` records are pruned after 24 hours. `PENDING` and
retryable `SENDING` records are not deleted by retention cleanup.

W3C trace context and `X-Trace-ID` are handled by the NATS adapter.
Never put raw Person payloads into logs or expose NATS/Watermill types
through public clients.

## Hard rules

- Business services import only the public Infrastructure client contracts.
- Infrastructure services remain vendor-neutral and operation-oriented.
- Concrete providers live under `externalService/<capability>/<provider>/`.
- Do not add a vendor-specific service strategy file.
- Never assume `Persist:false` guarantees delivery across a process crash.
- Never use Infrastructure's transactional outbox as permanent end-user audit history.

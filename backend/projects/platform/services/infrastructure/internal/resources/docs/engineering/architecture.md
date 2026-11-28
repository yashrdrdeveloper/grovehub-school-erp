# Infrastructure — Messaging Architecture

## 1. Module layout and ownership

```text
services/infrastructure/
├── infrastructure.go                      # New(config,db,logger), Migrate
├── client/messaging/messaging.go           # @Client + public contracts + EventStore
└── internal/
    ├── app/
    │   ├── bootstrap/impl/
    │   │   ├── bootstrap.go               # Build/Close
    │   │   ├── client.go                  # @Proxy client assembly
    │   │   ├── service.go                 # provider composition
    │   │   ├── repository.go              # GORM store + validator + mapper
    │   │   └── events.go                  # dispatcher worker lifecycle
    │   ├── config/messaging/messaging.go
    │   ├── dto/messaging/
    │   ├── entity/messaging.go
    │   ├── repository/messaging/
    │   ├── validators/messaging/
    │   ├── mappers/messaging/
    │   ├── services/messaging/
    │   │   ├── impl/messaging.go
    │   │   ├── strategies/                # provider-neutral registry
    │   │   └── operations/{publish,consume,health,dispatch}/
    │   └── externalService/messaging/nats/
    │       ├── nats.go                    # external service contract
    │       └── adapters/messaging.go      # NATS/Watermill implementation
    └── resources/migrations/
        └── 000001_create_messaging.sql
```

Infrastructure owns transport DTOs and client contracts in `services/infrastructure/client/messaging`. Configuration lives in `internal/app/config/messaging` and subscription context helpers in `internal/app/client/messaging/helpers`, exposed through the public client. Business modules import only the public client; no separate platform messaging library exists. NATS/Watermill types do not escape the adapter.

## 2. Runtime configuration

The public config is loaded from environment:

```dotenv
GROVEHUB_MESSAGING_PROVIDER=nats
GROVEHUB_NATS_URL=nats://127.0.0.1:4222
GROVEHUB_NATS_QUEUE_GROUP_PREFIX=grovehub
GROVEHUB_NATS_DURABLE_PREFIX=grovehub
```

These values are **local examples**, not Go production defaults. The development scripts can supply fallbacks; direct production execution must explicitly configure broker endpoints and validate all required fields. Names are configurable for deployment isolation.

## 3. Publish operations

- `Publish(ctx, topic, payload, metadata)` publishes directly through a registered strategy.
- `PublishWithOptions(..., PublishOptions{Persist:false})` also publishes directly and does not write PostgreSQL.
- `PublishWithOptions(..., PublishOptions{Persist:true})` writes a `PENDING` event to `infrastructure.messaging`, to be delivered later by the dispatcher.
- `EventStore.Enqueue` is the explicit transaction-aware outbox API for domain event producers.
- The same logical event must never be sent by both `Enqueue` and `PublishWithOptions(Persist:true)`.

Atomicity with business records exists only when `Enqueue` runs using the **same PostgreSQL transaction context** as business writes.

## 4. Consume and health

`MessagingClient.Consume` binds a caller-owned typed handler to a logical topic. Infrastructure owns provider subscription mechanics, but the consuming business module owns event schema validation and idempotent mutation. `Health` tests broker availability; it does not prove that every consumer has processed every event.

## 5. Dispatcher lifecycle and semantics

The dispatcher is started by `bootstrap/impl/events.go` and stopped/joined on runtime `Close`. It claims PostgreSQL outbox messages with `FOR UPDATE SKIP LOCKED`, temporary leases, bounded batches and retries. Publication to NATS and DB `MarkSent` are not a single transaction; consumers must deduplicate redelivery. Successfully published `SENT` rows are pruned after 24 hours. Unpublished events remain for retry/incident recovery.

No HTTP controllers or routes are created because the typed client and internal worker provide the API. Adding HTTP endpoints would require a genuine product requirement and explicit security model.

## 6. Migration and tests

Infrastructure owns the live messaging schema migration and any legacy data transfer. Applied SQL is immutable. Tests must verify both direct and persist modes, context-bound enqueue, retries, leases, prune safety, idempotent consumer expectations, shutdown cancellation, and broker connectivity.

## Test resource layout

Service-owned test directories live under `internal/resources/tests/`, including unit and integration suites, helpers and mocks. Generated suites mirror the resource paths under `.generated/`.

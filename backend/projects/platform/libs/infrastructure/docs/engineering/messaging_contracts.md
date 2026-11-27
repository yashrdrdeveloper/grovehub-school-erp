# Infrastructure Messaging — Shared Contract

## Public capability

`libs/infrastructure/messaging` defines:

| Type/file | Purpose |
|---|---|
| `Event` in `event.go` | Pending outbox event ID, topic, payload, metadata |
| `Message` and `Handler` in `consumer.go` | Transport-neutral broker delivery DTO and callback |
| `PublishOptions` in `producer.go` | Per-message database persistence choice |
| `Publisher` in `producer.go` | Minimal typed direct-publish contract |
| `Config`, `FromEnv`, `Validate` in `config.go` | Required provider/URL/queue/durable configuration |

The consuming application normally imports `services/infrastructure/client/messaging`, a typed public facade over these definitions, rather than depending on the runtime's `internal/app` packages. The Infrastructure service owns concrete NATS and PostgreSQL implementation, not the shared contract.

## Delivery choice

`Persist:false` means direct publish without a PostgreSQL outbox row; `Persist:true` means the Infrastructure service queues an outbox event. **Atomicity exists only when the call participates in the business transaction.** A pending event and a broker ACK are not the same guarantee; consumers must handle duplicate delivery.

## Configuration

`FromEnv` reads `GROVEHUB_MESSAGING_PROVIDER`, `GROVEHUB_NATS_URL`, `GROVEHUB_NATS_QUEUE_GROUP_PREFIX` and `GROVEHUB_NATS_DURABLE_PREFIX`. No fixed production endpoint is embedded in the library. Validate before building runtime dependencies; local development may set explicit values in its startup script.

## SRP restrictions

- Shared events are transport-neutral and contain opaque payload bytes.
- This package cannot define specific domain event types, names or business status transitions.
- It cannot import NATS, Watermill, GORM repository internals or a business service.
- Never place audit event retention or broker worker scheduling in the DTO library.
- Every public type change is a cross-module contract change; update consumers and tests.

## Verification

Test env parsing, required-field validation, generated event ID propagation at the runtime boundary, defensive copies of payload/metadata, direct/persistent behavior, and compile-time compatibility with public client aliases.

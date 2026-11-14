# Infrastructure service — resource-based internal layout

Infrastructure owns NATS transport, event delivery, SQL persistence and the
outbox lifecycle. Business services must not import NATS or a database
repository directly.

```text
backend/services/infrastructure/
├── infrastructure.go
├── client/messaging/
│   ├── messaging.go             # @Client broker facade
│   ├── outbox.go                # public transactional EventStore contract
│   └── config.go                # public config alias
└── internal/
    ├── app/
    │   ├── bootstrap/{bootstrap.go,impl/bootstrap.go}
    │   ├── client/messaging/impl/messaging.go
    │   ├── config/messaging.go
    │   ├── constants/error/messaging.go
    │   ├── dto/messaging/{request,response}.go
    │   ├── mappers/messaging/{messaging.go,impl/messaging.go}
    │   ├── entity/messaging.go  # @Entity infrastructure.messaging
    │   ├── externalService/nats/{nats.go,adapters/messaging.go}
    │   ├── repository/messaging/{messaging.go,impl/messaging.go}
    │   ├── services/messaging/{messaging.go,dispatcher.go,impl/,strategies/}
    │   └── validators/messaging/{messaging.go,impl/messaging.go}
    └── resources/
        ├── resources.go
        └── migrations/{migrate.go,000001_create_messaging.sql}
```

## Transactional guarantees and schema ownership

The transactional event outbox is an Infrastructure concern, not a reusable
SQL helper in `libs/core`. Foundation depends only on
`client/messaging.EventStore`; the same PostgreSQL GORM transaction context
is supplied to Infrastructure's `repository/messaging` implementation.
The Infrastructure schema is in the same PostgreSQL database. Moving the outbox
to an independent database would break atomicity with Person writes.

`entity.Messaging` stores one pending event delivery, not a chat conversation.
Its statuses are PENDING, SENDING and SENT. The dispatcher claims records with
`FOR UPDATE SKIP LOCKED`, publishes to JetStream, then marks the record SENT;
failures are retried with exponential backoff. A bounded cleanup removes only
SENT records older than 24 hours. Pending deliveries are never deleted.

Run `make db-schema-init` to apply Foundation, Infrastructure, Identity and
Integration migrations in order. Production: `groovehub-platform migrate all`
applies the embedded Foundation and Infrastructure migrations. The historical
Foundation outbox migration is retained to preserve Goose version tracking.
Infrastructure's first migration copies legacy records and drops the legacy
table **only after the copy succeeds**. Stop any old publisher before migrating.

## Broker client contracts

`client/messaging.MessagingClient` exposes Publish, Consume and Health,
delegated by generated `@Client` proxies to the private messaging service.
`@ExternalService` denotes the NATS provider's interface and is backed by a
compile-time adapter implementation assertion. NATS subscription handlers are
Go callbacks and cannot be encoded as REST `@RequestLine` operations without
introducing a separate HTTP controller and request DTO. We retain direct
typed delegation for broker callbacks rather than inventing HTTP semantics.

No new unit-test files are introduced by this structural refactor.

The Messaging mapper converts public event commands into private durable
entities; repository implementations are responsible only for persistence,
leasing, retries and pruning. DTOs are resource-oriented `request.go` and
`response.go`, consistent with Foundation's DTO layout.

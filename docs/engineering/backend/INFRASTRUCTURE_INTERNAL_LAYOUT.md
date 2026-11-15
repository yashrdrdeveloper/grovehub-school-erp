# Infrastructure service — resource-based internal layout

Infrastructure follows the same resource-based composition used by Foundation.
All public messaging contracts and DTOs are exposed via
`libs/infrastructure/messaging`, not `libs/core/events/messaging`.
Replace the three separate public client files with one `client/messaging/messaging.go`.

Messaging is a capability; NATS is only one implementation and appears **only**
under the Messaging external-service provider folder.

```text
backend/services/infrastructure/
├── infrastructure.go
├── client/messaging/messaging.go      # @Client, EventStore, Config, DTO aliases
└── internal/
    ├── app/
    │   ├── bootstrap/
    │   │   ├── bootstrap.go          # Bootstrap / Runtime interfaces
    │   │   └── impl/
    │   │       ├── bootstrap.go      # builder.Build, runtime.Close
    │   │       ├── client.go         # generated messaging client composition
    │   │       ├── service.go        # provider registration and service composition
    │   │       ├── repository.go     # repository, mapper and validator assembly
    │   │       └── events.go         # background dispatcher lifecycle
    │   ├── client/messaging/impl/messaging.go
    │   ├── config/messaging/messaging.go
    │   ├── constants/error/messaging.go
    │   ├── dto/messaging/{request,response}.go
    │   ├── entity/messaging.go
    │   ├── mappers/messaging/{messaging.go,impl/messaging.go}
    │   ├── repository/messaging/{messaging.go,impl/messaging.go}
    │   ├── validators/messaging/{messaging.go,impl/messaging.go}
    │   ├── externalService/messaging/
    │   │   └── nats/
    │   │       ├── nats.go          # @ExternalService contract
    │   │       └── adapters/messaging.go
    │   └── services/messaging/
    │       ├── messaging.go         # service contract
    │       ├── impl/messaging.go    # @Service composed operations
    │       ├── strategies/
    │       │   ├── strategy.go      # provider-neutral contract
    │       │   └── registry.go      # generic provider registry
    │       └── operations/
    │           ├── publish/{publish.go,helpers/,strategies/}
    │           ├── consume/{consume.go,strategies/}
    │           ├── health/{health.go,strategies/}
    │           └── dispatch/dispatcher.go
    └── resources/migrations/
        ├── migrate.go
        └── 000001_create_messaging.sql
```

There are no Infrastructure HTTP controllers or routes: messaging uses a typed
`@Client` whose generated `@Proxy` methods delegate into the private service.
This is the same `client/<resource>/impl` pattern as Foundation, but Foundation
uses generated in-process HTTP `@RequestLine` calls for its actual HTTP routes.

### Event persistence

- `Publish` or `PublishWithOptions(..., PublishOptions{Persist:false})` publishes
  directly to the provider without an outbox row.
- `PublishWithOptions(..., PublishOptions{Persist:true})` enqueues a `PENDING`
  message in `infrastructure.messaging`; the asynchronous worker publishes it.
- Existing Foundation Person producers use `EventStore.Enqueue` inside their
  UnitOfWork and are unchanged. Do not enqueue the same event twice.
- `Persist:true` is atomic with business writes **only when called inside their
  existing GORM transaction** using the same PostgreSQL database.
- The worker claims records with `FOR UPDATE SKIP LOCKED`, retries failures,
  reclaims expired leases, and removes only `SENT` records older than 24 hours.
  This is delivery state, not permanent user-visible audit history.

NATS-specific types and Watermill APIs cannot appear in the generic
`services/messaging`, `client/messaging`, or business modules. Add additional
providers under `externalService/messaging/<provider>/`; register them in
`bootstrap/impl/service.go` without creating vendor-named service strategies.

### Environment-driven broker configuration

`libs/infrastructure/messaging.FromEnv` reads
`GROVEHUB_MESSAGING_PROVIDER`, `GROVEHUB_NATS_URL`,
`GROVEHUB_NATS_QUEUE_GROUP_PREFIX` and
`GROVEHUB_NATS_DURABLE_PREFIX`. Application startup validates these values.
There is no hardcoded NATS endpoint or queue/durable name in Go.
`scripts/lib/env.sh` supplies developer defaults, and `.env.example`
documents the settings for deployments and direct binary execution.

### Why there is no Messaging HTTP controller

The public `@Client` interface is the application boundary for Publish,
PublishWithOptions, Consume and Health; its generated proxy forwards to the
operation-based Messaging service, which delegates to the selected
`externalService/messaging/<provider>` adapter. There is no
HTTP endpoint to register, so a controller or routes folder would have no
actual work. Domain event consumer handlers belong to their owning services.

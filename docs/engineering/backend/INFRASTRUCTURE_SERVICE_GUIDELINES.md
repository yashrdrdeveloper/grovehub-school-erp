# GroveHub Infrastructure Service Guidelines

> **Status:** Mandatory GroveHub backend infrastructure ownership standard  
> **Applies to:** integrations with SMTP/mail engines, cache engines, object storage, queues, search engines, external infrastructure APIs, and similar technical infrastructure capabilities  
> **Owner:** `backend/services/infrastructure`

---

## 1. Single Infrastructure Ownership Boundary

All GroveHub backend integration code for infrastructure capabilities lives under:

```text
backend/services/infrastructure/
```

Business services must not create provider-specific infrastructure adapters inside their own service roots.

The business service owns the business use case. Infrastructure owns the technical integration.

---

## 2. Canonical Infrastructure Structure

Resource-based modules follow Foundation's internal/app layout. See
`INFRASTRUCTURE_INTERNAL_LAYOUT.md` for the current Messaging tree.

```text
internal/app/
  bootstrap/{bootstrap.go,impl/{bootstrap,client,messaging,repository,events}.go}
  client/<capability>/impl/
  config/<capability>/
  entity/
  dto/<resource>/
  repository/<resource>/impl/
  services/<capability>/{impl,operations/<operation>/{helpers,strategies/impl}}
  externalService/<capability>/<provider>/adapters/
```

Only create folders with real runtime responsibility; an event transport
without HTTP endpoints should not invent `routes/` or `controller/`.

## 3. Client Rule

Consumers depend only on the public Infrastructure client contract.

```text
Consuming service
    ↓
infrastructure/client
    ↓
infrastructure/internal/app/client/<capability>/impl
    ↓
infrastructure/services interface
    ↓
infrastructure/services/impl
```

Client implementations are intentionally thin. They forward `context.Context` and delegate to the owning service. They must not duplicate provider behavior, retry policy, configuration resolution, or external API logic.

Business services must not import Infrastructure implementation, strategy, or external-service packages.

---

## 4. Service Rule

`internal/app/services/<capability>/<capability>.go` defines the capability service contract.

`internal/app/services/<capability>/impl/<capability>.go` composes operation-specific strategies; `bootstrap/impl/<capability>.go` resolves and registers providers.

`services/<capability>/operations/<operation>/helpers/` holds pure operation-specific helpers only.

Repositories remain persistence-only. Do not move provider selection, business decisions, or external API calls into repositories.

---

## 5. Strategy / Registry Rule

Provider selection must use a registry/strategy pattern instead of growing vendor `if/else` or `switch` trees.

```text
Capability Service
    ↓
Registry[CapabilityStrategy]
    ↓
configured provider key
    ↓
Capability Strategy
    ↓
External Service Adapter
```

A generic registry may be reused across capabilities:

```text
Registry[T]
    Register(key, strategy)
    Resolve(key)
```

Each capability still owns a strongly typed strategy interface. Do not force unrelated capabilities into one giant generic `Execute()` contract.

Example:

```text
MessagingService
    ↓
Registry[MessagingStrategy]
    ↓
externalService/messaging/nats/adapters (implements provider Strategy)
```

Adding another messaging provider should normally register another strategy during module composition without changing `MessagingService` provider branching logic.

---

## 6. External Service Rule

Concrete vendor/API mechanics live under:

```text
backend/services/infrastructure/internal/app/externalService/<capability>/<provider>/
```

Examples:

```text
externalService/messaging/nats/adapters/
externalService/mail/gmail/adapters/
externalService/mail/gmail/webhooks/
```

Provider-specific SDK types, request mapping, protocol behavior, webhook verification, and external API calls remain inside the external-service boundary.

Generic service/client contracts must not expose vendor SDK types.

---

## 7. Provider-Agnostic Naming

Public contracts use capability names, not implementation/vendor names.

Use:

```text
MessagingClient
MailClient
StorageClient
CacheClient
```

Avoid public contracts such as:

```text
NATSClient
GmailClient
RedisClient
S3Client
```

Consumers request the capability. Infrastructure resolves the configured provider internally.

---

## 8. Messaging Ownership

The first messaging implementation uses Watermill with NATS JetStream.

```text
Consumer Service
    ↓
MessagingClient
    ↓
MessagingService
    ↓
MessagingStrategy registry
    ↓
externalService/messaging/nats/adapters
    ↓
Watermill
    ↓
NATS JetStream
```

JetStream owns durable broker message and consumer state. A GroveHub database table is not required solely for broker registration. The existing `infrastructure.messaging` outbox records only explicitly persisted delivery messages.

The `infrastructure.messaging` transactional outbox already exists. It is used when `PublishWithOptions(..., PublishOptions{Persist:true})` is called within a business UnitOfWork; direct publish writes no outbox row.

---

## 9. Trace Context Across Messaging

Messaging publish/consume boundaries must preserve GroveHub trace context.

Publishing propagates W3C trace context in message metadata and may include `X-Trace-ID` for support correlation.

Consumers extract the W3C parent before invoking the handler so downstream work continues the same distributed trace.

Provider-specific adapters must not create disconnected traces for normal request-derived message flows.

---

## 10. Durable Broker State

NATS JetStream persistence must use persistent storage outside the disposable application container/process.

Application restarts must not remove broker state.

Server migration/rebuild must preserve or restore JetStream storage, or use a clustered deployment when higher availability is required.

Do not duplicate every event into PostgreSQL. Only persist events that require transactional delivery; retain user-facing audit records separately.

---

## 11. Hard Rules

```text
YES all backend infrastructure integration code is owned by backend/services/infrastructure
YES consumers use infrastructure/client contracts
YES client implementations delegate to infrastructure services
YES provider selection uses typed strategy + registry resolution
YES provider-specific API/SDK behavior stays under externalService/<capability>/<provider>
YES context.Context crosses messaging/provider boundaries
YES JetStream owns durable broker delivery state

NO vendor switch trees in business/service methods
NO provider-specific client exposed to consuming business services
NO business service importing infrastructure impl/strategy/externalService packages
NO database/network/filesystem work inside services/helpers
NO vendor SDK types in generic client/service contracts
NO outbox record for a non-persistent Publish
```

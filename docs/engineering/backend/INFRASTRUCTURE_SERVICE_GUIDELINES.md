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

Infrastructure follows the normal GroveHub flat service structure and adds two explicit infrastructure responsibilities when needed:

```text
backend/services/infrastructure/
├── infrastructure.go
├── config/
├── constants/
│   └── error/
├── entity/
├── repository/
│   └── impl/
├── services/
│   ├── <capability>.go
│   ├── impl/
│   │   └── <capability>.go
│   ├── helpers/
│   │   └── <capability>.go
│   └── strategy/
│       ├── registry.go
│       └── <capability>/
│           ├── strategy.go
│           └── <provider>.go
├── externalService/
│   └── <provider>/
│       ├── adapters/
│       └── webhooks/
├── client/
│   ├── <capability>.go
│   └── impl/
│       └── <capability>.go
├── controller/
│   └── impl/
├── routes/
└── migrations/
```

Create only folders that have a current responsibility.

---

## 3. Client Rule

Consumers depend only on the public Infrastructure client contract.

```text
Consuming service
    ↓
infrastructure/client
    ↓
infrastructure/client/impl
    ↓
infrastructure/services interface
    ↓
infrastructure/services/impl
```

Client implementations are intentionally thin. They forward `context.Context` and delegate to the owning service. They must not duplicate provider behavior, retry policy, configuration resolution, or external API logic.

Business services must not import Infrastructure implementation, strategy, or external-service packages.

---

## 4. Service Rule

`services/<capability>.go` defines the capability service contract.

`services/impl/<capability>.go` owns technical use-case orchestration and resolves the configured provider strategy.

`services/helpers/<capability>.go` contains pure capability-specific service helpers only.

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
NATS strategy
```

Adding another messaging provider should normally register another strategy during module composition without changing `MessagingService` provider branching logic.

---

## 6. External Service Rule

Concrete vendor/API mechanics live under:

```text
backend/services/infrastructure/externalService/<provider>/
```

Examples:

```text
externalService/nats/adapters/
externalService/gmail/adapters/
externalService/gmail/webhooks/
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
NATS strategy
    ↓
externalService/nats/adapters
    ↓
Watermill
    ↓
NATS JetStream
```

JetStream owns durable broker message and consumer state. A GroveHub database table is not required solely for event registration or broker delivery state.

A database may be introduced later only when GroveHub needs application-owned dynamic configuration such as tenant-specific runtime subscriptions, admin-managed bindings, or auditable runtime registration.

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

GroveHub application code must not duplicate all event payloads into PostgreSQL merely to compensate for ephemeral broker storage.

---

## 11. Hard Rules

```text
YES all backend infrastructure integration code is owned by backend/services/infrastructure
YES consumers use infrastructure/client contracts
YES client implementations delegate to infrastructure services
YES provider selection uses typed strategy + registry resolution
YES provider-specific API/SDK behavior stays under externalService/<provider>
YES context.Context crosses messaging/provider boundaries
YES JetStream owns durable broker delivery state

NO vendor switch trees in business/service methods
NO provider-specific client exposed to consuming business services
NO business service importing infrastructure impl/strategy/externalService packages
NO database/network/filesystem work inside services/helpers
NO vendor SDK types in generic client/service contracts
NO DB table solely for broker event registration/delivery state
```

# GroveHub Event Guidelines

> **Status:** Mandatory GroveHub backend standard  
> **Applies to:** all domain events published or consumed by backend services  
> **Purpose:** define event ownership, package placement, event contracts, publishing boundaries, and messaging integration.

---

## 1. Ownership

A domain service owns the events produced from changes to data/business state that it owns.

Examples:

```text
Foundation owns Person/GID events.
Identity owns Account/Credential/Session events.
Infrastructure owns transport mechanics, not business event meaning.
```

Do not place Foundation or Identity event contracts inside `services/infrastructure`.

Infrastructure provides the generic `MessagingClient`; domain services decide what an event means and when it must be emitted.

---

## 2. Canonical Event Structure

Aggregate-owned events use an event contract, pipelines, producers, and optional consumers.

```text
backend/
├── libs/core/events/
│   ├── envelope.go                     # Typed domain-independent envelope
│   ├── pipeline/rules.go                # Generic ShouldEmit/Build/Enqueue rules
│   ├── outbox/
│   │   ├── store.go                     # Transactional enqueue and leased claims
│   │   └── dispatcher.go                # Retry and broker delivery
│   └── messaging/
│       ├── producer.go                  # Transport-neutral publisher interface
│       └── consumer.go                  # Transport-neutral handler contract
└── services/foundation/events/person/
    ├── person.go                        # @Event Changed payload and schema metadata
    ├── pipelines/
    │   ├── changed.go                   # Typed event rule composition
    │   ├── flows/changed.go             # Change input state
    │   └── strategies/
    │       ├── strategy.go              # @Strategy contract
    │       └── impl/changed.go          # Eligibility and payload building strategy
    ├── producers/
    │   ├── producer.go                  # ChangedProducer domain contract
    │   ├── noop.go                      # Optional test-only producer
    │   └── impl/person.go               # Enqueue through ChangedPipeline
    ├── consumers/
    │   ├── consumer.go                  # Optional subscription contract
    │   └── impl/person.go               # MessagingClient.Consume integration
    ├── handlers/changed.go              # Decode/validate inbound envelope
    └── helpers/
        ├── changed.go                   # HasChanged eligibility
        ├── envelope.go                  # BuildChanged
        ├── message.go                   # Marshal outbox payload
        └── metadata.go                  # Broker headers
```

The lifecycle `Operation` type remains at `foundation/constants/events/operation.go`.
No consumer is registered automatically in Foundation: producer-side self-subscription
would generate unnecessary load and risk loops. Consumers are explicitly registered
by services that need incoming Person Changed events.


## 3. Responsibility Split

- `events/person/person.go`: stable domain event payload and annotation metadata.
- `events/person/pipelines/`: typed event policy; `Rules()` composes `ShouldEmit`,
  `Build`, and `Enqueue` as in the validator rule pattern.
- `events/person/producers/`: domain-level producer contract; enqueue must happen
  inside the successful aggregate write transaction, not after returning from the service.
- `libs/core/events/outbox/`: durable storage, nonblocking row claims, leases,
  retry scheduling, broker delivery, and published-state tracking.
- `events/person/consumers/`: optional registration with Infrastructure MessagingClient.
  `handlers/` decodes and validates the incoming typed envelope. Consumers must
  deduplicate event IDs and protect contact details.
- `services/person/operations/`: owns actual business state transitions and invokes
  the producer only for newly committed Create, changed Update, or new Deactivate.
  GetByGID and Search never enqueue events.

Event construction helpers must remain pure. They do not call repositories,
NATS adapters or HTTP. Infrastructure owns transport mechanics.


## 4. Event Envelope

All service events use the shared core envelope from:

```text
backend/libs/core/events
```

Canonical JSON envelope, implemented by `coreevents.Envelope[T any]` and `coreevents.NewEnvelope(ctx, metadata, payload)`:

```json
{
  "eventId": "evt_1001",
  "eventType": "identity.account.created.v1",
  "version": 1,
  "aggregateId": "account_123",
  "occurredAt": "2026-11-06T10:00:00Z",
  "correlationId": "request-or-trace-id",
  "causationId": "command-or-prior-event-id",
  "data": {}
}
```

- `eventId` identifies this unique event; `aggregateId` identifies the Person GID or Account ID that changed.
- `version` is the **event schema version** (currently `1`), not the entity record version; `data.version` may separately hold the record's version.
- `correlationId` identifies the end-to-end request/workflow. Use explicitly propagated context when available, otherwise the active OpenTelemetry trace ID.
- `causationId` identifies the originating command or prior event, never a freshly invented event ID. Registration supplies its validated idempotency key as the command ID.
- When correlation or causation is genuinely unavailable (for example, an internal background job with no trace or command context), the field is omitted instead of inventing a misleading value.
- `source` is **not** part of the envelope; it stays in broker message metadata as a transport detail.
- Keep existing domain event-type strings and stable NATS topics when migrating consumers; adding the top-level properties changes the JSON contract.
- The domain payload belongs to its owning service package. Only `Envelope[T]`, `Metadata`, context propagation, and generic construction live in `libs/core/events`.

---

## 5. Naming

Topics/subjects are stable and must not be randomized per message or per test run.

Examples:

```text
foundation_person_changed
identity_account_created
```

Event types are versioned:

```text
foundation.person.changed.v1
identity.account.created.v1
```

Message/event IDs are unique; topics remain stable.

Do not create a new JetStream subject/stream name just to make one event instance unique.

---

## 6. Create Semantics and Idempotency

When a service returns `(resource, created, error)`, publish the `created` event only when:

```text
error == nil
created == true
```

An idempotent replay that returns an already-existing resource with `created == false` must not emit a second created event.

This is required for Foundation Person creation and Identity Account creation.

---

## 7. Event Payload Policy and Foundation Person exception

Default to minimal, non-sensitive payloads unless a concrete domain contract requires otherwise. Never include credentials, tokens, secrets or unrelated personal data.

**Current Foundation Person contract is an explicit full-response-snapshot exception**, adopted to support consumer-side old/new comparison:

- Topic: `foundation_person_changed`; event type: `foundation.person.changed.v1`.
- Operation type and constants: `foundation/constants/events/operation.go`; JSON values remain unchanged.
- Envelope: `coreevents.Envelope[persons.Changed]` with `eventId`, `eventType`, `version`, `aggregateId`, `occurredAt`, optional `correlationId`/`causationId`, and `data`.
- `data`: `{ "operation": "UPDATED", "old": {...}, "new": {...} }`.
- Each non-null snapshot uses the public `foundation/dto/response.Person` representation (GID, names, contacts, status, version, createdAt, updatedAt).
- `CREATED`: old is null, new is the created Person. `UPDATED`: old/new are complete snapshots. `DEACTIVATED`: old/new are complete snapshots with status transition. A deactivation emits only DEACTIVATED, not UPDATED as well.
- No publish on failed operations, idempotent create replay, no-op update, or already-deactivated Person.
- `DELETED` is **not implemented** in the current Person service. New null is reserved for future deletion semantics; do not invent a delete event.

This change replaces the earlier `foundation_person_created` / `foundation.person.created.v1` consumer contract. Migrate subscribers instead of assuming backward compatibility.

**Security:** Complete Person DTOs contain names, email addresses and phone numbers. Scope subscribers by authorization, protect broker storage and retention, and never log raw event snapshots. Do not generalize full-snapshot publishing to every aggregate automatically. Consider data minimization and compliance when evolving the contract.

### Consumer-side diff

The producer emits `{operation, old, new}` only; it does not include a computed diff. Consumers decode into `Envelope[persons.Changed]` and can run:

```go
changes, err := diff.Compare(envelope.Data.Old, envelope.Data.New)
```

Use `grovehub/libs/core/diff`, which is generic, typed via Go type inference, and returns `([]diff.Change, error)`. The result is a flat list of JSON field paths with old and new values (for example `contacts.emails[0].value`); arrays compare by position. The domain event contract lives in Foundation, while the comparison algorithm lives in core. Consumers use the Infrastructure `MessagingClient.Consume` transport; no separate HTTP client is needed to receive messages. Ensure handlers are idempotent because broker redelivery can occur.

## 8. Composition

```text
Foundation.New
    ├─ PersonService (composed use-case operations)
    │  ├─ Create Persist    ┐
    │  ├─ Update Persist    ├─ ChangedProducer.Enqueue inside UnitOfWork
    │  └─ Deactivate Persist┘
    │       ↓
    │  ChangedPipeline.Rules() = ShouldEmit + Build + Enqueue
    │       ↓
    │  core outbox.Store.Enqueue (foundation.event_outbox)
    └─ Dispatcher (lifecycle-owned polling worker)
            ↓ pending / expired leases (SKIP LOCKED)
       Infrastructure MessagingClient.Publish
            ↓
       NATS JetStream (foundation_person_changed)
```

Only outbound event data is stored in Foundation's outbox table. Other services can
use the same core store with their own schema and migrations. `foundation.Module.Close`
stops the polling worker before Infrastructure and the DB are closed.


## 9. Current Delivery Model: Transactional Outbox

Schema migration: `backend/services/foundation/migrations/000002_create_event_outbox.sql`.

```text
Person write + Contact write + event outbox insert
           ↓ (same PostgreSQL transaction)
   Commit all or roll back all
           ↓ asynchronous polling
   Lock eligible messages with SKIP LOCKED
           ↓
   Publish with stable subject/metadata; mark SENT
   On failure: mark PENDING with capped exponential delay
   On worker crash: reclaim expired SENDING leases
```

This provides durable **at-least-once** delivery; it cannot guarantee exactly once
between PostgreSQL and NATS. If a broker publish succeeds but marking SENT fails,
the event may be delivered again. Consumers must persistently deduplicate by
`eventId`. Events may be delivered later than the HTTP response, and subscribers
must tolerate out-of-order events using aggregate version when needed.

Unlike the previous best-effort post-commit publisher, an outbox insert failure
now rolls back the Person mutation and surfaces an error to the business caller.
An idempotent Create replay, no-op Update or already-deactivated Person does not
enqueue a new event. Create/Update/Deactivate all record events in their existing
write transactions, and Deactivate now explicitly uses UnitOfWork.

Dispatcher lease and retry metadata must be monitored operationally. The
migration keeps SENT rows for audit/reconciliation; add retention and dead-letter
policy before production growth. This implementation does not invent an inbound
Person side effect: the consumer adapter is available for explicit integration.


## 10. Testing

Sanity tests must verify transactionally enqueued delivery via the outbox dispatcher into real Infrastructure messaging/NATS using stable event topics.

Current sanity coverage includes:

```text
Foundation Person Changed (CREATED, UPDATED, DEACTIVATED)
Identity Account Created
```

Use unique event/payload IDs to identify the message for a test run; never create random topic names for isolation because JetStream auto-provisioning can accumulate persistent resources.

---

## 11. Documentation Rule

Whenever event structure, event naming, envelope shape, event ownership, or publishing behavior changes, update this document in the same change.

If the change also affects generic backend placement/layering, update `GROVEHUB_BACKEND_STRUCTURE.md` as well.

> **Guiding principle:** the domain owns what happened; Infrastructure owns how the message travels.

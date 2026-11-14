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

Organize domain events first by **direction** (producers versus consumers), then
by aggregate, then by operation. Create only the packages that are used.

```text
backend/
├── libs/core/events/
│   ├── envelope.go                          # Generic typed event envelope
│   ├── pipeline/rules.go                     # ShouldEmit / Build / Enqueue
│   ├── outbox/{store,dispatcher}.go          # Transactional delivery and retry
│   └── messaging/{producer,consumer}.go      # Generic transport contracts
└── services/foundation/events/
    └── producers/
        └── person/
            ├── person.go                    # ChangedProducer interface
            ├── noop.go                      # Isolated-test implementation
            ├── impl/
            │   └── person.go                # Producer implementation
            └── operations/
                └── changed/
                    ├── changed.go           # @Event payload and generated constants
                    ├── pipeline.go          # ChangedPipeline / Rules / Enqueue
                    ├── flows/
                    │   └── changed.go       # Typed Change flow
                    ├── strategies/
                    │   ├── strategy.go      # @Strategy contract
                    │   └── impl/
                    │       └── changed.go   # Policy/serialization strategy
                    └── helpers/
                        ├── changed.go       # HasChanged
                        ├── envelope.go      # BuildChanged
                        ├── message.go       # BuildMessage
                        └── metadata.go      # Broker headers
```

The Foundation **consumers/** package does not exist yet: Foundation does not
currently consume any inbound Person event. Add
`services/foundation/events/consumers/<aggregate>/` only when a real inbound
subscription and handler exists; do not add empty directories merely for symmetry.

The `Operation` type and lifecycle values remain owned by
`foundation/constants/events/operation.go`. The `@Event` generator emits the
topic/type/source constants into the new Changed operation package, preserving
existing message topic and wire payload. Other service event layouts may retain
their current shape until separately migrated.


## 3. Responsibility Split

- `events/producers/person/person.go`: the domain producer interface used by
  Person Create, Update and Deactivate. The concrete adapter lives in `impl/`.
- `events/producers/person/operations/changed/changed.go`: stable Changed
  event contract and `@Event` metadata. The `pipeline.go` composes
  `ShouldEmit`, `Build` and `Enqueue` with core event pipeline rules.
- `operations/changed/flows/`: the event-specific typed Change state.
- `operations/changed/strategies/`: event eligibility and immutable envelope
  creation; `helpers/` owns pure Person Changed logic.
- `libs/core/events/outbox/`: durable transaction-bound inserts, leased claims,
  retry scheduling, transport delivery and sent-state tracking.
- `services/person/operations/`: owns business writes and calls the producer
  within the same UnitOfWork as successful Create, changed Update and
  Deactivate. GetByGID and Search never enqueue events.
- `events/consumers/` (future): inbound subscription and handler ownership,
  only when Foundation actually consumes another event.

Domain event helpers must not call persistence, NATS adapters or HTTP.
Infrastructure owns how messages travel, not what events mean.


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
policy before production growth. Foundation currently implements **no inbound event consumer**. Add a consumer
only when a real subscription and business handler are needed.


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

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

Events are grouped by aggregate/resource under the owning service.

```text
backend/services/<service>/events/
└── <aggregate>/
    ├── <aggregate>.go
    ├── impl/
    │   └── <aggregate>.go
    └── helpers/
        └── <aggregate>.go
```

Examples:

```text
backend/services/foundation/events/persons/
├── person.go
├── impl/person.go
└── helpers/person.go

backend/services/identity/events/accounts/
├── account.go
├── impl/account.go
└── helpers/account.go
```

Foundation lifecycle operation names and the `Operation` type live at `backend/services/foundation/constants/events/operation.go` (Go package `events`), not in `libs/core/constants` or `events/persons/person.go`. Event contracts import this Foundation-owned type; producers and consumers use `foundationevents.Created`, `foundationevents.Updated`, and `foundationevents.Deactivated`. Other domains define their own operations only when required; do not impose a single global operation enum.

Create only folders that have a current responsibility.

---

## 3. Responsibility Split

### `<aggregate>.go`

Public event contracts and stable event identifiers only.

Typical contents:

```text
topic/subject constant
event-type/version constant
source constant
event payload struct
```

Consumers may import this package to deserialize the event contract.

### `impl/<aggregate>.go`

Event publishing orchestration.

It may:

```text
decorate the owning domain service
publish after a successful business operation
generate event IDs
marshal the event envelope
call Infrastructure MessagingClient
log publishing failures
```

It must not contain the core business operation itself.

### `helpers/<aggregate>.go`

Pure event construction helpers.

Typical responsibilities:

```text
build the shared event envelope
map domain/service output into the public event payload
build standard event metadata
```

Helpers should not call NATS, repositories, HTTP, or other services.

---

## 4. Event Envelope

All service events use the shared core envelope from:

```text
backend/libs/core/events
```

Canonical shape:

```text
EventID
EventType
Source
OccurredAt
Data
```

The domain payload belongs to the owning service event package; only the generic envelope belongs in core.

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
- Envelope: `coreevents.Envelope[persons.Changed]` with `eventId`, `eventType`, `source`, `occurredAt`, and `data`.
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

The root launcher creates Infrastructure messaging once and injects its public `MessagingClient` into domain modules that publish events.

```text
launcher
  ↓
Infrastructure Module
  ↓ MessagingClient
Foundation / Identity / future domain module
  ↓
events/<aggregate>/impl
  ↓
MessagingClient.Publish
  ↓
NATS JetStream
```

Domain event implementations depend only on the public Infrastructure client contract; they must not import NATS/Watermill adapters directly.

---

## 9. Current Delivery Model

GroveHub currently uses direct post-success event publishing:

```text
business write succeeds
  ↓
publish event
```

A publish failure is logged and must not incorrectly report that an already-committed database write failed.

A transactional outbox may be introduced later when guaranteed database-to-event delivery becomes a production requirement. Do not add an outbox prematurely without that requirement.

---

## 10. Testing

Sanity tests must verify real Infrastructure messaging/NATS delivery using stable event topics.

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

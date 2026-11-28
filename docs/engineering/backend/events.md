# Events, Producers, Consumers and Delivery Guarantees

**Scope:** generic domain event contracts and asynchronous messaging across modules.

## 1. Separation of owners

A domain module owns **what happened** (event name, payload, eligibility), while the messaging Infrastructure capability owns **how to deliver** (broker adapter, outbox table, worker, retry, leases, transport telemetry). Event history/audit owns its own long-term retention; neither a broker nor a short-lived outbox is automatically a compliance archive.

```text
services/<module>/internal/app/events/
├── producers/<aggregate>/
│   ├── <aggregate>.go                       # ChangedProducer interface
│   ├── impl/<aggregate>.go                  # typed producer composition
│   └── operations/<event>/
│       ├── <event>.go                       # versioned event data, @Event
│       ├── pipeline.go                      # ShouldEmit, Build, Enqueue/Publish
│       ├── flows/                           # typed event state, if multistep
│       ├── helpers/                         # pure envelope/metadata/snapshot construction
│       └── strategies/{strategy.go,impl/}   # typed event-building strategy
└── consumers/<aggregate>/
    ├── <aggregate>.go                       # Consumer contract
    ├── impl/<aggregate>.go                  # handler/subscription composition
    └── operations/<event>/
        ├── <event>.go                       # handler entry
        ├── helpers/                         # decoder/dedup validation
        └── strategies/                      # when needed
```

Do not create `consumers/` until a genuine subscription exists. A single event can use a small package instead of artificial flow/strategy packages.

## 2. Event envelope

Use a generic, versioned `Envelope[T]` carrying at least:

| Field | Role |
|---|---|
| `eventId` | Unique emission identity for deduplication |
| `eventType` | Stable, namespaced semantic identity |
| `version` | **Event schema version**, not DB aggregate version |
| `aggregateId` | Stable identity of affected aggregate |
| `occurredAt` | UTC timestamp of domain occurrence |
| `correlationId` | Cross-call workflow correlation when known |
| `causationId` | Message/command that caused this event when known |
| `data` | Typed, domain-owned payload |

Separate broker metadata (`traceparent`, topic, source, headers) from the domain payload. Use a **stable topic per event category**, not a random topic per event instance.

## 3. Producer SRP

- Decide whether a state transition warrants emission (not a no-op or idempotent replay).
- Produce typed data that accurately reflects committed change; never invent a mutation.
- Build envelope and metadata with a stable event ID.
- Select delivery intent: durable outbox when atomic with business mutation, or direct broker notification when downstream retry is intentionally external.
- Respect permissions, PII minimization, data ownership and event schema compatibility.
- Do not import NATS/Watermill SDK in the domain producer.
- Do not put domain eligibility logic in Infrastructure.

Changes in event shape require schema-version policy and consumer compatibility review.

## 4. Durable event workflow

```text
UnitOfWork.Run(ctx)
  -> persist aggregate
  -> persist related entities
  -> EventStore.Enqueue(ctx, pending_event)
  -> COMMIT
worker: claim PENDING/expired lease -> direct provider Publish
       -> mark SENT; retry on failure
consumer: verify/decode -> deduplicate eventId -> apply its own transaction
```

The outbox shares the **same transaction and PostgreSQL database** with business mutations. Calling enqueue after committing the business record is **not atomic**; doing so inside the same context only matters if the repository uses the transaction attached to that context.

An outbox ensures an **at-least-once** delivery attempt, not exactly-once end-to-end execution. A crash after broker ACK and before `MarkSent` can redeliver. Consumers must persist event-ID idempotency keyed by handler/subscriber identity, protect ordering by aggregate version if required, and handle poison messages.

## 5. Direct publication

`Publish(ctx, topic, payload, metadata)` and `PublishWithOptions(..., PublishOptions{Persist:false})` publish to the configured broker directly. There is no PostgreSQL outbox write. Use for explicitly best-effort telemetry or for an external webhook whose sender retries after a non-successful response. Never acknowledge the upstream source before required broker delivery if that is the retry contract.

`PublishWithOptions(..., PublishOptions{Persist:true})` creates durable delivery intent in the outbox. It does **not** imply indefinite retention in PostgreSQL. Do not call both the domain producer's direct `EventStore.Enqueue` and `PublishWithOptions(Persist:true)` for the same logical event.

## 6. Consumer SRP

A consumer should:

1. Verify transport/authentication context and accepted event type/version.
2. Parse the known envelope and validate required fields; malformed data follows a bounded DLQ policy.
3. Acquire/record event-ID idempotency in the same transaction as its own state mutation.
4. Invoke its own use-case service; never directly mutate another module's schema.
5. ACK only after safe processing; NACK/retry transient failures within bounded backoff policy.
6. Preserve trace/correlation without dumping payloads into logs.

Do not treat subscription registration as proof of IAM or tenant scoping.

## 7. Ordering, concurrency, and replay

- No global ordering guarantee is assumed from a general broker.
- Choose ordering key by aggregate ID when causal sequence matters.
- Use aggregate version or expected state to reject stale updates.
- Define replay, late joiner, high-throughput and dead-letter behavior per topic.
- Keep retention (broker stream vs outbox rows vs audit storage) documented separately.
- Ensure worker shutdown drains or cancels in-flight operations within a bounded deadline.
- During migration, maintain topic/event compatibility and deterministic payload versioning.

## 8. Testing contract

Test: no emission for no-op, idempotency replay, transaction rollback without enqueue, worker retry/lease expiry, publish-success/mark-sent-failure redelivery, consumer duplicate, unsupported versions, out-of-order aggregate versions, poison payload, cancellation, and no secret/PII exposure. A mocked `Publish` is not proof of durable delivery; add broker/PostgreSQL sanity tests.

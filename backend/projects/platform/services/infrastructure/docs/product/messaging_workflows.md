# Infrastructure — Messaging Workflow Execution

## Preconditions

A running backend has connected to PostgreSQL and the configured NATS JetStream server, required migrations are applied, and the internal Messaging client is injected by bootstrap. Use a **test topic** and non-sensitive synthetic payload in an isolated environment; verify infrastructure health first.

## 1. Immediate publish (no PostgreSQL outbox record)

```go
err := messaging.Publish(ctx, topic, payload, metadata)
```

**Expected:** the direct publisher returns success or a broker error. There is no new `infrastructure.messaging` outbox row attributable to this call. A broker outage must not be silently reported as success.

Equivalent explicit option: `PublishWithOptions(ctx, topic, payload, metadata, PublishOptions{Persist:false})`.

## 2. Persist a business event for later delivery

Inside the same UnitOfWork used for the business record:

```go
err := messaging.PublishWithOptions(
    txCtx, topic, payload, metadata,
    messagingclient.PublishOptions{Persist:true},
)
```

**Expected:** one `PENDING` row in `infrastructure.messaging` as part of the transaction. On rollback neither the business row nor outbox event commits. On commit the dispatcher publishes asynchronously and transitions it to `SENT`. Capture a stable `event-id` in metadata for consumer deduplication.

Alternative: a domain producer calls `EventStore.Enqueue(txCtx, event)` directly. Use **one** of these routes for one logical event, not both.

## 3. Consume a message

A module registers `MessagingClient.Consume(ctx, topic, handler)` with a typed context-aware handler. Verify an authorized, version-compatible payload, idempotent side effects, and appropriate ACK/retry behavior. The consumer's business logic belongs to the subscribing module, never the Infrastructure dispatcher.

## 4. Reproduce retry without losing work

1. Start a test worker, enable test-topic producer and record emitted `event-id`.
2. Disable broker connectivity **only in the isolated environment**.
3. Queue a persistent event; verify PENDING/retrying state persists, not a false SENT.
4. Restore broker; verify eventual delivery.
5. Simulate broker ACK followed by `MarkSent` failure in a controlled test; expect possible duplicate delivery and confirm consumer deduplication.
6. Verify SENT cleanup after configured age; pending/retrying rows must not be deleted.

## 5. Failure acceptance matrix

| Scenario | Expected outcome |
|---|---|
| Direct publish, broker down | Return error; no DB outbox |
| Persist=true, DB insert fails | Return error, rollback mutation |
| Persist=true, broker down after commit | Outbox remains retryable |
| Worker restart | Pending/expired lease can be reclaimed |
| Broker ACK, MarkSent fails | At-least-once redelivery possible |
| Consumer receives duplicate `event-id` | No duplicated business side effect |
| Retention cleanup | Prune only old SENT rows |
| Missing config | Startup validation failure |

`Persist:true` is a **delivery persistence** choice, not a permanent audit or archive feature.

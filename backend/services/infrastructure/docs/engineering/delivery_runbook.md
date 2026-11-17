# Infrastructure — Delivery Operations Runbook

**Purpose:** diagnose broker and outbox delivery without exposing event payloads.

## 1. Startup checks

- Database migrations for the shared PostgreSQL schema are complete.
- `GROVEHUB_MESSAGING_PROVIDER`, `GROVEHUB_NATS_URL`, queue and durable prefixes are configured.
- NATS JetStream is reachable and the registered transport strategy passes `Health`.
- The worker starts only once per intended runtime, or multi-runtime concurrency is safe via row leasing.

## 2. Message lifecycle

```text
PENDING -> claim/lease (SENDING) -> direct broker publish -> SENT
                                   └─ failure -> retry scheduling
lease expires -> reclaim/retry
SENT + older than retention -> prune
```

Current implementation uses approximately 1-second worker ticks, claims batches of four, leases for around 90 seconds, schedules retry state durably and prunes up to 500 eligible sent rows per minute older than 24 hours. These are implementation settings, **not generic platform guarantees**; update this file if tuned.

## 3. Failure diagnostics

**Broker unavailable:** direct publish returns error; persistent messages continue pending/retrying. Confirm URL, JetStream server and credentials; do not mark records sent manually.

**Backlog grows:** inspect counts by `status`, next-attempt timing, attempts and lease expiry. Determine whether send errors, network saturation or DB contention dominate.

**Repeated event delivered:** expected possible at-least-once behavior. Verify consumer deduplication by event ID, not accidental extra producer calls. Check whether `MarkSent` failed after publish.

**Stuck sending:** inspect lease expiry, worker cancellation and DB time; expired leases are eligible to reclaim.

**Retention concern:** only `SENT` records are pruned; permanent audit history must live in another owner-controlled store.

## 4. Data safety and shutdown

Do not print raw payloads, contact data, message contents or credentials. Use event IDs, topic, status, attempts, lease and trace metadata for diagnostics. Stop producers/consumers safely before manual table repair; back up affected rows and rehearse a recovery plan. On shutdown, cancel and join the dispatcher before closing broker/database connections.

## 5. Acceptance checks

Verify startup health, direct publish, one persist operation inside UnitOfWork, worker delivery, restart with queued messages, crash between broker ACK and sent-state update, retry after transient failure, and pruning only sent/expired entries. A consumer must remain correct under duplication.

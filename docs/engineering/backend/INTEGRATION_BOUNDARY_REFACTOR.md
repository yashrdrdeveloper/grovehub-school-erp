# Integration ownership refactor — authoritative contract

> Supersedes the obsolete service-ownership descriptions in INTEGRATION_PHASE_3.md and INTEGRATION_PHASE_4.md. Those two files document historical implementations, not current behavior.

## Boundary

- **Integration**: provider catalogs, capability adapters, OAuth/secret references, connections, provider subscription/watch lifecycle, webhook verification, connection health tests and stateless provider API calls.
- **Communication (future, NOT implemented here)**: mailbox setup/selection, user-facing messages, threads, attachments, `connectionId + providerMessageId` deduplication, initial and incremental ingestion, per-mailbox synchronization cursors, retries/checkpoints, delivery outcomes and consumer idempotency.
- **Infrastructure/NATS**: the transport for versioned provider notifications. **Temporal is not required in Integration**. Long-running user workflows belong to a future orchestrator.

## Refactored provider API

Integration exposes a Gmail capability facade for connected accounts:
```
GET  /integrations/connections/{id}/mail/profile
GET  /integrations/connections/{id}/mail/messages?pageToken=
GET  /integrations/connections/{id}/mail/messages/{messageId}
GET  /integrations/connections/{id}/mail/history?startHistoryId=&pageToken=
POST /integrations/connections/{id}/mail/send   {"raw":"<base64url RFC 5322>"}
```
Every route checks ownership using a verified principal context (no untrusted header). Provider-specific adapter translates to Gmail API. The facade does not write email bodies, message IDs, threads or consumer cursors into Integration's database. `/mail/messages/{messageId}` currently returns **metadata** rather than full raw content: Communication ingestion will need expanded payload/attachment capabilities later.

The previous `POST /integrations/connections/{id}/sync` endpoint is **removed**. Communication will orchestrate sync using these capabilities when built.

## Webhook to consumer hand-off

1. Google Gmail Pub/Sub sends authenticated notification; Integration validates Google's signed OIDC JWT.
2. Integration maps mailbox address to active connected accounts and publishes **one event per connection**, `integration.mailbox.changed.v1`, on NATS subject `integration_mailbox_changed`.
3. Generic event envelope includes event identity, version, aggregateId=connectionId and timestamps. Event `data`: `{provider,connectionId,tenantId,ownerGid,mailbox,historyHint}`. The history hint is *not* an authoritative cursor.
4. A future Communication consumer uses **its own last-committed history cursor** to fetch provider history through Integration, upserts into Communication and advances the cursor *only after* its database work succeeds. Consumer must be idempotent and handle repeated/reordered notifications.
5. Integration acknowledges Pub/Sub only after the messaging publish reports success; failures return non-2xx so Pub/Sub can retry. **No end-to-end guarantee yet**: NATS retention/durable consumer, transactional outbox, and Communication ingestion must be configured before depending on automatic mail synchronization.

## Schema/data migration

- Existing schema migrations `000001` and `000003` are retained for migration history and backwards compatibility.
- New **`000005_isolate_legacy_mail_state.sql`** moves `integration.mail_references` and `integration.sync_states` into an `integration_legacy` schema **without dropping any data**. This is a staging archive and is not read or written by Integration.
- Before implementing Communication, design Communication's generic message schema; migrate historical rows and cursors from `integration_legacy` using a carefully reviewed, idempotent migration. Do not blindly recreate `MailReference` as an Integration entity.
- Connection records, encrypted OAuth credentials, provider subscription state, and connectivity test history **remain in Integration**.

## Known limitations

- Gmail outbound API accepts raw MIME and returns a provider response, not delivery confirmation.
- OAuth and protected actions are still gated by missing verified IAM middleware (existing Phase 1-4 issue); do not bypass to test.
- Outlook remains a catalog placeholder and advertises no implemented capabilities.
- New source files and documentation were changed on GitHub; the Go test suite and migration were not executed against a live deployment.

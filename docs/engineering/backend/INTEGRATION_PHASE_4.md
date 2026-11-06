# Integration Phase 4 — Gmail watch and push-triggered synchronization

## Implemented

- `POST /integrations/connections/{id}/watch` registers Gmail's `users.watch` using the configured Pub/Sub topic. Active Gmail connection + verified owner required.
- `POST /integrations/webhooks/gmail` verifies a Google-signed OIDC JWT using Google's HTTPS JWKS, RS256, issuer, audience, service account email, email verification, issued/expiry times, and then decodes the Pub/Sub-wrapped `emailAddress` + `historyId`.
- The webhook looks up *active* Gmail connections by the provider's mailbox email address and triggers the existing full/incremental sync for each.
- A Postgres **session advisory lock** on the `connectionId` serializes manual and webhook-triggered sync *across application instances*; concurrent attempts return a busy error and Pub/Sub can retry.
- Only successful processing returns HTTP 204; failed processing returns non-2xx. **Processing is synchronous**; if a large mailbox takes longer than Pub/Sub's push deadline, notifications can be redelivered. The existing message-reference ledger is idempotent.
- Gmail watch expiration is persisted in `integration.webhook_subscriptions`, and a background renewal loop checks near-expiring watches every six hours. Independent application instances can attempt renewal; the API and DB upsert make this safe.
- Migration `000004_mail_watch.sql` adds the watch uniqueness and lookup/expiry indexes.
- Integration still stores provider-side mail references only, not full bodies, attachments, or user-facing Communication messages. No Temporal.

## Google Cloud setup

Configure the Gmail API in the same Google Cloud project as the Pub/Sub topic, e.g. `projects/PROJECT_ID/topics/grovehub-gmail`. Grant Gmail's push service account `gmail-api-push@system.gserviceaccount.com` Pub/Sub Publisher permission for that topic.

Create a **push subscription** for the topic targeting:
`https://YOUR_HOST/integrations/webhooks/gmail`

Enable **authenticated push** using a dedicated Google service account. Set the OIDC audience to the exact expected value. Ensure the Pub/Sub service agent can mint tokens for that service account.

Environment variables:
- `GROVEHUB_GMAIL_PUBSUB_TOPIC`: full topic path (`projects/.../topics/...`)
- `GROVEHUB_GMAIL_PUBSUB_AUDIENCE`: exact push JWT audience
- `GROVEHUB_GMAIL_PUBSUB_SERVICE_ACCOUNT`: dedicated push service-account email

Also configure the Phase 2 Gmail OAuth credentials and encryption key. Once connected and initial sync has run, register the mailbox watch. Gmail watch must be renewed before its provider-defined expiration.

## Safety and operational limitations

1. The protected mailbox connection routes still require a trusted IAM principal; **no such launcher middleware is wired yet**. They return 403 until that dependency exists. The webhook has its own Google JWT verification.
2. Initial sync and push-triggered sync are synchronous, not queued. When pushing high-volume mailbox traffic, migrate to **durable enqueue then 2xx ACK** and asynchronous processing (NATS JetStream/transactional queue). For now, Pub/Sub retries on failure.
3. Connection deduplication and event processing are idempotent. A single webhook can trigger sync on multiple active connections referencing the same mailbox (across tenants).
4. The Postgres advisory lock uses a dedicated pooled connection; configure database pools with sufficient connections (**at least two connections per instance**), or long syncs will stall.
5. A push notification's historyId is a **hint**; syncing always reads the stored cursor to avoid skipping changes. A stored history cursor is advanced only on completed sync.
6. Watch renewals are periodic only while the process is running; missed renewals are not immediately repaired on startup. Large initial scans and Gmail errors remain bounded by Phase 3.
7. Google JWKS is cached for one hour. A newly rotated kid may initially be rejected until refresh; Pub/Sub retry is required.
8. No live Google requests, database migrations, or Go tests were executed in this change.

## Source references

- https://developers.google.com/workspace/gmail/api/reference/rest/v1/users/watch
- https://developers.google.com/workspace/gmail/api/guides/push
- https://cloud.google.com/pubsub/docs/authenticate-push-subscriptions

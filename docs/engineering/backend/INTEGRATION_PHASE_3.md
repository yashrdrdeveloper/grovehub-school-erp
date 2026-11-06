# Integration Phase 3 — Gmail full + incremental synchronization

## Implemented

- `POST /integrations/connections/{id}/sync` requires an authenticated, authorized connection owner (trusted IAM principal; wiring still pending).
- Gmail API `messages.list` scans mailbox pages; for every message, `messages.get?format=metadata` retrieves ID, thread ID, history ID and system labels. **No message content is stored**.
- Each provider message reference is uniquely keyed by `(connection_id,provider_message_id)`; idempotent upsert handles resync without duplicates.
- Mail direction derives from Gmail's `SENT` label, otherwise INBOUND; it is not a verified recipient identity.
- Full sync captures a history baseline **before** scanning to avoid losing concurrent changes; it advances the cursor only after all pages succeed.
- Incremental sync reads `history.list` from saved cursor, upserts added/relabeled messages and flags deleted messages. Cursor advances only after all pages succeed.
- Gmail history returning 404 triggers full sync, per official Gmail guidance.
- Both initial and incremental syncing allow at most 100 pages per invocation. A larger mailbox yields `INTEGRATION_SYNC_LIMIT` and leaves cursor unchanged; previous upserts are safe to replay.
- Records are **Integration's provider-reference ledger**, not Communication's user-facing messages or message bodies. A future contract/event bridge is required to populate Communication and associate people/threads/attachments.
- `000003_mail_sync.sql` introduces `integration.mail_references`. Apply with `scripts/database/migrate.sh init integration`.
- No Temporal needed.

## Not yet implemented / important operational limits

- No automatic scheduler, Gmail Pub/Sub watch or authenticated push notifications yet; sync is manual/on-demand.
- No IAM verified-principal middleware available in the launcher, so protected endpoints currently return 403 until that dependency is implemented.
- No cross-worker per-mailbox sync mutex. Parallel sync requests could race; serialize by connection before production enablement.
- No persistent pagination checkpoint: mailboxes over 10,000 messages cannot complete full sync using this version.
- No message bodies, attachments, Communication delivery events, or UI mailbox display.
- No live Gmail account/API tests or Go tests executed against the remote repository.
- A Gmail 404 on retrieving metadata for a message concurrently deleted may fail a sync; rerun to converge.

## Provider behavior references

Official Google documentation:
- https://developers.google.com/workspace/gmail/api/reference/rest/v1/users.messages/list
- https://developers.google.com/workspace/gmail/api/reference/rest/v1/users.history/list
- https://developers.google.com/workspace/gmail/api/reference/rest/v1/users.messages/get

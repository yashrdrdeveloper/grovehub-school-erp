# Integration Phase 2 — Gmail OAuth and mailbox connectivity tests

## Scope implemented
- Each connection has its own `connectionId`, tenant/owner, external mailbox address and independently encrypted Google OAuth credential.
- Multiple Gmail mailboxes can coexist under the same tenant/owner. Unique ACTIVE mailbox per tenant/provider prevents accidental duplicate active connections.
- Gmail OAuth start binds a pending connection to a cryptographically random state; only a SHA-256 digest is stored and the state expires in 10 minutes. Callback consumes it once, exchanges Google's code, resolves Gmail profile, and activates the connection.
- Access/refresh tokens are encrypted in PostgreSQL using AES-256-GCM; per-deployment 32-byte encryption key supplied through an environment variable. An external secret manager can replace this storage later.
- INBOUND test checks Gmail readonly mailbox access; **does not claim inbound message delivery**.
- OUTBOUND test sends one self-addressed test email through Gmail API; success means Gmail accepted the send request, **not confirmed delivery**.
- Every test saves a SUCCESS/FAILED row, timestamps and sanitized outcome; latest result is copied to connection.lastTestedAt/lastTestStatus. Errors are not stored verbatim.
- Disconnect deletes the saved token and pending OAuth states. Revocation at Google is not implemented.
- No Temporal needed; Gmail sync, notifications, webhook processing, message upsert and Communication Setup are later phases.

## Required Google configuration
Register an OAuth web application in Google Cloud, configure Gmail API, necessary scopes and redirect URI to exactly:
`https://YOUR_HOST/integrations/oauth/gmail/callback`

Set:
- `GROVEHUB_GMAIL_CLIENT_ID`
- `GROVEHUB_GMAIL_CLIENT_SECRET`
- `GROVEHUB_GMAIL_REDIRECT_URI`
- `GROVEHUB_INTEGRATION_ENCRYPTION_KEY` (base64 of 32 cryptographically random bytes).

If any setting is missing or invalid, Gmail OAuth/testing are disabled and other Integration APIs remain accessible. Do not log or commit these secrets. Keep the encryption key stable, protected, and backed up; rotating it needs a credential re-encryption procedure.

Apply `scripts/database/migrate.sh init integration` to run migrations 000001 and 000002.

## HTTP routes
- POST `/integrations/connections/{id}/oauth/start` -> `authorizationUrl`
- GET `/integrations/oauth/gmail/callback?state=...&code=...` -> activation result
- POST `/integrations/connections/{id}/test` body `{"direction":"INBOUND"}` or `{"direction":"OUTBOUND"}`
- GET `/integrations/connections/{id}/tests` -> latest 100 tests
- GET `/integrations/connections/{id}` -> includes `externalId`, `lastTestedAt`, `lastTestStatus`

**Current deployment blocker:** No verified IAM middleware yet populates Integration's trusted principal, so all protected connection endpoints return 403. OAuth callback itself accepts only valid one-time OAuth states created through the protected start endpoint. Connect IAM and tenant authorization before using the flow. **Do not replace this guard with unchecked headers or request-supplied tenant IDs.**

## Limitations and required follow-up
- OAuth scopes require Google consent/verification per Google's policy; credentials and configured callback needed to test live.
- Connection test is on-demand, not a periodic scheduler.
- Token refresh is supported while running tests. Refresh-token rotation/revocation, production rate limiting, richer test health, sync cursors, and Gmail Pub/Sub are subsequent work.
- OAuth state cleanup of expired entries should be scheduled by maintenance.
- Test history is capped on reads, but retention/pruning policy must be added.
- No live Google operations/tests were executed in this change.

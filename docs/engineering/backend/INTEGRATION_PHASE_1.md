# Integration Phase 1 — provider catalogue + connection persistence

## Boundaries

- Integration owns provider catalog/strategy contracts and provider connections.
- Communication (later) owns actual messages, threads, person communication history.
- Infrastructure owns NATS messaging. Temporal belongs to the future orchestrator, not Integration.
- Gmail/Outlook in this phase are **catalog descriptors only**. Their declared capabilities are planned, NOT functioning.
- The HTTP controller uses a typed, trusted `VerifiedPrincipal` in request context. No current launcher/IAM middleware supplies it, so all connection endpoints **deny by default** (403) until verified authorization is connected. Never accept tenant or actor from untrusted headers.
- The public `GET /integrations/providers` catalog is available without credentials.

## Current APIs

- `GET /integrations/providers`: provider catalog.
- `POST /integrations/connections`: create `PENDING_SETUP` connection for caller actor (available once verified principal middleware is installed).
- `GET /integrations/connections`: list actor's connections (capped to 200).
- `GET /integrations/connections/{id}`: retrieve a connection only for owner.
- `DELETE /integrations/connections/{id}`: mark disconnected, clear credential reference.

The GORM entity and SQL migration use schema `integration`; migration must be applied using the repository's migration runner. `sync_states` and `webhook_subscriptions` are reserved schemas for later real implementations.

## Next phases

1. Wire verified IAM/tenant middleware and authorization tests before enabling connection management.
2. Implement genuine Google OAuth with CSRF state, encrypted token reference, refresh, scopes and callback; mark ACTIVE only after verification.
3. Add mail send/read tests, incremental Gmail history sync, webhook delivery and deduplication into Communication.
4. Implement Outlook Graph as another strategy.

Never treat a PENDING_SETUP record as a connected mailbox. No secrets are accepted over the Phase 1 connection API. No SMTP/IMAP, Maddy, Temporal or mail sync is active.

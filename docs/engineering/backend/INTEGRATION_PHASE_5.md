# Integration Phase 5 — generic provider capabilities

## Architecture and ownership

**Integration owns**: provider catalog/registries, provider instances/connections, encrypted OAuth token storage, refresh, provider-specific API adapters, webhook subscriptions and verification, connectivity tests, and transport-neutral execution of provider actions.

**Communication (future, NOT implemented)** owns: mailbox configuration/selection, messages, threads, attachments, provider message ID deduplication, inbox/outbox history, synchronization cursor/checkpoint, delivery state, and all message persistence.

**Orchestrator (future)** owns Temporal workflow logic; Integration contains no Temporal dependency.

## Generic capability strategy

`providers.ActionRegistry` maps `providerId + action` to an optional `ActionHandler`. A handler is registered only if the provider client is configured. `providers.OAuthRegistry` provides separate optional OAuth strategies for provider-specific auth/token refresh. Provider-specific errors are normalized to safe HTTP categories. Unknown action returns a controlled error, never a vendor switch in the generic controller.

The **single generic action endpoint** is:

```
POST /integrations/connections/{id}/actions/{action}
{"input": { "... provider capability arguments ...": "..." }}
```

Integration checks verified connection ownership, active status and registered action, obtains a valid OAuth token from its encrypted vault, then invokes the selected provider strategy. No client-supplied access token, tenant ID, secret, or untrusted destination URL is accepted.

Provider catalog:
```
GET /integrations/providers
GET /integrations/providers/{provider}
```
The detail endpoint distinguishes `availableActions` from advertised provider capabilities and reports `oauthConfigured`.

## Supported actions

| Action | Gmail | Outlook | Input |
| --- | --- | --- | --- |
| `mail.profile` | Yes | Yes | `{}` |
| `mail.list` | Yes | Yes | `{"pageToken":""}` |
| `mail.get` | Yes | Yes | `{"messageId":"..."}` |
| `mail.attachment.get` | Yes | Yes | `{"messageId":"...","attachmentId":"..."}` |
| `mail.send` | Yes | Yes | `{"raw":"base64url RFC5322"}` |
| `mail.history` | Yes | No | `{"startHistoryId":"...","pageToken":""}` |
| `mail.delta` | No | Yes, Inbox | `{"cursor":""}` |

Gmail returns provider IDs for a successful send; Microsoft Graph `sendMail` returns HTTP 202 without a message ID, so Integration returns `{"accepted":true}`. Neither is proof of delivery. Provider results are transient; Communication must persist messages/attachments/cursors when implemented.

Gmail `mail.get` requests FULL format and includes the Gmail MIME part structure. Outlook `mail.get` returns a Microsoft Graph message representation. These outputs are provider-specific intentionally. A future Communication adapter will normalize them.

For Microsoft `mail.delta`, the cursor is the **entire** trusted Graph `@odata.nextLink` / `@odata.deltaLink` for the Inbox resource, not a database cursor managed by Integration. Untrusted hosts and paths are rejected. Delta is per folder, and other folders can be added later without changing Integration's ownership.

## Provider configuration

Gmail requires:
- `GROVEHUB_GMAIL_CLIENT_ID`
- `GROVEHUB_GMAIL_CLIENT_SECRET`
- `GROVEHUB_GMAIL_REDIRECT_URI` (ends with `/integrations/oauth/gmail/callback`)

Outlook requires an Azure/Entra application supporting organizational and/or personal Microsoft accounts:
- `GROVEHUB_OUTLOOK_CLIENT_ID`
- `GROVEHUB_OUTLOOK_CLIENT_SECRET`
- `GROVEHUB_OUTLOOK_REDIRECT_URI` (ends with `/integrations/oauth/outlook/callback`)
- Optional `GROVEHUB_OUTLOOK_TENANT` (default: `common`)

Outlook delegated permissions: `offline_access`, `User.Read`, `Mail.Read`, `Mail.Send`. Obtain consent and configure matching redirect URI in Entra.

Both require `GROVEHUB_INTEGRATION_ENCRYPTION_KEY` as base64-encoded 32 bytes. The encrypted credential is keyed to each connection and not exposed by APIs. `POST /integrations/connections/{id}/oauth/start` and `GET /integrations/oauth/{provider}/callback` cover both providers.

## Connection tests

Both Gmail and Outlook support `POST /integrations/connections/{id}/test` with `{"direction":"INBOUND"}` or `{"direction":"OUTBOUND"}`, saving history and latest test time/status. INBOUND checks mailbox read access only; OUTBOUND submits a self-addressed test message. Neither confirms delivery or receipt.

## Limitations and deployment prerequisites

- **IAM blocker:** The launcher still does not wire trusted principal middleware to Integration. All protected connection/action endpoints remain inaccessible (403) until verified IAM is implemented; do not use untrusted tenant headers.
- Outlook Graph **webhook subscriptions and lifecycle notifications are not implemented**. Gmail Pub/Sub watch and NATS mailbox-change event handling remain provider-specific. Outlook delta is available on demand only.
- Outlook and Gmail have different semantics; only implemented actions are registered. OAuth settings must be configured before their handlers become available.
- Gmail/Outlook external calls, Ginkgo tests and PostgreSQL migrations have **not been executed in this environment**. Production rollout requires real OAuth credentials, a mailbox integration test, and service authorization tests.
- Large attachments/body handling is constrained by provider API response limits; implement streaming/object storage integration later if required. No message content is persisted in Integration.
- Duplicate calls to `mail.send` may create duplicate messages. The future Communication sender must implement idempotency/intent tracking; Integration cannot promise exactly-once delivery.
- Existing legacy message data remain archived under `integration_legacy` and untouched.

## References

- https://developers.google.com/workspace/gmail/api/reference/rest/v1/users.messages
- https://learn.microsoft.com/en-us/entra/identity-platform/v2-oauth2-auth-code-flow
- https://learn.microsoft.com/en-us/graph/api/user-list-messages
- https://learn.microsoft.com/en-us/graph/api/user-sendmail
- https://learn.microsoft.com/en-us/graph/api/message-delta

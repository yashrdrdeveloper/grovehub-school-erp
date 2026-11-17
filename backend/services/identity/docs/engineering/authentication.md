# Identity — Registration, Login and Password Setup Implementation

## Authority and scope

General layering, SRP, transactions and annotation conventions are governed by `docs/engineering/backend/backend.md`. This document records Identity-specific decisions and remaining integration risks.

## HTTP boundaries

- Public controller: `controller/public/`: login and set-password.
- Protected controller: `controller/protected/`: admin-only registration.
- Routes are explicitly registered by `routes/http.go`; route naming itself is not an authorization mechanism.
- A verified, short-lived opaque access token becomes a trusted account ID in context only after session lookup and ACTIVE account verification.
- Interim admin permission is a configured allowlist of account UUIDs. Switch to a proper policy decision point/RBAC model before production deployment.

## Registration

Protected request -> authorizer -> validate idempotency -> normalize/fingerprint -> Foundation person resolution -> account PENDING_SETUP -> password setup token digest + encrypted delivery job -> commit -> email worker -> mail.send. No plaintext initial password, no immediate ACTIVE account.

The Foundation client may represent another transactional boundary; the local Identity transaction cannot prove cross-service ACID. Implement recovery and compensation as appropriate.

## Password setup

32 random bytes -> base64url bearer token in email -> SHA-256 digest stored in `identity.password_setup_tokens`. On POST set-password, a conditional UPDATE atomically marks the digest used only when unexpired and unused. Credential creation and account activation join the same GORM transaction. Concurrent redemption yields only one successful consumer.

## Sessions

Login validates email and password against the existing password credential service, requires ACTIVE status, and generates independent opaque access (32 random bytes) and refresh (48 random bytes) secrets. Refresh rotates both credentials and logout revokes the session. Only SHA-256 digests live in `identity.sessions`. The access-secret digest is checked on protected requests along with expiry, revocation and current account state. Token format is not JWT.

**Pending work:** logout-all, throttling, anomaly detection, access-token replay hardening, sessions invalidation after password reset and formal IAM/RBAC integration.

## Mail integration

`PasswordSetupInvitationQueue` writes token and encrypted message in the local transaction. `PasswordSetupMailWorker` claims jobs with short leases and calls `IntegrationMailSender`, which invokes the documented `mail.send` Integration action for a system-owned mailbox. Worker cancellation is tied to the launcher process.

Environment settings:

| Key | Meaning |
|---|---|
| `IDENTITY_ADMIN_ACCOUNT_IDS` | Comma-separated, verified admin account UUIDs |
| `IDENTITY_INVITATION_KEY_B64` | Base64-encoded 32-byte AES-GCM encryption key |
| `IDENTITY_SETUP_URL` | HTTPS frontend password-setup URL (no user-controlled destination) |
| `IDENTITY_SYSTEM_MAIL_BASE_URL` | HTTPS Integration service base URL |
| `IDENTITY_SYSTEM_MAIL_CONNECTION_ID` | Configured system-owned mailbox connection |
| `IDENTITY_SYSTEM_MAIL_ACCESS_TOKEN` | Credential authorized to use that connection's `mail.send` action |

If mail configuration is incomplete, invitation delivery is not started and registration fails closed. The setup URL must be on a trusted allowlisted host, and all secret configuration must use a secret store rather than source control. Validate Integration's actual service-auth permissions and mailbox ownership; `mail.send` is not an anonymous system endpoint.

## Database

Migrations `202611180001_password_setup.sql` and `202611180002_sessions.sql` create token, mail and session tables, and widen the accounts status constraint for PENDING_SETUP. A down migration with remaining pending accounts must not be applied without a documented data plan.

## Release gates

Run `gofmt`, build, DB migrations, repository unit tests, contract tests, concurrency tests, real Integration mail sandbox sends and failure injection before production. No Go compilation, SQL migration or mail send has been executed as part of this repository edit.

# Identity — Clean architecture (resource-first)

Identity is rebuilt to the same module topology as Foundation, Integration and Infrastructure. The exported root `identity.go` exposes only New and Migrate; dependency composition, routing, DTOs, entities and repositories are inside `internal/app`. All applied identity migration versions remain immutable and are embedded from `internal/resources/migrations`.

## Ownership

- Foundation owns Person and contacts; Identity stores an immutable PersonGID reference and does not modify Foundation tables directly.
- Identity owns accounts, tenant scope, username, credential hashes, session credentials, one-time proofs and queued invitation intent.
- Integration owns provider connections, Gmail/Outlook sends and OAuth credentials. Identity calls its public typed client from the registration delivery helper.
- Infrastructure owns event publication; account domain event producer work is still pending and must use the shared outbox.
- The Account search operation shares `libs/core/search` request, query, cursor and field allowlisting conventions with Foundation, plus a resource-owned `AccountConfig`.

## Module SRP

`bootstrap/impl/` composes dependencies and owns worker cancellation. `controller/account/{public,protected}` and `controller/authentication/{public,protected}` adapt typed REST requests. `services/account`, `services/authentication`, `services/credential` and `services/registration` own business operations. `repository/<resource>/impl` owns SQL. `internal/resources/templates/registration/password_setup.html` owns presentation, while `registration/operations/register/helpers/mail.go` renders and sends through Integration.

## Security

Opaque session secrets are returned only in secure HttpOnly cookies for web traffic, never JSON. Session checks verify token hash, expiry, revocation, and ACTIVE account. Cookie-authenticated writes require a matching Origin and CSRF header. Public recovery exchanges consume emailed purpose-bound proofs, issue temporary grant cookies, and allow one-time password reset or setup. Password hashes use Argon2id. Sensitive changes require recent reauthentication.

Admin actions are restricted to the verified current account, explicit admin allowlist, and matching tenant ID; no admin account in the allowlist or no tenant assignment means fail closed. Tenant-foreign accounts are not disclosed. Initial controlled admin provisioning remains a deployment prerequisite.

## Operation ordering

Account create is an explicit ordered flow: authorize -> normalize -> idempotency -> Foundation person resolution -> persist tenant account and encrypted invitation in one Identity UnitOfWork. Foundation may have an independent transaction; compensate or reconcile partial remote writes rather than claiming cross-module ACID. Invitation delivery occurs after commit in an at-least-once worker, using the Integration action interface.

The backend annotation processor checks `@ChainOfResponsibility` but does not provide business authorization or execute a chain automatically. Typed flow steps must remain explicit.

## Deployment and verification limitations

The current schema history contains a destructive initial identity migration that MUST NOT be rerun against an existing schema; only apply pending additive migrations using the existing `identity_schema_state` history. Set `IDENTITY_ALLOWED_ORIGIN`, `IDENTITY_ADMIN_ACCOUNT_IDS`, `IDENTITY_INVITATION_KEY_B64`, `IDENTITY_SETUP_URL`, `IDENTITY_RESET_URL`, `IDENTITY_SYSTEM_MAIL_TENANT_ID`, and `IDENTITY_SYSTEM_MAIL_CONNECTION_ID` using managed secrets/configuration. Password recovery is deliberately non-enumerating. Pending production release requirements include distributed rate limits, refresh-token replay family detection, admin bootstrap, data backfill/rollback, end-to-end Integration sandbox sends and PostgreSQL-backed concurrency tests.

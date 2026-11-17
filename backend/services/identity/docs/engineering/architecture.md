# Identity — Clean architecture (resource-first)

Identity is rebuilt to the same module topology as Foundation, Integration and Infrastructure. The exported root `identity.go` exposes only New and Migrate; dependency composition, routing, DTOs, entities and repositories are inside `internal/app`. All applied identity migration versions remain immutable and are embedded from `internal/resources/migrations`.

## Ownership

- Foundation owns Person and contacts; Identity stores an immutable PersonGID reference and does not modify Foundation tables directly.
- Identity owns accounts, tenant scope, username, credential hashes, session credentials, one-time proofs and queued invitation intent.
- Integration owns provider connections, Gmail/Outlook sends and OAuth credentials. Identity calls its public typed client from the registration delivery helper.
- Identity owns the `identity.account.changed.v1` producer; Infrastructure `EventStore.Enqueue` stores account events in the same local transaction as account creation, setup, updates, suspension and deactivation. Consumer deduplication remains required.
- The Account search operation shares `libs/core/search` request, query, cursor and field allowlisting conventions with Foundation, plus a resource-owned `AccountConfig`.

## Module SRP

`bootstrap/impl/` composes dependencies and owns worker cancellation. Identity controllers declare annotated HTTP verbs under the resource-specific access tiers; `routes/http.go` invokes generated registration groups only. Launcher owns tier classification and denies unauthorized external access. `controller/account/{public,protected}` and `controller/authentication/{public,protected}` adapt typed REST requests. `services/account`, `services/authentication`, `services/credential` and `services/registration` own business operations. `repository/<resource>/impl` owns SQL. `internal/resources/templates/registration/password_setup.html` owns presentation, while `registration/operations/register/helpers/mail.go` renders and sends through Integration.

## Security

Opaque session secrets are returned only in secure HttpOnly cookies for web traffic, never JSON. Identity provides a typed opaque access-token verifier to Launcher; the launcher gateway checks token hash, expiry, revocation and ACTIVE account through that verifier. The launcher also enforces trusted Origin and CSRF on stateful cookie operations. Public recovery exchanges consume emailed purpose-bound proofs, issue temporary grant cookies, and allow one-time password reset or setup. Password hashes use Argon2id. Sensitive changes require recent reauthentication.

Admin actions are restricted to the verified current account, explicit admin allowlist, and matching tenant ID; no admin account in the allowlist or no tenant assignment means fail closed. Tenant-foreign accounts are not disclosed. Initial controlled admin provisioning remains a deployment prerequisite.

## Operation ordering

Account create is an explicit ordered flow: authorize -> normalize -> idempotency -> Foundation person resolution -> persist tenant account and encrypted invitation in one Identity UnitOfWork. Foundation may have an independent transaction; compensate or reconcile partial remote writes rather than claiming cross-module ACID. Invitation delivery occurs after commit in an at-least-once worker, using the Integration action interface.

Authentication password login uses a narrow annotated `LoginStrategy` interface and `strategies/normal/operations/login/` implementation. A future OpenID Connect strategy has its own authorization-code / claims contract; it does not return null for inapplicable password-registration methods.\n\nThe backend annotation processor checks `@ChainOfResponsibility`, while the Create operation explicitly executes ordered stages with the shared `libs/core/pipeline/service.Run` runner. The annotation itself provides no authorization or transactional guarantee.

## Deployment and verification limitations

The current schema history contains a destructive initial identity migration that MUST NOT be rerun against an existing schema; only apply pending additive migrations using the existing `identity_schema_state` history. Set `IDENTITY_ALLOWED_ORIGIN`, `IDENTITY_ADMIN_ACCOUNT_IDS`, `IDENTITY_INVITATION_KEY_B64`, `IDENTITY_SETUP_URL`, `IDENTITY_RESET_URL`, `IDENTITY_SYSTEM_MAIL_TENANT_ID`, and `IDENTITY_SYSTEM_MAIL_CONNECTION_ID` using managed secrets/configuration. Password recovery is deliberately non-enumerating. Pending production release requirements include distributed rate limits (the local HTTP limiter is only per-instance), refresh-token replay family detection, admin bootstrap, tenant data backfill/rollback, end-to-end Integration sandbox sends and PostgreSQL-backed concurrency tests. Email-only login is enforced until a tenant-aware username-login contract is designed.

## Renamed entities

The Invitation entity is in `entity/invitation.go` and maps to `identity.invitations`; a forward Goose migration renames the existing invitation table without deleting queued emails. The shared one-time setup/reset credential entity remains `PasswordToken` in `entity/password_token.go`.

## Gateway limitations

The launcher supplies a verified Identity Principal to protected requests, plus an authorized tenant/PersonGID context to Integration. Restricted Gmail Pub/Sub push is verified through Integration's workload verifier. Unknown restricted calls are denied. Privileged internal gateway dispatch denies requests until a typed authorizer is configured; direct existing public service clients retain their own service invariants. OAuth sign-in for accounts is a future capability and is not represented as working by the strategy interface alone.

## Rate limit and operation contracts

Controller method `@RouteLimit N` metadata is validated by the shared generator and emitted as exact HTTP method+path policy. The `Identity` module's bootstrap returns generated rate limits to Launcher; no Identity route strings are maintained inside the generic limiter. Individual domain methods remain `services/authentication/operations/<action>` and `services/account/operations/<action>` contracts with `impl/` packages; `services/registration/registration.go` and `security/security.go` are their own upfront interfaces. HTTP cookies are exposed as `transport/cookies.Manager` with `transport/cookies/impl` owned by the transport boundary.

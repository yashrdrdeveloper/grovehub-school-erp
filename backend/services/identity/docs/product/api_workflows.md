# Identity — API surface

All success/error serialization uses the shared `libs/core/http/response` writer. Identity must not invent a second envelope. The fields of the current platform response writer take precedence over earlier chat mockups.

## Authentication

| Method | Route | Body |
|---|---|---|
| POST | /identity/public/auth/login | {"identifier":"john@example.test","password":"..."} |
| POST | /identity/public/auth/refresh | none, refresh cookie + CSRF |
| POST | /identity/public/auth/forgot-password | {"email":"john@example.test"} |
| GET | /identity/public/auth/recovery/validate | no body; Authorization: Recovery <proof> (non-consuming) |
| GET | /identity/public/auth/invitations/validate | no body; Authorization: Invitation <proof> (non-consuming) |
| POST | /identity/public/auth/recovery/exchange | none; Authorization: Recovery <proof> |
| GET | /identity/public/auth/reset-password/status | no body; valid HttpOnly reset grant cookie |
| POST | /identity/public/auth/reset-password | {"newPassword":"..."} + grant cookie |
| POST | /identity/public/auth/invitations/exchange | none; Authorization: Invitation <proof> |
| POST | /identity/protected/auth/logout | none |
| POST | /identity/protected/auth/logout-all | none; recent authentication required |
| POST | /identity/protected/auth/change-password | {"newPassword":"..."}; recent authentication required |
| POST | /identity/protected/auth/reauthenticate | {"password":"..."} |
| GET | /identity/protected/auth/sessions | none |
| DELETE | /identity/protected/auth/sessions/{sessionId} | none |
| GET | /identity/protected/auth/me | none |

## Accounts

| Method | Route | Body |
|---|---|---|
| POST | /identity/protected/accounts | {"person":{"firstName":"John","lastName":"Doe"},"email":"john@example.test"} + Idempotency-Key |
| POST | /identity/protected/accounts/search | coresearch.Request JSON |
| GET | /identity/protected/accounts/{accountId} | none |
| PATCH | /identity/protected/accounts/{accountId} | {"username":"john.doe"} |
| PATCH | /identity/protected/accounts/{accountId}/suspend | none |
| PATCH | /identity/protected/accounts/{accountId}/activate | none |
| PATCH | /identity/protected/accounts/{accountId}/deactivate | none |
| POST | /identity/protected/accounts/{accountId}/resend-invitation | none |
| GET | /identity/public/accounts/setup/status | no body; valid HttpOnly setup grant cookie |
| POST | /identity/public/accounts/setup | {"password":"..."} + setup grant cookie |

Admin creation requires a verified ACTIVE session, admin allowlist, and a tenant-scoped account. It resolves Foundation PersonGID, creates PENDING_SETUP, and queues an invitation. No caller-supplied PersonGID, status, role or password is accepted on creation.

## Cookies and exchanges

Access/refresh cookies are HttpOnly, Secure and SameSite-restricted. Public login sets a readable CSRF cookie; the frontend echoes it as `X-CSRF-Token` for authenticated state changes and refresh. Setup/recovery proof is placed in the fragment of an HTTPS email link; frontend exchanges it using a dedicated Authorization header, and the server sets a short-lived HttpOnly grant cookie. No emailed token is sent in the password request JSON.

Public forgot-password accepts all syntactically plausible emails equally to prevent account enumeration. Errors from authenticated administration can include typed safe metadata. This is a source-level implementation guide, not an observed live API test record.

## Invitation activation (two server-enforced submissions)

The public setup flow removes the proof fragment from browser history before issuing a GET validation request. Validation inspects token hash, purpose, expiration, used-at and account status **without consuming** the proof. The existing POST exchange consumes the proof and issues an HttpOnly setup grant cookie. Setup uses only `{"password":"..."}`; the administrator owns any username through protected account onboarding/update, and a missing username remains null.

A database migration adds `attempts` to `identity.password_tokens`. `ReserveSetupAttempt` atomically increments only for valid unconsumed `SETUP_GRANT` tokens with fewer than two attempts, and commits independently of the activation unit of work. A failed transaction cannot refund an attempt. On or after the third submission, the API responds `IDENTITY_SETUP_ATTEMPTS_EXCEEDED` (429) and refuses activation; the administrator must resend the invitation to create a new grant. Password input/confirmation validation in the browser does not call this endpoint and does not use a server attempt. Two legitimate submit requests are permitted, including the successful second one.

A validation GET only checks the proof at that instant; the token can expire or be consumed before POST exchange. The frontend must handle either failure and never display the raw token. Rate-limit controls remain necessary for public endpoints.

After a successful proof exchange, GET grant-status endpoints can resume a password form after page refresh without restoring or persisting the emailed proof in a readable URL or browser storage. The setup status response returns `{verified:true,expiresIn:N}`; exhausted setup grants return `IDENTITY_SETUP_ATTEMPTS_EXCEEDED`.

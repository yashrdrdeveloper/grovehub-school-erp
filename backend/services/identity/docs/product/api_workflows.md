# Identity — API surface

All success/error serialization uses the shared `libs/core/http/response` writer. Identity must not invent a second envelope. The fields of the current platform response writer take precedence over earlier chat mockups.

## Authentication

| Method | Route | Body |
|---|---|---|
| POST | /identity/public/auth/login | {"identifier":"john@example.test","password":"..."} |
| POST | /identity/public/auth/refresh | none, refresh cookie + CSRF |
| POST | /identity/public/auth/forgot-password | {"email":"john@example.test"} |
| POST | /identity/public/auth/recovery/exchange | none; Authorization: Recovery <proof> |
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
| POST | /identity/public/accounts/setup | {"username":"john.doe","password":"..."} + setup grant cookie |

Admin creation requires a verified ACTIVE session, admin allowlist, and a tenant-scoped account. It resolves Foundation PersonGID, creates PENDING_SETUP, and queues an invitation. No caller-supplied PersonGID, status, role or password is accepted on creation.

## Cookies and exchanges

Access/refresh cookies are HttpOnly, Secure and SameSite-restricted. Public login sets a readable CSRF cookie; the frontend echoes it as `X-CSRF-Token` for authenticated state changes and refresh. Setup/recovery proof is placed in the fragment of an HTTPS email link; frontend exchanges it using a dedicated Authorization header, and the server sets a short-lived HttpOnly grant cookie. No emailed token is sent in the password request JSON.

Public forgot-password accepts all syntactically plausible emails equally to prevent account enumeration. Errors from authenticated administration can include typed safe metadata. This is a source-level implementation guide, not an observed live API test record.

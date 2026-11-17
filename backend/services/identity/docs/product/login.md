# Public Login — API Workflow

## Endpoint

`POST /identity/public/auth/login` uses email/password authentication. The operation requires an ACTIVE account with a password credential. It verifies credentials before issuing a persisted session.

Request:

```bash
curl -i -X POST "$BASE_URL/identity/public/auth/login" \
  -H 'Content-Type: application/json' \
  --data '{"email":"person@example.test","password":"REPLACE_WITH_PASSWORD"}'
```

The JSON payload is `{"email":string,"password":string}`. Success is HTTP 200 with the normal REST envelope containing:

```json
{
  "accessToken": "<opaque-bearer-secret>",
  "refreshToken": "<opaque-refresh-secret>",
  "tokenType": "Bearer",
  "expiresIn": 900
}
```

Only password hashes and token digests are stored in the database; raw credentials must not be logged or included in telemetry. Access tokens expire after 15 minutes. Refresh-token expiry is 30 days, and the refresh endpoint rotates both tokens atomically. `POST /identity/public/auth/refresh` accepts `{"refreshToken":"<opaque-refresh-secret>"}`. `POST /identity/protected/auth/logout` requires `Authorization: Bearer <access-token>` and revokes the current session.

A wrong email, password, or non-ACTIVE account returns a generic invalid-credentials error to avoid disclosing account existence. Protected endpoints validate access-token hashes and the current account status. A deployment needs DB migrations and explicit rate limiting, anti-enumeration measures, lockout and rate limiting and comprehensive refresh/revocation integration tests before production release.

## Admin-only endpoint

An admin with a verified access token can attempt `POST /identity/protected/auth/register`. Authorization is additionally checked against `IDENTITY_ADMIN_ACCOUNT_IDS`. No configured admins means registration is forbidden. This allowlist is an interim bootstrap policy, not a replacement for RBAC.

## Test cases

Successful valid login, unknown email, wrong password, pending/suspended/deactivated account, expired token, malformed bearer token, revoked session, duplicate concurrent logins, DB failure and session persistence failure. Compare actual response envelope and error codes against the running API before treating examples as end-to-end verified.

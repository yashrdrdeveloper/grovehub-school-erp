# Admin Invitation and Password Setup

## Protected registration

`POST /identity/protected/auth/register` requires a verified ACTIVE Identity access session and membership in the configured admin account allowlist.

```bash
curl -i -X POST "$BASE_URL/identity/protected/auth/register" \
  -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
  -H "Idempotency-Key: $REGISTRATION_KEY" \
  -H "Content-Type: application/json" \
  --data '{"firstName":"Sample","middleName":"Q","lastName":"Example","email":"sample@example.test","phoneNumber":""}'
```

**No password field is allowed by the intended registration contract.** The account is persisted with status `PENDING_SETUP` and a single-use invitation is queued in the same local Identity transaction. On replay, the same idempotency key and normalized payload return the original account. The sender must not email a password.

## Email link

The configured system mailbox uses Integration `mail.send` to deliver an HTTPS password-setup URL. The URL includes a 256-bit, random, base64url token valid for 30 minutes. Identity stores only SHA-256 digests; the pending delivery job stores the URL under AES-GCM encryption using a configured 256-bit key.

Sending is asynchronous with leases and retry backoff. A provider accepting a send does not guarantee inbox delivery. Retries may produce duplicate emails containing the same token. Resend/reissue after token expiry is not implemented yet.

## Public password activation

`POST /identity/public/auth/set-password`:

```bash
curl -i -X POST "$BASE_URL/identity/public/auth/set-password" \
  -H "Content-Type: application/json" \
  --data '{"token":"REPLACE_WITH_EMAIL_TOKEN","password":"REPLACE_WITH_VALID_PASSWORD"}'
```

Expected successful data includes `{"status":"ACTIVE"}`. This operation atomically consumes the token, creates the password credential and transitions the account to ACTIVE. Expired, reused, invalid or malformed tokens are rejected. Do not put the link or password into logs.

## Critical limitations

- Foundation writes are executed during registration orchestration; independent Foundation transactions are not automatically rolled back when an Identity transaction fails.
- The dedicated Integration mailbox must be actually configured, authorized and able to execute `mail.send`; the code does not create that mailbox or grant permissions.
- Initial-admin bootstrap is external: a working ACTIVE administrator account and credential must exist before protected registration is usable.
- Replay, rate limiting, user activation after delayed email, expired invitation recovery and delivery monitoring need integration testing before release.

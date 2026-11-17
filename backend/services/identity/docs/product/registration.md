# Registration — API Workflow

## User goal

Create an Identity account and password credential associated with a Foundation person, preserving idempotency for retries.

**Evidence:** `controller/auth.go`, `controller/impl/auth.go`, `dto/request/auth.go`, `dto/response/auth.go`, `services/impl/auth.go`.

## Endpoint and access

- Controller annotation: `@Routes /auth`, method annotation: `@POST /register`.
- Controller-relative path: **`POST /auth/register`**. The externally reachable path may include an enclosing router prefix; confirm the deployed prefix before using this as a complete URL.
- Expected successful transport status: **201 Created**, based on `rest.Created(...)`.
- Registration's authentication and rate-limit configuration must be verified against runtime route middleware.

## Request

The request body shape in source is:

```json
{
  "firstName": "Sample",
  "middleName": "Q",
  "lastName": "Example",
  "email": "sample@example.test",
  "phoneNumber": "",
  "password": "example-password"
}
```

The JSON keys are sourced from `request.Register`. This is a **shape example**, not a claim the example password satisfies the deployed password policy.

Registration reads an **Idempotency-Key** from request context. Supply that header through the configured HTTP middleware. Its validation constraints are enforced by the shared idempotency utility.

Generic cURL shape (replace base URL, key and password):

```bash
curl -i -X POST "$IDENTITY_BASE_URL/auth/register" \
  -H 'Content-Type: application/json' \
  -H "Idempotency-Key: $IDEMPOTENCY_KEY" \
  --data '{
    "firstName":"Sample",
    "middleName":"Q",
    "lastName":"Example",
    "email":"sample@example.test",
    "phoneNumber":"",
    "password":"REPLACE_WITH_VALID_SECRET"
  }'
```

## Response DTO fields

The result is mapped to a `Register` object with the following JSON keys:

```json
{
  "accountId": "<uuid>",
  "personGid": "<uuid>",
  "gid": "<value>",
  "email": "sample@example.test",
  "phone": "<value-if-present>",
  "status": "ACTIVE"
}
```

The output schema is observed. The example values are placeholders; the `phone` property is omitted when empty. Confirm the **outer REST envelope**, whether `gid` aliases `personGid`, and exact actual status values from the mapper and integration tests before treating this as a wire-format golden example.

## Observed internal steps and side effects

1. Validate idempotency key and normalize/validate registration body.
2. Calculate a request fingerprint.
3. If the registration idempotency key exists, compare stored fingerprint: equal returns the mapped existing account after credential reconciliation; different produces a conflict.
4. Search Foundation for a person matching the email; create one if absent; fail if ambiguous.
5. Check whether the resolved person already has an Identity account.
6. Update the Foundation person with registration data.
7. Persist Identity account linked to person GID and store the password credential.
8. Return registration result.

**Consistency note:** Identity uses a GORM transaction around registration but Foundation changes may be outside that database transaction. Do not assume the entire multi-component workflow rolls back on error.

## Reproduction / verification checklist

- Successful first registration: record returned IDs; verify Identity account and credential persistence.
- Repeat identical request with same idempotency key: confirm stable account association and no duplicate row.
- Repeat changed payload with same key: confirm idempotency conflict and no unintended mutation.
- Repeat person/email with a different key: confirm business conflict semantics.
- Concurrent identical registration: confirm uniqueness, retry behavior and race handling.
- Unavailable Foundation, account write failure and credential write failure: inspect partial state and documented recovery.
- Confirm malformed JSON, missing required fields, invalid password and missing/invalid idempotency key produce stable error codes.

Concrete successful responses, error HTTP codes and test command output **must be captured from a running instance** before this guide is marked executable/end-to-end verified.

## Known constraints

At present this document covers registration only. It does not assert login, token issuance, refresh, logout, OIDC, SAML, role authorization or password reset endpoints exist.

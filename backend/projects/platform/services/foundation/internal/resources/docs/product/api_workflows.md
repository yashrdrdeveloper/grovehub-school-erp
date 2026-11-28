# Foundation — Person API Execution Guide

**Base URL:** set `BASE_URL` to your isolated backend HTTP endpoint (for example `http://127.0.0.1:<port>`). **Preconditions:** backend started, Foundation and Infrastructure migrations applied, PostgreSQL reachable, permitted development principal or IAM policy in place. Request/response envelopes follow the linked OpenAPI.

## 1. Create Person

**Use case:** Create a permanent Person GID, with optional email/phone contacts and primary flags.

```bash
curl -i -X POST "$BASE_URL/foundation/protected/persons" \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: demo-person-create-001" \
  -d '{
    "firstName":"Demo",
    "lastName":"Learner",
    "contacts":{
      "emails":[{"value":"demo.learner@example.test","isPrimary":true}],
      "phones":[]
    }
  }'
```

Expect HTTP **201** for a newly created aggregate and `success: true` with `data.gid`, `status`, `version`, timestamps and contacts. A request matched to an already existing normalized contact may return HTTP **200** with the existing Person; no new GID or event. A duplicate/conflicting contact may return **409**. Capture `data.gid` from the real response for subsequent steps; never fabricate a UUID.

**Retry behavior:** repeat with the same valid key and identical normalized input to verify idempotence. Reuse with different input should trigger the documented conflict policy, not another create.

## 2. Get Person by GID

```bash
curl -i "$BASE_URL/foundation/protected/persons/$GID"
```

Expect **200** with the complete Person and contacts. An invalid GID should yield **400**; an unknown valid GID **404**. The GID is a permanent UUIDv7 identifier.

## 3. Search Persons

```bash
curl -i -X POST "$BASE_URL/foundation/protected/persons/search" \
  -H "Content-Type: application/json" \
  -d '{
    "searchText":"Demo",
    "filters":[],
    "sort":[{"field":"createdAt","direction":"DESC"}],
    "pageInfo":{"rows":25}
  }'
```

Expect **200**, `data.items` and `data.pageInfo` with `rows`, `hasNext`, and optional `nextCursor`. For the next page, repeat with `pageInfo.cursor` equal to the returned **opaque** cursor. Check contact filters and sorting only from the documented allowlist; unknown fields/operators fail with **400**.

## 4. PATCH Person

**Use case:** Rename and/or replace contacts while retaining the same GID.

```bash
curl -i -X PATCH "$BASE_URL/foundation/protected/persons/$GID" \
  -H "Content-Type: application/json" \
  -d '{"middleName":"Test"}'
```

Expected **200** with updated complete Person. Omitted fields are unchanged. To clear middle name, send `{"middleName":null}`. To clear the full contacts set, send `{"contacts":null}` (use carefully in isolated test data); replacing contacts with `{"contacts":{"emails":[...],"phones":[...]}}` replaces the desired set. Do **not** send `isVerified` in user input. Invalid body **400**, unknown GID **404**, contact conflict **409**.

Verify unchanged PATCH does not increment business version or emit a second event when actual state is identical.

## 5. Deactivate Person

```bash
curl -i -X PATCH "$BASE_URL/foundation/protected/persons/$GID/deactivate"
```

Expect **200**, `data.status = "DEACTIVATED"`. The GID is preserved; calling again returns the already-deactivated Person without a new mutation. Unknown GID **404**, invalid GID **400**. There is no public delete endpoint.

## 6. Business acceptance checks

1. Create once; capture the server-returned GID.
2. Fetch and confirm contacts/primary flags.
3. Search for the new Person; paginate with returned cursor only.
4. PATCH optional name and contacts, then verify read/search and no unauthorized `isVerified` write.
5. Deactivate, fetch and verify stable GID/status.
6. Confirm duplicate create/no-op/deactivate do not emit duplicate changed events.
7. Confirm unauthorized callers are denied **by actual middleware/policy** before any production use; route labels alone are insufficient.

All examples use synthetic data; actual response and HTTP status details are governed by `internal/resources/api/openapi/v1/openapi.yaml`.

# Foundation — Product and API Workflows

Product documentation describes **verified business behavior**, not generic backend design. Link engineering questions to [module engineering docs](../engineering/README.md).

## Per-API reproducibility record

For **each implemented endpoint**, document:

| Field | Required content |
|---|---|
| Capability | Actor's goal and business outcome |
| Preconditions | Required account, permissions, tenant, records and feature flags |
| Contract | Method, path, API version, request/response schema source |
| Setup | Concrete local fixture or seed command |
| Request | Runnable curl/HTTP example with headers and JSON |
| Success | Observed status, response and persisted side effects |
| Failure | Invalid-input, unauthorized, not-found and conflict examples |
| Retry | Idempotency key, deduplication and concurrency behavior |
| Events | Actual published event types and eventually visible effects |
| Verification | Steps/queries to prove result and cleanup |

## End-to-end workflow record

Use this structure only for behaviors confirmed by implementation:

```text
Workflow:
Actor and permission:
Initial state:
Step 1 request / expected result:
Step 2 request / expected result:
Eventual effects and observation point:
Failure/rollback behavior:
How to reproduce locally:
How to reset fixtures:
```

## Source of truth

API handlers and OpenAPI contracts define callable paths. Entity/repository/migration code defines persistence effects. Automated tests prove examples. **Do not treat a proposed flow as released.**

**Status:** Documentation structure initialized. API-specific executable examples must be populated and verified from the actual current implementation.

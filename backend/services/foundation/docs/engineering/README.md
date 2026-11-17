# Foundation — Engineering Documentation

This module follows [the global backend standard](../../../../../docs/engineering/backend/backend.md), [search rules](../../../../../docs/engineering/backend/searchDetails.md) and [annotation rules](../../../../../docs/engineering/backend/annotation.md).

This document is the module-specific engineering entrypoint. It is **not** a duplicate of global SRP, DRY, ACID or folder conventions.

## Code-owned responsibilities

- Domain entities and their invariants belong to this module.
- Request handling, validation, application orchestration and persistence must respect the global dependency boundaries.
- Any configured search fields, joins, indexes, cursors and query budgets must be documented against the actual repository implementation.
- Any published events must list their actual type/version, delivery policy, outbox coupling and consumer compatibility.
- DB schema changes must list the actual migration path, constraints, transaction approach and backwards compatibility.

## Implementation inventory to maintain

| Topic | What this module must record |
|---|---|
| Entrypoints | Actual routes, consumers, schedulers and startup wiring |
| Storage | Concrete tables, keys, indexes, migration history |
| Invariants | Domain constraints and lifecycle rules |
| Transactions | Operations that must commit atomically; isolation expectations |
| Error handling | Stable module errors and transport mappings |
| Search | Supported fields/operators, relation strategies and cursor sort keys |
| Events | Producers, consumers, schema versions, retry/idempotency behavior |
| Runtime | Configuration, health, logs, metrics, tracing |
| Tests | Commands, coverage targets, contract and integration dependencies |

**Status:** This file establishes documentation ownership. Its implementation inventory must be filled from reviewed code; it does not assert that a hypothetical feature is implemented.

For user-facing behavior, see [product documentation](../product/README.md).

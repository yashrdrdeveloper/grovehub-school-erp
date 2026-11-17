# Service Documentation Template

Copy the **structure**, not the contents, into each implemented service. Do not fabricate APIs or business behavior.

## Required files

```text
<service>/
  docs/
    engineering/
      README.md       # implementation architecture and decisions
    product/
      README.md       # real features and reproducible business workflows
```

## Engineering README sections

1. Purpose, ownership and boundaries; link to `docs/engineering/backend/backend.md`.
2. Actual package/folder tree and SRP, omitting nonexistent packages.
3. Entrypoints and dependency graph.
4. Domain invariants and lifecycle.
5. Data stores, constraints, transaction isolation and migration ownership.
6. Internal/external integrations and adapters.
7. Published/consumed event contracts, delivery/idempotency and replay.
8. API contract sources and error catalog, without duplicating the product walkthrough.
9. Configuration, secrets, access control, tenant isolation.
10. Reliability, logs, metrics, traces, operational alerts and recovery.
11. Search configuration, pagination, indexes and query budgets.
12. Testing strategy, commands, known limitations and ADRs.
13. Exceptions to global conventions: explicit rationale, owner and review date.

## Product README sections

1. Capability overview: actor, goal, prerequisites and permissions.
2. API inventory: method, real path, contract/version and auth requirements.
3. Per use case: starting state, steps, exact request headers/body, successful response, state changes and side effects.
4. Validation/error cases: status, error code and steps to reproduce.
5. Retries/idempotency, race/conflict outcomes, pagination and concurrency as applicable.
6. Asynchronous behavior: events, eventual consistency, expected observable effects.
7. Complete end-to-end workflows and cleanup/reset steps.
8. How to run/verify locally with real setup commands and fixtures.
9. Known differences across environments, compatibility and migration notes.

**Separation:** Engineering = how the code works. Product = what the service does and how to verify it. A global standard is linked, not pasted into every module.

## Maintenance rule

Every shipped API or event behavior change updates the owning product documentation. Every internal architectural decision updates its engineering documentation. Mark proposed behavior as proposed and never invent expected responses.
